# Specification Refinement Workflow

The specification refinement workflow maintains consistency across authoritative specs and documents for a project.

## Proper Context For ALL Operations

When working on **anything**, always read:

- @README.md – project overview and issue summary
- @specs – all specifications and documents
- @issues/00-index.md – issue tracking and status

## Trigger Phrases Reference

| User Says | Action |
|-----------|--------|
| "research refinement" or "find issues" | Start [Refinement Research](#refinement-research) |
| "start refining" or "what's next?" | Start [Refinement Session](#refinement-session) |
| "close the refinement session" or "done for now" | Run [Closing a Refinement Session](#closing-a-refinement-session) |
| "archive issues" | Run [Archive Issues](#archive-issues) |

## Refinement Research

1. Review [Known Issues and Non-Issues](#known-issues-and-non-issues) to avoid flagging accepted items
2. Research all files in `@specs` for contradictions, functionality gaps, or ambiguous assumptions
3. Group issues by which spec documents they impact (e.g., "Tech Spec + Go Stack", "PRD + CLI Reference")
4. Create issue files in `issues/` following the [Issue File Format](#issue-file-format)
5. Determine recommended resolution order (foundational issues first)
6. Update `@README.md` with newly added issues
7. Update `@issues/00-index.md` with newly added issues
8. Ask user if they want to start a refinement session

### Issue File Format

Each issue group file (`issues/NN-title.md`) should follow this structure:

```markdown
# Issue Group N: [Title]

**Documents Affected**: [List]
**Suggested Order**: N of M ([rationale])
**Estimated Effort**: Small | Medium | Large

---

## Overview

[Brief description of the group's theme]

---

## Issues

### [ID]: [Title]

**Severity**: High | Medium | Low

**Issue**: [Description with spec line references]

**Questions**: (if user input needed)
1. [Question]

**Suggested Resolution**: [Recommendation]

**Response**:
> _[Your response here]_

---

## Checklist

- [ ] [Action item]
```

**Issue IDs follow these prefixes:**
- `C-N`: Contradictions between documents
- `M-N`: Missing functionality or documentation
- `A-N`: Ambiguous or unclear specifications
- `X-N`: Typos or minor fixes
- `N-N`: New issues discovered in subsequent review passes (after initial research)
- `L-N`: Low-priority improvements (optional enhancements, not blocking)

**Numbering**: Issue numbers (the `N` in each prefix) are sequential within each prefix category and reset to 1 after an archive operation.

## Refinement Session

1. Load [Proper Context for ALL Operations](#proper-context-for-all-operations)
2. Read `@issues/00-index.md` for issue overview and completion states
3. Select the first incomplete issue **group** from the recommended order
4. Read the selected issue group file
5. For each unfilled response (`_[Your response here]_`) in the issue group:
   - Display the question/options to the user
   - Wait for their response
   - Update the issue file with their response as a blockquote
   - If all responses are already filled, proceed directly to step 6
6. Apply resolutions to affected spec documents
7. Update the issue group's checklist items as completed
8. Mark the issue group as completed in `@issues/00-index.md` by adding `✅ COMPLETED` to the Rationale column
9. Ask: "Continue to next issue group or close session?"
   - **Continue**: Return to step 2
   - **Close**: Run [Closing a Refinement Session](#closing-a-refinement-session)

## Closing a Refinement Session

This process handles **spec version updates only**. Issue group completion status is managed during the session itself (step 8 of Refinement Session).

1. Check if any spec documents in `specs/` were modified during this session
2. If **no specs were modified**: No version updates needed. Skip to step 4.
3. If **specs were modified**, for each modified spec:
   - Increment the patch version (e.g., 0.5.0 → 0.5.1)
   - Add an entry to the Revision History table
4. Summarize the session:
   - Number of issue groups completed
   - Number of issue groups remaining
   - Specs modified (if any)

**Note**: Incomplete issue groups remain open and will appear in future refinement sessions. To defer an issue without resolving it, mark the issue group complete with a response explaining the deferral (e.g., "Deferred: requires further discussion").

## Archive Issues

Start the archive process if a user asks to "archive issues." Then, follow this process:

### Step 1: Check for Incomplete Issues

1. Read `@issues/00-index.md`
2. Identify any issues NOT marked as `✅ COMPLETED`
3. If incomplete issues exist:
   - List the incomplete issues to the user
   - Ask: "These issues are not complete. Do you want to archive them anyway?"
   - If **No**: End the process and tell the user to ask again once they have either completed all open issues OR are ready to archive incomplete issues
   - If **Yes**: Continue to Step 2

### Step 2: Create Archive Directory

1. Generate a timestamp in format `YYYYMMDDHHmmss` (e.g., `20251231134552` for December 31st, 2025 at 1:45:52 PM)
2. Create directory `archive/{timestamp}/`

### Step 3: Move and Update Issue Files

1. Move all files from `issues/` (including `00-index.md`) to `archive/{timestamp}/`
2. For each moved file, append an archive footer:

```markdown

---

## Archived

This issue was archived on YYYY-MM-DD at HH:MM:SS.
```

### Step 4: Reset README.md

Update `@README.md` to remove issue entries but preserve the structure:

1. Keep the directory structure section but update it to show only:
   ```
   project-spec-review/
   ├── README.md
   ├── specs/                          # Original specifications and documents
   │   └── ...
   ├── issues/                         # Review findings organized by group
   │   └── 00-index.md                 # Start here - overview and recommended order
   └── archive/                        # Archived issue sets
       └── {timestamp}/                # Previous review cycle
   ```
2. Clear the issue file list from the directory structure
3. Reset the summary section to show 0 issues
4. Clear the Recommended Order table but keep headers

### Step 5: Create Fresh Index

Create a new `@issues/00-index.md` with empty tables (the old one was moved in Step 3):

```markdown
# Refinement Issues — Index

**Created**: [Current Month Year]
**Total Issues**: 0
**Groups**: 0

---

## Recommended Order

| Order | File | Issues | Documents | Effort | Rationale |
|-------|------|--------|-----------|--------|-----------|

---

## By Severity

### High Severity (0 issues)
| Issue | Group | File |
|-------|-------|------|

### Medium Severity (0 issues)
| Issue | Group | File |
|-------|-------|------|

### Low Severity (0 issues)
| Issue | Group | File |
|-------|-------|------|

---

## By Document Impact

_No issues currently tracked._

---

## Quick Start

1. Run refinement research to identify new issues
2. Work through issues in recommended order
3. Archive when complete

---

## Files in This Directory

```
issues/
└── 00-index.md                    # This file
```
```

### Step 6: Confirm Completion

1. Report to the user:
   - How many issue files were archived
   - The archive location (`archive/{timestamp}/`)
   - Confirmation that `issues/` and `README.md` are reset
2. Remind user that:
   - All issue numbering (C-N, M-N, etc.) resets to 1 for the next cycle
   - They can run "research refinement" to start a new review cycle

## Version Management

Each spec document in `specs/` maintains its own independent version:

1. **Increment versions independently** — each document has its own version
2. **Update at session end** — increment patch version for each modified document
3. **Update revision history** — add entry to Revision History table

Example revision history entry:
```markdown
| 0.5.1-draft | Dec 2024 | Clarified validation timing, added error examples |
```

**Note**: Version mismatches across documents are intentional—each evolves independently.

## Known Issues and Non-Issues

### Non-Issues

Items that should NOT be flagged during reviews:

- **Spec Version Mismatch** — versions intentionally differ per document (see [Version Management](#version-management))

### Known Issues (Deferred)

_None currently tracked._
