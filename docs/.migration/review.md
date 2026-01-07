# Layered Documentation System (LDS) - Review

**For AI Agents**: This document contains instructions for the Review phase of manifest-based extraction. Read this entire file, then follow the process below.

**Status**: Active
**Implementation Spec**: `2026-01-05-manifest-based-extraction-review`

**Terminology**: Key terms used in this document:
- `LDS_DIST_DIR`: The directory containing this distribution (e.g., `layered-docs-system/`)
- `DOCS_DIR`: The target project documentation directory (e.g., `docs/`)
- `LEGACY_DOCS_DIR`: Source documentation being migrated

For complete definitions, see the [Glossary](../DOCUMENTATION-GUIDE.md#glossary).

---

## Overview

The Review phase is the second phase of the three-phase manifest-based extraction workflow:

1. **Analyze** (`analyze.md`) - Generate extraction manifest from legacy documentation
2. **Review** (this workflow) - Human review and approval of manifest entries
3. **Extract** (`extract.md`) - Deterministic extraction based on approved manifest

This workflow enables human reviewers to inspect, adjust, and approve extraction manifest entries before content extraction occurs. The review process uses direct YAML editing with confidence-based review workflows and detailed audit logging.

**Key Principles**:
- Manifests are edited directly as YAML files (no UI required)
- All changes are tracked in a review log for audit purposes
- Entries are grouped by confidence level for efficient review
- Review ensures approximately 100% source coverage before extraction

---

## Prerequisites

Before running this workflow:

1. **Analysis phase must be complete** - Execute `analyze.md` first
2. **Manifest files must exist** in `.migration/manifests/`
3. **Main manifest must exist** at `.migration/extraction-manifest.yaml`
4. **Understand layer classifications** - Review `DOCUMENTATION-GUIDE.md` for classification heuristics

---

## Step 0: Load Configuration

Read the migration configuration from `.migration/config.yaml`:

```yaml
docs_dir: "{DOCS_DIR}"
legacy_docs_dir: "{LEGACY_DOCS_DIR}"
created_at: "{timestamp}"
```

Store these values for use throughout the review phase.

---

## Confidence-Based Review Workflow

Manifest entries are grouped into three review tracks based on their confidence level. This allows reviewers to allocate appropriate time and attention to each category.

### Review Track Summary

| Confidence | Review Approach | Time per Entry | Destination |
|------------|-----------------|----------------|-------------|
| High | Quick scan for obvious errors | Seconds | Layer directly |
| Medium | Moderate review of classification and target | 1-2 minutes | Review required |
| Low | Careful individual review | 3-5 minutes | `blackhole/` by default |

---

### High-Confidence Review Track

**Approach**: Quick verification - scan for obvious errors only.

High-confidence entries have strong heuristic matches and clear layer placement. The analysis phase is confident these classifications are correct.

**What to Look For**:
- Incorrect layer assignment (e.g., reference content classified as specification)
- Obvious template mismatch (e.g., PRD content assigned to ADR template)
- Target path errors (e.g., wrong directory or filename)
- Missing or incorrect field mappings for template fields

**Expected Time**: A few seconds per entry

**Review Actions**:
1. Scan the `target` path - does it make sense for the content?
2. Check the `classification_reason` - does it align with the content?
3. If everything looks correct, mark as reviewed
4. If issues found, escalate to medium-confidence review

**When to Escalate to Medium Review**:
- Classification reason does not match your understanding of the content
- Target path seems incorrect but you are unsure of the correct location
- Field mappings reference content you want to verify

---

### Medium-Confidence Review Track

**Approach**: Moderate review - verify classification and target path.

Medium-confidence entries have recognizable patterns but uncertain placement. These require verification that the classification is correct and the target path is appropriate.

**What to Verify**:
1. **Classification correctness** - Read the source content and confirm it belongs in the proposed layer
2. **Target path appropriateness** - Verify the filename and directory are correct
3. **Template selection** - Ensure the assigned template matches the content type
4. **Field mappings** - Check that `field_mappings` make sense for the template

**Expected Time**: 1-2 minutes per entry

**Review Actions**:
1. Read the referenced source content (use the `source` field to locate it)
2. Compare against classification heuristics from `DOCUMENTATION-GUIDE.md`
3. Verify target layer is appropriate
4. Adjust confidence to `high` if classification is confirmed correct
5. Adjust to `low` or change target if classification is incorrect
6. Document decision rationale in `review_notes`

**When to Escalate to Low Review**:
- Content does not clearly fit any layer
- Content may need to be split across multiple entries
- Significant ambiguity about correct placement

---

### Low-Confidence Review Track

**Approach**: Careful individual review - determine correct placement or confirm blackhole routing.

Low-confidence entries have weak or no heuristic matches. By default, these will be extracted to `blackhole/` for later manual processing.

**What to Determine**:
1. **Correct layer placement** - Can this content be assigned to a specific layer?
2. **Split requirement** - Should this content be split into multiple entries?
3. **Blackhole confirmation** - Is `blackhole/` the appropriate destination?
4. **Exclusion decision** - Should this content be excluded from migration entirely?

**Expected Time**: 3-5 minutes per entry

**Review Actions**:
1. Read the full source content carefully
2. Apply classification heuristics manually
3. If layer can be determined:
   - Update `target` to appropriate layer path
   - Increase `confidence` to `medium` or `high`
   - Document reasoning in `review_notes`
4. If content should be split:
   - Follow the [Entry Splitting Instructions](#entry-splitting-instructions)
5. If blackhole is correct:
   - Add `review_notes` explaining why content is unclassifiable
   - Mark as reviewed to confirm the decision
6. Add detailed `review_notes` documenting your decision rationale

---

### Entry Ordering Within Confidence Groups

Within each confidence group, entries are ordered for efficient review:

1. **By source file path** (alphabetical) - Groups related content together
2. **Within same source file, by start_line** - Follows document flow

**Rationale**: This ordering allows reviewers to maintain context while reviewing entries from the same source document. You can read through the source file once while reviewing all entries that reference it.

**Example Order**:
```
legacy-docs/architecture.md:10-50   (entry-arch-001)
legacy-docs/architecture.md:55-120  (entry-arch-002)
legacy-docs/architecture.md:125-200 (entry-arch-003)
legacy-docs/cli-reference.md:1-100  (entry-cli-001)
legacy-docs/cli-reference.md:105-300 (entry-cli-002)
legacy-docs/prd.md:1-75             (entry-prd-001)
```

---

## Source Coverage Verification

Before approving the manifest, verify that approximately 100% of source content is accounted for. This prevents content from being silently dropped during migration.

### Manual Coverage Verification Technique

For each source file in `LEGACY_DOCS_DIR`, verify coverage using simple arithmetic:

**Step 1**: Determine total lines in source file
```bash
wc -l legacy-docs/PRD.md
# Output: 500 legacy-docs/PRD.md
```

**Step 2**: Sum line ranges from manifest entries for that source
- Look at all entries where `source.file` matches the source file
- Sum up `(end_line - start_line + 1)` for each entry

**Step 3**: Calculate coverage
```
Coverage = (Sum of entry line ranges) / (Total source lines) * 100%
```

**Target**: Approximately 100% coverage (some variation acceptable due to whitespace/headers)

---

### Gap Identification Technique

Identify gaps in coverage by checking line range continuity:

**Step 1**: List all entries for a source file, sorted by `start_line`

**Step 2**: For each consecutive pair of entries, verify:
- `end_line` of entry N + 1 equals `start_line` of entry N+1
- If not, there is a gap between entries

**Step 3**: Check boundaries:
- First entry should start at line 1 (or document why content before it is excluded)
- Last entry should end at the file's last line (or document why content after it is excluded)

**Example Gap Analysis**:
```
Source: legacy-docs/PRD.md (500 lines total)

Manifest entries:
  - entry-001: lines 1-200
  - entry-002: lines 350-500

Gap detected: lines 201-349 have no manifest entry

Action required:
  - Review lines 201-349 in source
  - Either add entry for this content or document exclusion reason
```

---

### Coverage Calculation Examples

**Example 1: 100% Coverage**
```
Source file: legacy-docs/architecture.md
Total lines: 300

Manifest entries:
  - entry-arch-001: lines 1-100    = 100 lines
  - entry-arch-002: lines 101-200  = 100 lines
  - entry-arch-003: lines 201-300  = 100 lines

Total covered: 100 + 100 + 100 = 300 lines
Coverage: 300/300 = 100%
Status: COMPLETE
```

**Example 2: Gap Detected**
```
Source file: legacy-docs/PRD.md
Total lines: 500

Manifest entries:
  - entry-prd-001: lines 1-200     = 200 lines
  - entry-prd-002: lines 350-500   = 151 lines

Total covered: 200 + 151 = 351 lines
Coverage: 351/500 = 70.2%
Gap: lines 201-349 (149 lines uncovered)
Status: INCOMPLETE - review gap content
```

**Example 3: Overlapping Ranges (Error)**
```
Source file: legacy-docs/cli.md
Total lines: 400

Manifest entries:
  - entry-cli-001: lines 1-150     = 150 lines
  - entry-cli-002: lines 100-250   = 151 lines  ERROR: overlaps with entry-cli-001
  - entry-cli-003: lines 251-400   = 150 lines

Issue: Lines 100-150 are covered by two entries
Action: Adjust line ranges to be non-overlapping
```

---

### Silently Dropped Content Detection

Content can be silently dropped if a source file has no manifest entries at all.

**Detection Process**:
1. List all files in `LEGACY_DOCS_DIR`
2. For each file, search manifest `sources` section
3. Files not listed in `sources` have no entries

**What to Do with Missing Files**:
- **If intentionally excluded**: Document the exclusion in `metadata.exclusions` with reason
- **If accidentally omitted**: Re-run analysis or manually add entries

**Example Check**:
```
Files in LEGACY_DOCS_DIR:
  - PRD.md
  - architecture.md
  - cli-reference.md
  - CONTRIBUTING.md
  - old-notes.md

Files in manifest sources:
  - PRD.md
  - architecture.md
  - cli-reference.md

Missing from manifest:
  - CONTRIBUTING.md - Action: Add entries or document exclusion
  - old-notes.md - Action: Add entries or document exclusion
```

---

## Structured Review Log

All significant changes made during review must be logged in the manifest's `metadata.review_log` section. This provides an audit trail of review decisions.

### Review Log Location and Structure

**Location**: `metadata.review_log` array in the manifest file

**Purpose**: Audit trail for all review decisions and changes

**When to Add Entries**:
- Changing an entry's `target` path
- Changing an entry's `confidence` level
- Changing an entry's `mode`
- Modifying `field_mappings`
- Splitting or merging entries
- Any other significant modification

### Required Fields for Each Log Entry

Each review log entry must include these fields:

| Field | Format | Description |
|-------|--------|-------------|
| `timestamp` | ISO 8601 (e.g., `2026-01-05T16:30:00Z`) | When the change was made |
| `entry_id` | String | ID of the manifest entry being changed |
| `change` | String | Description of what was changed |
| `reason` | String | Explanation of why the change was made |

### Review Log Format Example

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  # ... other metadata ...

  review_log:
    - timestamp: "2026-01-05T16:30:00Z"
      entry_id: "adr-0002-networking"
      change: "Changed target from decisions/0002-networking.md to specs/network-design.md"
      reason: "Content describes network behavior specification, not a decision record"

    - timestamp: "2026-01-05T16:35:00Z"
      entry_id: "vision-main"
      change: "Upgraded confidence from medium to high"
      reason: "Verified content matches PRD template structure after manual review"

    - timestamp: "2026-01-05T16:40:00Z"
      entry_id: "cli-detailed-examples"
      change: "Updated field_mappings.commands to include all subcommands"
      reason: "Original mapping missed nested subcommand documentation"

    - timestamp: "2026-01-05T16:45:00Z"
      entry_id: "mixed-content-001"
      change: "Split into entries mixed-content-001a and mixed-content-001b"
      reason: "Source content contains both ADR and specification content that belong in different layers"
```

---

## Template Field Verification Reference

When reviewing entries with `template` assignments, verify that `field_mappings` cover all required template fields.

### ADR Templates

#### MADR Minimal Template (`adr-madr-minimal.md`)

| Field | Required | Description |
|-------|----------|-------------|
| `Status` | Yes | Current decision status: Draft, Proposed, Accepted, Deprecated, or Superseded |
| `Date` | Yes | Date of the decision in YYYY-MM-DD format |
| `Decision-Makers` | Yes | List of people or roles involved in the decision |
| `Context` | Yes | Description of the situation requiring a decision, including forces at play |
| `Decision` | Yes | Clear, concise statement of the decision made |
| `Consequences` | Yes | Outcomes of the decision, categorized as Positive, Negative, and Neutral |
| `Notes` | No | Additional context, clarifications, or references |

#### MADR Full Template (`adr-madr-full.md`)

Includes all MADR Minimal fields, plus:

| Field | Required | Description |
|-------|----------|-------------|
| `Technical Story` | No | Link to related issue or ticket |
| `Context and Problem Statement` | Yes | Detailed description of the situation (replaces `Context` in full template) |
| `Decision Drivers` | Yes | List of factors that influenced the decision (e.g., performance, expertise, constraints) |
| `Considered Options` | Yes | List of alternatives that were evaluated |
| `Decision Outcome` | Yes | Clear statement of the chosen option (replaces `Decision` in full template) |
| `Pros/Cons` | Yes | Detailed analysis of pros and cons for each considered option |
| `Links` | No | Related documents, issues, or external references |

---

### PRD/Vision Template

#### Lean PRD Template (`prd-lean.md`)

| Field | Required | Description |
|-------|----------|-------------|
| `Version` | Yes | Document version number (e.g., 0.1.0) |
| `Date` | Yes | Document date in YYYY-MM-DD format |
| `Status` | Yes | Document status: Draft, Active, or Deprecated |
| `Problem Statement` | Yes | Description of the problem being solved, including scenario, existing solution gaps, and unmet need |
| `Vision` | Yes | Core value proposition with key principles |
| `Core Concepts` | Yes | Definitions and explanations of key concepts |
| `Non-Goals` | Yes | Explicit list of what the product does NOT do, with reasoning |
| `Success Criteria` | Yes | Measurable outcomes that indicate product success |
| `Glossary` | No | Definitions of domain-specific terms used in the document |
| `Related Documents` | No | Links to related ADRs, specs, or external resources |

---

### Feature Spec Template

#### Feature Specification Template (`spec-feature.md`)

| Field | Required | Description |
|-------|----------|-------------|
| `Version` | Yes | Specification version number |
| `Date` | Yes | Specification date in YYYY-MM-DD format |
| `Status` | Yes | Spec status: Draft, Review, Accepted, or Implemented |
| `Overview` | Yes | One-paragraph description of the feature and its purpose |
| `Behavior` | Yes | Detailed description of normal flow and state machine (if applicable) |
| `Data Schema` | Yes | Configuration/input schemas with field tables and validation rules |
| `Examples` | Yes | Usage examples with input, behavior, and expected result |
| `Edge Cases` | Yes | Documentation of unusual scenarios and how the system handles them |
| `Error Handling` | Yes | Error conditions, codes, messages, and recovery actions |
| `Integration Points` | No | External systems, APIs, or services this feature interacts with |
| `Implementation Notes` | No | Guidance for developers implementing the feature |
| `Open Questions` | No | Unresolved issues or decisions that need to be made |
| `Revision History` | No | Record of significant changes to the specification |

---

### Field Mapping Verification Checklist

When reviewing entries with templates:

- [ ] All required fields from the template are present in `field_mappings`
- [ ] Field names in `field_mappings` match the template's expected field names
- [ ] Source references (e.g., `source:legacy-docs/PRD.md:45-60`) point to correct content
- [ ] Literal values are appropriate for the field (e.g., Status values are valid enum options)
- [ ] No template fields are left unmapped without a documented reason

**Common Issues**:
- Missing required fields (e.g., `Decision-Makers` omitted from ADR)
- Incorrect field names (e.g., `context` instead of `Context`)
- Source references pointing to wrong line ranges
- Invalid enum values (e.g., `Status: In Progress` instead of `Status: Draft`)

---

## YAML Editing Examples and Patterns

Manifests are edited directly as YAML files. The following examples show common review actions with before/after YAML and corresponding review log entries.

### Pattern 1: Target Path Adjustment

**Scenario**: Entry is classified correctly but target path needs correction.

**Before**:
```yaml
entries:
  - id: "adr-0003-database"
    source:
      file: "legacy-docs/architecture.md"
      start_line: 200
      end_line: 275
    target: "decisions/0003-database.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains 'We decided to use...' decision pattern"
```

**After**:
```yaml
entries:
  - id: "adr-0003-database"
    source:
      file: "legacy-docs/architecture.md"
      start_line: 200
      end_line: 275
    target: "decisions/0003-postgresql-database-selection.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains 'We decided to use...' decision pattern"
    reviewed: true
    reviewed_by: "reviewer@example.com"
    reviewed_at: "2026-01-05T17:00:00Z"
    review_notes: "Updated filename to be more descriptive per ADR naming convention"
```

**Review Log Entry**:
```yaml
- timestamp: "2026-01-05T17:00:00Z"
  entry_id: "adr-0003-database"
  change: "Updated target from decisions/0003-database.md to decisions/0003-postgresql-database-selection.md"
  reason: "Filename updated to be more descriptive and follow ADR naming conventions"
```

---

### Pattern 2: Confidence Level Change

**Scenario**: After verification, medium-confidence entry can be upgraded to high.

**Before**:
```yaml
entries:
  - id: "spec-auth-flow"
    source:
      file: "legacy-docs/features.md"
      start_line: 100
      end_line: 200
    target: "specs/authentication-flow.md"
    mode: "verbatim"
    confidence: "medium"
    classification_reason: "Contains behavioral descriptions but mixed with some implementation details"
```

**After**:
```yaml
entries:
  - id: "spec-auth-flow"
    source:
      file: "legacy-docs/features.md"
      start_line: 100
      end_line: 200
    target: "specs/authentication-flow.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains behavioral descriptions but mixed with some implementation details"
    reviewed: true
    reviewed_by: "reviewer@example.com"
    reviewed_at: "2026-01-05T17:15:00Z"
    review_notes: "Verified content is specification - implementation details are acceptable context"
```

**Review Log Entry**:
```yaml
- timestamp: "2026-01-05T17:15:00Z"
  entry_id: "spec-auth-flow"
  change: "Upgraded confidence from medium to high"
  reason: "Manual review confirmed content is specification; implementation details provide helpful context"
```

---

### Pattern 3: Layer Assignment Update

**Scenario**: Entry incorrectly placed in migration/ needs to be assigned to proper layer.

**Before**:
```yaml
entries:
  - id: "uncertain-content-007"
    source:
      file: "legacy-docs/technical-notes.md"
      start_line: 50
      end_line: 150
    target: "migration/uncertain/technical-notes-50-150.md"
    mode: "verbatim"
    confidence: "low"
    classification_reason: "Mixed content format, unable to determine layer"
```

**After**:
```yaml
entries:
  - id: "uncertain-content-007"
    source:
      file: "legacy-docs/technical-notes.md"
      start_line: 50
      end_line: 150
    target: "reference/api-configuration.md"
    mode: "verbatim"
    confidence: "medium"
    classification_reason: "Mixed content format, unable to determine layer"
    template: "reference-config.md"
    reviewed: true
    reviewed_by: "reviewer@example.com"
    reviewed_at: "2026-01-05T17:30:00Z"
    review_notes: "Content describes API configuration options - belongs in reference layer"
```

**Review Log Entry**:
```yaml
- timestamp: "2026-01-05T17:30:00Z"
  entry_id: "uncertain-content-007"
  change: "Changed target from migration/uncertain/technical-notes-50-150.md to reference/api-configuration.md; upgraded confidence from low to medium; added template reference-config.md"
  reason: "Content is API configuration documentation - correctly classified as reference material"
```

---

### Pattern 4: Review Notes and Reviewed Flags

**Scenario**: Marking an entry as reviewed with observations.

**Before**:
```yaml
entries:
  - id: "vision-main"
    source:
      file: "legacy-docs/PRD.md"
      start_line: 1
      end_line: 150
    target: "product/vision.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains problem statement and product vision patterns"
    template: "prd-lean.md"
```

**After**:
```yaml
entries:
  - id: "vision-main"
    source:
      file: "legacy-docs/PRD.md"
      start_line: 1
      end_line: 150
    target: "product/vision.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains problem statement and product vision patterns"
    template: "prd-lean.md"
    reviewed: true
    reviewed_by: "jane.smith@example.com"
    reviewed_at: "2026-01-05T18:00:00Z"
    review_notes: "Approved. Contains all required PRD fields. Consider adding glossary entries post-extraction."
```

**Review Log Entry**:
```yaml
- timestamp: "2026-01-05T18:00:00Z"
  entry_id: "vision-main"
  change: "Marked as reviewed with no changes"
  reason: "Content correctly classified and all PRD template fields are present"
```

---

### Copy-Pasteable YAML Snippets

#### Review Log Entry Template
```yaml
- timestamp: "YYYY-MM-DDTHH:MM:SSZ"
  entry_id: "entry-id-here"
  change: "Description of what was changed"
  reason: "Explanation of why the change was made"
```

#### Reviewed Flags Template
```yaml
reviewed: true
reviewed_by: "your-identifier@example.com"
reviewed_at: "2026-01-05T16:30:00Z"
review_notes: "Your observations and notes about this entry"
```

#### Full Entry Review Block (Add to Entry)
```yaml
    reviewed: true
    reviewed_by: "reviewer@example.com"
    reviewed_at: "2026-01-05T16:30:00Z"
    review_notes: "Entry verified correct. No changes required."
```

#### Skip Review Metadata (Add to Manifest Metadata)
```yaml
metadata:
  # ... other metadata ...
  status: "approved"
  review_skipped: true
  review_skipped_reason: "Simple migration with clear classifications - proceeding with default proposals"
```

---

## Skip Review Procedure

If you need to bypass the review phase entirely, follow this formal procedure.

### Formal Skip Procedure

To skip review and proceed directly to extraction:

**Step 1**: Add skip flags to manifest metadata:
```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"

  # Skip review flags
  status: "approved"
  review_skipped: true
  review_skipped_reason: "Brief explanation of why review is being skipped"
```

**Step 2**: Proceed to extraction phase:
```
Execute the extraction workflow: extract.md
```

### Consequences of Skipping Review

When review is skipped, entries proceed based on their analyzed confidence level:

| Confidence | Destination | What Happens |
|------------|-------------|--------------|
| High | Layer directly as analyzed | Content extracted to target paths as proposed |
| Medium | `migration/` directory | Content placed in migration folder for post-extraction review |
| Low | `blackhole/` directory | Content placed in blackhole as unclassified |

**Important**: No human verification of classifications occurred. Incorrect classifications will not be caught until after extraction.

### Warnings About Skipping Review

**Not Recommended For**:
- Complex migrations with many source files
- First-time users unfamiliar with LDS layer classifications
- Sources with mixed content types (e.g., PRD + specs in one file)
- Migrations where classification accuracy is critical

**Risks**:
- Incorrect layer assignments go undetected
- Medium-confidence content accumulates in `migration/` requiring post-extraction cleanup
- Low-confidence content ends up in `blackhole/` when it could be properly classified
- Field mapping errors are not caught until extraction fails or produces incorrect output

**Minimum Recommendation**: If skipping full review, at minimum review low-confidence entries to prevent content loss.

### Bulk Marking Procedure

To mark all entries as reviewed without individual inspection:

**Option 1: Manual Addition**
Add `reviewed: true` to each entry in the manifest files. This is tedious but explicit.

**Option 2: Script Approach**
Use a text processing tool to add reviewed flags:
```bash
# Example using sed (adjust for your YAML structure)
# This adds reviewed: true after each entry's confidence field
sed -i '' '/confidence:/a\    reviewed: true' .migration/manifests/*.yaml
```

**Warning**: Bulk marking bypasses individual verification. Use only when:
- You have high confidence in the analysis quality
- You plan to review extracted content afterward
- Time constraints prevent thorough review

---

## Entry Splitting and Merging

### Recommendation: Keep Entries As-Is

**Prefer keeping entries unchanged** when possible. Splitting and merging:
- Adds complexity to the manifest
- Can introduce errors in line range calculations
- Makes audit trail harder to follow
- Requires careful ID management

**Only modify entries when clearly necessary** - when content genuinely serves multiple layers or when distinct entries should be combined for coherence.

### When to Split vs Avoid Splitting

**Split When**:
- Content genuinely serves multiple layers (e.g., embedded decision in a specification)
- Different parts of the range require different templates
- Parts of the content have different confidence levels that you want to track separately

**Avoid Splitting When**:
- Content is cohesive and belongs together conceptually
- Splitting would create fragments that are not useful standalone
- The complexity of managing multiple entries exceeds the benefit
- You are unsure - when in doubt, keep entries as-is and note the complexity in `review_notes`

### Entry Splitting Instructions

When content in a single source range genuinely serves multiple layers (e.g., contains both a decision and a specification), split into separate entries.

**Step 1: Duplicate the Original Entry**

Copy the entire entry in the manifest.

**Step 2: Assign New Unique IDs**

Both entries need unique IDs. Use a suffix pattern:

```yaml
# Original ID: mixed-content-001
# New IDs: mixed-content-001a, mixed-content-001b
```

**Step 3: Adjust Line Ranges**

Make line ranges non-overlapping:

**Before (single entry)**:
```yaml
- id: "mixed-content-001"
  source:
    file: "legacy-docs/technical.md"
    start_line: 100
    end_line: 200
  target: "migration/uncertain/technical-100-200.md"
  confidence: "low"
```

**After (split entries)**:
```yaml
- id: "mixed-content-001a"
  source:
    file: "legacy-docs/technical.md"
    start_line: 100
    end_line: 145
  target: "decisions/0004-caching-strategy.md"
  mode: "verbatim"
  confidence: "medium"
  classification_reason: "Decision content about caching approach"
  reviewed: true
  reviewed_by: "reviewer@example.com"
  reviewed_at: "2026-01-05T19:00:00Z"
  review_notes: "Split from mixed-content-001 - lines 100-145 contain ADR content"

- id: "mixed-content-001b"
  source:
    file: "legacy-docs/technical.md"
    start_line: 146
    end_line: 200
  target: "specs/cache-behavior.md"
  mode: "verbatim"
  confidence: "medium"
  classification_reason: "Specification content about cache behavior"
  reviewed: true
  reviewed_by: "reviewer@example.com"
  reviewed_at: "2026-01-05T19:00:00Z"
  review_notes: "Split from mixed-content-001 - lines 146-200 contain spec content"
```

**Step 4: Update Targets**

Assign appropriate target paths for each split entry.

**Step 5: Add Review Log Entries**

```yaml
- timestamp: "2026-01-05T19:00:00Z"
  entry_id: "mixed-content-001"
  change: "Split into mixed-content-001a (decisions layer) and mixed-content-001b (specs layer)"
  reason: "Source content contains both decision rationale and behavioral specification - belong in separate layers"
```

### Entry Merging Instructions

When multiple entries should produce a single output document (e.g., vision content spread across multiple source files), merge them.

**Step 1: Choose Primary Entry**

Select which entry to keep as the primary. Usually the one with the most content or clearest target path.

**Step 2: Update Source to Multi-File Format**

**Before (separate entries)**:
```yaml
- id: "vision-part-1"
  source:
    file: "legacy-docs/PRD.md"
    start_line: 1
    end_line: 100
  target: "product/vision.md"
  confidence: "high"

- id: "vision-part-2"
  source:
    file: "legacy-docs/roadmap.md"
    start_line: 15
    end_line: 80
  target: "product/vision-roadmap.md"
  confidence: "medium"
```

**After (merged entry)**:
```yaml
- id: "vision-combined"
  source:
    files:
      - file: "legacy-docs/PRD.md"
        start_line: 1
        end_line: 100
      - file: "legacy-docs/roadmap.md"
        start_line: 15
        end_line: 80
    separator: "\n\n---\n\n"
  target: "product/vision.md"
  mode: "verbatim"
  confidence: "high"
  classification_reason: "Combined vision content from PRD and roadmap"
  reviewed: true
  reviewed_by: "reviewer@example.com"
  reviewed_at: "2026-01-05T19:30:00Z"
  review_notes: "Merged vision-part-1 and vision-part-2 into single vision document"
```

**Step 3: Remove Secondary Entry**

Delete the entry that was merged (vision-part-2 in this example).

**Step 4: Update ID and Target**

Assign a descriptive ID and appropriate target path.

**Step 5: Add Review Log Entry**

```yaml
- timestamp: "2026-01-05T19:30:00Z"
  entry_id: "vision-combined"
  change: "Created by merging vision-part-1 and vision-part-2"
  reason: "Both entries contain vision content for the same document - combined for coherent output"
```

---

## Approval Process

After reviewing all entries (or explicitly skipping review), complete the approval process to authorize extraction.

### Required Metadata Updates

Add approval fields to the manifest metadata:

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"

  # Approval fields
  status: "approved"
  approved_by: "lead-reviewer@example.com"
  approved_at: "2026-01-05T20:00:00Z"
```

### Review Log Completion

Before approval, verify the `review_log` section contains all significant changes made during review:

- All target path changes should be logged
- All confidence level changes should be logged
- All split/merge operations should be logged
- Final approval can be noted as a log entry

**Example Final Log Entry**:
```yaml
- timestamp: "2026-01-05T20:00:00Z"
  entry_id: "manifest"
  change: "Manifest approved for extraction"
  reason: "All entries reviewed; coverage verified at 98%; 3 low-confidence entries confirmed for blackhole"
```

### Version Control Recommendation

**Commit the approved manifest to version control before proceeding to extraction.**

Benefits:
- Allows rollback if extraction issues occur
- Provides audit trail of approved state
- Enables comparison between analysis and reviewed versions
- Preserves record of review decisions

```bash
git add .migration/
git commit -m "Approve extraction manifest after review

- Reviewed N entries across M source files
- Coverage verified at X%
- N entries adjusted during review
- Ready for extraction phase"
```

### Pre-Approval Checklist

Before marking the manifest as approved, verify:

- [ ] **All entries reviewed** - Every entry has `reviewed: true` or review was explicitly skipped with documentation
- [ ] **Coverage verification complete** - Source coverage is approximately 100% (or gaps are documented)
- [ ] **Template field mappings verified** - Required fields are mapped for entries with templates
- [ ] **Review log complete** - All significant changes during review are logged
- [ ] **Low-confidence entries resolved** - Each low-confidence entry is either:
  - Upgraded to medium/high with correct target
  - Confirmed for blackhole with documented reason
- [ ] **No overlapping line ranges** - No source content is covered by multiple entries
- [ ] **No duplicate targets** - No two entries write to the same target file

---

## Next Phase

After review completes and the manifest is approved, output the next step command:

```
REVIEW PHASE COMPLETE
=====================

Manifest has been reviewed and approved.

Status: {metadata.status}
Approved by: {metadata.approved_by}
Approved at: {metadata.approved_at}

Entries reviewed:
  - High confidence: {N} entries
  - Medium confidence: {N} entries
  - Low confidence: {N} entries

┌─────────────────────────────────────────────────────────────────────┐
│ NEXT STEP                                                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ Run /clear, then execute the extraction phase:                      │
│                                                                     │
│   Execute @docs/.migration/extract.md                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

The extraction phase will:
  - Verify source file hashes match manifest
  - Generate migration step files
  - Output a command list for running each step
```

---

## Related Documentation

- `install.md` - Main installation workflow (orchestrates all phases)
- `analyze.md` - Analysis phase workflow (generates manifests)
- `extract.md` - Extraction phase workflow (executes approved manifest)
- `schemas/extraction-manifest.schema.yaml` - Manifest schema definition
- `DOCUMENTATION-GUIDE.md` - Classification heuristics and layer definitions
