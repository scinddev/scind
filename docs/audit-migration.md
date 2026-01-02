# Documentation Audit — Migration

**For AI Agents**: Use this audit process when the documentation was migrated from existing source content.

---

## Prerequisites

- Documentation has been migrated using `install.md` from source documentation
- All migration steps have been executed
- Read the common audit instructions in the Common Instructions section below

---

## Configuration

| Setting | Value |
|---------|-------|
| **Source Directory** | `./specs` |
| **Target Directory** | `./docs` |

### Source Files

The following files comprise the source documentation:

| File | Description |
|------|-------------|
| `specs/contrail-shell-integration.md` | Shell integration specification |
| `specs/contrail-technical-spec.md` | Technical specification |
| `specs/contrail-go-stack.md` | Go stack implementation details |
| `specs/contrail-cli-reference.md` | CLI reference documentation |
| `specs/contrail-prd.md` | Product requirements document |

---

## Critical Audit Principles

**IMPORTANT**: The migrated docs use a layered structure with sub-directories. Before claiming any content is missing, you MUST:

1. **Read ALL Markdown files** under `./docs` — not just README.md files, but every `.md` file in every subdirectory including `appendices/` folders
2. **Read the full content of each file** — not just the first section or summary
3. **Search across the entire documentation directory** using grep for key terms before concluding something is missing

**Why this matters**: Migrated content is reorganized, not reduced. Content that was in one large file may now be split across:
- Main `{topic}.md` files (overview, key concepts)
- Appendix files in `appendices/{topic}/` (detailed examples, code scaffolding, complete scripts)
- Multiple specification files (one per feature instead of one monolithic doc)

A successful migration should have **approximately the same or greater total line count** as the source.

---

## Migration Audit Process

### Step 1: Inventory Source Documentation

Read ALL files in the source documentation directory (`./specs`):

```
SOURCE_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "specs/contrail-shell-integration.md", lines: N },
    { path: "specs/contrail-technical-spec.md", lines: N },
    { path: "specs/contrail-go-stack.md", lines: N },
    { path: "specs/contrail-cli-reference.md", lines: N },
    { path: "specs/contrail-prd.md", lines: N }
  ]
}
```

**Command to count lines**:
```bash
wc -l specs/*.md
```

### Step 2: Inventory Migrated Documentation

Read ALL Markdown files in the documentation directory (`./docs`):

```
MIGRATED_INVENTORY = {
  total_files: N,
  total_lines: N,
  files: [
    { path: "docs/decisions/0001-docker-compose-project-name-isolation.md", lines: N },
    { path: "docs/specs/shell-integration.md", lines: N },
    { path: "docs/specs/appendices/configuration-schemas/complete-examples.md", lines: N },
    ...
  ]
}
```

**Command to count lines**:
```bash
find docs -name "*.md" -exec wc -l {} + | sort -n
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
- Check for content in `.migration/` or `blackhole/` directories

### Step 4: Identify Potentially Missing Content

For key terms and concepts from the source, search the migrated documentation:

1. Extract major section headings from source files
2. For each heading, grep the documentation directory
3. If not found by exact match, try related terms
4. Only mark as "missing" after exhaustive search

**Key terms to search for** (based on source files):

| Source File | Key Terms to Search |
|-------------|---------------------|
| `contrail-shell-integration.md` | shell integration, bash, zsh, fish, prompt, cd hook |
| `contrail-technical-spec.md` | Docker Compose, project name, networking, overlay, Traefik |
| `contrail-go-stack.md` | Go, Cobra, Viper, implementation, package |
| `contrail-cli-reference.md` | contrail up, contrail down, CLI, commands, options |
| `contrail-prd.md` | vision, problem, solution, user stories, roadmap |

**Search command**:
```bash
grep -r "search term" docs/
```

```
CONTENT_CHECK = [
  {
    source_term: "Shell Integration",
    search_results: ["docs/specs/shell-integration.md", "docs/specs/appendices/shell-integration/..."],
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

**Check for `.migration/` directory**:
```bash
ls -la docs/.migration/ 2>/dev/null || echo "No .migration directory"
```

**Check for `blackhole/` directory**:
```bash
ls -la docs/blackhole/ 2>/dev/null || echo "No blackhole directory"
```

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

## Common Audit Instructions

### Purpose

The audit workflow:
1. Compares documentation content against expectations
2. Calculates content distribution across layers
3. Identifies content that may have been missed or misclassified
4. Reports on documentation completeness

### Core Audit Process

#### Step A: Inventory Documentation

Scan the documentation directory and categorize all content:

**Count All Content**

For the entire documentation directory:
- List all Markdown files recursively
- Count total lines per file
- Calculate total lines across all files

```
DOCS_INVENTORY = {
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

**Categorize by Layer**

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

**Identify Appendix Content**

Separately track appendix content:

```
APPENDIX_INVENTORY = {
  total_files: N,
  total_lines: N,
  by_parent: {
    "docs/specs/appendices/configuration-schemas": { files: N, lines: N },
    "docs/reference/appendices/cli": { files: N, lines: N },
    "docs/reference/appendices/configuration": { files: N, lines: N },
    ...
  }
}
```

#### Step B: Analyze Content Distribution

**Calculate Layer Distribution**

For each layer, calculate:
- Percentage of total documentation
- Main content vs appendix content ratio
- Number of cross-references to other layers

**Check for Orphaned Content**

Identify files that:
- Are not linked from any index or README
- Don't follow naming conventions
- Appear to be duplicates

---

## Report Format

```markdown
# Migration Audit Report

**Date**: {timestamp}
**Source Documentation**: ./specs
**Migrated Documentation**: ./docs

---

## Executive Summary

| Metric | Value |
|--------|-------|
| Source total lines | {N} |
| Migrated total lines | {M} |
| Difference | {+/-X} ({Y}%) |
| Migration status | {Complete/Needs Review} |

---

## Line Count Comparison

### Source Files

| File | Lines |
|------|-------|
| specs/contrail-shell-integration.md | {N} |
| specs/contrail-technical-spec.md | {N} |
| specs/contrail-go-stack.md | {N} |
| specs/contrail-cli-reference.md | {N} |
| specs/contrail-prd.md | {N} |
| **Total** | **{N}** |

### Migrated Files

| File | Lines |
|------|-------|
| docs/decisions/0001-docker-compose-project-name-isolation.md | {N} |
| docs/specs/shell-integration.md | {N} |
| docs/specs/appendices/configuration-schemas/complete-examples.md | {N} |
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

### Documentation Summary

**Total Content**: {N} files, {M} lines

### Layer Distribution

| Layer | Files | Lines | % of Total |
|-------|-------|-------|------------|
| Decisions | {N} | {M} | {X}% |
| Product | {N} | {M} | {X}% |
| Architecture | {N} | {M} | {X}% |
| Specifications | {N} | {M} | {X}% |
| Reference | {N} | {M} | {X}% |
| Implementation | {N} | {M} | {X}% |

### Appendix Usage

| Parent Document | Appendix Files | Appendix Lines |
|-----------------|----------------|----------------|
| docs/specs/appendices/configuration-schemas | {N} | {M} |
| docs/reference/appendices/cli | {N} | {M} |
| docs/reference/appendices/configuration | {N} | {M} |
| ... | ... | ... |

**Total Appendix Content**: {N} files, {M} lines ({X}% of total)

---

## Recommendations

1. **{priority}**: {recommendation}
2. ...
```

---

## Verification Checklist

After completing the audit:

- [ ] All Markdown files in `./specs` have been read in full
- [ ] All Markdown files in `./docs` have been read in full (including appendices)
- [ ] Line counts are accurate (not estimated)
- [ ] Layer categorization is complete
- [ ] Appendix content is accounted for
- [ ] Content from `.migration/` directory reviewed
- [ ] Summary report is generated
