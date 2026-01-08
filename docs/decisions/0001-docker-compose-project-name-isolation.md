# Docker Compose Project Name Isolation

**Status**: Accepted

## Context

Need to run multiple instances of the same application simultaneously.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

## Consequences

This is Docker's official mechanism for running multiple copies of the same stack. It isolates containers, networks, and volumes without requiring modifications to the application.
