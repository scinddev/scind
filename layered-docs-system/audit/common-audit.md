# Documentation Audit — Common Instructions

**For AI Agents**: This section contains the core audit process shared by all audit types.

**Terminology**: See the [Glossary](../LAYERED-DOCUMENTATION-SYSTEM.md#glossary) for definitions of `DOCS_DIR` and other terms.

---

## Purpose

The audit workflow:
1. Compares documentation content against expectations
2. Calculates content distribution across layers
3. Identifies content that may have been missed or misclassified
4. Reports on documentation completeness

---

## Core Audit Process

### Step 1: Inventory Documentation

Scan the documentation root (`DOCS_DIR`) and categorize all content:

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
