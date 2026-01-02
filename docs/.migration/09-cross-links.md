# Migration Step: Layer Cross-Links

**Prerequisites**: Complete all layer migration steps (01-08)
**Estimated Time**: Final pass

---

## Overview

This final step verifies and adds cross-references between documents across all layers. The layered documentation system works best when documents reference each other appropriately.

---

## Cross-Reference Patterns

### ADRs → Specs

Each ADR that describes a decision should link to the specification that implements it.

| ADR | Should Link To |
|-----|----------------|
| 0001-docker-compose-project-name-isolation | specs/naming-conventions.md |
| 0002-two-layer-networking | specs/proxy-infrastructure.md |
| 0003-pure-overlay-design | specs/generated-override-files.md |
| 0004-convention-based-naming | specs/naming-conventions.md |
| 0005-structure-vs-state-separation | specs/configuration-schemas.md |
| 0006-three-configuration-schemas | specs/configuration-schemas.md, reference/configuration.md |
| 0007-port-type-system | specs/port-types.md |
| 0008-traefik-reverse-proxy | specs/proxy-infrastructure.md |
| 0009-flexible-tls-configuration | specs/proxy-infrastructure.md |
| 0010-up-down-command-semantics | specs/workspace-lifecycle.md, reference/cli.md |
| 0011-options-based-targeting | specs/context-detection.md, reference/cli.md |

### Specs → ADRs

Each specification should reference the ADRs that motivated its design.

| Spec | Should Link To |
|------|----------------|
| configuration-schemas | ADR-0005, ADR-0006 |
| context-detection | ADR-0011 |
| docker-labels | ADR-0008 |
| environment-variables | ADR-0003 |
| generated-override-files | ADR-0003 |
| naming-conventions | ADR-0001, ADR-0004 |
| port-types | ADR-0007 |
| proxy-infrastructure | ADR-0002, ADR-0008, ADR-0009 |
| shell-integration | (no specific ADR) |
| workspace-lifecycle | ADR-0010 |

### Architecture → Everything

The architecture overview should serve as a hub, linking to:
- All relevant ADRs for design decisions
- All specs for detailed behavior
- Reference docs for quick lookup

### Reference → Specs

Reference documentation should link to specs for detailed behavior:

| Reference | Should Link To |
|-----------|----------------|
| cli.md | specs/context-detection.md, specs/workspace-lifecycle.md |
| configuration.md | specs/configuration-schemas.md |

### Vision → Architecture

The vision document should link down to architecture for "how" questions.

---

## Verification Checklist

### ADRs

- [ ] Each ADR has "Related Documents" section
- [ ] Links to implementing specs are present
- [ ] Links use relative paths (`../specs/...`)

### Specifications

- [ ] Each spec has "Related Documents" section
- [ ] Links to relevant ADRs are present
- [ ] Links to related specs are present
- [ ] Appendix links work if appendices exist

### Architecture

- [ ] Links to all relevant ADRs
- [ ] Links to key specs
- [ ] Diagram references are explained

### Reference

- [ ] Links to detailed specs
- [ ] Links to relevant ADRs for design rationale

### Product (Vision)

- [ ] Links to architecture
- [ ] Links to key ADRs

---

## Link Syntax

Use relative paths from the document location:

```markdown
<!-- From specs/port-types.md -->
[ADR-0007](../decisions/0007-port-type-system.md)

<!-- From decisions/0007-port-type-system.md -->
[Port Types Spec](../specs/port-types.md)

<!-- From architecture/overview.md -->
[Configuration Schemas](../specs/configuration-schemas.md)
[CLI Reference](../reference/cli.md)
```

---

## README Updates

Ensure all README.md index files are updated:

### `decisions/README.md`

```markdown
| ADR | Title | Status |
|-----|-------|--------|
| [0001](./0001-docker-compose-project-name-isolation.md) | Docker Compose Project Name Isolation | Accepted |
| [0002](./0002-two-layer-networking.md) | Two-Layer Networking | Accepted |
...
```

### `specs/README.md`

```markdown
| Specification | Description |
|---------------|-------------|
| [Configuration Schemas](./configuration-schemas.md) | workspace.yaml, application.yaml, proxy.yaml |
| [Context Detection](./context-detection.md) | Directory walking and context resolution |
...
```

### `product/README.md`

```markdown
| Document | Description |
|----------|-------------|
| [Vision](./vision.md) | Product vision and core concepts |
| [Comparison](./comparison.md) | Comparison with alternative tools |
| [Roadmap](./roadmap.md) | Future considerations |
```

---

## Final Verification

After completing cross-links:

1. **Link Check**: Run a markdown link checker to find broken links
   ```bash
   # If using markdownlint-cli2 (Tier 2 tooling)
   npx markdownlint-cli2 "docs/**/*.md"
   ```

2. **Visual Review**: Open each README and verify the index tables are complete

3. **Navigation Test**: Starting from `docs/README.md`, verify you can navigate to any document within 2-3 clicks

---

## Completion Checklist

- [ ] All ADRs have Related Documents sections
- [ ] All Specs have Related Documents sections
- [ ] Architecture links to ADRs and Specs
- [ ] Reference links to Specs
- [ ] All README index files updated
- [ ] Link check passes
- [ ] Navigation test passes

