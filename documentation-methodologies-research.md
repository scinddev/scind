# Comprehensive Research: Documentation Methodologies for Evolving Software Specifications

**Date**: December 2024
**Context**: Research into documentation standards, formats, and systems for managing comprehensive and evolving software specifications

---

## The Problem Being Addressed

Looking at the Contrail specifications, there are currently:
1. **PRD** (contrail-prd.md) - Product vision, concepts, decisions
2. **Technical Spec** (contrail-technical-spec.md) - Architecture, schemas, operations
3. **CLI Reference** (contrail-cli-reference.md) - Command documentation
4. **Go Stack** (contrail-go-stack.md) - Implementation details
5. **Shell Integration** (contrail-shell-integration.md) - Shell scripts

The issue: these documents contain **overlapping concerns** and embedded decisions that can drift out of sync. The PRD already contains inline "Decision" blocks that are essentially ADRs mixed with requirements.

When working with agentic specification tools like Agent OS, many specs are created for features that are well-written individually but often conflict or contradict other design decisions in other specifications or with the reality of the project state.

---

## Part 1: ADRs and MADR - What They Solve (and Don't)

### Architectural Decision Records (ADR)

**Source**: [ADR Homepage](https://adr.github.io/)

ADRs are **decision-focused**, not comprehensive specification containers. Per [AWS's ADR guidance](https://aws.amazon.com/blogs/architecture/master-architecture-decision-records-adrs-best-practices-for-effective-decision-making/), ADRs should:
- Focus on a **single decision**
- Be **brief** (like meeting minutes, not functional specifications)
- Become **immutable** once accepted (new decisions supersede old ones)

### MADR (Markdown ADR)

**Source**: [MADR GitHub](https://github.com/adr/madr)

MADR provides four template variants:
- **Full** (`adr-template.md`) - All sections with explanations
- **Minimal** (`adr-template-minimal.md`) - Mandatory sections only
- **Bare** (`adr-template-bare.md`) - All sections, empty
- **Bare-Minimal** (`adr-template-bare-minimal.md`) - Just mandatory sections, blank

### The Y-Statement Format

**Source**: [Y-Statement Template](https://medium.com/olzzio/y-statements-10eb07b5a177)

An even more lightweight approach:

> "In the context of **[use case]**, facing **[concern]**, we decided for **[option]** to achieve **[quality]**, accepting **[downside]**."

**Long form**:
> "In the context of **[use case/user story]**, facing **[concern]**, we decided for **[option]** and neglected **[other options]**, to achieve **[system qualities/desired consequences]**, accepting **[downside/undesired consequences]**, because **[additional rationale]**."

**Example**:
> "In the context of the Web shop service, facing the need to keep user session data consistent and current across instances, we decided for the Database Session State Pattern (and against Client Session State or Server Session State) to achieve cloud elasticity, accepting that a session database needs to be designed, implemented, and replicated."

### What ADRs Don't Solve

ADRs capture **why** decisions were made but don't address:
- Product requirements (the "what")
- Technical specifications (the "how" in detail)
- API/CLI references
- Implementation guides
- Behavioral examples

---

## Part 2: Complementary Documentation Frameworks

### Diátaxis Framework

**Source**: [Diátaxis](https://diataxis.fr/)

Diátaxis identifies **four types of documentation** based on user needs:

| Type | Purpose | Contrail Analog |
|------|---------|-----------------|
| **Tutorials** | Learning-oriented, step-by-step | (Not present - future user guides) |
| **How-to Guides** | Task-oriented, solving problems | Shell integration examples |
| **Reference** | Information-oriented, technical descriptions | CLI Reference, schemas |
| **Explanation** | Understanding-oriented, context and rationale | PRD concepts, decision blocks |

**Key insight**: The PRD currently mixes **Explanation** (concepts, rationale) with **Reference** (configuration schemas). Diátaxis would suggest separating these.

### arc42 Template

**Source**: [arc42](https://arc42.org/)

arc42 is a **comprehensive architecture documentation template** with 12 sections:
1. Introduction and Goals
2. Constraints
3. Context and Scope
4. Solution Strategy
5. Building Block View
6. Runtime View
7. Deployment View
8. Cross-cutting Concepts
9. Architecture Decisions
10. Quality Requirements
11. Risks and Technical Debt
12. Glossary

**Key insight**: arc42 **includes ADRs as Section 9** but treats them as one piece of a larger documentation structure. Everything is optional—treat it like compartments in a cabinet.

### C4 Model

**Source**: [C4 Model](https://c4model.com/)

The C4 Model provides **four levels of abstraction** for architecture diagrams:
1. **System Context** - How your system relates to users and other systems
2. **Container** - High-level technology choices (apps, databases, etc.)
3. **Component** - Components within each container
4. **Code** - Class/function level (optional)

**Key insight**: The Contrail specs contain embedded C4-style diagrams (the ASCII network topology diagrams). [Structurizr](https://structurizr.com/) could generate these from a model and keep them consistent.

---

## Part 3: Living Documentation & Self-Healing Specs

### Living Documentation (Cyrille Martraire)

**Source**: [Living Documentation Book (O'Reilly)](https://www.oreilly.com/library/view/living-documentation-continuous/9780134689418/)

Four principles:
1. **Reliable** - Automatically updated, can't drift from code
2. **Low effort** - Generated, not manually maintained
3. **Collaborative** - Everyone contributes
4. **Insightful** - Reveals understanding, not just lists

**Implementation approaches**:
- **Executable specifications** via [Cucumber/Gherkin](https://cucumber.io/docs/) - Behavior specs that double as tests
- **Code-generated docs** - API docs from code annotations
- **Model-driven documentation** - Generate diagrams and text from a single model

### Docs-as-Code

**Source**: [Docs as Code](https://www.gitbook.com/blog/what-is-docs-as-code)

Key practices:
- Store docs in the **same repo** as code
- Use **Markdown/AsciiDoc**
- Apply **version control** workflows (branches, PRs, reviews)
- **Automate** builds and publishing
- Use **linters** (markdownlint, Vale) for consistency

**Tools**: [MkDocs](https://github.com/mkdocs/mkdocs), Hugo, Docusaurus, Jekyll

### AI-Powered Specification Sync

New tools address the problem of specs drifting from reality:

- **[Semcheck](https://www.xugj520.cn/en/archives/ai-code-documentation-sync-tool.html)** - LLM-powered CLI that verifies spec-implementation alignment
- **[GitHub Spec Kit](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/)** - Makes specifications the **center of the engineering process**, driving implementation and checklists

---

## Part 4: PRD Approaches

### Modern Agile PRDs

**Source**: [Atlassian PRD Guide](https://www.atlassian.com/agile/product-management/requirements)

Agile PRDs:
- Focus on **shared understanding**, not exhaustive detail
- Use **themes**, **epics**, and **user stories** rather than monolithic specs
- Are **living documents** that evolve with feedback

### Document Type Hierarchy

**Source**: [PRD vs BRD vs TRD](https://medium.com/@kokoproduct/decoding-product-development-documents-prd-vs-brd-vs-trd-a45ba984cd78)

| Document | Purpose | Audience |
|----------|---------|----------|
| **BRD** (Business Requirements) | Why we're building this | Executives, stakeholders |
| **PRD** (Product Requirements) | What we're building | Product, design, engineering |
| **TRD/TDD** (Technical Design) | How we're building it | Engineering, architecture |
| **ADR** | Why specific technical decisions | Engineering, future maintainers |

---

## Part 5: Single Source of Truth Patterns

### Single Source of Truth (SSOT)

**Source**: [Wikipedia - SSOT](https://en.wikipedia.org/wiki/Single_source_of_truth)

Key principle: **Every data element is mastered in only one place**.

Applied to specifications:
- One canonical location per concern
- Updates happen only in that location
- Other views are derived/generated

### Specification-First Development

**Source**: [Red Hat SSOT Architecture](https://www.redhat.com/en/blog/single-source-truth-architecture)

> "Instead of using the source code as the SSOT, you can use an API specification... both code and documentation can be generated from the specification document."

---

## Part 6: A Proposed Hybrid System for Contrail

Given the specific needs, here's a synthesis:

### Layer 1: Decisions (ADRs using MADR)

**Location**: `docs/decisions/`

Move the embedded "Decision" blocks from the PRD into standalone ADRs:
- `0001-docker-compose-project-name-isolation.md`
- `0002-two-layer-networking.md`
- `0003-pure-overlay-design.md`
- etc.

**Template**: Use MADR minimal, or even simpler Y-statements for quick capture.

**Lifecycle**: Immutable once accepted. New decisions supersede rather than edit.

### Layer 2: Product Requirements (Lightweight PRD)

**Location**: `docs/product/`

Focus on:
- **Vision** and **problem statement** (stable)
- **Core concepts** and **glossary** (evolving)
- **User stories/epics** (versioned, living)
- **Non-goals** (stable)
- **Success criteria** (versioned)

**What moves out**: Technical schemas, implementation details, CLI syntax.

### Layer 3: Architecture (arc42-lite or C4-based)

**Location**: `docs/architecture/`

Could use a simplified arc42 structure:
- Context and scope (diagrams)
- Building blocks view
- Cross-cutting concepts
- Quality requirements
- Risks

**Integration**: Link to ADRs for decision rationale.

### Layer 4: Reference (Diátaxis "Reference" type)

**Location**: `docs/reference/`

This is where CLI Reference, configuration schemas, and environment variables belong:
- **CLI Reference** - Generated or hand-maintained
- **Configuration Schema** - Could be generated from Go struct definitions
- **API/Labels Reference** - Docker labels, environment variables

### Layer 5: Implementation Specs (TDD-style)

**Location**: `docs/specs/` or inline with code

The Go Stack spec is essentially a Technical Design Document. It should:
- Reference ADRs for decision rationale
- Reference the architecture for context
- Be the source of truth for implementation specifics

### Layer 6: Executable Specifications (Optional - Gherkin/BDD)

For critical behaviors, consider [Cucumber/Gherkin](https://cucumber.io/docs/gherkin/) specs:

```gherkin
Feature: Workspace Isolation
  Scenario: Two workspaces run same application without conflict
    Given workspace "dev" exists with app "frontend"
    And workspace "review" exists with app "frontend"
    When I start both workspaces
    Then each workspace has isolated containers
    And network aliases don't conflict
```

These become **living documentation** that proves the spec matches reality.

---

## Part 7: Tools to Consider

### Documentation Generation & Management

| Tool | Purpose | Relevance |
|------|---------|-----------|
| [MkDocs](https://github.com/mkdocs/mkdocs) | Static site from Markdown | Docs publishing |
| [Structurizr](https://structurizr.com/) | C4 model diagrams from DSL | Architecture diagrams |
| [Log4brains](https://github.com/thomvaill/log4brains) | ADR static site generator | ADR management |
| [ADR Manager](https://github.com/adr/adr-manager) | VS Code plugin for ADRs | ADR authoring |

### Specification Consistency

| Tool | Purpose | Relevance |
|------|---------|-----------|
| [Semcheck](https://github.com/spec-check/semcheck) | LLM-based spec verification | Drift detection |
| [Vale](https://vale.sh/) | Prose linter | Terminology consistency |
| [markdownlint](https://github.com/DavidAnson/markdownlint) | Markdown linting | Format consistency |

### Specification-First Development

| Tool | Purpose | Relevance |
|------|---------|-----------|
| [GitHub Spec Kit](https://github.com/github/spec-kit) | Spec-driven development | AI-assisted implementation |
| [Cucumber](https://cucumber.io/) | BDD/executable specs | Behavior verification |

---

## Part 8: Key Recommendations

### 1. Separate Concerns by Document Type

The current PRD conflates:
- Vision/strategy (stable, rarely changes)
- Architectural decisions (captured once, immutable)
- Technical specifications (evolves with implementation)
- Reference documentation (derived from code/config)

### 2. Adopt ADRs for Decisions, but Don't Overload Them

ADRs answer "why did we choose X?" — they don't replace:
- Requirements ("what do we need?")
- Specifications ("how does it work?")
- References ("what are the options?")

### 3. Consider a Canonical Model

If adopting [Structurizr](https://structurizr.com/) or similar:
- Define architecture in a DSL
- Generate diagrams for multiple contexts
- Generate documentation fragments
- Keep diagrams consistent with reality

### 4. Embrace Layered Documentation

```
┌─────────────────────────────────────────────┐
│  Vision & Strategy (PRD-lite, stable)       │
├─────────────────────────────────────────────┤
│  Decisions (ADRs, immutable once accepted)  │
├─────────────────────────────────────────────┤
│  Architecture (C4/arc42, evolving)          │
├─────────────────────────────────────────────┤
│  Specifications (TDD, living documents)     │
├─────────────────────────────────────────────┤
│  Reference (Generated or hand-maintained)   │
├─────────────────────────────────────────────┤
│  Executable Specs (Gherkin, tests as docs)  │
└─────────────────────────────────────────────┘
```

### 5. Add Tooling for Consistency

- **Linting**: markdownlint, Vale for terminology
- **Cross-referencing**: Ensure ADRs are linked from specs
- **Automation**: Consider generating reference docs from code

### 6. Define Canonical Sources

For the Contrail project:

| Concern | Canonical Source | Derived Views |
|---------|------------------|---------------|
| CLI flags | Go code (Cobra definitions) | CLI Reference doc |
| Config schemas | Go structs + YAML | Technical spec examples |
| Decisions | ADRs in `docs/decisions/` | PRD mentions, linked |
| Behaviors | Executable specs (Gherkin) | PRD user stories |

---

## Sources

### ADR & Decision Documentation
- [ADR Homepage](https://adr.github.io/)
- [MADR GitHub](https://github.com/adr/madr)
- [AWS ADR Best Practices](https://aws.amazon.com/blogs/architecture/master-architecture-decision-records-adrs-best-practices-for-effective-decision-making/)
- [Google Cloud ADR Overview](https://cloud.google.com/architecture/architecture-decision-records)
- [Y-Statement Template](https://medium.com/olzzio/y-statements-10eb07b5a177)
- [Joel Parker Henderson ADR Examples](https://github.com/joelparkerhenderson/architecture-decision-record)

### Documentation Frameworks
- [Diátaxis Framework](https://diataxis.fr/)
- [arc42 Template](https://arc42.org/)
- [C4 Model](https://c4model.com/)
- [Structurizr](https://structurizr.com/)

### Living Documentation
- [Living Documentation Book (O'Reilly)](https://www.oreilly.com/library/view/living-documentation-continuous/9780134689418/)
- [Cucumber/Gherkin](https://cucumber.io/docs/)
- [Docs as Code](https://www.gitbook.com/blog/what-is-docs-as-code)
- [Modular Documentation (Red Hat)](https://redhat-documentation.github.io/modular-docs/)

### Design Docs & RFCs
- [Design Docs at Google](https://www.industrialempathy.com/posts/design-docs-at-google/)
- [RFCs and Design Docs (Pragmatic Engineer)](https://blog.pragmaticengineer.com/rfcs-and-design-docs/)
- [Document Types: PRD vs TRD](https://medium.com/@kokoproduct/decoding-product-development-documents-prd-vs-brd-vs-trd-a45ba984cd78)

### Single Source of Truth
- [Wikipedia - SSOT](https://en.wikipedia.org/wiki/Single_source_of_truth)
- [Red Hat SSOT Architecture](https://www.redhat.com/en/blog/single-source-truth-architecture)

### AI-Powered Tools
- [GitHub Spec Kit](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/)
- [AI Code-Documentation Sync](https://www.xugj520.cn/en/archives/ai-code-documentation-sync-tool.html)

---

## Conclusion

No single system handles all documentation types—the solution is a **layered approach** where each layer uses the right tool for its purpose, with clear ownership and linkages between layers. ADRs are excellent for decisions but should be combined with other frameworks (Diátaxis for structure, arc42 for architecture, docs-as-code for workflow) to create a cohesive, maintainable system.
