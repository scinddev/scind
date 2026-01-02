# ADR-0002: Two-Layer Networking

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Services need both external access (via reverse proxy) and internal access (between applications in the same workspace).

## Decision

Implement two network layers:
- `contrail-proxy` network: Host-wide, connects Traefik to public services
- `{workspace}-internal` network: Per-workspace, connects all applications for internal communication

## Consequences

### Positive

- Separating concerns allows public services to be routable via Traefik while protected services remain internal
- The workspace-internal network provides isolation between workspaces
- Applications in different workspaces cannot accidentally communicate

### Negative

- More complex network topology than a single flat network
- Debugging network issues requires understanding both layers

### Neutral

- Each workspace has its own internal network named `{workspace}-internal`

---

## Related Documents

- [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md) - Implements the two-layer networking architecture

---

<!-- Migrated from specs/contrail-prd.md:171-179 -->
