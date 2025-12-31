# Issue Group 21: Go Stack Missing App Exec Command

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-21: Missing `app exec` Command in Go Stack CLI Mapping

**Severity**: Low
**Category**: Missing Mapping

**Finding**: The CLI Reference documents an `app exec` command at the conceptual level (the `contrail-compose exec` provides this functionality), but there's no explicit `contrail app exec` command documented.

Looking at the Go Stack's CLI to Cobra Command Mapping (lines 879-924), there's no `app exec` listed, and the CLI Reference doesn't have a dedicated `contrail app exec` command section.

However, the PRD Quick Reference (lines 348-349) shows:
```
contrail-compose logs -f       # View logs for current app (via shell function)
```

This suggests the pattern is to use `contrail-compose` for exec, not a direct `contrail app exec`.

**Clarification Needed**: Is this intentional? The pattern seems to be:
- Use `contrail app up/down/restart/status` for lifecycle commands
- Use `contrail-compose exec` for container interaction

If intentional, this should be explicitly documented somewhere that `contrail app exec` is not provided because `contrail-compose exec` is the intended interface.

**Documents**:
- `contrail-cli-reference.md` - no `app exec` command
- `contrail-go-stack.md` - no `app exec` in mapping table

**Recommendation**: Either:
1. Add a note explaining that `contrail-compose exec` replaces `contrail app exec`, or
2. Decide if `contrail app exec` should be added as a shortcut

---

## Tasks

### Task 21.1: Clarify App Exec Design Decision

Document whether `contrail app exec` is intentionally absent (because `contrail-compose exec` exists) or if it should be added.

**Your Decision**:

> No action needed. There are no references to `contrail app exec` in the specs - `contrail-compose exec` is already the documented pattern. No additional explanation is necessary.

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-21: Missing `app exec` clarification | Low | ✅ No action needed |

