# Issue Group 20: Workspace Destroy and Registry Update

**Documents**: Technical Spec, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-20: Workspace Destroy Registry Removal Not Documented in CLI Reference

**Severity**: Low
**Category**: Documentation Gap

**Finding**: The Technical Spec's "Destroy Sequence" (lines 1168-1178) explicitly mentions:
```
6. Remove workspace from registry (`~/.config/contrail/workspaces.yaml`)
```

However, the CLI Reference's `contrail workspace destroy` documentation (lines 466-488) does not mention this registry removal step. The CLI Reference only mentions:
1. Run `workspace down --volumes`
2. Remove `.generated/` directory
3. Prompt before removing application directories
4. Remove `workspace.yaml`
5. Release any assigned ports

**Documents**:
- `contrail-technical-spec.md` (lines 1168-1178) - includes registry removal
- `contrail-cli-reference.md` (lines 466-488) - missing registry removal step

**Recommendation**: Add "Remove workspace from registry" as step 6 in the CLI Reference's `workspace destroy` behavior section.

---

## Tasks

### Task 20.1: Update CLI Reference with Registry Removal Step

Add a step to the CLI Reference's `workspace destroy` behavior to mention removing the workspace from the registry.

**Your Decision**:

> Yes

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-20: Registry removal missing from CLI docs | Low | ✅ Fixed |


---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
