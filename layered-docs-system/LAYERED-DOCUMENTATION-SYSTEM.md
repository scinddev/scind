# Layered Documentation System

**Version**: 3.0.0
**Purpose**: A cohesive, maintainable system for creating and managing software documentation using a layered approach with operational workflows for AI agents.

**What's New in 3.0**:
- Appendix support for large content (code blocks, examples, detailed references)
- Confidence-based classification with `migration/` and `blackhole/` fallbacks
- Configurable content thresholds
- Automatic migration comparison and loss analysis
- Directory-based structure with `README.md` + `appendices/` pattern

---

## Quick Start

### For AI Agents

**Fresh install** (empty templates):
```
Execute @path/to/layered-docs-system/install.md and install to @docs/
```

**Migration install** (reorganize existing docs):
```
Execute @path/to/layered-docs-system/install.md and install to @docs/ from existing docs in @specs/
```

> **Note**: For large documentation sets, the install phase generates migration step files in `docs/.migration/` that can be executed in separate sessions. This avoids context window exhaustion. After the install phase, execute each step file individually.

**Interactive** (prompts for all options):
```
Execute @path/to/layered-docs-system/install.md
```

### Operational Workflows

| Task | Guide | When to Use |
|------|-------|-------------|
| Set up documentation | `install.md` | Project initialization |
| Add new feature/decision | `create.md` | Adding new documentation |
| Code changed | `update.md` | After implementation changes |
| Periodic audit | `sync.md` | Before releases, periodic maintenance |
| Quality improvement | `refine.md` | Documentation review cycles |
| Compare docs vs source | `audit.md` | After migration, or to check for drift |

### For Humans

1. Read this document to understand the system
2. Use the [Quick Start Checklist](#quick-start-checklist) to set up manually
3. Or use the operational guides with an AI assistant

---

## Overview

This system organizes documentation into seven distinct layers, each with a specific purpose, ownership model, and lifecycle. The goal is to prevent drift between documents by ensuring each piece of information has exactly one canonical source.

**Note**: This entire system is **design documentation for developers**—not end-user documentation. The layers separate different types of design artifacts: Layers 1-6 describe *what to build*, while Layer 7 describes *how to build it*.

### The Seven Layers

| Layer | Purpose | Stability | Audience |
|-------|---------|-----------|----------|
| **1. Decisions** | Capture *why* choices were made | Immutable once accepted | Future maintainers, architects |
| **2. Vision** | Define *what* we're building and *why* | Stable, rarely changes | All stakeholders |
| **3. Architecture** | Show *how* components relate | Evolves with design | Engineers, architects |
| **4. Specifications** | Detail *how* features work | Living, evolves with implementation | Engineers |
| **5. Reference** | Provide lookup information | Generated or hand-maintained | Engineers |
| **6. Behaviors** | Verify system meets expectations | Executable, tied to tests | QA, engineers |
| **7. Implementation** | Describe *how to build* the system | Short-lived, absorbed into code | Engineers implementing |

### Core Principles

1. **Single Source of Truth**: Each piece of information lives in exactly one place
2. **Separation of Concerns**: Different document types serve different purposes
3. **Appropriate Stability**: Some documents are immutable; others evolve constantly
4. **Clear Ownership**: Each layer has defined maintainers and update triggers
5. **Linkage Over Duplication**: Reference other documents rather than copying content
6. **Preserve Full Content**: Migration means moving ALL technical details, code examples, error messages, and output samples—not summarizing them
7. **Appendix for Scale**: Large content lives in appendices, keeping main docs scannable (but see layer-specific guidance—ADRs rarely need appendices)
8. **Confidence-Based Fallback**: Content that can't be classified goes to `migration/` or `blackhole/`

---

## Directory Structure

Layers use a flat file structure where main documents are siblings in the layer directory, with appendices organized in a shared `appendices/` directory.

### Standard Pattern

```
docs/{layer}/
├── README.md              # Auto-generated layer index
├── {topic}.md             # Main document for topic
├── {another-topic}.md     # Main document for another topic
└── appendices/            # Shared appendices directory
    ├── {topic}/           # Appendices for {topic}.md
    │   ├── examples.md
    │   └── code-samples.md
    └── {another-topic}/   # Appendices for {another-topic}.md
        └── error-catalog.md
```

### Examples

```
docs/
├── DOCUMENTATION-GUIDE.md     # Project-specific guide (generated from install)
│
├── decisions/                 # Layer 1: ADRs (simple files, not directories)
│   ├── README.md              # Index of all ADRs
│   ├── 0001-docker-compose-isolation.md   # Single file per ADR
│   └── 0002-networking-model.md           # Include code examples inline
│
├── reference/                 # Layer 5: Reference
│   ├── README.md              # Auto-generated index
│   ├── cli.md                 # CLI reference (commands, flags, descriptions)
│   ├── configuration.md       # Config reference (options, defaults)
│   └── appendices/
│       ├── cli/                       # Appendices for cli.md
│       │   ├── detailed-examples.md   # Complete workflow examples
│       │   ├── error-messages.md      # Full error message catalog
│       │   └── output-formats.md      # Detailed output examples
│       └── configuration/             # Appendices for configuration.md
│           └── complete-schemas.md    # Full YAML/JSON schema examples
│
├── implementation/            # Layer 7: Implementation Guides
│   ├── README.md              # Auto-generated index
│   ├── tech-stack.md          # Tech stack overview
│   └── appendices/
│       └── tech-stack/                 # Appendices for tech-stack.md
│           ├── go-scaffolding/         # Complete code scaffolding
│           │   ├── cmd-main.go.md
│           │   ├── cli-root.go.md
│           │   └── cli-commands.go.md
│           └── shell-scripts/          # Complete shell scripts
│               ├── bash.sh.md
│               ├── zsh.zsh.md
│               └── fish.fish.md
│
├── migration/                 # Content classified but placement uncertain
│   ├── README.md              # Explains what's here and why
│   └── {category}/            # Organized by detected category
│       └── {content}.md
│
└── blackhole/                 # Content that couldn't be classified at all
    ├── README.md              # Explains what's here and why
    └── {source-filename}.md   # Raw content dump with source attribution
```

### Layer README.md Files

Each layer has an auto-generated `README.md` that serves as an index. This file:
- Lists all documents in the layer
- Provides brief guidance on when to create documents of this type
- Links to `DOCUMENTATION-GUIDE.md` for authoritative instructions

### Appendix Naming Convention

Appendices are stored in `appendices/{basename}/` where `{basename}` matches the main document's filename (without `.md`):
- `cli.md` → `appendices/cli/`
- `tech-stack.md` → `appendices/tech-stack/`
- `workspace-lifecycle.md` → `appendices/workspace-lifecycle/`

---

## Appendix System

### Purpose

Appendices keep main documents scannable while preserving valuable detailed content:
- Large code blocks (scaffolding, complete scripts)
- Detailed examples (multi-step workflows)
- Complete reference tables (all error messages, all options)
- Full schema definitions (complete YAML/JSON examples)

### When Content Goes to Appendix

Content moves to an appendix when it exceeds configured thresholds (see [Content Thresholds](#content-thresholds)) or matches specific patterns:

| Content Type | Main Document | Appendix |
|--------------|---------------|----------|
| Code blocks | < threshold lines | ≥ threshold lines |
| Complete file examples | Never | Always |
| Step-by-step instructions | < threshold steps | ≥ threshold steps |
| Error message catalogs | Summary/key errors | Complete catalog |
| Schema examples | Excerpt showing structure | Complete schema |
| Shell scripts | Synopsis/key functions | Full script |

### Appendix Document Structure

Each appendix file should be self-contained and reference-able:

```markdown
# [Appendix Title]

> **Parent**: [Link to main document](../../{topic}.md)
> **Purpose**: [What this appendix contains]

## Content

[The detailed content]

---

*This appendix supports [{Topic} Document](../../{topic}.md).*
```

### Linking to Appendices

From main documents, link to appendices for detail:

```markdown
## Command Examples

For basic usage, see the examples below. For complete workflow examples
including error handling, see [Detailed Examples](./appendices/cli/detailed-examples.md).

### Basic Example

\`\`\`bash
contrail workspace up
\`\`\`
```

Note the path structure: from `cli.md`, appendices are at `./appendices/cli/`.

---

## Confidence-Based Classification

During migration or content creation, the system classifies content based on confidence level:

### Classification Destinations

| Confidence | Destination | Meaning |
|------------|-------------|---------|
| **High** | Layer directly | Heuristics are confident about both category AND placement |
| **Medium** | `migration/{category}/` | Heuristics recognize the category but uncertain about exact placement |
| **Low/None** | `blackhole/` | Heuristics cannot classify the content meaningfully |

### The `migration/` Directory

Content in `migration/` has been categorized but needs human review for final placement:

```
migration/
├── README.md                  # Explains what's here
├── shell-integration/         # Detected as shell-related
│   └── complete-scripts.md    # Full bash/zsh/fish scripts
├── cli-reference/             # Detected as CLI-related
│   └── workflow-examples.md   # Detailed workflow examples
└── implementation/            # Detected as implementation-related
    └── scaffolding-code.md    # Go scaffolding code
```

**README.md content**:
```markdown
# Migration Content

This directory contains content that was classified during migration but
couldn't be placed with high confidence. Review each file and:

1. Move to the appropriate layer if placement is now clear
2. Move to `appendices/` of the relevant document
3. Keep here if still uncertain (will be reviewed in next audit)

## Contents

| File | Detected Category | Suggested Destination |
|------|-------------------|----------------------|
| shell-integration/complete-scripts.md | Shell Integration | `implementation/tech-stack/appendices/shell-scripts/` |
| ... | ... | ... |
```

### The `blackhole/` Directory

Content in `blackhole/` could not be classified at all. This represents a gap in the heuristics:

```
blackhole/
├── README.md                  # Explains what's here and suggests heuristic improvements
└── contrail-go-stack-lines-450-600.md  # Raw content with source attribution
```

**README.md content**:
```markdown
# Blackhole Content

This directory contains content that the classification heuristics could not
process. This represents gaps in the system that should be addressed.

## What This Means

Content ends up here when:
- The heuristics couldn't determine even a category
- The content format wasn't recognized
- The content doesn't match any known patterns

## Action Required

1. Review each file manually
2. Determine where it should go (layer, migration, or truly orphaned)
3. **Update the heuristics** in DOCUMENTATION-GUIDE.md to capture similar content in future

## Contents

| File | Source | Notes |
|------|--------|-------|
| contrail-go-stack-lines-450-600.md | contrail-go-stack.md:450-600 | Appears to be Cobra command scaffolding |

## Suggested Heuristic Updates

Based on content in this blackhole, consider adding these patterns:

- [ ] Pattern for Cobra command scaffolding → Implementation layer, appendix
- [ ] Pattern for complete Go file contents → Implementation layer, appendix
```

---

## Content Thresholds

These thresholds control how content is classified during migration and ongoing maintenance.

### Default Thresholds

| Threshold | Default | Purpose |
|-----------|---------|---------|
| `CODE_BLOCK_LINES` | 50 | Code blocks ≥ this many lines go to appendix |
| `STEP_LIST_ITEMS` | 10 | Step-by-step lists ≥ this many items go to appendix |
| `TABLE_ROWS` | 20 | Tables ≥ this many rows go to appendix |
| `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Complete file examples always go to appendix |
| `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Error message catalogs always go to appendix |
| `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Complete shell scripts always go to appendix |

### Pattern-Based Classification

In addition to size thresholds, these patterns trigger appendix placement:

| Pattern | Detection | Destination |
|---------|-----------|-------------|
| Complete file (has shebang or package declaration) | `#!/` or `package main` at start | Appendix |
| Scaffold code (multiple related code blocks) | Sequential code blocks with file paths in headings | Appendix |
| Error catalog (list of error codes/messages) | Headers like "Error Messages", "Exit Codes" with tables | Appendix |
| Workflow example (multi-step with context) | Numbered steps with explanatory text between commands | Appendix |

### Customizing Thresholds

After installation, thresholds are written to `DOCUMENTATION-GUIDE.md`. Edit the thresholds section to customize:

```markdown
## Content Thresholds (Configurable)

| Threshold | Value | Purpose |
|-----------|-------|---------|
| `CODE_BLOCK_LINES` | 30 | ← Lowered from default 50 for this project |
| `STEP_LIST_ITEMS` | 10 | |
| ...
```

The `audit.md` workflow reads these thresholds when analyzing content.

---

## Document Hierarchy and Authority

When documents conflict, higher-authority documents win. Update lower documents to match.

```
┌─────────────────────────────────────────────────────────────┐
│                      MOST AUTHORITATIVE                      │
├─────────────────────────────────────────────────────────────┤
│  ADRs (Architectural Decision Records)                      │
│  - Decisions are immutable once accepted                    │
│  - New decisions supersede (not modify) old ones            │
│  - If anything conflicts with ADR, ADR wins                 │
├─────────────────────────────────────────────────────────────┤
│  Gherkin Feature Files                                       │
│  - Executable specifications                                 │
│  - Tests enforce behavior                                    │
│  - If test passes, documentation is accurate                 │
├─────────────────────────────────────────────────────────────┤
│  Vision (PRD)                                                │
│  - High-level "what" and "why"                              │
│  - Rarely changes after design phase                         │
│  - References ADRs for decisions                             │
├─────────────────────────────────────────────────────────────┤
│  Technical Specification                                     │
│  - Architecture and schemas                                  │
│  - References ADRs for rationale                             │
│  - Updated as architecture evolves                           │
├─────────────────────────────────────────────────────────────┤
│  Reference Documentation (CLI, Config)                       │
│  - Factual, complete, lookup-oriented                        │
│  - Generated where possible                                  │
│  - Updated with implementation                               │
├─────────────────────────────────────────────────────────────┤
│  Implementation Guides (Tech Stack)                          │
│  - How to build, patterns to follow                          │
│  - Most volatile, changes with tech decisions                │
│                     LEAST AUTHORITATIVE                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Layer 1: Decisions (ADRs)

### Purpose

Capture the *why* behind significant technical and product decisions. ADRs become part of the project's institutional memory.

### Characteristics

- **Immutable**: Once accepted, ADRs are never edited (superseded by new ADRs instead)
- **Single-focused**: Each ADR addresses exactly one decision
- **Context-rich**: Explains the situation that necessitated the decision
- **Consequence-aware**: Documents both benefits and trade-offs

### When to Create an ADR

Create an ADR when:
- Choosing between multiple viable technical approaches
- Making a decision that would be expensive to reverse
- Establishing a pattern that will be followed throughout the codebase
- Deviating from common practice or industry norms
- Making a decision that team members might question later

Do NOT create an ADR for:
- Implementation details that don't affect architecture
- Decisions that are easily reversible
- Standard industry practices that need no justification

### Directory Structure

ADRs are simple single files. Unlike other layers, ADRs do NOT use the directory-with-README pattern because:
- ADRs are typically 50-150 lines—compact enough for a single file
- Code examples in ADRs should be included inline (they provide essential context for the decision)
- ADRs rarely need appendices; if an ADR is complex enough to need one, consider splitting it into multiple ADRs

```
docs/decisions/
├── README.md                                     # Index of all ADRs
├── 0001-docker-compose-orchestration.md          # Single file per ADR
├── 0002-two-layer-networking.md
└── 0003-pure-overlay-design.md
```

**Include code examples inline in ADRs.** Code examples in an ADR are part of the decision's context—they show what was decided, not just describe it. Do not move code examples to appendices.

### Lifecycle

```
Draft → Proposed → Accepted → [Superseded by NNNN]
                           → [Deprecated]
```

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | MADR Minimal | Most decisions — structured but lightweight |
| 2 | Y-Statement | Quick capture, smaller decisions, inline documentation |
| 3 | MADR Full | Complex cross-cutting decisions requiring detailed analysis |

See `templates/adr-*.md` for template files.

### Classification Heuristics

| Signal | → ADR Layer |
|--------|-------------|
| "We chose X over Y because..." | ✓ |
| "This pattern will be used throughout..." | ✓ |
| "We considered alternatives including..." | ✓ |
| "Future developers should know that..." | ✓ |
| "Here's how to configure X..." | ✗ (Reference) |
| "The system does X when Y happens..." | ✗ (Specification) |

---

## Layer 2: Vision (PRD-Lite)

### Purpose

Define the product's purpose, goals, and constraints. This is the stable foundation that rarely changes.

### Characteristics

- **Stable**: Changes infrequently (major pivots only)
- **Strategic**: Focuses on *what* and *why*, not *how*
- **Stakeholder-oriented**: Written for business and product audiences
- **Constraint-defining**: Establishes non-goals and boundaries

### What Belongs Here

- Problem statement and market context
- Product vision and success criteria
- Core concepts and glossary
- Non-goals (explicit boundaries)
- User personas and key use cases (high-level)

### What Does NOT Belong Here

- Technical implementation details
- Configuration schemas or API specs
- CLI syntax or command references
- Architecture decisions (those go in ADRs)

### Directory Structure

```
docs/product/
├── README.md              # Auto-generated index
├── vision.md              # Problem, solution, success criteria
├── comparison.md          # How this compares to alternatives (optional)
├── roadmap.md             # Future considerations and planned enhancements (optional)
├── concepts.md            # Core concepts and glossary
├── personas.md            # User personas and use cases (optional)
└── appendices/
    └── concepts/
        └── glossary-full.md  # Extended glossary if needed
```

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | Lean PRD | Most projects — vision + problem + concepts + non-goals |
| 2 | Epic-Based PRD | Agile teams with formal backlog, user story focus |

**Additional Product Documents** (auto-detected during migration):

| Document | Template | Use When |
|----------|----------|----------|
| `comparison.md` | `product-comparison.md` | Comparing with alternative tools |
| `roadmap.md` | `product-roadmap.md` | Future considerations, planned enhancements |

See `templates/prd-*.md` and `templates/product-*.md` for template files.

### Classification Heuristics

| Signal | → Vision Layer |
|--------|----------------|
| "The problem we're solving is..." | ✓ |
| "Success looks like..." | ✓ |
| "This product is NOT for..." | ✓ |
| "A [concept] is defined as..." | ✓ |
| "The architecture uses..." | ✗ (Architecture) |
| "When the user runs command X..." | ✗ (Specification/Reference) |

**Comparison Document Heuristics**:

| Signal | → Comparison |
|--------|--------------|
| Comparison tables with other tools/products | ✓ |
| "vs" language ("X vs Y", "compared to") | ✓ |
| Feature matrices comparing multiple products | ✓ |
| "Related Tools" or "Similar Projects" sections | ✓ |
| "Why choose X over Y" discussions | ✓ |
| Tool names in comparative context ("Unlike Docker Compose...") | ✓ |

**Roadmap Document Heuristics**:

| Signal | → Roadmap |
|--------|-----------|
| "Future Considerations" sections | ✓ |
| "Roadmap" or "Future Work" headings | ✓ |
| "Planned features" or "upcoming" language | ✓ |
| "Phase 2", "Phase 3", "v2.0" planning content | ✓ |
| "Not yet implemented" with intent to implement | ✓ |
| "Eventually" or "in the future" feature descriptions | ✓ |
| Enhancement proposals or feature wishlists | ✓ |

---

## Layer 3: Architecture

### Purpose

Show how the system's components relate to each other and to external systems. Provides context for understanding the codebase.

### Characteristics

- **Evolving**: Updates as the design matures
- **Visual**: Diagrams are primary, text supports
- **Hierarchical**: Multiple levels of abstraction
- **Decision-linked**: References ADRs for rationale

### What Belongs Here

- System context diagrams (C4 Level 1)
- Container diagrams (C4 Level 2)
- Component diagrams (C4 Level 3) for complex subsystems
- Cross-cutting concerns (security, logging, error handling patterns)
- Quality attributes and constraints
- Known risks and technical debt

### What Does NOT Belong Here

- API endpoint details (Reference layer)
- Step-by-step behavioral flows (Specification layer)
- Code-level class diagrams (rarely useful)

### Directory Structure

```
docs/architecture/
├── README.md              # Auto-generated index
├── overview.md            # System context, key containers
├── networking.md          # Network topology, communication patterns
├── data-flow.md           # How data moves through the system
├── cross-cutting.md       # Security, logging, error handling
└── appendices/
    ├── overview/
    │   └── c4-diagrams.md        # Full diagram source/exports
    └── cross-cutting/
        └── security-details.md   # Detailed security analysis if needed
```

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | C4-Lite | Most projects — Context + Container diagrams with narrative |
| 2 | arc42 Full | Large systems, formal documentation requirements, regulatory needs |
| 3 | Structurizr DSL | Model-driven diagrams generated from code (requires tooling) |

**Note**: C4 and arc42 are complementary. Use C4 diagrams within arc42 sections for comprehensive coverage:
- C4 Context Diagram → arc42 Section 3 (Context and Scope)
- C4 Container Diagram → arc42 Section 5 (Building Block View)

See `templates/architecture-*.md` for template files.

### Classification Heuristics

| Signal | → Architecture Layer |
|--------|---------------------|
| "The system is composed of..." | ✓ |
| "Component A communicates with B via..." | ✓ |
| "Security is handled by..." (pattern) | ✓ |
| "This diagram shows..." | ✓ |
| "Run this command to..." | ✗ (Reference) |
| "When X happens, the system does Y..." | ✗ (Specification) |

---

## Layer 4: Specifications

### Purpose

Detail how specific features work. These are living documents that evolve with the implementation.

### Characteristics

- **Living**: Updated as implementation proceeds
- **Detailed**: Includes schemas, state machines, algorithms
- **Testable**: Specific enough to verify
- **ADR-linked**: References decisions for rationale

### What Belongs Here

- Feature specifications with detailed behavior
- State machines and lifecycle diagrams
- Data schemas (YAML, JSON examples)
- Error handling and edge cases
- Integration points between components

### What Does NOT Belong Here

- API reference tables (Reference layer)
- Why decisions were made (ADR layer)
- User-facing tutorials (separate documentation)

### Directory Structure

```
docs/specs/
├── README.md                   # Auto-generated index
├── workspace-lifecycle.md      # Workspace states and transitions
├── port-assignment.md          # Port allocation algorithm
├── overlay-generation.md       # How overlays are generated
├── proxy-integration.md        # Traefik integration details
└── appendices/
    ├── workspace-lifecycle/
    │   └── state-diagrams.md   # Detailed state machine diagrams
    ├── port-assignment/
    │   └── algorithm-details.md  # Full algorithm pseudocode
    └── overlay-generation/
        └── complete-examples.md  # Full override file examples
```

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | Feature Spec | Most features — narrative + examples + edge cases |
| 2 | RFC-Style | Proposed changes needing team review before implementation |

See `templates/spec-*.md` for template files.

### Classification Heuristics

| Signal | → Specification Layer |
|--------|----------------------|
| "When X happens, the system does Y..." | ✓ |
| "The valid states are..." | ✓ |
| "Here's how feature X works..." | ✓ |
| "Edge case: if Y, then Z..." | ✓ |
| "Why we chose X..." | ✗ (ADR) |
| "Command syntax: ..." | ✗ (Reference) |

---

## Layer 5: Reference

### Purpose

Provide lookup information for users and developers. Quick answers to "what are the options?" questions.

### Characteristics

- **Exhaustive**: Covers all options/commands/configs
- **Scannable**: Tables, lists, clear headings
- **Accurate**: Often generated from code
- **Stable format**: Structure rarely changes

### What Belongs Here

- CLI command reference
- Configuration options and defaults
- API endpoint tables
- Environment variables
- Error codes and messages
- Docker labels and metadata

### What Does NOT Belong Here

- Why options exist (ADR layer)
- How features work conceptually (Specification layer)
- Step-by-step guides (Tutorials, if applicable)

### Directory Structure

```
docs/reference/
├── README.md              # Auto-generated index
├── cli.md                 # Command reference (commands, flags, basics)
├── configuration.md       # Config file options (schemas, defaults)
├── environment.md         # Environment variables
├── labels.md              # Docker labels reference
├── errors.md              # Error codes and meanings (summary)
└── appendices/
    ├── cli/
    │   ├── detailed-examples.md     # Complete workflow examples
    │   ├── error-messages.md        # Full error message catalog
    │   └── output-formats.md        # Detailed output examples
    ├── configuration/
    │   └── complete-schemas.md      # Full YAML/JSON schema examples
    └── errors/
        └── full-catalog.md          # Complete error catalog
```

### Generation Strategy

Reference docs are ideal candidates for generation:

| Source | Generated Doc |
|--------|---------------|
| Cobra CLI definitions | CLI reference |
| Go struct tags | Configuration schema |
| Error constants | Error code reference |
| Docker Compose examples | Label reference |

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1** | CLI Reference | Projects with CLI tools |
| **1** | Configuration Reference | Projects with config files |

Include whichever templates apply to your project.

See `templates/reference-*.md` for template files.

### Classification Heuristics

| Signal | → Reference Layer |
|--------|------------------|
| "The available options are..." | ✓ |
| "The default value is..." | ✓ |
| "Syntax: command [options]..." | ✓ |
| "Error CODE means..." | ✓ |
| "This works because..." | ✗ (Specification/ADR) |
| "The architecture consists of..." | ✗ (Architecture) |

---

## Layer 6: Behaviors (Executable Specifications)

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

```
features/
├── workspace/
│   ├── workspace-lifecycle.feature
│   └── workspace-isolation.feature
├── proxy/
│   └── proxy-routing.feature
└── support/
    └── step_definitions/
```

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | Gherkin | Standard BDD format, wide tooling support |
| 2 | Doctest-Style | Embedded in documentation, simple examples (language-specific) |

See `templates/behavior-gherkin.feature` for template file.

### Classification Heuristics

| Signal | → Behavior Layer |
|--------|-----------------|
| "Given X, when Y, then Z..." | ✓ |
| "This behavior must not regress..." | ✓ |
| "Example: user does X, sees Y..." | ✓ |
| "The algorithm works by..." | ✗ (Specification) |
| "Command options include..." | ✗ (Reference) |

---

## Layer 7: Implementation Guides

### Purpose

Describe *how to build* the system—technology stack, scaffolding, code templates, and implementation priorities. These are consumed during implementation and then archived or absorbed into code.

### Characteristics

- **Short-lived**: Consumed during implementation, then archived
- **Code-adjacent**: Contains executable code, templates, scaffolding
- **Version-specific**: Includes dependency versions, project structure
- **Actionable**: Step-by-step instructions for building

### What Belongs Here

- Technology stack with specific versions
- Project structure and directory layout
- Code scaffolding and templates
- Dependency rationale (why this library over that one)
- Implementation priority phases
- Build and development setup instructions

### What Does NOT Belong Here

- Behavioral specifications (Layer 4)
- Architecture decisions rationale (Layer 1: ADR)
- API reference (Layer 5)

### Directory Structure

```
docs/implementation/
├── README.md                  # Auto-generated index
├── tech-stack.md              # Dependencies, versions, rationale (overview)
├── scaffolding.md             # Project structure explanation
├── build-setup.md             # Development environment setup
└── appendices/
    ├── tech-stack/
    │   ├── go-scaffolding/              # Complete code scaffolding
    │   │   ├── cmd-main.go.md
    │   │   ├── cli-root.go.md
    │   │   ├── cli-workspace.go.md
    │   │   └── cli-commands.go.md
    │   └── shell-scripts/               # Complete shell scripts
    │       ├── bash.sh.md
    │       ├── zsh.zsh.md
    │       └── fish.fish.md
    └── scaffolding/
        └── directory-tree.md  # Full project structure with explanations
```

### Lifecycle

```
Draft → Active (during implementation) → Archived/Absorbed
```

Implementation guides are unique in that they have a planned end-of-life:
- Once implementation is complete, the content either:
  - Gets **archived** (moved to an `archive/` folder or deleted)
  - Gets **absorbed** into the codebase (as README files in packages, code comments, etc.)

### Template Options

| Priority | Template | Use When |
|----------|----------|----------|
| **1 (Default)** | Tech Stack | All projects — dependencies, versions, rationale |

See `templates/implementation-tech-stack.md` for template file.

### Classification Heuristics

| Signal | → Implementation Layer |
|--------|----------------------|
| "Install these dependencies..." | ✓ |
| "The project structure is..." | ✓ |
| "Use this code template..." | ✓ |
| "We chose library X because..." | ✓ (or ADR if significant) |
| "Implementation phases are..." | ✓ |
| "When X happens, the system does Y..." | ✗ (Specification) |
| "The architecture uses..." | ✗ (Architecture) |

---

## Classification Decision Tree

Use this flowchart to classify content:

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
└─ NO → Reconsider: may not need documentation
```

---

## Cross-Layer Linking

### Linking Conventions

Use relative Markdown links to connect layers:

```markdown
<!-- In a specification (ADRs are simple files, not directories) -->
This design follows the overlay pattern. See [ADR-0003](../decisions/0003-pure-overlay-design.md) for rationale.

<!-- In architecture -->
For command details, see the [CLI Reference](../reference/cli.md).

<!-- In a behavior -->
This scenario verifies the workspace lifecycle described in [Workspace Lifecycle Spec](../specs/workspace-lifecycle.md).

<!-- From main doc to appendix -->
For complete workflow examples, see [Detailed Examples](./appendices/cli/detailed-examples.md).
```

### Link Direction Guidelines

| From | To | When |
|------|----|------|
| Specification | ADR | Explaining *why* a design choice was made |
| Architecture | ADR | Justifying architectural patterns |
| Architecture | Specification | Deep-diving into component behavior |
| Behavior | Specification | Referencing the spec being verified |
| Reference | Specification | Providing conceptual context |
| Implementation | Specification | Referencing the spec being implemented |
| Implementation | ADR | Explaining technology choice rationale |

---

## Single Source of Truth (SSOT)

### Canonical Source Table

Every fact should be mastered in exactly one place:

| Information Type | Canonical Source | Derived/Referenced From |
|------------------|------------------|------------------------|
| Why we chose X over Y | ADR | PRD, Tech Spec, Architecture |
| Configuration schema | Tech Spec (or code) | Reference docs, examples |
| Command syntax | CLI Reference | Tech Spec operations |
| Feature behavior | Specification | Gherkin tests |
| API contracts | OpenAPI/Protobuf | Reference docs |
| Port assignment algorithm | Specification | ADR (for rationale) |

### Key Principle

When you find yourself copying information, create a reference instead:

```markdown
For details on the port assignment algorithm, see
[Port Assignment Specification](./specs/port-assignment.md).
```

When sources conflict, the canonical source wins. Update derived views to match.

---

## Preventing AI-Generated Conflicts

When AI tools generate feature specs, apply this workflow:

### Before Creating New Documentation

1. **Check ADRs first** — Does the generated spec contradict any accepted decision?
2. **Check Gherkin scenarios** — Would implementing this break existing tests?
3. **Check specifications** — Does this conflict with existing feature specs?

### When Conflicts Are Found

1. **If new spec contradicts ADR** — The ADR wins. Revise the new spec.
2. **If new spec contradicts old spec** — Resolve the conflict explicitly:
   - Update the old spec if requirements changed
   - Create a new ADR if a decision changed
   - Revise the new spec if it was generated incorrectly
3. **Identify new decisions** — If the generated spec implies new decisions, create ADRs

### Validation Checklist

- [ ] No contradiction with existing ADRs
- [ ] No contradiction with existing specifications
- [ ] Cross-links to related documents added
- [ ] Appropriate layer (use classification decision tree)
- [ ] Follows layer template

---

## Maintenance Workflow

### When Documents Update

| Trigger | Layers Affected |
|---------|-----------------|
| New decision made | Create new ADR (Layer 1) |
| Product pivot | Update Vision (Layer 2) |
| Design changes | Update Architecture (Layer 3), possibly new ADR |
| Feature implementation | Update Specification (Layer 4) |
| CLI/API changes | Update Reference (Layer 5) |
| Bug fix for critical path | Update Behavior (Layer 6) |
| Starting implementation | Create/update Implementation (Layer 7) |
| Implementation complete | Archive or absorb Implementation (Layer 7) |

### Version Management

- **ADRs**: Immutable. Supersede with new ADR, update status to "Superseded by NNNN"
- **Vision**: Version bumps for significant changes (1.0 → 2.0)
- **Architecture/Specs**: Semantic versioning (MAJOR.MINOR.PATCH)
- **Reference**: Version matches software version or is unversioned
- **Behaviors**: Tied to test suite, no separate versioning
- **Implementation**: Version matches implementation phase; archived when complete

### Review Checklist

When reviewing documentation changes:

- [ ] Content is in the correct layer
- [ ] Links to other layers are accurate
- [ ] No duplicate content across layers
- [ ] Templates are followed
- [ ] Decision rationale is in ADRs, not inline

---

## Tooling Recommendations

### Tier 1: Essential (Start Here)

| Tool | Purpose |
|------|---------|
| **markdownlint** | Format consistency |
| **Git + PRs** | Version control, review |
| **MkDocs or similar** | Publishing (optional) |

### Tier 2: Recommended (Add When Needed)

| Tool | Purpose |
|------|---------|
| **Vale** | Terminology consistency, style guide |
| **Log4brains** | ADR static site generation |
| **Structurizr** | C4 diagram generation from DSL |

### Tier 3: Advanced (For Larger Projects)

| Tool | Purpose |
|------|---------|
| **Cucumber/Gherkin** | Executable specifications |
| **Link checker** | Validate cross-references |
| **Custom generators** | Reference docs from code |

---

## Quick Start Checklist

### Initial Setup

- [ ] Create directory structure (see below)
- [ ] Choose templates for each layer
- [ ] Document initial decisions as ADRs
- [ ] Write initial Vision document
- [ ] Set up basic linting

### Recommended Directory Structure

```
project/
├── docs/
│   ├── DOCUMENTATION-GUIDE.md   # Project-specific guide (generated)
│   │
│   ├── decisions/               # Layer 1: ADRs (simple files, not directories)
│   │   ├── README.md            # Index of all ADRs
│   │   └── 0001-initial-architecture.md   # Single file per ADR
│   │
│   ├── product/                 # Layer 2: Vision
│   │   ├── README.md            # Auto-generated index
│   │   ├── vision.md
│   │   └── appendices/          # Only if needed
│   │
│   ├── architecture/            # Layer 3: Architecture
│   │   ├── README.md            # Auto-generated index
│   │   ├── overview.md
│   │   └── appendices/          # Only if needed
│   │
│   ├── specs/                   # Layer 4: Specifications
│   │   ├── README.md            # Auto-generated index
│   │   ├── [feature].md
│   │   └── appendices/
│   │       └── [feature]/       # Appendices for [feature].md
│   │
│   ├── reference/               # Layer 5: Reference
│   │   ├── README.md            # Auto-generated index
│   │   ├── cli.md
│   │   └── appendices/
│   │       └── cli/
│   │           └── detailed-examples.md
│   │
│   ├── implementation/          # Layer 7: Implementation Guides
│   │   ├── README.md            # Auto-generated index
│   │   ├── tech-stack.md
│   │   └── appendices/
│   │       └── tech-stack/
│   │           └── scaffolding/
│   │
│   ├── migration/               # Content awaiting final placement
│   │   └── README.md            # (only created if content exists)
│   │
│   └── blackhole/               # Unclassified content
│       └── README.md            # (only created if content exists)
│
├── features/                    # Layer 6: Behaviors
│   └── [feature].feature
└── ...
```

---

## Operational Guides

This system includes operational guides for common workflows:

| Guide | Purpose |
|-------|---------|
| `install.md` | Set up documentation system (fresh or migration) |
| `create.md` | Create new documentation with validation |
| `update.md` | Update documentation after code changes |
| `sync.md` | Audit and synchronize docs with code |
| `refine.md` | Improve documentation quality |
| `audit.md` | Compare docs vs source, analyze migration loss |

See each guide for detailed instructions.

### Migration and Audit Workflow

After a migration install, the system automatically:

1. **Generates documentation** into the layered structure
2. **Compares original vs generated** content
3. **Produces a summary** showing:
   - Content successfully migrated to layers
   - Content in `migration/` (classified but placement uncertain)
   - Content in `blackhole/` (could not be classified)
4. **Suggests heuristic updates** for blackhole content

The `audit.md` guide can be run manually at any time to repeat this analysis.

---

## Files in This System

```
layered-docs-system/
├── LAYERED-DOCUMENTATION-SYSTEM.md    # This file (full reference)
├── install.md                         # Installation workflow
├── create.md                          # Creating new documentation
├── update.md                          # Updating after changes
├── sync.md                            # Synchronization audit
├── refine.md                          # Quality improvement
├── audit.md                           # Compare docs vs source, migration analysis
└── templates/
    ├── adr-madr-minimal.md            # ADR template (recommended)
    ├── adr-y-statement.md             # ADR template (lightweight)
    ├── adr-madr-full.md               # ADR template (comprehensive)
    ├── prd-lean.md                    # Vision template (recommended)
    ├── prd-epic-based.md              # Vision template (agile)
    ├── architecture-c4-lite.md        # Architecture template (recommended)
    ├── architecture-arc42.md          # Architecture template (comprehensive)
    ├── spec-feature.md                # Specification template (recommended)
    ├── spec-rfc.md                    # Specification template (proposals)
    ├── reference-cli.md               # Reference template (CLI)
    ├── reference-config.md            # Reference template (configuration)
    ├── behavior-gherkin.feature       # Behavior template (Gherkin)
    ├── implementation-tech-stack.md   # Implementation template
    ├── appendix.md                    # Appendix template
    ├── migration-readme.md            # README template for migration/
    └── blackhole-readme.md            # README template for blackhole/
```
