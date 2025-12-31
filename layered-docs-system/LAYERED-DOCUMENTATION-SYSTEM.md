# Layered Documentation System

**Version**: 1.0.0
**Purpose**: A cohesive, maintainable system for creating and managing software documentation using a layered approach.

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

**Interactive** (prompts for all options):
```
Execute @path/to/layered-docs-system/install.md
```

The installer will:
1. Determine install location (or ask)
2. Detect existing documentation to migrate (or ask)
3. If migrating: analyze content and map to layers
4. Ask which layers and templates to include
5. Ask which tooling tier to set up
6. Create directory structure and migrate/copy content
7. Generate a project-specific documentation guide
8. Add cross-layer links (for migrations)

### For Humans

1. Read this document to understand the system
2. Use the [Quick Start Checklist](#quick-start-checklist) to set up manually
3. Or use the `install.md` with an AI assistant for guided setup

---

## Overview

This system organizes documentation into six distinct layers, each with a specific purpose, ownership model, and lifecycle. The goal is to prevent drift between documents by ensuring each piece of information has exactly one canonical source.

### The Six Layers

| Layer | Purpose | Stability | Audience |
|-------|---------|-----------|----------|
| **1. Decisions** | Capture *why* choices were made | Immutable once accepted | Future maintainers, architects |
| **2. Vision** | Define *what* we're building and *why* | Stable, rarely changes | All stakeholders |
| **3. Architecture** | Show *how* components relate | Evolves with design | Engineers, architects |
| **4. Specifications** | Detail *how* features work | Living, evolves with implementation | Engineers |
| **5. Reference** | Provide lookup information | Generated or hand-maintained | Users, developers |
| **6. Behaviors** | Verify system meets expectations | Executable, tied to tests | QA, engineers |

### Core Principles

1. **Single Source of Truth**: Each piece of information lives in exactly one place
2. **Separation of Concerns**: Different document types serve different purposes
3. **Appropriate Stability**: Some documents are immutable; others evolve constantly
4. **Clear Ownership**: Each layer has defined maintainers and update triggers
5. **Linkage Over Duplication**: Reference other documents rather than copying content

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

```
docs/decisions/
├── 0001-use-docker-compose-for-orchestration.md
├── 0002-two-layer-networking-model.md
├── 0003-pure-overlay-design-pattern.md
└── index.md                                      # Optional: links to all ADRs
```

### Lifecycle

```
Draft → Proposed → Accepted → [Superseded by NNNN]
                           → [Deprecated]
```

### Template Options

**Priority 1: MADR Minimal** (Recommended)
- Structured enough to be consistent
- Light enough to encourage adoption
- See: `templates/adr-madr-minimal.md`

**Priority 2: Y-Statement**
- Single-sentence format for quick capture
- Best for smaller decisions
- See: `templates/adr-y-statement.md`

**Priority 3: MADR Full**
- Comprehensive template with all sections
- Use for complex, cross-cutting decisions
- See: `templates/adr-madr-full.md`

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
├── vision.md              # Problem, solution, success criteria
├── concepts.md            # Core concepts and glossary
└── personas.md            # User personas and use cases (optional)
```

### Template Options

**Priority 1: Lean PRD**
- Vision + problem + concepts + non-goals
- Minimal structure, maximum clarity
- See: `templates/prd-lean.md`

**Priority 2: Epic-Based PRD**
- Organized around user stories and epics
- Better for agile teams
- See: `templates/prd-epic-based.md`

### Classification Heuristics

| Signal | → Vision Layer |
|--------|----------------|
| "The problem we're solving is..." | ✓ |
| "Success looks like..." | ✓ |
| "This product is NOT for..." | ✓ |
| "A [concept] is defined as..." | ✓ |
| "The architecture uses..." | ✗ (Architecture) |
| "When the user runs command X..." | ✗ (Specification/Reference) |

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
├── overview.md            # System context, key containers
├── networking.md          # Network topology, communication patterns
├── data-flow.md           # How data moves through the system
├── cross-cutting.md       # Security, logging, error handling
└── diagrams/              # Source files for diagrams (optional)
    └── workspace.structurizr.dsl
```

### Template Options

**Priority 1: C4-Lite**
- Context + Container diagrams with narrative
- Lightweight, easy to maintain
- See: `templates/architecture-c4-lite.md`

**Priority 2: arc42-Simplified**
- Structured sections from arc42
- More comprehensive coverage
- See: `templates/architecture-arc42.md`

**Priority 3: Structurizr DSL**
- Model-driven, diagrams generated from code
- Best for larger systems
- Requires tooling setup

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
- User-facing tutorials (Tutorial layer, if applicable)

### Directory Structure

```
docs/specs/
├── workspace-lifecycle.md      # Workspace states and transitions
├── port-assignment.md          # Port allocation algorithm
├── overlay-generation.md       # How overlays are generated
└── proxy-integration.md        # Traefik integration details
```

### Template Options

**Priority 1: Feature Spec**
- Narrative + examples + edge cases
- Good for most features
- See: `templates/spec-feature.md`

**Priority 2: RFC-Style**
- Problem + proposal + alternatives
- Good for proposed changes
- See: `templates/spec-rfc.md`

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
- Step-by-step guides (Tutorial layer)

### Directory Structure

```
docs/reference/
├── cli.md                 # Command reference
├── configuration.md       # Config file options
├── environment.md         # Environment variables
├── labels.md              # Docker labels reference
└── errors.md              # Error codes and meanings
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

**Priority 1: Command Reference**
- Per-command sections with options tables
- See: `templates/reference-cli.md`

**Priority 2: Configuration Reference**
- Hierarchical config with defaults
- See: `templates/reference-config.md`

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

Define expected behaviors in a way that can be automatically verified. Tests that double as documentation.

### Characteristics

- **Executable**: Run as part of CI/CD
- **Living**: Fail when behavior changes
- **Example-driven**: Concrete scenarios, not abstract
- **User-focused**: Written from user perspective

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

**Priority 1: Gherkin**
- Standard BDD format
- Wide tooling support (Cucumber, Behave, etc.)
- See: `templates/behavior-gherkin.feature`

**Priority 2: Doctest-Style**
- Embedded in documentation
- Good for simple examples
- Language-specific tooling

### Classification Heuristics

| Signal | → Behavior Layer |
|--------|-----------------|
| "Given X, when Y, then Z..." | ✓ |
| "This behavior must not regress..." | ✓ |
| "Example: user does X, sees Y..." | ✓ |
| "The algorithm works by..." | ✗ (Specification) |
| "Command options include..." | ✗ (Reference) |

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
└─ NO → Reconsider: may not need documentation
```

---

## Cross-Layer Linking

### Linking Conventions

Use relative Markdown links to connect layers:

```markdown
<!-- In a specification -->
This design follows the overlay pattern. See [ADR-0003](../decisions/0003-pure-overlay-design.md) for rationale.

<!-- In architecture -->
For command details, see the [CLI Reference](../reference/cli.md).

<!-- In a behavior -->
This scenario verifies the workspace lifecycle described in [Workspace Lifecycle Spec](../specs/workspace-lifecycle.md).
```

### Link Direction Guidelines

| From | To | When |
|------|----|------|
| Specification | ADR | Explaining *why* a design choice was made |
| Architecture | ADR | Justifying architectural patterns |
| Architecture | Specification | Deep-diving into component behavior |
| Behavior | Specification | Referencing the spec being verified |
| Reference | Specification | Providing conceptual context |

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

### Version Management

- **ADRs**: Immutable. Supersede with new ADR, update status to "Superseded by NNNN"
- **Vision**: Version bumps for significant changes (1.0 → 2.0)
- **Architecture/Specs**: Semantic versioning (MAJOR.MINOR.PATCH)
- **Reference**: Version matches software version or is unversioned
- **Behaviors**: Tied to test suite, no separate versioning

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
| **Semcheck** | AI-powered spec-implementation drift detection |
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
│   ├── decisions/           # Layer 1: ADRs
│   │   └── 0001-initial-architecture.md
│   ├── product/             # Layer 2: Vision
│   │   └── vision.md
│   ├── architecture/        # Layer 3: Architecture
│   │   └── overview.md
│   ├── specs/               # Layer 4: Specifications
│   │   └── [feature].md
│   └── reference/           # Layer 5: Reference
│       └── cli.md
├── features/                # Layer 6: Behaviors (optional)
│   └── [feature].feature
└── ...
```

---

## Appendix: Canonical Sources

Define which source is authoritative for each type of information:

| Information Type | Canonical Source | Derived Views |
|------------------|------------------|---------------|
| CLI flags and options | Code (Cobra/argparse definitions) | Reference docs |
| Configuration schema | Code (structs, types) | Reference docs, examples |
| Decision rationale | ADRs | Linked from specs |
| Feature behavior | Specifications | Behavior tests |
| API contracts | OpenAPI/Protobuf | Reference docs |

When sources conflict, the canonical source wins. Update derived views to match.

---

## Files in This System

```
layered-docs-system/
├── LAYERED-DOCUMENTATION-SYSTEM.md    # This file (full reference)
├── install.md                         # Interactive installer for AI agents
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
    └── behavior-gherkin.feature       # Behavior template (Gherkin)
```
