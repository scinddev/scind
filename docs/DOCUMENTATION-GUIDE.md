# Scind Documentation Guide

**For AI Agents and Contributors**: This guide explains how the documentation in this directory is organized and how to maintain it.

---

## Migration Note

This documentation was migrated from `specs/` on 2026-01-06 using the Layered Documentation System (LDS).

Migration step files are in `.migration/` — execute them in separate sessions to complete the migration.

---

## Glossary

This section defines key terms and concepts used throughout the documentation.

### Directory Terminology

| Term | Definition |
|------|------------|
| `DOCS_DIR` | The documentation root directory (`docs/`) — where all project documentation lives |
| `LEGACY_DOCS_DIR` | Source documentation being migrated (`specs/`) — the original documentation location |
| `LDS_DIST_DIR` | The Layered Documentation System distribution directory containing installation workflows and templates |

### Core System Terms

| Term | Definition |
|------|------------|
| **Layered Documentation System (LDS)** | A documentation framework that organizes software design documentation into seven distinct layers, each with specific purposes, ownership, and lifecycles |
| **Layer** | A distinct category of documentation with a specific purpose, stability level, audience, and lifecycle |
| **Operational Workflow** | A guided process that AI agents can execute to perform documentation tasks (e.g., `update.md`, `sync.md`) |
| **Fresh Install** | An installation mode that creates an empty documentation structure with templates but no migrated content |
| **Migration Install** | An installation mode that analyzes existing documentation and reorganizes it into the layered structure |

### Content Classification Terms

| Term | Definition |
|------|------------|
| **Confidence Level** | A rating (High, Medium, or Low) indicating how certain the classification heuristics are about where content belongs |
| **Classification Heuristics** | Pattern-matching rules used to determine which layer content belongs to |
| **Content Thresholds** | Configurable numeric limits that determine when content should be moved to an appendix |

### Document Structure Terms

| Term | Definition |
|------|------------|
| **Main Document** | The primary `{topic}.md` file in a layer directory containing overview and key information |
| **Appendix** | A supplementary file containing large content (code examples, detailed references) that exceeds thresholds, stored in `appendices/{topic}/` |
| **Appendix Directory** | The `appendices/` subdirectory within a layer, organized by topic matching main document names |
| **migration/ Directory** | A temporary directory for content classified with medium confidence that needs human review |
| **blackhole/ Directory** | A catch-all directory for content that could not be classified, indicating gaps in heuristics |

### Authority Terms

| Term | Definition |
|------|------------|
| **Document Hierarchy** | The order of authority when documents conflict — higher-authority documents win |
| **Canonical Source** | The single authoritative location where a piece of information is mastered |
| **Single Source of Truth (SSOT)** | The principle that each piece of information should live in exactly one place |
| **Supersede** | To replace an existing ADR with a new one; the old ADR's status becomes "Superseded by NNNN" |

### Canonical Layer Names

| Layer | Canonical Name | Aliases | Directory |
|-------|---------------|---------|-----------|
| 1 | **Decisions** | ADRs, Architectural Decision Records | `decisions/` |
| 2 | **Vision** | Product, PRD, PRD-Lite | `product/` |
| 3 | **Architecture** | (none) | `architecture/` |
| 4 | **Specifications** | Specs | `specs/` |
| 5 | **Reference** | (none) | `reference/` |
| 6 | **Behaviors** | Gherkin, Feature Files | `behaviors/` |
| 7 | **Implementation** | Implementation Guides, Scaffolding | `implementation/` |

### Layer-Specific Terms

| Term | Definition |
|------|------------|
| **ADR** | Architecture Decision Record — an immutable document capturing a significant technical or product decision |
| **MADR** | Markdown Any Decision Records — a template format for ADRs (available as "minimal" or "full" variants) |
| **Y-Statement** | A lightweight single-sentence ADR format for quick decision capture |
| **PRD / PRD-Lite** | Product Requirements Document — the Vision layer document defining product purpose, goals, and constraints |
| **Gherkin** | A domain-specific language for writing executable specifications using Given/When/Then syntax |
| **Feature File** | A `.feature` file containing Gherkin scenarios that serve as executable behavior specifications |

### Workflow Terms

| Term | Definition |
|------|------------|
| **Drift** | When documentation and implementation become out of sync |
| **Cross-Layer Link** | A relative Markdown link connecting documents in different layers |

---

## Content Thresholds (Configurable)

These thresholds control how content is classified. Edit to customize for this project.

| Threshold | Value | Purpose |
|-----------|-------|---------|
| `CODE_BLOCK_LINES` | 50 | Code blocks >= this go to appendix |
| `STEP_LIST_ITEMS` | 10 | Step lists >= this go to appendix |
| `TABLE_ROWS` | 20 | Tables >= this go to appendix |
| `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Complete file examples -> appendix |
| `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Error catalogs -> appendix |
| `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Shell scripts -> appendix |

**Rationale**: These values balance scannability with inline utility. Code blocks under 50 lines are typically readable in context; longer blocks benefit from dedicated appendix files. Complete files, error catalogs, and shell scripts are always appendices because they are reference material by nature.

---

## Core Principles

1. **Single Source of Truth**: Each piece of information lives in exactly one place
2. **Separation of Concerns**: Different document types serve different purposes
3. **Appropriate Stability**: Some documents are immutable; others evolve constantly
4. **Clear Ownership**: Each layer has defined maintainers and update triggers
5. **Linkage Over Duplication**: Reference other documents rather than copying content
6. **Preserve Full Content**: Migration means moving ALL technical details—not summarizing
7. **Appendix for Scale**: Large content lives in appendices, keeping main docs scannable
8. **Confidence-Based Fallback**: Content that can't be classified goes to `migration/` or `blackhole/`

---

## Layer Overview

This documentation system uses 7 layers, organized by purpose:

| Layer | Directory | Purpose | Stability | Audience |
|-------|-----------|---------|-----------|----------|
| 1. Decisions | `decisions/` | Capture WHY choices were made | Immutable | Future maintainers |
| 2. Vision | `product/` | Define WHAT we're building | Stable | All stakeholders |
| 3. Architecture | `architecture/` | Show HOW components relate | Evolving | Engineers, architects |
| 4. Specifications | `specs/` | Detail HOW features work | Living | Engineers |
| 5. Reference | `reference/` | Provide lookup tables | Generated/maintained | Engineers |
| 6. Behaviors | `behaviors/` | Verify expected behaviors | Executable | QA, engineers |
| 7. Implementation | `implementation/` | Guide HOW to build | Short-lived | Implementing engineers |

---

## Classification Decision Tree

When adding new content, use this tree to determine the correct layer:

```
Is this explaining WHY a choice was made?
├─ YES → Layer 1: Decisions (ADR)
└─ NO ↓

Is this about product vision, goals, or concepts?
├─ YES → Layer 2: Vision
└─ NO ↓

Is this showing how components relate (diagrams, flows)?
├─ YES → Layer 3: Architecture
└─ NO ↓

Is this detailing HOW a feature works (behavior, edge cases)?
├─ YES → Layer 4: Specifications
└─ NO ↓

Is this a lookup table (commands, options, configs)?
├─ YES → Layer 5: Reference
└─ NO ↓

Is this a concrete verifiable scenario (Given/When/Then)?
├─ YES → Layer 6: Behaviors
└─ NO ↓

Is this implementation scaffolding (code templates, dependencies)?
├─ YES → Layer 7: Implementation
└─ NO → May not need documentation
```

---

## Classification Heuristics

### Layer 1: Decisions (ADRs)

**Signals to look for**:
- "We chose X over Y because..."
- "We decided to..." / "Decision:"
- "This pattern will be used throughout..."
- "We considered alternatives including..."
- Trade-off discussions with rationale
- Explicit "Decision" blocks or sections

**ADR Structure**:
- Simple single files: `NNNN-{title}.md`
- No appendices for ADRs (include code inline)
- Status: Proposed | Accepted | Deprecated | Superseded

### Layer 2: Vision (PRD)

**Signals to look for**:
- Problem statements ("The problem we're solving...")
- Product vision or mission statements
- Success criteria or goals
- Non-goals or out-of-scope sections
- Core concepts and glossary

**Also in Layer 2**:
- `comparison.md` — How this compares to alternatives
- `roadmap.md` — Future considerations

**Comparison Document Signals**:
- Comparison tables with other tools/products
- "vs" language ("X vs Y", "compared to")
- Feature matrices comparing multiple products
- "Why choose X over Y" discussions

**Roadmap Document Signals**:
- "Future Considerations" sections
- "Roadmap" or "Future Work" headings
- "Planned features" or "upcoming" language
- "Phase 2", "v2.0" planning content

### Layer 3: Architecture

**Signals to look for**:
- System diagrams (ASCII art, Mermaid)
- Network topology descriptions
- Component relationship descriptions
- Cross-cutting concerns (security, logging patterns)
- "The system is composed of..."
- "Component A communicates with B via..."

### Layer 4: Specifications

**Signals to look for**:
- Detailed behavioral descriptions
- State machines or lifecycle descriptions
- Data schemas with field-level detail
- Edge case documentation
- Configuration file format specifications
- "When X happens, the system does Y..."
- "The valid states are..."

### Layer 5: Reference

**Signals to look for**:
- Command syntax documentation
- Option/flag tables
- Environment variable lists
- Configuration option tables with defaults
- Error code tables
- "The available options are..."
- "Syntax: command [options]..."

### Layer 6: Behaviors

**Signals to look for**:
- "Given/When/Then" scenarios
- Explicit test case descriptions
- User journey examples with expected outcomes
- "This behavior must not regress..."

**Signals that indicate other layers**:
- "The algorithm works by..." → Specification
- "Command options include..." → Reference

### Layer 7: Implementation

**Signals to look for**:
- Technology stack with versions
- Dependency lists with rationale
- Project scaffolding instructions
- Code templates
- "Install these dependencies..."
- "The project structure is..."

---

## Layer 6: Behaviors (Detailed Guidance)

### Purpose

Define expected behaviors in a way that can be automatically verified. Tests that double as documentation — **living documentation** that can never become outdated.

### Characteristics

- **Executable**: Run as part of CI/CD
- **Living**: Fail when behavior changes
- **Example-driven**: Concrete scenarios, not abstract
- **User-focused**: Written from user perspective

### The Living Documentation Advantage

Gherkin feature files serve triple duty:

1. **Specification** — Defines expected behavior before implementation
2. **Documentation** — Always accurate because it's tested
3. **Tests** — Executable validation that prevents regression

If a Gherkin test passes, the documentation is accurate. If behavior changes, the test fails, forcing documentation updates.

### What Belongs Here

- Critical user journeys
- Integration scenarios
- Edge case behaviors
- Regression-prevention tests

### When to Use

Use executable specs for:
- Behaviors that have historically broken
- Complex multi-step workflows
- Integration points between components
- Behaviors described in specifications that must not regress

### Directory Structure

Feature files are organized by domain:

```
behaviors/
├── README.md
├── {domain}/                    # e.g., workspace/, proxy/
│   ├── {feature}.feature
│   └── {another}.feature
└── support/                     # optional: step definitions
    └── step_definitions/
```

Example paths:
- `behaviors/workspace/workspace-lifecycle.feature`
- `behaviors/proxy/proxy-routing.feature`

### Template

```gherkin
Feature: [Feature Name]
  As a [role]
  I want [capability]
  So that [benefit]

  Background:
    Given [common precondition]

  Scenario: [Scenario Name]
    Given [initial context]
    When [action is taken]
    Then [expected outcome]

  # Link to specification
  # See: ../specs/{feature}.md
```

### Linking to Specifications

Every behavior file should reference the specification it verifies:

```gherkin
# This feature verifies behaviors from:
# See: ../specs/workspace-lifecycle.md

Feature: Workspace Lifecycle
  ...
```

---

## Cross-Layer Linking

Documents should link to related content in other layers:

| From | Link To | Purpose |
|------|---------|---------|
| Specifications | ADRs | Explain "why" for design choices |
| Specifications | Reference | Point to detailed syntax |
| Architecture | ADRs | Justify architectural patterns |
| Architecture | Specifications | Deep-dive into component behavior |
| Reference | Specifications | Provide conceptual context |
| Implementation | ADRs | Explain technology choices |
| Implementation | Specifications | Reference what's being implemented |
| Behaviors | Specifications | Reference the spec being verified |

### Link Format

```markdown
## Related Documents

- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md) — Design rationale
- [Proxy Specification](../specs/proxy-infrastructure.md) — Detailed behavior

## Appendices

- [Traefik Configuration](./appendices/proxy-infrastructure/traefik-config.yaml) — Complete config
```

### Inter-Appendix Linking

| From | To | Path Pattern |
|------|----|--------------|
| Same topic appendix | Same topic appendix | `./other-appendix.md` |
| One topic's appendix | Different topic's appendix | `../other-topic/appendix.md` |
| Appendix | Its main document | `../../main-doc.md` |

---

## Document Hierarchy (Authority Order)

When content appears in multiple places, this hierarchy determines the canonical source:

```
┌─────────────────────────────────────────────────────────────┐
│                      MOST AUTHORITATIVE                     │
├─────────────────────────────────────────────────────────────┤
│  ADRs (Architectural Decision Records)                      │
│  - Decisions are immutable once accepted                    │
│  - If anything conflicts with ADR, ADR wins                 │
├─────────────────────────────────────────────────────────────┤
│  Gherkin Feature Files                                      │
│  - Executable specifications                                │
│  - If test passes, documentation is accurate                │
├─────────────────────────────────────────────────────────────┤
│  Vision (PRD)                                               │
│  - High-level "what" and "why"                              │
├─────────────────────────────────────────────────────────────┤
│  Technical Specification                                    │
│  - Architecture and schemas                                 │
├─────────────────────────────────────────────────────────────┤
│  Reference Documentation (CLI, Config)                      │
│  - Factual, complete, lookup-oriented                       │
├─────────────────────────────────────────────────────────────┤
│  Implementation Guides (Tech Stack)                         │
│  - How to build, patterns to follow                         │
│                     LEAST AUTHORITATIVE                     │
└─────────────────────────────────────────────────────────────┘
```

For conflicts, defer to the higher-authority document.

---

## Tooling (Tier 2)

This project uses Tier 2 tooling:

| Tool | Purpose | Setup |
|------|---------|-------|
| **markdownlint** | Markdown linting | `npm install --save-dev markdownlint-cli` |
| **Vale** | Prose linting | `brew install vale` |
| **Log4brains** | ADR browsing | `npm install -g log4brains` |
| **Structurizr** | Architecture diagrams | Docker or cloud |

### Running Linters

```bash
# Markdown linting
npx markdownlint-cli2 "docs/**/*.md"

# Prose linting (requires .vale.ini)
vale docs/
```

---

## Directory Structure

```
docs/
├── README.md                    # Documentation index
├── DOCUMENTATION-GUIDE.md       # This file
│
├── decisions/                   # Layer 1: ADRs (simple files)
│   ├── README.md               # ADR index
│   ├── 0000-template.md        # Template
│   └── 0001-*.md ... 0012-*.md # ADR files
│
├── product/                     # Layer 2: Vision
│   ├── README.md
│   ├── vision.md
│   ├── glossary.md
│   ├── comparison.md
│   └── roadmap.md
│
├── architecture/                # Layer 3: Architecture
│   ├── README.md
│   └── overview.md
│
├── specs/                       # Layer 4: Specifications
│   ├── README.md
│   ├── {feature}.md            # Main spec files
│   └── appendices/
│       └── {feature}/          # Per-spec appendices
│
├── reference/                   # Layer 5: Reference
│   ├── README.md
│   ├── cli.md
│   ├── configuration.md
│   └── appendices/
│       ├── cli/
│       └── configuration/
│
├── behaviors/                   # Layer 6: Behaviors
│   ├── README.md
│   ├── {domain}/               # e.g., workspace/, proxy/
│   │   └── {feature}.feature
│   └── support/                # optional: step definitions
│       └── step_definitions/
│
├── implementation/              # Layer 7: Implementation
│   ├── README.md
│   ├── tech-stack.md
│   └── appendices/
│       └── tech-stack/
│
├── maintenance/                 # Maintenance workflows
│   ├── audit.md
│   ├── refine.md
│   ├── sync.md
│   └── update.md
│
└── .migration/                  # Migration step files (temporary)
    ├── README.md
    ├── common-instructions.md
    └── 01-decisions.md ... 09-cross-links.md
```

---

## Appendix Guidelines

### When to Use Appendices

Move content to appendices when it exceeds thresholds:
- Code blocks >= 50 lines
- Step lists >= 10 items
- Tables >= 20 rows
- Complete file examples (always)
- Error catalogs (always)
- Shell scripts (always)

**Exception**: ADRs never use appendices — all content including code stays inline.

### Appendix Directory Structure

```
specs/
├── shell-integration.md          # Main document (overview, key concepts)
└── appendices/
    └── shell-integration/        # Named after main document
        ├── bash-setup.sh         # Complete bash script
        ├── zsh-setup.zsh         # Complete zsh script
        └── fish-setup.fish       # Complete fish script
```

### Linking to Appendices

From main document:
```markdown
For complete shell scripts, see:
- [Bash Setup](./appendices/shell-integration/bash-setup.sh)
- [Zsh Setup](./appendices/shell-integration/zsh-setup.zsh)
```

From appendix (back-link):
```markdown
> **Parent**: [Shell Integration](../../shell-integration.md)
```

---

## Confidence-Based Classification

During migration or content creation, the system classifies content based on confidence level:

| Confidence | Destination | Meaning |
|------------|-------------|---------|
| **High** | Layer directly | Heuristics confident about both category AND placement |
| **Medium** | `migration/{category}/` | Category recognized but exact placement uncertain |
| **Low/None** | `blackhole/` | Heuristics cannot classify the content meaningfully |

### The `migration/` Directory

Content here has been categorized but needs human review for final placement. Review each file and:
1. Move to the appropriate layer if placement is now clear
2. Move to `appendices/` of the relevant document
3. Keep here if still uncertain (will be reviewed in next audit)

### The `blackhole/` Directory

Content here could not be classified at all. This represents gaps in the heuristics:
1. Review each file manually
2. Determine where it should go
3. **Update the heuristics** in this file to capture similar content in future

---

## Maintenance Workflows

Four workflows are available in `maintenance/`:

| Workflow | Purpose | When to Use |
|----------|---------|-------------|
| `audit.md` | Verify documentation completeness | After migration, periodically |
| `refine.md` | Improve quality without code changes | Documentation reviews |
| `sync.md` | Verify docs match implementation | Pre-release, after refactoring |
| `update.md` | Update docs after code changes | When implementation changes |

To execute a workflow:
```
Execute the workflow in @docs/maintenance/audit.md
```

---

## Single Source of Truth (SSOT)

Every fact should be mastered in exactly one place:

| Information Type | Canonical Source | Referenced From |
|------------------|------------------|-----------------|
| Why we chose X over Y | ADR | PRD, Specs, Architecture |
| Configuration schema | Tech Spec (or code) | Reference docs |
| Command syntax | CLI Reference | Tech Spec operations |
| Feature behavior | Specification | Gherkin tests |
| Port assignment algorithm | Specification | ADR (for rationale) |

**Key Principle**: When you find yourself copying information, create a reference instead.
