# Configuration Reference

**Version**: 0.5.0
**Date**: December 2024

---

## Configuration Files

| File | Location | Purpose |
|------|----------|---------|
| `proxy.yaml` | `~/.config/contrail/` | Global proxy settings |
| `workspace.yaml` | Workspace root | Workspace definition |
| `application.yaml` | Application directory | Service contract |

---

## proxy.yaml

Global proxy configuration.

```yaml
proxy:
  domain: contrail.test
  traefik_image: traefik:v3.2.3
  dashboard:
    enabled: true
    port: 8080
  tls:
    mode: auto                   # auto | custom | disabled
    cert_file: path/to/cert      # for mode: custom
    key_file: path/to/key        # for mode: custom
```

### Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `proxy.domain` | string | `contrail.test` | TLD for generated hostnames |
| `proxy.traefik_image` | string | `traefik:v3.2.3` | Docker image for Traefik |
| `proxy.dashboard.enabled` | bool | `true` | Enable Traefik dashboard |
| `proxy.dashboard.port` | int | `8080` | Dashboard port |
| `proxy.tls.mode` | string | `auto` | TLS mode: auto, custom, disabled |
| `proxy.tls.cert_file` | string | - | Path to certificate (for custom mode) |
| `proxy.tls.key_file` | string | - | Path to key (for custom mode) |

---

## workspace.yaml

Workspace definition.

```yaml
workspace:
  name: dev
  network: dev-internal          # optional
  templates:                     # optional
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
  applications:
    app-one:
      repository: git@github.com:org/app-one.git
    app-two:
      repository: git@github.com:org/app-two.git
      path: ./custom-path
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `workspace.name` | string | Yes | Workspace identifier |
| `workspace.network` | string | No | Internal network name (default: `{name}-internal`) |
| `workspace.templates` | object | No | Template overrides |
| `workspace.applications` | map | No | Application definitions |
| `applications.{name}.repository` | string | No | Git URL for cloning |
| `applications.{name}.path` | string | No | Custom path (default: `./{name}`) |

---

## application.yaml

Application service contract.

```yaml
default_flavor: full

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
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `default_flavor` | string | No | Default flavor (default: `"default"`) |
| `flavors` | map | No | Named configurations |
| `flavors.{name}.compose_files` | list | Yes | Compose files for this flavor |
| `exported_services` | map | Yes | Services to export |
| `exported_services.{name}.service` | string | No | Compose service name (default: map key) |
| `exported_services.{name}.ports` | list | Yes | Port configurations |

### Port Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | `proxied` or `assigned` |
| `protocol` | string | For proxied | `http`, `https`, `tcp`, etc. |
| `port` | int | No | Container port (inferred if omitted) |
| `visibility` | string | No | `public` or `protected` (default) |

---

## Template Variables

Available in `workspace.templates`:

| Variable | Scope | Example |
|----------|-------|---------|
| `%PROXY_DOMAIN%` | Proxy | `contrail.test` |
| `%WORKSPACE_NAME%` | Workspace | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | `dev-internal` |
| `%APPLICATION_NAME%` | Application | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | `default` |
| `%EXPORTED_SERVICE%` | Export | `web` |
| `%SERVICE_NAME%` | Export | `nginx` |
| `%SERVICE_PORT%` | Export | `8080` |
| `%SERVICE_PROTOCOL%` | Export | `https` |

---

## State Files

### ~/.config/contrail/state.yaml

Global port assignments.

```yaml
assigned_ports:
  dev:
    app-one:
      db: 5432

port_inventory:
  5432:
    status: assigned
    assignment:
      workspace: dev
      application: app-one
      exported_service: db
```

### .generated/state.yaml

Workspace runtime state.

```yaml
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite
```

---

## Related Documentation

- [Configuration Schemas Spec](../../specs/configuration-schemas/README.md)
- [Naming Conventions Spec](../../specs/naming-conventions/README.md)
- [ADR-0006: Three Configuration Schemas](../../decisions/0006-three-configuration-schemas/README.md)
