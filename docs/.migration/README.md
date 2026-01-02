# Migration Plan

**Generated**: 2026-01-02
**Source**: `specs/`
**Target**: `docs/`

## Overview

This migration is split into 9 step files to avoid context window issues. Each step can be executed in a separate session.

## Content Summary

| Step | Layer | Files | Est. Lines | Status |
|------|-------|-------|------------|--------|
| [01](./01-decisions.md) | Decisions | 11 ADRs | ~550 lines | ⬜ Pending |
| [02](./02-vision.md) | Vision | 1 file | ~200 lines | ⬜ Pending |
| [03](./03-comparison.md) | Comparison | 1 file | ~100 lines | ⬜ Pending |
| [04](./04-roadmap.md) | Roadmap | 1 file | ~100 lines | ⬜ Pending |
| [05](./05-architecture.md) | Architecture | 1 file | ~200 lines | ⬜ Pending |
| [06](./06-specifications.md) | Specifications | 10 files | ~2,500 lines | ⬜ Pending |
| [07](./07-reference.md) | Reference | 2 files + appendices | ~1,700 lines | ⬜ Pending |
| [08](./08-implementation.md) | Implementation | 1 file + appendices | ~1,700 lines | ⬜ Pending |
| [09](./09-cross-links.md) | Cross-Links | — | — | ⬜ Pending |

## Step Files

- [common-instructions.md](./common-instructions.md) — Read this first (migration principles)
- [01-decisions.md](./01-decisions.md) — 11 ADRs with pre-extracted content
- [02-vision.md](./02-vision.md) — Product vision PRD-Lite
- [03-comparison.md](./03-comparison.md) — Alternative tools comparison
- [04-roadmap.md](./04-roadmap.md) — Future considerations
- [05-architecture.md](./05-architecture.md) — Architecture overview with diagrams
- [06-specifications.md](./06-specifications.md) — 10 specification files
- [07-reference.md](./07-reference.md) — CLI and configuration reference
- [08-implementation.md](./08-implementation.md) — Go tech stack documentation
- [09-cross-links.md](./09-cross-links.md) — Final cross-reference verification

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

Steps 01-08 can be executed in any order (or in parallel across multiple sessions), but dependencies may exist if content references other layers.

## After All Steps Complete

1. Run the audit workflow (see `docs/audit.md`)
2. Review content in `migration/` directory if it exists and move to appropriate layers
3. Review content in `blackhole/` directory if it exists
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

