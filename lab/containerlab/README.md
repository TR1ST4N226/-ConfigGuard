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
