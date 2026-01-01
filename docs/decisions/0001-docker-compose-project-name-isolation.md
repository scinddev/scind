# ADR-0001: Docker Compose Project Name Isolation

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Contrail needs to run multiple instances of the same application simultaneously. A developer might run `dev`, `review`, and `control` workspaces, each containing the same set of applications. Without isolation, these instances would conflict with each other—containers, networks, and volumes would collide.

Docker Compose provides a `--project-name` flag (or `name:` in compose files) that creates isolated namespaces for resources. This is the official mechanism for running multiple copies of the same stack.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

For example:
- `dev-app-one`
- `dev-app-two`
- `review-app-one`

## Consequences

### Positive

- Uses Docker's official mechanism for multi-instance support
- Containers, networks, and volumes are automatically namespaced
- No modifications required to the application's own compose files
- Works with any Docker Compose application without vendor lock-in

### Negative

- Project names become longer and more verbose in Docker listings
- Need to track the mapping between workspace/app names and project names

### Neutral

- Naming convention must be documented and consistently applied
- Generated override files set the `name:` field explicitly
