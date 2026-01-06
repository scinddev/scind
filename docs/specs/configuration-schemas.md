# Configuration Schemas Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail uses three configuration schemas to define system behavior, separating structure (configuration) from state (runtime). The system uses a clear boundary between what is configured and what is computed or observed.

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/contrail/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

**Related Documents**:
- [ADR-0005: Structure vs State Separation](../decisions/0005-structure-vs-state-separation.md)
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md)
- [Architecture: Overview](../architecture/README.md)

**Appendices**:
- [Complete Examples](./appendices/configuration-schemas/complete-examples.md)

---

## Proxy Configuration

**Location**: `~/.config/contrail/proxy.yaml` (global/per-user)

This file configures the global Traefik reverse proxy settings that apply to all workspaces.

```yaml
proxy:
  domain: contrail.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image (defaults to pinned version)
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard (default: true)
    port: 8080                           # Dashboard port (default: 8080)
  tls:
    mode: auto                           # auto | custom | disabled
    # For mode: custom (e.g., enterprise CA certificates)
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `proxy.domain` | string | No | `contrail.test` | TLD for generated hostnames |
| `proxy.traefik_image` | string | No | `traefik:v3.2.3` | Traefik Docker image to use |
| `proxy.dashboard.enabled` | boolean | No | `true` | Enable Traefik dashboard |
| `proxy.dashboard.port` | integer | No | `8080` | Dashboard port |
| `proxy.tls.mode` | string | No | `auto` | TLS mode: auto, custom, or disabled |
| `proxy.tls.cert_file` | string | Conditional | - | Certificate file path (required if mode=custom) |
| `proxy.tls.key_file` | string | Conditional | - | Key file path (required if mode=custom) |

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

### Validation Rules

- `proxy.domain` must be a valid domain name (lowercase alphanumeric with dots)
- `proxy.dashboard.port` must be in range 1-65535
- `proxy.tls.mode` must be one of: `auto`, `custom`, `disabled`
- If `proxy.tls.mode` is `custom`, both `cert_file` and `key_file` are required
- Certificate and key files must exist and be readable when mode is `custom`

---

## Proxy Infrastructure

**Location**: `~/.config/contrail/proxy/` (global/per-user)

The proxy is implemented as a Docker Compose project managed by Contrail. It runs a Traefik instance that handles reverse proxying for all workspaces on the host.

### Directory Structure

Created by `contrail proxy init`:

```
~/.config/contrail/proxy/
├── docker-compose.yaml    # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

### docker-compose.yaml (generated)

```yaml
name: contrail-proxy

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
      - contrail-proxy
    restart: unless-stopped
    labels:
      - "contrail.managed=true"
      - "contrail.component=proxy"

networks:
  contrail-proxy:
    external: true
```

### traefik.yaml (generated)

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
    network: contrail-proxy
    watch: true
  file:
    directory: /etc/traefik/dynamic
    watch: true

log:
  level: INFO

accessLog: {}
```

### Lifecycle

- `proxy init`: Creates the directory structure and configuration files
- `proxy up`: Starts the Traefik container (creates `contrail-proxy` network if needed)
- `proxy down`: Stops the Traefik container
- `workspace up`: Automatically runs `proxy up` if proxy is not running

### Recovery

If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

---

## Workspace Registry

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

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `workspaces.{name}.path` | string | Yes | Absolute path to workspace directory |
| `workspaces.{name}.registered_at` | datetime | Yes | When workspace was first registered |
| `workspaces.{name}.last_seen` | datetime | Yes | When workspace was last accessed |

### Registry Behavior

- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

### Docker Label Fallback

If the registry file is missing or corrupted, Contrail can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels. This provides resilience against accidental deletion of `~/.config/contrail/workspaces.yaml`.

---

## Global State

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

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `assigned_ports.{workspace}.{app}.{service}` | integer | Assigned host port for the service |
| `port_inventory.{port}.status` | string | `assigned`, `unavailable`, or `released` |
| `port_inventory.{port}.first_seen` | datetime | When port was first tracked |
| `port_inventory.{port}.last_checked` | datetime | Last availability check time |
| `port_inventory.{port}.assignment.workspace` | string | Owning workspace (if assigned) |
| `port_inventory.{port}.assignment.application` | string | Owning application (if assigned) |
| `port_inventory.{port}.assignment.exported_service` | string | Owning service (if assigned) |

### Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

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

### Port Status Transitions

- `unavailable` → `assigned`: Port became free, Contrail claimed it
- `assigned` → `released`: Workspace/app removed, port freed
- `unavailable` → `released`: External process stopped, `contrail port gc` cleaned it up

### Port Availability Checking

`contrail port scan` and `contrail port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

---

## Workspace Configuration

**Location**: `{workspace}/workspace.yaml`

This file defines a workspace and its applications. The workspace name is used as a prefix for project names and hostnames.

```yaml
workspace:
  name: dev                           # Required. Used as prefix for project names and hostnames
  # network: dev-custom               # Optional. Defaults to {name}-internal
  templates:                          # Optional. Customizable naming patterns
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
  applications:
    app-one:
      repository: git@github.com:company/app-one.git  # Optional. For initial cloning
    app-two:
      repository: git@github.com:company/app-two.git
    app-three:
      repository: git@github.com:company/app-three.git
      path: ./custom-path             # Optional. Defaults to ./{app-name}
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `workspace.name` | string | Yes | - | Workspace identifier used for prefixes |
| `workspace.network` | string | No | `{name}-internal` | Internal network name |
| `workspace.templates.hostname` | string | No | See schema | Public hostname pattern |
| `workspace.templates.alias` | string | No | See schema | Internal network alias pattern |
| `workspace.templates.project-name` | string | No | See schema | Docker Compose project name pattern |
| `workspace.applications.{name}.repository` | string | No | - | Git repository URL for cloning |
| `workspace.applications.{name}.path` | string | No | `./{name}` | Path relative to workspace |

### Single-Application Workspace

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                         # Application in workspace root directory
```

### Conventions

- Application name = directory path (e.g., `app-one` → `./{workspace}/app-one/`)
- Network name defaults to `{workspace.name}-internal`

**Note**: Branch (`ref`) and flavor are runtime state, not configuration.

### Validation Rules

- `workspace.name` is required and must be lowercase alphanumeric with hyphens
- Application names must be lowercase alphanumeric with hyphens
- `path` must be a valid relative path within the workspace
- If `repository` is specified, it must be a valid git URL

---

## Template Variables

Contrail uses template variables to generate hostnames, aliases, and other computed values. Templates can be customized at the workspace level.

### Default Templates (built-in)

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
| `%PROXY_DOMAIN%` | Proxy | Domain from `proxy.yaml` | `contrail.test` |
| `%WORKSPACE_NAME%` | Workspace | Workspace name | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | Internal network name | `dev-internal` |
| `%APPLICATION_NAME%` | Application | Application identifier | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | Resolved flavor | `default` |
| `%EXPORTED_SERVICE%` | Export | Key from `exported_services` | `web-debug` |
| `%SERVICE_NAME%` | Export | Underlying Compose service | `web` |
| `%SERVICE_PORT%` | Export | Container port number (see note) | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

**Note on `%SERVICE_PORT%`**: This variable provides the container's internal port number. While not used in the default templates, it's available for advanced customization such as adding debugging labels or custom routing rules that need to reference the original container port (e.g., `contrail.debug.port=%SERVICE_PORT%`).

### Template Resolution Timing

Template variables are resolved at **generation time** (when `workspace generate` or `workspace up` runs). The resolved values are written into the generated override files.

### Flavor Changes

When `contrail flavor set FLAVOR` is executed, it:
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

### Running Application Considerations

Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `contrail up` to start new services |
| Flavor removes services | Services still running but not in override | Run `contrail up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `contrail app restart` to pick up changes |

### Orphaned Service Handling

When `contrail up` is run after a flavor change that removes services, Contrail passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

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

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `default_flavor` | string | No | `default` | Default flavor to use |
| `flavors.{name}.compose_files` | array | Yes | - | List of compose files for this flavor |
| `exported_services.{name}.service` | string | No | Same as key | Underlying Compose service name |
| `exported_services.{name}.ports` | array | Yes | - | Port configurations |
| `exported_services.{name}.ports[].type` | string | Yes | - | `proxied` or `assigned` |
| `exported_services.{name}.ports[].protocol` | string | Conditional | - | Required for proxied: `http`, `https` |
| `exported_services.{name}.ports[].port` | integer | Conditional | - | Container port (can be inferred) |
| `exported_services.{name}.ports[].visibility` | string | No | `protected` | `public` or `protected` |

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

### Port Type Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

### Port Inference Rules

- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

### Visibility

| Visibility | Description |
|------------|-------------|
| `public` | This port is intended for external/production use |
| `protected` (default) | This port exists for development/debugging but should not be depended on in production |

Visibility is primarily **documentation** to communicate intent to collaborators. It does not change Contrail's core behavior. Both public and protected proxied services route through Traefik. Visibility is included in generated Docker labels for external tool integration.

### Validation Rules

- Flavor names must be lowercase alphanumeric with hyphens
- All compose files referenced in flavors must exist
- Exported service keys must be lowercase alphanumeric with hyphens
- If `service` is specified, it must match a service in the referenced compose files
- Port type must be `proxied` or `assigned`
- Protocol is required when type is `proxied`
- Port number must be in range 1-65535

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

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `applications.{name}.flavor` | string | Currently selected flavor for the application |

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

## Examples

### Example 1: Minimal Proxy Configuration

**Input**:
```yaml
proxy:
  domain: local.test
```

**Behavior**: Uses default Traefik image, enables dashboard on port 8080, auto TLS mode.

**Result**: Proxy configured for `*.local.test` domains with automatic TLS.

### Example 2: Multi-Application Workspace

**Input** (workspace.yaml):
```yaml
workspace:
  name: dev
  applications:
    frontend:
      repository: git@github.com:org/frontend.git
    backend:
      repository: git@github.com:org/backend.git
    api:
      repository: git@github.com:org/api.git
```

**Behavior**: Creates workspace `dev` with three applications, using default naming patterns.

**Result**:
- Network: `dev-internal`
- Projects: `dev-frontend`, `dev-backend`, `dev-api`
- Hostnames: `dev-frontend-{service}.contrail.test`, etc.

### Example 3: Service with Multiple Port Types

**Input** (application.yaml):
```yaml
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

**Behavior**: Web service is accessible via HTTPS (public), HTTP (protected), and has a debug port assigned directly.

**Result**: Traefik routes for both HTTP and HTTPS, plus direct host port binding for debugging.

---

## Edge Cases

### Missing Default Flavor

**Scenario**: Application has no `default_flavor` and no flavor named `default`.

**Behavior**: Uses the first defined flavor alphabetically.

**Rationale**: Provides deterministic fallback behavior while encouraging explicit configuration.

### Single Port Compose Service

**Scenario**: Compose service has exactly one port defined, application.yaml omits `port:`.

**Behavior**: Port is inferred from the Compose service definition.

**Rationale**: Reduces boilerplate for simple services. See [ADR-0007](../decisions/0007-port-type-system.md) for port type system rationale.

### Custom Network Name

**Scenario**: Workspace specifies custom network name that conflicts with existing Docker network.

**Behavior**: Error on workspace creation with guidance to choose different name.

**Rationale**: Prevents accidental network collisions with external resources.

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Missing workspace.name | `workspace.name is required` | Add required field |
| Invalid TLS mode | `tls.mode must be one of: auto, custom, disabled` | Use valid mode value |
| Missing cert files | `Certificate file not found: {path}` | Create/move cert files or change TLS mode |
| Invalid compose file reference | `Flavor "{name}" references non-existent file: {file}` | Fix compose_files list |
| Invalid service reference | `Exported service "{name}" references non-existent Compose service: {service}` | Fix service reference |
| Workspace name collision | `Workspace "{name}" already registered at {path}` | Use different name or remove stale registration |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
