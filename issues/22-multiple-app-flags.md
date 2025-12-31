# Issue Group 22: Multiple --app Flag Behavior

**Documents**: CLI Reference, Go Stack
**Effort**: Small

---

## Issues

### Issue N-22: Repeatable --app Flag Not Reflected in Go Stack

**Severity**: Medium
**Category**: Implementation Gap

**Finding**: The CLI Reference documents `--app` as repeatable for several commands:
- `workspace up` (line 372): `| -a, --app | Start specific app(s) only (repeatable) |`
- `workspace down` (line 407): `| -a, --app | Stop specific app(s) only (repeatable) |`
- `workspace clone` (line 313): `| -a, --app | Clone specific app(s) only (repeatable) |`
- `workspace generate` (line 348): `| -a, --app | Generate for specific app(s) only (repeatable) |`
- `workspace restart` (line 432): `| -a, --app | Restart specific app(s) only (repeatable) |`

However, the Go Stack scaffolding (line 334) defines:
```go
rootCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "specify application (overrides context detection)")
```

This uses `StringVarP` which only captures a single value. For repeatable flags, it should use `StringSliceVarP` or `StringArrayVarP`:
```go
rootCmd.PersistentFlags().StringSliceVarP(&apps, "app", "a", nil, "specify application(s) (overrides context detection)")
```

**Documents**:
- `contrail-cli-reference.md` - documents repeatable `-a` flag
- `contrail-go-stack.md` (line 334) - uses StringVarP (single value)

**Recommendation**: Update Go Stack to use `StringSliceVarP` for the `--app` flag to support multiple applications.

---

## Tasks

### Task 22.1: Update Go Stack App Flag Definition

Change the `--app` flag from `StringVarP` to `StringSliceVarP` to support repeatable usage.

**Your Decision**:

> Yes, update Go Stack to ensure --app flag is repeatable in the places that it is listed as repeatable in the CLI Reference.

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-22: Repeatable --app flag implementation | Medium | ✅ Fixed |

