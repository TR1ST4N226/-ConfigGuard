# ConfigGuard

## Français

ConfigGuard détecte les écarts entre la configuration réseau attendue et la configuration actuellement appliquée à un équipement ou un serveur.

### État du projet

MVP initial. La CLI compare actuellement deux fichiers de configuration en lecture seule.

### Prérequis

- Go 1.27.3 ou version ultérieure

### Utilisation

```powershell
go run . version
go run . snapshot --output .\examples\actual.conf
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
```

La commande `snapshot` est actuellement disponible sur Linux et collecte l’état des interfaces et des routes en lecture seule.

### Codes de sortie

- `0` : commande terminée ; la configuration peut être conforme ou contenir un écart
- `1` : erreur de lecture ou d'exécution
- `2` : utilisation invalide de la commande

## English

ConfigGuard detects drift between the expected network configuration and the configuration currently applied to an equipment or server.

### Project status

Early MVP. The CLI currently compares two configuration files in read-only mode.

### Requirements

- Go 1.27.3 or newer

### Usage

```powershell
go run . version
go run . snapshot --output .\examples\actual.conf
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
```

The `snapshot` command is currently available on Linux and collects interface and route state in read-only mode.

### Exit codes

- `0`: command completed; the configuration may be compliant or contain drift
- `1`: input or execution error
- `2`: invalid command usage
