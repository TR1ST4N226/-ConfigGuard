# Scenarios

## Français

Les scénarios de test doivent suivre ce cycle :

1. appliquer une configuration conforme ;
2. exécuter ConfigGuard ;
3. introduire une dérive volontaire ;
4. exécuter ConfigGuard à nouveau ;
5. vérifier le code de sortie et le rapport.

Exemples de dérive :

- supprimer une route ;
- modifier une adresse IP ;
- activer Telnet ;
- supprimer un serveur NTP.

## English

Test scenarios should follow this cycle:

1. apply a compliant configuration;
2. run ConfigGuard;
3. introduce intentional drift;
4. run ConfigGuard again;
5. verify the exit code and report.

Examples:

- remove a route;
- change an IP address;
- enable Telnet;
- remove an NTP server.
