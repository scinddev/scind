# Comprehensive Analysis: Documentation Methodologies for Evolving Developer Documentation

**Date**: January 2026
**Purpose**: Analysis of documentation methodologies for the Contrail project specifications
**Related Research**: See `documentation-methodologies-research.md` for detailed methodology research
**Implementation**: See `LAYERED-DOCS-SYSTEM-ANALYSIS.md` for the merged system and `layered-docs-system/` for operational guides

> **Note**: This analysis has been incorporated into the Layered Documentation System v3.0.
> Key additions in v3.0 include appendix support, confidence-based classification, and
> content thresholds to prevent information loss during migration.

---

## Table of Contents

1. [Current Documentation Inventory](#current-documentation-inventory)
2. [The Core Problem: Specification Drift](#the-core-problem-specification-drift)
3. [Recommended Layered Documentation System](#recommended-layered-documentation-system)
4. [Document Hierarchy and Authority](#document-hierarchy-and-authority)
5. [Single Source of Truth (SSOT) Strategy](#single-source-of-truth-ssot-strategy)
6. [Preventing Agent-Generated Conflicts](#preventing-agent-generated-conflicts)
7. [Recommended Directory Structure](#recommended-directory-structure)
8. [Methodology Blending Recommendation](#methodology-blending-recommendation)
9. [Evolution Path for Contrail](#evolution-path-for-contrail)
10. [Key Takeaways](#key-takeaways)

---

## Current Documentation Inventory

Looking at the `specs/` folder, these document types currently exist:

| Document | Type | Purpose |
|----------|------|---------|
| `contrail-prd.md` | Product Requirements | What to build and why |
| `contrail-technical-spec.md` | Technical Specification | How it works (architecture, schemas, operations) |
| `contrail-cli-reference.md` | Reference Documentation | Complete CLI command reference |
| `contrail-shell-integration.md` | Feature Specification | Shell function behavior and scripts |
| `contrail-go-stack.md` | Implementation Guide | Technology stack and scaffolding |

The documents are well-structured but organically evolved rather than following a formal methodology. This is common—and the good news is that the existing documents map fairly cleanly onto established methodologies.

---

## The Core Problem: Specification Drift

The central challenge: **agentic tools generate feature specs that conflict with each other and with implementation reality**. This happens because:

1. **No single source of truth** for decisions
2. **Duplicated information** across documents that drifts
3. **No explicit hierarchy** showing which document is authoritative for what
4. **Design-phase documents become stale** as implementation proceeds

---

## Recommended Layered Documentation System

Based on the research, here's a system designed for these specific needs:

### Layer 1: Architectural Decision Records (ADRs)

**Format**: MADR (Markdown Architectural Decision Records)

**Why ADR is crucial for this problem**: ADRs capture the *why* behind decisions. When Agent OS or similar tools generate conflicting specs, the ADRs serve as the authoritative source for design decisions. If a generated spec contradicts an ADR, the ADR wins.

**What belongs in ADRs**:
- Technology choices (Docker Compose project name isolation, Two-layer networking)
- Design patterns (Pure overlay, Convention-based naming)
- Tradeoff decisions (Structure vs. state separation)

**The PRD already contains these!** The "Key Architectural Decisions" section in `contrail-prd.md` (lines 162-314) is essentially a set of ADRs embedded in the PRD. These should be extracted into individual ADR files.

**Recommended structure**:
```
docs/
├── decisions/
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
```

**MADR template** (simplified):

```markdown
# [short title describing the decision]

**Status**: Proposed | Accepted | Deprecated | Superseded by [ADR-XXXX]
**Date**: YYYY-MM-DD

## Context

[What is the issue we're addressing? What forces are at play?]

## Decision

[What is our response to these forces?]

## Consequences

### Positive
- [Good outcome 1]
- [Good outcome 2]

### Negative
- [Tradeoff or downside 1]

## Alternatives Considered

### [Alternative 1]
- **Pros**: [advantages]
- **Cons**: [disadvantages]
- **Why rejected**: [reason]
```

---

### Layer 2: Product Requirements Document (Lean PRD)

**Keep the existing `contrail-prd.md`** but refactor it:

**What stays in PRD**:
- Executive Summary
- Problem Statement
- Product Vision
- Core Concepts (terminology)
- Success Criteria
- Non-Goals
- Known Limitations

**What moves out**:
- Key Architectural Decisions → Individual ADR files
- Detailed Configuration Examples → Technical Spec
- CLI Quick Reference → CLI Reference

**Why**: The PRD should answer "What are we building and why?" It should NOT contain implementation details. When implementation changes, you shouldn't need to update the PRD.

---

### Layer 3: Technical Specification (Architecture + Behavior)

The `contrail-technical-spec.md` is doing double-duty:
1. **Architecture documentation** (network topology, directory structure)
2. **Detailed behavior specification** (operations, generation logic)

**Consider splitting using arc42-inspired sections**:

| Section | Current Location | Recommendation |
|---------|-----------------|----------------|
| Architecture diagrams | Technical Spec | Keep in Tech Spec |
| Configuration schemas | Technical Spec | Keep in Tech Spec |
| Operations (startup/shutdown sequences) | Technical Spec | Could become Gherkin feature files |
| Naming conventions | Technical Spec | Keep in Tech Spec |
| Docker labels | Technical Spec | Could be separate feature spec |

**The key insight**: The technical spec is currently ~1400 lines. That's maintainable, but consider which parts could be **executable specifications** (Gherkin) vs. **reference documentation**.

---

### Layer 4: Behavior Specifications (Gherkin Feature Files)

This is the **highest-value addition** possible. Gherkin files serve triple duty:

1. **Specification** — defines expected behavior
2. **Documentation** — always accurate because it's tested
3. **Tests** — executable validation

**What belongs in Gherkin**:

```gherkin
# features/workspace-lifecycle.feature

Feature: Workspace Lifecycle
  As a developer
  I want to manage workspace environments
  So that I can run isolated development stacks

  Rule: Workspace names must be unique

    Scenario: Initialize new workspace
      Given no workspace named "dev" exists
      When I run "contrail workspace init --workspace=dev"
      Then a workspace.yaml file should be created
      And the workspace should be registered in the global registry

    Scenario: Prevent duplicate workspace names
      Given a workspace named "dev" exists at "/workspaces/dev"
      When I run "contrail workspace init --workspace=dev" in "/other-location"
      Then the command should fail with exit code 3
      And the error should mention the existing workspace path

  Rule: Proxy must be running for workspace up

    Scenario: Auto-start proxy on workspace up
      Given the proxy is not running
      When I run "contrail workspace up --workspace=dev"
      Then the proxy should be started automatically
      And the workspace network should be created
      And all applications should be started
```

**What this solves**:
- **Living documentation** — if the behavior changes, the test breaks
- **Contradiction detection** — if a Gherkin scenario contradicts an ADR, you catch it immediately
- **Specification by example** — concrete examples clarify ambiguous requirements

**Recommended Gherkin files for Contrail**:

```
features/
├── workspace-lifecycle.feature
├── application-management.feature
├── flavor-switching.feature
├── port-assignment.feature
├── proxy-management.feature
├── context-detection.feature
├── override-generation.feature
├── shell-integration.feature
```

---

### Layer 5: Reference Documentation

The `contrail-cli-reference.md` and `contrail-shell-integration.md` are excellent reference documents. Keep them, but ensure they're structured for **lookup, not reading**.

**CLI Reference best practices** (from clig.dev):
- Every command documented
- Every flag with description and default
- Examples for common operations
- Structured consistently

The CLI reference already follows these patterns well.

**Configuration Reference pattern** (for schemas):

```markdown
## Configuration: `workspace.yaml`

### `workspace.name`
**Type**: string (required)
**Pattern**: lowercase alphanumeric with hyphens
**Example**: `dev`, `feature-x`, `pr-123`

Used as prefix for Docker Compose project names and hostnames.

### `workspace.applications`
**Type**: map<string, ApplicationRef>
**Required**: At least one application

Maps application names to their configuration...
```

---

### Layer 6: Implementation Guides

The `contrail-go-stack.md` is an implementation guide. These are valuable for:
- Onboarding new developers
- Ensuring consistency across implementation
- Scaffolding patterns

**Keep these separate from specifications**. Implementation details change more frequently than behavior specifications.

---

## Document Hierarchy and Authority

**Critical for preventing drift**: Establish a clear hierarchy of authority.

```
┌─────────────────────────────────────────────────────────┐
│                      MOST AUTHORITATIVE                  │
├─────────────────────────────────────────────────────────┤
│  ADRs (Architectural Decision Records)                  │
│  - Decisions are immutable once accepted                │
│  - New decisions supersede (not modify) old ones        │
│  - If anything conflicts with ADR, ADR wins             │
├─────────────────────────────────────────────────────────┤
│  Gherkin Feature Files                                   │
│  - Executable specifications                             │
│  - Tests enforce behavior                                │
│  - If test passes, documentation is accurate             │
├─────────────────────────────────────────────────────────┤
│  PRD (Product Requirements)                              │
│  - High-level "what" and "why"                          │
│  - Rarely changes after design phase                     │
│  - References ADRs for decisions                         │
├─────────────────────────────────────────────────────────┤
│  Technical Specification                                 │
│  - Architecture and schemas                              │
│  - References ADRs for rationale                         │
│  - Updated as architecture evolves                       │
├─────────────────────────────────────────────────────────┤
│  Reference Documentation (CLI, Config)                   │
│  - Factual, complete, lookup-oriented                    │
│  - Generated where possible                              │
│  - Updated with implementation                           │
├─────────────────────────────────────────────────────────┤
│  Implementation Guides (Go Stack)                        │
│  - How to build, patterns to follow                      │
│  - Most volatile, changes with tech decisions            │
│                     LEAST AUTHORITATIVE                  │
└─────────────────────────────────────────────────────────┘
```

---

## Single Source of Truth (SSOT) Strategy

**Every fact should be mastered in exactly one place**:

| Information | Mastered In | Referenced From |
|-------------|-------------|-----------------|
| Why we chose Docker Compose project names | ADR-0001 | PRD, Tech Spec |
| The workspace.yaml schema | Tech Spec (Configuration Schemas) | CLI Reference, Go Stack |
| What `contrail workspace up` does | CLI Reference | Tech Spec (Operations) |
| Shell function implementation | Shell Integration Spec | Go Stack |
| Port assignment algorithm | Tech Spec (Global State) | ADR-0007 |

**Key principle**: When you find yourself copying information, create a reference instead. Markdown supports this:

```markdown
For details on the port assignment algorithm, see
[Global State Configuration](./contrail-technical-spec.md#global-state).
```

---

## Preventing Agent-Generated Conflicts

When tools like Agent OS generate feature specs, apply this workflow:

1. **Check ADRs first** — Does the generated spec contradict any accepted decision?
2. **Check Gherkin scenarios** — Would implementing this break existing tests?
3. **Identify new decisions** — If the generated spec implies new decisions, create ADRs
4. **Update one source** — If the spec is valid, update the authoritative source (not multiple places)

**Automated validation** that could be added:
- Link checking (ensure cross-references are valid)
- ADR consistency checking (no contradictions between ADRs)
- Gherkin test execution (behavior matches specification)

---

## Recommended Directory Structure

> **Updated Structure**: See `LAYERED-DOCS-SYSTEM-ANALYSIS.md` for the v2.0 directory structure
> which uses `{topic}/README.md` + `appendices/` pattern for all document types.

```
contrail/
├── docs/
│   ├── DOCUMENTATION-GUIDE.md        # Thresholds and project config (NEW)
│   │
│   ├── decisions/                    # ADRs (Layer 1)
│   │   └── NNNN-{title}/
│   │       ├── README.md
│   │       └── appendices/           # Extended analysis, background
│   │
│   ├── product/                      # PRD (Layer 2)
│   │   └── vision/
│   │       └── README.md
│   │
│   ├── architecture/                 # Technical Spec (Layer 3)
│   │   └── overview/
│   │       ├── README.md
│   │       └── appendices/           # Detailed diagrams
│   │
│   ├── specs/                        # Feature Specifications (Layer 4)
│   │   └── {feature}/
│   │       ├── README.md
│   │       └── appendices/           # Schemas, examples, error catalogs
│   │
│   ├── reference/                    # Reference Docs (Layer 5)
│   │   ├── cli/
│   │   │   ├── README.md
│   │   │   └── appendices/           # Detailed examples, errors
│   │   └── configuration/
│   │       ├── README.md
│   │       └── appendices/           # Complete examples, schemas
│   │
│   ├── implementation/               # Implementation Guides (Layer 7)
│   │   └── tech-stack/
│   │       ├── README.md
│   │       └── appendices/           # Full scaffold scripts
│   │
│   ├── migration/                    # Medium-confidence content (NEW)
│   └── blackhole/                    # Unclassified content (NEW)
│
└── features/                         # Gherkin specs (Layer 6)
    ├── workspace-lifecycle.feature
    └── ...
```

---

## Methodology Blending Recommendation

Based on these needs, here's the recommended blend:

| Need | Methodology | Why |
|------|-------------|-----|
| Record decisions | **MADR** (Markdown ADR) | Structured, includes alternatives considered |
| Define behavior | **Gherkin** | Executable, prevents drift |
| Architecture docs | **C4 Model** (Levels 1-2) + **arc42-lite** | Visual + structured text |
| Reference docs | **Diátaxis (Reference quadrant)** | Factual, lookup-oriented |
| Process | **Docs-as-Code** | Same repo, same PR, same review |

---

## Evolution Path for Contrail

### Phase 1: Extract ADRs (Immediate)

Extract the 11 architectural decisions from the PRD into individual ADR files. This is quick and provides immediate value.

**Decisions to extract**:
1. Docker Compose Project Name Isolation
2. Two-Layer Networking
3. Pure Overlay (Applications Remain Workspace-Agnostic)
4. Convention-Based Naming
5. Structure vs State Separation
6. Three Configuration Schemas
7. Port Type System for Exported Services
8. Traefik for Reverse Proxy
9. Flexible TLS Configuration
10. `up`/`down` Command Semantics
11. Options-Based Targeting with Context Detection

### Phase 2: Add Gherkin Specifications (Near-term)

Start with the most complex/contentious behaviors:
- Context detection algorithm
- Port assignment and conflict resolution
- Workspace up/down sequences

### Phase 3: Refactor Document Hierarchy (Medium-term)

- Slim down PRD to vision/goals only
- Ensure Tech Spec references ADRs (doesn't duplicate rationale)
- Add explicit cross-references between documents

### Phase 4: Automate Validation (Longer-term)

- CI pipeline runs Gherkin tests
- Link checking validates cross-references
- ADR tooling generates index/graph

---

## Key Takeaways

1. **ADRs are essential** — They're the antidote to specification drift. When anything conflicts, the ADR wins.

2. **Gherkin is your living documentation** — Executable specifications can't lie. They're always accurate or they fail.

3. **Single Source of Truth** — Every fact in one place, everything else references it.

4. **Clear hierarchy** — Know which document is authoritative for what.

5. **Docs-as-Code** — Documentation in the same repo, reviewed in the same PR, deployed the same way.

6. **The existing docs are mostly there** — The main gaps are extracted ADRs and executable specifications.

---

## Next Steps

Potential follow-up actions:

1. Create a template ADR from one of the existing architectural decisions
2. Draft a Gherkin feature file for a specific Contrail behavior
3. Create a more detailed migration plan for restructuring the current docs
4. Explore any specific methodology in more depth

---

## Related Documents

- `documentation-methodologies-research.md` — Comprehensive research on documentation methodologies
- `LAYERED-DOCS-SYSTEM-ANALYSIS.md` — Analysis and recommendations for merged system
- `layered-docs-system/` — Implementation of the layered documentation system (v3.0)
- `specs/contrail-prd.md` — Current Product Requirements Document
- `specs/contrail-technical-spec.md` — Current Technical Specification
- `specs/contrail-cli-reference.md` — Current CLI Reference
- `specs/contrail-shell-integration.md` — Current Shell Integration Spec
- `specs/contrail-go-stack.md` — Current Go Stack Specification
