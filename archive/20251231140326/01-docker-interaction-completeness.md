# Issue Group 1: Docker Interaction Completeness

**Documents Affected**: Technical Spec, Go Stack
**Suggested Order**: 1 of 5 (foundational patterns affect implementation)
**Estimated Effort**: Medium

---

## Overview

Several details about Docker/Compose interaction patterns need clarification to ensure consistent implementation.

---

## Issues

### M-1: Missing `docker compose` Error Handling Pattern

**Severity**: Medium

**Issue**: The Go Stack describes shelling out to `docker compose` via `exec.Command` (contrail-go-stack.md:147-152), but doesn't define how to handle errors from docker compose commands. The Technical Spec also doesn't specify error handling for failed docker compose operations.

**Questions**:
1. Should Contrail capture and reformat docker compose errors, or pass them through directly to the user?

**Suggested Resolution**: Add error handling pattern to Go Stack that specifies pass-through with optional context prefix (e.g., "Failed to start app-one: <docker compose error>").

**Response**:
> Pass-through with context — Prefix errors with Contrail context (e.g., "Failed to start app-one:") then show full docker compose output.

---

### M-2: Missing `--remove-orphans` Consistency

**Severity**: Low

**Issue**: The Technical Spec mentions passing `--remove-orphans` to Docker Compose when flavors change (contrail-technical-spec.md:526), but this behavior isn't documented for `workspace down` or `app down` operations. Should orphaned containers always be removed, or only during flavor changes?

**Questions**:
1. Should `--remove-orphans` be passed to all `docker compose down` invocations by default?

**Suggested Resolution**: Document that `--remove-orphans` is passed for all `up` operations to handle flavor changes and manual compose file edits, but not for `down` (where it's the default behavior).

**Response**:
> Yes, always for `up` — Pass `--remove-orphans` on all `up` operations.

---

### A-1: Ambiguous Container Restart Behavior

**Severity**: Low

**Issue**: The `contrail app restart` and `contrail workspace restart` commands are documented as "equivalent to `down` followed by `up`" (CLI Reference line 436), but this doesn't specify whether volumes are preserved. The `down` command has a `--volumes` flag, but restart doesn't document this.

**Questions**:
1. Should `restart` preserve volumes by default (i.e., not pass `--volumes` to the internal `down`)?

**Suggested Resolution**: Document that `restart` always preserves volumes (does not pass `--volumes` to internal down operation).

**Response**:
> Yes, always preserve volumes — `restart` never removes volumes.

---

## Checklist

- [x] Add docker compose error handling pattern to Go Stack
- [x] Document `--remove-orphans` behavior in Technical Spec
- [x] Clarify `restart` volume preservation in CLI Reference


---

## Archived

This issue was archived on 2024-12-31 at 14:03:26.
