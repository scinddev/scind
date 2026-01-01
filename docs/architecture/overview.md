# Contrail: Architecture Overview

**Version**: 0.5.0
**Date**: December 2024
**Status**: Active

<!-- Migrated from specs/contrail-prd.md and specs/contrail-technical-spec.md -->

---

## System Context

### Overview

Contrail is a workspace orchestration system that enables developers to run multiple isolated instances of multi-application Docker Compose stacks simultaneously on a single host. It provides a thin coordination layer over Docker Compose that generates integration configuration while keeping applications workspace-agnostic.

### Context Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              HOST                                       │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        PROXY LAYER                               │   │
│  │  ┌──────────┐                                                    │   │
│  │  │ Traefik  │◄─────── contrail-proxy (external network)         │   │
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

### External Dependencies

| System | Purpose | Protocol |
|--------|---------|----------|
| Docker Engine | Container runtime | Docker API |
| Docker Compose | Application orchestration | CLI / API |
| Traefik | Reverse proxy for HTTP/HTTPS routing | Docker labels |
| Git | Repository cloning for applications | SSH/HTTPS |

---

## Containers

### Container Diagram

```
┌───────────────────────────────────────────────────────────────┐
│                        CONTRAIL SYSTEM                         │
│                                                                │
│  ┌────────────────┐                                           │
│  │   Contrail     │     Generates override files              │
│  │     CLI        │     Manages workspace lifecycle           │
│  │   (Go binary)  │     Controls proxy and networks           │
│  └───────┬────────┘                                           │
│          │                                                     │
│          │ docker compose commands                             │
│          ▼                                                     │
│  ┌────────────────────┐     ┌─────────────────────┐           │
│  │  Docker Compose    │────▶│   Docker Engine     │           │
│  │  (orchestration)   │     │   (containers)      │           │
│  └────────────────────┘     └──────────┬──────────┘           │
│                                        │                       │
│          ┌─────────────────────────────┼───────────────┐      │
│          │                             │               │      │
│          ▼                             ▼               ▼      │
│  ┌──────────────┐          ┌──────────────┐   ┌──────────────┐│
│  │   Traefik    │          │  Application │   │  Application ││
│  │   (proxy)    │◄────────▶│  Containers  │   │  Containers  ││
│  └──────────────┘          └──────────────┘   └──────────────┘│
│                                                                │
└───────────────────────────────────────────────────────────────┘
```

### Container Descriptions

| Container | Technology | Purpose |
|-----------|------------|---------|
| Contrail CLI | Go binary | User interface, configuration management, override file generation |
| Traefik | Docker image | Reverse proxy for HTTP/HTTPS routing, dynamic configuration via Docker labels |
| Application Containers | Various | User applications running via Docker Compose, managed per-workspace |

---

## Networks

### Network Types

| Network | Scope | Purpose | Created By |
|---------|-------|---------|------------|
| `contrail-proxy` | Host-wide | Connects Traefik to services needing external access | `proxy init` / `proxy up` |
| `{workspace}-internal` | Per-workspace | Inter-application communication within workspace | `workspace up` (lazy) |
| Application networks | Per-application | Internal service communication | Docker Compose (automatic) |

### Network Diagram

```
┌────────────────────────────────────────────────────────────────────────┐
│                              HOST                                       │
│                                                                         │
│  ┌──────────────────── contrail-proxy (shared) ───────────────────┐    │
│  │                                                                 │    │
│  │  ┌──────────┐                                                   │    │
│  │  │ Traefik  │◄────────────────────────────┐                     │    │
│  │  └──────────┘                             │                     │    │
│  │       ▲                                   │                     │    │
│  └───────┼───────────────────────────────────┼─────────────────────┘    │
│          │                                   │                          │
│  ┌───────┼─── dev-internal ──────────────────┼─────────────────────┐    │
│  │       │                                   │                     │    │
│  │  ┌────┴────┐   ┌─────────┐   ┌─────────┐ │                     │    │
│  │  │ app-one │   │ app-two │   │app-three│ │                     │    │
│  │  │   web   │   │   api   │   │   web   │─┘  (connects to proxy) │    │
│  │  └─────────┘   └─────────┘   └─────────┘                       │    │
│  │      ▲              ▲            ▲                              │    │
│  │      └──────────────┴────────────┘                              │    │
│  │           Internal communication via aliases                    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  ┌───────── review-internal (isolated) ────────────────────────────┐    │
│  │  (Same structure, completely separate namespace)                │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
└────────────────────────────────────────────────────────────────────────┘
```

---

## Communication Patterns

### Internal Communication

| From | To | Method | Purpose |
|------|----|--------|---------|
| Application A | Application B | Internal network alias | Service-to-service calls within workspace |
| Contrail CLI | Docker Engine | Docker API | Container and network management |
| Contrail CLI | Docker Compose | CLI subprocess | Application lifecycle operations |

### External Communication

| Direction | Endpoint | Purpose |
|-----------|----------|---------|
| Inbound | Traefik (:80/:443) | HTTP/HTTPS requests routed to applications |
| Inbound | Assigned ports | Direct database/service access |
| Outbound | Git repositories | Application source code cloning |

---

## Key Components

### Contrail CLI Components

| Component | Purpose |
|-----------|---------|
| Context Detector | Walks directory tree to find workspace.yaml and application.yaml |
| Config Loader | Merges proxy, workspace, and application configuration |
| Override Generator | Creates Docker Compose override files with networks, aliases, labels |
| Manifest Builder | Computes and caches resolved hostnames, ports, environment variables |
| Port Manager | Tracks and assigns host ports across all workspaces |

### Generated Files

| File | Location | Purpose |
|------|----------|---------|
| `{app}.override.yaml` | `.generated/` | Docker Compose override with networks, labels, env vars |
| `state.yaml` | `.generated/` | Active flavor per application |
| `manifest.yaml` | `.generated/` | Computed read-only view of workspace topology |

---

## Cross-Cutting Concerns

### Isolation

Workspaces achieve isolation through:
- **Docker Compose project names**: `{workspace}-{application}` ensures unique container/volume names
- **Dedicated networks**: Each workspace has its own internal network
- **Convention-based naming**: Hostnames include workspace prefix for uniqueness

See [ADR-0001: Docker Compose Project Name Isolation](../decisions/0001-docker-compose-project-name-isolation.md) and [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md).

### Service Discovery

Applications discover each other via injected environment variables:
- `CONTRAIL_{APP}_{SERVICE}_HOST` — Hostname or alias
- `CONTRAIL_{APP}_{SERVICE}_PORT` — Port number
- `CONTRAIL_{APP}_{SERVICE}_URL` — Full URL (for proxied services)

### Configuration Layering

Configuration is loaded with precedence (highest to lowest):
1. Command-line flags
2. Environment variables (`CONTRAIL_*`)
3. Workspace config (`workspace.yaml`)
4. Application config (`application.yaml`)
5. Global config (`~/.config/contrail/proxy.yaml`)
6. Built-in defaults

See [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md).

---

## Quality Attributes

| Attribute | Requirement | How Achieved |
|-----------|-------------|--------------|
| Portability | Applications run standalone without Contrail | Pure overlay design—no app modifications required |
| Isolation | Workspaces don't interfere with each other | Separate networks, unique project names |
| Predictability | Given names, all derived values are deterministic | Convention-based naming throughout |
| Developer Experience | Minimal typing for common operations | Context detection from current directory |

---

## Known Risks & Technical Debt

| Risk/Debt | Impact | Mitigation |
|-----------|--------|------------|
| Concurrent operations | Race conditions if multiple terminals run commands simultaneously | Document: use one terminal per workspace for operations |
| Port garbage collection | Orphaned ports from filesystem-deleted workspaces | `contrail port gc` command for manual cleanup |
| Docker dependency | Requires Docker Engine and Compose | `contrail doctor` checks prerequisites |

---

## Directory Structure

```
~/.config/contrail/
├── proxy.yaml                        # Global proxy configuration
├── state.yaml                        # Global port assignments
├── workspaces.yaml                   # Workspace registry
└── proxy/                            # Traefik Docker Compose project
    ├── docker-compose.yaml
    ├── traefik.yaml
    ├── dynamic/
    └── certs/

{workspace}/
├── workspace.yaml                    # Workspace configuration
├── overrides/                        # Manual overrides (optional)
├── .generated/                       # Generated files (gitignored)
│   ├── state.yaml
│   ├── manifest.yaml
│   └── {app}.override.yaml
└── {app}/                            # Application directories
    ├── application.yaml              # Service contract
    └── docker-compose.yaml
```

---

## Related Documents

- [Vision](../product/vision.md)
- [Decisions](../decisions/)
- [Specifications](../specs/)
- [CLI Reference](../reference/cli.md)
