# ADR-0007: Port Type System for Exported Services

## Status

Accepted

## Context

Services need different handling based on how they're accessed - some need HTTP proxying, others need direct port binding.

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

### Positive

- `type` determines routing: `proxied` (through Traefik) or `assigned` (direct port binding)
- `protocol` specifies how proxied traffic is handled: `http`, `https`, or future SNI types
- Supports multiple protocols on the same exported service (both HTTP and HTTPS)
- Environment variables use proxy values (port 80/443) for proxied services
- Enables future plugin system for additional protocols (postgresql, mysql SNI routing)

### Negative

- More complex configuration than a single port type
- Users must understand the distinction between type and protocol

### Neutral

- `visibility` remains as documentation for collaborators (public vs protected)
- `proxied`: Traffic routed through Traefik; protocol specifies how (http, https, future SNI types)
- `assigned`: Direct port binding; if port unavailable, Contrail finds next available

## Related Documents

- [Port Types Spec](../specs/port-types.md) - Full specification of port types and their behavior

<!-- Migrated from specs/contrail-prd.md:234-264 -->
