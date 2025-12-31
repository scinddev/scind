# Issue Group 4: CLI Output and Feedback Behavior

**Documents Affected**: CLI Reference, Go Stack
**Suggested Order**: 4 of 5 (user experience consistency)
**Estimated Effort**: Small

---

## Overview

Several CLI output behaviors need clarification for consistent user experience.

---

## Issues

### A-6: Progress Output During Long Operations

**Severity**: Low

**Issue**: The CLI Reference doesn't specify what output the user sees during potentially long operations like `workspace clone` (cloning multiple repositories) or `workspace up` (starting multiple applications).

**Questions**:
1. Should these operations show progress per-application (e.g., "Cloning app-one... done", "Cloning app-two...")?
2. Should `--quiet` suppress all progress output?

**Suggested Resolution**: Show per-application progress by default, with `--quiet` suppressing everything except final success/error. Add a sentence to CLI Reference Global Flags section.

**Response**:
> Per-app progress, quiet suppresses — Show "Starting app-one... done" by default, `--quiet` shows nothing except errors.

---

### A-7: Context Feedback Suppression

**Severity**: Low

**Issue**: The CLI Reference states that `--quiet` suppresses "context indicators" (line 105), but it's unclear what other output is suppressed. For example, does `contrail app status --quiet` show just "running" or show nothing?

**Questions**:
1. Should `--quiet` on status commands show just the status value, or nothing?

**Suggested Resolution**: `--quiet` shows minimal machine-readable output (status value, names-only for list commands). Add examples to CLI Reference.

**Response**:
> Minimal machine-readable output — `--quiet` suppresses progress and context indicators, but status commands output just the value (e.g., "running"), list commands output names only (one per line). Errors are always shown.

---

### M-4: Missing `--color` Flag

**Severity**: Low

**Issue**: The CLI Reference documents `NO_COLOR` environment variable (line 1421) for disabling colored output, but doesn't mention a `--color` or `--no-color` flag. Some CI environments may need explicit control.

**Questions**:
1. Should Contrail support `--color=auto|always|never` flag?

**Suggested Resolution**: Add `--color` flag for explicit control, defaulting to `auto` (detect terminal). Document alongside `NO_COLOR` in Environment Variables section.

**Response**:
> Add `--color` flag — Add flag with `auto|always|never` values, defaulting to `auto`.

---

## Checklist

- [x] Document progress output behavior in CLI Reference
- [x] Clarify `--quiet` semantics with examples
- [x] Add `--color` flag to CLI Reference and Go Stack


---

## Archived

This issue was archived on 2024-12-31 at 14:03:26.
