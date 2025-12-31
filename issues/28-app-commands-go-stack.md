# Issue Group 28: App Commands Missing from Go Stack Scaffolding

**Documents**: Go Stack, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-28: App Commands Not Scaffolded in Go Stack

**Severity**: Low
**Category**: Missing Scaffolding

**Finding**: The Go Stack's CLI to Cobra Command Mapping table (lines 879-924) lists all the app commands:
- `contrail app list` → `appCmd` → `appListCmd`
- `contrail app show` → `appCmd` → `appShowCmd`
- `contrail app init` → `appCmd` → `appInitCmd`
- `contrail app add` → `appCmd` → `appAddCmd`
- `contrail app remove` → `appCmd` → `appRemoveCmd`
- `contrail app up` → `appCmd` → `appUpCmd`
- `contrail app down` → `appCmd` → `appDownCmd`
- `contrail app restart` → `appCmd` → `appRestartCmd`
- `contrail app status` → `appCmd` → `appStatusCmd`

However, the Go Stack only provides complete scaffolding for:
- Workspace commands (`internal/cli/workspace.go`)
- Top-level aliases (`internal/cli/aliases.go`)
- compose-prefix (`internal/cli/compose_prefix.go`)
- Proxy commands (`internal/cli/proxy.go`)

There's no scaffolded `internal/cli/app.go` file like there is for workspace commands.

**Documents**:
- `contrail-go-stack.md` (lines 879-898) - lists app commands in mapping
- `contrail-go-stack.md` - no `app.go` scaffolding provided

**Recommendation**: Add scaffolding for app commands similar to workspace commands, including proper flag definitions from the CLI Reference.

---

## Tasks

### Task 28.1: Add App Commands Scaffolding

Add scaffolding for app commands to the Go Stack specification.

**Your Decision**:

> Yes, add it

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-28: Missing app commands scaffolding | Low | ✅ Added `internal/cli/app.go` scaffolding with all app commands and flags |

