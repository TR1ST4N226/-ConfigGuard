# ConfigGuard network lab

## Français

### Windows + WSL2 + Docker Desktop

Sous Windows, utilisez Docker Compose. Il évite les limitations de namespaces de Containerlab avec Docker Desktop :

```bash
cd lab
docker compose up -d --build
docker compose ps
```

Le conteneur `configguard` utilise la même image locale que le nœud réseau et monte le dépôt dans `/workspace`.

Arrêt :

```bash
docker compose down
```

Ce laboratoire fournit une topologie Linux minimale pour tester la collecte et la détection de drift sans équipement propriétaire.

### Test SSH avec Docker Compose

Le service `network-device` expose un serveur SSH de test sur `localhost:2222`.
Cette procédure utilise une clé locale non versionnée :

```powershell
cd C:\TR1ST4N_ROG_ZEPHYRUS\ConfigGuard
cmd.exe /c 'ssh-keygen -t ed25519 -f C:\TR1ST4N_ROG_ZEPHYRUS\ConfigGuard\lab\ssh\id_ed25519 -N "" -C configguard-lab'
Copy-Item .\lab\ssh\id_ed25519.pub .\lab\ssh\authorized_keys
docker compose -f .\lab\docker-compose.yml up -d --build --force-recreate
```

Construisez ensuite le fichier `known_hosts` avec les clés du conteneur :

```powershell
$hostKeys = docker exec configguard-network-device sh -c "cat /etc/ssh/ssh_host_*_key.pub"
$knownHosts = foreach ($hostKey in $hostKeys) {
  $parts = $hostKey -split '\s+'
  if ($parts.Length -ge 2) {
    "[127.0.0.1]:2222 $($parts[0]) $($parts[1])"
  }
}
$knownHosts | Set-Content -Encoding ascii .\lab\ssh\known_hosts
```

Collecte SSH en lecture seule :

```powershell
go run . snapshot `
  --source=ssh `
  --ssh-host=127.0.0.1:2222 `
  --ssh-user=root `
  --ssh-key=.\lab\ssh\id_ed25519 `
  --ssh-known-hosts=.\lab\ssh\known_hosts `
  --output=.\lab\ssh\snapshot.conf
```

Les fichiers de clé, `known_hosts` et les snapshots du laboratoire sont ignorés par Git.

Audit via l'inventaire :

```powershell
go run . audit `
  --inventory .\examples\lab-inventory.yaml `
  --device docker-device `
  --desired .\examples\desired.conf `
  --output .\lab\ssh\audit.conf `
  --report-format json `
  --report-output .\lab\ssh\audit.json
```

Smoke test automatisé :

```powershell
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
```

Prérequis :

- Linux ou WSL2 ;
- Docker ;
- Containerlab.

Construire l'image locale :

```bash
docker build -t configguard-lab:local lab/containerlab
```

Démarrage :

```bash
sudo containerlab deploy -t topology.yml
```

Le laboratoire est volontairement limité. Il utilise l'adresse fournie par le réseau Docker partagé et ne simule pas encore un équipement Cisco ou Juniper.

Sous WSL2, les avertissements concernant LLDP, iptables et les ACL du dossier peuvent apparaître avec Docker Desktop. Ils ne sont pas bloquants si les conteneurs restent actifs.

Suppression :

```bash
sudo containerlab destroy -t topology.yml
```

## English

### Windows + WSL2 + Docker Desktop

On Windows, use Docker Compose. It avoids namespace limitations affecting Containerlab with Docker Desktop:

```bash
cd lab
docker compose up -d --build
docker compose ps
```

The `configguard` container uses the same local image as the network node and mounts the repository at `/workspace`.

Stop:

```bash
docker compose down
```

This lab provides a minimal Linux topology for testing collection and drift detection without proprietary equipment.

### SSH test with Docker Compose

The `network-device` service exposes a test SSH server on `localhost:2222`.
Use a local key that is not committed:

```powershell
cd C:\TR1ST4N_ROG_ZEPHYRUS\ConfigGuard
cmd.exe /c 'ssh-keygen -t ed25519 -f C:\TR1ST4N_ROG_ZEPHYRUS\ConfigGuard\lab\ssh\id_ed25519 -N "" -C configguard-lab'
Copy-Item .\lab\ssh\id_ed25519.pub .\lab\ssh\authorized_keys
docker compose -f .\lab\docker-compose.yml up -d --build --force-recreate
```

Then build `known_hosts` from the container host keys using the equivalent PowerShell procedure above. Run the read-only SSH snapshot with the `go run . snapshot` command shown above.

Run a complete inventory-driven audit:

```powershell
go run . audit `
  --inventory .\examples\lab-inventory.yaml `
  --device docker-device `
  --desired .\examples\desired.conf `
  --output .\lab\ssh\audit.conf `
  --report-format json `
  --report-output .\lab\ssh\audit.json
```

Keys, `known_hosts`, and lab snapshots are ignored by Git.

Automated smoke test:

```powershell
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
```

Requirements:

- Linux or WSL2;
- Docker;
- Containerlab.

Build the local image:

```bash
docker build -t configguard-lab:local lab/containerlab
```

Start:

```bash
sudo containerlab deploy -t topology.yml
```

The lab is intentionally small. It uses the address provided by the shared Docker network and does not simulate Cisco or Juniper yet.

Under WSL2, Docker Desktop may report LLDP, iptables, and directory ACL warnings. They are not blocking if the containers remain running.

Destroy:

```bash
sudo containerlab destroy -t topology.yml
```
