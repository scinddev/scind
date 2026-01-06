# ADR-0005: Structure vs State Separation

## Status

Accepted

## Context

Configuration could include runtime choices (which branch, which flavor) or only structural definitions.

## Decision

Separate structure (what exists) from state (what's active):

| Aspect | Structure (config files) | State (runtime) |
|--------|--------------------------|-----------------|
| What apps exist | workspace.yaml | - |
| Available flavors | application.yaml | - |
| Active flavor | - | .generated/state.yaml or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

## Consequences

### Positive

- Configuration files describe the system's shape, not its current state
- State changes frequently; structure changes rarely
- Avoids polluting config files with transient information
- Branch management stays with git where it belongs

### Negative

- State must be tracked separately (in `.generated/state.yaml`)
- Users must understand the distinction between structure and state

### Neutral

- State file is gitignored to avoid conflicts

## Related Documents

- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Implements the structural configuration layer

<!-- Migrated from specs/contrail-prd.md:203-222 -->
