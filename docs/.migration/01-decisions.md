# Migration Step: Layer 1 - Decisions (ADRs)

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 11 ADRs, approximately 450 lines total

---

## Overview

Extract 11 architectural decisions from `specs/contrail-prd.md` (lines 162-315) into individual ADR files using the MADR Minimal template.

**Source Section**: "Key Architectural Decisions" in contrail-prd.md

---

## ADR 1: `decisions/0001-docker-compose-project-name-isolation.md`

**Source**: `specs/contrail-prd.md:162-169`

### Content

```markdown
# ADR-0001: Docker Compose Project Name Isolation

## Status

Accepted

## Context

Need to run multiple instances of the same application simultaneously.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

## Consequences

- This is Docker's official mechanism for running multiple copies of the same stack
- It isolates containers, networks, and volumes without requiring modifications to the application
- Applications can run independently with different project names

<!-- Migrated from specs/contrail-prd.md:162-169 -->
```

---

## ADR 2: `decisions/0002-two-layer-networking.md`

**Source**: `specs/contrail-prd.md:170-179`

### Content

```markdown
# ADR-0002: Two-Layer Networking

## Status

Accepted

## Context

Services need both external access (via reverse proxy) and internal access (between applications).

## Decision

Implement two network layers:
- `contrail-proxy` network: Host-wide, connects Traefik to public services
- `{workspace}-internal` network: Per-workspace, connects all applications for internal communication

## Consequences

- Separating concerns allows public services to be routable via Traefik while protected services remain internal
- The workspace-internal network provides isolation between workspaces
- Applications in different workspaces cannot communicate directly

<!-- Migrated from specs/contrail-prd.md:170-179 -->
```

---

## ADR 3: `decisions/0003-pure-overlay-design.md`

**Source**: `specs/contrail-prd.md:180-191`

### Content

```markdown
# ADR-0003: Pure Overlay Design (Applications Remain Workspace-Agnostic)

## Status

Accepted

## Context

Applications could embed workspace configuration, or it could be applied externally.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

## Consequences

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems
- Override files must be regenerated when configuration changes

<!-- Migrated from specs/contrail-prd.md:180-191 -->
```

---

## ADR 4: `decisions/0004-convention-based-naming.md`

**Source**: `specs/contrail-prd.md:192-202`

### Content

```markdown
# ADR-0004: Convention-Based Naming

## Status

Accepted

## Context

Hostnames and aliases could be explicitly configured or derived from conventions.

## Decision

Derive names from conventions:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Protected alias: `{app}-{service}`
- Network name: `{workspace}-internal`

## Consequences

- Conventions reduce configuration and ensure consistency
- System behavior is predictable - given workspace and app names, hostnames are deterministic
- Explicit overrides were considered but removed to keep the schema simple
- Users must follow naming conventions for names to avoid collisions

<!-- Migrated from specs/contrail-prd.md:192-202 -->
```

---

## ADR 5: `decisions/0005-structure-vs-state-separation.md`

**Source**: `specs/contrail-prd.md:203-222`

### Content

```markdown
# ADR-0005: Structure vs State Separation

## Status

Accepted

## Context

Configuration could include runtime choices (which branch, which flavor) or only structural definitions.

## Decision

Separate structure (what exists) from state (what's active):

| Aspect | Structure (config files) | State (runtime) |
|--------|--------------------------|-----------------|
| What apps exist | workspace.yaml | - |
| Available flavors | application.yaml | - |
| Active flavor | - | .generated/state.yaml or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

## Consequences

- Configuration files describe the system's shape, not its current state
- State changes frequently; structure changes rarely
- Avoids polluting config files with transient information
- Branch management stays with git where it belongs

<!-- Migrated from specs/contrail-prd.md:203-222 -->
```

---

## ADR 6: `decisions/0006-three-configuration-schemas.md`

**Source**: `specs/contrail-prd.md:223-233`

### Content

```markdown
# ADR-0006: Three Configuration Schemas

## Status

Accepted

## Context

Configuration could be in one monolithic file or separated by concern.

## Decision

Three schema types that can be combined:
- `proxy`: Global/per-user (domain, Traefik settings)
- `workspace`: Per-workspace (name, application list)
- `application`: Per-application (flavors, exported services)

## Consequences

- Separation of concerns - proxy config rarely changes, workspace config defines the environment, application config is owned by the application team
- Configuration is distributed across files, requiring multiple reads
- Clear ownership boundaries between teams

<!-- Migrated from specs/contrail-prd.md:223-233 -->
```

---

## ADR 7: `decisions/0007-port-type-system.md`

**Source**: `specs/contrail-prd.md:234-264`

### Content

```markdown
# ADR-0007: Port Type System for Exported Services

## Status

Accepted

## Context

Services need different handling based on how they're accessed - some need HTTP proxying, others need direct port binding.

## Decision

Each exported service declares ports with a `type` (routing mechanism) and optionally a `protocol`:

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
      - type: proxied
        protocol: http
        visibility: protected
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

## Consequences

- `type` determines routing: `proxied` (through Traefik) or `assigned` (direct port binding)
- `protocol` specifies how proxied traffic is handled: `http`, `https`, or future SNI types
- Supports multiple protocols on the same exported service (both HTTP and HTTPS)
- Environment variables use proxy values (port 80/443) for proxied services
- Enables future plugin system for additional protocols (postgresql, mysql SNI routing)
- `visibility` remains as documentation for collaborators

<!-- Migrated from specs/contrail-prd.md:234-264 -->
```

---

## ADR 8: `decisions/0008-traefik-reverse-proxy.md`

**Source**: `specs/contrail-prd.md:265-271`

### Content

```markdown
# ADR-0008: Traefik for Reverse Proxy

## Status

Accepted

## Context

Need a reverse proxy that can dynamically route to containers.

## Decision

Use Traefik with Docker provider, reading labels from containers.

## Consequences

- Traefik's Docker integration allows dynamic routing without config file changes
- Labels on containers (added via generated overrides) define routing rules
- Single shared proxy instance serves all workspaces
- Requires Traefik knowledge for advanced customization

<!-- Migrated from specs/contrail-prd.md:265-271 -->
```

---

## ADR 9: `decisions/0009-flexible-tls-configuration.md`

**Source**: `specs/contrail-prd.md:272-289`

### Content

```markdown
# ADR-0009: Flexible TLS Configuration

## Status

Accepted

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints (personal dev machines, enterprise networks with managed CAs).

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development - uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments - user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

## Consequences

- `auto` provides zero-config HTTPS for most users with mkcert installed
- `custom` supports enterprise environments where developers already have CA-signed certs
- Avoids mandating a specific certificate tool while still enabling secure-by-default development

<!-- Migrated from specs/contrail-prd.md:272-289 -->
```

---

## ADR 10: `decisions/0010-up-down-command-semantics.md`

**Source**: `specs/contrail-prd.md:290-303`

### Content

```markdown
# ADR-0010: up/down Command Semantics

## Status

Accepted

## Context

Commands could use `start`/`stop` or `up`/`down` terminology.

## Decision

Use `up` and `down` as primary commands, matching Docker Compose semantics:
- `up`: Build, create networks/volumes, generate overrides, start containers
- `down`: Stop containers, remove containers/networks, optionally remove volumes

## Consequences

- Semantic alignment with Docker Compose, which users already know
- `up` conveys "bring the environment into existence" (more than just starting)
- `down` conveys "tear down" rather than just pausing
- Matches the underlying `docker compose up/down` commands Contrail invokes

<!-- Migrated from specs/contrail-prd.md:290-303 -->
```

---

## ADR 11: `decisions/0011-options-based-targeting.md`

**Source**: `specs/contrail-prd.md:304-315`

### Content

```markdown
# ADR-0011: Options-Based Targeting with Context Detection

## Status

Accepted

## Context

Commands need to target specific workspaces and applications.

## Decision

Use `--workspace` and `--app` options (not positional arguments) with automatic context detection from current directory.

## Consequences

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options

<!-- Migrated from specs/contrail-prd.md:304-315 -->
```

---

## Also Create: `decisions/README.md`

```markdown
# Architectural Decision Records

This directory contains Architecture Decision Records (ADRs) documenting significant technical and product decisions for Contrail.

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [0001](./0001-docker-compose-project-name-isolation.md) | Docker Compose Project Name Isolation | Accepted |
| [0002](./0002-two-layer-networking.md) | Two-Layer Networking | Accepted |
| [0003](./0003-pure-overlay-design.md) | Pure Overlay Design | Accepted |
| [0004](./0004-convention-based-naming.md) | Convention-Based Naming | Accepted |
| [0005](./0005-structure-vs-state-separation.md) | Structure vs State Separation | Accepted |
| [0006](./0006-three-configuration-schemas.md) | Three Configuration Schemas | Accepted |
| [0007](./0007-port-type-system.md) | Port Type System | Accepted |
| [0008](./0008-traefik-reverse-proxy.md) | Traefik for Reverse Proxy | Accepted |
| [0009](./0009-flexible-tls-configuration.md) | Flexible TLS Configuration | Accepted |
| [0010](./0010-up-down-command-semantics.md) | up/down Command Semantics | Accepted |
| [0011](./0011-options-based-targeting.md) | Options-Based Targeting | Accepted |

## Creating New ADRs

Use `0000-template.md` as a starting point for new decisions.
```

---

## Also Create: `decisions/0000-template.md`

```markdown
# ADR-NNNN: [Title]

## Status

[Proposed | Accepted | Deprecated | Superseded]

## Context

[Describe the context and problem statement]

## Decision

[Describe the decision and its justification]

## Consequences

[Describe the consequences of this decision]

<!-- Template for new ADRs -->
```

---

## Completion Checklist

- [ ] All 11 ADR files created
- [ ] README.md created with index
- [ ] 0000-template.md created
- [ ] Source attributions present in each file
