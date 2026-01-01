# ADR-0011: Options-Based Targeting with Context Detection

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Commands need to target specific workspaces and applications. Approaches considered:

1. **Positional arguments**: `contrail workspace up dev app-one`
2. **Option flags**: `contrail workspace up --workspace=dev --app=app-one`
3. **Context detection**: Infer from current directory

Positional arguments are concise but become confusing with optional values. Flags are verbose but explicit and composable.

## Decision

Use `--workspace` / `-w` and `--app` / `-a` options (not positional arguments) with automatic context detection from current directory.

**Context detection algorithm**:
1. Walk up from current directory looking for `workspace.yaml` (establishes workspace root)
2. Walk up toward workspace root looking for `application.yaml` (within workspace tree only)
3. Never traverse above workspace root (prevents vendor package hijacking)

**Override behavior**:
- Explicit flags completely replace context detection (no merging)
- When `-a` flag is specified, context-detected application is ignored

## Consequences

### Positive

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options
- Workspace boundary prevents accidental detection in vendor packages

### Negative

- More verbose than positional arguments when explicit targeting is needed
- Context detection adds complexity to command implementation

### Neutral

- Global commands (`port`, `proxy`, `config`) ignore directory context
- Error messages explain what context was or wasn't detected
