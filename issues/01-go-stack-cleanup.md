# Issue Group 1: Go Stack Cleanup

**Documents Affected**: Go Stack only  
**Suggested Order**: 1 of 10 (quick wins, no decisions needed)  
**Estimated Effort**: Small

---

## Overview

These are straightforward fixes to the Go Stack document—missing command mappings and a typo. No cross-document coordination required.

---

## Issues

### M-4: `contrail app exec` Missing from Go Stack Command Mapping

**Severity**: Low

**Issue**: `contrail app exec` is documented in CLI Reference (lines 606-627) but absent from the Go Stack Cobra command mapping table (lines 746-756).

**Fix**: Add to Go Stack command mapping table:
```
| `contrail app exec` | `appCmd` → `appExecCmd` | |
```

**Response**:  
> Either remove references to `contrail app exec` or replace them with the new `contrail-compose exec` we landed on.

---

### M-5: `contrail flavor reset` Missing from Go Stack Command Mapping

**Severity**: Low

**Issue**: `contrail flavor reset` is documented in CLI Reference (lines 705-718) but missing from Go Stack mapping.

**Fix**: Add to Go Stack command mapping table:
```
| `contrail flavor reset` | `flavorCmd` → `flavorResetCmd` | |
```

**Response**:  
> I'm not sure what `contrail flavor reset` was meant to do, and I am not sure we need it. Unless you have a compelling reason to think otherwise, remove references to `contrail flavor reset`.

---

### X-1: Go Stack Has Typo in `BoolVarP` Call

**Severity**: Low

**Issue**: Line 335 has `BoolVarP(&quiet, "quiet", "q", "", ...)` where the empty string should be `false`.

**Fix**: Change to:
```go
rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output, suppress context indicators")
```

**Response**:  
> Go ahead and make this change.

---

## Checklist

- [x] ~~Add `appExecCmd` to command mapping table~~ Removed `contrail app exec` from CLI Reference (use `contrail-compose exec` instead)
- [x] ~~Add `flavorResetCmd` to command mapping table~~ Removed `contrail flavor reset` from CLI Reference (not needed)
- [x] Fix `BoolVarP` typo on line 335
