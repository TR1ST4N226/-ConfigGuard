# ConfigGuard

[![CI](https://github.com/TR1ST4N226/-ConfigGuard/actions/workflows/ci.yml/badge.svg)](https://github.com/TR1ST4N226/-ConfigGuard/actions/workflows/ci.yml)
[![Release](https://github.com/TR1ST4N226/-ConfigGuard/actions/workflows/release.yml/badge.svg)](https://github.com/TR1ST4N226/-ConfigGuard/actions/workflows/release.yml)

```text
░██████                            ░████ ░██             ░██████                                        ░██ 
 ░██   ░██                          ░██                   ░██   ░██                                       ░██ 
░██         ░███████  ░████████  ░████████ ░██ ░████████ ░██        ░██    ░██  ░██████   ░██░████  ░████████ 
░██        ░██    ░██ ░██    ░██    ░██    ░██░██    ░██ ░██  █████ ░██    ░██       ░██  ░███     ░██    ░██ 
░██        ░██    ░██ ░██    ░██    ░██    ░██░██    ░██ ░██     ██ ░██    ░██  ░███████  ░██      ░██    ░██ 
 ░██   ░██ ░██    ░██ ░██    ░██    ░██    ░██░██   ░███  ░██  ░███ ░██   ░███ ░██   ░██  ░██      ░██   ░███ 
  ░██████   ░███████  ░██    ░██    ░██    ░██ ░█████░██   ░█████░█  ░█████░██  ░█████░██ ░██       ░█████░██ 
                                                     ░██                                                      
                                               ░███████
```

## Français

ConfigGuard détecte les écarts entre la configuration réseau attendue et la configuration actuellement appliquée à un équipement ou un serveur.

### État du projet

Version 0.2.2. La CLI prend en charge la collecte locale et SSH, les inventaires,
les audits, les règles de conformité et les rapports JSON/SARIF en lecture seule.

### Prérequis

- Go 1.27.3 ou version ultérieure
- Docker Desktop ou Docker Engine, uniquement pour le laboratoire SSH
- PowerShell 7 ou Windows PowerShell, pour exécuter le smoke test sous Windows

### Utilisation

```powershell
go run . version
go run . snapshot --output .\examples\actual.conf
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
go run . check --config .\examples\actual.conf
go run . check --config .\examples\actual.conf --rules .\examples\rules.yaml
go run . check --config .\examples\actual.conf --report-format json --report-output .\report.json
go run . check --config .\examples\actual.conf --report-format sarif --report-output .\report.sarif
```

La commande `snapshot` est actuellement disponible sur Linux et collecte l’état des interfaces et des routes en lecture seule.
La collecte est limitée à 30 secondes par défaut ; utilisez `--timeout 60s` pour modifier cette limite.

La commande `check` applique trois contrôles par défaut : SSH obligatoire, Telnet interdit et NTP configuré.
Un fichier YAML personnalisé peut remplacer ces contrôles :

```yaml
rules:
  - id: ssh-required
    description: SSH must be enabled
    severity: HIGH
    contains:
      - ssh enabled
```

Les rapports JSON et SARIF sont disponibles pour les outils CI et GitHub Code Scanning.
La CI du projet exécute également un contrôle de conformité et conserve le rapport SARIF comme artefact téléchargeable.

### Laboratoire Docker/SSH

Le laboratoire reproductible se trouve dans `lab/` et utilise Docker Compose.
Il fournit un équipement Linux simulé avec un serveur SSH, sans image réseau propriétaire :

```powershell
docker compose -f .\lab\docker-compose.yml up -d --build
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
docker compose -f .\lab\docker-compose.yml down
```

Le smoke test vérifie la collecte SSH par clé, un audit conforme, la détection
d'un drift et les codes de sortie `0`, `3` et `4`. Il crée et supprime ses
clés et snapshots temporaires automatiquement.

Des règles étendues optionnelles sont disponibles dans `examples/extended-rules.yaml`
pour contrôler le DNS, SNMP et le durcissement SSH :

```powershell
go run . check `
  --config .\examples\actual.conf `
  --rules .\examples\extended-rules.yaml
```

Voir [CHANGELOG.md](CHANGELOG.md) pour l'historique des versions.

Pour une machine distante, utilisez une clé SSH et un fichier `known_hosts` vérifié :

```powershell
go run . snapshot --source ssh --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\actual.conf
```

Le mot de passe SSH n’est pas accepté par le MVP. Les commandes distantes doivent rester en lecture seule.

Pour Cisco IOS, ajoutez `--platform cisco-ios`. Le profil utilise des commandes `show` en lecture seule :

```powershell
go run . snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\cisco-actual.conf
```

Pour éviter de répéter les paramètres SSH, utilisez un inventaire YAML :

```powershell
go run . snapshot --inventory .\examples\inventory.yaml --device switch-01 --output .\examples\cisco-actual.conf
```

Les chemins relatifs de `key` et `known_hosts` sont résolus depuis le dossier de l'inventaire.
Pour afficher les équipements sans exposer les paramètres SSH :

```powershell
go run . list --inventory .\examples\inventory.yaml
```

Pour exécuter collecte, comparaison et conformité en une commande :

```powershell
go run . audit `
  --inventory .\examples\inventory.yaml `
  --device switch-01 `
  --desired .\examples\desired.conf `
  --output .\examples\actual.conf
```

Pour auditer tous les équipements de l'inventaire :

```powershell
go run . audit `
  --inventory .\examples\inventory.yaml `
  --all `
  --desired .\examples\desired.conf `
  --output-dir .\audit-output
```

Avec un rapport JSON ou SARIF par équipement, ajoutez `--report-format` et
`--report-output` en indiquant un dossier de rapports.

### Codes de sortie

- `0` : configuration conforme
- `1` : erreur de lecture ou d'exécution
- `2` : utilisation invalide de la commande
- `3` : une ou plusieurs règles de conformité ont échoué
- `4` : un écart de configuration a été détecté

### Validation automatique

GitHub Actions exécute automatiquement le formatage, les tests, `go vet` et la compilation Linux/Windows sur chaque push et chaque Pull Request.
Le workflow vérifie également la conformité SARIF et exécute le smoke test
Docker/SSH. Voir [docs/VALIDATION.md](docs/VALIDATION.md) pour le protocole
complet de validation locale et les preuves à conserver.

## English

ConfigGuard detects drift between the expected network configuration and the configuration currently applied to an equipment or server.

### Project status

Version 0.2.2. ConfigGuard supports read-only local and SSH collection, YAML
inventories, single-device and multi-device audits, compliance rules, drift
detection, and JSON/SARIF reports.

### Requirements

- Go 1.27.3 or newer

### Usage

```powershell
go run . version
go run . snapshot --output .\examples\actual.conf
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
go run . check --config .\examples\actual.conf
go run . check --config .\examples\actual.conf --rules .\examples\rules.yaml
go run . check --config .\examples\actual.conf --report-format json --report-output .\report.json
go run . check --config .\examples\actual.conf --report-format sarif --report-output .\report.sarif
```

The `snapshot` command is currently available on Linux and collects interface and route state in read-only mode.
Collection is limited to 30 seconds by default; use `--timeout 60s` to change the limit.

The `check` command applies three default controls: SSH required, Telnet disabled and NTP configured.
A custom YAML file can replace these controls.

JSON and SARIF reports are available for CI tools and GitHub Code Scanning.
The project CI also runs a compliance check and keeps the SARIF report as a downloadable artifact.

### Docker/SSH laboratory

The reproducible laboratory in `lab/` uses Docker Compose and a Linux
container running OpenSSH. It does not require proprietary network images:

```powershell
docker compose -f .\lab\docker-compose.yml up -d --build
powershell -ExecutionPolicy Bypass -File .\lab\smoke-test.ps1
docker compose -f .\lab\docker-compose.yml down
```

The smoke test verifies SSH key-based collection, a compliant audit, drift
detection, and exit codes `0`, `3`, and `4`. Temporary keys and snapshots are
removed automatically.

For a remote machine, use an SSH private key and a verified `known_hosts` file:

```powershell
go run . snapshot --source ssh --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\actual.conf
```

SSH passwords are not accepted by the MVP. Remote commands must remain read-only.

For Cisco IOS, add `--platform cisco-ios`. The profile uses read-only `show` commands:

```powershell
go run . snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\cisco-actual.conf
```

To avoid repeating SSH parameters, use a YAML inventory:

```powershell
go run . snapshot --inventory .\examples\inventory.yaml --device switch-01 --output .\examples\cisco-actual.conf
```

Relative `key` and `known_hosts` paths are resolved from the inventory directory.
To list devices without exposing SSH parameters:

```powershell
go run . list --inventory .\examples\inventory.yaml
```

To run collection, comparison and compliance in one command:

```powershell
go run . audit `
  --inventory .\examples\inventory.yaml `
  --device switch-01 `
  --desired .\examples\desired.conf `
  --output .\examples\actual.conf
```

To audit every device in the inventory:

```powershell
go run . audit `
  --inventory .\examples\inventory.yaml `
  --all `
  --desired .\examples\desired.conf `
  --output-dir .\audit-output
```

For one JSON or SARIF report per device, add `--report-format` and
`--report-output` with a report directory.

### Exit codes

- `0`: configuration is compliant
- `1`: input or execution error
- `2`: invalid command usage
- `3`: one or more compliance rules failed
- `4`: configuration drift was detected

### Automated validation

GitHub Actions automatically runs formatting checks, tests, `go vet`, and Linux/Windows builds on every push and Pull Request.
It also runs the SARIF compliance check and the Docker/SSH smoke test. See
[docs/VALIDATION.md](docs/VALIDATION.md) for the complete local validation
procedure and evidence checklist.
