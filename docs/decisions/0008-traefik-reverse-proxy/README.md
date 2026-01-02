# ADR-0008: Traefik for Reverse Proxy

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Contrail needs a reverse proxy that can dynamically route to containers based on hostname. The proxy must support automatic configuration as containers start and stop, without requiring manual configuration file updates.

## Decision

Use Traefik with Docker provider, reading labels from containers for routing configuration.

## Consequences

### Positive

- Traefik's Docker integration allows dynamic routing without config file changes
- Labels on containers (added via generated overrides) define routing rules
- Traefik is well-documented and widely used
- Supports both HTTP/HTTPS routing and future TCP/SNI routing
- Dashboard provides visibility into routing configuration

### Negative

- Adds Traefik as a dependency
- Users unfamiliar with Traefik may need to learn its concepts for debugging
- Traefik must be running before proxied services can be accessed

### Neutral

- Single shared Traefik instance serves all workspaces on the host
- Contrail manages proxy lifecycle (`proxy up`, `proxy down`)

---

## Notes

The Traefik proxy is bootstrapped via `contrail proxy init` and auto-started by `workspace up` if not running. Alternative proxies (nginx, Caddy) were considered but Traefik's native Docker label integration made it the clear choice.
