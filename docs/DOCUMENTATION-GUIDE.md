# Contrail Documentation Guide

This guide explains how to navigate and contribute to Contrail's documentation.

---

## Documentation Structure

Contrail uses a **7-layer documentation system** that separates concerns by type and audience:

```
docs/
├── DOCUMENTATION-GUIDE.md      # This file
├── decisions/                  # Layer 1: ADRs
├── product/                    # Layer 2: Vision
├── architecture/               # Layer 3: Architecture
├── specs/                      # Layer 4: Specifications
├── reference/                  # Layer 5: Reference
├── implementation/             # Layer 7: Implementation
└── migration/                  # Migration notes

features/                       # Layer 6: Behaviors (Gherkin)
```

---

## Layer Overview

### Layer 1: Decisions (`docs/decisions/`)

**Architectural Decision Records (ADRs)** capture the *why* behind significant technical choices.

- **Format**: MADR Minimal
- **When to use**: Recording decisions with lasting impact
- **Key files**: `README.md` (index), `0000-template/` (template)

**Current ADRs**:
| # | Title | Status |
|---|-------|--------|
| 0001 | Docker Compose Project Name Isolation | Accepted |
| 0002 | Two-Layer Networking | Accepted |
| 0003 | Pure Overlay Design | Accepted |
| 0004 | Convention-Based Naming | Accepted |
| 0005 | Structure vs State Separation | Accepted |
| 0006 | Three Configuration Schemas | Accepted |
| 0007 | Port Type System | Accepted |
| 0008 | Traefik for Reverse Proxy | Accepted |
| 0009 | Flexible TLS Configuration | Accepted |
| 0010 | up/down Command Semantics | Accepted |
| 0011 | Options-Based Targeting with Context Detection | Accepted |
| 0012 | Layered Documentation System | Accepted |

### Layer 2: Vision (`docs/product/`)

**Product vision and requirements** at the strategic level.

- **Format**: Lean PRD
- **When to use**: Understanding product goals and scope
- **Key files**: `vision/README.md`

### Layer 3: Architecture (`docs/architecture/`)

**System structure and component relationships**.

- **Format**: C4-Lite (Context + Container diagrams)
- **When to use**: Understanding system design
- **Key files**: `overview/README.md`

### Layer 4: Specifications (`docs/specs/`)

**Precise behavioral definitions** for features.

- **Format**: Feature Spec
- **When to use**: Implementing or testing specific behaviors
- **Key files**: See `docs/specs/README.md` for index

**Current specifications**:
- Configuration Schemas
- Context Detection
- Docker Labels
- Environment Variables
- Generated Override Files
- Naming Conventions
- Port Types
- Proxy Infrastructure
- Shell Integration
- Workspace Lifecycle

### Layer 5: Reference (`docs/reference/`)

**Lookup documentation** for commands and configuration.

- **Format**: Reference tables and structured lists
- **When to use**: Quick lookup of specific options
- **Key files**: `cli/README.md`, `configuration/README.md`

### Layer 6: Behaviors (`features/`)

**Executable specifications** in Gherkin format.

- **Format**: Standard BDD (Given/When/Then)
- **When to use**: Defining acceptance criteria
- **Key files**: `_template.feature`

### Layer 7: Implementation (`docs/implementation/`)

**Technology stack and implementation guidance**.

- **Format**: Tech Stack
- **When to use**: Setting up development environment
- **Key files**: `go-stack/README.md`

---

## Finding Information

### "Why was this decision made?"
→ Check **Layer 1: Decisions** for relevant ADRs

### "What is Contrail trying to solve?"
→ Check **Layer 2: Vision** for product context

### "How does the system work?"
→ Check **Layer 3: Architecture** for system design

### "How does this feature behave?"
→ Check **Layer 4: Specifications** for precise definitions

### "What command should I use?"
→ Check **Layer 5: Reference** for CLI and config options

### "What should happen in this scenario?"
→ Check **Layer 6: Behaviors** for Gherkin scenarios

### "How do I build this?"
→ Check **Layer 7: Implementation** for tech stack details

---

## Contributing to Documentation

### Adding a New ADR

1. Copy `docs/decisions/0000-template/README.md`
2. Increment the number
3. Fill in Context, Decision, Consequences
4. Update `docs/decisions/README.md` index

### Adding a New Specification

1. Create directory under `docs/specs/{spec-name}/`
2. Create `README.md` following the spec template
3. Update `docs/specs/README.md` index
4. Cross-reference related ADRs

### Adding a New Feature Behavior

1. Copy `features/_template.feature`
2. Define scenarios using Given/When/Then
3. Follow Gherkin best practices

---

## Navigation Tips

- **Cross-references**: Documents link to related content in other layers
- **README.md files**: Each directory has an index
- **Version headers**: Major documents include version and date

---

## Maintenance

This documentation system was installed using the Layered Documentation System v3.0.

For system documentation, see: `layered-docs-system/LAYERED-DOCUMENTATION-SYSTEM.md`
