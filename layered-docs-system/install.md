# Layered Documentation System — Install

**For AI Agents**: This document contains instructions for installing and customizing the Layered Documentation System for a project. Read this entire file, then follow the process below.

---

## Prerequisites

Before starting, read the complete system documentation:

- `LAYERED-DOCUMENTATION-SYSTEM.md` — Full system guide with all layers, heuristics, and templates

Pay special attention to:
- **Directory Structure** — Layers use flat file structure with `{topic}.md` files and shared `appendices/{topic}/` directories; ADRs are simple single files
- **Appendix System** — Large content goes to appendices (except for ADRs, which include code inline)
- **Confidence-Based Classification** — Content goes to layers, `migration/`, or `blackhole/`
- **Content Thresholds** — Default thresholds for appendix classification

---

## Critical Migration Principle

**Migration means PRESERVING content, not summarizing it.**

When migrating existing documentation:
- **Preserve ALL technical details** — Every code example, error message, output sample, and configuration snippet must be included in the migrated documentation
- **The migrated documents should contain the same depth of information as the source** — If the source has a 200-line code example, the migrated document (or its appendix) must contain that 200-line code example
- **Do not abbreviate or summarize** — Move content, don't rewrite it
- **Reorganization, not reduction** — The goal is to put content in the right place, not to make it shorter

A successful migration should result in approximately the same total line count as the source (possibly slightly higher due to added structure and cross-references).

---

## Installation Modes

This installer supports two modes:

1. **Fresh Install** — Create an empty documentation system with templates
2. **Migration Install** — Analyze existing documentation and reorganize into layers

The mode is determined by whether existing documentation is specified.

### Migration Install Includes

When migrating existing documentation, the installer:
1. Classifies content into layers with confidence scoring
2. Places high-confidence content directly into layers
3. Places medium-confidence content in `migration/` for review
4. Places low/no-confidence content in `blackhole/` as a safety net
5. Runs an automatic comparison audit (same as `audit.md`)
6. Reports what was migrated, what needs review, and what was lost

---

## Installation Process

### Step 1: Determine Install Location

If the user specified an install directory (e.g., `@docs/`), use that path.

If no directory was specified, ask:

> **Install Location**
>
> Where would you like to install the documentation system? Common locations:
> - `docs/` — Standard documentation directory
> - `documentation/` — Alternative naming
> - `.docs/` — Hidden directory (less common)
>
> Please specify a path, or press enter for the default (`docs/`).

Store the chosen path as `INSTALL_DIR`.

---

### Step 2: Determine Source Documentation

If the user specified existing documentation (e.g., `@specs/`, `@existing-docs/`), use that path.

If no source was specified, ask:

> **Existing Documentation**
>
> Do you have existing documentation to migrate into the layered system?
>
> 1. **Yes** — I have existing docs to reorganize (specify path)
> 2. **No** — Start fresh with empty templates
>
> If yes, please provide the path to your existing documentation (e.g., `specs/`, `docs/legacy/`).

Store the result:
- If path provided → `SOURCE_DOCS = {path}`, `MIGRATION_MODE = true`
- If "no" or empty → `SOURCE_DOCS = null`, `MIGRATION_MODE = false`

---

### Step 3: Load and Analyze Source Documentation (Migration Mode Only)

**Skip this step if `MIGRATION_MODE = false`.**

#### 3a: Read All Source Files

Read ALL files in `SOURCE_DOCS`. For each file:
- Record the filename and full content
- Note the file type (Markdown, YAML, etc.)

Present a summary to the user:

> **Source Documentation Loaded**
>
> Found {N} files in `{SOURCE_DOCS}`:
> - `filename1.md` — {first 50 chars or heading}
> - `filename2.md` — {first 50 chars or heading}
> - ...
>
> Total content: approximately {X} lines across {N} files.

#### 3b: Analyze Content Against Layer Heuristics

For each piece of content in the source files, apply the classification heuristics from `LAYERED-DOCUMENTATION-SYSTEM.md`:

**Layer 1 — Decisions (ADRs)**
Signals to look for:
- "We chose X over Y because..."
- "We decided to..." / "Decision:"
- "This pattern will be used throughout..."
- "We considered alternatives including..."
- Explicit "Decision" blocks or sections
- Trade-off discussions with rationale

**Layer 2 — Vision (PRD-Lite)**
Signals to look for:
- Problem statements ("The problem we're solving...")
- Product vision or mission statements
- Success criteria or goals
- Non-goals or out-of-scope sections
- Core concepts and glossary/definitions
- User personas or target audience
- "Executive Summary" sections

**Layer 2 — Comparison (Product Positioning)**
Signals to look for:
- Comparison tables with other tools/products
- "vs" language ("X vs Y", "compared to", "alternatives")
- Feature matrices comparing multiple products
- "Related Tools" or "Similar Projects" sections
- Competitive analysis content
- "Why choose X over Y" discussions
- Tool/product name mentions in comparative context (e.g., "Unlike Docker Compose...")

**Layer 2 — Roadmap (Future Direction)**
Signals to look for:
- "Future Considerations" sections
- "Roadmap" or "Future Work" headings
- "Planned features" or "upcoming" language
- "Phase 2", "Phase 3", "v2.0" planning content
- "Not yet implemented" with intent to implement
- "Eventually" or "in the future" feature descriptions
- Enhancement proposals or feature wishlists

**Layer 3 — Architecture**
Signals to look for:
- System diagrams (ASCII art, Mermaid, etc.)
- Network topology descriptions
- Component relationship descriptions ("A communicates with B via...")
- Cross-cutting concerns (security patterns, logging strategies)
- "Architecture" section headings
- Container/service descriptions at a high level

**Layer 4 — Specifications**
Signals to look for:
- Detailed behavioral descriptions ("When X happens, the system does Y...")
- State machines or lifecycle descriptions
- Data schemas with field-level detail
- Edge case documentation
- Feature-specific sections with implementation detail
- Configuration file format specifications
- Error handling descriptions

**Layer 5 — Reference**
Signals to look for:
- Command syntax documentation
- Option/flag tables
- Environment variable lists
- Configuration option tables with defaults
- API endpoint lists
- Error code tables
- Quick-reference tables

**Layer 6 — Behaviors**
Signals to look for:
- "Given/When/Then" scenarios
- Explicit test case descriptions
- User journey examples with expected outcomes

**Layer 7 — Implementation**
Signals to look for:
- Technology stack with versions
- Dependency lists
- Project scaffolding instructions
- Code templates
- Build setup instructions

#### 3c: Apply Content Thresholds

For each piece of content, check against thresholds from `LAYERED-DOCUMENTATION-SYSTEM.md`:

| Threshold | Default | Check |
|-----------|---------|-------|
| `CODE_BLOCK_LINES` | 50 | Code blocks ≥ 50 lines → mark as appendix content |
| `STEP_LIST_ITEMS` | 10 | Step lists ≥ 10 items → mark as appendix content |
| `TABLE_ROWS` | 20 | Tables ≥ 20 rows → mark as appendix content |
| `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Complete file examples → always appendix |
| `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Error catalogs → always appendix |
| `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Complete shell scripts → always appendix |

**Pattern Detection**:
- `#!/` or `package main` at start of code block → complete file, appendix
- Sequential code blocks with file paths in headings → scaffold code, appendix
- Headers like "Error Messages", "Exit Codes" with tables → error catalog, appendix
- Numbered steps with explanatory text between commands → workflow example, appendix

#### 3d: Assign Confidence Scores

For each content section, assign a confidence level:

| Confidence | Criteria | Destination |
|------------|----------|-------------|
| **High** | Clear heuristic match, obvious layer, clear main vs appendix | Layer directly |
| **Medium** | Category recognized but placement uncertain OR unclear if main vs appendix | `migration/{category}/` |
| **Low/None** | No heuristic match, unrecognized format, ambiguous category | `blackhole/` |

**Confidence indicators**:
- **High**: Multiple strong heuristic signals, content fits template cleanly
- **Medium**: Single heuristic match, or content spans multiple categories
- **Low**: No clear heuristic match, unusual format, mixed content types

#### 3e: Create Content Mapping

Build a mapping of source content to target locations:

```
CONTENT_MAP = {
  "Layer 1: Decisions": [
    { source: "prd.md", section: "Decision: Docker Compose", lines: "45-62", confidence: "high", appendix: false },
    { source: "prd.md", section: "Decision: Two-Layer Networking", lines: "78-95", confidence: "high", appendix: false },
    ...
  ],
  "Layer 5: Reference": [
    { source: "cli-reference.md", section: "Commands", lines: "1-200", confidence: "high", appendix: false },
    { source: "cli-reference.md", section: "Detailed Examples", lines: "201-800", confidence: "high", appendix: true },
    { source: "cli-reference.md", section: "Error Messages", lines: "801-1000", confidence: "high", appendix: true },
    ...
  ],
  "Layer 7: Implementation": [
    { source: "go-stack.md", section: "Overview", lines: "1-100", confidence: "high", appendix: false },
    { source: "go-stack.md", section: "Cobra Scaffolding", lines: "101-1200", confidence: "high", appendix: true },
    { source: "go-stack.md", section: "Shell Scripts", lines: "1201-1600", confidence: "high", appendix: true },
    ...
  ],
  "migration": [
    { source: "shell-integration.md", section: "Alternative Approaches", lines: "600-650", confidence: "medium", category: "decisions-or-specs" },
    ...
  ],
  "blackhole": [
    { source: "go-stack.md", section: "Unknown Format Block", lines: "1617-1650", confidence: "low", reason: "Format not recognized" },
    ...
  ]
}
```

#### 3f: Present Analysis to User

Present the content mapping for review:

> **Content Analysis Complete**
>
> Based on the classification heuristics, here's how your existing content maps:
>
> ## High Confidence (will be migrated directly)
>
> **Layer 1: Decisions** — {N} sections
> - Decision blocks from `prd.md`
> - Architecture rationale from `technical-spec.md`
>
> **Layer 5: Reference** — {N} sections
> - Main content: Commands, options, flags
> - **Appendix**: Detailed examples ({X} lines)
> - **Appendix**: Error message catalog ({Y} lines)
>
> **Layer 7: Implementation** — {N} sections
> - Main content: Overview, dependencies
> - **Appendix**: Go scaffolding code ({X} lines)
> - **Appendix**: Shell scripts ({Y} lines)
>
> [Continue for all layers...]
>
> ## Medium Confidence (will go to `migration/`)
>
> - `shell-integration.md:600-650` → `migration/decisions-or-specs/` — Alternative approaches section
> - {Other medium-confidence items}
>
> ## Low/No Confidence (will go to `blackhole/`)
>
> - `go-stack.md:1617-1650` → `blackhole/` — Unknown format block
> - {Other low-confidence items}
>
> ---
>
> **Summary**:
> - High confidence: {N} sections → layers
> - Medium confidence: {N} sections → migration/
> - Low confidence: {N} sections → blackhole/
>
> Does this mapping look correct? (yes/no/adjust)

If "adjust", allow the user to manually reassign sections or change confidence levels.

---

### Step 4: Gather Layer Selections

**For Fresh Install**: Present all options.

**For Migration Install**: Pre-select layers based on content analysis. Only layers with identified content are pre-selected, but user can add additional empty layers.

#### Layer 1: Decisions (ADRs)

> **Layer 1: Decisions (ADRs)**
>
> ADRs capture *why* significant technical and product decisions were made.
>
> Template options:
> 1. **MADR Minimal** (Recommended) — Structured but lightweight
> 2. **Y-Statement** — Single-sentence format
> 3. **MADR Full** — Comprehensive with pros/cons analysis
> 4. **Skip** — Do not include this layer
>
> Which template? (1/2/3/4)

Store selection as `DECISION_ADR_TEMPLATE`.

#### Layer 2: Vision (PRD-Lite)

> **Layer 2: Vision (PRD-Lite)**
>
> Vision documents define the product's purpose, goals, and constraints.
>
> Template options:
> 1. **Lean PRD** (Recommended) — Vision + problem + concepts + non-goals
> 2. **Epic-Based PRD** — Organized around user stories and epics
> 3. **Skip** — Do not include this layer
>
> Which template? (1/2/3)

Store selection as `DECISION_VISION_TEMPLATE`.

#### Layer 2: Additional Product Documents

> **Layer 2: Additional Product Documents**
>
> Beyond the core vision, the product layer can include:
>
> Optional documents:
> 1. **Comparison** — How this product compares to alternatives (auto-created if comparison content detected)
> 2. **Roadmap** — Future considerations and planned enhancements (auto-created if roadmap content detected)
>
> These documents are created automatically during migration if matching content is found. For fresh installs, they can be added later using `create.md`.

Store selection as `DECISION_PRODUCT_EXTRAS` (note: typically auto-detected during migration).

#### Layer 3: Architecture

> **Layer 3: Architecture**
>
> Architecture documents show how components relate.
>
> Template options:
> 1. **C4-Lite** (Recommended) — Context + Container diagrams with narrative
> 2. **arc42 Full** — Comprehensive 12-section template
> 3. **Skip** — Do not include this layer
>
> Which template? (1/2/3)

Store selection as `DECISION_ARCHITECTURE_TEMPLATE`.

#### Layer 4: Specifications

> **Layer 4: Specifications**
>
> Specifications detail how features work.
>
> Template options:
> 1. **Feature Spec** (Recommended) — Narrative + examples + edge cases
> 2. **RFC-Style** — Problem + proposal + alternatives (for reviews)
> 3. **Both** — Include both templates
> 4. **Skip** — Do not include this layer
>
> Which template? (1/2/3/4)

Store selection as `DECISION_SPEC_TEMPLATE`.

#### Layer 5: Reference

> **Layer 5: Reference**
>
> Reference documents provide lookup information.
>
> Template options:
> 1. **CLI Reference** — Command documentation
> 2. **Configuration Reference** — Config options and defaults
> 3. **Both** — Include both templates
> 4. **Skip** — Do not include this layer
>
> Which template? (1/2/3/4)

Store selection as `DECISION_REFERENCE_TEMPLATE`.

#### Layer 6: Behaviors (Executable Specifications)

> **Layer 6: Behaviors (Executable Specifications)**
>
> Behavior specs define expected behaviors in Gherkin format.
>
> Template options:
> 1. **Gherkin** — Standard BDD format
> 2. **Skip** — Do not include this layer
>
> Which template? (1/2)

Store selection as `DECISION_BEHAVIOR_TEMPLATE`.

#### Layer 7: Implementation Guides

> **Layer 7: Implementation Guides**
>
> Implementation guides describe how to build the system.
>
> Template options:
> 1. **Tech Stack** — Dependencies, versions, rationale
> 2. **Skip** — Do not include this layer
>
> Which template? (1/2)

Store selection as `DECISION_IMPLEMENTATION_TEMPLATE`.

---

### Step 5: Gather Tooling Selections

> **Tooling Tier Selection**
>
> Which level of tooling would you like documented?
>
> 1. **Tier 1: Essential** — markdownlint, Git + PRs only
> 2. **Tier 2: Recommended** — Adds Vale, Log4brains, Structurizr
> 3. **Tier 3: Advanced** — Adds Cucumber, link checking, custom generators
> 4. **None** — No tooling recommendations
>
> Which tier? (1/2/3/4)

Store selection as `DECISION_TOOLING_TIER`.

---

### Step 6: Confirm Selections

Present a summary:

> **Installation Summary**
>
> **Mode**: {Fresh Install | Migration from `{SOURCE_DOCS}`}
> **Install Location**: `{INSTALL_DIR}`
>
> **Layers to Include**:
> - Decisions: {template} {if migration: "+ N extracted ADRs"}
> - Vision: {template} {if migration: "+ migrated content"}
> - Architecture: {template} {if migration: "+ migrated content"}
> - Specifications: {template} {if migration: "+ migrated content"}
> - Reference: {template} {if migration: "+ migrated content"}
> - Behaviors: {template}
> - Implementation: {template}
>
> **Tooling Tier**: {tier}
>
> Proceed with installation? (yes/no)

If "no", ask which selections to change.

---

### Step 7: Create Directory Structure

Create directories based on selected layers. Note that **ADRs use simple single files** (not directories), while other layers use flat file structure with `{topic}.md` files and shared `appendices/{topic}/` directories:

```
{INSTALL_DIR}/
├── DOCUMENTATION-GUIDE.md          # Always created (includes thresholds)
├── audit.md                         # Always created (audit instructions for this install)
│
├── decisions/                       # If Layer 1 selected (simple files, NOT directories)
│   ├── README.md                   # Auto-generated index of all ADRs
│   └── 0000-template.md            # Template ADR (single file)
│
├── product/                         # If Layer 2 selected
│   ├── README.md                   # Auto-generated index
│   ├── vision.md                   # Template or migrated content
│   ├── comparison.md               # Only if comparison content detected (migration)
│   └── roadmap.md                  # Only if roadmap content detected (migration)
│
├── architecture/                    # If Layer 3 selected
│   ├── README.md                   # Auto-generated index
│   └── overview.md                 # Template or migrated content
│
├── specs/                           # If Layer 4 selected
│   ├── README.md                   # Auto-generated index
│   └── _template.md                # Template file
│
├── reference/                       # If Layer 5 selected
│   ├── README.md                   # Auto-generated index
│   ├── cli.md                      # CLI reference
│   ├── configuration.md            # Config reference
│   └── appendices/                 # Only if migrating content
│       ├── cli/                    # Appendices for cli.md
│       │   └── ...
│       └── configuration/          # Appendices for configuration.md
│           └── ...
│
├── implementation/                  # If Layer 7 selected
│   ├── README.md                   # Auto-generated index
│   ├── tech-stack.md               # Template or migrated content
│   └── appendices/                 # Only if migrating content
│       └── tech-stack/             # Appendices for tech-stack.md
│           └── ...
│
├── migration/                       # Only if MIGRATION_MODE and has medium-confidence content
│   └── README.md                   # Explains what's here
│
├── blackhole/                       # Only if MIGRATION_MODE and has low-confidence content
│   └── README.md                   # Explains what's here and suggests heuristic updates
│
└── features/                        # If Layer 6 selected (at project root)
    └── _template.feature           # Template file
```

**Note**: `migration/` and `blackhole/` directories are ONLY created if there is content to place there. Empty directories are not created.

---

### Step 8: Generate Migration Plan (Migration Mode Only)

**Skip this step if `MIGRATION_MODE = false`.**

Rather than performing migration inline (which risks context window exhaustion), this step generates a migration plan with step files that can be executed in separate sessions.

#### 8a: Create Migration Directory

Create the migration workspace:

```
{INSTALL_DIR}/.migration/
├── README.md                    # Overview and execution instructions
├── common-instructions.md       # Shared migration principles
├── 01-decisions.md              # ADR extraction steps
├── 02-vision.md                 # Vision content migration
├── 03-architecture.md           # Architecture content migration
├── 04-specs.md                  # Specification content migration
├── 05-reference.md              # Reference content migration
├── 06-implementation.md         # Implementation content migration
├── 07-medium-confidence.md      # Medium-confidence content for review
├── 08-low-confidence.md         # Low-confidence content (blackhole)
└── 09-cross-links.md            # Cross-layer linking pass
```

**Note**: The `.migration/` directory uses a leading dot to indicate it's a working directory, not part of the final documentation.

#### 8b: Generate Common Instructions File

Create `{INSTALL_DIR}/.migration/common-instructions.md`:

```markdown
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
- `specs/prd.md:45-62` means lines 45-62 of `specs/prd.md`
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
- `{INSTALL_DIR}/decisions/0000-template.md` for ADRs
- `{INSTALL_DIR}/specs/_template.md` for specifications
- Other layer templates as applicable

## Cross-Layer Links

Do NOT add cross-layer links during individual step execution. A separate pass (09-cross-links.md) handles this after all content is migrated.
```

#### 8c: Generate Migration Step Files

For each layer with mapped content, create a step file. Each step file should be self-contained and executable in a separate session.

**Deciding Content vs Pointers**:
- If the total content for a step file is under ~2000 lines, include the pre-extracted content directly
- If larger, use source pointers and let the executing agent read the source
- Err on the side of pre-extracting content when feasible, as it produces more reliable results

**Step File Format**:

```markdown
# Migration Step: {Layer Name}

**Prerequisites**:
- Read `common-instructions.md` first
- Directory `{target_dir}` must exist

**Estimated Size**: {N} lines of content to create

---

## Output File 1: `{relative/path/to/output.md}`

### Content

{Either pre-extracted content OR source pointer}

<!-- If pre-extracted: -->
```markdown
{Full content to write to the file}
```

<!-- If source pointer: -->
**Source**: `{source_file}:{start_line}-{end_line}`
**Template**: Use `{template_name}` structure
**Sections to extract**:
- Context: lines {X}-{Y}
- Decision: lines {A}-{B}
- Code example: lines {C}-{D} (include inline)

### Post-Creation

- [ ] Add source attribution: `<!-- Migrated from {source_file}:{lines} -->`
- [ ] Verify line count is approximately {expected}

---

## Output File 2: `{relative/path/to/output2.md}`

{Repeat format...}

---

## Completion Checklist

- [ ] All output files created
- [ ] Source attributions added
- [ ] Line counts verified (not summarized)
```

#### 8d: Generate Decisions Step File

Create `{INSTALL_DIR}/.migration/01-decisions.md`:

For each ADR identified in Step 3:
1. Include the full content to write (pre-extracted) if total is under ~2000 lines
2. Otherwise, include precise source pointers with section mapping
3. ADRs are single files — specify exact filename: `NNNN-{slugified-title}.md`
4. Include any code examples inline (ADRs don't use appendices)

Example structure:

```markdown
# Migration Step: Layer 1 — Decisions (ADRs)

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: {N} ADRs, approximately {M} lines total

---

## ADR 1: `decisions/0001-docker-compose-orchestration.md`

### Content

```markdown
# ADR-0001: Docker Compose as Orchestration Layer

## Status

Accepted

## Context

{Pre-extracted context from source, preserving full detail...}

## Decision

{Pre-extracted decision...}

## Consequences

{Pre-extracted consequences...}

<!-- Migrated from specs/technical-spec.md:45-92 -->
```

### Verification

- [ ] Contains full context (not summarized)
- [ ] Code examples included inline
- [ ] Source attribution present

---

## ADR 2: `decisions/0002-two-layer-networking.md`

{Continue for each ADR...}
```

#### 8e: Generate Vision Step File

Create `{INSTALL_DIR}/.migration/02-vision.md`:

1. Target file: `product/vision.md`
2. Pre-extract content if feasible, otherwise provide section-by-section source pointers
3. Include appendix file generation if needed (e.g., `appendices/vision/glossary.md`)

#### 8f: Generate Architecture Step File

Create `{INSTALL_DIR}/.migration/03-architecture.md`:

1. Target file: `architecture/overview.md`
2. Preserve diagrams exactly — include as pre-extracted content to avoid corruption
3. Include appendix file generation if needed (e.g., `appendices/overview/c4-diagrams.md`)

#### 8g: Generate Specifications Step File

Create `{INSTALL_DIR}/.migration/04-specs.md`:

1. Multiple target files: `specs/{feature-name}.md`
2. For each specification, include:
   - Main content (pre-extracted or source pointer)
   - Appendix content if needed (at `appendices/{feature-name}/`)
   - Notes about removing decision rationale (now in ADRs)

#### 8h: Generate Reference Step File

Create `{INSTALL_DIR}/.migration/05-reference.md`:

1. Target files: `reference/cli.md`, `reference/configuration.md`
2. Appendix files: `reference/appendices/cli/*.md`, `reference/appendices/configuration/*.md`
3. Emphasize: Include ALL content — every command, flag, error message
4. Pre-extract when possible since reference docs are critical

#### 8i: Generate Implementation Step File

Create `{INSTALL_DIR}/.migration/06-implementation.md`:

1. Target file: `implementation/tech-stack.md`
2. Appendix files:
   - `implementation/appendices/tech-stack/go-scaffolding/*.md`
   - `implementation/appendices/tech-stack/shell-scripts/*.md`
3. **Critical**: Scaffolding code and shell scripts must be preserved in full
4. Pre-extract scaffolding code to ensure no loss

#### 8j: Generate Medium-Confidence Step File

Create `{INSTALL_DIR}/.migration/07-medium-confidence.md`:

For content with medium confidence from Step 3:
1. Target directory: `migration/{category}/`
2. Include source pointers with uncertainty reasons
3. Include suggested destinations for human review

#### 8k: Generate Low-Confidence Step File

Create `{INSTALL_DIR}/.migration/08-low-confidence.md`:

For content with low/no confidence from Step 3:
1. Target directory: `blackhole/`
2. Include source pointers with classification failure reasons
3. Include suggested heuristic updates

#### 8l: Generate Cross-Links Step File

Create `{INSTALL_DIR}/.migration/09-cross-links.md`:

After all content is migrated, this step adds cross-layer links:

```markdown
# Migration Step: Cross-Layer Links

**Prerequisites**:
- All previous migration steps (01-08) must be complete
- Read `common-instructions.md`

## Instructions

Scan all created files and add appropriate cross-layer links:

### Link Patterns to Add

**In specifications** — Link to ADRs for rationale:
```markdown
See [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md) for rationale.
```

**In main documents** — Link to appendices for details:
```markdown
For complete workflow examples, see [Detailed Examples](./appendices/cli/detailed-examples.md).
```

**In architecture** — Link to specs for detail:
```markdown
For detailed port assignment behavior, see [Port Types Spec](../specs/port-types.md).
```

### Files to Update

{List of files that need cross-links based on content mapping}

### Verification

- [ ] All ADR references added where decisions are mentioned
- [ ] All appendix links added where detailed content exists
- [ ] All architecture-to-spec links added
- [ ] Links use correct relative paths
```

#### 8m: Generate Migration README

Create `{INSTALL_DIR}/.migration/README.md`:

```markdown
# Migration Plan

**Generated**: {timestamp}
**Source**: `{SOURCE_DOCS}`
**Target**: `{INSTALL_DIR}`

## Overview

This migration is split into {N} step files to avoid context window issues. Each step can be executed in a separate session.

## Content Summary

| Step | Layer | Files | Lines | Status |
|------|-------|-------|-------|--------|
| 01 | Decisions | {N} ADRs | ~{M} lines | ⬜ Pending |
| 02 | Vision | 1 file | ~{M} lines | ⬜ Pending |
| 03 | Architecture | 1 file | ~{M} lines | ⬜ Pending |
| 04 | Specifications | {N} files | ~{M} lines | ⬜ Pending |
| 05 | Reference | {N} files | ~{M} lines | ⬜ Pending |
| 06 | Implementation | {N} files | ~{M} lines | ⬜ Pending |
| 07 | Medium-Confidence | {N} files | ~{M} lines | ⬜ Pending |
| 08 | Low-Confidence | {N} files | ~{M} lines | ⬜ Pending |
| 09 | Cross-Links | — | — | ⬜ Pending |

## Execution Instructions

### For Each Step:

1. Start a new session (to have fresh context)
2. Read `common-instructions.md`
3. Read the step file (e.g., `01-decisions.md`)
4. Execute the migration instructions in the step file
5. Mark the step complete in this README
6. Proceed to next step (in same or new session as needed)

### Recommended Order

Execute steps 01-08 in order. Step 09 (cross-links) must be executed last.

Steps 01-06 (layer content) can potentially be parallelized if running multiple sessions, but dependencies may exist if content references other layers.

## After All Steps Complete

1. Run the audit workflow (see Step 14 in install.md, or use `{INSTALL_DIR}/audit.md`)
2. Review content in `migration/` directory and move to appropriate layers
3. Review content in `blackhole/` directory
4. Delete this `.migration/` directory:
   ```bash
   rm -rf {INSTALL_DIR}/.migration/
   ```

## Step Execution Log

| Step | Executed By | Date | Notes |
|------|-------------|------|-------|
| 01 | | | |
| 02 | | | |
| ... | | | |
```

---

### Step 9: Execute Migration Steps (Migration Mode Only)

**Skip this step if `MIGRATION_MODE = false`.**

This step is executed in separate sessions after the install phase completes. Each migration step file can be executed independently.

#### How to Execute Migration Steps

The install phase (Steps 1-8) generates migration step files in `{INSTALL_DIR}/.migration/`. To execute them:

1. **Start a new session** (to have fresh context)
2. **Provide the step file to the agent**:
   ```
   Execute the migration step in @docs/.migration/01-decisions.md
   ```
3. **The agent should**:
   - Read `common-instructions.md` first
   - Read the specified step file
   - Create each output file as instructed
   - Verify line counts match expectations
4. **Mark the step complete** in `.migration/README.md`
5. **Repeat** for each step file in order

#### Execution Order

| Step | File | Dependencies |
|------|------|--------------|
| 01 | `01-decisions.md` | None |
| 02 | `02-vision.md` | None |
| 03 | `03-architecture.md` | None |
| 04 | `04-specs.md` | None (but may reference ADRs) |
| 05 | `05-reference.md` | None |
| 06 | `06-implementation.md` | None |
| 07 | `07-medium-confidence.md` | None |
| 08 | `08-low-confidence.md` | None |
| 09 | `09-cross-links.md` | All of the above |

Steps 01-08 can be executed in any order (or in parallel across multiple sessions). Step 09 must be executed last.

#### Verifying Each Step

After each step, verify:
- [ ] All output files were created
- [ ] Line counts approximately match source (not summarized)
- [ ] Source attribution comments are present
- [ ] Any code examples are preserved in full

#### After All Steps Complete

Proceed to Step 14 (Run Migration Audit) to validate the migration.

---

### Step 10: Generate Customized Guide

Create `{INSTALL_DIR}/DOCUMENTATION-GUIDE.md` containing:

1. **Header** with project-specific title
2. **Migration Note** (if applicable): "This documentation was migrated from `{SOURCE_DOCS}` on {date}"
3. **Content Thresholds** — Configurable thresholds section:
   ```markdown
   ## Content Thresholds (Configurable)

   These thresholds control how content is classified. Edit to customize for this project.

   | Threshold | Value | Purpose |
   |-----------|-------|---------|
   | `CODE_BLOCK_LINES` | 50 | Code blocks ≥ this go to appendix |
   | `STEP_LIST_ITEMS` | 10 | Step lists ≥ this go to appendix |
   | `TABLE_ROWS` | 20 | Tables ≥ this go to appendix |
   | `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Complete file examples → appendix |
   | `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Error catalogs → appendix |
   | `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Shell scripts → appendix |
   ```
4. **Layer Overview** — Only selected layers
5. **Classification Heuristics** — Only for selected layers
6. **Decision Tree** — Simplified for included layers
7. **Cross-Layer Linking** — Adjusted for selected layers
8. **Tooling Section** — Selected tier only
9. **Directory Structure** — Actual created structure
10. **Appendix Guidelines** — How to create and link appendices

---

### Step 11: Generate Audit File

Create `{INSTALL_DIR}/audit.md` by combining the appropriate audit templates.

#### Template Sources

The audit templates are located in `layered-docs-system/audit/`:
- `common-audit.md` — Core audit process (always included)
- `audit-fresh-install.md` — For fresh installs (no source comparison)
- `audit-migration.md` — For migrations (includes source comparison)

#### 11a: Generate Fresh Install Audit

**If `MIGRATION_MODE = false`**, create `{INSTALL_DIR}/audit.md` by combining:
1. `common-audit.md`
2. `audit-fresh-install.md`

The file should be customized with:
- `{INSTALL_DIR}` replaced with the actual documentation path

#### 11b: Generate Migration Audit

**If `MIGRATION_MODE = true`**, create `{INSTALL_DIR}/audit.md` by combining:
1. `audit-migration.md` (which includes migration-specific instructions)
2. `common-audit.md`

The file should be customized with:
- `{SOURCE_DOCS}` replaced with the actual source documentation path
- `{INSTALL_DIR}` replaced with the actual documentation path

**Migration Audit Template**:

```markdown
# Documentation Audit

**For AI Agents**: This document contains instructions for auditing the documentation in this directory.

---

## Audit Configuration

- **Documentation Directory**: `{INSTALL_DIR}`
- **Source Documentation**: `{SOURCE_DOCS}`
- **Install Type**: Migration
- **Install Date**: {timestamp}

---

## Critical Audit Principles

**IMPORTANT**: The migrated docs use a layered structure with sub-directories. Before claiming any content is missing, you MUST:

1. **Read ALL Markdown files** under `{INSTALL_DIR}/` — not just README.md files, but every `.md` file in every subdirectory including `appendices/` folders
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

Read ALL files in `{SOURCE_DOCS}`:

For each file, record:
- File path
- Total line count
- Major section headings

Calculate total source lines.

### Step 2: Inventory Migrated Documentation

Read ALL Markdown files in `{INSTALL_DIR}/`:

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
1. Search `{INSTALL_DIR}/` for the term using grep
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

{COMMON_AUDIT_CONTENT}
```

Where `{COMMON_AUDIT_CONTENT}` is replaced with the content from `common-audit.md`.

---

### Step 12: Copy Templates

For each selected layer, copy the appropriate template file:
- Decisions: `0000-template.md`
- Specs: `_template-feature.md` and/or `_template-rfc.md`
- Reference: `_template-cli.md` and/or `_template-config.md`
- Behaviors: `_template.feature`
- Implementation: `tech-stack.md`

---

### Step 13: Create Initial ADR (Optional)

If Layer 1 was selected, ask:

> Would you like to create an ADR documenting the decision to adopt this documentation system? (yes/no)

If yes, create `0001-adopt-layered-documentation-system.md`:
- **Context**: Documentation needs structure for maintainability
- **Decision**: Adopted Layered Documentation System
- **Consequences**: Clear separation of concerns, but requires discipline to maintain

---

### Step 14: Run Migration Audit (Migration Mode Only)

**Skip this step if `MIGRATION_MODE = false`.**

After migration, automatically run the audit workflow using the generated `{INSTALL_DIR}/audit.md` to compare original content against generated content:

#### 14a: Calculate Content Statistics

For original source:
- Count total lines
- Count lines per document
- Identify major sections

For generated output:
- Count total lines in layers (main content)
- Count total lines in appendices
- Count total lines in migration/
- Count total lines in blackhole/

#### 14b: Generate Comparison Summary

Create a detailed comparison showing where content went:

```
MIGRATION AUDIT SUMMARY
=======================

Source: {SOURCE_DOCS}
Generated: {INSTALL_DIR}
Date: {timestamp}

CONTENT DISTRIBUTION
--------------------
Total source lines: {N}

Successfully migrated to layers:
  - Main content: {N} lines ({X}%)
  - Appendices: {N} lines ({X}%)
  - Total: {N} lines ({X}%)

Needs review (migration/):
  - {N} lines ({X}%)
  - Categories: {list of migration subdirectories}

Unclassified (blackhole/):
  - {N} lines ({X}%)
  - Reason: {summary of classification failures}

LAYER BREAKDOWN
---------------
Layer 1 (Decisions): {N} ADRs created
Layer 2 (Vision): {N} lines
Layer 3 (Architecture): {N} lines
Layer 4 (Specs): {N} files, {N} lines
Layer 5 (Reference): {N} files, {N} lines (+ {N} lines in appendices)
Layer 7 (Implementation): {N} lines (+ {N} lines in appendices)

APPENDIX CONTENT
----------------
reference/appendices/cli/detailed-examples.md: {N} lines
reference/appendices/cli/error-messages.md: {N} lines
implementation/appendices/tech-stack/go-scaffolding/: {N} files, {N} lines
implementation/appendices/tech-stack/shell-scripts/: {N} files, {N} lines

MIGRATION DIRECTORY CONTENTS
----------------------------
{List each file in migration/ with line count and suggested destination}

BLACKHOLE CONTENTS
------------------
{List each file in blackhole/ with line count and suggested heuristic update}

THRESHOLDS USED
---------------
CODE_BLOCK_LINES: 50
STEP_LIST_ITEMS: 10
TABLE_ROWS: 20
...
```

#### 14c: Save Audit Report

Save the audit summary to `{INSTALL_DIR}/MIGRATION-AUDIT.md` for reference.

---

### Step 15: Final Report

**For Fresh Install**:

> **Installation Complete**
>
> Created the following structure:
> ```
> {directory tree}
> ```
>
> **Next Steps**:
> 1. Review `{INSTALL_DIR}/DOCUMENTATION-GUIDE.md` for usage instructions
> 2. Start with your Vision document in `{INSTALL_DIR}/product/vision.md`
> 3. Document key decisions as ADRs in `{INSTALL_DIR}/decisions/`
> 4. Use `create.md` when adding new documentation
> 5. Use `{INSTALL_DIR}/audit.md` to verify documentation completeness

**For Migration Install (After Install Phase — Steps 1-8)**:

> **Install Phase Complete — Migration Plan Generated**
>
> Created the following structure:
> ```
> {directory tree including .migration/}
> ```
>
> **Migration Plan Generated**:
> - **{N} step files** created in `{INSTALL_DIR}/.migration/`
> - **Estimated {M} total lines** to migrate across all steps
>
> **Content Classification Summary**:
> - High confidence (layers): {N} sections, ~{M} lines
> - Medium confidence (migration/): {N} sections, ~{M} lines
> - Low confidence (blackhole/): {N} sections, ~{M} lines
>
> **Next Steps — Execute Migration**:
> 1. Review the migration plan in `{INSTALL_DIR}/.migration/README.md`
> 2. Execute each step file in a separate session:
>    ```
>    Execute the migration step in @docs/.migration/01-decisions.md
>    ```
> 3. Mark steps complete as you go
> 4. After all steps: Run the audit (Step 14)
>
> **Migration Step Files**:
> | Step | File | Est. Lines |
> |------|------|------------|
> | 01 | `01-decisions.md` | ~{N} |
> | 02 | `02-vision.md` | ~{N} |
> | 03 | `03-architecture.md` | ~{N} |
> | 04 | `04-specs.md` | ~{N} |
> | 05 | `05-reference.md` | ~{N} |
> | 06 | `06-implementation.md` | ~{N} |
> | 07 | `07-medium-confidence.md` | ~{N} |
> | 08 | `08-low-confidence.md` | ~{N} |
> | 09 | `09-cross-links.md` | — |
>
> **Original files are unchanged** — source documentation is preserved for comparison.

**For Migration Install (After All Migration Steps Complete)**:

> **Migration Complete**
>
> All migration steps have been executed. Final structure:
> ```
> {directory tree}
> ```
>
> **Migration Summary**:
> - **{N} ADRs** extracted from embedded decisions
> - **Vision** consolidated from {N} source sections
> - **Architecture** consolidated from {N} source sections
> - **{N} Specification files** created
> - **{N} Reference files** created
> - **{N} Appendix files** created with detailed content
>
> **Content Distribution**:
> - Successfully migrated: {X}% of source content
> - In appendices: {X}% (detailed examples, code, etc.)
> - Needs review (`migration/`): {X}%
> - Unclassified (`blackhole/`): {X}%
>
> **Action Items**:
> 1. Review `{INSTALL_DIR}/MIGRATION-AUDIT.md` for detailed breakdown
> 2. Review content in `migration/` and move to appropriate layers
> 3. Review content in `blackhole/` and either:
>    - Move to appropriate layer/appendix
>    - Update heuristics in `DOCUMENTATION-GUIDE.md` for future
> 4. Check cross-layer links are correct
> 5. Delete the `.migration/` working directory:
>    ```bash
>    rm -rf {INSTALL_DIR}/.migration/
>    ```
> 6. Remove or archive original source files if desired
>
> **Original files are unchanged** — you can compare and verify before removing them.
>
> **To re-run audit later**: Use `{INSTALL_DIR}/audit.md`

---

## Quick Reference

### Layer → Directory Mapping

| Layer | Directory | Structure |
|-------|-----------|-----------|
| Decisions | `decisions/` | `NNNN-title.md` (simple single files, NOT directories) |
| Vision | `product/` | `vision.md` |
| Comparison | `product/` | `comparison.md` (auto-detected during migration) |
| Roadmap | `product/` | `roadmap.md` (auto-detected during migration) |
| Architecture | `architecture/` | `overview.md`, `{topic}.md` |
| Specifications | `specs/` | `{feature-name}.md` + `appendices/{feature-name}/` |
| Reference | `reference/` | `cli.md` + `appendices/cli/`, `configuration.md` + `appendices/configuration/` |
| Behaviors | `features/` (project root) | `{feature}.feature` |
| Implementation | `implementation/` | `tech-stack.md` + `appendices/tech-stack/` |
| Migration | `migration/` | `{category}/{content}.md` (temporary) |
| Blackhole | `blackhole/` | `{source}-lines-{N}-{M}.md` (unclassified) |

### Confidence → Destination Mapping

| Confidence | Destination |
|------------|-------------|
| High | Layer directly (main or appendix) |
| Medium | `migration/{category}/` |
| Low/None | `blackhole/` |

### Tooling Setup Commands

**Tier 1 (Essential)**:
```bash
npm install --save-dev markdownlint-cli
```

**Tier 2 (Recommended)**:
```bash
npm install --save-dev markdownlint-cli
brew install vale
npm install -g log4brains
```

**Tier 3 (Advanced)**:
```bash
# All of Tier 2, plus:
npm install --save-dev @cucumber/cucumber
```
