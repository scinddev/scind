# Contrail Technical Specification

**Version**: 0.5.2-draft
**Date**: December 2024
**Status**: Design Phase

---

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
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        PROXY LAYER                               │   │
│  │  ┌──────────┐                                                    │   │
│  │  │ Traefik  │◄─────── proxy (external network)                   │   │
│  │  └──────────┘              │                                     │   │
│  └────────────────────────────┼─────────────────────────────────────┘   │
│                               │                                         │
│  ┌────────────────────────────┼─────────────────────────────────────┐   │
│  │                            │         WORKSPACE: dev              │   │
│  │                            ▼                                     │   │
│  │            ┌─── dev-internal (workspace network) ───┐            │   │
│  │            │                                        │            │   │
│  │    ┌───────┴───────┐ ┌───────┴───────┐ ┌───────────┴───┐        │   │
│  │    │   app-one     │ │   app-two     │ │   app-three   │        │   │
│  │    │ (dev-app-one) │ │ (dev-app-two) │ │(dev-app-three)│        │   │
│  │    │               │ │               │ │               │        │   │
│  │    │ ┌───┐ ┌───┐   │ │ ┌───┐ ┌───┐   │ │ ┌───┐ ┌───┐   │        │   │
│  │    │ │web│ │ db│   │ │ │web│ │api│   │ │ │web│ │wrk│   │        │   │
│  │    │ └───┘ └───┘   │ │ └───┘ └───┘   │ │ └───┘ └───┘   │        │   │
│  │    └───────────────┘ └───────────────┘ └───────────────┘        │   │
│  │                                                                  │   │
│  │    Aliases on dev-internal:                                      │   │
│  │      app-one-web, app-two-web, app-two-api, app-three-web, ...   │   │
│  │                                                                  │   │
│  │    External hostnames (via Traefik):                             │   │
│  │      dev-app-one-web.contrail.test, dev-app-two-web.contrail.test│   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                        WORKSPACE: review                         │   │
│  │                            ...                                   │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                        WORKSPACE: control                        │   │
│  │                            ...                                   │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Networks

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

## Port Types and Proxying

Exported services declare ports with a `type` that determines how the port is routed, and optionally a `protocol` for proxied services:

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `*_HOST`, `*_PORT` |

### Port Type Descriptions

- **proxied**: Traffic is routed through Traefik. The exported service gets a hostname (`{workspace}-{app}-{export}.{domain}`) and Traefik labels are generated. Environment variables contain the **proxy values** (hostname and proxy port 80/443), not the container port.
- **assigned**: The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Contrail increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

### Protocol (for proxied type)

When `type: proxied`, the `protocol` field is **required** and specifies how Traefik routes the traffic:

- **https**: Routes through Traefik's `websecure` entrypoint (port 443) with TLS termination
- **http**: Routes through Traefik's `web` entrypoint (port 80)
- **tcp**, **postgresql**, **mysql**, etc. (future): SNI-based TCP routing for database connections. Plugins will handle generating appropriate Traefik TCP router configuration.

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent to collaborators:

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Contrail's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

**Docker label exposure**: Visibility is included in the generated Docker labels (`workspace.visibility=public` or `workspace.visibility=protected`), enabling external tools (such as Servlo) to distinguish between public and protected services for display or filtering purposes.

### Private Services

Services not listed in `exported_services` remain private (standard Docker Compose behavior—only accessible within the application's own compose network).

---

## Directory Structure

### Standard Multi-Application Workspace

```
~/.config/contrail/
├── proxy.yaml                        # Global proxy configuration
└── state.yaml                        # Global port assignments and inventory

workspaces/
├── proxy/
│   ├── docker-compose.yaml           # Traefik service definition
│   ├── traefik.yaml                  # Traefik static configuration
│   └── .env                          # Proxy-level environment variables
│
├── dev/                              # Workspace root
│   ├── workspace.yaml                # Workspace configuration (structure)
│   │
│   ├── overrides/                    # Manual overrides (optional, workspace-specific)
│   │   └── app-two.yaml              # Merged after generated config
│   │
│   ├── .generated/                   # Generated files (gitignored)
│   │   ├── state.yaml                # Runtime state (active flavors)
│   │   ├── manifest.yaml             # Computed values (read-only)
│   │   ├── app-one.override.yaml     # Generated compose override
│   │   ├── app-two.override.yaml
│   │   └── app-three.override.yaml
│   │
│   ├── app-one/                      # Cloned application repository
│   │   ├── docker-compose.yaml       # Application's compose file (app-owned)
│   │   ├── application.yaml          # Service contract + flavors (app-owned)
│   │   └── ...
│   ├── app-two/
│   │   ├── docker-compose.yaml
│   │   ├── docker-compose.worker.yaml
│   │   ├── docker-compose.extras.yaml
│   │   ├── application.yaml
│   │   └── ...
│   └── app-three/
│       ├── docker-compose.yaml
│       ├── application.yaml
│       └── ...
│
├── review/                           # Another workspace (same structure)
│   └── ...
│
└── control/                          # Another workspace
    └── ...
```

### Single-Application Workspace

When promoting an existing Docker Compose project:

```
~/my-project/                         # Workspace AND application directory
├── workspace.yaml                    # workspace.name = "dev"
├── application.yaml                  # Application service contract
├── docker-compose.yaml               # Existing compose file (unchanged)
├── docker-compose.worker.yaml        # Optional additional compose files
├── .generated/                       # Generated files (gitignored)
│   ├── state.yaml
│   ├── manifest.yaml
│   └── my-project.override.yaml
├── overrides/                        # Manual overrides (optional)
└── src/                              # Application source code
```

In this configuration, `workspace.yaml` references the application with `path: .`:

```yaml
workspace:
  name: dev
  applications:
    my-project:
      path: .                         # Application is in workspace root
```

---

## Configuration Schemas

The system uses three schema types, separating structure (configuration) from state (runtime):

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/contrail/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

### Proxy Configuration

**Location**: `~/.config/contrail/proxy.yaml` (global/per-user)

```yaml
proxy:
  domain: contrail.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image (defaults to pinned version)
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard (default: true)
    # Future: password support via environment variable
  # Future: TLS settings, entrypoints, etc.
```

### Proxy Infrastructure

**Location**: `~/.config/contrail/proxy/` (global/per-user)

The proxy is implemented as a Docker Compose project managed by Contrail. It runs a Traefik instance that handles reverse proxying for all workspaces on the host.

**Directory structure** (created by `contrail proxy init`):
```
~/.config/contrail/proxy/
├── docker-compose.yaml    # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
└── certs/                # TLS certificates (future)
```

**docker-compose.yaml** (generated):
```yaml
name: contrail-proxy

services:
  traefik:
    image: ${TRAEFIK_IMAGE:-traefik:v3.2.3}  # Configurable via proxy.yaml
    command:
      - "--configFile=/etc/traefik/traefik.yaml"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yaml:/etc/traefik/traefik.yaml:ro
      - ./dynamic:/etc/traefik/dynamic:ro
      - ./certs:/etc/traefik/certs:ro
    networks:
      - proxy
    restart: unless-stopped
    labels:
      - "contrail.managed=true"
      - "contrail.component=proxy"

networks:
  contrail-proxy:
    external: true
```

**traefik.yaml** (generated):
```yaml
api:
  dashboard: false

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    exposedByDefault: false
    network: contrail-proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

**Lifecycle**:
- `proxy init`: Creates the directory structure and configuration files
- `proxy up`: Starts the Traefik container (creates `contrail-proxy` network if needed)
- `proxy down`: Stops the Traefik container
- `workspace up`: Automatically runs `proxy up` if proxy is not running

**Recovery**: If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

### Workspace Registry

**Location**: `~/.config/contrail/workspaces.yaml` (global/per-user)

This file tracks all known workspaces across the system, enabling `workspace list` and preventing workspace name collisions.

```yaml
# AUTO-GENERATED - Managed by Contrail
# Records known workspaces and their locations

workspaces:
  dev:
    path: /home/user/workspaces/dev
    registered_at: 2024-12-28T10:30:00Z
    last_seen: 2024-12-29T14:22:00Z
  review:
    path: /home/user/workspaces/review
    registered_at: 2024-12-28T11:00:00Z
    last_seen: 2024-12-29T13:15:00Z
  myapp-dev:
    path: /home/user/projects/myapp
    registered_at: 2024-12-29T09:00:00Z
    last_seen: 2024-12-29T14:30:00Z
```

**Registry behavior**:
- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

**Docker label fallback**: If the registry file is missing or corrupted, Contrail can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels. This provides resilience against accidental deletion of `~/.config/contrail/workspaces.yaml`.

### Global State

**Location**: `~/.config/contrail/state.yaml` (global/per-user)

This file tracks port assignments for `assigned` type ports across all workspaces, plus an inventory of port availability for garbage collection and debugging.

```yaml
# AUTO-GENERATED - Managed by Contrail
# Records assigned ports and port availability inventory

assigned_ports:
  dev:                                  # Workspace name
    app-one:                            # Application name
      db: 5432                          # Exported service: assigned host port
    app-two:
      db: 5433                          # Incremented because 5432 was taken
      cache: 6379
  review:
    app-one:
      db: 5434                          # Different workspace, different port

port_inventory:
  5432:
    status: assigned                    # assigned | unavailable | released
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:                         # Present only if status=assigned
      workspace: dev
      application: app-one
      exported_service: db
  5433:
    status: assigned
    first_seen: 2025-12-28T17:53:58Z
    last_checked: 2025-12-29T13:01:37Z
    assignment:
      workspace: dev
      application: app-two
      exported_service: db
  5434:
    status: unavailable                 # External process using this port
    first_seen: 2025-12-28T17:54:00Z
    last_checked: 2025-12-29T13:01:40Z
    # No assignment - taken by external process
```

**Port assignment rules**:
1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)

**Port status transitions**:
- `unavailable` → `assigned`: Port became free, Contrail claimed it
- `assigned` → `released`: Workspace/app removed, port freed
- `unavailable` → `released`: External process stopped, `contrail port gc` cleaned it up

**Port availability checking**: `contrail port scan` and `contrail port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

### Workspace Configuration

**Location**: `{workspace}/workspace.yaml`

```yaml
workspace:
  name: dev                           # Required. Used as prefix for project names and hostnames
  # network: dev-custom               # Optional. Defaults to {name}-internal
  applications:
    app-one:
      repository: git@github.com:company/app-one.git  # Optional. For initial cloning
    app-two:
      repository: git@github.com:company/app-two.git
    app-three:
      repository: git@github.com:company/app-three.git
      path: ./custom-path             # Optional. Defaults to ./{app-name}
```

**Single-application workspace**:

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                         # Application in workspace root directory
```

**Conventions**:
- Application name = directory path (e.g., `app-one` → `./{workspace}/app-one/`)
- Network name defaults to `{workspace.name}-internal`

**Note**: Branch (`ref`) and flavor are runtime state, not configuration.

### Template Variables

Contrail uses template variables to generate hostnames, aliases, and other computed values. Templates can be customized at the workspace level.

**Default templates** (built-in):

```yaml
workspace:
  name: dev
  templates:
    # Public hostname pattern
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    
    # Internal network alias pattern
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    
    # Docker Compose project name
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

**Available variables**:

| Variable | Scope | Description | Example |
|----------|-------|-------------|---------|
| `%PROXY_DOMAIN%` | Proxy | Domain from `proxy.yaml` | `contrail.test` |
| `%WORKSPACE_NAME%` | Workspace | Workspace name | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | Internal network name | `dev-internal` |
| `%APPLICATION_NAME%` | Application | Application identifier | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | Resolved flavor | `default` |
| `%EXPORTED_SERVICE%` | Export | Key from `exported_services` | `web-debug` |
| `%SERVICE_NAME%` | Export | Underlying Compose service | `web` |
| `%SERVICE_PORT%` | Export | Container port number | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

### Template Resolution Timing

Template variables are resolved at **generation time** (when `workspace generate` or `workspace up` runs). The resolved values are written into the generated override files.

**Flavor changes**: When `contrail flavor set FLAVOR` is executed, it:
1. Updates `.generated/state.yaml` with the new flavor
2. Immediately regenerates the affected application's override file
3. If the application is currently running, displays a warning:
   ```
   Warning: Application "app-name" is currently running.
   The new flavor has been applied to the configuration, but running
   containers still use the previous flavor.

   To apply the flavor change:
     contrail app restart -a app-name
   ```

This ensures override files always reflect the current flavor without requiring a separate `generate` step.

**Running application considerations**: Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `contrail up` to start new services |
| Flavor removes services | Services still running but not in override | Run `contrail up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `contrail app restart` to pick up changes |

**Orphaned service handling**: When `contrail up` is run after a flavor change that removes services, Contrail passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

**Example scenario**:
1. User runs `contrail generate` with flavor "lite"
2. Override files are generated with "lite" values (e.g., `%APPLICATION_FLAVOR%` = "lite")
3. User runs `contrail flavor set full`
4. Contrail regenerates the override with "full" values immediately
5. User runs `contrail up` — override is already up-to-date

### Application Configuration (Service Contract)

**Location**: `{application}/application.yaml` (lives in the application's git repository)

This file defines the application's exported services and flavors. The application name is inferred from the directory name.

```yaml
default_flavor: full                    # Optional. Defaults to "default" if not specified

flavors:
  lite:
    compose_files:
      - docker-compose.yaml
  full:
    compose_files:
      - docker-compose.yaml
      - docker-compose.worker.yaml
      - docker-compose.extras.yaml

exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
      - type: proxied
        protocol: http
        visibility: protected
  api:
    ports:
      - type: proxied
        protocol: https
        visibility: public
  worker:
    ports:
      - type: assigned
        port: 9000
        visibility: protected
```

#### Exported Service Configuration

Each key in `exported_services` is the "exported service name" used for hostname generation, network aliases, and environment variables. By default, this key maps to a Compose service of the same name.

**Mapping to a different Compose service**: Use the `service:` property when the exported name differs from the Compose service name:

```yaml
exported_services:
  db:
    service: postgres                   # Maps to Compose service "postgres", exported as "db"
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

#### Port Configuration

Each exported service declares one or more ports:

```yaml
ports:
  - type: proxied                       # Required: proxied or assigned
    protocol: https                     # Required for proxied: http, https, or future SNI types
    port: 8080                          # Optional: container port (see inference rules below)
    visibility: public                  # Optional: public or protected (documentation only)
```

**Port type constraints**:
- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

**Port inference rules**:
- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

#### Examples

**Simple web service** (single port in Compose, inferred):

```yaml
# docker-compose.yaml
services:
  web:
    ports:
      - "8080"
```

```yaml
# application.yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
```

Result: HTTPS proxy to container port 8080. Environment variables will use proxy port 443.

**Database with direct port** (assigned port, auto-assigned if unavailable):

```yaml
# application.yaml
exported_services:
  db:
    service: mysql
    ports:
      - type: assigned
        port: 3306
        visibility: protected
```

**Service with both proxy and direct ports**:

```yaml
# application.yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        port: 443
        visibility: public
      - type: proxied
        protocol: http
        port: 80
        visibility: protected
      - type: assigned
        port: 9229
        visibility: protected           # Node.js debug port
```

---

## State Management

**Location**: `{workspace}/.generated/state.yaml` (gitignored)

Runtime state is tracked separately from configuration. State represents explicit choices made by the user (e.g., which flavor to use), not computed values.

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite                      # Overridden from default
```

**Flavor resolution order**:
1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

The manifest is a computed, read-only view of the workspace's current state. It captures all resolved values (hostnames, aliases, environment variables) derived from configuration and state. This serves several purposes:

- **Discoverability**: Humans and tools can inspect one file to understand the workspace topology
- **Tool integration**: Dashboards, DNS updaters, or service discovery tools can consume this structured data
- **Debugging**: Inspect computed hostnames and environment variables without reconstructing from templates
- **Caching**: Contrail can compare the manifest against configuration to determine if regeneration is needed

```yaml
# AUTO-GENERATED - Computed from configuration and state
# Generated: 2024-12-27T10:30:00Z

workspace:
  name: dev
  network: dev-internal

proxy:
  domain: contrail.test

applications:
  app-one:
    flavor: default
    project: dev-app-one
    exported_services:
      web:
        service: web                      # Underlying Compose service
        alias: app-one-web                # Internal network alias
        ports:
          - type: proxied
            protocol: https
            container_port: 443
            visibility: public
            hostname: dev-app-one-web.contrail.test
        environment:
          CONTRAIL_APP_ONE_WEB_HOST: dev-app-one-web.contrail.test
          CONTRAIL_APP_ONE_WEB_PORT: 443
          CONTRAIL_APP_ONE_WEB_SCHEME: https
          CONTRAIL_APP_ONE_WEB_URL: https://dev-app-one-web.contrail.test
          CONTRAIL_APP_ONE_WEB_HTTPS_HOST: dev-app-one-web.contrail.test
          CONTRAIL_APP_ONE_WEB_HTTPS_PORT: 443
          CONTRAIL_APP_ONE_WEB_HTTPS_URL: https://dev-app-one-web.contrail.test

      db:
        service: postgres
        alias: app-one-db
        ports:
          - type: assigned
            container_port: 5432
            host_port: 5432
            visibility: protected
        environment:
          CONTRAIL_APP_ONE_DB_HOST: app-one-db
          CONTRAIL_APP_ONE_DB_PORT: 5432
```

---

## Generated Override File

**Location**: `{workspace}/.generated/{application-name}.override.yaml`

These files are generated by the workspace tooling and should not be edited directly. They are gitignored.

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + app-one/application.yaml
# Flavor: default
# Generated: 2024-12-27T10:30:00Z

name: dev-app-one                         # Explicit project name to prevent conflicts

services:
  web:
    networks:
      dev-internal:
        aliases:
          - app-one-web
      contrail-proxy: {}                   # Connected to proxy for Traefik routing
    labels:
      # Traefik HTTPS router
      - "traefik.enable=true"
      - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.contrail.test`)"
      - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
      - "traefik.http.routers.dev-app-one-web-https.tls=true"
      - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"
      # Contrail context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/home/user/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/home/user/workspaces/dev/app-one"
      # Contrail export labels
      - "contrail.export.web.host=dev-app-one-web.contrail.test"
      - "contrail.export.web.proxy.https.visibility=public"
      - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"
    environment:
      - CONTRAIL_WORKSPACE_NAME=dev
      - CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
      - CONTRAIL_APP_ONE_WEB_PORT=443
      - CONTRAIL_APP_ONE_WEB_SCHEME=https
      - CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test
      # ... additional environment variables ...

  postgres:
    ports:
      - "5432:5432"                        # host:container - assigned port mapping
    networks:
      dev-internal:
        aliases:
          - app-one-db
    labels:
      # Contrail context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/home/user/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/home/user/workspaces/dev/app-one"
      # Contrail export labels
      - "contrail.export.db.host=dev-app-one-db.contrail.test"
      - "contrail.export.db.port.5432.visibility=protected"
      - "contrail.export.db.port.5432.assigned=5432"
    environment:
      - CONTRAIL_WORKSPACE_NAME=dev
      - CONTRAIL_APP_ONE_DB_HOST=app-one-db
      - CONTRAIL_APP_ONE_DB_PORT=5432
      # ... additional environment variables ...

networks:
  dev-internal:
    external: true
  contrail-proxy:
    external: true
```

### Manual Override File

**Location**: `{workspace}/overrides/{application-name}.yaml`

Optional. If present, merged after the generated override file. Useful for workspace-specific customizations that can't be expressed in the application config.

**Preservation guarantee**: Files in `{workspace}/overrides/` are **never modified by Contrail**. They persist across all regeneration operations, including `workspace generate --force`. Only the `.generated/` directory contents are affected by regeneration.

**Merge order**: Docker Compose files are merged in this order:
```
docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
```

This allows workspace-specific customizations (extra environment variables, volume mounts, middleware) that persist across regeneration.

```yaml
# Manual overrides for app-one in dev workspace
services:
  web:
    environment:
      - DEBUG=true
      - FEATURE_FLAG_X=enabled
    labels:
      - "traefik.http.routers.dev-app-one-web-https.middlewares=dev-auth@docker"

  postgres:
    volumes:
      - ./local-db-init:/docker-entrypoint-initdb.d:ro
```

---

## Docker Labels

Contrail uses Docker labels for workspace discovery, service routing, and external tool integration. All labels use the `contrail.` namespace prefix with kebab-case for multi-word segments.

### Context Labels

Applied to all application containers for workspace discovery and registry reconstruction:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.workspace.name` | Workspace identifier | `dev` |
| `contrail.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `contrail.app.name` | Application identifier | `app-one` |
| `contrail.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

### Export Labels

Applied to containers with exported services. Labels are keyed by export name for consistency, supporting multiple exports per container.

**Proxied exports** (HTTP/HTTPS through Traefik):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.proxy.http.visibility={public|protected|private}
contrail.export.{name}.proxy.http.url={url}
contrail.export.{name}.proxy.https.visibility={public|protected|private}
contrail.export.{name}.proxy.https.url={url}
```

**Assigned port exports** (direct port mapping):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.port.{internal-port}.visibility={public|protected|private}
contrail.export.{name}.port.{internal-port}.assigned={external-port}
```

**Example** — web service with proxied HTTP/HTTPS and debug port:
```yaml
labels:
  # Context
  - "contrail.workspace.name=dev"
  - "contrail.workspace.path=/Users/beau/workspaces/dev"
  - "contrail.app.name=app-one"
  - "contrail.app.path=/Users/beau/workspaces/dev/app-one"
  # Proxied export: web
  - "contrail.export.web.host=dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.http.visibility=public"
  - "contrail.export.web.proxy.http.url=http://dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.https.visibility=public"
  - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"
  # Assigned port export: debug
  - "contrail.export.debug.host=dev-app-one-debug.contrail.test"
  - "contrail.export.debug.port.9000.visibility=protected"
  - "contrail.export.debug.port.9000.assigned=9003"
```

### Proxy Container Labels

Applied to the Contrail-managed Traefik proxy container:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.managed` | Indicates Contrail manages this container | `true` |
| `contrail.component` | Component type | `proxy` |

### External Tool Integration

External tools can discover Contrail workspaces and services by querying Docker labels:

```bash
# Find all Contrail-managed containers
docker ps --filter "label=contrail.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=contrail.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "contrail.workspace.path"}}' $(docker ps -q --filter "label=contrail.workspace.name")
```

---

## Environment Variable Injection

All exported services receive environment variables for service discovery. This enables applications to reference other services without hardcoding hostnames.

### Naming Convention

Environment variables use a `CONTRAIL_` prefix to avoid conflicts with application-defined variables.

**Base variables** (always generated for each exported service):
```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

**Protocol-specific variables** (generated for each proxied protocol):
```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

### Variable Generation Rules

**For `proxied` type ports**:
- `*_HOST` contains the fully qualified proxied hostname (e.g., `dev-app-one-web.contrail.test`)
- `*_PORT` contains the proxy port (443 for HTTPS, 80 for HTTP)—**not** the container port
- `*_SCHEME` and `*_URL` are generated
- Protocol-specific variables (`*_HTTPS_*`, `*_HTTP_*`) are also generated

**For `assigned` type ports**:
- `*_HOST` contains the internal alias (e.g., `app-one-db`)
- `*_PORT` contains the assigned host port (which may differ from the requested port)
- No `*_SCHEME` or `*_URL` variables
- No protocol-specific variables

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | ✓ | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | ✓ | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | ✓ | Both |
| `assigned` | - | Internal alias | Assigned port | ✗ | ✗ | ✗ |

**HTTPS-default rationale**: When both HTTP and HTTPS are configured for an exported service, base variables (`*_PORT`, `*_SCHEME`, `*_URL`) default to HTTPS (port 443) following security-by-default principles. Applications should prefer HTTPS for service-to-service communication. Use protocol-specific variables (`*_HTTP_*`) when HTTP is explicitly required.

### Usage in Applications

```php
// PHP example - using URL directly for proxied services
$apiUrl = getenv('CONTRAIL_APP_TWO_API_URL') ?: 'https://app-two-api.contrail.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// PHP example - building connection for assigned port services
$dbHost = getenv('CONTRAIL_APP_ONE_DB_HOST') ?: 'app-one-db';
$dbPort = getenv('CONTRAIL_APP_ONE_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

```javascript
// Node.js example - using URL directly
const apiUrl = process.env.CONTRAIL_APP_TWO_API_URL || 'https://app-two-api.contrail.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Node.js example - building connection manually
const dbHost = process.env.CONTRAIL_APP_ONE_DB_HOST || 'app-one-db';
const dbPort = process.env.CONTRAIL_APP_ONE_DB_PORT || '5432';
```

---

## Proxy Layer

### Traefik Configuration

**proxy/docker-compose.yaml:**
```yaml
services:
  traefik:
    image: ${TRAEFIK_IMAGE:-traefik:v3.2.3}  # Configurable via proxy.yaml
    container_name: traefik
    restart: unless-stopped
    command:
      - "--api.dashboard=true"
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--providers.docker.network=contrail-proxy"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    networks:
      - contrail-proxy

networks:
  contrail-proxy:
    name: contrail-proxy
    driver: bridge
```

### DNS Configuration

For local development, configure DNS resolution for the workspace domains. Options include:

1. **dnsmasq**: Route all `*.contrail.test` to `127.0.0.1`
   ```
   address=/contrail.test/127.0.0.1
   ```
2. **/etc/hosts**: Manual entries for each hostname
3. **Local DNS server**: More complex but flexible

**Note**: The `.test` TLD is reserved by RFC 2606 for testing purposes and will not conflict with real domains or mDNS (unlike `.local`).

---

## CLI Interface

Contrail provides a comprehensive command-line interface for managing workspaces, applications, flavors, ports, and proxy configuration.

For complete CLI documentation, see **[Contrail CLI Reference](./contrail-cli-reference.md)**.

### Key Concepts

**Command Structure**: `contrail [resource] [action] [--options...]`

**Context Detection**: Contrail automatically detects workspace and application context from the current directory by walking up the tree looking for `workspace.yaml` and `application.yaml`.

**Options-Based Targeting**: All targeting uses `--workspace` / `-w` and `--app` / `-a` flags rather than positional arguments, making commands consistent and composable with context detection.

### Context Detection Algorithm

Context detection uses a **workspace boundary** approach to prevent accidental detection of config files in vendor packages or nested test fixtures.

**Detection steps**:
1. **Find workspace root**: Walk up from current working directory looking for `workspace.yaml`. The first one found establishes the workspace root.
2. **Find application context**: Walk up from current working directory toward the workspace root looking for `application.yaml`. Only consider `application.yaml` files that are within the workspace directory tree.
3. **Never traverse above workspace root** for application detection—this prevents vendor packages or nested fixtures from hijacking context.

**Error handling** (when no workspace found):
- If `application.yaml` found but no `workspace.yaml`:
  ```
  No workspace found (workspace.yaml) in current directory or any parent directories,
  but found an application (application.yaml) at: /path/to/application.yaml
  ```
- If neither found:
  ```
  No workspace found (workspace.yaml) in current directory or any parent directories,
  and no application (application.yaml) found either
  ```

**Edge cases**:
- **Nested vendor packages**: If working in `app-one/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored because the workspace's `app-one/application.yaml` is found first when walking toward the workspace root.
- **Workspace within workspace**: If a test fixture has its own `workspace.yaml` nested inside a workspace, the closest (innermost) `workspace.yaml` wins—this is typically the test fixture, which is the expected behavior.
- **Single-app workspaces**: With `path: .`, both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

### Quick Reference

```bash
# Workspace operations
contrail workspace init --workspace=dev
contrail workspace up [-w NAME]
contrail workspace down [-w NAME]
contrail workspace status [-w NAME]

# Application operations
contrail app add --app=NAME --repo=URL
contrail app up [-a NAME]

# Flavor management
contrail flavor set FLAVOR [-a NAME]

# Port management
contrail port list
contrail port gc

# Top-level aliases (with context detection)
contrail up
contrail down

# Docker Compose passthrough (shell function)
contrail-compose exec php bash
contrail-compose logs -f
contrail-compose -a app-two ps
```

---

## Operations

### Startup Sequence (`workspace up`)

1. Ensure proxy is running
2. Create workspace network if it doesn't exist
3. Check if override files are stale; regenerate if needed (see Staleness Detection)
4. For each application:
   - Resolve active flavor
   - Execute `docker compose up -d` with compose files + override

### Staleness Detection

Contrail uses **mtime comparison** to determine if generated override files need to be regenerated. Override files are considered stale if any of the following source files have a newer modification time than the generated override:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`, `docker-compose.worker.yaml`)

**Behavior**:
- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Touch a file accidentally? Use `--force` to ensure clean state

**Note**: mtime comparison is simple and fast but may trigger unnecessary regeneration if files are touched without content changes. The `--force` flag provides explicit control when needed.

### Generation Logic (`workspace generate`)

1. **Resolve flavor** for each application (CLI → state → default_flavor → "default")
2. **Get compose files** from resolved flavor's `compose_files` list
3. **Validate compose files exist** on disk; if any are missing, report error with available alternatives:
   ```
   Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
     Application: app-two
     Available compose files: docker-compose.yaml, docker-compose.dev.yaml
   ```
4. **Validate service references** in `exported_services` point to actual Compose services:
   ```
   Error: Exported service "api" references non-existent Compose service: backend
     Application: my-app
     Available services in docker-compose.yaml: web, db, redis
   ```
5. **Infer port values** for any exported services with omitted `port:` field (see Port Configuration)
6. **Default service names** for any exported services with omitted `service:` field
7. **Collect all exported services** across all applications in workspace
8. **Generate override file** with networks, aliases, labels, and environment variables
9. **Update state file** with resolved flavors
10. **Update manifest** with computed values

### Shutdown Sequence (`workspace down`)

1. For each application:
   - Execute `docker compose down`
2. Optionally remove workspace network
3. If `--volumes` specified, remove associated volumes

### Destroy Sequence (`workspace destroy`)

Completely removes a workspace and optionally its application directories:

1. Run `workspace down --volumes` to stop all containers and remove volumes
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force`)
4. Remove `workspace.yaml`
5. Release any assigned ports from global state
6. Remove workspace from registry (`~/.config/contrail/workspaces.yaml`)

**Flags**:
- `--force`: Skip confirmation prompts for application directory removal

### Viewing Logs

Using `contrail-compose` (recommended):
```bash
# All logs for an application (context-aware)
contrail-compose logs -f

# Specific service
contrail-compose logs -f web

# Different app from workspace root
contrail-compose -a app-two logs -f
```

Using raw Docker Compose:
```bash
# All logs for an application
docker compose -p dev-app-two logs -f

# Specific service
docker compose -p dev-app-two logs -f web

# All containers in a workspace (using labels)
docker logs $(docker ps -q --filter "label=contrail.workspace.name=dev")
```

### Listing Workspace Status

Using `contrail-compose`:
```bash
contrail-compose ps
contrail-compose -a app-two ps
```

Using raw Docker Compose:
```bash
# All containers in a workspace
docker ps --filter "label=contrail.workspace.name=dev"

# All containers for an application
docker ps --filter "label=contrail.app.name=app-two"
```

---

## Conventions and Best Practices

### Application Requirements

For an application to work well within a workspace, it should:

1. **Include an `application.yaml`**: This file defines the service contract - which services the application exports to the workspace. This is owned and maintained by the application developers.
2. **Use environment variables for external service URLs**: Don't hardcode hostnames for dependencies. Use the injected `CONTRAIL_{APP}_{EXPORTED_SERVICE}_*` variables.
3. **Expose ports without host bindings**: Use `ports: ["8080"]` not `ports: ["8080:8080"]` to avoid conflicts.
4. **Use relative volume paths**: Ensure builds and mounts work regardless of absolute path.

### The Service Contract (`application.yaml`)

The `application.yaml` file is the interface between the application and the workspace system:

- **Owned by**: Application developers (committed to the application's git repository)
- **Application name**: Inferred from the directory name (no explicit `name:` field required)
- **Purpose**: 
  - Declares which services the application exports and how they should be exposed
  - Defines flavors (different ways to run the application with different compose file combinations)
- **Consumed by**: Workspace tooling (to generate override files)

Application developers should update `application.yaml` when:
- Adding a new service that other applications need to access
- Changing the type (proxied, assigned) or protocol (http, https) of an exported service
- Renaming services that are exposed to the workspace
- Adding a new flavor (e.g., a "lite" mode that excludes optional services)
- Changing which compose files are needed for a flavor

### Naming Conventions

- **Workspace names**: Lowercase alphanumeric with hyphens (e.g., `dev`, `feature-x`, `pr-123`)
- **Application names**: Lowercase alphanumeric with hyphens, inferred from directory name
- **Exported service names**: The key in `exported_services`, may differ from the underlying Compose service name
- **Proxied hostnames** (proxied type): `{workspace}-{application}-{exported_service}.{domain}` (e.g., `dev-app-one-web.contrail.test`)
- **Internal aliases** (all types): `{application}-{exported_service}` (e.g., `app-one-web`, `app-one-db`)
- **Environment variables**: `CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}` in SCREAMING_SNAKE_CASE
- **Traefik router names**: `{workspace}-{application}-{exported_service}-{protocol}` (e.g., `dev-app-one-web-https`)

**Collision warning**: Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Creative naming that produces identical project names could cause conflicts:
- **Traefik routers**: Conflicting router names cause routing failures
- **Docker volumes**: Conflicting volume names (e.g., `dev-app-one_postgres_data`) could cause data to be shared unexpectedly or overwritten
- **Docker networks**: Conflicting network names could connect unrelated services

Example collision: workspace `dev-app` with app `one` and workspace `dev` with app `app-one` both produce project name `dev-app-one`.

Follow the lowercase-alphanumeric-with-hyphens convention and avoid names that could produce ambiguous concatenations.

### Git Strategy

**Workspace repository** (optional - can be version controlled):
- `workspace.yaml` - workspace definition
- `overrides/` - manual overrides

**Generated files** (gitignored):
- `.generated/` - generated override files

**Application directories** (cloned repositories):
- Each application is its own git repository
- `application.yaml` lives in the application repo and defines its service contract
- Can be managed as submodules, or simply cloned separately

Example workspace `.gitignore`:
```
.generated/
app-*/
```

---

## Future Enhancements

### Application-Level Dependencies

Allow specifying startup order between applications, similar to Docker Compose's `depends_on`:

```yaml
# workspace.yaml
applications:
  app-one:
    path: ./app-one
    
  app-two:
    path: ./app-two
    depends_on:
      - app-one
```

### Health Checks

Integration with Docker health checks to wait for dependent services.

### HTTPS Support

Automatic TLS certificate generation for local development:

```yaml
# workspace.yaml
workspace:
  name: dev
  tls:
    enabled: true
    provider: mkcert  # or: self-signed, acme
```

### Volume Sharing

Shared volumes between applications for scenarios like shared uploads:

```yaml
# workspace.yaml
volumes:
  shared-uploads:
    driver: local

applications:
  app-one:
    volumes:
      - shared-uploads:/app/uploads
  app-two:
    volumes:
      - shared-uploads:/app/uploads:ro
```

---

## Related Documentation

- **[Contrail PRD](./contrail-prd.md)**: Product requirements, concepts, and architectural decisions
- **[Contrail CLI Reference](./contrail-cli-reference.md)**: Complete CLI command documentation
- **[Contrail Shell Integration](./contrail-shell-integration.md)**: Shell functions, completion, and Docker Compose passthrough
- **[Contrail Go Stack](./contrail-go-stack.md)**: Go technology stack, dependencies, and project scaffolding

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0-draft | Dec 2024 | Initial technical specification |
| 0.2.0-draft | Dec 2024 | Port type system, global state, environment variables |
| 0.3.0-draft | Dec 2024 | Type/protocol split, assigned port binding, port inventory |
| 0.4.0-draft | Dec 2024 | Single-app workspace support, CLI redesign reference, extracted CLI documentation |
| 0.5.0-draft | Dec 2024 | Updated operations examples to show `contrail-compose` usage; linked shell integration specification |
| 0.5.1-draft | Dec 2024 | Spec review: validation rules, staleness detection, context detection algorithm, proxy init, Docker labels, flavor set behavior |
| 0.5.2-draft | Dec 2024 | Spec review: renamed proxy network to contrail-proxy, expanded collision warning, fixed Docker label prefixes |
