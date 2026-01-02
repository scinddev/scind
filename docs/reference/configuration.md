# Configuration Reference

Quick reference for all configuration files and options.

<!-- Source: specs/contrail-technical-spec.md, specs/contrail-cli-reference.md -->

---

## File Locations

| File | Location | Purpose |
|------|----------|---------|
| proxy.yaml | `~/.config/contrail/proxy.yaml` | Global proxy settings |
| state.yaml | `~/.config/contrail/state.yaml` | Global port assignments and inventory |
| workspaces.yaml | `~/.config/contrail/workspaces.yaml` | Workspace registry |
| workspace.yaml | `{workspace}/workspace.yaml` | Workspace definition |
| application.yaml | `{app}/application.yaml` | Application settings |

---

## proxy.yaml

Global Traefik and TLS configuration.

**Location**: `~/.config/contrail/proxy.yaml` (global/per-user)

```yaml
proxy:
  domain: contrail.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image (defaults to pinned version)
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard (default: true)
    port: 8080                           # Dashboard port (default: 8080)
    # Future: password support via environment variable
  tls:
    mode: auto                           # auto | custom | disabled
    # For mode: custom (e.g., enterprise CA certificates)
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `proxy.domain` | string | No | `contrail.test` | Base domain for public hostnames |
| `proxy.traefik_image` | string | No | `traefik:v3.2.3` | Traefik Docker image |
| `proxy.dashboard.enabled` | bool | No | `true` | Enable Traefik dashboard |
| `proxy.dashboard.port` | int | No | `8080` | Dashboard port |
| `proxy.tls.mode` | enum | No | `auto` | TLS mode: `auto`, `custom`, `disabled` |
| `proxy.tls.cert_file` | string | If custom | - | Path to certificate file |
| `proxy.tls.key_file` | string | If custom | - | Path to private key file |

### TLS Modes

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available to generate locally-trusted certificates; falls back to Traefik's default self-signed certificate (browser warnings) |
| `custom` | Uses user-provided certificate and key files (for enterprise CA or manually generated certs) |
| `disabled` | HTTP only, no HTTPS entrypoint (not recommended for production-like testing) |

### Certificate Setup by Mode

- **auto with mkcert**: Run `mkcert -install` once per machine to add the local CA to your trust store, then `mkcert "*.contrail.test"` to generate a wildcard certificate. Contrail will detect and use these automatically.
- **custom (enterprise CA)**: Obtain a wildcard certificate signed by your enterprise CA for `*.contrail.test` (or your configured domain). Place the cert and key files at the configured paths.
- **auto without mkcert**: Traefik serves its default self-signed certificate. Browsers will show security warnings.

---

## state.yaml (Global)

Global port assignments and inventory.

**Location**: `~/.config/contrail/state.yaml` (global/per-user)

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

### Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

### Port Status Values

| Status | Description |
|--------|-------------|
| `assigned` | Port is assigned to a Contrail service |
| `unavailable` | Port is in use by an external process |
| `released` | Port was previously assigned but is now free |

### Port Conflict at Startup

If a previously assigned port has become unavailable (e.g., taken by an external process) when `workspace up` runs, Contrail fails with a clear error:

```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

---

## workspaces.yaml

Workspace registry tracking all known workspaces.

**Location**: `~/.config/contrail/workspaces.yaml` (global/per-user)

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

### Registry Behavior

- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

**Docker label fallback**: If the registry file is missing or corrupted, Contrail can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels.

---

## workspace.yaml

Workspace definition and application references.

**Location**: `{workspace}/workspace.yaml`

```yaml
workspace:
  name: main                            # Workspace name (required)
  # network: main-custom                # Optional. Defaults to {name}-internal
  applications:
    frontend:
      repository: git@github.com:company/frontend.git  # Optional. For initial cloning
    backend:
      repository: git@github.com:company/backend.git
    shared-db:
      repository: git@github.com:company/shared-db.git
      path: ./custom-path               # Optional. Defaults to ./{app-name}
  templates:                            # Optional. Template customization
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `workspace.name` | string | Yes | - | Workspace name (used in hostnames) |
| `workspace.network` | string | No | `{name}-internal` | Custom network name |
| `workspace.applications` | object | Yes | - | Map of applications |
| `workspace.applications.<name>.repository` | string | No | - | Git repository URL for cloning |
| `workspace.applications.<name>.path` | string | No | `./{name}` | Path to application directory |
| `workspace.templates` | object | No | (defaults) | Template variable patterns |

### Single-Application Workspace

When promoting an existing Docker Compose project:

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                           # Application in workspace root directory
```

### Template Variables

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

---

## application.yaml

Application-specific configuration (service contract).

**Location**: `{application}/application.yaml` (lives in the application's git repository)

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
        container_port: 80
      - type: proxied
        protocol: http
        visibility: protected
        container_port: 80
  api:
    service: node                       # Map to different Compose service name
    ports:
      - type: proxied
        protocol: http
        visibility: protected
        container_port: 3000
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `default_flavor` | string | No | `default` | Default flavor to use |
| `flavors` | object | No | `{default: {}}` | Named configurations |
| `flavors.<name>.compose_files` | array | No | `[docker-compose.yaml]` | Docker Compose files for this flavor |
| `exported_services` | object | No | `{}` | Services to expose |
| `exported_services.<name>.service` | string | No | key name | Compose service name |
| `exported_services.<name>.ports` | array | Yes | - | Port configurations |

### Port Configuration

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | enum | Yes | - | `proxied` or `assigned` |
| `protocol` | string | If proxied | - | `http`, `https`, or future SNI types |
| `visibility` | enum | No | `protected` | `public` or `protected` |
| `container_port` | int | See rules | - | Container port |
| `port` | int | If assigned | - | Preferred host port |

### Port Inference Rules

- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `container_port:` or `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

### Port Type Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

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
- **assigned**: The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Contrail increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent:

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Contrail's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

**Docker label exposure**: Visibility is included in the generated Docker labels (`contrail.export.<name>.proxy.<protocol>.visibility`), enabling external tools to distinguish between public and protected services.

---

## Flavor Resolution

Flavor is resolved in this order:

1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Generated Files

### .generated/state.yaml

Runtime state tracking active flavors.

**Location**: `{workspace}/.generated/state.yaml` (gitignored)

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite                      # Overridden from default
```

---

### .generated/manifest.yaml

Computed, read-only view of workspace topology.

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

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

### .generated/*.override.yaml

Generated Docker Compose override files.

**Location**: `{workspace}/.generated/{application-name}.override.yaml`

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

---

## Manual Override Files

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

**Preservation guarantee**: Files in `{workspace}/overrides/` are **never modified by Contrail**. They persist across all regeneration operations, including `workspace generate --force`.

**Merge order**: Docker Compose files are merged in this order:
```
docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
```

---

## Environment Variable Injection

All exported services receive environment variables for service discovery.

### Naming Convention

Environment variables use a `CONTRAIL_` prefix to avoid conflicts. Hyphens in names are converted to underscores, and names are uppercased.

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

## Docker Labels

Contrail uses Docker labels for workspace discovery, service routing, and external tool integration.

### Context Labels

Applied to all application containers:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.workspace.name` | Workspace identifier | `dev` |
| `contrail.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `contrail.app.name` | Application identifier | `app-one` |
| `contrail.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

### Export Labels

**Proxied exports** (HTTP/HTTPS through Traefik):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.proxy.http.visibility={public|protected}
contrail.export.{name}.proxy.http.url={url}
contrail.export.{name}.proxy.https.visibility={public|protected}
contrail.export.{name}.proxy.https.url={url}
```

**Assigned port exports** (direct port mapping):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.port.{internal-port}.visibility={public|protected}
contrail.export.{name}.port.{internal-port}.assigned={external-port}
```

### Proxy Container Labels

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.managed` | Indicates Contrail manages this container | `true` |
| `contrail.component` | Component type | `proxy` |

### External Tool Integration

```bash
# Find all Contrail-managed containers
docker ps --filter "label=contrail.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=contrail.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "contrail.workspace.path"}}' \
  $(docker ps -q --filter "label=contrail.workspace.name")
```

---

## Directory Structure

### Standard Multi-Application Workspace

```
~/.config/contrail/
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

## Naming Conventions

- **Workspace names**: Lowercase alphanumeric with hyphens (e.g., `dev`, `feature-x`, `pr-123`)
- **Application names**: Lowercase alphanumeric with hyphens, inferred from directory name
- **Exported service names**: The key in `exported_services`, may differ from Compose service name
- **Proxied hostnames**: `{workspace}-{application}-{exported_service}.{domain}`
- **Internal aliases**: `{application}-{exported_service}` (e.g., `app-one-web`)
- **Environment variables**: `CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}` in SCREAMING_SNAKE_CASE
- **Traefik router names**: `{workspace}-{application}-{exported_service}-{protocol}`

**Collision warning**: Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Avoid names that could produce ambiguous concatenations.

---

## Related Documents

- [CLI Reference](./cli.md)
- [Configuration Schemas Spec](../specs/configuration-schemas.md)
- [Port Types Spec](../specs/port-types.md)
- [Docker Labels Spec](../specs/docker-labels.md)
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md)

<!-- See appendices/configuration/ for complete examples and JSON schemas -->
