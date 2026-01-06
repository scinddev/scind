# Configuration Schemas - Complete Examples

This appendix contains complete example configurations for all Contrail configuration files.

---

## Proxy Configuration

**Location**: `~/.config/contrail/proxy.yaml`

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
    # cert_file: ~/.config/contrail/certs/wildcard.crt
    # key_file: ~/.config/contrail/certs/wildcard.key
```

---

## Workspace Configuration

**Location**: `{workspace}/workspace.yaml`

### Multi-Application Workspace

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

### With Custom Templates

```yaml
workspace:
  name: dev
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
  applications:
    app-one:
      repository: git@github.com:company/app-one.git
```

---

## Application Configuration

**Location**: `{application}/application.yaml`

### Full Example

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
  debug:
    compose_files:
      - docker-compose.yaml
      - docker-compose.debug.yaml

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
    service: postgres                   # Maps to Compose service "postgres"
    ports:
      - type: assigned
        port: 5432
        visibility: protected
  worker:
    ports:
      - type: assigned
        port: 9000
        visibility: protected
```

### Minimal Example

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
```

---

## Global State

**Location**: `~/.config/contrail/state.yaml`

```yaml
# AUTO-GENERATED - Managed by Contrail
# Records assigned ports and port availability inventory

assigned_ports:
  dev:
    app-one:
      db: 5432
    app-two:
      db: 5433
      cache: 6379
  review:
    app-one:
      db: 5434

port_inventory:
  5432:
    status: assigned
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:
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
```

---

## Workspace Registry

**Location**: `~/.config/contrail/workspaces.yaml`

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

---

## Generated State

**Location**: `{workspace}/.generated/state.yaml`

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite
```

---

## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml`

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
        service: web
        alias: app-one-web
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
```
