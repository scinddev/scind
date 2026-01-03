# Layered Documentation System (LDS) — Create

**For AI Agents**: This document contains instructions for creating new documentation within an established Layered Documentation System. Follow this process when adding new features, decisions, or specifications.

**Terminology**: See the [Glossary](./LAYERED-DOCUMENTATION-SYSTEM.md#glossary) for definitions of key terms.

---

## Prerequisites

Before creating new documentation:

1. Read `LAYERED-DOCUMENTATION-SYSTEM.md` for system context
2. Read `DOCUMENTATION-GUIDE.md` in the docs directory for project-specific thresholds
3. Understand the existing documentation structure
4. Know what you're documenting (feature, decision, behavior, etc.)

---

## Creation Process

### Step 1: Classify the Content

Use the classification decision tree to determine which layer the new content belongs in:

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
└─ NO ↓

Is this implementation scaffolding, code templates, or dependency lists?
├─ YES → Layer 7: Implementation Guides
└─ NO → Reconsider: may not need documentation
```

Store the result as `TARGET_LAYER`.

---

### Step 2: Check for Conflicts

**Before writing any content**, verify the new documentation won't contradict existing docs.

#### 2a: Check Existing ADRs

Read all ADR files in `docs/decisions/`. For each ADR:
- Does the new content contradict any accepted decision?
- Does the new content imply a decision that should be an ADR?

If conflicts found:

> **Conflict Detected**
>
> The proposed documentation conflicts with existing ADR(s):
> - [ADR-NNNN: Title](path/to/adr.md) — {describe conflict}
>
> Options:
> 1. Revise the new documentation to align with the ADR
> 2. Create a new ADR to supersede the old decision
> 3. Abort and discuss with the team
>
> Which option? (1/2/3)

#### 2b: Check Existing Specifications

Read relevant specification files in `docs/specs/`. For each spec:
- Does the new content contradict existing specifications?
- Is there overlap that should be consolidated?

If conflicts found:

> **Specification Overlap Detected**
>
> The proposed documentation overlaps with:
> - [Specification Name](path/to/spec.md) — {describe overlap}
>
> Options:
> 1. Update the existing specification instead of creating new
> 2. Create new documentation with explicit cross-references
> 3. Merge content into a single specification
>
> Which option? (1/2/3)

#### 2c: Check Gherkin Behaviors

If the new content affects behavior, check `features/` for existing scenarios:
- Would implementing this break existing tests?
- Are there scenarios that need updating?

---

### Step 3: Select Template

Based on `TARGET_LAYER`, use the appropriate template:

| Layer | Template File | Priority |
|-------|---------------|----------|
| Decisions | `templates/adr-madr-minimal.md` | Default |
| Decisions | `templates/adr-madr-full.md` | Complex decisions |
| Decisions | `templates/adr-y-statement.md` | Quick capture |
| Vision | `templates/prd-lean.md` | Default |
| Vision | `templates/prd-epic-based.md` | Agile teams |
| Architecture | `templates/architecture-c4-lite.md` | Default |
| Architecture | `templates/architecture-arc42.md` | Comprehensive |
| Specifications | `templates/spec-feature.md` | Default |
| Specifications | `templates/spec-rfc.md` | Proposals needing review |
| Reference | `templates/reference-cli.md` | CLI documentation |
| Reference | `templates/reference-config.md` | Configuration docs |
| Behaviors | `templates/behavior-gherkin.feature` | Default |
| Implementation | `templates/implementation-tech-stack.md` | Default |

---

### Step 4: Determine File Location

Based on `TARGET_LAYER`, determine the file path using the flat file structure with shared `appendices/` directory:

| Layer | Directory Structure |
|-------|---------------------|
| Decisions | `docs/decisions/NNNN-{title}.md` (simple single files) |
| Vision | `docs/product/{topic}.md` |
| Architecture | `docs/architecture/{topic}.md` |
| Specifications | `docs/specs/{feature-name}.md` |
| Reference | `docs/reference/{topic}.md` |
| Behaviors | `features/{feature}.feature` |
| Implementation | `docs/implementation/{topic}.md` |

For ADRs, find the next available number:
```bash
ls docs/decisions/[0-9]*.md | sort | tail -1
# If last is 0005-*.md, next is 0006-*.md
```

**Directory structure**:
```
docs/{layer}/
├── README.md              # Auto-generated layer index
├── {topic}.md             # Main document
└── appendices/            # Only if needed (see Step 5b)
    └── {topic}/           # Matches main document basename
        └── {appendix-name}.md
```

---

### Step 5: Create the Document

#### 5a: Copy Template Structure

1. Create the file: `docs/{layer}/{topic}.md`
2. Use the appropriate template structure
3. Fill in:
   - Title/heading
   - Version (start at 0.1.0 for new docs)
   - Date
   - Status (Draft for new content)

#### 5b: Check Content Against Thresholds

Before writing, read the thresholds from `DOCUMENTATION-GUIDE.md`:

| Threshold | Default | If Exceeded |
|-----------|---------|-------------|
| `CODE_BLOCK_LINES` | 50 | Move to `appendices/` |
| `STEP_LIST_ITEMS` | 10 | Move to `appendices/` |
| `TABLE_ROWS` | 20 | Move to `appendices/` |
| `EXAMPLE_FILE_ALWAYS_APPENDIX` | true | Always in `appendices/` |
| `ERROR_CATALOG_ALWAYS_APPENDIX` | true | Always in `appendices/` |
| `SHELL_SCRIPT_ALWAYS_APPENDIX` | true | Always in `appendices/` |

**If your content will exceed thresholds**:
1. Create the `appendices/{topic}/` subdirectory (where `{topic}` matches your main document's basename)
2. Create appropriately named appendix files in that directory
3. Link from main document to appendices

**Appendix file structure**:
```markdown
# [Appendix Title]

> **Parent**: [Link to main document](../../{topic}.md)
> **Purpose**: [What this appendix contains]

## Content

[The detailed content]

---

*This appendix supports [{Topic} Document](../../{topic}.md).*
```

#### 5c: Write Content

Follow the template sections. Key guidelines:

**For ADRs**:
- Context should explain the situation clearly
- Decision should be a clear statement
- Consequences should include both positive and negative
- Status should be "Proposed" until reviewed

**For Specifications**:
- Include concrete examples
- Document edge cases
- Reference ADRs for rationale (don't duplicate "why")
- Be testable — specific enough to verify

**For Reference**:
- Be exhaustive — cover all options
- Use tables for scanability
- Include defaults and types
- Add examples for complex options

**For Behaviors**:
- Write from user perspective
- Use Given/When/Then format
- One scenario per behavior
- Include both positive and negative cases

---

### Step 6: Add Cross-Layer Links

Add links to related documents in other layers:

#### From Specifications → ADRs

```markdown
## Design Rationale

This specification implements the decisions from:
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [ADR-0007: Port Type System](../decisions/0007-port-type-system.md)
```

#### From Main Document → Appendix

```markdown
## Examples

For basic usage, see below. For complete workflow examples,
see [Detailed Examples](./appendices/workspace-lifecycle/detailed-examples.md).
```

#### From Specifications → Reference

```markdown
## Related Reference

For CLI command details, see [CLI Reference](../reference/cli.md#workspace-up).
```

#### From Behaviors → Specifications

```markdown
# This feature verifies behavior defined in:
# - docs/specs/workspace-lifecycle.md

Feature: Workspace Lifecycle
  ...
```

#### From Architecture → Specifications

```markdown
For detailed behavior of the proxy component, see:
- [Proxy Integration Spec](../specs/proxy-integration.md)
```

---

### Step 7: Validate the Document

Before finalizing, run through this checklist:

- [ ] **Correct layer**: Content matches the layer's purpose
- [ ] **No conflicts**: Doesn't contradict ADRs or existing specs
- [ ] **Template followed**: All required sections completed
- [ ] **File structure**: Using `{topic}.md` with `appendices/{topic}/` if needed
- [ ] **Thresholds respected**: Large content in appendices
- [ ] **Cross-links added**: References to related documents and appendices
- [ ] **No duplication**: Information not copied from other sources
- [ ] **Testable** (for specs): Specific enough to verify
- [ ] **ADR rationale moved**: "Why" explanations are in ADRs, not inline

---

### Step 8: Update Related Documents

After creating the new document, update related docs:

#### For New ADRs

If the decision affects existing specifications:
1. Update affected specs to reference the new ADR
2. Remove any now-redundant rationale from specs

#### For New Specifications

If the spec should be tested:
1. Create corresponding Gherkin feature files
2. Link behavior files to the specification

#### For New Reference Docs

If reference content was extracted from specifications:
1. Update specs to link to the reference doc
2. Remove duplicated content from specs

---

### Step 9: Final Report

Present a summary of what was created:

> **Documentation Created**
>
> **New Document**: `{file_path}`
> **Layer**: {layer name}
> **Template Used**: {template name}
>
> **Cross-Links Added**:
> - From: {file} → To: {new doc}
> - From: {new doc} → To: {related doc}
>
> **Related Documents Updated**:
> - {list of modified files}
>
> **Next Steps**:
> - Review the document for accuracy
> - If ADR: Move status from "Proposed" to "Accepted" after review
> - If Behavior: Implement step definitions

---

## Quick Reference

### Content Type → Layer Mapping

| If you're documenting... | Use Layer... |
|--------------------------|--------------|
| A choice between options | 1: Decisions (ADR) |
| Product goals or vision | 2: Vision |
| System diagrams or topology | 3: Architecture |
| How a feature works | 4: Specifications |
| Command syntax or config options | 5: Reference |
| Test scenarios | 6: Behaviors |
| Build setup or dependencies | 7: Implementation |

### ADR Numbering

ADRs are numbered sequentially: `0001`, `0002`, etc.
- Never reuse numbers
- Never renumber existing ADRs
- If an ADR is superseded, keep the old number and add status note

### When to Use Appendices

| Content Type | Always Appendix? | Threshold |
|--------------|------------------|-----------|
| Complete file examples | Yes | - |
| Error message catalogs | Yes | - |
| Shell scripts | Yes | - |
| Code blocks | No | ≥ 50 lines |
| Step lists | No | ≥ 10 items |
| Tables | No | ≥ 20 rows |

### Document Status Values

| Status | Meaning |
|--------|---------|
| Draft | Work in progress, not yet proposed |
| Proposed | Ready for review |
| Accepted | Approved and authoritative |
| Deprecated | No longer applicable |
| Superseded by NNNN | Replaced by another document |
