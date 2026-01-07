# Layered Documentation System (LDS) - Extract

**For AI Agents**: This document contains instructions for the Extraction phase of manifest-based extraction. Read this entire file, then follow the process below.

**Terminology**: Key terms used in this document:
- `DOCS_DIR`: The target project documentation directory (e.g., `docs/`)
- `LEGACY_DOCS_DIR`: Source documentation being migrated
- `Manifest`: YAML file defining extraction entries with source references, targets, and modes

For complete definitions, see the [Glossary](../DOCUMENTATION-GUIDE.md#glossary).

---

## Overview

The Extraction phase is the third and final phase of the three-phase manifest-based extraction workflow:

1. **Analyze** (`analyze.md`) - Generate extraction manifest from legacy documentation
2. **Review** (`review.md`) - Human review and approval of manifest entries
3. **Extract** (this workflow) - Deterministic extraction based on approved manifest

This workflow performs deterministic content extraction based on reviewed and approved manifests. The core principle is simple: **the same manifest will always produce identical output**.

---

## Prerequisites

Before running this workflow, verify the following conditions are met:

1. **Manifest status must be "approved"**
   - The manifest `metadata.status` field must equal `"approved"`
   - If not approved, this workflow will refuse to proceed and direct you to `review.md`

2. **Analysis phase must be complete**
   - Manifest files must exist in `.migration/manifests/`
   - Main manifest must exist at `.migration/extraction-manifest.yaml`

3. **Review phase should be complete**
   - Entries should have `reviewed: true` flags (or review explicitly skipped)
   - For manifest approval process, see `review.md`

4. **Source files must be accessible**
   - All source files referenced in the manifest must exist
   - Source file integrity will be verified against stored hashes

5. **Directory structure must exist**
   - Target directories must be created (typically done by `install.md` Step 7)
   - This includes layer directories like `decisions/`, `product/`, `architecture/`, etc.

---

## Step 0: Load Configuration

Read the migration configuration from `.migration/config.yaml`:

```yaml
docs_dir: "{DOCS_DIR}"
legacy_docs_dir: "{LEGACY_DOCS_DIR}"
created_at: "{timestamp}"
```

Store these values for use throughout the extraction phase.

---

## Critical Principles

### Extraction is DETERMINISTIC

**Running extraction with the same manifest produces identical output every time.**

This is the core value proposition of manifest-based extraction:

- **No heuristic decisions during extraction** - All classification decisions were made during analysis and review phases
- **No AI interpretation of content** - The manifest is the complete specification
- **Manifest is the single source of truth** - What the manifest says is what gets extracted
- **Reproducible migrations** - The same manifest run at different times yields the same result

### Extraction is MECHANICAL

The extraction phase is purely mechanical execution. It does not:

- Make classification decisions
- Interpret content meaning
- Adjust line ranges based on content
- Skip or modify entries based on content analysis

If an entry is in the manifest and the manifest is approved, it gets extracted exactly as specified.

---

## Workflow Steps

### Step 1: Validate Manifest Status

Before any extraction proceeds, verify the manifest is approved.

#### 1.1 Load Main Manifest

Read the main manifest at `.migration/extraction-manifest.yaml`:

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  status: "approved"  # <-- MUST be "approved" to proceed
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"
```

#### 1.2 Check Status Field

**If `metadata.status` equals `"approved"`**: Proceed to Step 2.

**If `metadata.status` does NOT equal `"approved"`**: Stop with error message:

```
EXTRACTION CANNOT PROCEED
=========================

The manifest has not been approved for extraction.

Current status: "{current_status}"
Required status: "approved"

The manifest must be reviewed and approved before extraction can proceed.

To approve the manifest:
1. Run the review workflow: review.md
2. Review all manifest entries
3. Set metadata.status to "approved"
4. Return to this workflow

Extraction aborted. No files have been created or modified.
```

#### 1.3 Load Layer-Specific Manifests

After verifying approval, load all layer-specific manifests from `.migration/manifests/`:

| Manifest File | Layer |
|---------------|-------|
| `01-decisions-manifest.yaml` | Decisions (ADRs) |
| `02-vision-manifest.yaml` | Vision/Product |
| `03-architecture-manifest.yaml` | Architecture |
| `04-specs-manifest.yaml` | Specifications |
| `05-reference-manifest.yaml` | Reference |
| `06-behaviors-manifest.yaml` | Behaviors |
| `07-implementation-manifest.yaml` | Implementation |

#### 1.4 Validate Schema Version

Check that all manifests have compatible `schema_version` values:

- **Compatible**: Schema version `1.x.x` (any 1.x version is compatible)
- **Incompatible**: Schema version `2.x.x` or higher requires updated extract.md

If incompatible schema version detected:

```
SCHEMA VERSION INCOMPATIBLE
===========================

Manifest schema version: {manifest_version}
Supported versions: 1.x.x

This version of extract.md does not support schema version {manifest_version}.
Please update to a compatible version of the LDS tooling.

Extraction aborted.
```

---

### Step 2: Verify Source File Hashes

Source hash verification ensures the manifest is still accurate before extraction begins. This step must complete **before any files are written**.

#### 2.1 Build Source File List

Collect all unique source files referenced in:

1. The `sources` registry (if present in manifest)
2. Entry `source` specifications
3. Multi-file entry file lists

Example manifest `sources` section:

```yaml
sources:
  "legacy-docs/PRD.md":
    hash: "a1b2c3d4e5f6789..."
    hash_algorithm: "sha256"
    size_bytes: 15234
  "legacy-docs/architecture.md":
    hash: "b2c3d4e5f6g7890..."
    hash_algorithm: "sha256"
    size_bytes: 8567
```

#### 2.2 Calculate Current Hashes

For each source file:

1. **Check if file exists** - Missing file is a critical error
2. **Calculate current hash** using the stored algorithm (default: SHA-256)
3. **Compare with stored hash** in manifest
4. **Record result**: `MATCH`, `MISMATCH`, or `MISSING`

**Hash Calculation Method**:

- **Whole file hash**: Read entire file as bytes, compute SHA-256 digest
- **Line range hash**: Extract specified lines (1-indexed, inclusive), join with newlines, encode as UTF-8, compute SHA-256
- **Combined hash**: For multi-file sources, concatenate content with separators, then hash

#### 2.3 Display Verification Results

**When all sources match**:

```
SOURCE FILE VERIFICATION
========================

Checking source files against manifest hashes...

File: legacy-docs/PRD.md
  Manifest hash:  a1b2c3d4e5f6789...
  Current hash:   a1b2c3d4e5f6789...
  Status:         MATCH

File: legacy-docs/architecture.md
  Manifest hash:  b2c3d4e5f6g7890...
  Current hash:   b2c3d4e5f6g7890...
  Status:         MATCH

File: legacy-docs/CLI.md
  Manifest hash:  c3d4e5f6g7h8901...
  Current hash:   c3d4e5f6g7h8901...
  Status:         MATCH

----------------------------------------
Verification complete: 3/3 sources match
Proceeding with extraction...
```

**When mismatches are detected**:

```
SOURCE FILE VERIFICATION
========================

Checking source files against manifest hashes...

File: legacy-docs/PRD.md
  Manifest hash:  a1b2c3d4e5f6789...
  Current hash:   a1b2c3d4e5f6789...
  Status:         MATCH

File: legacy-docs/architecture.md
  Manifest hash:  b2c3d4e5f6g7890...
  Current hash:   x9y8z7w6v5u4321...
  Status:         MISMATCH - file content has changed

File: legacy-docs/CLI.md
  Status:         MISSING - file not found

----------------------------------------
Verification FAILED: 1/3 match, 1 mismatch, 1 missing
```

#### 2.4 Handle Hash Mismatches

**CRITICAL**: When mismatch is detected, **NEVER silently proceed**.

Present the user with three options:

```
SOURCE FILE MISMATCH DETECTED

The following source files have changed since the manifest was generated:

File: legacy-docs/architecture.md
  Manifest hash:  b2c3d4e5f6g7890...
  Current hash:   x9y8z7w6v5u4321...
  Status:         CONTENT CHANGED

File: legacy-docs/CLI.md
  Status:         FILE MISSING

Extraction cannot proceed with outdated manifest. The manifest was
generated on 2026-01-05T14:30:00Z based on the source files as they
existed at that time. Since then, the files have been modified.

Choose an option:
[1] STOP - Abort extraction (no files will be written)
[2] CONTINUE - Proceed with current manifest despite changes (risky)
[3] RE-ANALYZE - Re-run analysis phase to regenerate manifest

Your choice:
```

**Option 1: Stop Extraction Entirely**

- Aborts the extraction process immediately
- No output files are created
- Manifest and source files remain unchanged
- User can investigate what changed and retry later

**Option 2: Continue Anyway**

- Proceeds with extraction using existing manifest despite changes
- Actual content extracted comes from current (changed) source files
- **Risks**:
  - Extracted content may differ from what was reviewed
  - Line ranges may no longer align with intended content
  - Classification decisions were based on previous file state
- Warning logged in extraction report
- Affected entries flagged in output attribution
- **Manual verification strongly recommended**

Warning displayed when selecting Option 2:

```
WARNING: PROCEEDING WITH CHANGED SOURCES

You have chosen to continue extraction despite source file changes.

RISKS:
- Extracted content may differ from what was reviewed
- Line ranges may no longer align with intended content
- Classification decisions were based on previous file state

The extraction report will flag all entries affected by changed sources.
Manual verification of extracted content is STRONGLY RECOMMENDED.

Proceeding with extraction...
```

**Option 3: Re-Run Analysis Phase**

- Aborts extraction without writing output files
- User directed to run `analyze.md` to regenerate manifest
- New manifest will reflect current source file state
- Review phase will need to be repeated

Message displayed when selecting Option 3:

```
RE-ANALYSIS RECOMMENDED

Source files have changed significantly since manifest generation.
To ensure accurate extraction, the analysis phase should be re-run.

Steps to proceed:
1. Run 'analyze.md' to regenerate the manifest
2. Review the new manifest entries via 'review.md'
3. Return to 'extract.md' to complete extraction

Note: Previous manifest will be overwritten. If you need to preserve
the previous manifest, back it up before re-running analysis.

Returning to analysis phase...
```

#### 2.5 The `--force` Flag

For automated/scripted scenarios, the `--force` flag auto-selects Option 2 (continue anyway):

**When `--force` is specified**:

- Hash verification still runs
- If mismatches detected, warning is displayed but extraction proceeds automatically
- Warning is logged in extraction report
- Use case: CI/CD pipelines, automated migration scripts

**Warning when `--force` is used with mismatches**:

```
SOURCE FILE MISMATCH DETECTED (--force flag set)

The following source files have changed since the manifest was generated:

File: legacy-docs/architecture.md
  Status: CONTENT CHANGED

WARNING: Proceeding automatically due to --force flag.
Extraction may produce unexpected results.

This warning will be recorded in the extraction report.
Proceeding with extraction...
```

---

### Step 3: Process Extraction Entries

Process each manifest entry in sequential order by entry ID.

#### 3.1 Extraction Mode: Verbatim (Default)

Verbatim mode extracts exact lines from source files without modification.

**Manifest Example**:

```yaml
- id: "vision-main"
  source:
    file: "legacy-docs/PRD.md"
    start_line: 45
    end_line: 120
  target: "product/vision.md"
  mode: "verbatim"
  confidence: "high"
  classification_reason: "Contains problem statement and product vision patterns"
```

**Extraction Process**:

1. Read lines 45-120 (1-indexed, inclusive) from `legacy-docs/PRD.md`
2. Write content unchanged to `{DOCS_DIR}/product/vision.md`
3. Add source attribution comments

**This is the default mode** - if no `mode` field is specified, use verbatim extraction.

#### 3.2 Extraction Mode: Content

Content mode uses the `inline_content` field directly without reading any source file.

**Manifest Example**:

```yaml
- id: "readme-decisions"
  target: "decisions/README.md"
  mode: "content"
  confidence: "high"
  classification_reason: "Auto-generated layer index"
  inline_content: |
    # Architectural Decision Records

    This directory contains the project's ADRs.

    ## Index

    | ADR | Title | Status |
    |-----|-------|--------|
    | [0001](../decisions/0001-docker-compose-project-name-isolation.md) | Docker Compose Orchestration | Accepted |

    ## Adding New ADRs

    See [DOCUMENTATION-GUIDE.md](../DOCUMENTATION-GUIDE.md) for instructions.
```

**Extraction Process**:

1. Use the `inline_content` field directly
2. Write content to `{DOCS_DIR}/decisions/README.md`
3. Add extraction ID attribution (no source attribution since content is inline)

**No source file reference is required** for content mode entries.

#### 3.3 Extraction Mode: Transform

Transform mode applies built-in transformations to source content before writing.

**Built-in Transform Types**:

| Transform Type | Description | Options |
|---------------|-------------|---------|
| `markdown_format` | Format/clean markdown content | None |
| `heading_adjust` | Adjust heading levels | `level_offset`, `skip_first` |
| `link_rewrite` | Rewrite internal documentation links | `base_path` |
| `custom` | Custom transformation | OUT OF SCOPE |

**Heading Adjust Example**:

```yaml
- id: "cli-reference-transformed"
  source:
    file: "legacy-docs/CLI.md"
  target: "reference/cli.md"
  mode: "transform"
  confidence: "high"
  classification_reason: "CLI documentation needing heading level adjustment"
  transform_config:
    type: "heading_adjust"
    options:
      level_offset: 1    # Add 1 to all heading levels (# -> ##)
      skip_first: true   # Do not adjust the first heading
```

**Heading Adjust Behavior**:

- `level_offset: 1` transforms `# Heading` to `## Heading`, `## Sub` to `### Sub`, etc.
- `level_offset: -1` transforms `## Heading` to `# Heading` (reduces levels)
- `skip_first: true` leaves the first heading unchanged (useful for document titles)

**Link Rewrite Example**:

```yaml
- id: "spec-with-links"
  source:
    file: "legacy-docs/features.md"
    start_line: 100
    end_line: 250
  target: "specs/workspace-lifecycle.md"
  mode: "transform"
  confidence: "high"
  classification_reason: "Feature spec with internal links needing rewrite"
  transform_config:
    type: "link_rewrite"
    options:
      base_path: "../"
```

**Link Rewrite Behavior**:

- Rewrites internal documentation links to match new path structure
- Preserves external links unchanged
- Updates relative paths based on `base_path` option

**Custom Transform Type**:

The `custom` transform type is **OUT OF SCOPE** for this implementation. If an entry specifies `transform_config.type: "custom"`:

1. Log an error for the entry
2. Skip the entry (do not extract)
3. Record in extraction report as skipped
4. Continue processing remaining entries

```
WARNING: Custom transform type not supported

Entry: custom-transform-entry
Transform type: custom
Status: SKIPPED - custom transforms not implemented

This entry will not be extracted. To extract this content,
change the transform type to a supported type (markdown_format,
heading_adjust, link_rewrite) or use verbatim mode.
```

#### 3.4 Extraction Mode: Summarize

Summarize mode generates a summary of source content. **This mode is NEVER automatic** - it requires explicit user instruction during the review phase.

**Manifest Example**:

```yaml
- id: "arch-overview-summary"
  source:
    file: "legacy-docs/detailed-architecture.md"
  target: "architecture/overview.md"
  mode: "summarize"
  confidence: "medium"
  classification_reason: "User requested summary of 500-line architecture doc"
  review_notes: "USER EXPLICITLY REQUESTED SUMMARY - original at legacy-docs/detailed-architecture.md"
  reviewed: true
  reviewed_by: "user@example.com"
  reviewed_at: "2026-01-05T15:00:00Z"
```

**Review Marker Validation**:

Entries with `mode: summarize` should have:
- `reviewed: true`
- `reviewed_by: "{identifier}"`

**If review markers are missing**, issue a warning (do NOT fail extraction):

```
WARNING: Summarize entry without review markers

Entry: arch-overview-summary
Mode: summarize
reviewed: false (or missing)
reviewed_by: (missing)

Summarize mode should only be used when explicitly requested by a user
during the review phase. This entry lacks the expected review markers.

Proceeding with summarization, but please verify this was intentional.
```

**Summarization Process**:

1. Read source content from specified file/lines
2. Generate summary preserving key technical details
3. Write summary to target file
4. Add source attribution with note about summarization

#### 3.5 Extraction Mode Comparison

| Mode | Description | Source Required | When to Use |
|------|-------------|-----------------|-------------|
| `verbatim` | Extract exact lines without modification | Yes | Default - preserves content exactly |
| `content` | Use inline content directly | No (uses `inline_content`) | Auto-generated content, indexes |
| `transform` | Apply built-in transformation | Yes | Heading adjustments, link rewrites |
| `summarize` | Generate summary (requires explicit approval) | Yes | Only when user explicitly requests |

---

### Step 4: Apply Templates

When an entry specifies a `template`, load and apply it during extraction.

#### 4.1 Load Templates

Templates are loaded from layer-specific locations:

| Layer | Template Location |
|-------|-------------------|
| Decisions | `{DOCS_DIR}/decisions/0000-template.md` |
| Specs | `{DOCS_DIR}/specs/_template.md` |
| Reference | `{DOCS_DIR}/reference/_template-cli.md`, `_template-config.md` |
| Other | As specified in `templates/` directory |

**Manifest Example with Template**:

```yaml
- id: "adr-0001-docker"
  source:
    file: "legacy-docs/architecture.md"
    start_line: 45
    end_line: 95
  target: "decisions/0001-docker-compose-orchestration.md"
  mode: "verbatim"
  confidence: "high"
  classification_reason: "Contains decision pattern with alternatives considered"
  template: "adr-madr-minimal.md"
  field_mappings:
    title: "Docker Compose Orchestration"
    status: "Accepted"
    date: "2026-01-05"
    context: "source:legacy-docs/architecture.md:45-60"
    decision: "source:legacy-docs/architecture.md:62-75"
```

#### 4.2 Perform Variable Substitution

Templates contain placeholders that are filled using `field_mappings`:

**Template Placeholders** (example):

```markdown
# {title}

**Status**: {status}
**Date**: {date}

## Context

{context}

## Decision

{decision}

## Consequences

{consequences}
```

**Field Mapping Types**:

1. **Static values**: Direct string values
   ```yaml
   field_mappings:
     title: "Docker Compose Orchestration"
     status: "Accepted"
   ```

2. **Source references**: Content from source files
   ```yaml
   field_mappings:
     context: "source:legacy-docs/architecture.md:45-60"
     decision: "source:legacy-docs/architecture.md:62-75"
   ```

   Format: `source:{file}:{start_line}-{end_line}`

#### 4.3 Handle Missing Template Fields

If a template field is not provided in `field_mappings`, insert "UNDEFINED" marker text:

**Before** (template):
```markdown
## Consequences

{consequences}
```

**After** (if `consequences` not in field_mappings):
```markdown
## Consequences

UNDEFINED
```

This makes missing fields visible for manual completion rather than silently omitting them.

---

### Step 5: Handle Multi-File Sources

Some entries combine content from multiple source files.

#### 5.1 Multi-File Source Structure

**Manifest Example**:

```yaml
- id: "vision-combined"
  source:
    files:
      - file: "legacy-docs/PRD.md"
        start_line: 1
        end_line: 100
        hash: "a1b2c3..."
      - file: "legacy-docs/roadmap.md"
        start_line: 15
        end_line: 80
        hash: "d4e5f6..."
    separator: "\n\n---\n\n"
    combined_hash: "g7h8i9..."
  target: "product/vision.md"
  mode: "verbatim"
  confidence: "high"
  classification_reason: "PRD contains vision, roadmap contains future direction"
```

#### 5.2 Multi-File Extraction Process

1. **Read each file in order** as specified in `source.files` list
2. **Extract specified line ranges** (or entire file if no range specified)
3. **Join content** using `separator` field (default: `\n\n`)
4. **Treat concatenated result** as a single content block
5. **Verify `combined_hash`** if provided in manifest

**Separator Behavior**:

- Default separator: `\n\n` (blank line between files)
- Custom separators can include markdown (e.g., `\n\n---\n\n` for horizontal rule)

#### 5.3 Combined Hash Verification

If `combined_hash` is provided:

1. Calculate hash of concatenated content (after joining with separator)
2. Compare with stored `combined_hash`
3. If mismatch, treat as source change (same options as file-level mismatch)

---

### Step 6: Execute Extraction

Process entries and create output files.

#### 6.1 Sequential Processing

Process entries in order by entry ID. For each entry:

1. Verify source content is accessible
2. Extract content according to mode
3. Apply template if specified
4. Create target file
5. Add attribution comments

#### 6.2 Create Target Files

Write extracted content to target paths relative to `DOCS_DIR`:

- Entry `target: "decisions/0001-docker-compose.md"`
- Full path: `{DOCS_DIR}/decisions/0001-docker-compose.md`

Ensure parent directories exist before writing.

#### 6.3 Add Source Attribution

Every extracted file includes source attribution comments at the top:

```markdown
<!-- Migrated from legacy-docs/PRD.md:45-120 -->
<!-- Extraction ID: vision-main -->

# Vision Document

[Extracted content here...]
```

**Attribution Format**:

```
<!-- Migrated from {source_file}:{start_line}-{end_line} -->
<!-- Extraction ID: {entry_id} -->
```

For multi-file sources:

```markdown
<!-- Migrated from multiple sources:
     - legacy-docs/PRD.md:1-100
     - legacy-docs/roadmap.md:15-80
-->
<!-- Extraction ID: vision-combined -->
```

For content mode (no source file):

```markdown
<!-- Generated content - no source file -->
<!-- Extraction ID: readme-decisions -->
```

#### 6.4 The `--resume` Mode

The `--resume` flag enables recovery from partial extraction failures.

**When `--resume` is specified**:

1. Before processing each entry, check if target file already exists
2. If target exists, skip the entry (do not overwrite)
3. Continue with remaining entries
4. Report shows which entries were skipped vs processed

**Use Case**: If extraction fails after processing 25 of 47 entries, use `--resume` to complete the remaining 22 entries without re-extracting the first 25.

**Resume Output Example**:

```
EXTRACTION PROGRESS (--resume mode)
===================================

Entry: adr-0001-docker
Target: decisions/0001-docker-compose.md
Status: SKIPPED (file exists)

Entry: adr-0002-networking
Target: decisions/0002-networking.md
Status: SKIPPED (file exists)

Entry: vision-main
Target: product/vision.md
Status: EXTRACTED (75 lines)

...
```

---

### Step 7: Generate Extraction Report

Generate a comprehensive report even on partial failure.

#### 7.1 Report Contents

The extraction report includes:

**Metadata**:
- Manifest file reference
- Manifest hash (for verification)
- Execution timestamp
- Executor info (agent/user identifier)

**Summary**:
- Total entries in manifest
- Successful extractions count
- Failed extractions count
- Skipped entries count (if `--resume` used)

**Line Statistics**:
- Source lines processed
- Output lines generated
- Variance percentage

**Understanding Line Statistics Variance**:

Line counts between source and output are expected to differ. This is normal and does not indicate content loss. Common causes of variance:

| Source of Variance | Direction | Typical Impact |
|-------------------|-----------|----------------|
| Template boilerplate | + Output | +10-30 lines per templated file |
| Source attribution comments | + Output | +2-3 lines per file |
| Whitespace normalization | +/- | Variable, typically small |
| Removed duplicate content | - Output | Depends on duplication amount |
| Appendix restructuring | Neutral | Lines moved, not lost |

**Acceptable variance threshold**: +/-10% is normal. Variance >15% should trigger manual review to verify no content was unintentionally omitted.

**Important**: Line statistics measure *extraction completeness*, not *content preservation*. Use the audit workflow (`maintenance/audit.md`) to verify actual content coverage.

**Layer Breakdown**:
- Files created per layer
- Total lines per layer

**File List**:
- Each created file with:
  - Source file(s) and line range(s)
  - Output line count
  - Template used (if any)

**Verification Results**:
- Line count check status
- File existence check status
- Content hash check (for verbatim extractions)

#### 7.2 Report Format (YAML)

Save report to `{DOCS_DIR}/extraction-report.yaml`:

```yaml
extraction_report:
  metadata:
    manifest_file: ".migration/extraction-manifest.yaml"
    manifest_hash: "a1b2c3d4e5f6789..."
    executed_at: "2026-01-05T16:00:00Z"
    executed_by: "claude-agent"
    flags:
      force: false
      resume: false

  summary:
    total_entries: 47
    successful: 45
    failed: 1
    skipped: 1

  line_statistics:
    source_lines_processed: 2450
    output_lines_generated: 2512
    variance_percentage: 2.5

  layer_breakdown:
    decisions:
      files_created: 5
      total_lines: 312
    product:
      files_created: 2
      total_lines: 245
    architecture:
      files_created: 1
      total_lines: 156
    specs:
      files_created: 8
      total_lines: 890
    reference:
      files_created: 3
      total_lines: 524
    implementation:
      files_created: 2
      total_lines: 385

  files_created:
    - target: "decisions/0001-docker-compose.md"
      source: "legacy-docs/architecture.md:45-95"
      source_lines: 51
      output_lines: 75
      template: "adr-madr-minimal.md"
      status: "success"

    - target: "decisions/0002-networking.md"
      source: "legacy-docs/architecture.md:100-145"
      source_lines: 46
      output_lines: 62
      template: "adr-madr-minimal.md"
      status: "success"

    - target: "product/vision.md"
      source: "legacy-docs/PRD.md:1-150"
      source_lines: 150
      output_lines: 148
      template: null
      status: "success"

    - target: "specs/custom-feature.md"
      source: "legacy-docs/features.md:200-350"
      source_lines: 151
      output_lines: 0
      template: null
      status: "failed"
      error: "Custom transform type not supported"

  verification_results:
    line_count_check:
      status: "pass"
      source_total: 2450
      output_total: 2512
      threshold_percent: 10

    file_existence_check:
      status: "pass"
      expected_files: 21
      created_files: 20
      missing_files:
        - "specs/custom-feature.md"

    content_hash_check:
      status: "pass"
      verbatim_entries: 35
      hash_verified: 35
      hash_failed: 0

  warnings:
    - entry: "arch-overview-summary"
      type: "summarize_without_review_markers"
      message: "Summarize mode entry lacks reviewed_by field"

  errors:
    - entry: "custom-feature"
      type: "unsupported_transform"
      message: "Custom transform type not supported - entry skipped"

  source_changes:
    detected: false
    files_changed: []
    force_flag_used: false
```

#### 7.3 Report on Partial Failure

The report is generated even when some entries fail:

- Failed entries are listed with error details
- Successful entries are recorded normally
- Summary reflects actual success/failure counts
- User can review failures and retry with `--resume`

#### 7.4 Human-Readable Summary

Also output a human-readable summary to the terminal:

```
EXTRACTION REPORT
=================

Date: 2026-01-05T16:00:00Z
Manifest: .migration/extraction-manifest.yaml

FILES CREATED
-------------
decisions/0001-docker-compose.md           (75 lines)
decisions/0002-networking.md               (62 lines)
decisions/0003-api-design.md               (58 lines)
product/vision.md                          (148 lines)
product/comparison.md                      (95 lines)
architecture/overview.md                   (156 lines)
specs/workspace-lifecycle.md               (245 lines)
specs/networking.md                        (180 lines)
reference/cli.md                           (312 lines)
reference/configuration.md                 (212 lines)
implementation/tech-stack.md               (185 lines)
...

LAYER SUMMARY
-------------
Decisions:      5 files,   312 lines
Product:        2 files,   245 lines
Architecture:   1 file,    156 lines
Specifications: 8 files,   890 lines
Reference:      3 files,   524 lines
Implementation: 2 files,   385 lines

SUMMARY
-------
Total entries processed: 47
Files created: 45
Errors: 1 (custom-feature - unsupported transform type)
Warnings: 1

Source lines: 2,450
Output lines: 2,512 (+2.5%)

VERIFICATION
------------
Line count check:     PASS (within 10% threshold)
File existence check: PASS (20/21 files - 1 failed entry)
Content hash check:   PASS (35/35 verbatim entries verified)

Report saved to: docs/extraction-report.yaml

Status: COMPLETED WITH ERRORS
```

---

### Step 8: Output Next Steps

After generating the extraction report, output a summary of the next steps with copy-paste commands.

#### 8.1 Generate Command List

Based on the migration step files that were created, output a command list for the extraction phase:

```
EXTRACTION PHASE COMPLETE
=========================

Migration step files have been generated. Run each command below in a
separate session, with /clear between each.

Copy and paste each command, one at a time:

┌─────────────────────────────────────────────────────────────────────┐
│ EXTRACTION COMMANDS                                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ Execute @docs/.migration/01-decisions.md                            │
│                                                                     │
│ Execute @docs/.migration/02-vision.md                               │
│                                                                     │
│ Execute @docs/.migration/03-architecture.md                         │
│                                                                     │
│ Execute @docs/.migration/04-specs.md                                │
│                                                                     │
│ Execute @docs/.migration/05-reference.md                            │
│                                                                     │
│ Execute @docs/.migration/06-implementation.md                       │
│                                                                     │
│ Execute @docs/.migration/07-medium-confidence.md                    │
│                                                                     │
│ Execute @docs/.migration/08-low-confidence.md                       │
│                                                                     │
│ Execute @docs/.migration/09-cross-links.md                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

IMPORTANT: Run /clear between each command to free context.

After all extraction steps complete, run the audit:

  Execute @docs/maintenance/audit.md
```

#### 8.2 Skip Empty Step Files

Only include commands for step files that have manifest entries. If a layer has no content to extract, omit it from the list:

```
┌─────────────────────────────────────────────────────────────────────┐
│ EXTRACTION COMMANDS (4 steps - some layers empty)                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ Execute @docs/.migration/01-decisions.md                            │
│                                                                     │
│ Execute @docs/.migration/04-specs.md                                │
│                                                                     │
│ Execute @docs/.migration/05-reference.md                            │
│                                                                     │
│ Execute @docs/.migration/09-cross-links.md                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

Skipped (no content): 02-vision, 03-architecture, 06-implementation,
                      07-medium-confidence, 08-low-confidence
```

#### 8.3 Include Estimated Sizes

When known, include estimated line counts to help users anticipate session size:

```
┌─────────────────────────────────────────────────────────────────────┐
│ EXTRACTION COMMANDS                                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ Execute @docs/.migration/01-decisions.md          (~320 lines)      │
│                                                                     │
│ Execute @docs/.migration/02-vision.md             (~150 lines)      │
│                                                                     │
│ Execute @docs/.migration/05-reference.md          (~980 lines)      │
│                                                                     │
│ Execute @docs/.migration/06-implementation.md     (~1200 lines)     │
│                                                                     │
│ Execute @docs/.migration/09-cross-links.md                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Execution Order

Execute layer steps 01-08 in order. The cross-links step (`09-cross-links.md`) **MUST execute after all other layer steps complete** because it reads all generated files to identify linking opportunities.

---

## Out of Scope

The following are explicitly **NOT part of the extraction phase**:

1. **Custom transform type implementation**
   - If encountered, error and skip the entry
   - May be implemented in future versions

2. **Cross-layer linking**
   - Handled by separate `09-cross-links.md` step after all layers extracted
   - Not performed during individual layer extraction

3. **Manifest editing or modification**
   - Manifest is read-only during extraction
   - All edits happen during review phase

4. **Re-running analysis phase**
   - Extraction only directs user to run `analyze.md` when needed
   - Does not invoke analysis directly

5. **UI for extraction progress**
   - Terminal output only
   - No web-based or graphical interface

6. **Automatic re-analysis when source changes**
   - Requires explicit user action
   - Changes are detected and reported, not automatically handled

---

## Related Documentation

| Document | Purpose |
|----------|---------|
| `install.md` | Main installation workflow (orchestrates all phases) |
| `analyze.md` | Analysis phase workflow (generates manifests) |
| `review.md` | Review phase workflow (approves manifests) |
| `schemas/extraction-manifest.schema.yaml` | Manifest schema definition with examples |
| `DOCUMENTATION-GUIDE.md` | Classification heuristics and layer definitions |

---

## Quick Reference

### Manifest Status Check

```yaml
metadata:
  status: "approved"  # Must be "approved" to proceed
```

### Extraction Modes

| Mode | Source Required | Description |
|------|-----------------|-------------|
| `verbatim` | Yes | Extract exact lines (default) |
| `content` | No | Use inline content |
| `transform` | Yes | Apply transformation |
| `summarize` | Yes | Generate summary (explicit approval required) |

### Built-in Transform Types

| Type | Description |
|------|-------------|
| `markdown_format` | Format markdown content |
| `heading_adjust` | Adjust heading levels (`level_offset`, `skip_first`) |
| `link_rewrite` | Rewrite internal links |
| `custom` | **Not supported** - will error and skip |

### Command Flags

| Flag | Effect |
|------|--------|
| `--force` | Auto-continue on source mismatch |
| `--resume` | Skip existing target files |

### Attribution Comment Format

```markdown
<!-- Migrated from {source_file}:{start_line}-{end_line} -->
<!-- Extraction ID: {entry_id} -->
```

### Report Location

```
{DOCS_DIR}/extraction-report.yaml
```
