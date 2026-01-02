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
- [Architecture: Overview](../architecture/overview.md)

**Appendices**:
- [Complete Examples](./appendices/configuration-schemas/complete-examples.md)

---

## proxy.yaml Schema

**Location**: `~/.config/contrail/proxy.yaml` (global/per-user)

This file configures the global Traefik reverse proxy settings that apply to all workspaces.

### Schema

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

## workspace.yaml Schema

**Location**: `{workspace}/workspace.yaml`

This file defines a workspace and its applications. The workspace name is used as a prefix for project names and hostnames.

### Schema

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
| `workspace.templates.hostname` | string | No | See above | Public hostname pattern |
| `workspace.templates.alias` | string | No | See above | Internal network alias pattern |
| `workspace.templates.project-name` | string | No | See above | Docker Compose project name pattern |
| `workspace.applications.{name}.repository` | string | No | - | Git repository URL for cloning |
| `workspace.applications.{name}.path` | string | No | `./{name}` | Path relative to workspace |

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

### Single-Application Workspace

When promoting an existing Docker Compose project, `workspace.yaml` can reference the application with `path: .`:

```yaml
workspace:
  name: dev
  applications:
    my-project:
      path: .                         # Application is in workspace root
```

### Validation Rules

- `workspace.name` is required and must be lowercase alphanumeric with hyphens
- Application names must be lowercase alphanumeric with hyphens
- `path` must be a valid relative path within the workspace
- If `repository` is specified, it must be a valid git URL

---

## application.yaml Schema

**Location**: `{application}/application.yaml` (lives in the application's git repository)

This file defines the application's service contract - which services it exports and how they should be exposed. The application name is inferred from the directory name.

### Schema

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
  db:
    service: postgres                   # Maps to Compose service "postgres", exported as "db"
    ports:
      - type: assigned
        port: 5432
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

**Rationale**: Reduces boilerplate for simple services. See ADR-0007 for port type system rationale.

### Custom Network Name

**Scenario**: Workspace specifies custom network name that conflicts with existing Docker network.

**Behavior**: Error on workspace creation with guidance to choose different name.

**Rationale**: Prevents accidental network collisions with external resources.

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Missing workspace.name | YAML_VALIDATION | `workspace.name is required` | Add required field |
| Invalid TLS mode | YAML_VALIDATION | `tls.mode must be one of: auto, custom, disabled` | Use valid mode value |
| Missing cert files | FILE_NOT_FOUND | `Certificate file not found: {path}` | Create/move cert files or change TLS mode |
| Invalid compose file reference | FILE_NOT_FOUND | `Flavor "{name}" references non-existent file: {file}` | Fix compose_files list |
| Invalid service reference | YAML_VALIDATION | `Exported service "{name}" references non-existent Compose service: {service}` | Fix service reference |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
