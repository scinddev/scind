# Configuration Schemas Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail uses three configuration schemas that separate structure from state:

| Schema | Location | Owner | Purpose |
|--------|----------|-------|---------|
| `proxy.yaml` | `~/.config/contrail/` | User | Global proxy settings |
| `workspace.yaml` | Workspace root | Workspace owner | Workspace definition |
| `application.yaml` | Application directory | Application team | Service contract |

---

## Proxy Configuration

**Location**: `~/.config/contrail/proxy.yaml`

```yaml
proxy:
  domain: contrail.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard
    port: 8080                           # Dashboard port
  tls:
    mode: auto                           # auto | custom | disabled
    # For mode: custom
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

### TLS Modes

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available; falls back to Traefik self-signed |
| `custom` | Uses user-provided certificate and key files |
| `disabled` | HTTP only, no HTTPS entrypoint |

---

## Workspace Configuration

**Location**: `{workspace}/workspace.yaml`

```yaml
workspace:
  name: dev                              # Required. Prefix for project names
  # network: dev-custom                  # Optional. Defaults to {name}-internal
  applications:
    app-one:
      repository: git@github.com:company/app-one.git
    app-two:
      repository: git@github.com:company/app-two.git
      path: ./custom-path                # Optional. Defaults to ./{app-name}
```

### Single-Application Workspace

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                            # Application in workspace root
```

---

## Application Configuration (Service Contract)

**Location**: `{application}/application.yaml`

```yaml
default_flavor: full                     # Optional. Defaults to "default"

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
      - type: proxied
        protocol: http
        visibility: protected
  api:
    ports:
      - type: proxied
        protocol: https
        visibility: public
  db:
    service: postgres                    # Maps to different Compose service
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

### Exported Service Configuration

Each key in `exported_services` is the "exported service name" used for hostname generation, network aliases, and environment variables.

**Mapping to a different Compose service**:
```yaml
exported_services:
  db:
    service: postgres                    # Compose service name
    ports:
      - type: assigned
        port: 5432
```

### Port Configuration

```yaml
ports:
  - type: proxied                        # Required: proxied or assigned
    protocol: https                      # Required for proxied
    port: 8080                           # Optional: inferred from Compose
    visibility: public                   # Optional: public or protected
```

**Port type constraints**:
- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports

**Port inference rules**:
- If Compose service has exactly one port, that port is used
- If multiple ports exist, `port:` must be explicitly specified

---

## State Files

### Global State

**Location**: `~/.config/contrail/state.yaml`

```yaml
assigned_ports:
  dev:
    app-one:
      db: 5432
    app-two:
      db: 5433                           # Incremented because 5432 taken

port_inventory:
  5432:
    status: assigned                     # assigned | unavailable | released
    assignment:
      workspace: dev
      application: app-one
      exported_service: db
```

### Workspace State

**Location**: `{workspace}/.generated/state.yaml`

```yaml
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite                         # Overridden from default
```

**Flavor resolution order**:
1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Template Variables

Templates can be customized at the workspace level:

```yaml
workspace:
  name: dev
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

### Available Variables

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

## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml`

Read-only computed view of the workspace's current state:

```yaml
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
        service: web
        alias: app-one-web
        ports:
          - type: proxied
            protocol: https
            hostname: dev-app-one-web.contrail.test
        environment:
          CONTRAIL_APP_ONE_WEB_HOST: dev-app-one-web.contrail.test
          CONTRAIL_APP_ONE_WEB_PORT: 443
          CONTRAIL_APP_ONE_WEB_URL: https://dev-app-one-web.contrail.test
```

---

## Related Documentation

- [ADR-0006: Three Configuration Schemas](../../decisions/0006-three-configuration-schemas/README.md)
- [Naming Conventions Spec](../naming-conventions/README.md)
- [Environment Variables Spec](../environment-variables/README.md)
