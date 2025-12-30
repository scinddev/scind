# Contrail: Product Requirements Document

**Version**: 0.5.0-draft  
**Date**: December 2024  
**Status**: Design Phase

---

## Executive Summary

Contrail is a workspace orchestration system for Docker Compose that enables developers to run multiple isolated instances of multi-application stacks simultaneously on a single host. It solves the problem of needing complete, independent environments for development, code review, and testing without requiring Kubernetes or cloud infrastructure.

The name "Contrail" evokes the trails left by aircraft—parallel paths that don't intersect, much like the isolated workspaces the system creates.

---

## Problem Statement

### The Scenario

A developer works on a system composed of multiple applications (e.g., `app-one`, `app-two`, `app-three`) that communicate with each other. They need to run multiple complete copies of this stack simultaneously:

- **dev**: Active development with local changes
- **review**: Clean checkout of a PR branch for code review  
- **control**: Stable baseline for comparison

Each environment needs:
1. Internal communication between applications (app-one can reach app-two's API)
2. External access via unique hostnames (dev-app-one-web.contrail.test, review-app-one-web.contrail.test)
3. Complete isolation (dev's database is separate from review's database)

### Why Existing Solutions Fall Short

| Solution | Limitation |
|----------|------------|
| Docker Compose alone | No built-in multi-instance orchestration; manual project naming |
| Docker Compose `include` | Merges into single application model; doesn't handle parallel instances |
| DDEV / Lando / Docksal | Single-application focused (one Drupal site, not multi-app stacks) |
| Skaffold / Tilt / Garden | Kubernetes-focused, not Docker Compose |
| Manual scripts | Error-prone, hard to maintain, no conventions |

### The Gap

There is no existing tool that orchestrates **multiple isolated instances of multi-application Docker Compose stacks** with **generated integration configuration** while keeping **applications workspace-agnostic**.

---

## Product Vision

Contrail provides a thin coordination layer over Docker Compose that:

1. **Preserves application independence**: Applications don't know they're in a workspace
2. **Uses pure overlay**: All integration happens via generated Docker Compose override files
3. **Follows conventions**: Predictable naming for hostnames, aliases, and networks
4. **Separates structure from state**: Configuration describes what exists; runtime state describes what's active
5. **Enables direct Docker Compose access**: The `contrail-compose` shell function provides context-aware passthrough to Docker Compose with full tab completion

---

## Core Concepts

### Workspace

A logical grouping of applications that run together, sharing an internal network for communication while isolated from other workspaces.

```
workspace: dev
├── app-one (project: dev-app-one)
├── app-two (project: dev-app-two)
└── app-three (project: dev-app-three)
    └── all connected via: dev-internal network
```

### Application

A Docker Compose-based service that can participate in workspaces. Applications define a **service contract** (`application.yaml`) that declares what they export to the workspace, but their core `docker-compose.yaml` remains workspace-agnostic.

### Single-Application Workspace

Contrail supports workspaces where the workspace directory is also the application directory. This allows promoting existing Docker Compose projects without restructuring:

```bash
cd ~/my-existing-project
contrail workspace init --workspace=dev
contrail app init --app=myapp
```

### Port Types and Proxying

| Type | Protocol | Behavior | External Access | Internal Access | Use Case |
|------|----------|----------|-----------------|-----------------|----------|
| **proxied** | `https` | Traefik HTTPS proxy | Via hostname (port 443) | Via alias | Web frontends, public APIs |
| **proxied** | `http` | Traefik HTTP proxy | Via hostname (port 80) | Via alias | HTTP services, redirects |
| **assigned** | - | Direct port binding | Via assigned host port | Via alias | Databases, caches, debug ports |
| **proxied** | `tcp`, etc. | SNI TCP proxy (future) | Via hostname | Via alias | Database GUIs, external tools |

Services not listed in `exported_services` remain **private**—only accessible within the application's own Docker Compose network.

### Visibility

Each port can have a `visibility` of `public` or `protected`. This is primarily **documentation** to communicate intent to collaborators—it does not change Contrail's behavior. Both public and protected proxied services route through Traefik.

### Flavor

A named configuration that specifies which Docker Compose files to use when running an application. Enables "lite" vs "full" modes without duplicating configuration.

---

## Architecture

### Network Topology

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              HOST                                       │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        PROXY LAYER                               │   │
│  │  ┌──────────┐                                                    │   │
│  │  │ Traefik  │◄─────── proxy (external network)                   │   │
│  │  │(proxied  │              │                                     │   │
│  │  │  types)  │              │                                     │   │
│  │  └──────────┘              │                                     │   │
│  └────────────────────────────┼─────────────────────────────────────┘   │
│                               │                                         │
│  ┌────────────────────────────┼─────────────────────────────────────┐   │
│  │         WORKSPACE: dev     │                                     │   │
│  │                            ▼                                     │   │
│  │            ┌─── dev-internal (workspace network) ───┐            │   │
│  │            │                                        │            │   │
│  │    ┌───────┴───────┐ ┌───────┴───────┐ ┌───────────┴───┐        │   │
│  │    │   app-one     │ │   app-two     │ │   app-three   │        │   │
│  │    │ (dev-app-one) │ │ (dev-app-two) │ │(dev-app-three)│        │   │
│  │    └───────────────┘ └───────────────┘ └───────────────┘        │   │
│  │                                                                  │   │
│  │    Internal aliases: app-one-web, app-two-api, app-one-db, ...  │   │
│  │    Proxied hostnames: dev-app-one-web.contrail.test, ...        │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │         WORKSPACE: review  (completely isolated)                 │   │
│  │    Internal aliases: app-one-web, app-two-api, ... (same names!) │   │
│  │    Proxied hostnames: review-app-one-web.contrail.test, ...      │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

### Key Architectural Decisions

#### Decision: Docker Compose Project Name Isolation

**Context**: Need to run multiple instances of the same application simultaneously.

**Decision**: Use Docker Compose's native `--project-name` (or `name:` in compose file) to create isolated namespaces. Each application in a workspace gets project name `{workspace}-{application}`.

**Rationale**: This is Docker's official mechanism for running multiple copies of the same stack. It isolates containers, networks, and volumes without requiring modifications to the application.

#### Decision: Two-Layer Networking

**Context**: Services need both external access (via reverse proxy) and internal access (between applications).

**Decision**: 
- `proxy` network: Host-wide, connects Traefik to public services
- `{workspace}-internal` network: Per-workspace, connects all applications for internal communication

**Rationale**: Separating concerns allows public services to be routable via Traefik while protected services remain internal. The workspace-internal network provides isolation between workspaces.

#### Decision: Pure Overlay (Applications Remain Workspace-Agnostic)

**Context**: Applications could embed workspace configuration, or it could be applied externally.

**Decision**: Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

**Rationale**: 
- Applications can run standalone without Contrail
- No vendor lock-in or special conventions in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems

#### Decision: Convention-Based Naming

**Context**: Hostnames and aliases could be explicitly configured or derived from conventions.

**Decision**: Derive names from conventions:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Protected alias: `{app}-{service}`
- Network name: `{workspace}-internal`

**Rationale**: Conventions reduce configuration, ensure consistency, and make the system predictable. Explicit overrides were considered but removed to keep the schema simple.

#### Decision: Structure vs State Separation

**Context**: Configuration could include runtime choices (which branch, which flavor) or only structural definitions.

**Decision**: Separate structure (what exists) from state (what's active):

| Aspect | Structure (config files) | State (runtime) |
|--------|--------------------------|-----------------|
| What apps exist | workspace.yaml | - |
| Available flavors | application.yaml | - |
| Active flavor | - | .generated/state.yaml or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

**Rationale**: 
- Configuration files describe the system's shape, not its current state
- State changes frequently; structure changes rarely
- Avoids polluting config files with transient information
- Branch management stays with git where it belongs

#### Decision: Three Configuration Schemas

**Context**: Configuration could be in one monolithic file or separated by concern.

**Decision**: Three schema types that can be combined:
- `proxy`: Global/per-user (domain, Traefik settings)
- `workspace`: Per-workspace (name, application list)
- `application`: Per-application (flavors, exported services)

**Rationale**: Separation of concerns—proxy config rarely changes, workspace config defines the environment, application config is owned by the application team.

#### Decision: Port Type System for Exported Services

**Context**: Services need different handling based on how they're accessed—some need HTTP proxying, others need direct port binding.

**Decision**: Each exported service declares ports with a `type` (routing mechanism) and optionally a `protocol`:
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

**Rationale**: 
- `type` determines routing: `proxied` (through Traefik) or `assigned` (direct port binding)
- `protocol` specifies how proxied traffic is handled: `http`, `https`, or future SNI types
- Supports multiple protocols on the same exported service (both HTTP and HTTPS)
- Environment variables use proxy values (port 80/443) for proxied services
- Enables future plugin system for additional protocols (postgresql, mysql SNI routing)
- `visibility` remains as documentation for collaborators

#### Decision: Traefik for Reverse Proxy

**Context**: Need a reverse proxy that can dynamically route to containers.

**Decision**: Use Traefik with Docker provider, reading labels from containers.

**Rationale**: Traefik's Docker integration allows dynamic routing without config file changes. Labels on containers (added via generated overrides) define routing rules.

#### Decision: `up`/`down` Command Semantics

**Context**: Commands could use `start`/`stop` or `up`/`down` terminology.

**Decision**: Use `up` and `down` as primary commands, matching Docker Compose semantics:
- `up`: Build, create networks/volumes, generate overrides, start containers
- `down`: Stop containers, remove containers/networks, optionally remove volumes

**Rationale**:
- Semantic alignment with Docker Compose, which users already know
- `up` conveys "bring the environment into existence" (more than just starting)
- `down` conveys "tear down" rather than just pausing
- Matches the underlying `docker compose up/down` commands Contrail invokes

#### Decision: Options-Based Targeting with Context Detection

**Context**: Commands need to target specific workspaces and applications.

**Decision**: Use `--workspace` and `--app` options (not positional arguments) with automatic context detection from current directory.

**Rationale**:
- Consistent command structure across all commands
- Context detection reduces typing for common workflows
- Explicit flags available when needed
- Easy to extend with additional targeting options

---

## CLI Interface

Contrail provides a comprehensive command-line interface following the pattern:

```
contrail [resource] [action] [--options...]
```

**Resources**: `workspace`, `app`, `flavor`, `port`, `proxy`, `config`

**Key features**:
- Context detection from current directory (workspace and application)
- Consistent `--workspace` / `-w` and `--app` / `-a` flags
- Top-level aliases for common operations (`up`, `down`, `ps`, `logs`)
- Multiple output formats (`--json`, `--yaml`, `--quiet`)

For complete CLI documentation, see **[Contrail CLI Reference](./contrail-cli-reference.md)**.

### Quick Reference

```bash
# Workspace lifecycle
contrail workspace init --workspace=dev
contrail workspace up [-w NAME]
contrail workspace down [-w NAME]
contrail workspace status [-w NAME]

# Application management
contrail app add --app=NAME --repo=URL
contrail app up [-a NAME]
contrail app logs [-a NAME]

# Flavor management
contrail flavor list [-a NAME]
contrail flavor set FLAVOR [-a NAME]

# Port management
contrail port list
contrail port gc

# Shortcuts (with context detection)
contrail up                    # Bring up current workspace
contrail down                  # Tear down current workspace
contrail logs                  # View logs for current app
```

---

## Configuration Schema

For detailed configuration schemas, see **[Contrail Technical Specification](./contrail-technical-spec.md)**.

### Overview

| File | Location | Purpose |
|------|----------|---------|
| `proxy.yaml` | `~/.config/contrail/` | Global proxy settings (domain, etc.) |
| `state.yaml` | `~/.config/contrail/` | Global port assignments |
| `workspace.yaml` | Workspace root | Workspace definition and applications |
| `application.yaml` | Application directory | Service contract and flavors |

### Workspace Configuration Example

```yaml
workspace:
  name: dev
  applications:
    app-one:
      repository: git@github.com:org/app-one.git
    app-two:
      repository: git@github.com:org/app-two.git
    app-three:
      path: .  # Single-app workspace: app in workspace root
```

### Application Configuration Example

```yaml
default_flavor: default

flavors:
  default:
    compose_files:
      - docker-compose.yaml
  full:
    compose_files:
      - docker-compose.yaml
      - docker-compose.worker.yaml

exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

---

## Service Discovery

All exported services receive environment variables for discovering other services:

**Base pattern**: `CONTRAIL_{APP}_{EXPORTED_SERVICE}_{SUFFIX}`

**Protocol-specific pattern**: `CONTRAIL_{APP}_{EXPORTED_SERVICE}_{PROTOCOL}_{SUFFIX}`

**Variables by port type**:

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | ✓ | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | ✓ | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | ✓ | Both |
| `assigned` | - | Internal alias | Assigned port | ✗ | ✗ | ✗ |

**Examples**:
```bash
# Proxied service with HTTPS (base variables)
CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_PORT=443
CONTRAIL_APP_ONE_WEB_SCHEME=https
CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test

# Protocol-specific variables
CONTRAIL_APP_ONE_WEB_HTTPS_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTPS_PORT=443
CONTRAIL_APP_ONE_WEB_HTTPS_URL=https://dev-app-one-web.contrail.test

# Assigned port service
CONTRAIL_APP_ONE_DB_HOST=app-one-db
CONTRAIL_APP_ONE_DB_PORT=5432
```

Applications should use these environment variables rather than hardcoding hostnames.

---

## Directory Structure

```
~/.config/contrail/
├── proxy.yaml                        # Global proxy configuration
└── state.yaml                        # Global port assignments

project-root/
├── proxy/
│   └── docker-compose.yaml           # Traefik configuration
│
├── dev/                              # Workspace
│   ├── workspace.yaml                # Workspace configuration
│   ├── overrides/                    # Manual overrides (optional)
│   │   └── app-one.yaml
│   ├── .generated/                   # Generated files (gitignored)
│   │   ├── state.yaml                # Flavor choices
│   │   ├── manifest.yaml             # Computed values (read-only)
│   │   ├── app-one.override.yaml
│   │   └── app-two.override.yaml
│   ├── app-one/                      # Cloned application
│   │   ├── application.yaml          # Service contract
│   │   └── docker-compose.yaml       # Application compose
│   └── app-two/
│
├── review/                           # Another workspace
│   └── ...
│
└── control/                          # Another workspace
    └── ...
```

### Single-Application Workspace

```
~/my-project/                         # Workspace AND application
├── workspace.yaml                    # workspace.name = "dev"
├── application.yaml                  # Application config
├── docker-compose.yaml               # Existing compose file
├── .generated/                       # Generated files
└── src/                              # Application source
```

---

## Future Considerations

### Port Type Plugins (Future)

**Context**: Different services need different proxying strategies—HTTP/HTTPS is handled by Traefik's HTTP routers, but databases need TCP routing with SNI.

**Consideration**: Plugin system where protocols can register handlers:
```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: proxied
        protocol: postgresql          # Plugin handles this protocol
        port: 5432
        visibility: public
```

Plugins would generate appropriate Traefik configuration (TCP routers, SNI rules) for their protocol.

### Application Dependencies (Future)

**Context**: Some applications may need others to be running first.

**Consideration**: Add dependency ordering:
```yaml
workspace:
  applications:
    app-two:
      depends_on:
        - app-one
```

### Shared Volumes (Future)

**Context**: Applications might need to share files (uploads, assets).

**Consideration**: Workspace-level volume definitions that can be mounted into multiple applications.

### Health Checks (Future)

**Context**: Starting applications in order isn't sufficient if they need warm-up time.

**Consideration**: Integration with Docker health checks to wait for readiness.

---

## Non-Goals

1. **Kubernetes support**: Contrail is specifically for Docker Compose environments
2. **Production deployment**: Focused on local development and testing
3. **CI/CD integration**: May come later, but not initial focus
4. **GUI**: CLI-first; GUI could be added later
5. **Image building**: Uses existing images; doesn't manage builds
6. **Secret management**: Uses Docker Compose's existing mechanisms

---

## Success Criteria

1. **Zero application changes**: Existing Docker Compose applications work without modification
2. **Parallel environments**: Can run dev, review, and control simultaneously without conflicts
3. **Predictable naming**: Given workspace and app names, hostnames and aliases are deterministic
4. **Fast iteration**: Switching flavors or regenerating config takes seconds
5. **Debuggable**: Generated files are human-readable; easy to understand what's happening

---

## Related Documentation

- **[Contrail CLI Reference](./contrail-cli-reference.md)**: Complete CLI documentation
- **[Contrail Technical Specification](./contrail-technical-spec.md)**: Architecture, schemas, and implementation details
- **[Contrail Shell Integration](./contrail-shell-integration.md)**: Shell functions, completion, and Docker Compose passthrough
- **[Contrail Go Stack](./contrail-go-stack.md)**: Go technology stack, dependencies, and project scaffolding

---

## Appendix: Terminology

| Term | Definition |
|------|------------|
| **Workspace** | An isolated environment containing multiple applications |
| **Application** | A Docker Compose-based service that participates in workspaces |
| **Flavor** | A named configuration specifying which compose files to use |
| **Service Contract** | The `application.yaml` file defining what an application exports |
| **Exported Service** | A named export in `application.yaml`, may map to a different Compose service |
| **Override File** | Generated Docker Compose file that adds workspace integration |
| **Manifest** | Generated read-only file showing computed hostnames, ports, and environment variables |
| **Port Type** | How a port is routed: `proxied` (through Traefik) or `assigned` (direct port binding) |
| **Protocol** | For proxied types, the traffic protocol: `http`, `https`, or future SNI types |
| **Visibility** | Documentation-only flag (`public`/`protected`) indicating intended use |
| **Alias** | A DNS name on the workspace-internal network |
| **Context Detection** | Automatic discovery of workspace/app from current directory |

---

## Appendix: Template Customization (Advanced)

While Contrail provides sensible defaults for hostname and alias generation, advanced users can customize templates at the workspace level. See the Technical Specification for details on available template variables.

Default templates:
- Hostname: `{workspace}-{app}-{export}.{domain}` → `dev-app-one-web.contrail.test`
- Alias: `{app}-{export}` → `app-one-web`
- Project name: `{workspace}-{app}` → `dev-app-one`

---

## Appendix: Comparison with Related Tools

| Feature | Contrail | Docker `include` | DDEV/Lando | Tilt/Garden |
|---------|----------|------------------|------------|-------------|
| Multi-app orchestration | ✓ | ✓ (merged) | ✗ | ✓ |
| Parallel workspace instances | ✓ | ✗ | ✗ | ✗ |
| Apps remain agnostic | ✓ | N/A | N/A | ✗ |
| Docker Compose native | ✓ | ✓ | ✓ | ✗ (K8s) |
| Generated integration | ✓ | ✗ | ✓ | ✓ |
| Service discovery | ✓ | Manual | Limited | ✓ |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0-draft | Dec 2024 | Initial PRD based on design session |
| 0.2.0-draft | Dec 2024 | Port type system, manifest file, global port state, CONTRAIL_* env vars |
| 0.3.0-draft | Dec 2024 | Type/protocol split, protocol-specific env vars, assigned port binding, port inventory, contrail.test domain |
| 0.4.0-draft | Dec 2024 | CLI redesign (options-based, context detection, up/down semantics), single-app workspace support, extracted CLI reference document |
| 0.5.0-draft | Dec 2024 | Added `contrail-compose` shell function concept to product vision; linked shell integration specification |
