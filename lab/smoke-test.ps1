$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent $PSScriptRoot
$composeFile = Join-Path $PSScriptRoot "docker-compose.yml"
$tempDirectory = Join-Path ([IO.Path]::GetTempPath()) ("configguard-smoke-" + [Guid]::NewGuid().ToString("N"))
$isWindowsHost = $env:OS -eq "Windows_NT"
$composeStarted = $false

function Invoke-ExpectedExitCode {
    param(
        [string]$Command,
        [string[]]$Arguments,
        [int]$ExpectedExitCode = 0
    )

    & $Command @Arguments
    if ($LASTEXITCODE -ne $ExpectedExitCode) {
        throw "$Command returned $LASTEXITCODE; expected $ExpectedExitCode"
    }
}

try {
    New-Item -ItemType Directory -Path $tempDirectory | Out-Null
    $privateKey = Join-Path $tempDirectory "id_ed25519"
    $publicKey = "$privateKey.pub"
    $authorizedKeys = Join-Path $tempDirectory "authorized_keys"
    $knownHosts = Join-Path $tempDirectory "known_hosts"
    $binary = Join-Path $tempDirectory "configguard.exe"
    $inventory = Join-Path $tempDirectory "inventory.yaml"
    $driftInventory = Join-Path $tempDirectory "drift-inventory.yaml"

    if ($isWindowsHost) {
        $keygenCommand = 'ssh-keygen -t ed25519 -f "' + $privateKey + '" -N "" -C configguard-smoke'
        & cmd.exe /c $keygenCommand
    } else {
        $keygenCommand = "ssh-keygen -t ed25519 -f '$privateKey' -N '' -C configguard-smoke"
        & bash -c $keygenCommand
    }
    if ($LASTEXITCODE -ne 0) {
        throw "ssh-keygen returned $LASTEXITCODE"
    }
    Copy-Item $publicKey $authorizedKeys

    $env:CONFIGGUARD_SSH_DIR = $tempDirectory
    Invoke-ExpectedExitCode "docker" @("compose", "-f", $composeFile, "up", "-d", "--build", "--force-recreate")
    $composeStarted = $true

    $hostKeys = & docker exec configguard-network-device sh -c "cat /etc/ssh/ssh_host_*_key.pub"
    $knownHostLines = foreach ($hostKey in $hostKeys) {
        $parts = $hostKey -split "\s+"
        if ($parts.Length -ge 2) {
            "[127.0.0.1]:2222 $($parts[0]) $($parts[1])"
        }
    }
    if (-not $knownHostLines) {
        throw "No SSH host keys were returned by the network-device container"
    }
    $knownHostLines | Set-Content -Encoding ascii $knownHosts

    Push-Location $repoRoot
    try {
        Invoke-ExpectedExitCode "go" @("build", "-o", $binary, ".")
        @"
devices:
  - name: docker-device
    source: ssh
    platform: linux
    ssh:
      host: 127.0.0.1:2222
      user: root
      key: '$privateKey'
      known_hosts: '$knownHosts'
      command: >-
        printf 'hostname switch-01\nntp server 10.0.0.10\nssh enabled\n'
"@ | Set-Content -Encoding utf8 $inventory

        Invoke-ExpectedExitCode $binary @(
            "audit",
            "--inventory=$inventory",
            "--device=docker-device",
            "--desired=.\examples\desired.conf",
            "--output=$tempDirectory\audit.conf",
            "--report-format=json",
            "--report-output=$tempDirectory\audit.json"
        )

        @"
devices:
  - name: docker-device
    source: ssh
    platform: linux
    ssh:
      host: 127.0.0.1:2222
      user: root
      key: '$privateKey'
      known_hosts: '$knownHosts'
      command: >-
        printf 'hostname switch-01\nntp server 10.0.0.20\ntelnet enabled\n'
"@ | Set-Content -Encoding utf8 $driftInventory

        Invoke-ExpectedExitCode $binary @(
            "audit",
            "--inventory=$driftInventory",
            "--device=docker-device",
            "--desired=.\examples\desired.conf",
            "--output=$tempDirectory\audit-drift.conf",
            "--report-format=json",
            "--report-output=$tempDirectory\audit-drift.json"
        ) 4

        $commonArguments = @(
            "snapshot",
            "--source=ssh",
            "--ssh-host=127.0.0.1:2222",
            "--ssh-user=root",
            "--ssh-key=$privateKey",
            "--ssh-known-hosts=$knownHosts"
        )

        Invoke-ExpectedExitCode $binary ($commonArguments + @(
            "--ssh-command=printf 'hostname switch-01\nntp server 10.0.0.10\nssh enabled\n'",
            "--output=$tempDirectory\compliant.conf"
        ))
        Invoke-ExpectedExitCode $binary @(
            "diff",
            "--desired=.\examples\desired.conf",
            "--actual=$tempDirectory\compliant.conf"
        )
        Invoke-ExpectedExitCode $binary @(
            "check",
            "--config=$tempDirectory\compliant.conf"
        )

        Invoke-ExpectedExitCode $binary ($commonArguments + @(
            "--ssh-command=printf 'hostname switch-01\nntp server 10.0.0.20\ntelnet enabled\n'",
            "--output=$tempDirectory\drift.conf"
        ))
        Invoke-ExpectedExitCode $binary @(
            "diff",
            "--desired=.\examples\desired.conf",
            "--actual=$tempDirectory\drift.conf"
        ) 4
        Invoke-ExpectedExitCode $binary @(
            "check",
            "--config=$tempDirectory\drift.conf"
        ) 3
    }
    finally {
        Pop-Location
    }

    Write-Host "ConfigGuard SSH smoke test passed."
}
finally {
    if ($composeStarted) {
        & docker compose -f $composeFile down
    }
    Remove-Item Env:CONFIGGUARD_SSH_DIR -ErrorAction SilentlyContinue
    Remove-Item $tempDirectory -Recurse -Force -ErrorAction SilentlyContinue
}
