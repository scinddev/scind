# ADR-0007: Port Type System for Exported Services

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Services need different handling based on how they're accessed. Web frontends need HTTP/HTTPS proxying with hostname routing. Databases need direct port binding for tools that don't support proxied connections. A unified approach cannot handle both cases optimally.

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
- `proxied`: Traffic routed through Traefik with hostname-based routing
- `assigned`: Direct port binding on host, auto-incremented if unavailable

**Protocols** (for proxied type):
- `http`, `https`: HTTP routing through web/websecure entrypoints
- Future: `tcp`, `postgresql`, `mysql` for SNI-based TCP routing

## Consequences

### Positive

- Clear distinction between routing mechanisms
- Supports multiple protocols on the same exported service
- Environment variables use appropriate values (proxy ports for proxied, host ports for assigned)
- Enables future plugin system for additional protocols
- `visibility` flag provides documentation for collaborators

### Negative

- More complex configuration than simple port lists
- Users must understand the distinction between types

### Neutral

- Visibility is documentation only; doesn't affect routing behavior

---

## Notes

This design anticipates future SNI-based TCP routing for databases while keeping the initial implementation focused on HTTP/HTTPS and direct port binding.
