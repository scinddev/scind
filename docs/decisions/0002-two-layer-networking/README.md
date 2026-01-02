# ADR-0002: Two-Layer Networking

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Services need both external access (via reverse proxy for browsers/API clients) and internal access (between applications within a workspace). A single network cannot satisfy both requirements cleanly.

## Decision

Implement a two-layer networking architecture:

- **`contrail-proxy` network**: Host-wide, connects Traefik to services that need external access
- **`{workspace}-internal` network**: Per-workspace, connects all applications for internal communication

## Consequences

### Positive

- Clear separation of concerns between public and internal traffic
- Public services routable via Traefik while protected services remain internal-only
- Workspace-internal network provides complete isolation between workspaces
- Same internal aliases can exist in multiple workspaces without conflict

### Negative

- Services needing external access must be connected to two networks
- Additional complexity in generated override files

### Neutral

- The proxy network is shared across all workspaces on the host

---

## Notes

This architecture enables the internal alias pattern where `app-one-web` always resolves to the correct application within the current workspace context.
