# Documentation Guide

**Project**: Contrail
**Generated**: 2026-01-02
**Tooling**: Tier 2 (markdownlint, Vale, Log4brains, Structurizr)

---

## Overview

This documentation follows the Layered Documentation System, organizing content by type and purpose. Each layer serves a specific audience and use case.

---

## Structure

```
docs/
├── decisions/          # Layer 1: ADRs (Architectural Decision Records)
├── product/            # Layer 2: Vision, comparison, roadmap
├── architecture/       # Layer 3: System architecture and diagrams
├── specs/              # Layer 4: Detailed feature specifications
├── reference/          # Layer 5: CLI and configuration reference
├── implementation/     # Layer 7: Technology stack and patterns
└── .migration/         # Migration plan (delete after migration)

features/               # Layer 6: Gherkin behavior scenarios (at project root)
```

---

## Layers

### Layer 1: Decisions (`decisions/`)

**Purpose**: Capture architectural decisions and their rationale.

**Audience**: Developers, architects, future maintainers.

**Format**: MADR Minimal (Status, Context, Decision, Consequences).

**When to write**: When making significant architectural or design decisions.

**Example**: [ADR-0001: Docker Compose Project Name Isolation](./decisions/0001-docker-compose-project-name-isolation.md)

### Layer 2: Product (`product/`)

**Purpose**: Communicate vision, value proposition, and roadmap.

**Audience**: Product managers, stakeholders, new team members.

**Contents**:
- Vision — Why Contrail exists, core concepts
- Comparison — How Contrail differs from alternatives
- Roadmap — Future considerations

### Layer 3: Architecture (`architecture/`)

**Purpose**: Explain how the system works at a high level.

**Audience**: Developers needing system understanding.

**Contents**:
- System overview
- Component diagrams
- Data flow descriptions
- Network topology

**Diagram format**: ASCII (inline) or Structurizr DSL (in appendices).

### Layer 4: Specifications (`specs/`)

**Purpose**: Define detailed behavior for implementation.

**Audience**: Developers implementing or maintaining features.

**Format**: Structured template with:
- Overview
- Behavior (normal flow)
- Data Schema (with tables)
- Examples
- Edge Cases
- Error Handling

**When to write**: Before implementing a feature or when documenting existing behavior.

### Layer 5: Reference (`reference/`)

**Purpose**: Quick lookup for commands and configuration.

**Audience**: Users running commands, writing configuration.

**Format**: Dense tables, minimal prose. Optimized for scanning.

**Contents**:
- CLI reference — All commands, options, exit codes
- Configuration reference — All config files and fields

### Layer 6: Behaviors (`features/`)

**Purpose**: Executable specifications in Gherkin.

**Audience**: QA, developers, stakeholders wanting to understand behavior.

**Location**: `features/` at project root (standard convention).

**Format**: Gherkin `.feature` files with Given/When/Then scenarios.

### Layer 7: Implementation (`implementation/`)

**Purpose**: Document technology choices and patterns.

**Audience**: Developers working on the codebase.

**Contents**:
- Tech stack — Dependencies and rationale
- Code patterns — Error handling, testing strategy
- Scaffolding — In appendices

---

## Content Thresholds

Move large content to appendices (`appendices/{feature-name}/`) when:

| Content Type | Threshold |
|--------------|-----------|
| Code blocks | >50 lines |
| Step lists | >10 items |
| Tables | >20 rows |
| Embedded images | Any (use appendices) |

---

## Cross-References

Documents should link to related content:

- ADRs → link to implementing specs
- Specs → link to motivating ADRs
- Architecture → links to both ADRs and specs
- Reference → links to specs for detailed behavior

Use relative paths: `../specs/context-detection.md`

---

## Tooling (Tier 2)

### markdownlint

Enforces Markdown consistency.

```bash
npx markdownlint-cli2 "docs/**/*.md"
```

Configuration: `.markdownlint.yaml` or `.markdownlint-cli2.yaml`

### Vale

Prose linting for style and terminology.

```bash
vale docs/
```

Configuration: `.vale.ini` and `styles/` directory.

### Log4brains

ADR management and visualization.

```bash
log4brains preview     # Local preview
log4brains build       # Static site
```

Configuration: `.log4brains.yaml`

### Structurizr

Architecture diagrams as code.

```bash
# Using Structurizr Lite (Docker)
docker run -it --rm -p 8080:8080 -v $(pwd)/docs/architecture:/usr/local/structurizr structurizr/lite
```

DSL files: `docs/architecture/appendices/` or `docs/architecture/*.dsl`

---

## Writing Conventions

### Tone

- Technical and direct
- Use imperative mood for instructions ("Run the command", not "You should run")
- Avoid marketing language

### Formatting

- Use code fences with language hints
- Use tables for structured data
- Keep paragraphs short (3-5 sentences)
- Use headers to enable scanning

### Naming

- Files: `kebab-case.md`
- Directories: `kebab-case/`
- ADRs: `NNNN-title.md` (e.g., `0001-docker-compose-project-name-isolation.md`)

---

## Migration Status

This documentation was migrated from existing specs. See:

- [Migration Plan](./.migration/README.md) — Step-by-step migration instructions
- [Audit](./audit.md) — Post-migration verification checklist

After migration is complete, delete the `.migration/` directory.

---

## Related

- [Layered Documentation System](../layered-docs-system/LAYERED-DOCUMENTATION-SYSTEM.md) — The system this guide implements

