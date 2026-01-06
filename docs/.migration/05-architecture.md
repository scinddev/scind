# Migration Step: Layer 3 - Architecture

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 file, approximately 200 lines

---

## Overview

Create `architecture/overview.md` using the C4-Lite template, consolidating architecture content.

**Source Sections**:
- `specs/contrail-prd.md:111-159` - Architecture diagram, proxy layer
- `specs/contrail-technical-spec.md:22-94` - Architecture overview, networks

---

## Output File: `architecture/overview.md`

**Source**: Multiple files (see sections)

### Content

```markdown
# Contrail Architecture Overview

This document provides a high-level view of Contrail's architecture using C4-Lite approach.

---

## System Context

Contrail orchestrates Docker Compose applications within isolated workspaces, providing:
- Network isolation between workspaces
- Reverse proxy routing via Traefik
- Generated configuration for workspace integration

### Key Components

| Component | Description |
|-----------|-------------|
| **Contrail CLI** | Command-line tool for managing workspaces and applications |
| **Traefik Proxy** | Shared reverse proxy routing external traffic to containers |
| **Docker Compose** | Underlying container orchestration |
| **Workspace Network** | Per-workspace internal network for application communication |

---

## Container Diagram

```
                              HOST

   PROXY LAYER
   ┌──────────┐
   │ Traefik  │◄─────── contrail-proxy (external network)
   │(proxied  │              │
   │  types)  │              │
   └──────────┘              │
                             │
          WORKSPACE: dev     │
                             ▼
             ┌─── dev-internal (workspace network) ───┐
             │                                        │
     ┌───────┴───────┐ ┌───────┴───────┐ ┌───────────┴───┐
     │   app-one     │ │   app-two     │ │   app-three   │
     │ (dev-app-one) │ │ (dev-app-two) │ │(dev-app-three)│
     └───────────────┘ └───────────────┘ └───────────────┘

     Internal aliases: app-one-web, app-two-api, app-one-db, ...
     Proxied hostnames: dev-app-one-web.contrail.test, ...

          WORKSPACE: review  (completely isolated)
     Internal aliases: app-one-web, app-two-api, ... (same names!)
     Proxied hostnames: review-app-one-web.contrail.test, ...
```

---

## Network Topology

### Proxy Network

- **Name**: `contrail-proxy`
- **Scope**: Host-level, shared across all workspaces
- **Purpose**: Connects Traefik to services that need external access
- **Created by**: Proxy layer setup (once per host)

### Workspace Internal Network

- **Name**: `{workspace-name}-internal` (e.g., `dev-internal`)
- **Scope**: Per-workspace
- **Purpose**: Enables inter-application communication within a workspace using stable aliases
- **Created by**: `workspace up` (lazy, idempotent - created if not exists)

### Application Default Networks

- **Name**: Managed by Docker Compose per application
- **Scope**: Per-application
- **Purpose**: Internal communication between services within a single application
- **Created by**: Docker Compose (automatic)

---

## Proxy Layer Bootstrap

The proxy layer (Traefik) is a prerequisite for proxied services. Contrail manages the proxy lifecycle:

1. **Bootstrap**: `contrail proxy init` creates the proxy configuration at `~/.config/contrail/proxy/`
2. **Auto-start**: `workspace up` automatically starts the proxy if not running
3. **Shared instance**: A single proxy instance serves all workspaces on the host

The proxy is implemented as a Docker Compose project managed by Contrail. Users generally don't interact with it directly - workspaces connect to it automatically via Docker labels.

---

## Data Flow

### External Request Flow (Proxied Services)

```
Browser → *.contrail.test:443
    → Traefik (contrail-proxy network)
    → Container (via Docker labels)
```

### Internal Request Flow (Between Applications)

```
App-One Container → app-two-api (alias)
    → {workspace}-internal network
    → App-Two API Container
```

### Assigned Port Flow

```
Database Client → localhost:5432
    → Host Port Binding
    → Container Port
```

---

## Goals

- **Application independence**: Individual applications remain unaware of the workspace system
- **Pure overlay**: All workspace integration is achieved through Docker Compose override files
- **Inter-application communication**: Applications within a workspace can communicate via stable, predictable internal hostnames
- **External access**: A shared reverse proxy routes external requests to the appropriate workspace and service
- **Multiple workspaces**: The same set of applications can be instantiated multiple times with different workspace names

<!-- Migrated from specs/contrail-prd.md:111-159, specs/contrail-technical-spec.md:22-94 -->
```

---

## Also Create: `architecture/README.md`

```markdown
# Architecture Documentation

This directory contains architecture documentation for Contrail.

## Contents

| Document | Description |
|----------|-------------|
| [overview.md](./overview.md) | System context and container diagrams |
```

---

## Completion Checklist

- [ ] `architecture/overview.md` created
- [ ] `architecture/README.md` created
- [ ] All diagrams preserved exactly
- [ ] Source attribution present
