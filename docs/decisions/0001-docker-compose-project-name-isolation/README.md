# ADR-0001: Docker Compose Project Name Isolation

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Contrail needs to run multiple instances of the same application simultaneously across different workspaces (e.g., development, code review, and stable/control environments). Without isolation, container names, networks, and volumes would collide.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}` (e.g., `dev-app-one`).

## Consequences

### Positive

- Uses Docker's official mechanism for running multiple copies of the same stack
- Isolates containers, networks, and volumes without requiring modifications to the application
- Predictable naming makes debugging straightforward
- No special conventions needed in application docker-compose.yaml files

### Negative

- Project names become longer (workspace prefix)
- Creative naming could produce collision (e.g., workspace `dev-app` with app `one` vs workspace `dev` with app `app-one` both produce `dev-app-one`)

### Neutral

- Applications remain unaware of the isolation mechanism

---

## Notes

This decision enables the core value proposition of Contrail: running multiple complete copies of multi-application stacks simultaneously.
