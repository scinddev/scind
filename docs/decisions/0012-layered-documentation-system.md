# ADR-0012: Layered Documentation System

**Status**: Accepted

## Context

As the Scind project grew, documentation accumulated organically without a consistent organizational structure. Different types of content—decisions, specifications, reference material, and implementation guides—were intermixed or duplicated across files. This created several problems:

- **Discoverability**: Finding authoritative information required searching multiple files
- **Maintenance burden**: Updates needed to be made in multiple places
- **Unclear authority**: When documents conflicted, there was no defined hierarchy to resolve disputes
- **Inconsistent depth**: Some topics had excessive detail inline while others lacked sufficient context

We needed a systematic approach to organize documentation that:
1. Defines clear categories for different types of content
2. Establishes authority hierarchy when content conflicts
3. Supports both quick reference and deep-dive reading
4. Scales as the project grows without becoming unwieldy

## Decision

Adopt the Layered Documentation System (LDS) with seven distinct layers, each serving a specific purpose:

| Layer | Directory | Purpose | Stability |
|-------|-----------|---------|-----------|
| 1. Decisions | `decisions/` | Capture WHY choices were made | Immutable |
| 2. Vision | `product/` | Define WHAT we're building | Stable |
| 3. Architecture | `architecture/` | Show HOW components relate | Evolving |
| 4. Specifications | `specs/` | Detail HOW features work | Living |
| 5. Reference | `reference/` | Provide lookup tables | Generated/maintained |
| 6. Behaviors | `features/` | Verify expected behaviors | Executable |
| 7. Implementation | `implementation/` | Guide HOW to build | Short-lived |

Key principles of the system:
- **Single Source of Truth**: Each fact lives in exactly one place
- **Linkage Over Duplication**: Reference other documents rather than copying
- **Appendix for Scale**: Large content (code blocks >50 lines, complete files) moves to appendices
- **Authority Hierarchy**: ADRs > Gherkin > Vision > Specifications > Reference > Implementation

## Consequences

### Positive

- Clear placement rules via classification decision tree and heuristics
- Defined authority hierarchy resolves conflicts unambiguously
- Appendix system keeps main documents scannable while preserving detail
- Cross-layer linking creates navigable documentation graph
- Maintenance workflows (audit, refine, sync, update) provide operational guidance
- AI agents can follow classification heuristics consistently

### Negative

- Initial migration effort required to reorganize existing content
- Contributors must learn the layer system before adding documentation
- More files to navigate compared to monolithic documentation
- Appendix indirection adds one level of lookup for detailed content

### Neutral

- Documentation structure is now convention-driven rather than ad-hoc
- The `DOCUMENTATION-GUIDE.md` file becomes the authoritative reference for documentation practices

## Related Documents

- [Documentation Guide](../DOCUMENTATION-GUIDE.md) — Complete LDS reference with classification heuristics and thresholds
