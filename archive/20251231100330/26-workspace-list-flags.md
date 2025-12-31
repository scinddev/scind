# Issue Group 26: Workspace List Flags in Go Stack

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-26: Workspace List Missing Flags in Go Stack

**Severity**: Low
**Category**: Missing Implementation Detail

**Finding**: The CLI Reference documents `workspace list` with two flags (lines 193-199):
```
**Flags**:
| Flag | Description |
|------|-------------|
| `--validate` | Check that registered paths still contain `workspace.yaml` |
| `--rebuild` | Rebuild registry from Docker labels (useful if registry is missing) |
```

The Go Stack includes `workspaceListCmd` in the command scaffolding (lines 424-431), but does not add these flags in the `init()` function.

**Documents**:
- `contrail-cli-reference.md` (lines 193-199) - documents `--validate` and `--rebuild` flags
- `contrail-go-stack.md` (lines 424-431) - missing these flags

**Recommendation**: Add flags to `workspaceListCmd`:
```go
workspaceListCmd.Flags().Bool("validate", false, "check that registered paths still contain workspace.yaml")
workspaceListCmd.Flags().Bool("rebuild", false, "rebuild registry from Docker labels")
```

---

## Tasks

### Task 26.1: Add Flags to Workspace List

Add the `--validate` and `--rebuild` flags to the `workspaceListCmd` in the Go Stack scaffolding.

**Your Decision**:

> Yes

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-26: Missing workspace list flags in Go Stack | Low | ✅ Added `--validate` and `--rebuild` flags to `workspaceListCmd` |


---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
