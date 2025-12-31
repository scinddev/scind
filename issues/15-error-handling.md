# Issue Group 15: Error Handling & Edge Cases

**Documents Affected**: Technical Spec + CLI Reference + Go Stack
**Suggested Order**: 15 of 18 (important for implementation, can defer)
**Estimated Effort**: Medium

---

## Overview

These issues concern undefined behavior in error scenarios and edge cases that will surface during implementation.

---

## Issues

### N-8: Concurrent Workspace Operations Undefined

**Severity**: Medium

**Issue**: No document addresses concurrent operation behavior:
- Two terminals running `workspace up` simultaneously
- One terminal running `workspace down` while another runs `workspace up`
- Multiple users on shared host with conflicting port assignments

**Questions**:
1. Should Contrail use file locking for workspace operations?
2. Should port assignment use atomic file operations?
3. Should there be explicit error messages for race conditions?
4. Is this out of scope for v1 (document as known limitation)?

**Your Response**:
> Out of scope for v1—document as a known limitation.

---

### N-9: Git Clone Failure Behavior

**Severity**: Low

**Issue**: `workspace clone` behavior is undefined when:
- A repository URL is invalid or unreachable
- Clone fails mid-way (network interruption)
- Repository exists but has no default branch
- Authentication fails

**Location**: CLI Reference lines 307-327

**Questions**:
1. Should partial clones be cleaned up on failure?
2. Should the command continue with other apps or abort entirely?
3. What exit code should be used?

**Your Response**:
> Abort on first failure and clean up partial clones on failure.

---

### N-10: Docker/Docker Compose Not Installed

**Severity**: Low

**Issue**: The `doctor` command checks for Docker (CLI Reference line 1231) but behavior of other commands when Docker is missing is unspecified.

**Questions**:
1. Should all commands that need Docker check for it upfront?
2. What error message should be displayed?
3. Should Contrail provide installation guidance?

**Suggested Resolution**: Add a common check in commands that need Docker, with a clear error:
```
Error: Docker is not installed or not running.
Run 'contrail doctor' for setup guidance.
```

**Your Response**:
> Clear error message pointing to `contrail doctor`.

---

### N-11: Non-Existent Compose Service Referenced

**Severity**: Medium

**Issue**: If `application.yaml` references a Compose service that doesn't exist in the compose files, the behavior is undefined.

Example:
```yaml
# application.yaml
exported_services:
  api:
    service: backend  # <-- Compose service "backend" doesn't exist
```

**Questions**:
1. Should this be validated at `generate` time?
2. Should the error message suggest available service names?
3. Or let Docker Compose fail with its own error?

**Suggested Resolution**: Validate at `generate` time with helpful error:
```
Error: Exported service "api" references non-existent Compose service: backend
  Application: my-app
  Available services in docker-compose.yaml: web, db, redis
```

**Your Response**:
> Validate at `generate` time with helpful error listing available services.

---

### N-12: Orphaned Port Cleanup Timing

**Severity**: Low

**Issue**: The `port gc` command exists but there's no guidance on:
- How often should it be run?
- Should `workspace up` trigger automatic cleanup?
- What happens to ports from workspaces deleted via filesystem (not `workspace destroy`)?

**Location**: CLI Reference lines 857-874

**Questions**:
1. Should `port gc` run automatically on `workspace up`?
2. Should there be a warning threshold (e.g., >10 stale ports)?
3. Is manual `port gc` sufficient for v1?

**Your Response**:
> Manual `port gc` is sufficient for v1—document as a known limitation.

---

## Checklist

- [x] Document concurrent operation handling — added to PRD Known Limitations
- [x] Document git clone failure behavior — added to CLI Reference workspace clone
- [x] Add Docker availability check guidance — added Error Messages section to CLI Reference
- [x] Document Compose service validation approach — added to Tech Spec Generation Logic
- [x] Clarify port garbage collection timing — added to PRD Known Limitations
