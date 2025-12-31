# Issue Group 25: Workspace Prune Command in Go Stack

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-25: Workspace Prune Missing Dry-Run Flag in Go Stack

**Severity**: Low
**Category**: Missing Implementation Detail

**Finding**: The CLI Reference documents `workspace prune` with a `--dry-run` flag (lines 218-224):
```
**Flags**:
| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be removed without making changes |
```

The Go Stack includes `workspacePruneCmd` in the command scaffolding (lines 467-475), but does not include the `--dry-run` flag in its `init()` function (lines 514-542).

**Documents**:
- `contrail-cli-reference.md` (lines 218-224) - documents `--dry-run` flag
- `contrail-go-stack.md` (lines 467-475, 514-542) - missing `--dry-run` flag

**Recommendation**: Add the `--dry-run` flag to `workspacePruneCmd` in the Go Stack:
```go
workspacePruneCmd.Flags().Bool("dry-run", false, "show what would be removed without making changes")
```

---

## Tasks

### Task 25.1: Add Dry-Run Flag to Workspace Prune

Add the `--dry-run` flag to the `workspacePruneCmd` in the Go Stack scaffolding.

**Your Decision**:

> Yes, add it

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-25: Missing --dry-run flag in Go Stack | Low | ✅ Added `--dry-run` flag to `workspacePruneCmd` in Go Stack |


---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
