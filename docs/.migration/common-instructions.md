# Migration Common Instructions

**For AI Agents**: Read this file before executing any migration step file.

## Critical Migration Principle

**Migration means PRESERVING content, not summarizing it.**

When migrating content:
- **Preserve ALL technical details** — Every code example, error message, output sample, and configuration snippet
- **The migrated documents should contain the same depth of information as the source**
- **Do not abbreviate or summarize** — Move content, don't rewrite it
- **Reorganization, not reduction** — The goal is to put content in the right place, not to make it shorter

A successful migration of a section should have approximately the same line count as the source (possibly slightly higher due to added structure).

## Source Pointer Format

Content is referenced using `file:start-end` format:
- `specs/contrail-prd.md:45-62` means lines 45-62 of `specs/contrail-prd.md`
- Read the EXACT lines specified — do not read more or less

## How to Execute a Migration Step

1. Read this file (common-instructions.md)
2. Read the migration step file (e.g., `01-decisions.md`)
3. For each content block in the step file:
   a. If content is pre-extracted: Use the content directly
   b. If content is a source pointer: Read the specified lines from source
4. Create the target file with the content
5. Add the source attribution comment
6. Mark the section as complete (if tracking)

## Templates

When creating files, use the templates installed in:
- `docs/decisions/0000-template.md` for ADRs
- `docs/specs/_template.md` for specifications
- Other layer templates as applicable

## Cross-Layer Links

Do NOT add cross-layer links during individual step execution. A separate pass (09-cross-links.md) handles this after all content is migrated.

## Content Thresholds

Apply these thresholds to determine main content vs appendix:

| Threshold | Value | Action |
|-----------|-------|--------|
| `CODE_BLOCK_LINES` | 50 | Code blocks ≥ 50 lines → appendix |
| `STEP_LIST_ITEMS` | 10 | Step lists ≥ 10 items → appendix |
| `TABLE_ROWS` | 20 | Tables ≥ 20 rows → appendix |
| Complete file examples | Always | → appendix |
| Error catalogs | Always | → appendix |
| Shell scripts | Always | → appendix |
