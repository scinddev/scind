# Issue Group 8: Shell Integration

**Documents Affected**: CLI Reference + Shell Integration + Go Stack  
**Suggested Order**: 8 of 10 (depends on CLI commands being finalized)  
**Estimated Effort**: Small

---

## Overview

These issues are specific to shell integration—the `contrail-compose` function and its interaction with Docker Compose's own flags.

---

## Issues

### C-5: `-f` Flag Collision Between Contrail and Docker Compose

**Severity**: High

**Issue**: The Shell Integration spec (line 89) gives `contrail compose-prefix` a `-f`/`--flavor` flag:

```bash
contrail compose-prefix -f lite
```

But `-f` is Docker Compose's standard flag for specifying compose files:

```bash
docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up
```

This creates ambiguity in `contrail-compose`:

```bash
contrail-compose -f lite up        # Is -f for flavor or compose file?
contrail-compose -f myfile.yaml up # Same ambiguity
```

The shell function currently parses `-w`/`--workspace` and `-a`/`--app` as Contrail flags, passing everything else to Docker. But `-f` would be captured by Contrail when the user meant it for Docker.

**Location**: Shell Integration line 89

**Options**:

**A) Use Different Short Flag for Flavor**
- Change to `-F`/`--flavor` (capital F)
- Pro: No collision, clear distinction
- Con: Capital flags are unconventional

**B) Long Flag Only for Flavor**
- Remove `-f` short form, keep only `--flavor`
- Pro: No collision, follows principle of least surprise
- Con: More typing for flavor changes

**C) Change Flag Name Entirely**
- Use `--with-flavor` or `--use-flavor`
- Pro: Very explicit
- Con: Verbose

**D) Positional Flavor in compose-prefix Only**
- `contrail compose-prefix --flavor=lite` (internal command)
- `contrail-compose` doesn't expose flavor flag at all
- Users change flavor with `contrail flavor set` before using `contrail-compose`
- Pro: Clean separation, no collision
- Con: Less convenient for one-off flavor testing

**Suggested Resolution**: Option B (long flag only). Update:
- Shell Integration: Remove `-f` short form
- CLI Reference: Document `--flavor` only for `compose-prefix`
- Go Stack: Remove short flag registration

**Response**:
> Use Option D. Flavor should not be changed at runtime via `contrail-compose` as it would require regenerating everything and could have major impact on a running application. Users must use `contrail flavor set` before using `contrail-compose`. Remove the `--flavor` flag from `compose-prefix` entirely. See also: new issue added for documenting `flavor set` behavior when application is already running.

---

### X-3: Exit Code 5 Used But Not Cross-Referenced

**Severity**: Low

**Issue**: Exit code 5 (context detection failed) is documented in CLI Reference (line 1281) but Shell Integration uses it (lines 100-114) without explicit documentation.

**CLI Reference**:
```
| 5 | Context detection failed (workspace/app not found) |
```

**Shell Integration** (shows error, references exit but doesn't cite code):
```bash
$ contrail compose-prefix
Error: No application context detected.
...
$ echo $?
5
```

**Action**: Add exit code reference to Shell Integration doc for completeness.

**Suggested Addition** to Shell Integration:

```markdown
### Exit Codes

`contrail compose-prefix` uses standard Contrail exit codes:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 5 | Context detection failed (no workspace/app found) |

The shell function checks for empty output to detect failures and re-runs 
the command to display the error message.
```

**Response**:
> Yes, add exit codes reference table to Shell Integration doc.

---

## Checklist

- [x] Remove `--flavor` flag from `compose-prefix` command entirely (Option D)
- [x] Update Shell Integration spec to remove flavor flag
- [x] Update CLI Reference `compose-prefix` documentation
- [x] Update Go Stack flag registration
- [x] Add exit code documentation to Shell Integration
- [x] Create new issue (Group 12) for flavor set behavior when app is running

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
