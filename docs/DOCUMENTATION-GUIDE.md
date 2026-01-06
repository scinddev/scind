# Contrail Documentation Guide

**For AI Agents and Contributors**: This guide explains how the documentation in this directory is organized and how to maintain it.

---

## Migration Note

This documentation was migrated from `specs/` on 2026-01-05 using the Layered Documentation System (LDS).

Migration step files are in `.migration/` — execute them in separate sessions to complete the migration.

---

## Glossary

### Directory Terminology

| Term | Definition |
|------|------------|
| `DOCS_DIR` | The documentation root directory (`docs/`) |
| `LEGACY_DOCS_DIR` | Source documentation being migrated (`specs/`) |
| `LDS_DIST_DIR` | The Layered Documentation System distribution directory |

### Core System Terms

| Term | Definition |
|------|------------|
| **Layer** | A category of documentation organized by purpose (Decisions, Vision, Specifications, etc.) |
| **Canonical Source** | The authoritative location for a piece of information |
| **SSOT** | Single Source of Truth — each fact is mastered in exactly one place |
| **Cross-Link** | A reference from one document to another that provides related context |

### Document Structure Terms

| Term | Definition |
|------|------------|
| **Main Document** | The primary `{topic}.md` file in a layer directory |
| **Appendix** | Supporting content in `appendices/{topic}/` that contains detailed examples, code, etc. |
| **Appendix Directory** | The `appendices/` subdirectory within a layer, organized by topic |

### Layer-Specific Terms

| Term | Definition |
|------|------------|
| **ADR** | Architecture Decision Record — documents why a significant decision was made |
| **MADR** | Markdown Architecture Decision Record — a template format for ADRs |
| **PRD** | Product Requirements Document — documents product vision and goals |
| **Gherkin** | A language for writing executable specifications (Given/When/Then format) |

---

## Content Thresholds (Configurable)

These thresholds control how content is classified. Edit to customize for this project.

| Threshold | Value | Purpose |
|-----------|-------|---------|
| `CODE_BLOCK_LINES` | 50 | Code blocks >= this go to appendix |
| `STEP_LIST_ITEMS` | 10 | Step lists >= this go to appendix |
| `TABLE_ROWS` | 20 | Tables >= this go to appendix |
| `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Complete file examples -> appendix |
| `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Error catalogs -> appendix |
| `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Shell scripts -> appendix |

---

## Layer Overview

This documentation system uses 7 layers, organized by purpose:

| Layer | Directory | Purpose | Template |
|-------|-----------|---------|----------|
| 1. Decisions | `decisions/` | Capture WHY significant choices were made | MADR Minimal |
| 2. Vision | `product/` | Define WHAT we're building and WHY it matters | Lean PRD |
| 3. Architecture | `architecture/` | Show HOW components relate | C4-Lite |
| 4. Specifications | `specs/` | Detail HOW features work | Feature Spec |
| 5. Reference | `reference/` | Provide lookup tables | CLI + Config |
| 6. Behaviors | `behaviors/features/` | Define verifiable scenarios | Gherkin |
| 7. Implementation | `implementation/` | Guide HOW to build | Tech Stack |

---

## Classification Decision Tree

When adding new content, use this tree to determine the correct layer:

```
Is this explaining WHY a choice was made?
├─ YES → Layer 1: Decisions (ADR)
└─ NO ↓

Is this about product vision, goals, or concepts?
├─ YES → Layer 2: Vision
└─ NO ↓

Is this showing how components relate (diagrams, flows)?
├─ YES → Layer 3: Architecture
└─ NO ↓

Is this detailing HOW a feature works (behavior, edge cases)?
├─ YES → Layer 4: Specifications
└─ NO ↓

Is this a lookup table (commands, options, configs)?
├─ YES → Layer 5: Reference
└─ NO ↓

Is this a concrete verifiable scenario (Given/When/Then)?
├─ YES → Layer 6: Behaviors
└─ NO ↓

Is this implementation scaffolding (code templates, dependencies)?
├─ YES → Layer 7: Implementation
└─ NO → May not need documentation
```

---

## Classification Heuristics

### Layer 1: Decisions (ADRs)

**Signals to look for**:
- "We chose X over Y because..."
- "We decided to..." / "Decision:"
- "This pattern will be used throughout..."
- Trade-off discussions with rationale
- Explicit "Decision" blocks or sections

**ADR Structure**:
- Simple single files: `NNNN-{title}.md`
- No appendices for ADRs (include code inline)
- Status: Proposed | Accepted | Deprecated | Superseded

### Layer 2: Vision (PRD)

**Signals to look for**:
- Problem statements ("The problem we're solving...")
- Product vision or mission statements
- Success criteria or goals
- Non-goals or out-of-scope sections
- Core concepts and glossary

**Also in Layer 2**:
- `comparison.md` — How this compares to alternatives
- `roadmap.md` — Future considerations

### Layer 3: Architecture

**Signals to look for**:
- System diagrams (ASCII art, Mermaid)
- Network topology descriptions
- Component relationship descriptions
- Cross-cutting concerns

### Layer 4: Specifications

**Signals to look for**:
- Detailed behavioral descriptions
- State machines or lifecycle descriptions
- Data schemas with field-level detail
- Edge case documentation
- Configuration file format specifications

### Layer 5: Reference

**Signals to look for**:
- Command syntax documentation
- Option/flag tables
- Environment variable lists
- Configuration option tables with defaults
- Error code tables

### Layer 6: Behaviors

**Signals to look for**:
- "Given/When/Then" scenarios
- Explicit test case descriptions
- User journey examples with expected outcomes

### Layer 7: Implementation

**Signals to look for**:
- Technology stack with versions
- Dependency lists with rationale
- Project scaffolding instructions
- Code templates

---

## Cross-Layer Linking

Documents should link to related content in other layers:

| From | Link To | Purpose |
|------|---------|---------|
| Specifications | ADRs | Explain "why" for design choices |
| Specifications | Reference | Point to detailed syntax |
| Architecture | ADRs | Justify architectural patterns |
| Architecture | Specifications | Deep-dive into component behavior |
| Reference | Specifications | Provide conceptual context |
| Implementation | ADRs | Explain technology choices |
| Implementation | Specifications | Reference what's being implemented |

### Link Format

```markdown
## Related Documents

- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md) — Design rationale
- [Proxy Specification](../specs/proxy-infrastructure.md) — Detailed behavior

## Appendices

- [Traefik Configuration](./appendices/proxy-infrastructure/traefik-config.yaml) — Complete config
```

---

## Document Hierarchy (Authority Order)

When content appears in multiple places, this hierarchy determines the canonical source:

1. **ADRs** (highest authority) — Decisions are final
2. **Vision** — Product direction
3. **Architecture** — System design
4. **Specifications** — Feature behavior
5. **Reference** — Lookup data
6. **Behaviors** — Test scenarios
7. **Implementation** (lowest authority) — Build guidance

For conflicts, defer to the higher-authority document.

---

## Tooling (Tier 2)

This project uses Tier 2 tooling:

| Tool | Purpose | Setup |
|------|---------|-------|
| **markdownlint** | Markdown linting | `npm install --save-dev markdownlint-cli` |
| **Vale** | Prose linting | `brew install vale` |
| **Log4brains** | ADR browsing | `npm install -g log4brains` |
| **Structurizr** | Architecture diagrams | Docker or cloud |

### Running Linters

```bash
# Markdown linting
npx markdownlint-cli2 "docs/**/*.md"

# Prose linting (requires .vale.ini)
vale docs/
```

---

## Directory Structure

```
docs/
├── README.md                    # Documentation index
├── DOCUMENTATION-GUIDE.md       # This file
│
├── decisions/                   # Layer 1: ADRs
│   ├── README.md               # ADR index
│   ├── 0000-template.md        # Template
│   └── 0001-*.md ... 0011-*.md # ADR files
│
├── product/                     # Layer 2: Vision
│   ├── README.md
│   ├── vision.md
│   ├── comparison.md
│   └── roadmap.md
│
├── architecture/                # Layer 3: Architecture
│   ├── README.md
│   └── overview.md
│
├── specs/                       # Layer 4: Specifications
│   ├── README.md
│   ├── _template.md
│   ├── {feature}.md            # Main spec files
│   └── appendices/
│       └── {feature}/          # Per-spec appendices
│
├── reference/                   # Layer 5: Reference
│   ├── README.md
│   ├── cli.md
│   ├── configuration.md
│   └── appendices/
│       ├── cli/
│       └── configuration/
│
├── behaviors/                   # Layer 6: Behaviors
│   └── features/
│       └── *.feature
│
├── implementation/              # Layer 7: Implementation
│   ├── README.md
│   ├── tech-stack.md
│   └── appendices/
│       └── tech-stack/
│
├── maintenance/                 # Maintenance workflows
│   ├── audit.md
│   ├── refine.md
│   ├── sync.md
│   └── update.md
│
└── .migration/                  # Migration step files (temporary)
    ├── README.md
    ├── common-instructions.md
    └── 01-decisions.md ... 09-cross-links.md
```

---

## Appendix Guidelines

### When to Use Appendices

Move content to appendices when it exceeds thresholds:
- Code blocks >= 50 lines
- Step lists >= 10 items
- Tables >= 20 rows
- Complete file examples (always)
- Error catalogs (always)
- Shell scripts (always)

### Appendix Directory Structure

```
specs/
├── shell-integration.md          # Main document (overview, key concepts)
└── appendices/
    └── shell-integration/        # Named after main document
        ├── bash-setup.sh         # Complete bash script
        ├── zsh-setup.zsh         # Complete zsh script
        └── fish-setup.fish       # Complete fish script
```

### Linking to Appendices

From main document:
```markdown
For complete shell scripts, see:
- [Bash Setup](./appendices/shell-integration/bash-setup.sh)
- [Zsh Setup](./appendices/shell-integration/zsh-setup.zsh)
```

From appendix (back-link):
```markdown
> **Parent**: [Shell Integration](../../shell-integration.md)
```

---

## Maintenance Workflows

Four workflows are available in `maintenance/`:

| Workflow | Purpose | When to Use |
|----------|---------|-------------|
| `audit.md` | Verify documentation completeness | After migration, periodically |
| `refine.md` | Improve quality without code changes | Documentation reviews |
| `sync.md` | Verify docs match implementation | Pre-release, after refactoring |
| `update.md` | Update docs after code changes | When implementation changes |

To execute a workflow:
```
Execute the workflow in @docs/maintenance/audit.md
```
