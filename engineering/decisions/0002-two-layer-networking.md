# Two-Layer Networking

**Status**: Accepted

## Context

Services need both external access (via reverse proxy) and internal access (between applications).

## Decision

- `scind-proxy` network: Host-wide, connects Traefik to public services
- `{workspace}-internal` network: Per-workspace, connects all applications for internal communication

## Consequences

Separating concerns allows public services to be routable via Traefik while protected services remain internal. The workspace-internal network provides isolation between workspaces.
