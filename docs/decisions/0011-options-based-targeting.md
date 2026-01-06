# ADR-0011: Options-Based Targeting with Context Detection

## Status

Accepted

## Context

Commands need to target specific workspaces and applications.

## Decision

Use `--workspace` and `--app` options (not positional arguments) with automatic context detection from current directory.

## Consequences

### Positive

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options

### Negative

- Slightly more verbose than positional arguments
- Context detection behavior must be documented and understood

### Neutral

- Global flags are always available: `-w/--workspace`, `-a/--app`

## Related Documents

- [Context Detection Spec](../specs/context-detection.md) - Full specification of context detection behavior
- [CLI Reference](../reference/cli.md) - Command documentation

<!-- Migrated from specs/contrail-prd.md:304-315 -->
