# Architecture Overview

> Scind system architecture and component relationships.

---

## Overview

A **workspace** is a logical grouping of Docker Compose-based applications that run together on a single host, sharing internal networking for inter-service communication while maintaining isolation from other workspaces. This enables running multiple complete copies of the same application stack simultaneously (e.g., development, code review, and stable/control environments).

## Goals

The architecture achieves the product vision (see [Vision](../product/vision.md)) through these principles:

- **Application independence**: Individual applications remain unaware of the workspace system. No special labels, naming conventions, or workspace-specific configuration required in the application's own `docker-compose.yaml`.
- **Pure overlay**: All workspace integration is achieved through Docker Compose override files that are generated and managed externally.
- **Inter-application communication**: Applications within a workspace can communicate via stable, predictable internal hostnames that don't change based on the workspace name.
- **External access**: A shared reverse proxy (Traefik) routes external requests to the appropriate workspace and service based on hostname.
- **Multiple workspaces**: The same set of applications can be instantiated multiple times with different workspace names, running simultaneously without conflict.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              HOST                                       │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                        PROXY LAYER                                │  │
│  │  ┌──────────┐                                                     │  │
│  │  │ Traefik  │◄─────── scind-proxy (external network)              │  │
│  │  └──────────┘              │                                      │  │
│  └────────────────────────────┼──────────────────────────────────────┘  │
│                               │                                         │
│  ┌────────────────────────────┼──────────────────────────────────────┐  │
│  │                            │         WORKSPACE: dev               │  │
│  │                            ▼                                      │  │
│  │            ┌─── dev-internal (workspace network) ───┐             │  │
│  │            │                                        │             │  │
│  │    ┌───────┴───────┐  ┌───────┴───────┐  ┌──────────┴────┐        │  │
│  │    │   frontend    │  │   backend     │  │   shared-db   │        │  │
│  │    │(dev-frontend) │  │ (dev-backend) │  │(dev-shared-db)│        │  │
│  │    │               │  │               │  │               │        │  │
│  │    │ ┌───┐ ┌───┐   │  │ ┌───┐ ┌───┐   │  │ ┌───┐ ┌───┐   │        │  │
│  │    │ │web│ │ db│   │  │ │web│ │api│   │  │ │web│ │wrk│   │        │  │
│  │    │ └───┘ └───┘   │  │ └───┘ └───┘   │  │ └───┘ └───┘   │        │  │
│  │    └───────────────┘  └───────────────┘  └───────────────┘        │  │
│  │                                                                   │  │
│  │    Aliases on dev-internal:                                       │  │
│  │      frontend-web, backend-web, backend-api, shared-db-db, ...    │  │
│  │                                                                   │  │
│  │    External hostnames (via Traefik):                              │  │
│  │      dev-frontend-web.scind.test, dev-backend-api.scind.test      │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                        WORKSPACE: review                          │  │
│  │                            ...                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                        WORKSPACE: control                         │  │
│  │                            ...                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Networks

### Proxy Network

- **Name**: `scind-proxy`
- **Scope**: Host-level, shared across all workspaces
- **Purpose**: Connects Traefik to services that need external access
- **Created by**: Proxy layer setup (once per host)

### Workspace Internal Network

- **Name**: `{workspace}-internal` (e.g., `dev-internal`)
- **Scope**: Per-workspace
- **Purpose**: Enables inter-application communication within a workspace using stable aliases
- **Created by**: `workspace up` (lazy, idempotent—created if not exists)

### Application Default Networks

- **Name**: Managed by Docker Compose per application
- **Scope**: Per-application
- **Purpose**: Internal communication between services within a single application
- **Created by**: Docker Compose (automatic)

---

## Proxy Layer

The proxy layer (Traefik) is a prerequisite for proxied services. A single shared Traefik instance serves all workspaces on the host, routing external requests to the appropriate workspace and service based on hostname.

TLS termination occurs at the Traefik proxy. See [ADR-0009](../decisions/0009-flexible-tls-configuration.md) for configuration modes.

For bootstrap procedures and lifecycle management, see [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md).

---

## Docker Compose Integration

Scind generates Docker Compose override files that extend application compose files without modifying them. This implements the pure overlay design (see [ADR-0003](../decisions/0003-pure-overlay-design.md)).

For each application in a workspace, Scind generates:
- Network attachments for workspace internal network and proxy network
- Service aliases for internal hostname resolution
- Traefik labels for proxied services
- Port mappings for assigned services

See [Generated Override Files Spec](../specs/generated-override-files.md) for complete details.

---

## CLI Architecture

The `scind` CLI orchestrates all workspace operations.

```
┌─────────────────────────────────────────────────┐
│                   scind CLI                      │
├─────────────────────────────────────────────────┤
│  Commands: proxy | workspace | version | help   │
├──────────┬──────────┬──────────┬───────────────┤
│ Config   │ State    │ Generate │ Docker        │
│ Reader   │ Manager  │ Engine   │ Compose       │
└──────────┴──────────┴──────────┴───────────────┘
      │          │          │           │
      ▼          ▼          ▼           ▼
  *.yaml     state.yaml  .generated/  docker-compose
  configs    workspaces  overrides    commands
```

### Responsibilities
- **Config Reader**: Parses proxy.yaml, workspace.yaml, application.yaml
- **State Manager**: Maintains port assignments, workspace registry
- **Generate Engine**: Creates override files from configuration
- **Docker Compose**: Executes container operations

---

## Configuration Architecture

Scind uses a three-schema configuration system (see [ADR-0006](../decisions/0006-three-configuration-schemas.md)).

```
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Flow                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ~/.config/scind/           workspace/              app/     │
│  ┌──────────────┐    ┌─────────────────┐    ┌─────────────┐ │
│  │  proxy.yaml  │    │ workspace.yaml  │    │application. │ │
│  │  (global)    │    │ (per-workspace) │    │   yaml      │ │
│  └──────┬───────┘    └────────┬────────┘    └──────┬──────┘ │
│         │                     │                     │        │
│         └──────────┬──────────┴──────────┬─────────┘        │
│                    ▼                     ▼                   │
│              ┌───────────┐        ┌──────────────┐          │
│              │  Generate │───────>│ .generated/  │          │
│              │  Engine   │        │ override.yaml│          │
│              └───────────┘        └──────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### Schema Purposes
- **proxy.yaml**: Domain, TLS, Traefik settings (machine-wide)
- **workspace.yaml**: Workspace name, application list (per-workspace)
- **application.yaml**: Flavors, exported services (per-application)

---

## State Management

Scind separates structure (configuration) from state (runtime). See [ADR-0005](../decisions/0005-structure-vs-state-separation.md).

```
┌─────────────────────────────────────────────────────────────┐
│                      State Files                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Global State (~/.config/scind/)                            │
│  ┌────────────────────┐  ┌────────────────────┐             │
│  │   state.yaml       │  │  workspaces.yaml   │             │
│  │ - port_inventory   │  │ - registered       │             │
│  │ - assigned_ports   │  │   workspaces       │             │
│  └────────────────────┘  └────────────────────┘             │
│                                                              │
│  Per-Workspace State (workspace/.generated/)                │
│  ┌────────────────────┐  ┌────────────────────┐             │
│  │   state.yaml       │  │   manifest.yaml    │             │
│  │ - active_flavor    │  │ - computed values  │             │
│  │ - last_generated   │  │ - resolved ports   │             │
│  └────────────────────┘  └────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
```

---

## Port Type System

Services expose ports via two mechanisms (see [ADR-0007](../decisions/0007-port-type-system.md)).

```
┌─────────────────────────────────────────────────────────────┐
│                     Port Types                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  PROXIED (through Traefik)         ASSIGNED (direct bind)   │
│  ┌──────────────────────┐          ┌──────────────────────┐ │
│  │   Browser/Client     │          │   Browser/Client     │ │
│  │         │            │          │         │            │ │
│  │         ▼            │          │         │            │ │
│  │  ┌─────────────┐     │          │         │            │ │
│  │  │   Traefik   │     │          │         │            │ │
│  │  │  :443/:80   │     │          │         │            │ │
│  │  └──────┬──────┘     │          │         │            │ │
│  │         │            │          │         │            │ │
│  │         ▼            │          │         ▼            │ │
│  │  ┌─────────────┐     │          │  ┌─────────────┐     │ │
│  │  │  Container  │     │          │  │  Container  │     │ │
│  │  │   :8080     │     │          │  │  :5432      │     │ │
│  │  └─────────────┘     │          │  └─────────────┘     │ │
│  │                      │          │  (host port 15432)   │ │
│  └──────────────────────┘          └──────────────────────┘ │
│                                                              │
│  Use for: HTTP/HTTPS apps          Use for: Databases,      │
│                                    debug ports, non-HTTP    │
└─────────────────────────────────────────────────────────────┘
```

---

## Related Decisions

- [ADR-0001: Docker Compose Project Name Isolation](../decisions/0001-docker-compose-project-name-isolation.md)
- [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md)
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [ADR-0004: Convention-Based Naming](../decisions/0004-convention-based-naming.md)
- [ADR-0005: Structure vs State Separation](../decisions/0005-structure-vs-state-separation.md)
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md)
- [ADR-0007: Port Type System](../decisions/0007-port-type-system.md)
- [ADR-0008: Traefik for Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)

---

## Related Specifications

For detailed behavioral specifications:

- [Configuration Schemas](../specs/configuration-schemas.md)
- [Generated Override Files](../specs/generated-override-files.md)
- [Workspace Lifecycle](../specs/workspace-lifecycle.md)
- [Proxy Infrastructure](../specs/proxy-infrastructure.md)
- [State Management](../specs/state-management.md)
