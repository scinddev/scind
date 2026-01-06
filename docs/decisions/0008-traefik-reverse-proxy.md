# ADR-0008: Traefik for Reverse Proxy

## Status

Accepted

## Context

Need a reverse proxy that can dynamically route to containers.

## Decision

Use Traefik with Docker provider, reading labels from containers.

## Consequences

### Positive

- Traefik's Docker integration allows dynamic routing without config file changes
- Traefik automatically discovers containers and their routing configuration
- Labels on containers (added via generated overrides) define routing rules
- Industry-standard tool with good documentation

### Negative

- Adds a dependency on Traefik
- Traefik must be running for proxied services to be accessible
- Users unfamiliar with Traefik may need to learn its concepts

### Neutral

- Single shared proxy instance serves all workspaces
- Traefik dashboard available for debugging routing issues

## Related Documents

- [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md) - Implements the Traefik-based proxy layer
- [Docker Labels Spec](../specs/docker-labels.md) - Defines the labels used for Traefik routing

<!-- Migrated from specs/contrail-prd.md:265-271 -->
