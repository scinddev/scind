# ADR-0010: `up`/`down` Command Semantics

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Commands for starting and stopping workspaces could use various terminology: `start`/`stop`, `up`/`down`, `run`/`halt`, etc. The choice affects user expectations about what the commands do.

## Decision

Use `up` and `down` as primary commands, matching Docker Compose semantics:

- **`up`**: Build, create networks/volumes, generate overrides, start containers
- **`down`**: Stop containers, remove containers/networks, optionally remove volumes

## Consequences

### Positive

- Semantic alignment with Docker Compose, which users already know
- `up` conveys "bring the environment into existence" (more than just starting)
- `down` conveys "tear down" rather than just pausing
- Matches the underlying `docker compose up/down` commands Contrail invokes

### Negative

- Some users may expect `start`/`stop` terminology

### Neutral

- `down --volumes` follows Docker Compose pattern for volume removal

---

## Notes

This decision reinforces Contrail's philosophy as a thin layer over Docker Compose rather than an abstraction that hides the underlying tooling.
