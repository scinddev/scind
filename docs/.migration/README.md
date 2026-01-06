# Migration Plan

**Generated**: 2026-01-05
**Legacy Documentation**: `specs/`
**Documentation Root**: `docs/`

## Overview

This migration is split into 9 step files to avoid context window issues. Each step can be executed in a separate session.

## Content Summary

| Step | Layer | Files | Lines | Status |
|------|-------|-------|-------|--------|
| 01 | Decisions | 11 ADRs | ~450 lines | Pending |
| 02 | Vision | 3 files | ~400 lines | Pending |
| 03 | Comparison | 1 file | ~80 lines | Pending |
| 04 | Roadmap | 1 file | ~150 lines | Pending |
| 05 | Architecture | 1 file | ~200 lines | Pending |
| 06 | Specifications | 10 files | ~2,500 lines | Pending |
| 07 | Reference | 2 files + appendices | ~1,800 lines | Pending |
| 08 | Implementation | 1 file + appendices | ~1,600 lines | Pending |
| 09 | Cross-Links | - | - | Pending |

**Total estimated content**: ~6,100 lines (matching source)

## Execution Instructions

### For Each Step:

1. Start a new session (to have fresh context)
2. Read `common-instructions.md`
3. Read the step file (e.g., `01-decisions.md`)
4. Execute the migration instructions in the step file
5. Mark the step complete in this README
6. Proceed to next step (in same or new session as needed)

### Recommended Order

Execute steps 01-08 in order. Step 09 (cross-links) must be executed last.

Steps 01-08 can be parallelized if running multiple sessions.

## After All Steps Complete

1. Run the audit workflow using `docs/maintenance/audit.md`
2. Review any content in `migration/` directory and move to appropriate layers
3. Review any content in `blackhole/` directory
4. Delete this `.migration/` directory:
   ```bash
   rm -rf docs/.migration/
   ```

## Step Execution Log

| Step | Executed By | Date | Notes |
|------|-------------|------|-------|
| 01 | | | |
| 02 | | | |
| 03 | | | |
| 04 | | | |
| 05 | | | |
| 06 | | | |
| 07 | | | |
| 08 | | | |
| 09 | | | |
