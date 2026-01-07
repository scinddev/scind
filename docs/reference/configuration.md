# Configuration Reference

**Status**: Draft

This document is the authoritative reference for all Scind configuration files, their schemas, and their locations.

---

## Overview

Scind uses three schema types, separating structure (configuration) from state (runtime):

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/scind/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

---

## Configuration Files

### File Locations Summary

| File | Location | Scope | Purpose |
|------|----------|-------|---------|
| `proxy.yaml` | `~/.config/scind/proxy.yaml` | Global | Proxy domain, TLS, dashboard settings |
| `state.yaml` | `~/.config/scind/state.yaml` | Global | Port assignments, port inventory |
| `workspaces.yaml` | `~/.config/scind/workspaces.yaml` | Global | Workspace registry |
| `workspace.yaml` | `{workspace}/workspace.yaml` | Per-workspace | Workspace definition, applications |
| `application.yaml` | `{app}/application.yaml` | Per-application | Flavors, exported services |
| `state.yaml` | `{workspace}/.generated/state.yaml` | Per-workspace | Active flavors (runtime) |
| `manifest.yaml` | `{workspace}/.generated/manifest.yaml` | Per-workspace | Computed values (read-only) |

---

## Proxy Configuration

**Location**: `~/.config/scind/proxy.yaml` (global/per-user)

```yaml
proxy:
  domain: scind.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image (defaults to pinned version)
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard (default: true)
    port: 8080                           # Dashboard port (default: 8080)
    # Future: password support via environment variable
  tls:
    mode: auto                           # auto | custom | disabled
    # For mode: custom (e.g., enterprise CA certificates)
    cert_file: ~/.config/scind/certs/wildcard.crt
    key_file: ~/.config/scind/certs/wildcard.key
```

### TLS Modes

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available to generate locally-trusted certificates; falls back to Traefik's default self-signed certificate (browser warnings) |
| `custom` | Uses user-provided certificate and key files (for enterprise CA or manually generated certs) |
| `disabled` | HTTP only, no HTTPS entrypoint (not recommended for production-like testing) |

### Certificate Setup by Mode

- **auto with mkcert**: Run `mkcert -install` once per machine to add the local CA to your trust store, then `mkcert "*.scind.test"` to generate a wildcard certificate. Scind will detect and use these automatically.
- **custom (enterprise CA)**: Obtain a wildcard certificate signed by your enterprise CA for `*.scind.test` (or your configured domain). Place the cert and key files at the configured paths.
- **auto without mkcert**: Traefik serves its default self-signed certificate. Browsers will show security warnings.

---

## Proxy Infrastructure

**Location**: `~/.config/scind/proxy/` (global/per-user)

The proxy is implemented as a Docker Compose project managed by Scind. It runs a Traefik instance that handles reverse proxying for all workspaces on the host.

### Directory Structure

Created by `scind proxy init`:

```
~/.config/scind/proxy/
├── docker-compose.yaml   # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

### Generated docker-compose.yaml

```yaml
name: scind-proxy

services:
  traefik:
    image: ${TRAEFIK_IMAGE:-traefik:v3.2.3}  # Configurable via proxy.yaml
    command:
      - "--configFile=/etc/traefik/traefik.yaml"
      - "--api.dashboard=true"               # Set to false if dashboard.enabled: false
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"                          # Only included if dashboard.enabled: true (port from dashboard.port)
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yaml:/etc/traefik/traefik.yaml:ro
      - ./dynamic:/etc/traefik/dynamic:ro
      - ./certs:/etc/traefik/certs:ro
    networks:
      - scind-proxy
    restart: unless-stopped
    labels:
      - "scind.managed=true"
      - "scind.component=proxy"

networks:
  scind-proxy:
    external: true
```

### Generated traefik.yaml

```yaml
api:
  dashboard: true                          # Set based on proxy.yaml dashboard.enabled (default: true)

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    exposedByDefault: false
    network: scind-proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

### Proxy Lifecycle

- `proxy init`: Creates the directory structure and configuration files
- `proxy up`: Starts the Traefik container (creates `scind-proxy` network if needed)
- `proxy down`: Stops the Traefik container
- `workspace up`: Automatically runs `proxy up` if proxy is not running

**Recovery**: If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

---

## Workspace Registry

**Location**: `~/.config/scind/workspaces.yaml` (global/per-user)

This file tracks all known workspaces across the system, enabling `workspace list` and preventing workspace name collisions.

```yaml
# AUTO-GENERATED - Managed by Scind
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

### Registry Behavior

- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

**Docker label fallback**: If the registry file is missing or corrupted, Scind can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels. This provides resilience against accidental deletion of `~/.config/scind/workspaces.yaml`.

---

## Global State

**Location**: `~/.config/scind/state.yaml` (global/per-user)

This file tracks port assignments for `assigned` type ports across all workspaces, plus an inventory of port availability for garbage collection and debugging.

```yaml
# AUTO-GENERATED - Managed by Scind
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

### Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

### Port Conflict at Startup

If a previously assigned port has become unavailable (e.g., taken by an external process) when `workspace up` runs, Scind fails with a clear error:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  scind port scan       # Check which ports are conflicting
  scind port release 5432   # Release the conflicting assignment
  scind generate --force    # Regenerate with new port assignment
```

### Port Status Transitions

- `unavailable` -> `assigned`: Port became free, Scind claimed it
- `assigned` -> `released`: Workspace/app removed, port freed
- `unavailable` -> `released`: External process stopped, `scind port gc` cleaned it up

### Port Availability Checking

`scind port scan` and `scind port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

---

## Workspace Configuration

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

### Single-Application Workspace

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                         # Application in workspace root directory
```

### Conventions

- Application name = directory path (e.g., `app-one` -> `./{workspace}/app-one/`)
- Network name defaults to `{workspace.name}-internal`

**Note**: Branch (`ref`) and flavor are runtime state, not configuration.

---

## Template Variables

Scind uses template variables to generate hostnames, aliases, and other computed values. Templates can be customized at the workspace level.

### Default Templates

Built-in defaults:

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

### Available Variables

| Variable | Scope | Description | Example |
|----------|-------|-------------|---------|
| `%PROXY_DOMAIN%` | Proxy | Domain from `proxy.yaml` | `scind.test` |
| `%WORKSPACE_NAME%` | Workspace | Workspace name | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | Internal network name | `dev-internal` |
| `%APPLICATION_NAME%` | Application | Application identifier | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | Resolved flavor | `default` |
| `%EXPORTED_SERVICE%` | Export | Key from `exported_services` | `web-debug` |
| `%SERVICE_NAME%` | Export | Underlying Compose service | `web` |
| `%SERVICE_PORT%` | Export | Container port number (see note) | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

**Note on `%SERVICE_PORT%`**: This variable provides the container's internal port number. While not used in the default templates, it's available for advanced customization such as adding debugging labels or custom routing rules that need to reference the original container port (e.g., `scind.debug.port=%SERVICE_PORT%`).

### Template Resolution Timing

Template variables are resolved at **generation time** (when `workspace generate` or `workspace up` runs). The resolved values are written into the generated override files.

**Flavor changes**: When `scind flavor set FLAVOR` is executed, it:
1. Updates `.generated/state.yaml` with the new flavor
2. Immediately regenerates the affected application's override file
3. If the application is currently running, displays a warning:
   ```
   Warning: Application "app-name" is currently running.
   The new flavor has been applied to the configuration, but running
   containers still use the previous flavor.

   To apply the flavor change:
     scind app restart -a app-name
   ```

This ensures override files always reflect the current flavor without requiring a separate `generate` step.

### Running Application Considerations

Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `scind up` to start new services |
| Flavor removes services | Services still running but not in override | Run `scind up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `scind app restart` to pick up changes |

**Orphaned service handling**: When `scind up` is run after a flavor change that removes services, Scind passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

---

## Application Configuration (Service Contract)

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

### Exported Service Configuration

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

### Port Configuration

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

### Application Configuration Examples

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

For complete configuration examples, see [appendices/configuration/complete-examples.md](./appendices/configuration/complete-examples.md).

---

## Workspace State

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

### Flavor Resolution Order

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
- **Caching**: Scind can compare the manifest against configuration to determine if regeneration is needed

```yaml
# AUTO-GENERATED - Computed from configuration and state
# Generated: 2024-12-27T10:30:00Z

workspace:
  name: dev
  network: dev-internal

proxy:
  domain: scind.test

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
            hostname: dev-app-one-web.scind.test
        environment:
          SCIND_APP_ONE_WEB_HOST: dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_PORT: 443
          SCIND_APP_ONE_WEB_SCHEME: https
          SCIND_APP_ONE_WEB_URL: https://dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_HTTPS_HOST: dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_HTTPS_PORT: 443
          SCIND_APP_ONE_WEB_HTTPS_URL: https://dev-app-one-web.scind.test

      db:
        service: postgres
        alias: app-one-db
        ports:
          - type: assigned
            container_port: 5432
            host_port: 5432
            visibility: protected
        environment:
          SCIND_APP_ONE_DB_HOST: app-one-db
          SCIND_APP_ONE_DB_PORT: 5432
```

---

## Generated Override Files

**Location**: `{workspace}/.generated/{application-name}.override.yaml`

Generated Docker Compose override files wire applications into the Scind infrastructure. These files are regenerated when configuration changes.

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
      scind-proxy: {}                   # Connected to proxy for Traefik routing
    labels:
      # Traefik HTTPS router
      - "traefik.enable=true"
      - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.scind.test`)"
      - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
      - "traefik.http.routers.dev-app-one-web-https.tls=true"
      - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"
      # Scind context labels
      - "scind.workspace.name=dev"
      - "scind.workspace.path=/home/user/workspaces/dev"
      - "scind.app.name=app-one"
      - "scind.app.path=/home/user/workspaces/dev/app-one"
      # Scind export labels
      - "scind.export.web.host=dev-app-one-web.scind.test"
      - "scind.export.web.proxy.https.visibility=public"
      - "scind.export.web.proxy.https.url=https://dev-app-one-web.scind.test"
    environment:
      - SCIND_WORKSPACE_NAME=dev
      - SCIND_APP_ONE_WEB_HOST=dev-app-one-web.scind.test
      - SCIND_APP_ONE_WEB_PORT=443
      - SCIND_APP_ONE_WEB_SCHEME=https
      - SCIND_APP_ONE_WEB_URL=https://dev-app-one-web.scind.test
      # ... additional environment variables ...

  postgres:
    ports:
      - "5432:5432"                        # host:container - assigned port mapping
    networks:
      dev-internal:
        aliases:
          - app-one-db
    labels:
      # Scind context labels
      - "scind.workspace.name=dev"
      - "scind.workspace.path=/home/user/workspaces/dev"
      - "scind.app.name=app-one"
      - "scind.app.path=/home/user/workspaces/dev/app-one"
      # Scind export labels
      - "scind.export.db.host=dev-app-one-db.scind.test"
      - "scind.export.db.port.5432.visibility=protected"
      - "scind.export.db.port.5432.assigned=5432"
    environment:
      - SCIND_WORKSPACE_NAME=dev
      - SCIND_APP_ONE_DB_HOST=app-one-db
      - SCIND_APP_ONE_DB_PORT=5432
      # ... additional environment variables ...

networks:
  dev-internal:
    external: true
  scind-proxy:
    external: true
```

### Manual Override Files

**Location**: `{workspace}/overrides/{application-name}.yaml`

Optional files merged after generated overrides for workspace-specific customizations.

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

**Preservation guarantee**: Files in `{workspace}/overrides/` are **never modified by Scind**. They persist across all regeneration operations, including `workspace generate --force`.

**Merge order**: Docker Compose files are merged in this order:
```
docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
```

---

## Port Types and Proxying

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `*_HOST`, `*_PORT` |

### Type Descriptions

- **proxied**: Traffic is routed through Traefik. The exported service gets a hostname (`{workspace}-{app}-{export}.{domain}`) and Traefik labels are generated. Environment variables contain the **proxy values** (hostname and proxy port 80/443), not the container port.
- **assigned**: The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Scind increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

---

## Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent:

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Scind's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

**Docker label exposure**: Visibility is included in the generated Docker labels (`scind.export.<name>.proxy.<protocol>.visibility`), enabling external tools to distinguish between public and protected services.

---

## Environment Variable Injection

All exported services receive environment variables for service discovery.

### Naming Convention

Environment variables use a `SCIND_` prefix to avoid conflicts. Hyphens in names are converted to underscores, and names are uppercased.

**Base variables** (always generated for each exported service):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

**Protocol-specific variables** (generated for each proxied protocol):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

### Variable Generation Rules

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | Yes | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | Yes | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | Yes | Both |
| `assigned` | - | Internal alias | Assigned port | No | No | No |

**HTTPS-default rationale**: When both HTTP and HTTPS are configured, base variables default to HTTPS (port 443) following security-by-default principles.

### Usage Examples

```php
// PHP example - using URL directly for proxied services
$apiUrl = getenv('SCIND_APP_TWO_API_URL') ?: 'https://app-two-api.scind.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// PHP example - building connection for assigned port services
$dbHost = getenv('SCIND_APP_ONE_DB_HOST') ?: 'app-one-db';
$dbPort = getenv('SCIND_APP_ONE_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

```javascript
// Node.js example - using URL directly
const apiUrl = process.env.SCIND_APP_TWO_API_URL || 'https://app-two-api.scind.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Node.js example - building connection manually
const dbHost = process.env.SCIND_APP_ONE_DB_HOST || 'app-one-db';
const dbPort = process.env.SCIND_APP_ONE_DB_PORT || '5432';
```

---

## Docker Labels

Scind uses Docker labels for workspace discovery, service routing, and external tool integration.

### Context Labels

Applied to all application containers:

| Label | Description | Example |
|-------|-------------|---------|
| `scind.workspace.name` | Workspace identifier | `dev` |
| `scind.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `scind.app.name` | Application identifier | `app-one` |
| `scind.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

### Export Labels

**Proxied exports** (HTTP/HTTPS through Traefik):
```
scind.export.{name}.host={hostname}
scind.export.{name}.proxy.http.visibility={public|protected}
scind.export.{name}.proxy.http.url={url}
scind.export.{name}.proxy.https.visibility={public|protected}
scind.export.{name}.proxy.https.url={url}
```

**Assigned port exports** (direct port mapping):
```
scind.export.{name}.host={hostname}
scind.export.{name}.port.{internal-port}.visibility={public|protected}
scind.export.{name}.port.{internal-port}.assigned={external-port}
```

### Proxy Container Labels

| Label | Description | Example |
|-------|-------------|---------|
| `scind.managed` | Indicates Scind manages this container | `true` |
| `scind.component` | Component type | `proxy` |

### External Tool Integration

```bash
# Find all Scind-managed containers
docker ps --filter "label=scind.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=scind.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "scind.workspace.path"}}' \
  $(docker ps -q --filter "label=scind.workspace.name")
```

---

## Directory Structure

### Standard Multi-Application Workspace

```
~/.config/scind/
├── proxy.yaml                        # Global proxy configuration
├── state.yaml                        # Global port assignments and inventory
└── workspaces.yaml                   # Workspace registry

workspaces/
├── proxy/
│   ├── docker-compose.yaml           # Traefik service definition
│   ├── traefik.yaml                  # Traefik static configuration
│   └── .env                          # Proxy-level environment variables
│
├── dev/                              # Workspace root
│   ├── workspace.yaml                # Workspace configuration (structure)
│   │
│   ├── overrides/                    # Manual overrides (optional)
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
│   │   ├── docker-compose.yaml       # Application's compose file
│   │   ├── application.yaml          # Service contract + flavors
│   │   └── ...
│   ├── app-two/
│   │   ├── docker-compose.yaml
│   │   ├── docker-compose.worker.yaml
│   │   ├── application.yaml
│   │   └── ...
│   └── app-three/
│       └── ...
│
├── review/                           # Another workspace
│   └── ...
│
└── control/                          # Another workspace
    └── ...
```

### Single-Application Workspace

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

## Naming Conventions

- **Workspace names**: Lowercase alphanumeric with hyphens (e.g., `dev`, `feature-x`, `pr-123`)
- **Application names**: Lowercase alphanumeric with hyphens, inferred from directory name
- **Exported service names**: The key in `exported_services`, may differ from Compose service name
- **Proxied hostnames**: `{workspace}-{application}-{exported_service}.{domain}`
- **Internal aliases**: `{application}-{exported_service}` (e.g., `app-one-web`)
- **Environment variables**: `SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}` in SCREAMING_SNAKE_CASE
- **Traefik router names**: `{workspace}-{application}-{exported_service}-{protocol}`

**Collision warning**: Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Avoid names that could produce ambiguous concatenations.

---

## Related Documentation

- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Detailed schema specification
- [Port Types Spec](../specs/port-types.md) - Port type behaviors
- [Docker Labels Spec](../specs/docker-labels.md) - Label conventions
- [CLI Reference](./cli.md) - Commands that use these configurations
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md) - Design rationale

---

## Source Attribution

<!-- Source: specs/scind-technical-spec.md:216-765 -->
