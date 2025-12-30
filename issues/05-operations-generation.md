# Issue Group 5: Operations & Generation

**Documents Affected**: Technical Spec + CLI Reference  
**Suggested Order**: 5 of 10 (depends on workspace concepts from Group 4)  
**Estimated Effort**: Medium

---

## Overview

These issues relate to the generate/up/down lifecycle—how Contrail detects when regeneration is needed, what happens during operations, and how state is managed.

---

## Issues

### A-1: Override File Staleness Detection Criteria Undefined

**Severity**: High

**Issue**: Multiple documents mention checking if override files are "stale" before regenerating, but never define the criteria.

**References**:
- CLI Reference line 261: "Check if override files are stale; regenerate if needed"
- Technical Spec line 843: "Check if override files are stale; regenerate if needed"

**Questions**:
1. What files are compared?
2. What comparison method is used (mtime? content hash?)?
3. Should compose files be included in staleness checks?

**Options**:

**A) mtime Comparison**
```
Stale if any of these are newer than .generated/*.override.yaml:
- workspace.yaml
- {app}/application.yaml
- {app}/docker-compose*.yaml (from active flavor)
```
- Pro: Simple, fast
- Con: Touch a file accidentally → unnecessary regeneration

**B) Content Hash**
```yaml
# .generated/state.yaml
checksums:
  workspace.yaml: sha256:abc123...
  app-one/application.yaml: sha256:def456...
```
- Pro: Accurate, no false positives
- Con: More complex, need to read files to compute hash

**C) Manifest Comparison**
- Regenerate config, compare to existing manifest
- If different, write new overrides
- Pro: Catches any change that affects output
- Con: Does most of the work before knowing if it's needed

**Suggested Resolution**: Option A (mtime) for simplicity, with `--force` flag to override. Document explicitly in Tech Spec.

**Response**:
> Use mtime comparison for staleness detection, with `--force` flag to override.

---

### A-4: Template Variable Resolution Timing Unclear

**Severity**: Medium

**Issue**: `%APPLICATION_FLAVOR%` is listed as available for templates (Tech Spec lines 349-361), but when is it resolved?

**Scenario**:
1. User runs `contrail generate` with flavor "lite"
2. Override files are generated with "lite" values
3. User runs `contrail flavor set full`
4. User runs `contrail up` (without explicit regeneration)

**Questions**:
1. Are templates resolved at `generate` time only?
2. Does `flavor set` automatically trigger regeneration?
3. Should `up` detect flavor changes and auto-regenerate?

**Options**:

**A) Generate-Time Only**
- Templates resolved when `generate` runs
- `flavor set` just updates state file
- User must run `generate` or `up` (which auto-generates if stale)
- Staleness check includes state.yaml flavor vs generated flavor

**B) Flavor Set Triggers Regeneration**
- `flavor set` immediately regenerates affected app's override
- Pro: Always consistent
- Con: Might regenerate unnecessarily if user is about to change other things

**C) Up-Time Resolution**
- `up` always checks if flavor in state differs from flavor in generated files
- Regenerates if mismatch
- Pro: Lazy, only regenerates when actually starting
- Con: Adds to startup time

**Suggested Resolution**: Option A with staleness check. Document that `flavor set` marks the app as needing regeneration, and `up` will handle it.

**Response**:
> Use Option A (generate-time resolution) combined with Option B (flavor set triggers regeneration). Templates are resolved when `generate` runs, and `flavor set` immediately regenerates affected app's override to ensure consistency.

---

### M-3: `workspace destroy` Missing from Technical Spec

**Severity**: Low

**Issue**: CLI Reference documents `workspace destroy` (lines 368-390) with a detailed sequence, but the Technical Spec's Operations section only covers `up` and `down`.

**CLI Reference describes**:
1. Run `workspace down --volumes`
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force`)
4. Remove `workspace.yaml`
5. Release any assigned ports

**Action**: Add to Technical Spec Operations section.

**Response**:
> Add destroy sequence to Tech Spec Operations section.

---

### M-8: Port Inventory External Release Detection Undefined

**Severity**: Medium

**Issue**: The spec shows port status transitions including `unavailable` → `released` via `contrail port gc`, but doesn't explain how Contrail determines an external process has released a port.

**Location**: Technical Spec lines 290-294

**Current documentation**:
> - `unavailable` → `released`: External process stopped, `contrail port gc` cleaned it up

**Questions**:
1. How does `port gc` (or `port scan`) check if a port is actually bound?
2. What method is used (attempt to bind? check /proc/net? netstat parsing?)?

**Suggested Resolution**: Document in Tech Spec:
> `contrail port scan` checks port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state.

**Response**:
> Use the suggested `net.Listen("tcp", ":PORT")` method for port availability checking. Document this in the Tech Spec.

---

### M-9: Manual Override Behavior During Regeneration Undefined

**Severity**: Medium

**Issue**: What happens to `overrides/{app}.yaml` files during `workspace generate --force`? The relationship between manual overrides and regeneration isn't clarified.

**Location**: Technical Spec lines 647-666

**Questions**:
1. Are manual overrides in `overrides/` always preserved?
2. Is the merge order explicitly documented (base → generated → manual)?
3. Should `generate` warn if manual overrides exist?
4. What if manual overrides conflict with generated config?

**Current implicit behavior** (from directory structure):
- `overrides/` is separate from `.generated/`
- `--force` regenerates `.generated/*`, shouldn't touch `overrides/`
- Compose command includes both: `-f generated.yaml -f overrides/app.yaml`

**Suggested Resolution**: Document explicitly in Tech Spec:

> **Manual Overrides**
> 
> Files in `{workspace}/overrides/{app}.yaml` are never modified by Contrail. They are merged after generated overrides in the Docker Compose command:
> 
> ```
> docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
> ```
> 
> This allows workspace-specific customizations (extra environment variables, volume mounts, middleware) that persist across regeneration.
>
> `workspace generate --force` only affects `.generated/` contents.

**Response**:
> Document manual override preservation and merge order explicitly in Tech Spec. Files in `overrides/` are never modified by Contrail and are merged after generated overrides.

---

## Checklist

- [x] Define and document staleness detection criteria in Tech Spec
- [x] Document template variable resolution timing
- [x] Clarify relationship between `flavor set` and regeneration
- [x] Add `workspace destroy` sequence to Tech Spec Operations
- [x] Document port availability checking method
- [x] Document manual override preservation and merge order
