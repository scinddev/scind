# ADR-0007: Port Type System for Exported Services

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Services need different handling based on how they're accessed:

- **Web services**: Need HTTP/HTTPS proxying through Traefik with hostname-based routing
- **Databases and caches**: Need direct port binding for client tools that don't support HTTP proxying
- **Debug ports**: Need direct access for debuggers

A one-size-fits-all approach would either force everything through HTTP proxying (breaking database clients) or require manual port management (error-prone).

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

**Port types**:
- `proxied`: Traffic routed through Traefik. Gets hostname and Traefik labels.
- `assigned`: Direct port binding to host. Auto-assigned if port is unavailable.

**Protocols** (for proxied type):
- `http`: Routes through Traefik's web entrypoint (port 80)
- `https`: Routes through Traefik's websecure entrypoint (port 443)
- Future: `tcp`, `postgresql`, `mysql` for SNI-based routing

## Consequences

### Positive

- Clear distinction between proxy-routed and directly-bound services
- Environment variables use appropriate values (proxy port vs assigned port)
- Supports multiple protocols on the same exported service
- Enables future plugin system for additional protocols
- Visibility flag documents intent without changing behavior

### Negative

- More complex schema than a simple port list
- Users must understand the distinction between types

### Neutral

- Each exported service may have at most one `http` and one `https` proxied port
- Assigned ports are tracked globally to prevent conflicts across workspaces
