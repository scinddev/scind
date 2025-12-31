# Issue Group 29: Port Commands Missing from Go Stack Scaffolding

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-29: Port Commands Not Fully Documented in Go Stack

**Severity**: Low
**Category**: Missing Scaffolding

**Finding**: The CLI Reference documents these port commands:
- `port list` (lines 779-809)
- `port show` (lines 813-826)
- `port release` (lines 828-846)
- `port assign` (lines 848-862)
- `port gc` (lines 864-882)
- `port scan` (lines 884-896)

The Go Stack's mapping table (lines 902-904) only shows:
- `contrail port list` → `portCmd` → `portListCmd`
- `contrail port release` → `portCmd` → `portReleaseCmd`
- `contrail port gc` → `portCmd` → `portGcCmd`

Missing from the mapping:
- `port show`
- `port assign`
- `port scan`

**Documents**:
- `contrail-cli-reference.md` (lines 779-896) - documents all 6 port commands
- `contrail-go-stack.md` (lines 902-904) - only lists 3 port commands

**Recommendation**: Add `portShowCmd`, `portAssignCmd`, and `portScanCmd` to the Go Stack mapping table.

---

## Tasks

### Task 29.1: Add Missing Port Commands to Go Stack

Update the CLI to Cobra Command Mapping to include all 6 port commands.

**Your Decision**:

> Yes, add it

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-29: Missing port commands in mapping | Low | ✅ Added `port show`, `port assign`, `port scan` to Go Stack mapping |


---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
