# ADR-0005: Structure vs State Separation

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Configuration files could include:

1. **Only structure**: What applications exist, what flavors are available, what services are exported
2. **Structure and state**: Also includes runtime choices like active flavor, current branch, running status

Mixing structure and state leads to configuration files that change frequently with transient information.

## Decision

Separate structure (what exists) from state (what's active):

| Aspect | Structure (config files) | State (runtime) |
|--------|--------------------------|-----------------|
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

## Consequences

### Positive

- Configuration files describe the system's shape, not its current state
- Version-controlled files remain stable (no noisy commits)
- State changes don't pollute config history
- Branch management stays with git where it belongs
- Clear separation of concerns

### Negative

- State must be tracked separately (in `.generated/state.yaml`)
- Some queries require checking both config and state

### Neutral

- State files are gitignored
- CLI commands can override state temporarily without persisting
