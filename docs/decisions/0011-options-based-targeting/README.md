# ADR-0011: Options-Based Targeting with Context Detection

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Commands need to target specific workspaces and applications. Targeting could use positional arguments (`contrail up dev app-one`) or options (`contrail up --workspace=dev --app=app-one`). With context detection, explicit targeting is often unnecessary.

## Decision

Use `--workspace` / `-w` and `--app` / `-a` options (not positional arguments) with automatic context detection from current directory.

**Context detection algorithm**:
1. Walk up from current directory looking for `workspace.yaml` to establish workspace root
2. Walk up toward workspace root looking for `application.yaml` for app context
3. Never traverse above workspace root for application detection

## Consequences

### Positive

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed for scripts or CI
- Easy to extend with additional targeting options
- Prevents vendor packages or nested fixtures from hijacking context

### Negative

- Longer flag names compared to short positional args
- Context detection may be surprising if users don't expect it

### Neutral

- Works seamlessly with both multi-app and single-app workspaces

---

## Notes

Context detection uses a workspace boundary approach: `application.yaml` files outside the detected workspace root are ignored, preventing accidental detection of config files in vendor packages.
