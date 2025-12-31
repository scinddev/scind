# Issue Group 27: --keep-apps Flag in Workspace Destroy

**Documents**: Technical Spec, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-27: Missing --keep-apps Flag from Workspace Destroy

**Severity**: Low
**Category**: Missing Feature

**Finding**: The Technical Spec's Destroy Sequence (lines 1168-1181) mentions:
```
3. Prompt before removing application directories (unless `--force`)
```

And the flags section shows only:
```
**Flags**:
- `--force`: Skip confirmation prompts for application directory removal
```

However, there's no `--keep-apps` flag that would allow preserving application directories without requiring an interactive prompt. The CLI Reference similarly only shows `--force` (lines 474-479).

Common use case: A developer wants to destroy a workspace but keep the application directories because they have local changes or want to recreate the workspace with different settings. Currently they would need to:
1. Run `workspace destroy`
2. Answer "no" to the prompt about removing apps

A `--keep-apps` flag would make this non-interactive.

**Documents**:
- `contrail-technical-spec.md` (lines 1168-1181) - no `--keep-apps` flag
- `contrail-cli-reference.md` (lines 466-488) - no `--keep-apps` flag

**Recommendation**: Consider adding a `--keep-apps` flag to `workspace destroy` that skips application directory removal without prompting.

---

## Tasks

### Task 27.1: Consider Adding --keep-apps Flag

Decide whether to add a `--keep-apps` flag to `workspace destroy`.

**Your Decision**:

> Yes, add --keep-apps to both Technical Spec and CLI Reference

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-27: Missing --keep-apps flag | Low | ✅ Added `--keep-apps` flag to Tech Spec and CLI Reference |

