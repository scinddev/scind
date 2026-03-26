# Configuration Complete Examples

Full working configuration examples for Scind.

---

## Multi-Application Workspace

A complete multi-application workspace with frontend, backend, and shared database.

### workspace.yaml

```yaml
workspace:
  name: dev
  # network: dev-custom              # Optional. Defaults to dev-internal
  applications:
    frontend:
      repository: git@github.com:company/frontend.git
    backend:
      repository: git@github.com:company/backend.git
    shared-db:
      repository: git@github.com:company/shared-db.git
      path: ./database                # Custom path (default: ./shared-db)
  templates:                          # Optional. Template customization
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
    hostname-apex: "%WORKSPACE_NAME%-%APPLICATION_NAME%.%PROXY_DOMAIN%"
    alias-apex: "%APPLICATION_NAME%"
```

### frontend/application.yaml

```yaml
default_flavor: full

flavors:
  lite:
    compose_files:
      - docker-compose.yaml
  full:
    compose_files:
      - docker-compose.yaml
      - docker-compose.dev.yaml

exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
        port: 80
      - type: proxied
        protocol: http
        visibility: protected
        port: 80
```

### backend/application.yaml

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
  debug:
    compose_files:
      - docker-compose.yaml
      - docker-compose.debug.yaml

exported_services:
  api:
    primary: true                     # Explicit primary — gets apex URL
    service: node                     # Maps to Compose service "node"
    ports:
      - type: proxied
        protocol: https
        visibility: public
        port: 3000
      - type: assigned
        port: 9229                    # Node.js debug port
        visibility: protected
  worker:
    ports:
      - type: assigned
        port: 9000
        visibility: protected
```

### shared-db/application.yaml

Multiple exports, no primary — no apex URL generated:

```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
  cache:
    service: redis
    ports:
      - type: assigned
        port: 6379
        visibility: protected
```

---

## Single-Application Workspace

For promoting an existing Docker Compose project:

### workspace.yaml

```yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .                         # Application in workspace root directory
```

### application.yaml

```yaml
default_flavor: default

flavors:
  default:
    compose_files:
      - docker-compose.yaml
  full:
    compose_files:
      - docker-compose.yaml
      - docker-compose.worker.yaml

exported_services:
  web:
    service: nginx
    ports:
      - type: proxied
        protocol: https
        visibility: public
        port: 80
  api:
    service: php
    ports:
      - type: proxied
        protocol: https
        visibility: public
        port: 9000
  db:
    service: mysql
    ports:
      - type: assigned
        port: 3306
        visibility: protected
```

---

## Proxy Configuration

**Location**: `~/.config/scind/proxy.yaml`

```yaml
proxy:
  domain: scind.test
  traefik_image: traefik:v3.2.3
  dashboard:
    enabled: true
    port: 8080
  tls:
    mode: auto                        # auto | custom | disabled
    # For mode: custom
    # cert_file: ~/.config/scind/certs/wildcard.crt
    # key_file: ~/.config/scind/certs/wildcard.key
```

### Custom TLS Configuration

For enterprise environments with custom CA certificates:

```yaml
proxy:
  domain: dev.company.local
  dashboard:
    enabled: true
    port: 8080
  tls:
    mode: custom
    cert_file: ~/.config/scind/certs/wildcard.crt
    key_file: ~/.config/scind/certs/wildcard.key
```

---

## Service with Multiple Port Types

An exported service can combine proxied and assigned ports:

```yaml
exported_services:
  web:
    service: nginx
    ports:
      # Public HTTPS endpoint
      - type: proxied
        protocol: https
        visibility: public
        port: 443
      # Protected HTTP endpoint for internal tools
      - type: proxied
        protocol: http
        visibility: protected
        port: 80
      # Direct access for debugging
      - type: assigned
        port: 8080
        visibility: protected
```

---

## Manual Override Example

**Location**: `{workspace}/overrides/frontend.yaml`

Workspace-specific customizations that persist across regeneration:

```yaml
services:
  web:
    environment:
      - DEBUG=true
      - FEATURE_FLAG_X=enabled
    labels:
      - "traefik.http.routers.dev-frontend-web-https.middlewares=dev-auth@docker"

  node:
    volumes:
      - ./local-dev-data:/app/data:rw
```

---

## Global State

**Location**: `~/.config/scind/state.yaml`

```yaml
# AUTO-GENERATED - Managed by Scind
# Records assigned ports and port availability inventory

assigned_ports:
  dev:
    frontend:
      web: 8080
    backend:
      api: 3000
      worker: 9000
    shared-db:
      db: 5432
      cache: 6379
  review:
    frontend:
      web: 8081

port_inventory:
  5432:
    status: assigned
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:
      workspace: dev
      application: shared-db
      exported_service: db
  6379:
    status: assigned
    first_seen: 2025-12-28T17:53:58Z
    last_checked: 2025-12-29T13:01:37Z
    assignment:
      workspace: dev
      application: shared-db
      exported_service: cache
```

---

## Workspace Registry

**Location**: `~/.config/scind/workspaces.yaml`

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

---

## Generated State

**Location**: `{workspace}/.generated/state.yaml`

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  frontend:
    flavor: full
  backend:
    flavor: full
  shared-db:
    flavor: default
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
  domain: scind.test

applications:
  frontend:
    flavor: full
    project: dev-frontend
    exported_services:
      web:
        service: web
        alias: frontend-web
        primary: true                     # Implicit (single export)
        apex_alias: frontend
        ports:
          - type: proxied
            protocol: https
            container_port: 443
            host_port: 443
            visibility: public
            hostname: dev-frontend-web.scind.test
            apex_hostname: dev-frontend.scind.test
        environment:
          SCIND_FRONTEND_WEB_HOST: dev-frontend-web.scind.test
          SCIND_FRONTEND_WEB_PORT: 443
          SCIND_FRONTEND_WEB_SCHEME: https
          SCIND_FRONTEND_WEB_URL: https://dev-frontend-web.scind.test
          SCIND_FRONTEND_APEX_HOST: dev-frontend.scind.test
          SCIND_FRONTEND_APEX_PORT: 443
          SCIND_FRONTEND_APEX_SCHEME: https
          SCIND_FRONTEND_APEX_URL: https://dev-frontend.scind.test
```

---

## Complete Directory Structure Example

```
~/.config/scind/
├── proxy.yaml                        # Global proxy configuration
├── state.yaml                        # Global port assignments
└── workspaces.yaml                   # Workspace registry

~/workspaces/
├── proxy/
│   ├── docker-compose.yaml           # Traefik service definition
│   ├── traefik.yaml                  # Traefik static configuration
│   ├── dynamic/
│   │   └── tls.yaml                  # TLS certificate configuration
│   └── certs/
│       ├── wildcard.crt
│       └── wildcard.key
│
└── dev/
    ├── workspace.yaml
    ├── overrides/
    │   └── frontend.yaml
    ├── .generated/
    │   ├── state.yaml
    │   ├── manifest.yaml
    │   ├── frontend.override.yaml
    │   ├── backend.override.yaml
    │   └── shared-db.override.yaml
    ├── frontend/
    │   ├── docker-compose.yaml
    │   ├── docker-compose.dev.yaml
    │   ├── application.yaml
    │   └── src/
    ├── backend/
    │   ├── docker-compose.yaml
    │   ├── docker-compose.worker.yaml
    │   ├── docker-compose.debug.yaml
    │   ├── application.yaml
    │   └── src/
    └── database/
        ├── docker-compose.yaml
        ├── application.yaml
        └── init/
```

---

## Related Documents

- [Configuration Reference](../../configuration.md)
- [Generated Override Files](../../../specs/appendices/generated-override-files/complete-override-example.yaml)
