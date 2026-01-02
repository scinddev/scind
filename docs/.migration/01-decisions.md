# Migration Step: Layer 1 — Decisions (ADRs)

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 11 ADRs, approximately 550 lines total

---

## Overview

Extract decision blocks from the source documents and create individual ADR files. ADRs are single files (not directories) using the MADR Minimal template.

**Source documents containing decisions**:
- `specs/contrail-prd.md` lines 161-315 (Key Architectural Decisions section)
- `specs/contrail-shell-integration.md` lines 818-825 (Alternative Approaches Considered)
- `specs/contrail-go-stack.md` lines 59-85 (Dependency Rationale, Intentionally Excluded)

---

## ADR 1: `decisions/0001-docker-compose-project-name-isolation.md`

**Source**: `specs/contrail-prd.md:162-169`

### Content

```markdown
# ADR-0001: Docker Compose Project Name Isolation

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Contrail needs to run multiple instances of the same application simultaneously. Each instance must have isolated containers, networks, and volumes.

## Decision

Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

## Consequences

### Positive

- Uses Docker's official mechanism for running multiple copies of the same stack
- Isolates containers, networks, and volumes without requiring modifications to the application
- Works with any existing Docker Compose project

### Negative

- Project names must be unique across all workspaces on a host
- Creative naming that produces identical project names could cause conflicts

### Neutral

- Naming follows convention: `{workspace}-{app}` (e.g., `dev-app-one`)

---

## Notes

Example collision to avoid: workspace `dev-app` with app `one` and workspace `dev` with app `app-one` both produce project name `dev-app-one`.

<!-- Migrated from specs/contrail-prd.md:162-169 -->
```

---

## ADR 2: `decisions/0002-two-layer-networking.md`

**Source**: `specs/contrail-prd.md:171-179`

### Content

```markdown
# ADR-0002: Two-Layer Networking

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Services need both external access (via reverse proxy) and internal access (between applications in the same workspace).

## Decision

Implement two network layers:
- `contrail-proxy` network: Host-wide, connects Traefik to public services
- `{workspace}-internal` network: Per-workspace, connects all applications for internal communication

## Consequences

### Positive

- Separating concerns allows public services to be routable via Traefik while protected services remain internal
- The workspace-internal network provides isolation between workspaces
- Applications in different workspaces cannot accidentally communicate

### Negative

- More complex network topology than a single flat network
- Debugging network issues requires understanding both layers

### Neutral

- Each workspace has its own internal network named `{workspace}-internal`

---

<!-- Migrated from specs/contrail-prd.md:171-179 -->
```

---

## ADR 3: `decisions/0003-pure-overlay-design.md`

**Source**: `specs/contrail-prd.md:181-192`

### Content

```markdown
# ADR-0003: Pure Overlay Design (Applications Remain Workspace-Agnostic)

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Applications could embed workspace configuration (labels, network definitions, environment variables), or this integration could be applied externally.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

## Consequences

### Positive

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems

### Negative

- Requires a generation step before running
- Generated files must be kept in sync with source configuration

### Neutral

- Override files are stored in `.generated/` and gitignored

---

<!-- Migrated from specs/contrail-prd.md:181-192 -->
```

---

## ADR 4: `decisions/0004-convention-based-naming.md`

**Source**: `specs/contrail-prd.md:194-202`

### Content

```markdown
# ADR-0004: Convention-Based Naming

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Hostnames, network aliases, and project names could be explicitly configured per-service, or derived from conventions.

## Decision

Derive names from conventions:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Protected alias: `{app}-{service}`
- Network name: `{workspace}-internal`
- Project name: `{workspace}-{app}`

## Consequences

### Positive

- Reduces configuration burden
- Ensures consistency across workspaces
- Makes the system predictable and debuggable
- Given workspace and app names, hostnames are deterministic

### Negative

- Less flexibility for unusual naming requirements
- Explicit overrides were considered but removed to keep the schema simple

### Neutral

- Templates can be customized at the workspace level for advanced use cases

---

<!-- Migrated from specs/contrail-prd.md:194-202 -->
```

---

## ADR 5: `decisions/0005-structure-vs-state-separation.md`

**Source**: `specs/contrail-prd.md:204-222`

### Content

```markdown
# ADR-0005: Structure vs State Separation

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Configuration could include runtime choices (which branch to use, which flavor is active) or only structural definitions.

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

### Positive

- Configuration files describe the system's shape, not its current state
- State changes frequently; structure changes rarely
- Avoids polluting config files with transient information
- Branch management stays with git where it belongs

### Negative

- State must be tracked separately (in `.generated/state.yaml`)
- Users must understand the distinction

### Neutral

- State file is gitignored to avoid conflicts

---

<!-- Migrated from specs/contrail-prd.md:204-222 -->
```

---

## ADR 6: `decisions/0006-three-configuration-schemas.md`

**Source**: `specs/contrail-prd.md:224-233`

### Content

```markdown
# ADR-0006: Three Configuration Schemas

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Configuration could be in one monolithic file or separated by concern and ownership.

## Decision

Three schema types that can be combined:
- `proxy`: Global/per-user (domain, Traefik settings)
- `workspace`: Per-workspace (name, application list)
- `application`: Per-application (flavors, exported services)

## Consequences

### Positive

- Separation of concerns—proxy config rarely changes, workspace config defines the environment, application config is owned by the application team
- Application config (`application.yaml`) can live in the application's own repository
- Changes to one layer don't require touching others

### Negative

- Multiple files to understand and maintain
- Configuration hierarchy must be documented

### Neutral

- Proxy config: `~/.config/contrail/proxy.yaml`
- Workspace config: `{workspace}/workspace.yaml`
- Application config: `{app}/application.yaml`

---

<!-- Migrated from specs/contrail-prd.md:224-233 -->
```

---

## ADR 7: `decisions/0007-port-type-system.md`

**Source**: `specs/contrail-prd.md:235-264`

### Content

```markdown
# ADR-0007: Port Type System for Exported Services

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Services need different handling based on how they're accessed—some need HTTP proxying through Traefik, others need direct port binding for database connections.

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

**Port types**:
- `proxied`: Traffic routed through Traefik; protocol specifies how (http, https, future SNI types)
- `assigned`: Direct port binding; if port unavailable, Contrail finds next available

## Consequences

### Positive

- `type` determines routing: `proxied` (through Traefik) or `assigned` (direct port binding)
- `protocol` specifies how proxied traffic is handled: `http`, `https`, or future SNI types
- Supports multiple protocols on the same exported service (both HTTP and HTTPS)
- Environment variables use proxy values (port 80/443) for proxied services
- Enables future plugin system for additional protocols (postgresql, mysql SNI routing)

### Negative

- More complex configuration than a single port type
- Users must understand the distinction between type and protocol

### Neutral

- `visibility` remains as documentation for collaborators (public vs protected)

---

<!-- Migrated from specs/contrail-prd.md:235-264 -->
```

---

## ADR 8: `decisions/0008-traefik-reverse-proxy.md`

**Source**: `specs/contrail-prd.md:266-272`

### Content

```markdown
# ADR-0008: Traefik for Reverse Proxy

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Contrail needs a reverse proxy that can dynamically route to containers based on hostname, without requiring config file changes when containers are started or stopped.

## Decision

Use Traefik with Docker provider, reading labels from containers.

## Consequences

### Positive

- Traefik's Docker integration allows dynamic routing without config file changes
- Labels on containers (added via generated overrides) define routing rules
- Traefik automatically discovers containers and their routing configuration
- Industry-standard tool with good documentation

### Negative

- Adds a dependency on Traefik
- Traefik must be running for proxied services to be accessible
- Users unfamiliar with Traefik may need to learn its concepts

### Neutral

- Traefik dashboard available for debugging routing issues

---

<!-- Migrated from specs/contrail-prd.md:266-272 -->
```

---

## ADR 9: `decisions/0009-flexible-tls-configuration.md`

**Source**: `specs/contrail-prd.md:274-289`

### Content

```markdown
# ADR-0009: Flexible TLS Configuration

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints—personal dev machines may use mkcert, while enterprise networks may have managed CAs.

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development—uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments—user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

## Consequences

### Positive

- `auto` provides zero-config HTTPS for most users with mkcert installed
- `custom` supports enterprise environments where developers already have CA-signed certs
- Avoids mandating a specific certificate tool while still enabling secure-by-default development

### Negative

- Multiple modes add complexity to documentation
- Users must understand which mode fits their environment

### Neutral

- Default mode is `auto` for simplicity

---

<!-- Migrated from specs/contrail-prd.md:274-289 -->
```

---

## ADR 10: `decisions/0010-up-down-command-semantics.md`

**Source**: `specs/contrail-prd.md:291-303`

### Content

```markdown
# ADR-0010: up/down Command Semantics

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Commands could use `start`/`stop` or `up`/`down` terminology. The naming affects user expectations about what the commands do.

## Decision

Use `up` and `down` as primary commands, matching Docker Compose semantics:
- `up`: Build, create networks/volumes, generate overrides, start containers
- `down`: Stop containers, remove containers/networks, optionally remove volumes

## Consequences

### Positive

- Semantic alignment with Docker Compose, which users already know
- `up` conveys "bring the environment into existence" (more than just starting)
- `down` conveys "tear down" rather than just pausing
- Matches the underlying `docker compose up/down` commands Contrail invokes

### Negative

- Users expecting `start`/`stop` may be initially confused
- `down` removes containers by default (unlike `stop` which just pauses)

### Neutral

- Aliases (`start`/`stop`) could be added later if needed

---

<!-- Migrated from specs/contrail-prd.md:291-303 -->
```

---

## ADR 11: `decisions/0011-options-based-targeting.md`

**Source**: `specs/contrail-prd.md:305-315`

### Content

```markdown
# ADR-0011: Options-Based Targeting with Context Detection

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Commands need to target specific workspaces and applications. This could be done via positional arguments or named options.

## Decision

Use `--workspace` and `--app` options (not positional arguments) with automatic context detection from current directory.

## Consequences

### Positive

- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options

### Negative

- Slightly more verbose than positional arguments
- Context detection behavior must be documented and understood

### Neutral

- Global flags are always available: `-w/--workspace`, `-a/--app`

---

<!-- Migrated from specs/contrail-prd.md:305-315 -->
```

---

## Completion Checklist

- [ ] All 11 ADR files created
- [ ] Source attributions present in each file
- [ ] Template structure followed (MADR Minimal)
- [ ] Update `decisions/README.md` index with new entries
