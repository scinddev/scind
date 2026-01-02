# ADR-0001: Docker Compose Project Name Isolation

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Contrail needs to run multiple instances of the same application simultaneously. Each instance must have isolated containers, networks, and volumes.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

## Consequences

### Positive

- Uses Docker's official mechanism for running multiple copies of the same stack
- Isolates containers, networks, and volumes without requiring modifications to the application
- Works with any existing Docker Compose project

### Negative

- Project names must be unique across all workspaces on a host
- Creative naming that produces identical project names could cause conflicts

### Neutral

- Naming follows convention: `{workspace}-{app}` (e.g., `dev-app-one`)

---

## Related Documents

- [Naming Conventions Spec](../specs/naming-conventions.md) - Implements the naming patterns derived from this decision

---

## Notes

Example collision to avoid: workspace `dev-app` with app `one` and workspace `dev` with app `app-one` both produce project name `dev-app-one`.

<!-- Migrated from specs/contrail-prd.md:162-169 -->
