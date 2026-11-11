# Changelog

## [0.2.2] - 2026-11-11

### Fixed

- Synchronized the CLI, SARIF reports, and documentation version to `0.2.2`.

## [0.2.0] - 2026-11-10

### Added

- Read-only SSH collection with verified `known_hosts`.
- YAML inventory-driven audits.
- JSON and SARIF reports, including configuration drift findings.
- Docker Compose SSH laboratory for Windows, WSL2, and Linux.
- Automated Docker SSH smoke test.
- CI validation for tests, vet, cross-builds, compliance, and integration.
- Optional extended rules for DNS, SNMP, and SSH hardening.

### Security

- SSH password authentication is not supported.
- Private keys, `known_hosts`, and lab snapshots are excluded from Git.

