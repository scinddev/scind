<!-- Migrated from multiple sources:
     - specs/scind-technical-spec.md:9-94
     - specs/scind-prd.md:150-159
-->
<!-- Extraction ID: architecture-overview -->

## Overview

A **workspace** is a logical grouping of Docker Compose-based applications that run together on a single host, sharing internal networking for inter-service communication while maintaining isolation from other workspaces. This enables running multiple complete copies of the same application stack simultaneously (e.g., development, code review, and stable/control environments).

## Goals

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
│  │    │   app-one     │  │   app-two     │  │   app-three   │        │  │
│  │    │ (dev-app-one) │  │ (dev-app-two) │  │(dev-app-three)│        │  │
│  │    │               │  │               │  │               │        │  │
│  │    │ ┌───┐ ┌───┐   │  │ ┌───┐ ┌───┐   │  │ ┌───┐ ┌───┐   │        │  │
│  │    │ │web│ │ db│   │  │ │web│ │api│   │  │ │web│ │wrk│   │        │  │
│  │    │ └───┘ └───┘   │  │ └───┘ └───┘   │  │ └───┘ └───┘   │        │  │
│  │    └───────────────┘  └───────────────┘  └───────────────┘        │  │
│  │                                                                   │  │
│  │    Aliases on dev-internal:                                       │  │
│  │      app-one-web, app-two-web, app-two-api, app-three-web, ...    │  │
│  │                                                                   │  │
│  │    External hostnames (via Traefik):                              │  │
│  │      dev-app-one-web.scind.test, dev-app-two-web.scind.test       │  │
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

- **Name**: `{workspace-name}-internal` (e.g., `dev-internal`)
- **Scope**: Per-workspace
- **Purpose**: Enables inter-application communication within a workspace using stable aliases
- **Created by**: `workspace up` (lazy, idempotent—created if not exists)

### Application Default Networks

- **Name**: Managed by Docker Compose per application
- **Scope**: Per-application
- **Purpose**: Internal communication between services within a single application
- **Created by**: Docker Compose (automatic)

---

## Proxy Layer Bootstrap

The proxy layer (Traefik) is a prerequisite for proxied services. Scind manages the proxy lifecycle:

1. **Bootstrap**: `scind proxy init` creates the proxy configuration at `~/.config/scind/proxy/`
2. **Auto-start**: `workspace up` automatically starts the proxy if not running
3. **Shared instance**: A single proxy instance serves all workspaces on the host

The proxy is implemented as a Docker Compose project managed by Scind. Users generally don't interact with it directly—workspaces connect to it automatically via Docker labels.
