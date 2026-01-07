# Documentation Audit — Migration

**For AI Agents**: Use this audit process when the documentation was migrated from legacy documentation.

**Terminology**: See the [Glossary](../DOCUMENTATION-GUIDE.md#glossary) for definitions of `DOCS_DIR`, `LEGACY_DOCS_DIR`, and other terms.

---

## Audit Configuration

- **Documentation Root** (`DOCS_DIR`): `docs/`
- **Legacy Documentation** (`LEGACY_DOCS_DIR`): `specs/`
- **Install Type**: Migration
- **Install Date**: 2026-01-06

---

## Prerequisites

- Documentation has been migrated using `install.md` from legacy documentation
- All migration steps have been executed
- Read the common audit instructions in the Common Instructions section below

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

Read ALL files in the legacy documentation directory (`specs/`):

```
LEGACY_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "specs/file1.md", lines: N },
    { path: "specs/file2.md", lines: N },
    ...
  ]
}
```

### Step 2: Inventory Migrated Documentation

Read ALL Markdown files in the documentation root (`docs/`):

```
MIGRATED_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "docs/decisions/0001-example.md", lines: N },
    { path: "docs/specs/feature.md", lines: N },
    { path: "docs/specs/appendices/feature/details.md", lines: N },
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
**Legacy Documentation** (`LEGACY_DOCS_DIR`): specs/
**Documentation Root** (`DOCS_DIR`): docs/

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
| specs/file1.md | {N} |
| specs/file2.md | {N} |
| ... | ... |
| **Total** | **{N}** |

### Migrated Files

| File | Lines |
|------|-------|
| docs/decisions/0001-example.md | {N} |
| docs/specs/feature.md | {N} |
| docs/specs/appendices/feature/details.md | {N} |
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

---

## Common Audit Process

### Purpose

The audit workflow:
1. Compares documentation content against expectations
2. Calculates content distribution across layers
3. Identifies content that may have been missed or misclassified
4. Reports on documentation completeness

---

### Step 1: Inventory Documentation

Scan the documentation root (`docs/`) and categorize all content:

#### 1a: Count All Content

For the entire documentation directory:
- List all Markdown files recursively
- Count total lines per file
- Calculate total lines across all files

```
DOCS_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "decisions/0001-example.md", lines: N },
    { path: "specs/feature.md", lines: N },
    { path: "specs/appendices/feature/details.md", lines: N },
    ...
  ]
}
```

#### 1b: Categorize by Layer

Group files by their layer:

```
LAYER_INVENTORY = {
  "decisions": {
    documents: N,
    total_lines: N,
    files: [...]
  },
  "product": {
    documents: N,
    total_lines: N,
    files: [...]
  },
  "architecture": { ... },
  "specs": { ... },
  "reference": { ... },
  "implementation": { ... }
}
```

#### 1c: Identify Appendix Content

Separately track appendix content:

```
APPENDIX_INVENTORY = {
  total_files: N,
  total_lines: N,
  by_parent: {
    "specs/appendices/feature-name": { files: N, lines: N },
    "reference/appendices/cli": { files: N, lines: N },
    ...
  }
}
```

---

### Step 2: Analyze Content Distribution

#### 2a: Calculate Layer Distribution

For each layer, calculate:
- Percentage of total documentation
- Main content vs appendix content ratio
- Number of cross-references to other layers

#### 2b: Check for Orphaned Content

Identify files that:
- Are not linked from any index or README
- Don't follow naming conventions
- Appear to be duplicates

---

### Step 3: Generate Summary Report

Create a summary of the documentation state:

```markdown
## Documentation Summary

**Total Content**: {N} files, {M} lines

### Layer Distribution

| Layer | Files | Lines | % of Total |
|-------|-------|-------|------------|
| Decisions | {N} | {M} | {X}% |
| Vision | {N} | {M} | {X}% |
| Architecture | {N} | {M} | {X}% |
| Specifications | {N} | {M} | {X}% |
| Reference | {N} | {M} | {X}% |
| Implementation | {N} | {M} | {X}% |

### Appendix Usage

| Parent Document | Appendix Files | Appendix Lines |
|-----------------|----------------|----------------|
| {path} | {N} | {M} |
| ... | ... | ... |

**Total Appendix Content**: {N} files, {M} lines ({X}% of total)
```

---

## Verification Checklist

After completing the audit:

- [ ] All Markdown files have been read in full
- [ ] Line counts are accurate (not estimated)
- [ ] Layer categorization is complete
- [ ] Appendix content is accounted for
- [ ] Summary report is generated
