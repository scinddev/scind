# Migration Step: Cross-Layer Links

**Prerequisites**:
- Read `common-instructions.md`
- **All previous steps (01-08) must be completed first**

---

## Overview

After all content is migrated, add cross-layer links to improve navigation. This step adds "See Also" sections to key documents.

---

## Links to Add

### `decisions/README.md`

Add to bottom:
```markdown
## Related Documentation

- [Architecture Overview](../architecture/overview.md) - How these decisions manifest in the system design
- [Product Vision](../product/vision.md) - The problem these decisions solve
```

### `product/vision.md`

Add to bottom:
```markdown
## Related Documentation

- [Architecture Overview](../architecture/overview.md) - System design implementing this vision
- [Architectural Decisions](../decisions/README.md) - Key technical decisions
- [Comparison](./comparison.md) - How Contrail compares to alternatives
```

### `product/comparison.md`

Add to bottom:
```markdown
## Related Documentation

- [Product Vision](./vision.md) - What Contrail aims to achieve
- [Roadmap](./roadmap.md) - Planned enhancements
```

### `architecture/overview.md`

Add to bottom:
```markdown
## Related Documentation

- [Product Vision](../product/vision.md) - The problem Contrail solves
- [Architectural Decisions](../decisions/README.md) - Rationale for design choices
- [Configuration Schemas](../specs/configuration-schemas.md) - Detailed configuration specification
- [Workspace Lifecycle](../specs/workspace-lifecycle.md) - How operations work
```

### `specs/configuration-schemas.md`

Add to bottom:
```markdown
## Related Documentation

- [Architecture Overview](../architecture/overview.md) - System context for configuration
- [CLI Reference](../reference/cli.md) - Commands that use these configurations
- [Configuration Reference](../reference/configuration.md) - User-facing configuration guide
```

### `specs/shell-integration.md`

Add to bottom:
```markdown
## Related Documentation

- [CLI Reference](../reference/cli.md) - Full CLI documentation
- [Implementation: Tech Stack](../implementation/tech-stack.md) - Shell script embedding details
```

### `reference/cli.md`

Add to bottom:
```markdown
## Related Documentation

- [Context Detection Spec](../specs/context-detection.md) - How context detection works
- [Shell Integration Spec](../specs/shell-integration.md) - Shell function details
- [Configuration Reference](./configuration.md) - Configuration file reference
```

### `reference/configuration.md`

Add to bottom:
```markdown
## Related Documentation

- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Detailed schema specification
- [CLI Reference](./cli.md) - Commands that use these configurations
```

### `implementation/tech-stack.md`

Add to bottom:
```markdown
## Related Documentation

- [CLI Reference](../reference/cli.md) - Complete CLI documentation
- [Shell Integration Spec](../specs/shell-integration.md) - Shell scripts to embed
- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Config types to implement
```

---

## Also Create: `docs/README.md`

The root documentation index:

```markdown
# Contrail Documentation

Welcome to the Contrail documentation. Contrail is a workspace orchestration system for Docker Compose.

## Quick Navigation

| Layer | Description | Entry Point |
|-------|-------------|-------------|
| **Decisions** | Architectural decision records | [decisions/](./decisions/README.md) |
| **Product** | Vision, comparison, roadmap | [product/](./product/README.md) |
| **Architecture** | System design overview | [architecture/](./architecture/README.md) |
| **Specifications** | Detailed feature specs | [specs/](./specs/README.md) |
| **Reference** | CLI and configuration reference | [reference/](./reference/README.md) |
| **Implementation** | Developer guides | [implementation/](./implementation/README.md) |

## Getting Started

1. Read the [Product Vision](./product/vision.md) to understand what Contrail does
2. Review the [Architecture Overview](./architecture/overview.md) for system design
3. See the [CLI Reference](./reference/cli.md) for command usage

## Documentation Guide

For contributors, see the [Documentation Guide](./DOCUMENTATION-GUIDE.md) for how to maintain and extend this documentation.
```

---

## Completion Checklist

- [ ] All "Related Documentation" sections added
- [ ] `docs/README.md` created
- [ ] All links verified to exist
- [ ] No broken relative paths
