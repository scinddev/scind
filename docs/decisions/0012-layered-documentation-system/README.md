# ADR-0012: Layered Documentation System

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

The Contrail project has extensive documentation covering product vision, technical specifications, CLI reference, implementation details, and architectural decisions. Without clear organization, documentation becomes difficult to navigate and maintain, especially for AI agents working with the codebase.

## Decision

Adopt a 7-layer documentation system that separates concerns by type and audience:

| Layer | Purpose | Example Content |
|-------|---------|-----------------|
| 1. Decisions (ADRs) | Why choices were made | Architecture decisions, design rationale |
| 2. Vision | What and why at product level | PRD, product vision, goals |
| 3. Architecture | How the system is structured | Network topology, component diagrams |
| 4. Specifications | Precise behavioral definitions | Configuration schemas, environment variables |
| 5. Reference | Lookup documentation | CLI commands, config options |
| 6. Behaviors | Executable specifications | Gherkin scenarios (future) |
| 7. Implementation | How to build it | Technology stack, code scaffolding |

## Consequences

### Positive

- Clear navigation path for different information needs
- AI agents can efficiently locate relevant documentation
- Reduces duplication by having a single source for each type of content
- Supports documentation-as-code workflows

### Negative

- Initial migration effort to restructure existing documentation
- Requires discipline to maintain layer boundaries

### Neutral

- Content that spans layers is placed in the most specific applicable layer

---

## Notes

This ADR documents the adoption of the Layered Documentation System. The system itself is documented in `layered-docs-system/LAYERED-DOCUMENTATION-SYSTEM.md`.
