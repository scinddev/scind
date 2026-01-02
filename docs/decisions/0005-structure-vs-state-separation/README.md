# ADR-0005: Structure vs State Separation

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Configuration files could include both structural definitions (what exists) and runtime choices (which branch, which flavor). Mixing these concerns leads to configuration files that change frequently and contain transient information.

## Decision

Separate structure from state:

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
- Cleaner diffs when reviewing configuration changes

### Negative

- State must be tracked separately in `.generated/state.yaml`
- Requires understanding which file holds which information

### Neutral

- Git working directory is the source of truth for branch state

---

## Notes

This separation enables workspace.yaml to be committed to version control without constant changes from flavor switching or branch changes.
