<!-- Migrated from specs/scind-technical-spec.md:216-683 -->
<!-- Extraction ID: spec-configuration-schemas -->

## Configuration Schemas

The system uses three schema types, separating structure (configuration) from state (runtime):

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/scind/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

### Proxy Configuration

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

**TLS Modes**:

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available to generate locally-trusted certificates; falls back to Traefik's default self-signed certificate (browser warnings) |
| `custom` | Uses user-provided certificate and key files (for enterprise CA or manually generated certs) |
| `disabled` | HTTP only, no HTTPS entrypoint (not recommended for production-like testing) |

**Certificate Setup by Mode**:

- **auto with mkcert**: Run `mkcert -install` once per machine to add the local CA to your trust store, then `mkcert "*.scind.test"` to generate a wildcard certificate. Scind will detect and use these automatically.
- **custom (enterprise CA)**: Obtain a wildcard certificate signed by your enterprise CA for `*.scind.test` (or your configured domain). Place the cert and key files at the configured paths.
- **auto without mkcert**: Traefik serves its default self-signed certificate. Browsers will show security warnings.

### Proxy Infrastructure

**Location**: `~/.config/scind/proxy/` (global/per-user)

The proxy is implemented as a Docker Compose project managed by Scind. It runs a Traefik instance that handles reverse proxying for all workspaces on the host.

**Directory structure** (created by `scind proxy init`):
```
~/.config/scind/proxy/
├── docker-compose.yaml   # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

**docker-compose.yaml** (generated):
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

**traefik.yaml** (generated):
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

**Lifecycle**:
- `proxy init`: Creates the directory structure and configuration files
- `proxy up`: Starts the Traefik container (creates `scind-proxy` network if needed)
- `proxy down`: Stops the Traefik container
- `workspace up`: Automatically runs `proxy up` if proxy is not running

**Recovery**: If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

### Workspace Registry

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

**Registry behavior**:
- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

**Docker label fallback**: If the registry file is missing or corrupted, Scind can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels. This provides resilience against accidental deletion of `~/.config/scind/workspaces.yaml`.

### Global State

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

**Port assignment rules**:
1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

**Port conflict at startup**: If a previously assigned port has become unavailable (e.g., taken by an external process) when `workspace up` runs, Scind fails with a clear error:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  scind port scan       # Check which ports are conflicting
  scind port release 5432   # Release the conflicting assignment
  scind generate --force    # Regenerate with new port assignment
```

**Port status transitions**:
- `unavailable` → `assigned`: Port became free, Scind claimed it
- `assigned` → `released`: Workspace/app removed, port freed
- `unavailable` → `released`: External process stopped, `scind port gc` cleaned it up

**Port availability checking**: `scind port scan` and `scind port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

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

Scind uses template variables to generate hostnames, aliases, and other computed values. Templates can be customized at the workspace level.

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

**Running application considerations**: Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `scind up` to start new services |
| Flavor removes services | Services still running but not in override | Run `scind up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `scind app restart` to pick up changes |

**Orphaned service handling**: When `scind up` is run after a flavor change that removes services, Scind passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

**Example scenario**:
1. User runs `scind generate` with flavor "lite"
2. Override files are generated with "lite" values (e.g., `%APPLICATION_FLAVOR%` = "lite")
3. User runs `scind flavor set full`
4. Scind regenerates the override with "full" values immediately
5. User runs `scind up` — override is already up-to-date

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
