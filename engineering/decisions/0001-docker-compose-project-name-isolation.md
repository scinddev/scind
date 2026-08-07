# Docker Compose Project Name Isolation

**Status**: Accepted

## Context

Need to run multiple instances of the same application simultaneously.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

## Consequences

This is Docker's official mechanism for running multiple copies of the same stack. It isolates containers, networks, and volumes without requiring modifications to the application.

**Note**: This derivation implicitly assumes one working copy per project. When the same repository is checked out into multiple linked git worktrees, every worktree derives the same project name and internal network — a collision, not isolation. [ADR-0016](./0016-per-instance-isolation.md) extends this scheme with a per-instance disambiguation token to cover the multi-working-copy case (empty token = the behavior described here).
