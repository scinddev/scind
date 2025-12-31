# Layered Documentation System — Installer

**For AI Agents**: This document contains instructions for installing and customizing the Layered Documentation System for a project. Read this entire file, then follow the process below.

---

## Prerequisites

Before starting, read the complete system documentation:

- `LAYERED-DOCUMENTATION-SYSTEM.md` — Full system guide with all layers, heuristics, and templates

---

## Installation Process

### Step 1: Determine Install Location

If the user specified an install directory (e.g., `@docs/`), use that path.

If no directory was specified, ask:

> Where would you like to install the documentation system? Common locations:
> - `docs/` — Standard documentation directory
> - `documentation/` — Alternative naming
> - `.docs/` — Hidden directory (less common)
>
> Please specify a path, or press enter for the default (`docs/`).

Store the chosen path as `INSTALL_DIR`.

---

### Step 2: Gather Layer Selections

Present each layer and its template options. For each layer, the user selects ONE template (or skips the layer entirely).

#### Layer 1: Decisions (ADRs)

Ask:

> **Layer 1: Decisions (ADRs)**
>
> ADRs capture *why* significant technical and product decisions were made. They become immutable once accepted.
>
> Template options:
> 1. **MADR Minimal** (Recommended) — Structured but lightweight. Good balance of rigor and ease.
> 2. **Y-Statement** — Single-sentence format. Best for quick capture of smaller decisions.
> 3. **MADR Full** — Comprehensive with pros/cons analysis. Best for complex, cross-cutting decisions.
> 4. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2/3/4)

Store selection as `DECISION_ADR_TEMPLATE`:
- `1` → `adr-madr-minimal.md`
- `2` → `adr-y-statement.md`
- `3` → `adr-madr-full.md`
- `4` → `null` (skip layer)

---

#### Layer 2: Vision (PRD-Lite)

Ask:

> **Layer 2: Vision (PRD-Lite)**
>
> Vision documents define the product's purpose, goals, and constraints. They change infrequently.
>
> Template options:
> 1. **Lean PRD** (Recommended) — Vision + problem + concepts + non-goals. Minimal structure.
> 2. **Epic-Based PRD** — Organized around user stories and epics. Better for agile teams.
> 3. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2/3)

Store selection as `DECISION_VISION_TEMPLATE`:
- `1` → `prd-lean.md`
- `2` → `prd-epic-based.md`
- `3` → `null` (skip layer)

---

#### Layer 3: Architecture

Ask:

> **Layer 3: Architecture**
>
> Architecture documents show how system components relate. They include diagrams and cross-cutting concerns.
>
> Template options:
> 1. **C4-Lite** (Recommended) — Context + Container diagrams with narrative. Lightweight.
> 2. **arc42-Simplified** — Structured sections from arc42 template. More comprehensive.
> 3. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2/3)

Store selection as `DECISION_ARCHITECTURE_TEMPLATE`:
- `1` → `architecture-c4-lite.md`
- `2` → `architecture-arc42.md`
- `3` → `null` (skip layer)

---

#### Layer 4: Specifications

Ask:

> **Layer 4: Specifications**
>
> Specifications detail how features work. They are living documents that evolve with implementation.
>
> Template options:
> 1. **Feature Spec** (Recommended) — Narrative + examples + edge cases. Good for most features.
> 2. **RFC-Style** — Problem + proposal + alternatives. Good for proposed changes.
> 3. **Both** — Include both templates for different use cases.
> 4. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2/3/4)

Store selection as `DECISION_SPEC_TEMPLATE`:
- `1` → `spec-feature.md`
- `2` → `spec-rfc.md`
- `3` → `both`
- `4` → `null` (skip layer)

---

#### Layer 5: Reference

Ask:

> **Layer 5: Reference**
>
> Reference documents provide lookup information (CLI commands, configuration options, etc.).
>
> Template options:
> 1. **CLI Reference** — Command documentation with options tables.
> 2. **Configuration Reference** — Config file schema documentation.
> 3. **Both** — Include both templates.
> 4. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2/3/4)

Store selection as `DECISION_REFERENCE_TEMPLATE`:
- `1` → `reference-cli.md`
- `2` → `reference-config.md`
- `3` → `both`
- `4` → `null` (skip layer)

---

#### Layer 6: Behaviors (Executable Specifications)

Ask:

> **Layer 6: Behaviors (Executable Specifications)**
>
> Behavior specs define expected behaviors in Gherkin format that can be verified automatically.
>
> Template options:
> 1. **Gherkin** — Standard BDD format with Cucumber/Behave support.
> 2. **Skip** — Do not include this layer.
>
> Which template would you like? (1/2)

Store selection as `DECISION_BEHAVIOR_TEMPLATE`:
- `1` → `behavior-gherkin.feature`
- `2` → `null` (skip layer)

---

### Step 3: Gather Tooling Selections

Ask:

> **Tooling Tier Selection**
>
> Which level of tooling would you like documented?
>
> 1. **Tier 1: Essential** — markdownlint, Git + PRs only. Minimal setup.
> 2. **Tier 2: Recommended** — Adds Vale (prose linting), Log4brains (ADR site), Structurizr (diagrams).
> 3. **Tier 3: Advanced** — Adds Cucumber/Gherkin, Semcheck (AI drift detection), custom generators.
> 4. **None** — No tooling recommendations.
>
> Which tier? (1/2/3/4)

Store selection as `DECISION_TOOLING_TIER`:
- `1` → `essential`
- `2` → `recommended`
- `3` → `advanced`
- `4` → `none`

---

### Step 4: Confirm Selections

Present a summary of all selections:

> **Installation Summary**
>
> **Install Location**: `{INSTALL_DIR}`
>
> **Layers to Include**:
> - Decisions: {template name or "Skipped"}
> - Vision: {template name or "Skipped"}
> - Architecture: {template name or "Skipped"}
> - Specifications: {template name or "Skipped"}
> - Reference: {template name or "Skipped"}
> - Behaviors: {template name or "Skipped"}
>
> **Tooling Tier**: {tier name}
>
> Proceed with installation? (yes/no)

If "no", ask which selections to change and loop back to the relevant step.

---

### Step 5: Create Directory Structure

Create the following directories based on selected layers:

```
{INSTALL_DIR}/
├── DOCUMENTATION-GUIDE.md          # Always created (customized system guide)
├── decisions/                       # If Layer 1 selected
│   └── 0000-template.md            # Template file
├── product/                         # If Layer 2 selected
│   └── vision.md                   # Template file (renamed)
├── architecture/                    # If Layer 3 selected
│   └── overview.md                 # Template file (renamed)
├── specs/                           # If Layer 4 selected
│   └── _template-feature.md        # Template file(s), prefixed with _
│   └── _template-rfc.md            # (if both selected)
├── reference/                       # If Layer 5 selected
│   └── _template-cli.md            # Template file(s), prefixed with _
│   └── _template-config.md         # (if both selected)
└── features/                        # If Layer 6 selected (at project root, not in docs)
    └── _template.feature           # Template file
```

**Note**: Layer 6 (Behaviors) should be created at the project root level as `features/`, not inside the docs directory, since these are typically run by test frameworks that expect them in a standard location.

---

### Step 6: Generate Customized Guide

Create `{INSTALL_DIR}/DOCUMENTATION-GUIDE.md` containing:

1. **Header** with project-specific title
2. **Layer Overview** — Only include sections for selected layers
3. **Classification Heuristics** — Only for selected layers
4. **Decision Tree** — Simplified to only show selected layers
5. **Cross-Layer Linking** — Adjusted for selected layers
6. **Tooling Section** — Only include selected tier
7. **Directory Structure** — Reflect actual created structure

Use the content from `LAYERED-DOCUMENTATION-SYSTEM.md` as the source, but:
- Remove all sections for skipped layers
- Remove template options that weren't selected
- Update the "Files in This System" section to match actual files
- Simplify the decision tree to only branch to included layers

---

### Step 7: Copy Selected Templates

For each selected template:

1. Read the template from `templates/{filename}`
2. Copy to the appropriate location per Step 5
3. For template files that will be used as starting points (vision.md, overview.md), keep the template content as-is
4. For template files in `specs/` and `reference/`, prefix with `_template-` to indicate they are templates, not actual docs

---

### Step 8: Create Initial ADR (Optional)

If Layer 1 (Decisions) was selected, ask:

> Would you like to create an initial ADR documenting the decision to adopt this documentation system? (yes/no)

If yes, create `{INSTALL_DIR}/decisions/0001-adopt-layered-documentation.md` using the selected ADR template, pre-filled with:

- **Context**: The project needed a structured approach to documentation
- **Decision**: Adopted the Layered Documentation System with [selected layers]
- **Consequences**:
  - Positive: Clear separation of concerns, single source of truth
  - Negative: Initial setup overhead, learning curve

---

### Step 9: Final Report

Report to the user:

> **Installation Complete**
>
> Created the following structure:
> ```
> {show actual directory tree}
> ```
>
> **Next Steps**:
> 1. Review `{INSTALL_DIR}/DOCUMENTATION-GUIDE.md` for usage instructions
> 2. Start with your Vision document in `{INSTALL_DIR}/product/vision.md`
> 3. Document key decisions as ADRs in `{INSTALL_DIR}/decisions/`
> 4. {If tooling selected: Set up recommended tooling per the guide}
>
> **Tips**:
> - Template files prefixed with `_template-` should be copied and renamed when creating new docs
> - Use the classification heuristics in the guide to determine where content belongs
> - Link between layers rather than duplicating content

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

### Tooling Setup Commands

**Tier 1 (Essential)**:
```bash
npm install --save-dev markdownlint-cli
# Add to package.json scripts: "lint:docs": "markdownlint docs/**/*.md"
```

**Tier 2 (Recommended)**:
```bash
npm install --save-dev markdownlint-cli
brew install vale  # or: choco install vale
npm install -g log4brains
# See https://structurizr.com for diagram tooling
```

**Tier 3 (Advanced)**:
```bash
# All of Tier 2, plus:
npm install --save-dev @cucumber/cucumber
# See https://github.com/semcheck/semcheck for AI drift detection
```
