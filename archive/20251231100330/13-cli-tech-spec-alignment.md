# Issue Group 13: CLI & Tech Spec Alignment

**Documents Affected**: Technical Spec + CLI Reference
**Suggested Order**: 13 of 18 (straightforward alignment fixes)
**Estimated Effort**: Small

---

## Overview

These issues involve commands or flags documented in one spec but missing from another. Straightforward alignment fixes.

---

## Issues

### N-1: `workspace destroy --keep-apps` Missing from CLI Reference

**Severity**: Low

**Issue**: The Technical Spec (lines 1170-1172) documents a `--keep-apps` flag for `workspace destroy`:

```
**Flags**:
- `--force`: Skip confirmation prompts for application directory removal
- `--keep-apps`: Preserve application directories (only remove workspace configuration)
```

However, the CLI Reference (lines 466-474) only documents `--force`:

```markdown
**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `--force` | Skip all confirmation prompts |
```

**Questions**:
1. Should `--keep-apps` be added to CLI Reference?
2. Or should it be removed from Technical Spec?

**Your Response**:
> Remove it from the Technical Spec.

---

### N-2: `config edit` Command in Go Stack but Not CLI Reference

**Severity**: Low

**Issue**: The Go Stack command mapping (line 904) lists:

```
| `contrail config edit` | `configCmd` → `configEditCmd` | Global (no context) |
```

However, the CLI Reference Config Commands section (lines 1021-1088) only documents:
- `config show`
- `config get`
- `config set`
- `config path`

There is no `config edit` command documented.

**Questions**:
1. Should `config edit` be added to CLI Reference?
2. What would `config edit` do—open the config file in `$EDITOR`?
3. Or should it be removed from Go Stack as unneeded?

**Your Response**:
> Add it to the CLI Reference.

---

### N-3: `proxy logs` Command Clarification

**Severity**: Low

**Issue**: The CLI Reference revision history (line 1505) mentions "removed logs command" but `proxy logs` is still fully documented in CLI Reference (lines 1005-1018):

```markdown
### `contrail proxy logs`

View Traefik proxy logs.

**Flags**:
| Flag | Description |
|------|-------------|
| `-f, --follow` | Follow log output |
| `--tail` | Number of lines from end |
```

Meanwhile, the Go Stack aliases.go description (line 182) mentions "logs" in the aliases list but the actual scaffolded code and command mapping don't include it.

**Questions**:
1. Is `proxy logs` intended to stay (it's proxy-specific, not app logs)?
2. If staying, should Go Stack add `proxyLogsCmd` to the mapping table?
3. Or should CLI Reference remove the `proxy logs` documentation?

**Your Response**:
> Remove `proxy logs` from CLI Reference.

---

## Checklist

- [x] Resolve `--keep-apps` flag documentation — removed from Tech Spec
- [x] Resolve `config edit` command documentation — added to CLI Reference
- [x] Resolve `proxy logs` command status — removed from CLI Reference, fixed Go Stack aliases comment

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
