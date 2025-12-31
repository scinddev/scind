# Issue Group 2: Network Lifecycle Management

**Documents Affected**: Technical Spec, PRD
**Suggested Order**: 2 of 5 (network management is core to workspace isolation)
**Estimated Effort**: Small

---

## Overview

The workspace-internal network lifecycle has some ambiguous edge cases that need clarification.

---

## Issues

### A-2: Workspace Network Removal Timing

**Severity**: Medium

**Issue**: The Technical Spec states `workspace down` "optionally removes workspace network" (line 1166), but doesn't specify when this happens. Is the network removed when the last app is stopped, or only when all apps are stopped simultaneously via `workspace down`?

**Questions**:
1. If I run `workspace down -a app-one` and `workspace down -a app-two` separately, is the network removed after the second command?
2. Or is the network only removed when `workspace down` (without `-a`) is run?

**Suggested Resolution**: Clarify that the workspace network is only removed when `workspace down` is run without the `-a` flag (full workspace teardown), and document this explicitly.

**Response**:
> Only on full teardown — Network removed only with `workspace down` (no `-a` flag).

---

### A-3: Network Recreation on Partial Up

**Severity**: Low

**Issue**: The Technical Spec states the workspace network is created by `workspace up` if it doesn't exist (line 85). But if the network was removed during a previous `workspace down`, and the user runs `workspace up -a app-one` (single app), should the network be recreated?

**Questions**:
1. Should `workspace up -a app-one` recreate the workspace network if it doesn't exist?

**Suggested Resolution**: Yes, any `workspace up` (with or without `-a`) should ensure the network exists, since exported services may need cross-application communication later.

**Response**:
> Yes, always ensure network exists — Any `up` operation creates the network if missing.

---

## Checklist

- [x] Document network removal timing in Technical Spec Operations section
- [x] Clarify network creation behavior for partial `up` operations


---

## Archived

This issue was archived on 2024-12-31 at 14:03:26.
