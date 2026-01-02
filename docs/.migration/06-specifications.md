# Migration Step: Layer 4 — Specifications

**Prerequisites**: Read `common-instructions.md`, complete Layer 3 steps
**Estimated Size**: 10 files, approximately 2,500 lines total

---

## Overview

Extract detailed specifications from source documents. Each spec follows the template in `specs/_template.md`. Large content (code blocks >50 lines, tables >20 rows) goes to appendices.

**Source documents**:
- `specs/contrail-technical-spec.md` (primary)
- `specs/contrail-cli-reference.md` (context detection, port handling)
- `specs/contrail-shell-integration.md` (shell integration spec)

---

## Spec 1: `specs/configuration-schemas.md`

**Source**: `specs/contrail-technical-spec.md:200-400` (Configuration section)

### Summary

Documents the three configuration file schemas: `proxy.yaml`, `workspace.yaml`, `application.yaml`.

### Content Structure

```markdown
# Configuration Schemas Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Three configuration schemas define Contrail's behavior...

## proxy.yaml Schema
[Full schema with field table]

## workspace.yaml Schema
[Full schema with field table]

## application.yaml Schema
[Full schema with field table]

## Validation Rules
[Combined validation rules]

## Examples
[Minimal and complete examples for each]

<!-- See appendices/configuration-schemas/ for complete schema files -->
```

### Appendix Content

Create `specs/appendices/configuration-schemas/`:
- `proxy-schema.yaml` — JSON Schema for proxy.yaml
- `workspace-schema.yaml` — JSON Schema for workspace.yaml
- `application-schema.yaml` — JSON Schema for application.yaml
- `complete-examples.md` — Full working examples

---

## Spec 2: `specs/context-detection.md`

**Source**: `specs/contrail-cli-reference.md:100-200` (Context Detection section)

### Summary

Documents how Contrail detects workspace and application from current directory.

### Content Structure

```markdown
# Context Detection Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Context detection allows commands to infer workspace/app from directory...

## Detection Algorithm
1. Start from current directory
2. Walk up looking for markers
3. Return context or empty

## Marker Files
- workspace.yaml → workspace root
- application.yaml → application root
- docker-compose.yaml → application root (fallback)

## Examples
[Directory structure examples with detected context]

## Edge Cases
[Nested workspaces, missing files, etc.]
```

---

## Spec 3: `specs/docker-labels.md`

**Source**: `specs/contrail-technical-spec.md:500-600` (Labels section)

### Summary

Documents the Traefik labels and workspace labels added to containers.

### Content Structure

```markdown
# Docker Labels Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Contrail adds labels to containers for routing and identification...

## Traefik Labels
[Label patterns for HTTP, HTTPS, routers, services]

## Workspace Labels
[contrail.workspace, contrail.app, contrail.service labels]

## Label Generation
[How labels are generated from configuration]

## Examples
[Complete label sets for different scenarios]
```

---

## Spec 4: `specs/environment-variables.md`

**Source**: `specs/contrail-technical-spec.md:600-700` (Environment section)

### Summary

Documents environment variable injection for service discovery.

### Content Structure

```markdown
# Environment Variables Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Contrail injects environment variables for service discovery...

## Variable Naming Convention
{APP}_{SERVICE}_{PROTOCOL}_URL
{APP}_{SERVICE}_HOST
{APP}_{SERVICE}_PORT

## Variable Values
[How values are determined based on port type]

## Override Behavior
[How explicit env vars in compose override generated ones]

## Examples
[Variable sets for different service configurations]
```

---

## Spec 5: `specs/generated-override-files.md`

**Source**: `specs/contrail-technical-spec.md:400-500` (Override Generation section)

### Summary

Documents the structure and generation of Docker Compose override files.

### Content Structure

```markdown
# Generated Override Files Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Contrail generates override files to integrate apps into workspaces...

## File Location
.generated/contrail-override.yaml

## File Structure
[YAML structure with networks, services, labels, environment]

## Generation Process
[Step-by-step generation algorithm]

## Merge Behavior
[How overrides merge with base compose file]

## Examples
[Before/after showing base + override = effective config]
```

### Appendix Content

Create `specs/appendices/generated-override-files/`:
- `complete-override-example.yaml` — Full generated override file

---

## Spec 6: `specs/naming-conventions.md`

**Source**: `specs/contrail-technical-spec.md:300-350` (Naming section)

### Summary

Documents all naming conventions for projects, hostnames, aliases, and networks.

### Content Structure

```markdown
# Naming Conventions Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Contrail uses convention-based naming for predictability...

## Project Names
Pattern: {workspace}-{app}
Example: main-frontend → project name "main-frontend"

## Hostnames (Public)
Pattern: {workspace}-{app}-{service}.{domain}
Example: main-frontend-web.test

## Network Aliases (Protected)
Pattern: {app}-{service}
Example: frontend-web

## Network Names
- Proxy: contrail-proxy
- Internal: {workspace}-internal

## Validation Rules
[Character restrictions, length limits]
```

---

## Spec 7: `specs/port-types.md`

**Source**: `specs/contrail-technical-spec.md:700-800` (Port Types section)

### Summary

Documents the port type system: proxied vs assigned.

### Content Structure

```markdown
# Port Types Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Each exported service port has a type determining how it's exposed...

## Type: proxied
- Traffic routes through Traefik
- Protocols: http, https
- Uses standard ports (80, 443)

## Type: assigned
- Direct host port binding
- Port assignment algorithm (preferred → next available)
- Persistence across restarts

## Visibility
- public: Accessible via hostname
- protected: Internal network only

## Examples
[Configuration and resulting behavior]
```

---

## Spec 8: `specs/proxy-infrastructure.md`

**Source**: `specs/contrail-technical-spec.md:800-900` (Proxy section)

### Summary

Documents Traefik configuration and the contrail-proxy network.

### Content Structure

```markdown
# Proxy Infrastructure Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Traefik provides reverse proxy capabilities...

## Traefik Configuration
[Static and dynamic config]

## contrail-proxy Network
[Network creation, attachment, purpose]

## TLS Handling
[Certificate modes: auto, custom, disabled]

## Health Checks
[Traefik health and container health]

## Examples
[Traefik docker-compose.yaml]
```

### Appendix Content

Create `specs/appendices/proxy-infrastructure/`:
- `traefik-compose.yaml` — Complete Traefik docker-compose file
- `traefik-config.yaml` — Traefik static configuration

---

## Spec 9: `specs/shell-integration.md`

**Source**: `specs/contrail-shell-integration.md:1-600` (Shell Integration sections)

### Summary

Documents the `contrail-compose` shell function and shell integration.

### Content Structure

```markdown
# Shell Integration Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Shell integration provides context-aware docker compose commands...

## contrail-compose Function
[Function behavior, when it activates, pass-through mode]

## Shell Setup
[bash, zsh, fish installation]

## Completion
[Tab completion for commands and options]

## Environment Variables
[CONTRAIL_WORKSPACE, CONTRAIL_APP, etc.]

## Examples
[Usage examples in different shells]
```

### Appendix Content

Create `specs/appendices/shell-integration/`:
- `bash-setup.sh` — Bash integration script
- `zsh-setup.zsh` — Zsh integration script
- `fish-setup.fish` — Fish integration script

---

## Spec 10: `specs/workspace-lifecycle.md`

**Source**: `specs/contrail-technical-spec.md:900-1000` (Lifecycle section)

### Summary

Documents workspace and application lifecycle: up, down, restart semantics.

### Content Structure

```markdown
# Workspace Lifecycle Specification

**Version**: 1.0.0
**Status**: Accepted

## Overview
Lifecycle commands manage workspace and application state...

## up Command
[What happens: generate, network, start]

## down Command
[What happens: stop, remove containers, optionally volumes]

## restart Command
[Behavior: down + up vs reload]

## Status Tracking
[How Contrail tracks what's running]

## Partial Operations
[Starting/stopping individual apps]

## Examples
[Command sequences and resulting states]
```

---

## Completion Checklist

- [ ] All 10 spec files created
- [ ] Appendix directories created where needed
- [ ] Large content moved to appendices
- [ ] Cross-references to ADRs added
- [ ] Update `specs/README.md` to list all specifications

