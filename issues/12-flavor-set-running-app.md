# Issue Group 12: Flavor Set Behavior for Running Applications

**Documents Affected**: Tech Spec, CLI Reference
**Suggested Order**: 12 of 12 (discovered during Group 8 review)
**Estimated Effort**: Small

---

## Overview

This issue was identified during the C-5 resolution (Group 8). When deciding to remove the `--flavor` flag from `compose-prefix`, it became clear that changing flavors at runtime has implications that aren't fully documented.

---

## Issues

### A-11: Flavor Set Behavior When Application Is Running Undefined

**Severity**: Medium

**Issue**: The `contrail flavor set` command can change an application's active flavor, which affects which Docker Compose files are used. However, the documentation doesn't clearly explain what happens when the application is already running.

**Current documentation** (Tech Spec, Template Resolution Timing):
> When `contrail flavor set FLAVOR` is executed, it:
> 1. Updates `.generated/state.yaml` with the new flavor
> 2. Immediately regenerates the affected application's override file

**Questions**:
1. Does the running application automatically pick up the new flavor?
2. Does the user need to run `contrail app restart` or `contrail workspace up` to apply changes?
3. What if the new flavor adds/removes services? (e.g., "lite" has no worker, "full" has worker)
4. Should `flavor set` warn if the application is running?
5. Should `flavor set` offer to restart the application?

**Potential Scenarios**:

**Scenario A: Flavor adds a service**
```bash
contrail flavor set full  # full flavor includes worker service
# Worker service is defined in override but not running
# User must run `contrail up` to start the worker
```

**Scenario B: Flavor removes a service**
```bash
contrail flavor set lite  # lite flavor excludes worker service
# Worker service is still running but no longer in override
# What happens on next `contrail up`? Does worker get stopped?
```

**Scenario C: Flavor changes environment variables**
```bash
contrail flavor set debug  # debug flavor sets DEBUG=true
# Running containers still have old environment
# Restart required to pick up changes
```

**Options**:

**A) Document Only (Minimal)**
- Document that `flavor set` regenerates overrides but doesn't affect running containers
- User must manually restart/up to apply changes
- Pro: Simple, explicit
- Con: Users may be confused when changes don't take effect

**B) Warn When Running**
- `flavor set` checks if application is running
- If running, warns: "Application is running. Run `contrail app restart` to apply flavor changes."
- Pro: Helpful guidance
- Con: More implementation complexity

**C) Offer to Restart**
- `flavor set` offers: "Application is running. Restart now? [y/N]"
- Or add `--restart` flag: `contrail flavor set full --restart`
- Pro: Convenient
- Con: May be unexpected; adds complexity

**Suggested Resolution**: Option B (warn when running). Keep `flavor set` focused on its single responsibility, but provide helpful guidance when the application is running.

**Response**:
> Use Option B (warn when running). When `flavor set` detects that the application is running, it should display a warning with guidance to restart. For Scenario B (flavor removes a service), `contrail up` should stop orphaned services—services that are running but no longer defined in the active configuration.

---

## Checklist

- [x] Document `flavor set` behavior with running applications in Tech Spec
- [x] Add warning about running applications to CLI Reference
- [x] Document orphaned service handling (`--remove-orphans` on `contrail up`)
- [ ] Consider adding `--restart` flag (optional enhancement, deferred to implementation)

