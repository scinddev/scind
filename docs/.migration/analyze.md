# Layered Documentation System (LDS) - Analyze

**For AI Agents**: This document contains instructions for the Analysis phase of manifest-based extraction. Read this entire file, then follow the process below.

**Status**: COMPLETE

**Terminology**: Key terms used in this document:
- `LDS_DIST_DIR`: The directory containing this distribution (e.g., `layered-docs-system/`)
- `DOCS_DIR`: The target project documentation directory (e.g., `docs/`)
- `LEGACY_DOCS_DIR`: Source documentation being migrated

For complete definitions, see the [Glossary](../DOCUMENTATION-GUIDE.md#glossary).

---

## Overview

The Analysis phase is the first phase of the three-phase manifest-based extraction workflow:

1. **Analyze** (this workflow) - Generate extraction manifest from legacy documentation
2. **Review** (`review.md`) - Human review and approval of manifest entries
3. **Extract** (`extract.md`) - Deterministic extraction based on approved manifest

This workflow analyzes legacy documentation and generates an extraction manifest that defines exactly what content will be extracted and where it will be placed.

---

## Prerequisites

Before running this workflow:

1. The setup phase of `install.md` must have been completed, which creates:
   - `.migration/config.yaml` with path configuration
   - `.migration/manifests/` directory for output

2. Read `DOCUMENTATION-GUIDE.md` for classification heuristics

---

## Step 0: Load Configuration

Read the migration configuration from `.migration/config.yaml`:

```yaml
docs_dir: "{DOCS_DIR}"
legacy_docs_dir: "{LEGACY_DOCS_DIR}"
created_at: "{timestamp}"
```

Store these values for use throughout the analysis phase.

---

## Step 1: Read All Legacy Files

Read ALL files in `LEGACY_DOCS_DIR`. For each file:
- Record the filename and full content
- Note the file type (Markdown, YAML, etc.)
- Calculate file hash for manifest

Present a summary to the user:

```
LEGACY DOCUMENTATION LOADED
===========================

Found {N} files in {LEGACY_DOCS_DIR}:

  filename1.md          ({X} lines) - {first heading or summary}
  filename2.md          ({X} lines) - {first heading or summary}
  subdirectory/file.md  ({X} lines) - {first heading or summary}
  ...

Total: {N} files, approximately {X} lines
```

Proceed to section boundary detection and classification.

---

## Section Boundary Detection

Before classifying content, identify extractable sections using these boundary detection methods (in priority order):

### Primary Boundaries: Markdown Headers (H1-H6)

Headers establish the primary section structure:
- Each heading starts a new section
- Heading level establishes parent-child relationships (H2 under H1, H3 under H2, etc.)
- Section ends at the next heading of equal or higher level
- **Edge cases**:
  - Duplicate headings: Treat as separate sections (use line numbers to distinguish)
  - Skipped levels (H1 to H3): Treat H3 as direct child of H1

### Secondary Boundaries

When headers are absent or additional granularity is needed:

| Boundary Type | Pattern | Significance |
|--------------|---------|--------------|
| Horizontal rules | `---`, `***`, `___` | Major section dividers between topics |
| Multiple blank lines | 2+ consecutive blank lines | Major topic transitions |
| Semantic markers | "Note:", "Warning:", "Example:", "See also:" | Subsection boundaries within a section |
| List restarts | Numbered list restarts after prose | Section transitions within content |

### Code Block Boundaries

Code blocks may form self-contained extractable units:

| Pattern | Treatment |
|---------|-----------|
| Large code blocks (50+ lines) | Self-contained extractable unit, typically appendix content |
| Sequential code blocks with file paths | Scaffold section (group as single unit for appendix) |
| Shebang (`#!/`) at start | Complete file example, always appendix |
| Package declaration (`package main`, etc.) | Complete file example, always appendix |

### File/Path Reference Anchoring

File and path references act as anchors for content:
- Content mentioning a file path often describes that file
- Group content by the file/path it references
- Use path mentions to determine section relevance

---

## Layer Classification Heuristics

Apply these heuristics to classify content into the seven LDS layers. Use the Classification Decision Tree as primary routing logic.

### Classification Decision Tree

Apply these questions in sequence to classify content:

```
Is this explaining WHY a choice was made?
├─ YES → Layer 1: Decisions (ADR)
└─ NO ↓

Is this about product vision, goals, or concepts?
├─ YES → Layer 2: Vision (PRD-Lite)
└─ NO ↓

Is this showing how components relate (diagrams, topology)?
├─ YES → Layer 3: Architecture
└─ NO ↓

Is this detailing HOW a feature works (states, flows, schemas)?
├─ YES → Layer 4: Specifications
└─ NO ↓

Is this a lookup table (commands, options, codes)?
├─ YES → Layer 5: Reference
└─ NO ↓

Is this a concrete scenario that should be verified?
├─ YES → Layer 6: Behaviors
└─ NO ↓

Is this implementation scaffolding, code templates, or dependency lists?
├─ YES → Layer 7: Implementation Guides
└─ NO → Consider: may not need documentation, or route to blackhole/
```

### Layer 1: Decisions (ADR) Heuristics

**Signals that indicate ADR content**:
- "We chose X over Y because..."
- "We decided to..." / "Decision:"
- "This pattern will be used throughout..."
- "We considered alternatives including..."
- Trade-off discussions with rationale
- Explicit "Decision" blocks or sections

**Counter-examples (NOT ADR content)**:
- "Here's how to configure X..." → Reference
- "The system does X when Y happens..." → Specification

### Layer 2: Vision (PRD-Lite) Heuristics

**Signals that indicate Vision content**:
- "The problem we're solving is..."
- "Success looks like..."
- "This product is NOT for..." (non-goals)
- "A [concept] is defined as..." (glossary)
- Problem statements, product vision, mission statements
- Success criteria or goals
- User personas or target audience

**Comparison sub-document signals**:
- Comparison tables with other tools/products
- "vs" language ("X vs Y", "compared to")
- Feature matrices comparing multiple products
- "Why choose X over Y" discussions

**Roadmap sub-document signals**:
- "Future Considerations" sections
- "Roadmap" or "Future Work" headings
- "Phase 2", "Phase 3", "v2.0" planning content
- "Not yet implemented" with intent to implement

**Counter-examples (NOT Vision content)**:
- "The architecture uses..." → Architecture
- "When the user runs command X..." → Specification/Reference

### Layer 3: Architecture Heuristics

**Signals that indicate Architecture content**:
- "The system is composed of..."
- "Component A communicates with B via..."
- "Security is handled by..." (patterns)
- "This diagram shows..."
- System diagrams (ASCII art, Mermaid, etc.)
- Network topology descriptions
- Cross-cutting concerns (security, logging patterns)

**Counter-examples (NOT Architecture content)**:
- "Run this command to..." → Reference
- "When X happens, the system does Y..." → Specification

### Layer 4: Specifications Heuristics

**Signals that indicate Specification content**:
- "When X happens, the system does Y..."
- "The valid states are..."
- "Here's how feature X works..."
- "Edge case: if Y, then Z..."
- State machines or lifecycle descriptions
- Data schemas with field-level detail
- Error handling descriptions

**Counter-examples (NOT Specification content)**:
- "Why we chose X..." → ADR
- "Command syntax: ..." → Reference

### Layer 5: Reference Heuristics

**Signals that indicate Reference content**:
- "The available options are..."
- "The default value is..."
- "Syntax: command [options]..."
- "Error CODE means..."
- Command syntax documentation
- Option/flag tables
- Environment variable lists
- Configuration option tables with defaults
- Error code tables

**Counter-examples (NOT Reference content)**:
- "This works because..." → Specification/ADR
- "The architecture consists of..." → Architecture

### Layer 6: Behaviors Heuristics

**Signals that indicate Behavior content**:
- "Given X, when Y, then Z..." (Gherkin format)
- "This behavior must not regress..."
- "Example: user does X, sees Y..."
- Explicit test case descriptions
- User journey examples with expected outcomes

**Counter-examples (NOT Behavior content)**:
- "The algorithm works by..." → Specification
- "Command options include..." → Reference

### Layer 7: Implementation Heuristics

**Signals that indicate Implementation content**:
- "Install these dependencies..."
- "The project structure is..."
- "Use this code template..."
- "Implementation phases are..."
- Technology stack with versions
- Project scaffolding instructions
- Build setup instructions

**Counter-examples (NOT Implementation content)**:
- "When X happens, the system does Y..." → Specification
- "The architecture uses..." → Architecture

---

## Content Threshold and Appendix Routing

After layer classification, apply content thresholds to determine main document vs appendix placement.

### Content Thresholds

| Threshold | Default | Action |
|-----------|---------|--------|
| `CODE_BLOCK_LINES` | 50 | Code blocks >= 50 lines route to appendices |
| `TABLE_ROWS` | 20 | Tables >= 20 rows route to appendices |
| `STEP_LIST_ITEMS` | 10 | Step lists >= 10 items route to appendices |

### Always-Appendix Patterns

These patterns ALWAYS route to appendices regardless of size:

| Pattern | Detection Method |
|---------|-----------------|
| Complete file examples | `#!/` (shebang) or package declaration at code block start |
| Shell scripts | Code block with `#!/bin/bash`, `#!/bin/sh`, etc. |
| Error catalogs | Headers like "Error Messages", "Exit Codes" with tables |

### Scaffold Code Detection

Sequential code blocks representing project scaffolding should be grouped:
- Look for file paths in headings preceding code blocks
- Group related scaffold code for single appendix placement
- Example: Multiple code blocks with headings like `cmd/main.go`, `cli/root.go`

### Main vs Appendix Decision Logic

```
Is content >= threshold OR matches always-appendix pattern?
├─ YES → Route to appendix: {layer}/appendices/{topic}/
└─ NO → Route to main document: {layer}/{topic}.md
```

Appendix naming convention: `{layer}/appendices/{main-doc-basename}/{descriptive-name}.md`

---

## Extraction Mode Assignment

After classifying content into layers, determine the extraction mode for each entry. Most entries use `verbatim` mode.

### Mode Assignment Decision Tree

Apply these questions to determine the extraction mode:

```
Is the content auto-generated (README, index, table of contents)?
├─ YES → mode: "content" (provide inline_content)
└─ NO ↓

Does the content need structural transformation?
├─ YES → Check transformation type:
│   ├─ Headings need level adjustment (e.g., H1 → H2)?
│   │   → mode: "transform" + transform_config.type: "heading_adjust"
│   ├─ Internal links need path rewriting?
│   │   → mode: "transform" + transform_config.type: "link_rewrite"
│   └─ Markdown needs formatting/cleanup?
│       → mode: "transform" + transform_config.type: "markdown_format"
└─ NO ↓

Is explicit summarization requested by user during review?
├─ YES → mode: "summarize" (NEVER assign automatically)
└─ NO → mode: "verbatim" (default)
```

### Mode Descriptions

| Mode | When to Use | Notes |
|------|-------------|-------|
| `verbatim` | Default for most content | Extracts exact lines without modification |
| `content` | Auto-generated content like README indexes | Requires `inline_content` field |
| `transform` | When source structure doesn't match target | Requires `transform_config` |
| `summarize` | Only when user explicitly requests | NEVER assign automatically |

### Transform Mode Configuration

When assigning `transform` mode, include the `transform_config`:

```yaml
transform_config:
  type: "heading_adjust"
  options:
    level_offset: 1        # Add 1 to all heading levels (# → ##)
    skip_first: true       # Don't adjust the first heading
```

**heading_adjust**: Use when source headings don't match target document structure.
- Common when extracting a subsection that uses H2 but should be H1 in target

**link_rewrite**: Use when internal documentation links need path updates.
- Common when file locations change during migration
- Requires `options.base_path` to specify the new path root

**markdown_format**: Use for general markdown cleanup (rarely needed).

### Content Mode Usage

Use `content` mode for entries that don't extract from source files:

```yaml
- id: "readme-decisions"
  target: "decisions/README.md"
  mode: "content"
  confidence: "high"
  classification_reason: "Auto-generated index for decisions directory"
  inline_content: |
    # Architecture Decisions

    This directory contains Architecture Decision Records (ADRs).

    ## Index

    - [ADR-0001: Docker Compose](../decisions/0001-docker-compose-project-name-isolation.md)
```

### Summarize Mode Restrictions

**CRITICAL**: Never automatically assign `summarize` mode during analysis.

Summarize mode:
- Requires explicit user request during the Review phase
- Must include `reviewed: true` and `reviewed_by` fields
- Should have `review_notes` explaining the summarization request

The analysis phase should assign `verbatim` mode and add a `review_notes` suggestion if content might benefit from summarization:

```yaml
- id: "long-background-section"
  source:
    file: "legacy-docs/PRD.md"
    start_line: 500
    end_line: 800
  target: "product/appendices/vision/background-research.md"
  mode: "verbatim"
  confidence: "medium"
  classification_reason: "Lengthy background section, may be candidate for summarization"
  review_notes: "Consider: This 300-line section might benefit from summarization. Request summarize mode during review if appropriate."
```

---

## Multi-File Entry Creation

Some target documents require content from multiple source files. Use multi-file entries when content is naturally split across sources.

### When to Create Multi-File Entries

Create multi-file entries when:
- A single logical document is spread across multiple source files
- Related content needs consolidation into one target document
- Vision content is split between PRD, README, and other files

### Multi-File Source Specification

```yaml
- id: "vision-consolidated"
  source:
    files:
      - file: "legacy-docs/PRD.md"
        start_line: 1
        end_line: 100
      - file: "legacy-docs/README.md"
        start_line: 10
        end_line: 50
      - file: "legacy-docs/goals.md"
        start_line: 1
        end_line: 75
    separator: "\n\n---\n\n"
    combined_hash: "g7h8i9..."
  target: "product/vision.md"
  mode: "verbatim"
  confidence: "medium"
  classification_reason: "Vision content consolidated from multiple source files"
  review_notes: "Combines PRD intro (1-100), README overview (10-50), and goals document"
```

### Multi-File Entry Guidelines

1. **Order matters**: Files are concatenated in the order listed
2. **Separator**: Use `separator` to define content between file excerpts (default: `\n\n`)
3. **Hash verification**: Include `combined_hash` for integrity checking
4. **Documentation**: Use `review_notes` to explain why files were combined
5. **Confidence**: Usually `medium` since consolidation involves judgment

### Avoid Creating Multi-File Entries When

- Content from different files serves different purposes (create separate entries instead)
- Line ranges would overlap significantly (indicates possible misclassification)
- Content doesn't logically belong together (keep separate)

---

## Confidence Scoring System

Assign confidence scores to classify content routing destination.

### High Confidence

**Criteria**:
- Multiple strong heuristic signals match
- Content fits cleanly into a single layer
- Clear main vs appendix determination

**Routing**: Directly to layer directories (e.g., `decisions/`, `specs/`)

**Reasoning example**: "Contains 'We chose X because...' pattern and discusses trade-offs with alternatives considered. Clear ADR content."

### Medium Confidence

**Criteria**:
- Category recognized but exact placement uncertain
- Single heuristic match
- Content spans multiple categories
- Unclear if main vs appendix

**Routing**: To `migration/{category}/` directory for review

**Reasoning example**: "Discusses feature behavior but mixed with implementation details. Could be Specification or Implementation layer."

### Low Confidence

**Criteria**:
- No clear heuristic match
- Unrecognized format
- Mixed content types that don't fit patterns
- Ambiguous category

**Routing**: To `blackhole/` directory

**Reasoning example**: "Format not recognized. Contains prose that doesn't match any layer heuristics."

### Reasoning Field Requirements

Every manifest entry MUST include a detailed `classification_reason` field explaining:
- Which heuristics matched (or failed to match)
- Why this confidence level was assigned
- Any uncertainty or alternative classifications considered

---

## Non-Markdown File Handling

Non-Markdown files (YAML, JSON, plain text, images, etc.) use backlink detection for layer placement.

### File Type Detection

1. Scan `LEGACY_DOCS_DIR` for all files
2. Identify file extensions: `.yaml`, `.yml`, `.json`, `.txt`, `.png`, `.jpg`, `.svg`, etc.
3. Record each non-Markdown file in the source registry

### Backlink Detection Process

For each non-Markdown file:

1. **Scan Markdown files** for references to the non-Markdown file
   - Look for file path mentions in text
   - Look for image references: `![alt](path/to/file.png)`
   - Look for link references: `[text](path/to/file.yaml)`
   - Look for code block comments mentioning the file

2. **Build reference map**: Track which Markdown files reference which non-Markdown files

3. **Determine layer from referring content**:
   - If non-Markdown file is referenced by content destined for Layer X, place in Layer X's appendices
   - If referenced by multiple layers, place in the layer with most references (document the overlap)

### Backlink-Based Layer Inference Rules

| Scenario | Action |
|----------|--------|
| Referenced by single layer's content | Place in that layer's `appendices/` |
| Referenced by multiple layers | Place in layer with most references; document other references |
| No backlinks found | Route to `blackhole/` with note about orphaned file |

**Critical**: Apply backlink inference ONLY to non-Markdown files. Markdown files use standard heuristics.

### Non-Markdown Manifest Entry Format

```yaml
- id: "schema-config-yaml"
  source: "legacy-docs/config.schema.yaml"
  target: "reference/appendices/configuration/config.schema.yaml"
  mode: "verbatim"
  confidence: "high"
  classification_reason: "YAML schema file referenced by configuration.md which is Reference layer content"
```

---

## Field Mapping Inference

For content destined for templates (ADRs, PRDs, etc.), infer field mappings.

### Primary Method: Explicit Marker Detection

Look for explicit field markers in the source content:

**ADR markers**:
- "Context:" → `context` field
- "Decision:" → `decision` field
- "Consequences:" → `consequences` field
- "Status:" → `status` field

**PRD markers**:
- "Problem:" → `problem` field
- "Goals:" → `goals` field
- "Non-goals:" → `non_goals` field

### Secondary Method: Paragraph Proximity

When explicit markers are absent:
- Content belongs to the last seen field marker until another marker appears
- Unmarked content at document start: assign to appropriate intro/context field
- Multi-paragraph fields: Continue assignment until new marker or section boundary

### Correction Mechanism: Semantic Similarity Override

Apply semantic similarity override ONLY when:
1. Confidence is very high that proximity assignment is wrong
2. Content semantically belongs to a different field
3. Override improves accuracy significantly

**Example**: A "Consequences" paragraph that appears before the "Decision" section marker but clearly discusses outcomes should be reassigned to `consequences`.

### Field Mapping Manifest Format

```yaml
field_mappings:
  title: "Docker Compose Orchestration"
  status: "Accepted"
  context: "source:legacy-docs/PRD.md:45-60"
  decision: "source:legacy-docs/PRD.md:62-75"
  consequences: "source:legacy-docs/PRD.md:77-90"
```

Values can be:
- Literal strings: Direct value assignment
- Source references: `source:{file}:{start_line}-{end_line}` for content extraction

---

## Overlapping Content Handling

Content may legitimately belong to multiple layers. Handle overlapping extractions explicitly.

### Detecting Overlapping Content

Content belongs to multiple layers when:
- A decision discussion includes both rationale (ADR) and behavioral description (Specification)
- Architecture overview contains both component diagrams (Architecture) and API contracts (Reference)
- Implementation guide includes both setup instructions (Implementation) and configuration schema (Reference)

### Creating Separate Entries

When overlap is detected:

1. **Create separate extraction entries** for each layer
2. **Allow overlapping line ranges** across entries
3. Each entry operates **independently** with its own:
   - Target path
   - Extraction mode
   - Confidence level
   - Classification reason

### Manifest Representation

Overlapping entries should have related IDs for visibility:

```yaml
entries:
  - id: "networking-adr"
    source:
      file: "legacy-docs/architecture.md"
      start_line: 100
      end_line: 180
    target: "decisions/0002-two-layer-networking.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains decision rationale with 'We chose' pattern"
    review_notes: "Overlaps with networking-arch entry (lines 120-180)"

  - id: "networking-arch"
    source:
      file: "legacy-docs/architecture.md"
      start_line: 120
      end_line: 200
    target: "architecture/networking.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains network topology diagram and component relationships"
    review_notes: "Overlaps with networking-adr entry (lines 120-180)"
```

### Entry ID Conventions for Overlapping Content

Use related naming patterns:
- `{topic}-adr` and `{topic}-arch` for ADR/Architecture overlap
- `{topic}-spec` and `{topic}-ref` for Specification/Reference overlap
- Include `review_notes` documenting the overlap for review phase visibility

---

## Error Handling and Warnings

Handle unparseable, unreadable, and excluded files explicitly.

### Unparseable File Handling

Files that cannot be classified (understood but no heuristic match):

```yaml
- id: "unknown-format-001"
  source: "legacy-docs/mysterious-file.md"
  target: "blackhole/mysterious-file.md"
  mode: "verbatim"
  confidence: "low"
  classification_reason: "Content format not recognized. Contains prose that doesn't match any layer heuristics. Manual review required."
```

**Requirements**:
- Route to `blackhole/` directory
- Set `confidence: low`
- Include detailed `classification_reason` explaining why classification failed

### Unreadable File Handling

Files with access errors go to the manifest `warnings` section:

```yaml
warnings:
  - file: "legacy-docs/locked-file.md"
    error_type: "permission_denied"
    message: "Cannot read file: Permission denied"
    suggested_resolution: "Check file permissions or run with elevated privileges"

  - file: "legacy-docs/binary-data.bin"
    error_type: "binary_file"
    message: "Binary file detected, cannot extract as documentation"
    suggested_resolution: "If this should be documentation, convert to text format"

  - file: "legacy-docs/corrupted.md"
    error_type: "encoding_error"
    message: "File encoding not recognized (not UTF-8)"
    suggested_resolution: "Convert file to UTF-8 encoding"
```

**Warning fields**:
- `file`: Path to the problematic file
- `error_type`: Category of error (`permission_denied`, `binary_file`, `encoding_error`, `file_not_found`)
- `message`: Human-readable error description
- `suggested_resolution`: Actionable fix recommendation

### Exclusion Pattern Handling

By default, scan ALL files in `LEGACY_DOCS_DIR` with no automatic exclusions.

User-specified exclusion patterns are documented in manifest metadata:

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"
  exclusions:
    patterns:
      - "*.bak"
      - ".git/**"
      - "node_modules/**"
    reason: "User-specified exclusions for backup files and dependency directories"
```

### Error Visibility in Subsequent Phases

Errors and warnings must be visible in review and extract phases:
- Review phase displays warnings prominently before entry review
- Extract phase checks for unresolved warnings and may halt if critical
- Final extraction report includes warning summary

---

## Main Manifest Generation

Generate the main manifest at `{DOCS_DIR}/.migration/extraction-manifest.yaml`.

### Manifest File Structure

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"
  description: "Main extraction manifest for project documentation migration"

sources:
  "legacy-docs/PRD.md":
    hash: "a1b2c3d4e5f6789..."
    hash_algorithm: "sha256"
    size_bytes: 15234
    last_modified: "2026-01-04T10:00:00Z"
  "legacy-docs/architecture.md":
    hash: "b2c3d4e5f6g7890..."
    hash_algorithm: "sha256"
    size_bytes: 8567

warnings:
  - file: "legacy-docs/binary.dat"
    error_type: "binary_file"
    message: "Binary file detected"
    suggested_resolution: "Remove from documentation directory"

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
  # ... additional entries

summary:
  total_entries: 42
  by_layer:
    decisions: 5
    vision: 3
    architecture: 4
    specs: 12
    reference: 8
    behaviors: 2
    implementation: 6
    blackhole: 2
  by_confidence:
    high: 35
    medium: 5
    low: 2
```

### Metadata Section

| Field | Required | Description |
|-------|----------|-------------|
| `schema_version` | Yes | Always "1.0.0" |
| `created_at` | Yes | ISO 8601 timestamp |
| `created_by` | Yes | Always "analyze.md" |
| `docs_dir` | Yes | Target documentation directory |
| `legacy_docs_dir` | No | Source directory (migration only) |
| `description` | No | Human-readable manifest description |
| `exclusions` | No | User-specified exclusion patterns |

### Sources Section

Register all source files with content hashes for change detection:

```yaml
sources:
  "{relative-path-from-project-root}":
    hash: "{sha256-hash-of-file-content}"
    hash_algorithm: "sha256"
    size_bytes: {file-size-in-bytes}
    last_modified: "{iso8601-timestamp}"
```

**Hash calculation**: Read entire file as bytes, compute SHA-256 digest.

### Entries Section

Each entry defines a single extraction operation:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier within manifest |
| `source` | Conditional | Source specification (required unless mode is `content`) |
| `target` | Yes | Target file path relative to DOCS_DIR |
| `mode` | Yes | Extraction mode: `verbatim`, `content`, `transform`, `summarize` |
| `confidence` | Yes | Confidence level: `high`, `medium`, `low` |
| `classification_reason` | Yes | Explanation of classification decision |
| `template` | No | Template name for structured documents |
| `field_mappings` | No | Template field value mappings |
| `inline_content` | Conditional | Content for `content` mode entries |
| `review_notes` | No | Notes for human reviewers |

**Line number format**: All line numbers are 1-indexed and inclusive.

### Summary Section

Provide statistics for review phase:

```yaml
summary:
  total_entries: {count}
  by_layer:
    decisions: {count}
    vision: {count}
    architecture: {count}
    specs: {count}
    reference: {count}
    behaviors: {count}
    implementation: {count}
    migration: {count}
    blackhole: {count}
  by_confidence:
    high: {count}
    medium: {count}
    low: {count}
```

---

## Layer-Specific Manifest Generation

Generate layer-specific manifests in `{DOCS_DIR}/.migration/manifests/`.

### Directory Structure

```
{DOCS_DIR}/.migration/manifests/
├── 01-decisions-manifest.yaml
├── 02-vision-manifest.yaml
├── 03-architecture-manifest.yaml
├── 04-specs-manifest.yaml
├── 05-reference-manifest.yaml
├── 06-behaviors-manifest.yaml
└── 07-implementation-manifest.yaml
```

### Naming Convention

Pattern: `{NN}-{layer}-manifest.yaml`

| File | Layer |
|------|-------|
| `01-decisions-manifest.yaml` | ADR/Decisions extraction entries |
| `02-vision-manifest.yaml` | Vision/Product extraction entries |
| `03-architecture-manifest.yaml` | Architecture extraction entries |
| `04-specs-manifest.yaml` | Specification extraction entries |
| `05-reference-manifest.yaml` | Reference extraction entries |
| `06-behaviors-manifest.yaml` | Behavior extraction entries |
| `07-implementation-manifest.yaml` | Implementation extraction entries |

### Layer Manifest Structure

Each layer manifest is self-contained for parallel execution:

```yaml
metadata:
  schema_version: "1.0.0"
  created_at: "2026-01-05T14:30:00Z"
  created_by: "analyze.md"
  layer: "decisions"
  docs_dir: "docs"
  legacy_docs_dir: "legacy-docs"
  description: "Decisions layer extraction manifest"

sources:
  # Only sources referenced by this layer's entries
  "legacy-docs/architecture.md":
    hash: "b2c3d4e5f6g7890..."
    hash_algorithm: "sha256"

entries:
  # Only entries targeting this layer
  - id: "adr-0001-docker"
    source:
      file: "legacy-docs/architecture.md"
      start_line: 45
      end_line: 95
      hash: "x1y2z3..."
    target: "decisions/0001-docker-compose-orchestration.md"
    mode: "verbatim"
    confidence: "high"
    classification_reason: "Contains decision pattern with alternatives considered"
    template: "adr-madr-minimal.md"
```

**Key differences from main manifest**:
- `metadata.layer` field identifies the layer
- `sources` section contains only files referenced by this layer's entries
- `entries` section contains only entries targeting this layer

### Parallel Execution Capability

Layer manifests are designed for parallel execution:
- Each manifest is self-contained with its own sources subset
- No cross-manifest dependencies during extraction
- Order of layer extraction does not matter
- Cross-layer linking happens in a separate post-extraction pass

---

## Analysis Workflow

Execute these steps in sequence to complete the analysis phase.

### Step 1: Load Source Documentation

1. **Read all files** from `LEGACY_DOCS_DIR`:
   - Record file path relative to project root
   - Read file content
   - Note file type (Markdown, YAML, JSON, plain text, binary)

2. **Compute content hashes**:
   - For each readable file, compute SHA-256 hash
   - Record file size in bytes
   - Record last modified timestamp

3. **Build source registry**:
   ```yaml
   sources:
     "legacy-docs/file1.md":
       hash: "{sha256}"
       size_bytes: {size}
       last_modified: "{timestamp}"
   ```

4. **Handle unreadable files**:
   - Record in warnings list with error type and suggested resolution
   - Continue processing remaining files

5. **Present summary to user**:
   > **Source Documentation Loaded**
   >
   > Found {N} files in `{LEGACY_DOCS_DIR}`:
   > - {N} Markdown files
   > - {N} YAML/JSON files
   > - {N} other files
   > - {N} files with read errors (see warnings)
   >
   > Total readable content: approximately {X} lines across {N} files.

### Step 2: Apply Section Boundary Detection

For each Markdown file:

1. **Identify headers** (H1-H6) as primary boundaries
2. **Record section information**:
   - Heading text
   - Heading level
   - Start line number
   - End line number (next heading or EOF)
3. **Identify secondary boundaries** within sections:
   - Horizontal rules
   - Multiple blank lines
   - Semantic markers
4. **Identify code block boundaries**:
   - Mark large code blocks (50+ lines) as potential appendix content
   - Group sequential code blocks with file path headings

**Output**: List of extractable sections with line ranges

### Step 3: Apply Layer Classification Heuristics

For each identified section:

1. **Apply Classification Decision Tree** (see above)
2. **Scan for layer-specific signals**:
   - Check for ADR patterns
   - Check for Vision patterns
   - Check for Architecture patterns
   - Check for Specification patterns
   - Check for Reference patterns
   - Check for Behavior patterns
   - Check for Implementation patterns
3. **Record classification result**:
   - Assigned layer
   - Signals that matched
   - Signals that were absent

### Step 4: Apply Content Thresholds

For each classified section:

1. **Check code block thresholds**:
   - Code blocks >= 50 lines → mark for appendix
   - Shebang or package declaration → mark for appendix
2. **Check table thresholds**:
   - Tables >= 20 rows → mark for appendix
3. **Check step list thresholds**:
   - Step lists >= 10 items → mark for appendix
4. **Determine target path**:
   - Main document: `{layer}/{topic}.md`
   - Appendix: `{layer}/appendices/{topic}/{descriptive-name}.md`

### Step 5: Handle Non-Markdown Files

For each non-Markdown file:

1. **Build backlink map**:
   - Scan all Markdown files for references to this file
   - Record which Markdown sections reference it
2. **Determine layer from backlinks**:
   - If referenced by classified Markdown content, use that layer
   - If multiple layers reference it, use layer with most references
3. **Create extraction entry**:
   - Target: `{layer}/appendices/{topic}/{filename}`
   - Mode: `verbatim` (copy file as-is)
4. **Handle orphaned files** (no backlinks):
   - Route to `blackhole/` with note about missing references

### Step 6: Assign Confidence Scores

For each extraction entry:

1. **Evaluate signal strength**:
   - Multiple strong signals → `high`
   - Single signal or ambiguity → `medium`
   - No clear signals → `low`
2. **Determine routing**:
   - `high` → layer directory
   - `medium` → `migration/{category}/`
   - `low` → `blackhole/`
3. **Write detailed reasoning**:
   - Document which heuristics matched
   - Document any uncertainty
   - Note alternative classifications considered

### Step 7: Infer Field Mappings

For entries with templates (ADRs, PRDs, etc.):

1. **Scan for explicit markers**:
   - "Context:", "Decision:", "Consequences:" for ADRs
   - "Problem:", "Goals:", "Non-goals:" for PRDs
2. **Apply paragraph proximity**:
   - Assign unmarked content to last seen marker
3. **Build field_mappings**:
   ```yaml
   field_mappings:
     title: "{extracted or inferred title}"
     context: "source:{file}:{start}-{end}"
     decision: "source:{file}:{start}-{end}"
   ```
4. **Apply semantic override** (conservatively):
   - Only when proximity assignment is clearly wrong
   - Document the override in reasoning

### Step 8: Generate Manifests

1. **Create `.migration/` directory** in `DOCS_DIR`
2. **Generate main manifest** at `.migration/extraction-manifest.yaml`:
   - Include full metadata section
   - Include complete sources registry
   - Include all entries
   - Include warnings section
   - Include summary statistics
3. **Create `.migration/manifests/` directory**
4. **Generate layer-specific manifests**:
   - For each layer, create `{NN}-{layer}-manifest.yaml`
   - Include only entries and sources for that layer
   - Each manifest self-contained

### Step 9: Present Analysis Summary

Present the analysis results to the user:

> **Analysis Complete**
>
> Generated extraction manifest at `.migration/extraction-manifest.yaml`
>
> ## Summary
>
> | Metric | Count |
> |--------|-------|
> | Total entries | {N} |
> | High confidence | {N} |
> | Medium confidence | {N} |
> | Low confidence | {N} |
> | Warnings | {N} |
>
> ## Entries by Layer
>
> | Layer | Entries |
> |-------|---------|
> | Decisions (ADR) | {N} |
> | Vision | {N} |
> | Architecture | {N} |
> | Specifications | {N} |
> | Reference | {N} |
> | Behaviors | {N} |
> | Implementation | {N} |
> | migration/ | {N} |
> | blackhole/ | {N} |
>
> ## Next Steps
>
> 1. Review the generated manifest: `.migration/extraction-manifest.yaml`
> 2. Proceed to the Review phase: `review.md`
> 3. After review approval, proceed to Extract phase: `extract.md`
>
> **Proceed to review phase?** Execute: `review.md`

---

## Output Files

This workflow produces:

| File | Description |
|------|-------------|
| `.migration/extraction-manifest.yaml` | Main manifest with all entries |
| `.migration/manifests/01-decisions-manifest.yaml` | ADR extraction entries |
| `.migration/manifests/02-vision-manifest.yaml` | Vision/product extraction entries |
| `.migration/manifests/03-architecture-manifest.yaml` | Architecture extraction entries |
| `.migration/manifests/04-specs-manifest.yaml` | Specification extraction entries |
| `.migration/manifests/05-reference-manifest.yaml` | Reference extraction entries |
| `.migration/manifests/06-behaviors-manifest.yaml` | Behavior extraction entries |
| `.migration/manifests/07-implementation-manifest.yaml` | Implementation extraction entries |

---

## Manifest Schema Reference

For the complete manifest schema definition, see:
- `schemas/extraction-manifest.schema.yaml`

Key manifest components:
- **metadata**: Schema version, creation time, directories, layer (for layer-specific)
- **sources**: Source file registry with content hashes
- **warnings**: File access errors and issues
- **entries**: List of extraction entries with source, target, mode, and confidence
- **summary**: Entry counts by layer and confidence

---

## Next Phase

After analysis completes, output the next step command:

```
ANALYSIS PHASE COMPLETE
=======================

Generated extraction manifests in .migration/manifests/

Summary:
  - Total entries: {N}
  - High confidence: {N} (will extract to layers directly)
  - Medium confidence: {N} (needs review, will go to migration/)
  - Low confidence: {N} (needs review, will go to blackhole/)

Manifest files created:
  - .migration/extraction-manifest.yaml (main manifest)
  - .migration/manifests/01-decisions-manifest.yaml
  - .migration/manifests/02-vision-manifest.yaml
  - .migration/manifests/03-architecture-manifest.yaml
  - .migration/manifests/04-specs-manifest.yaml
  - .migration/manifests/05-reference-manifest.yaml
  - .migration/manifests/06-behaviors-manifest.yaml
  - .migration/manifests/07-implementation-manifest.yaml

┌─────────────────────────────────────────────────────────────────────┐
│ NEXT STEP                                                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ Run /clear, then execute the review phase:                          │
│                                                                     │
│   Execute @docs/.migration/review.md                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

The review phase allows you to:
  - Inspect all proposed extractions
  - Adjust confidence levels
  - Change target locations
  - Approve the manifest for extraction
```

---

## Related Documentation

- `install.md` - Main installation workflow (orchestrates all phases)
- `review.md` - Review phase workflow
- `extract.md` - Extraction phase workflow
- `schemas/extraction-manifest.schema.yaml` - Manifest schema definition
- `DOCUMENTATION-GUIDE.md` - Classification heuristics and layer definitions
