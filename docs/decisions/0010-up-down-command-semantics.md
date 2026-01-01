# ADR-0010: `up`/`down` Command Semantics

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

CLI commands for starting and stopping workspaces could use various naming conventions:

1. **`start`/`stop`**: Common but implies simple state toggle
2. **`up`/`down`**: Docker Compose convention, implies full lifecycle management
3. **`create`/`destroy`**: Implies persistence changes
4. **`run`/`halt`**: Vagrant convention

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

- None significant—this is the expected convention for Docker Compose tooling

### Neutral

- Top-level aliases (`contrail up`, `contrail down`) delegate to workspace commands
- Both commands support `--app` flag for targeting specific applications
