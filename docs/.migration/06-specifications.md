# Migration Step: Layer 4 - Specifications

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 10 spec files + appendices, approximately 2,500 lines

---

## Overview

Create specification files from `specs/contrail-technical-spec.md` and `specs/contrail-shell-integration.md`.

**Target Files**:
1. `specs/configuration-schemas.md` - Configuration file schemas
2. `specs/naming-conventions.md` - Naming patterns and conventions
3. `specs/port-types.md` - Port type system specification
4. `specs/proxy-infrastructure.md` - Traefik proxy configuration
5. `specs/context-detection.md` - Directory context detection algorithm
6. `specs/workspace-lifecycle.md` - Workspace operations (up/down/generate)
7. `specs/generated-override-files.md` - Override file generation
8. `specs/docker-labels.md` - Docker label conventions
9. `specs/environment-variables.md` - Environment variable injection
10. `specs/shell-integration.md` - Shell integration and contrail-compose

---

## Spec 1: `specs/configuration-schemas.md`

**Source**: `specs/contrail-technical-spec.md:216-712`

Create this file by reading and copying the entire Configuration Schemas section, including:
- Proxy Configuration (lines 233-263)
- Proxy Infrastructure (lines 264-337)
- Workspace Registry (lines 339-369)
- Global State (lines 371-442)
- Workspace Configuration (lines 443-476)
- Template Variables (lines 477-548)
- Application Configuration (lines 549-683)
- State Management (lines 685-707)
- Generated Manifest (lines 709-765)

**Include complete YAML examples inline** - these are reference material.

---

## Spec 2: `specs/naming-conventions.md`

**Source**: `specs/contrail-technical-spec.md:1247-1295`

Create with the full "Conventions and Best Practices" content including:
- Application requirements
- Service contract description
- Naming conventions table
- Collision warning

---

## Spec 3: `specs/port-types.md`

**Source**: `specs/contrail-technical-spec.md:96-135`

Create with the complete "Port Types and Proxying" section including:
- Type/protocol table
- Port type descriptions
- Protocol details
- Visibility documentation
- Private services

---

## Spec 4: `specs/proxy-infrastructure.md`

**Source**: `specs/contrail-technical-spec.md:264-337`

Create with appendix for complete configuration files:
- Main spec content in `specs/proxy-infrastructure.md`
- Create `specs/appendices/proxy-infrastructure/traefik-compose.yaml`
- Create `specs/appendices/proxy-infrastructure/traefik-config.yaml`

---

## Spec 5: `specs/context-detection.md`

**Source**: `specs/contrail-technical-spec.md:1053-1122`, `specs/contrail-cli-reference.md:30-147`

Combine context detection content from both sources into a single specification.

---

## Spec 6: `specs/workspace-lifecycle.md`

**Source**: `specs/contrail-technical-spec.md:1125-1245`

Create with the complete "Operations" section including:
- Startup sequence
- Staleness detection
- Generation logic
- Shutdown sequence
- Destroy sequence
- Viewing logs
- Listing status

---

## Spec 7: `specs/generated-override-files.md`

**Source**: `specs/contrail-technical-spec.md:768-873`

Create with complete override file example:
- Main spec with explanation
- Create `specs/appendices/generated-override-files/complete-override-example.yaml`

---

## Spec 8: `specs/docker-labels.md`

**Source**: `specs/contrail-technical-spec.md:875-954`

Create with complete label documentation including:
- Context labels
- Export labels
- Proxy container labels
- External tool integration examples

---

## Spec 9: `specs/environment-variables.md`

**Source**: `specs/contrail-technical-spec.md:955-1026`

Create with complete environment variable documentation including:
- Naming convention
- Variable generation rules
- Type/protocol table
- Usage examples in PHP and JavaScript

---

## Spec 10: `specs/shell-integration.md`

**Source**: `specs/contrail-shell-integration.md` (entire file ~844 lines)

This is the largest spec. Create with appendices for shell scripts:
- Main spec in `specs/shell-integration.md` (overview, architecture, compose-prefix, usage)
- Create `specs/appendices/shell-integration/bash-setup.sh` (lines 152-344)
- Create `specs/appendices/shell-integration/zsh-setup.zsh` (lines 350-555)
- Create `specs/appendices/shell-integration/fish-setup.fish` (lines 561-707)

---

## Also Create: `specs/README.md`

```markdown
# Specifications

Detailed specifications for Contrail's features and behaviors.

## Index

| Specification | Description |
|---------------|-------------|
| [configuration-schemas.md](./configuration-schemas.md) | Configuration file schemas (proxy, workspace, application) |
| [naming-conventions.md](./naming-conventions.md) | Naming patterns for hostnames, aliases, and projects |
| [port-types.md](./port-types.md) | Port type system (proxied vs assigned) |
| [proxy-infrastructure.md](./proxy-infrastructure.md) | Traefik proxy configuration |
| [context-detection.md](./context-detection.md) | Directory context detection algorithm |
| [workspace-lifecycle.md](./workspace-lifecycle.md) | Workspace operations (up/down/generate) |
| [generated-override-files.md](./generated-override-files.md) | Override file generation |
| [docker-labels.md](./docker-labels.md) | Docker label conventions |
| [environment-variables.md](./environment-variables.md) | Environment variable injection |
| [shell-integration.md](./shell-integration.md) | Shell integration and contrail-compose |

## Appendices

Large code examples and complete file contents are in `appendices/{spec-name}/`.
```

---

## Also Create: `specs/_template.md`

```markdown
# [Feature Name] Specification

**Status**: Draft | Review | Accepted

---

## Overview

[Brief description of the feature]

---

## Specification

[Detailed specification content]

---

## Examples

[Usage examples]

---

## Edge Cases

[Edge case handling]

<!-- Template for new specifications -->
```

---

## Appendix Directories to Create

```
specs/appendices/
  proxy-infrastructure/
    traefik-compose.yaml
    traefik-config.yaml
  generated-override-files/
    complete-override-example.yaml
  shell-integration/
    bash-setup.sh
    zsh-setup.zsh
    fish-setup.fish
  configuration-schemas/
    complete-examples.md
```

---

## Completion Checklist

- [ ] All 10 spec files created
- [ ] `specs/README.md` created
- [ ] `specs/_template.md` created
- [ ] All appendix directories and files created
- [ ] Complete code examples in appendices (not truncated)
- [ ] Source attributions present
