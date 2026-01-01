# Contrail Configuration Reference

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Configuration Files

| File | Location | Purpose |
|------|----------|---------|
| `proxy.yaml` | `~/.config/contrail/proxy.yaml` | Global proxy settings |
| `workspace.yaml` | `{workspace}/workspace.yaml` | Workspace definition |
| `application.yaml` | `{app}/application.yaml` | Application service contract |
| `state.yaml` | `~/.config/contrail/state.yaml` | Global port assignments |
| `workspaces.yaml` | `~/.config/contrail/workspaces.yaml` | Workspace registry |

---

## Proxy Configuration

**File**: `~/.config/contrail/proxy.yaml`

```yaml
proxy:
  domain: contrail.test                  # Domain for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image
  dashboard:
    enabled: true                        # Enable Traefik dashboard
    port: 8080                           # Dashboard port
  tls:
    mode: auto                           # auto | custom | disabled
    cert_file: ""                        # For mode: custom
    key_file: ""                         # For mode: custom
```

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `proxy.domain` | string | `contrail.test` | TLD for generated hostnames |
| `proxy.traefik_image` | string | `traefik:v3.2.3` | Traefik Docker image |
| `proxy.dashboard.enabled` | bool | `true` | Enable Traefik dashboard |
| `proxy.dashboard.port` | int | `8080` | Dashboard port |
| `proxy.tls.mode` | string | `auto` | TLS mode: `auto`, `custom`, `disabled` |
| `proxy.tls.cert_file` | string | | Path to certificate file (for `custom` mode) |
| `proxy.tls.key_file` | string | | Path to key file (for `custom` mode) |

---

## Workspace Configuration

**File**: `{workspace}/workspace.yaml`

```yaml
workspace:
  name: dev                              # Workspace identifier
  applications:
    app-one:
      repository: git@github.com:org/app-one.git
    app-two:
      repository: git@github.com:org/app-two.git
      path: ./custom-path                # Custom relative path
  templates:                             # Optional template overrides
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

| Key | Type | Required | Description |
|-----|------|----------|-------------|
| `workspace.name` | string | Yes | Workspace identifier (used in project names) |
| `workspace.applications` | map | No | Map of application names to config |
| `workspace.applications.*.repository` | string | No | Git repository URL for cloning |
| `workspace.applications.*.path` | string | No | Custom path relative to workspace (default: `./{app-name}`) |
| `workspace.templates.hostname` | string | No | Hostname template |
| `workspace.templates.alias` | string | No | Internal alias template |
| `workspace.templates.project-name` | string | No | Docker Compose project name template |

---

## Application Configuration

**File**: `{app}/application.yaml`

```yaml
default_flavor: full                     # Default flavor to use

flavors:
  lite:
    compose_files:
      - docker-compose.yaml
  full:
    compose_files:
      - docker-compose.yaml
      - docker-compose.worker.yaml

exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
  db:
    service: postgres                    # Maps to Compose service "postgres"
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

| Key | Type | Required | Description |
|-----|------|----------|-------------|
| `default_flavor` | string | No | Default flavor name (default: `default`) |
| `flavors` | map | No | Map of flavor names to config |
| `flavors.*.compose_files` | array | Yes | List of compose files for this flavor |
| `exported_services` | map | Yes | Map of exported service names to config |
| `exported_services.*.service` | string | No | Compose service name (default: map key) |
| `exported_services.*.ports` | array | Yes | List of port configurations |
| `exported_services.*.ports[].type` | string | Yes | `proxied` or `assigned` |
| `exported_services.*.ports[].protocol` | string | Conditional | Required for `proxied`: `http`, `https` |
| `exported_services.*.ports[].port` | int | No | Container port (inferred from Compose if single port) |
| `exported_services.*.ports[].visibility` | string | No | `public` or `protected` (default: `protected`) |

---

## Template Variables

| Variable | Scope | Description | Example |
|----------|-------|-------------|---------|
| `%PROXY_DOMAIN%` | Proxy | Domain from `proxy.yaml` | `contrail.test` |
| `%WORKSPACE_NAME%` | Workspace | Workspace name | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | Internal network name | `dev-internal` |
| `%APPLICATION_NAME%` | Application | Application identifier | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | Resolved flavor | `default` |
| `%EXPORTED_SERVICE%` | Export | Key from `exported_services` | `web` |
| `%SERVICE_NAME%` | Export | Underlying Compose service | `web` |
| `%SERVICE_PORT%` | Export | Container port number | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

**Note on `%SERVICE_PORT%`**: This variable provides the container's internal port number. While not used in the default templates, it's available for advanced customization such as adding debugging labels or custom routing rules that need to reference the original container port.

**Syntax note**: The PRD and conceptual documentation may use `{placeholder}` syntax for readability, but actual configuration uses `%VARIABLE_NAME%` syntax.

---

## Environment Variables

Contrail respects these environment variables:

| Variable | Description |
|----------|-------------|
| `CONTRAIL_CONFIG` | Path to config file |
| `CONTRAIL_WORKSPACE` | Default workspace name |
| `CONTRAIL_APP` | Default application name |
| `CONTRAIL_DOMAIN` | Override proxy domain |

---

## Related Documents

- [Configuration Schemas Specification](../specs/configuration-schemas.md)
- [Port Types and Proxying](../specs/port-types.md)
- [Environment Variable Injection](../specs/environment-variables.md)
