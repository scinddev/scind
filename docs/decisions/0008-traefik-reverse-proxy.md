# ADR-0008: Traefik for Reverse Proxy

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Contrail needs a reverse proxy that can dynamically route to containers based on hostname, without requiring config file changes when containers are started or stopped.

## Decision

Use Traefik with Docker provider, reading labels from containers.

## Consequences

### Positive

- Traefik's Docker integration allows dynamic routing without config file changes
- Labels on containers (added via generated overrides) define routing rules
- Traefik automatically discovers containers and their routing configuration
- Industry-standard tool with good documentation

### Negative

- Adds a dependency on Traefik
- Traefik must be running for proxied services to be accessible
- Users unfamiliar with Traefik may need to learn its concepts

### Neutral

- Traefik dashboard available for debugging routing issues

---

## Related Documents

- [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md) - Implements the Traefik-based proxy layer
- [Docker Labels Spec](../specs/docker-labels.md) - Defines the labels used for Traefik routing

---

<!-- Migrated from specs/contrail-prd.md:266-272 -->
