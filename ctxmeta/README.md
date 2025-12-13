# ctxmeta

A tiny, dependency-free module that standardizes correlation metadata in `context.Context`
for Router Architects microservices.

This module is intentionally independent of logging/tracing/transport libraries.

## Supported identifiers
- request_id
- tenant_id
- user_id

## Usage
```go
ctx = ctxmeta.WithRequestID(ctx, "req-123")
ctx = ctxmeta.WithTenantID(ctx, "t-1")
ctx = ctxmeta.WithUserID(ctx, "u-99")

id, ok := ctxmeta.RequestID(ctx)

fields := ctxmeta.Fields(ctx) // map[string]any for logging/metrics
```
