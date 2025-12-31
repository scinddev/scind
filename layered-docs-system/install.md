# Layered Documentation System — Installer

**For AI Agents**: This document contains instructions for installing and customizing the Layered Documentation System for a project. Read this entire file, then follow the process below.

---

## Prerequisites

Before starting, read the complete system documentation:

- `LAYERED-DOCUMENTATION-SYSTEM.md` — Full system guide with all layers, heuristics, and templates

---

## Installation Modes

This installer supports two modes:

1. **Fresh Install** — Create an empty documentation system with templates
2. **Migration Install** — Analyze existing documentation and reorganize into layers

The mode is determined by whether existing documentation is specified.

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

#### 3c: Create Content Mapping

Build a mapping of source content to target layers:

```
CONTENT_MAP = {
  "Layer 1: Decisions": [
    { source: "prd.md", section: "Decision: Docker Compose", lines: "45-62" },
    { source: "prd.md", section: "Decision: Two-Layer Networking", lines: "78-95" },
    ...
  ],
  "Layer 2: Vision": [
    { source: "prd.md", section: "Executive Summary", lines: "1-20" },
    { source: "prd.md", section: "Problem Statement", lines: "22-44" },
    { source: "prd.md", section: "Core Concepts", lines: "100-150" },
    ...
  ],
  "Layer 3: Architecture": [
    { source: "technical-spec.md", section: "Architecture", lines: "1-80" },
    { source: "prd.md", section: "Network Topology", lines: "160-200" },
    ...
  ],
  "Layer 4: Specifications": [
    { source: "technical-spec.md", section: "Port Types", lines: "85-150" },
    { source: "go-stack.md", section: "All content", lines: "1-500" },
    ...
  ],
  "Layer 5: Reference": [
    { source: "cli-reference.md", section: "All content", lines: "1-300" },
    { source: "technical-spec.md", section: "Environment Variables", lines: "200-250" },
    ...
  ],
  "Layer 6: Behaviors": [
    // Usually empty unless source has Gherkin or explicit test scenarios
  ],
  "Unclassified": [
    // Content that doesn't clearly fit any layer
  ]
}
```

#### 3d: Present Analysis to User

Present the content mapping for review:

> **Content Analysis Complete**
>
> Based on the classification heuristics, here's how your existing content maps to layers:
>
> **Layer 1: Decisions** — {N} sections identified
> - Decision blocks from `prd.md`
> - Architecture rationale from `technical-spec.md`
>
> **Layer 2: Vision** — {N} sections identified
> - Executive summary, problem statement, concepts from `prd.md`
>
> **Layer 3: Architecture** — {N} sections identified
> - Network topology diagrams
> - Component overview sections
>
> **Layer 4: Specifications** — {N} sections identified
> - Detailed behavioral specs from `technical-spec.md`
> - Implementation details from `go-stack.md`
>
> **Layer 5: Reference** — {N} sections identified
> - CLI documentation from `cli-reference.md`
> - Configuration tables
>
> **Layer 6: Behaviors** — {N} sections identified
> - {List or "None found — you may want to add these later"}
>
> **Unclassified** — {N} sections
> - {List any content that doesn't clearly fit}
>
> Does this mapping look correct? (yes/no/adjust)

If "adjust", allow the user to manually reassign sections.

---

### Step 4: Gather Layer Selections

**For Fresh Install**: Present all options as before.

**For Migration Install**: Pre-select layers based on content analysis. Only layers with identified content are pre-selected, but user can add additional empty layers.

#### Layer 1: Decisions (ADRs)

**Migration mode** (if decisions were found):
> **Layer 1: Decisions (ADRs)**
>
> Found {N} decision sections in your existing docs. These will be extracted as individual ADRs.
>
> Template for NEW decisions:
> 1. **MADR Minimal** (Recommended)
> 2. **Y-Statement**
> 3. **MADR Full**
>
> Which template for future ADRs? (1/2/3)

**Fresh install or no decisions found**:
> **Layer 1: Decisions (ADRs)**
>
> ADRs capture *why* significant technical and product decisions were made. They become immutable once accepted.
>
> Template options:
> 1. **MADR Minimal** (Recommended) — Structured but lightweight.
> 2. **Y-Statement** — Single-sentence format.
> 3. **MADR Full** — Comprehensive with pros/cons analysis.
> 4. **Skip** — Do not include this layer.
>
> Which template? (1/2/3/4)

Store selection as `DECISION_ADR_TEMPLATE`.

---

#### Layer 2: Vision (PRD-Lite)

**Migration mode** (if vision content found):
> **Layer 2: Vision (PRD-Lite)**
>
> Found vision/product content in your existing docs. This will be reorganized into the Vision layer.
>
> Template style:
> 1. **Lean PRD** (Recommended) — Vision + problem + concepts + non-goals.
> 2. **Epic-Based PRD** — Organized around user stories and epics.
>
> Which template style? (1/2)

**Fresh install or no vision content found**:
> **Layer 2: Vision (PRD-Lite)**
>
> Vision documents define the product's purpose, goals, and constraints.
>
> Template options:
> 1. **Lean PRD** (Recommended)
> 2. **Epic-Based PRD**
> 3. **Skip** — Do not include this layer.
>
> Which template? (1/2/3)

Store selection as `DECISION_VISION_TEMPLATE`.

---

#### Layer 3: Architecture

**Migration mode** (if architecture content found):
> **Layer 3: Architecture**
>
> Found architecture content (diagrams, component descriptions). This will be reorganized into the Architecture layer.
>
> Template style:
> 1. **C4-Lite** (Recommended) — Context + Container diagrams with narrative.
> 2. **arc42-Simplified** — Structured sections from arc42.
>
> Which template style? (1/2)

**Fresh install or no architecture content found**:
> **Layer 3: Architecture**
>
> Template options:
> 1. **C4-Lite** (Recommended)
> 2. **arc42-Simplified**
> 3. **Skip** — Do not include this layer.
>
> Which template? (1/2/3)

Store selection as `DECISION_ARCHITECTURE_TEMPLATE`.

---

#### Layer 4: Specifications

**Migration mode** (if spec content found):
> **Layer 4: Specifications**
>
> Found detailed specification content. This will be reorganized into the Specifications layer.
>
> Template for NEW specs:
> 1. **Feature Spec** (Recommended)
> 2. **RFC-Style**
> 3. **Both**
>
> Which template(s)? (1/2/3)

**Fresh install or no spec content found**:
> **Layer 4: Specifications**
>
> Template options:
> 1. **Feature Spec** (Recommended)
> 2. **RFC-Style**
> 3. **Both**
> 4. **Skip** — Do not include this layer.
>
> Which template? (1/2/3/4)

Store selection as `DECISION_SPEC_TEMPLATE`.

---

#### Layer 5: Reference

**Migration mode** (if reference content found):
> **Layer 5: Reference**
>
> Found reference documentation (CLI, configuration, etc.). This will be reorganized into the Reference layer.
>
> Template for NEW reference docs:
> 1. **CLI Reference**
> 2. **Configuration Reference**
> 3. **Both**
>
> Which template(s)? (1/2/3)

**Fresh install or no reference content found**:
> **Layer 5: Reference**
>
> Template options:
> 1. **CLI Reference**
> 2. **Configuration Reference**
> 3. **Both**
> 4. **Skip** — Do not include this layer.
>
> Which template? (1/2/3/4)

Store selection as `DECISION_REFERENCE_TEMPLATE`.

---

#### Layer 6: Behaviors (Executable Specifications)

> **Layer 6: Behaviors (Executable Specifications)**
>
> Behavior specs define expected behaviors in Gherkin format.
> {If migration: "No Gherkin content was found in your existing docs."}
>
> Template options:
> 1. **Gherkin** — Standard BDD format.
> 2. **Skip** — Do not include this layer.
>
> Which template? (1/2)

Store selection as `DECISION_BEHAVIOR_TEMPLATE`.

---

### Step 5: Gather Tooling Selections

Ask:

> **Tooling Tier Selection**
>
> Which level of tooling would you like documented?
>
> 1. **Tier 1: Essential** — markdownlint, Git + PRs only.
> 2. **Tier 2: Recommended** — Adds Vale, Log4brains, Structurizr.
> 3. **Tier 3: Advanced** — Adds Cucumber, Semcheck, custom generators.
> 4. **None** — No tooling recommendations.
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
>
> **Tooling Tier**: {tier}
>
> Proceed with installation? (yes/no)

If "no", ask which selections to change.

---

### Step 7: Create Directory Structure

Create directories based on selected layers:

```
{INSTALL_DIR}/
├── DOCUMENTATION-GUIDE.md          # Always created
├── decisions/                       # If Layer 1 selected
│   └── 0000-template.md            # Template file
├── product/                         # If Layer 2 selected
│   └── vision.md                   # Template or migrated content
├── architecture/                    # If Layer 3 selected
│   └── overview.md                 # Template or migrated content
├── specs/                           # If Layer 4 selected
│   └── _template-feature.md        # Template file(s)
├── reference/                       # If Layer 5 selected
│   └── _template-cli.md            # Template file(s)
└── features/                        # If Layer 6 selected (project root)
    └── _template.feature           # Template file
```

---

### Step 8: Migrate Content (Migration Mode Only)

**Skip this step if `MIGRATION_MODE = false`.**

For each layer with mapped content:

#### 8a: Extract Decisions as Individual ADRs

For each decision identified in source docs:
1. Create a new ADR file: `{INSTALL_DIR}/decisions/NNNN-{slugified-title}.md`
2. Use the selected ADR template structure
3. Fill in:
   - **Context**: Extract surrounding context from source
   - **Decision**: Extract the decision statement
   - **Consequences**: Extract or infer consequences
4. Add a note: `<!-- Migrated from {source_file}:{lines} -->`

Number ADRs sequentially starting from `0001`.

#### 8b: Consolidate Vision Content

Combine all vision-related content into `{INSTALL_DIR}/product/vision.md`:
1. Use the selected PRD template as the structure
2. Map source content to template sections:
   - Executive Summary → Overview section
   - Problem Statement → Problem Statement section
   - Core Concepts → Core Concepts section
   - Non-goals → Non-Goals section
   - Glossary → Glossary section
3. Preserve original wording where possible
4. Add section comments: `<!-- Source: {file}:{lines} -->`
5. Create `concepts.md` if content is extensive

#### 8c: Consolidate Architecture Content

Combine architecture content into `{INSTALL_DIR}/architecture/overview.md`:
1. Use the selected architecture template as the structure
2. Map source content to template sections:
   - Diagrams → Appropriate diagram sections
   - Component descriptions → Building Blocks section
   - Network topology → dedicated section or integrate into diagrams
3. Preserve diagrams exactly as-is
4. Add source comments

Create additional files if needed:
- `networking.md` for detailed network documentation
- `{topic}.md` for other major architectural areas

#### 8d: Organize Specification Content

For specification content, decide on file organization:

1. **If source is already well-organized by feature**: Create one spec file per source section
   - `{INSTALL_DIR}/specs/{feature-name}.md`

2. **If source is monolithic**: Split by logical feature boundaries
   - Identify natural feature boundaries in the content
   - Create separate spec files for each feature area

3. Use the selected spec template as the structure for each file
4. Preserve detailed behavioral descriptions, schemas, examples
5. Add source comments
6. Remove any decision rationale (that's now in ADRs — add links instead)

#### 8e: Organize Reference Content

For reference content:

1. **CLI Reference**: Consolidate into `{INSTALL_DIR}/reference/cli.md`
2. **Configuration**: Consolidate into `{INSTALL_DIR}/reference/configuration.md`
3. **Environment Variables**: `{INSTALL_DIR}/reference/environment.md`
4. **Other reference tables**: Create appropriate files

Use selected reference template structure. Preserve tables and syntax exactly.

#### 8f: Add Cross-Layer Links

After migration, scan all created files and add links:

1. In specifications, link to relevant ADRs:
   ```markdown
   See [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md) for rationale.
   ```

2. In architecture, link to specifications for detail:
   ```markdown
   For detailed port assignment behavior, see [Port Assignment Spec](../specs/port-assignment.md).
   ```

3. In reference, link to specifications for context:
   ```markdown
   See [Workspace Lifecycle](../specs/workspace-lifecycle.md) for behavioral details.
   ```

---

### Step 9: Generate Customized Guide

Create `{INSTALL_DIR}/DOCUMENTATION-GUIDE.md` containing:

1. **Header** with project-specific title
2. **Migration Note** (if applicable): "This documentation was migrated from `{SOURCE_DOCS}` on {date}"
3. **Layer Overview** — Only selected layers
4. **Classification Heuristics** — Only for selected layers
5. **Decision Tree** — Simplified for included layers
6. **Cross-Layer Linking** — Adjusted for selected layers
7. **Tooling Section** — Selected tier only
8. **Directory Structure** — Actual created structure

---

### Step 10: Copy Templates

For each selected layer, copy the template file:
- Decisions: `0000-template.md` (always include for new ADRs)
- Specs: `_template-feature.md` and/or `_template-rfc.md`
- Reference: `_template-cli.md` and/or `_template-config.md`
- Behaviors: `_template.feature`

---

### Step 11: Create Initial ADR (Optional)

If Layer 1 was selected, ask:

> Would you like to create an ADR documenting the decision to adopt this documentation system? (yes/no)

If yes, create the ADR. For migration mode, include:
- **Context**: Existing documentation was scattered/overlapping
- **Decision**: Migrated to Layered Documentation System
- **Consequences**: Clear separation, but required reorganization effort

---

### Step 12: Final Report

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

**For Migration Install**:

> **Migration Complete**
>
> Migrated content from `{SOURCE_DOCS}` into:
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
>
> **Files Created**:
> - `decisions/0001-{title}.md`
> - `decisions/0002-{title}.md`
> - ...
> - `product/vision.md`
> - `architecture/overview.md`
> - `specs/{name}.md`
> - ...
>
> **Action Items**:
> 1. Review migrated content for accuracy
> 2. Check cross-layer links are correct
> 3. Remove or archive original source files if desired
> 4. Fill in any template sections that couldn't be auto-populated
>
> **Original files are unchanged** — you can compare and verify before removing them.

---

## Customization Reference

### Layer → Directory Mapping

| Layer | Directory | Template Naming |
|-------|-----------|-----------------|
| Decisions | `decisions/` | `NNNN-title.md` |
| Vision | `product/` | `vision.md`, `concepts.md` |
| Architecture | `architecture/` | `overview.md`, `{topic}.md` |
| Specifications | `specs/` | `{feature-name}.md` |
| Reference | `reference/` | `cli.md`, `config.md`, etc. |
| Behaviors | `features/` (project root) | `{feature}.feature` |

### Content Classification Quick Reference

| If you see... | It belongs in... |
|---------------|------------------|
| "We decided...", "We chose X over Y" | Layer 1: Decisions |
| Problem statement, vision, goals, non-goals | Layer 2: Vision |
| Diagrams, component relationships, topology | Layer 3: Architecture |
| "When X happens, Y occurs", state machines, schemas | Layer 4: Specifications |
| Command tables, option lists, config defaults | Layer 5: Reference |
| Given/When/Then, test scenarios | Layer 6: Behaviors |

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
