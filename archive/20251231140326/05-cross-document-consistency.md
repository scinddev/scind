# Issue Group 5: Cross-Document Consistency

**Documents Affected**: All specs
**Suggested Order**: 5 of 5 (polish and alignment)
**Estimated Effort**: Small

---

## Overview

Minor inconsistencies between documents that should be aligned for clarity.

---

## Issues

### C-1: Inconsistent Default Flavor Naming

**Severity**: Low

**Issue**: The Technical Spec states the default flavor is resolved as: "CLI → state → default_flavor → 'default'" (line 688-691), implying the fallback is literally the string "default". However, example configurations show `default_flavor: full` or `default_flavor: default` explicitly.

The PRD example (line 404) doesn't show a `default_flavor` key at all, which could be confusing since readers might not realize "default" is the implicit fallback.

**Questions**:
1. Should examples consistently show `default_flavor` to make the resolution chain clear?

**Suggested Resolution**: Add `default_flavor: default` to the PRD example for consistency, and add a note that this key is optional.

**Response**:
> Add to example — Add `default_flavor: default` with comment noting it's optional.

---

### C-2: Manual Override Path Inconsistency

**Severity**: Low

**Issue**: The Technical Spec shows manual overrides at `{workspace}/overrides/{application-name}.yaml` (line 834-835), but the PRD Directory Structure (line 483-484) shows `overrides/app-one.yaml` without explaining the naming convention.

The example override filename `app-one.yaml` matches the application name, but this should be explicitly documented.

**Questions**:
1. Should the naming convention for override files be explicitly documented?

**Suggested Resolution**: Add a note in the PRD Directory Structure section: "Override files must be named `{application-name}.yaml` to be merged."

**Response**:
> Add note — Explicitly document the naming convention.

---

### X-1: Revision History Date Format Inconsistency

**Severity**: Low

**Issue**: The Shell Integration spec (line 842) uses "Dec 2024" in revision history, while the Technical Spec (line 1375-1376) uses "Dec 2024". This is actually consistent, but the PRD revision history (line 656) uses the same format. However, some entries say "December 2024" (Technical Spec created date line 4) vs "Dec 2024" in tables.

**Questions**:
1. Should all dates use a consistent format?

**Suggested Resolution**: Use "Dec 2024" in all revision history tables, and "December 2024" in document headers for full month names.

**Response**:
> Keep current pattern — Headers use full month ("December 2024"), tables use abbreviated ("Dec 2024"). This is already consistent.

---

### A-8: Single-App Workspace `application.yaml` Location

**Severity**: Medium

**Issue**: For single-app workspaces (`path: .`), both `workspace.yaml` and `application.yaml` are in the same directory. The Technical Spec mentions this (line 193-203), but doesn't clarify which file is read first for context detection.

The context detection algorithm (line 1055-1060) says it finds `workspace.yaml` first, then looks for `application.yaml` "within the workspace directory tree". For `path: .`, both files are in the same directory, so both are found immediately. This is correct, but should be explicitly mentioned as an edge case.

**Questions**:
1. Should the context detection documentation explicitly mention single-app workspace behavior?

**Suggested Resolution**: Add a note to the Context Detection Algorithm section: "For single-app workspaces (path: .), both files are in the same directory and are detected simultaneously."

**Response**:
> Add note — Explicitly document single-app workspace behavior in context detection.

---

## Checklist

- [x] Add `default_flavor: default` to PRD example (already present, added "Optional" comment)
- [x] Document override file naming convention in PRD
- [x] Standardize date format across revision histories (already consistent - no changes needed)
- [x] Clarify single-app workspace context detection (already documented in Tech Spec edge cases)


---

## Archived

This issue was archived on 2024-12-31 at 14:03:26.
