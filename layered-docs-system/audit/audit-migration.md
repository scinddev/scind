# Documentation Audit — Migration

**For AI Agents**: Use this audit process when the documentation was migrated from legacy documentation.

**Terminology**: See the [Glossary](../LAYERED-DOCUMENTATION-SYSTEM.md#glossary) for definitions of `DOCS_DIR`, `LEGACY_DOCS_DIR`, and other terms.

---

## Prerequisites

- Documentation has been migrated using `install.md` from legacy documentation
- All migration steps have been executed
- Read the common audit instructions in `common-audit.md` (or the Common Instructions section below)

---

## Critical Audit Principles

**IMPORTANT**: The migrated docs use a layered structure with sub-directories. Before claiming any content is missing, you MUST:

1. **Read ALL Markdown files** under the documentation directory — not just README.md files, but every `.md` file in every subdirectory including `appendices/` folders
2. **Read the full content of each file** — not just the first section or summary
3. **Search across the entire documentation directory** using grep for key terms before concluding something is missing

**Why this matters**: Migrated content is reorganized, not reduced. Content that was in one large file may now be split across:
- Main `{topic}.md` files (overview, key concepts)
- Appendix files in `appendices/{topic}/` (detailed examples, code scaffolding, complete scripts)
- Multiple specification files (one per feature instead of one monolithic doc)

A successful migration should have **approximately the same or greater total line count** as the source.

---

## Migration Audit Process

### Step 1: Inventory Legacy Documentation

Read ALL files in the legacy documentation directory (`LEGACY_DOCS_DIR`):

```
LEGACY_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "{LEGACY_DOCS_DIR}/file1.md", lines: N },
    { path: "{LEGACY_DOCS_DIR}/file2.md", lines: N },
    ...
  ]
}
```

### Step 2: Inventory Migrated Documentation

Read ALL Markdown files in the documentation root (`DOCS_DIR`):

```
MIGRATED_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "{DOCS_DIR}/decisions/0001-example.md", lines: N },
    { path: "{DOCS_DIR}/specs/feature.md", lines: N },
    { path: "{DOCS_DIR}/specs/appendices/feature/details.md", lines: N },
    ...
  ]
}
```

### Step 3: Compare Line Counts

Calculate the overall content comparison:

```
COMPARISON = {
  source_total_lines: N,
  migrated_total_lines: M,
  difference: M - N,
  ratio: M / N  // Should be ~1.0 or higher
}
```

**Expected outcome**: `migrated_total_lines >= source_total_lines`

If migrated content is significantly less than source:
- Re-check that ALL files were read (especially appendices)
- Verify no directories were skipped
- Check for content in `migration/` or `blackhole/` directories

### Step 4: Identify Potentially Missing Content

For key terms and concepts from the source, search the migrated documentation:

1. Extract major section headings from source files
2. For each heading, grep the documentation directory
3. If not found by exact match, try related terms
4. Only mark as "missing" after exhaustive search

```
CONTENT_CHECK = [
  {
    source_term: "Shell Integration",
    search_results: ["specs/shell-integration.md", "specs/appendices/shell-integration/..."],
    status: "found"
  },
  {
    source_term: "Revision History",
    search_results: [],
    status: "missing"  // Only after grep confirms no matches
  },
  ...
]
```

### Step 5: Check Migration and Blackhole Directories

If these directories exist, inventory their contents:

```
PENDING_CONTENT = {
  migration: {
    files: N,
    lines: M,
    items: [...]
  },
  blackhole: {
    files: N,
    lines: M,
    items: [...]
  }
}
```

### Step 6: Run Common Audit

Execute the common audit process for full content inventory.

---

## Report Format

```markdown
# Migration Audit Report

**Date**: {timestamp}
**Legacy Documentation** (`LEGACY_DOCS_DIR`): {legacy_docs_path}
**Documentation Root** (`DOCS_DIR`): {docs_path}

---

## Executive Summary

| Metric | Value |
|--------|-------|
| Legacy total lines | {N} |
| Migrated total lines | {M} |
| Difference | {+/-X} ({Y}%) |
| Migration status | {Complete/Needs Review} |

---

## Line Count Comparison

### Legacy Files

| File | Lines |
|------|-------|
| {LEGACY_DOCS_DIR}/file1.md | {N} |
| {LEGACY_DOCS_DIR}/file2.md | {N} |
| ... | ... |
| **Total** | **{N}** |

### Migrated Files

| File | Lines |
|------|-------|
| {DOCS_DIR}/decisions/0001-example.md | {N} |
| {DOCS_DIR}/specs/feature.md | {N} |
| {DOCS_DIR}/specs/appendices/feature/details.md | {N} |
| ... | ... |
| **Total** | **{M}** |

---

## Content Verification

### Confirmed Present (restructured but preserved)

| Source Content | New Location(s) |
|----------------|-----------------|
| {description} | {path(s)} |
| ... | ... |

### Confirmed Missing

| Source Content | Source Location | Notes |
|----------------|-----------------|-------|
| {description} | {file:lines} | {why it may be intentionally excluded, or action needed} |
| ... | ... | ... |

### Expanded or Improved

| Content | Original Lines | New Lines | Notes |
|---------|----------------|-----------|-------|
| {description} | {N} | {M} | {what was added} |
| ... | ... | ... | ... |

---

## Pending Review

### Migration Directory

| File | Lines | Suggested Destination |
|------|-------|-----------------------|
| {path} | {N} | {destination} |
| ... | ... | ... |

### Blackhole Directory

| File | Lines | Reason | Suggested Action |
|------|-------|--------|------------------|
| {path} | {N} | {reason} | {action} |
| ... | ... | ... | ... |

---

## Layer Distribution

{Include common audit layer summary here}

---

## Recommendations

1. **{priority}**: {recommendation}
2. ...
```
