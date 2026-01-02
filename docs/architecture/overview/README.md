# Contrail Architecture Overview

**Version**: 0.5.0
**Date**: December 2024

---

## System Context

Contrail is a CLI tool that orchestrates Docker Compose-based applications into isolated workspaces. It generates Docker Compose override files and manages networking to enable multiple parallel instances of multi-application stacks.

---

## Network Topology

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

---

## Network Layers

### Proxy Network

- **Name**: `contrail-proxy`
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

## Component Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User's Shell                                │
│                                                                     │
│   contrail workspace up          contrail-compose exec php bash     │
│         │                                    │                      │
│         ▼                                    ▼                      │
│   ┌───────────┐                    ┌─────────────────┐              │
│   │ contrail  │                    │ contrail-compose│              │
│   │  binary   │                    │ shell function  │              │
│   └───────────┘                    └────────┬────────┘              │
│         │                                   │                       │
│         ▼                                   ▼                       │
│   ┌─────────────┐                 ┌──────────────────────┐          │
│   │ Override    │                 │ contrail compose-    │          │
│   │ Generator   │                 │ prefix command       │          │
│   └─────────────┘                 └──────────┬───────────┘          │
│         │                                    │                      │
│         ▼                                    ▼                      │
│   ┌─────────────────────────────────────────────────────┐           │
│   │                  docker compose                      │           │
│   └─────────────────────────────────────────────────────┘           │
│         │                                                           │
│         ▼                                                           │
│   ┌─────────────────────────────────────────────────────┐           │
│   │                   Docker Engine                      │           │
│   └─────────────────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

### Standard Multi-Application Workspace

```
~/.config/contrail/
├── proxy.yaml                        # Global proxy configuration
├── state.yaml                        # Global port assignments
└── workspaces.yaml                   # Workspace registry

workspaces/
├── proxy/
│   ├── docker-compose.yaml           # Traefik service definition
│   ├── traefik.yaml                  # Traefik static configuration
│   └── dynamic/                      # Dynamic configuration
│
├── dev/                              # Workspace root
│   ├── workspace.yaml                # Workspace configuration
│   ├── overrides/                    # Manual overrides (optional)
│   ├── .generated/                   # Generated files (gitignored)
│   │   ├── state.yaml                # Runtime state (active flavors)
│   │   ├── manifest.yaml             # Computed values (read-only)
│   │   └── *.override.yaml           # Generated compose overrides
│   ├── app-one/                      # Cloned application
│   │   ├── application.yaml          # Service contract
│   │   └── docker-compose.yaml       # Application compose
│   └── app-two/
│       └── ...
│
├── review/                           # Another workspace
└── control/                          # Another workspace
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

## Port Types

| Type | Protocol | Behavior | External Access | Internal Access |
|------|----------|----------|-----------------|-----------------|
| **proxied** | `https` | Traefik HTTPS proxy | Via hostname (port 443) | Via alias |
| **proxied** | `http` | Traefik HTTP proxy | Via hostname (port 80) | Via alias |
| **assigned** | - | Direct port binding | Via assigned host port | Via alias |

---

## Proxy Layer

The proxy layer (Traefik) is a prerequisite for proxied services:

1. **Bootstrap**: `contrail proxy init` creates the proxy configuration
2. **Auto-start**: `workspace up` automatically starts the proxy if not running
3. **Shared instance**: A single proxy instance serves all workspaces on the host

---

## Key Design Principles

1. **Pure Overlay**: Applications remain workspace-agnostic; all integration via generated overrides
2. **Convention-Based Naming**: Predictable hostnames, aliases, and project names
3. **Structure vs State Separation**: Configuration describes shape; runtime manages state
4. **Two-Layer Networking**: Proxy network for external access; internal network for app communication

---

## Related Documentation

- [Product Vision](../../product/vision/README.md)
- [Configuration Schemas Spec](../../specs/configuration-schemas/README.md)
- [Naming Conventions Spec](../../specs/naming-conventions/README.md)
