# ADR-0008: Traefik for Reverse Proxy

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Contrail needs a reverse proxy to route external requests to the appropriate workspace and service based on hostname. Options considered:

1. **Traefik**: Docker-native with label-based configuration
2. **nginx**: Traditional, requires config file generation
3. **Caddy**: Simple config, good TLS support
4. **HAProxy**: High performance, more complex configuration

## Decision

Use Traefik with Docker provider, reading labels from containers.

Traefik is configured to:
- Use `contrail-proxy` as the Docker network for routing
- Expose HTTP (port 80) and HTTPS (port 443) entrypoints
- Read routing rules from container labels
- Support TLS with configurable certificate sources (mkcert, custom, or self-signed)

## Consequences

### Positive

- Docker-native: labels on containers define routing rules
- Dynamic routing: no config file changes needed when containers change
- Labels are added via generated overrides, keeping applications workspace-agnostic
- Mature, well-documented, widely used
- Built-in dashboard for debugging

### Negative

- Additional container to run and manage
- Learning curve for Traefik-specific label syntax

### Neutral

- Traefik is managed as a Docker Compose project by Contrail
- Single instance serves all workspaces on the host
- Auto-started by `workspace up` if not running
