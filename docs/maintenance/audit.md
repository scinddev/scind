# Documentation Audit

**For AI Agents**: This document contains instructions for auditing the documentation in this directory.

---

## Audit Configuration

- **Documentation Root** (`DOCS_DIR`): `docs/`
- **Legacy Documentation** (`LEGACY_DOCS_DIR`): `specs/`
- **Install Type**: Migration
- **Install Date**: 2026-01-05

---

## Critical Audit Principles

**IMPORTANT**: The migrated docs use a layered structure with sub-directories. Before claiming any content is missing, you MUST:

1. **Read ALL Markdown files** under `docs/` — not just README.md files, but every `.md` file in every subdirectory including `appendices/` folders
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

Read ALL files in `specs/`:

For each file, record:
- File path
- Total line count
- Major section headings

Calculate total source lines.

### Step 2: Inventory Migrated Documentation

Read ALL Markdown files in `docs/`:

For each file, record:
- File path
- Total line count

Calculate total migrated lines.

### Step 3: Compare Line Counts

| Metric | Value |
|--------|-------|
| Source total lines | {N} |
| Migrated total lines | {M} |
| Difference | {M - N} |
| Ratio | {M / N} |

**Expected**: Migrated lines should be approximately equal to or greater than source lines.

If migrated content is significantly less than source:
- Re-check that ALL files were read (especially in `appendices/` directories)
- Verify no directories were skipped
- Check for content in `migration/` or `blackhole/` directories if they exist

### Step 4: Verify Content Presence

For each major section heading from the source:
1. Search `docs/` for the term using grep
2. If not found by exact match, try related terms
3. Only mark as "missing" after exhaustive search confirms no matches

### Step 5: Generate Report

Produce a report that includes:

1. **Line Count Summary**
   - Total lines in source
   - Total lines in migrated docs
   - Difference and ratio

2. **Confirmed Missing Content** (verified via grep search)
   - List each item with source location
   - Note if intentionally excluded

3. **Content Restructured but Preserved**
   - List items that moved to new locations
   - Show old location → new location(s)

4. **Content Expanded or Improved**
   - List items where migrated content exceeds source
   - Note what was added

---

## Common Audit Process

### Inventory Documentation

Scan the documentation root (`DOCS_DIR`) and categorize all content:

#### Count All Content

For the entire documentation directory:
- List all Markdown files recursively
- Count total lines per file
- Calculate total lines across all files

#### Categorize by Layer

Group files by their layer:

| Layer | Directory |
|-------|-----------|
| Decisions | `decisions/` |
| Vision | `product/` |
| Architecture | `architecture/` |
| Specifications | `specs/` |
| Reference | `reference/` |
| Implementation | `implementation/` |

#### Identify Appendix Content

Separately track appendix content in each layer's `appendices/` subdirectory.

---

### Analyze Content Distribution

#### Calculate Layer Distribution

For each layer, calculate:
- Percentage of total documentation
- Main content vs appendix content ratio
- Number of cross-references to other layers

#### Check for Orphaned Content

Identify files that:
- Are not linked from any index or README
- Don't follow naming conventions
- Appear to be duplicates

---

### Generate Summary Report

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
