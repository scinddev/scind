# Options-Based Targeting with Context Detection

**Status**: Accepted

## Context

Commands need to target specific workspaces and applications.

## Decision

Use `--workspace` and `--app` options (not positional arguments) with automatic context detection from current directory.

## Consequences

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options
