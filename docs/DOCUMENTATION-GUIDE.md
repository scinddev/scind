# Contrail Documentation Guide

**Version**: 1.1.0
**Date**: December 2024

This documentation was migrated from `specs/` on December 2024 using the Layered Documentation System.

---

## Overview

This guide explains how the Contrail documentation is organized and how to maintain it.

**Audience**: This entire `docs/` directory is **design documentation for developers building Contrail**—not end-user documentation. The layers separate different types of design artifacts by their purpose and stability, helping developers understand *what to build* (Layers 1-6) and *how to build it* (Layer 7).

---

## Documentation Layers

### Layer 1: Decisions (ADRs)

**Location**: `docs/decisions/`
**Purpose**: Capture *why* significant decisions were made

ADRs (Architecture Decision Records) document the context, decision, and consequences of significant technical and product choices. They become immutable once accepted.

**When to create an ADR**:
- Choosing between multiple viable approaches
- Making irreversible architectural decisions
- Adopting new technologies or patterns
- Changing existing patterns or conventions

**Template**: `docs/decisions/0000-template.md` (MADR Minimal)

### Layer 2: Vision

**Location**: `docs/product/`
**Purpose**: Define the product's purpose, goals, and constraints

The vision layer contains high-level product documentation that changes infrequently. It answers "what are we building and why?"

**Files**:
- `vision.md` — Problem statement, product vision, core concepts, glossary

### Layer 3: Architecture

**Location**: `docs/architecture/`
**Purpose**: Describe system structure and relationships

Architecture documentation shows how components fit together without getting into behavioral details.

**Files**:
- `overview.md` — Context and container diagrams, network topology, key components

### Layer 4: Specifications

**Location**: `docs/specs/`
**Purpose**: Define detailed behavior and rules

Specifications describe *what* the system does in detail—state machines, data schemas, edge cases, error handling.

**Files**:
- `port-types.md` — Port type system and assignment
- `configuration-schemas.md` — All configuration file formats
- `generated-override-files.md` — Override file format and generation
- `docker-labels.md` — Label conventions for containers
- `environment-variables.md` — Injected environment variables
- `workspace-lifecycle.md` — Startup, shutdown, generation sequences
- `context-detection.md` — Directory-based context detection
- `shell-integration.md` — Shell function and completion architecture
- `naming-conventions.md` — Naming patterns and collision warnings
- `proxy-infrastructure.md` — Traefik proxy setup and configuration

**Template**: `docs/specs/_template-feature.md`

### Layer 5: Reference

**Location**: `docs/reference/`
**Purpose**: Quick lookup tables and syntax

Reference documentation is optimized for scanning—command tables, option lists, configuration schemas.

**Files**:
- `cli.md` — All CLI commands and flags
- `configuration.md` — Configuration file reference with defaults

**Templates**:
- `docs/reference/_template-cli.md`
- `docs/reference/_template-config.md`

### Layer 6: Behaviors

**Location**: `features/` (project root)
**Purpose**: Executable specifications in Gherkin format

Behavior specifications define expected behavior in Given/When/Then format. These can be executed as automated tests.

**Template**: `features/_template.feature`

### Layer 7: Implementation Guides

**Location**: `docs/implementation/`
**Purpose**: Technical scaffolding and implementation details for developers building Contrail

Implementation guides describe *how to build* the system. They have a shorter lifespan than other layers—consumed during implementation, then potentially archived or absorbed into code.

**Contents**:
- Technology stack decisions with version numbers
- Project structure and scaffolding
- Executable code examples and templates
- Dependency rationale
- Implementation priority phases

**Files**:
- `go-stack.md` — Go dependencies, project structure, Cobra scaffolding, testing strategy
- `shell-scripts.md` — Complete shell scripts for Bash/Zsh/Fish (embedded code)

**Lifecycle**: These files may be archived or moved into the codebase itself (as comments, `internal/` docs, or README files in relevant packages) once implementation is complete.

---

## Layer Summary

| Layer | What to Build | How to Build It |
|-------|---------------|-----------------|
| 1. Decisions | Why we chose this approach | — |
| 2. Vision | What problem we're solving | — |
| 3. Architecture | How components relate | — |
| 4. Specifications | How features behave | — |
| 5. Reference | What commands/options exist | — |
| 6. Behaviors | What scenarios must work | — |
| 7. Implementation | — | Technology stack, scaffolding, code templates |

---

## Classification Heuristics

When adding new documentation, use these signals to determine the appropriate layer:

| If you see... | It belongs in... |
|---------------|------------------|
| "We decided...", "We chose X over Y" | Layer 1: Decisions |
| Problem statement, vision, goals, non-goals | Layer 2: Vision |
| Diagrams, component relationships, topology | Layer 3: Architecture |
| "When X happens, Y occurs", state machines, schemas | Layer 4: Specifications |
| Command tables, option lists, config defaults | Layer 5: Reference |
| Given/When/Then, test scenarios | Layer 6: Behaviors |
| Dependency versions, code scaffolding, project structure | Layer 7: Implementation Guides |

---

## Decision Tree

```
Is this explaining WHY a choice was made?
├─ Yes → Layer 1: Decisions (ADR)
└─ No ↓

Is this about product purpose, goals, or concepts?
├─ Yes → Layer 2: Vision
└─ No ↓

Is this about system structure and component relationships?
├─ Yes → Layer 3: Architecture
└─ No ↓

Is this about detailed behavior or data schemas?
├─ Yes → Layer 4: Specifications
└─ No ↓

Is this a quick-reference table or lookup?
├─ Yes → Layer 5: Reference
└─ No ↓

Is this a testable scenario?
├─ Yes → Layer 6: Behaviors
└─ No ↓

Is this implementation scaffolding, code templates, or dependency lists?
├─ Yes → Layer 7: Implementation Guides
└─ No → Probably doesn't belong in docs
```

---

## Cross-Layer Linking

Documents should link to related content in other layers:

- **Specifications → ADRs**: Link to decisions that explain *why* a behavior exists
- **Architecture → Specifications**: Link to specs for detailed behavioral descriptions
- **Reference → Specifications**: Link to specs for context beyond the quick-reference
- **Implementation → Specifications**: Link to specs that the code implements

Example:
```markdown
See [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md) for rationale.
```

---

## Tooling (Tier 2: Recommended)

### Linting

**markdownlint** — Enforces consistent Markdown formatting
```bash
npm install --save-dev markdownlint-cli
npx markdownlint "docs/**/*.md"
```

### Prose Quality

**Vale** — Checks writing style and terminology
```bash
brew install vale
vale docs/
```

### ADR Management

**Log4brains** — Visualize and manage ADRs
```bash
npm install -g log4brains
log4brains preview
```

### Architecture Diagrams

**Structurizr** — Render C4 diagrams from DSL
```bash
# Use Structurizr Lite Docker container
docker run -it --rm -p 8080:8080 -v $(pwd)/docs/architecture:/usr/local/structurizr structurizr/lite
```

---

## Directory Structure

```
docs/
├── DOCUMENTATION-GUIDE.md          # This file
├── decisions/                       # Layer 1: ADRs
│   ├── 0000-template.md
│   ├── 0001-docker-compose-project-name-isolation.md
│   ├── 0002-two-layer-networking.md
│   ├── 0003-pure-overlay-design.md
│   ├── 0004-convention-based-naming.md
│   ├── 0005-structure-vs-state-separation.md
│   ├── 0006-three-configuration-schemas.md
│   ├── 0007-port-type-system.md
│   ├── 0008-traefik-reverse-proxy.md
│   ├── 0009-flexible-tls-configuration.md
│   ├── 0010-up-down-command-semantics.md
│   └── 0011-options-based-targeting.md
├── product/                         # Layer 2: Vision
│   └── vision.md
├── architecture/                    # Layer 3: Architecture
│   └── overview.md
├── specs/                           # Layer 4: Specifications
│   ├── _template-feature.md
│   ├── port-types.md
│   ├── configuration-schemas.md
│   ├── generated-override-files.md
│   ├── docker-labels.md
│   ├── environment-variables.md
│   ├── workspace-lifecycle.md
│   ├── context-detection.md
│   ├── shell-integration.md
│   ├── naming-conventions.md
│   └── proxy-infrastructure.md
├── reference/                       # Layer 5: Reference
│   ├── _template-cli.md
│   ├── _template-config.md
│   ├── cli.md
│   └── configuration.md
└── implementation/                  # Layer 7: Implementation Guides
    ├── go-stack.md
    └── shell-scripts.md

features/                            # Layer 6: Behaviors (project root)
└── _template.feature
```

---

## Maintenance Guidelines

1. **Keep layers separate**: Don't mix decision rationale into specs, or behavioral details into architecture
2. **Link don't duplicate**: Reference other layers instead of copying content
3. **Update all affected layers**: A feature change may require updates to specs, reference, and behaviors
4. **ADRs are immutable**: Supersede rather than edit accepted ADRs
5. **Templates are guides**: Use templates for consistency, but adapt as needed
6. **Implementation guides have short lifespans**: Archive or migrate to code once implementation is complete

---

## Migration Notes

The following files were migrated from the original `specs/` directory:

| Original File | Migrated To |
|---------------|-------------|
| `contrail-prd.md` | Decisions (11 ADRs), Vision, Architecture, Specs |
| `contrail-technical-spec.md` | Architecture, Specs (multiple) |
| `contrail-cli-reference.md` | Reference (`cli.md`), Specs (`context-detection.md`) |
| `contrail-shell-integration.md` | Specs (`shell-integration.md`), Implementation (`shell-scripts.md`) |
| `contrail-go-stack.md` | Implementation (`go-stack.md`) |

**Original files are preserved** in `specs/` at project root for reference during the transition period.
