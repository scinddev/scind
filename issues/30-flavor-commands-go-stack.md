# Issue Group 30: Flavor Commands Missing from Go Stack Scaffolding

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-30: Flavor Commands Not Scaffolded in Go Stack

**Severity**: Low
**Category**: Missing Scaffolding

**Finding**: The Go Stack's CLI to Cobra Command Mapping table (lines 899-901) lists:
- `contrail flavor list` → `flavorCmd` → `flavorListCmd`
- `contrail flavor show` → `flavorCmd` → `flavorShowCmd`
- `contrail flavor set` → `flavorCmd` → `flavorSetCmd`

However, there's no scaffolding provided for these commands. The CLI Reference (lines 687-771) documents:
- `flavor list` with workspace and app flags
- `flavor show` with workspace and app flags
- `flavor set` with positional `<flavor>` argument plus workspace and app flags

The Go Stack should provide scaffolding similar to workspace.go showing:
- The flavor argument as positional for `set`
- Proper flag definitions

**Documents**:
- `contrail-cli-reference.md` (lines 687-771) - complete flavor command documentation
- `contrail-go-stack.md` (lines 899-901) - lists commands but no scaffolding

**Recommendation**: Add scaffolding for flavor commands, especially noting the positional argument pattern for `flavor set`.

---

## Tasks

### Task 30.1: Add Flavor Commands Scaffolding

Add scaffolding for flavor commands to the Go Stack specification.

**Your Decision**:

> Yes

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-30: Missing flavor commands scaffolding | Low | ✅ Added `internal/cli/flavor.go` scaffolding with all flavor commands |

