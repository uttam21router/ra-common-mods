# ra-common-mods/logging

Structured Logging module (v1) per `STRUCTURED_LOGGING_DESIGN_v1.md`.

Simple package layout:
- `logging` (root): public package + core implementation
- `api/`: stable contracts (types + backend interface)
- `adapters/logrus/`: default backend implementation (imports `api` only)

No bridge/wrapper types are used.

## Quick start
```go
logging.InitService("user-service")
log := logging.Subsystem("repo")

log.Error(ctx, "query failed", logging.Err(err).With(logging.Fields{
  "table": "users",
}))
```
