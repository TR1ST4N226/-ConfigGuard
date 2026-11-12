# ConfigGuard validation runbook

This runbook defines a reproducible validation sequence for ConfigGuard. It is
also the recommended checklist for collecting screenshots or terminal captures
for a release demonstration.

## 1. Toolchain

Required for the CLI:

- Go 1.27.3 or newer
- Git

Required only for the integration lab:

- Docker Desktop or Docker Engine
- PowerShell 7 or Windows PowerShell

## 2. Static and unit validation

Run these commands from the repository root:

```powershell
gofmt -l .
go test ./...
go vet ./...
go build -o .\dist\configguard.exe .
```

The formatting command must produce no output. The test, vet, and build
commands must exit successfully.

## 3. CLI behavior

```powershell
go run . version
go run . check --config .\examples\compliant.conf
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
go run . check --config .\examples\actual.conf --report-format json --report-output .\dist\report.json
go run . check --config .\examples\actual.conf --report-format sarif --report-output .\dist\report.sarif
```

Expected evidence:

- the displayed version is `configguard 0.2.2`;
- the compliant check succeeds;
- the drift example reports the changed lines;
- the drift check returns exit code `4`, while the non-compliant rules check
  returns exit code `3`;
- both report files are created and contain structured output.

## 4. Docker/SSH integration test

```powershell
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
```

The script builds the lab, creates temporary SSH credentials, performs
read-only collection, verifies compliant and non-compliant audits, checks the
expected exit codes, and removes the containers, network, and temporary files.

For a visible manual run:

```powershell
docker compose -f .\lab\docker-compose.yml up -d --build
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
docker compose -f .\lab\docker-compose.yml down
```

## 5. Multi-device audit

Use an inventory with at least two devices and run:

```powershell
go run . audit `
  --inventory .\examples\inventory.yaml `
  --all `
  --desired .\examples\desired.conf `
  --output-dir .\dist\audit-output `
  --report-format json `
  --report-output .\dist\audit-reports
```

Verify that each device has a snapshot in `dist\audit-output` and a matching
JSON report in `dist\audit-reports`.

## 6. Evidence checklist

Capture screenshots with the command and its result visible. Recommended
evidence:

1. `go test ./...`, `go vet ./...`, and the successful build.
2. `configguard version` and a compliant check.
3. Drift detection with the changed lines and exit status.
4. JSON/SARIF report files in the output directory.
5. Docker containers starting, the SSH smoke test passing, and cleanup.
6. GitHub Actions showing all CI checks successful.

Do not include private keys, `known_hosts`, tokens, or credentials in screenshots
or uploaded artifacts.
