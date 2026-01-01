# ADR-0002: Two-Layer Networking Model

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Services in a Contrail workspace need two types of network access:

1. **External access**: Users and external tools need to reach services via hostnames through a reverse proxy
2. **Internal access**: Applications within a workspace need to communicate with each other directly

A single network could serve both purposes, but this would mix concerns and make isolation between workspaces more complex.

## Decision

Implement a two-layer networking model:

- **`contrail-proxy` network**: Host-wide, shared across all workspaces. Connects Traefik to services that need external access.
- **`{workspace}-internal` network**: Per-workspace (e.g., `dev-internal`). Connects all applications for internal communication with stable aliases.

## Consequences

### Positive

- Clear separation of concerns between external and internal access
- Workspace-internal networks provide strong isolation between workspaces
- Same internal aliases (e.g., `app-one-web`) can exist in multiple workspaces without collision
- Public services are explicitly connected to the proxy network; protected services can remain internal-only

### Negative

- Two networks to manage per workspace
- Services needing external access must be connected to both networks

### Neutral

- The proxy network is created once per host and shared
- Workspace networks are created lazily when the workspace is brought up
