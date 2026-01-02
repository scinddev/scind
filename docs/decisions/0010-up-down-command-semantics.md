# ADR-0010: up/down Command Semantics

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Commands could use `start`/`stop` or `up`/`down` terminology. The naming affects user expectations about what the commands do.

## Decision

Use `up` and `down` as primary commands, matching Docker Compose semantics:
- `up`: Build, create networks/volumes, generate overrides, start containers
- `down`: Stop containers, remove containers/networks, optionally remove volumes

## Consequences

### Positive

- Semantic alignment with Docker Compose, which users already know
- `up` conveys "bring the environment into existence" (more than just starting)
- `down` conveys "tear down" rather than just pausing
- Matches the underlying `docker compose up/down` commands Contrail invokes

### Negative

- Users expecting `start`/`stop` may be initially confused
- `down` removes containers by default (unlike `stop` which just pauses)

### Neutral

- Aliases (`start`/`stop`) could be added later if needed

---

## Related Documents

- [Workspace Lifecycle Spec](../specs/workspace-lifecycle.md) - Full specification of up/down behavior
- [CLI Reference](../reference/cli.md) - Command documentation

---

<!-- Migrated from specs/contrail-prd.md:291-303 -->
