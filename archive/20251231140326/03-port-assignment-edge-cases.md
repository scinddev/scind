# Issue Group 3: Port Assignment Edge Cases

**Documents Affected**: Technical Spec, CLI Reference
**Suggested Order**: 3 of 5 (port management is critical for workspace isolation)
**Estimated Effort**: Medium

---

## Overview

The port assignment system has edge cases and interaction patterns that need clarification.

---

## Issues

### A-4: Port Conflict During Generation

**Severity**: Medium

**Issue**: The Technical Spec describes port assignment during generation (lines 417-421), but doesn't specify what happens if a port is assigned in `state.yaml` but has become unavailable (e.g., another process took it) by the time containers start.

**Questions**:
1. Should Contrail re-check port availability at `workspace up` time and reassign if needed?
2. Or should it fail with an error pointing the user to run `contrail port gc` and regenerate?

**Suggested Resolution**: Fail at `up` time with a clear error message and suggest running `contrail port scan && contrail generate --force` to reassign ports.

**Response**:
> Fail with guidance — Fail at startup with clear error and suggest `port scan && generate --force` to reassign ports.

---

### A-5: Port Assignment Persistence Across Regeneration

**Severity**: Low

**Issue**: The Technical Spec states port assignments are "sticky" (line 420-421), meaning subsequent runs use the recorded port. But it's unclear whether `workspace generate --force` preserves existing port assignments or reassigns ports fresh.

**Questions**:
1. Should `--force` regeneration preserve sticky port assignments?

**Suggested Resolution**: Yes, `--force` should only regenerate override files, not reassign ports. Port reassignment should require explicit `contrail port release` or `contrail port gc`.

**Response**:
> Preserve ports on `--force` — Force only affects override files, not port assignments.

---

### M-3: Missing Port Range Configuration

**Severity**: Low

**Issue**: The Go Stack port command scaffolding shows `--available` flag to "show only available ports in configured range" (line 1075), but no port range configuration is documented in any spec. The proxy.yaml schema doesn't include port range settings.

**Questions**:
1. Should there be a configurable port range for auto-assignment (e.g., `port_range: 5400-5500`)?
2. Or should ports always be assigned starting from the requested port and incrementing?

**Suggested Resolution**: Add optional `port_range` configuration to proxy.yaml with sensible defaults (e.g., 10000-60000 for auto-increment). If not configured, use the increment-from-requested behavior as documented.

**Response**:
> Keep increment behavior only — Remove the `--available` flag reference, keep simple increment from requested port.

---

## Checklist

- [x] Document port conflict handling at startup time in Technical Spec
- [x] Clarify `--force` behavior with port assignments
- [x] Remove `--available` flag from Go Stack (no port range config needed)


---

## Archived

This issue was archived on 2024-12-31 at 14:03:26.
