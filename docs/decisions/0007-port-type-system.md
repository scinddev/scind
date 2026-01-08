# Port Type System for Exported Services

**Status**: Accepted

## Context

Services need different handling based on how they're accessed—some need HTTP proxying, others need direct port binding.

## Decision

Each exported service declares ports with a `type` (routing mechanism) and optionally a `protocol`:

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
      - type: proxied
        protocol: http
        visibility: protected
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

## Consequences

- `type` determines routing: `proxied` (through Traefik) or `assigned` (direct port binding)
- `protocol` specifies how proxied traffic is handled: `http`, `https`, or future SNI types
- Supports multiple protocols on the same exported service (both HTTP and HTTPS)
- Environment variables use proxy values (port 80/443) for proxied services
- Enables future plugin system for additional protocols (postgresql, mysql SNI routing)
- `visibility` remains as documentation for collaborators

## Related Decisions

- [ADR-0008: Traefik for Reverse Proxy](0008-traefik-reverse-proxy.md) - Traefik handles proxied port routing
- [ADR-0009: Flexible TLS Configuration](0009-flexible-tls-configuration.md) - TLS for proxied ports
