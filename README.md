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
go run . check --config .\examples\actual.conf
```

La commande `snapshot` est actuellement disponible sur Linux et collecte l’état des interfaces et des routes en lecture seule.

La commande `check` applique trois contrôles par défaut : SSH obligatoire, Telnet interdit et NTP configuré.

Pour une machine distante, utilisez une clé SSH et un fichier `known_hosts` vérifié :

```powershell
go run . snapshot --source ssh --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\actual.conf
```

Le mot de passe SSH n’est pas accepté par le MVP. Les commandes distantes doivent rester en lecture seule.

Pour Cisco IOS, ajoutez `--platform cisco-ios`. Le profil utilise des commandes `show` en lecture seule :

```powershell
go run . snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\cisco-actual.conf
```

### Codes de sortie

- `0` : commande terminée ; la configuration peut être conforme ou contenir un écart
- `1` : erreur de lecture ou d'exécution
- `2` : utilisation invalide de la commande
- `3` : une ou plusieurs règles de conformité ont échoué

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
go run . check --config .\examples\actual.conf
```

The `snapshot` command is currently available on Linux and collects interface and route state in read-only mode.

The `check` command applies three default controls: SSH required, Telnet disabled and NTP configured.

For a remote machine, use an SSH private key and a verified `known_hosts` file:

```powershell
go run . snapshot --source ssh --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\actual.conf
```

SSH passwords are not accepted by the MVP. Remote commands must remain read-only.

For Cisco IOS, add `--platform cisco-ios`. The profile uses read-only `show` commands:

```powershell
go run . snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output .\examples\cisco-actual.conf
```

### Exit codes

- `0`: command completed; the configuration may be compliant or contain drift
- `1`: input or execution error
- `2`: invalid command usage
- `3`: one or more compliance rules failed
