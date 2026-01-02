# Documentation Audit — Fresh Install

**For AI Agents**: Use this audit process when the documentation was created fresh (not migrated from existing content).

---

## Prerequisites

- Documentation has been installed using `install.md` with no source migration
- Read the common audit instructions in `common-audit.md` (or the Common Instructions section below)

---

## Fresh Install Audit Process

For fresh installs, the audit focuses on:
1. Verifying the structure was created correctly
2. Checking that templates are in place
3. Ensuring indexes link to all content

### Step 1: Verify Structure

Check that the expected directories exist based on layer selections during install:

```
Expected directories (based on install selections):
- [ ] decisions/ (if Layer 1 selected)
- [ ] product/ (if Layer 2 selected)
- [ ] architecture/ (if Layer 3 selected)
- [ ] specs/ (if Layer 4 selected)
- [ ] reference/ (if Layer 5 selected)
- [ ] implementation/ (if Layer 7 selected)
- [ ] DOCUMENTATION-GUIDE.md (always)
```

### Step 2: Verify Templates

Check that template files exist:

```
Expected templates:
- [ ] decisions/0000-template.md (if Layer 1)
- [ ] specs/_template/README.md (if Layer 4)
- [ ] Other templates based on selections
```

### Step 3: Verify Indexes

Check that README.md index files exist in each layer directory and contain appropriate links.

### Step 4: Run Common Audit

Execute the common audit process to inventory all content.

---

## Report Format

```markdown
# Fresh Install Audit Report

**Date**: {timestamp}
**Documentation Directory**: {path}

## Structure Verification

| Directory | Expected | Found | Status |
|-----------|----------|-------|--------|
| decisions/ | Yes | Yes/No | ✓/✗ |
| product/ | Yes | Yes/No | ✓/✗ |
| ... | ... | ... | ... |

## Template Verification

| Template | Expected | Found | Status |
|----------|----------|-------|--------|
| decisions/0000-template.md | Yes | Yes/No | ✓/✗ |
| ... | ... | ... | ... |

## Content Summary

{Include common audit summary here}

## Status

- Structure: {Complete/Incomplete}
- Templates: {Complete/Incomplete}
- Ready for use: {Yes/No}
```
