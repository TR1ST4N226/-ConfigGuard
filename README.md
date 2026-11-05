# ConfigGuard

ConfigGuard detects drift between an expected network configuration and the configuration currently applied to an equipment or server.

## Status

Early MVP. The current CLI compares two configuration files in read-only mode.

## Requirements

- Go 1.27.3 or newer

## Run

```powershell
go run . version
go run . diff --desired .\examples\desired.conf --actual .\examples\actual.conf
```

Exit codes:

- `0`: command completed; configuration may be compliant or contain drift
- `1`: input or execution error
- `2`: invalid command usage
