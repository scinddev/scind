# ADR-0012: Layered Documentation System

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Team

## Context

The original Contrail documentation consisted of five monolithic specification files:

- `contrail-prd.md` — Product requirements, decisions, concepts, and some technical details
- `contrail-technical-spec.md` — Architecture, configuration schemas, behaviors, and reference tables
- `contrail-cli-reference.md` — CLI commands with embedded context detection specification
- `contrail-shell-integration.md` — Shell function specifications
- `contrail-go-stack.md` — Implementation guide

This structure had several issues:

1. **Mixed concerns**: Decision rationale was embedded in specifications, making it hard to find "why" something was designed a certain way
2. **Difficult navigation**: Finding specific information required searching through long files
3. **Overlapping content**: Similar information appeared in multiple documents with subtle differences
4. **No clear ownership**: It was unclear which document to update when adding new features

## Decision

Adopt the Layered Documentation System, organizing documentation into six distinct layers based on purpose and stability:

| Layer | Purpose | Location |
|-------|---------|----------|
| 1. Decisions | Why choices were made (ADRs) | `docs/decisions/` |
| 2. Vision | Product purpose, goals, concepts | `docs/product/` |
| 3. Architecture | System structure and relationships | `docs/architecture/` |
| 4. Specifications | Detailed behavior and rules | `docs/specs/` |
| 5. Reference | Quick-lookup tables | `docs/reference/` |
| 6. Behaviors | Executable Gherkin specs | `features/` |

Migrate existing content from the five source files into this layered structure.

## Consequences

### Positive

- Clear separation of "why" (decisions) from "what" (specs) from "how to use" (reference)
- Easier to find specific information
- Each layer has clear ownership and update triggers
- Supports different reading patterns (quick lookup vs. deep understanding)
- Templates ensure consistency for new documentation
- Cross-layer linking connects related information

### Negative

- Initial migration effort required
- More files to manage
- Contributors need to understand the layer system
- Risk of content fragmentation if not maintained properly

### Neutral

- Original specification files are preserved for reference during transition
- The Go Stack document (`contrail-go-stack.md`) was not migrated as it's an implementation guide rather than product documentation
- Tooling recommendations (Tier 2) are documented but not enforced

---

## Notes

Migration performed December 2024. The following content was extracted:

- **11 ADRs** from embedded decision blocks in `contrail-prd.md`
- **Vision document** consolidated from PRD executive summary, problem statement, and concepts
- **Architecture overview** consolidated from PRD network topology and technical spec architecture
- **8 specification documents** organized by feature area
- **2 reference documents** for CLI and configuration
- **Gherkin template** for future behavior specifications
