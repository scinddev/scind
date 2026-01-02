# Configuration Schemas: Complete Examples

This appendix provides complete, working examples of all three configuration schemas.

---

## proxy.yaml — Complete Example

```yaml
# ~/.config/contrail/proxy.yaml
# Global proxy configuration for all workspaces

proxy:
  # TLD for generated hostnames
  domain: contrail.test

  # Traefik Docker image (pinned for reproducibility)
  traefik_image: traefik:v3.2.3

  # Dashboard configuration
  dashboard:
    enabled: true
    port: 8080

  # TLS configuration
  tls:
    mode: auto
    # Uncomment for custom certificates:
    # mode: custom
    # cert_file: ~/.config/contrail/certs/wildcard.crt
    # key_file: ~/.config/contrail/certs/wildcard.key
```

---

## workspace.yaml — Multi-Application Workspace

```yaml
# ~/workspaces/dev/workspace.yaml
# Multi-application development workspace

workspace:
  name: dev

  # Optional: custom network name (defaults to dev-internal)
  # network: dev-custom

  # Optional: custom naming templates
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"

  applications:
    frontend:
      repository: git@github.com:company/frontend.git

    backend:
      repository: git@github.com:company/backend.git

    api:
      repository: git@github.com:company/api.git
      path: ./services/api    # Custom path

    shared-db:
      repository: git@github.com:company/shared-db.git
```

---

## workspace.yaml — Single-Application Workspace

```yaml
# ~/my-project/workspace.yaml
# Single-application workspace where app is at workspace root

workspace:
  name: dev
  applications:
    myapp:
      path: .    # Application is in workspace root
```

---

## application.yaml — Web Application

```yaml
# ~/workspaces/dev/frontend/application.yaml
# Frontend application with multiple flavors

default_flavor: full

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

  assets:
    service: nginx
    ports:
      - type: proxied
        protocol: https
        visibility: public
```

---

## application.yaml — API Service

```yaml
# ~/workspaces/dev/api/application.yaml
# API service with database and cache

default_flavor: default

flavors:
  default:
    compose_files:
      - docker-compose.yaml

exported_services:
  api:
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

  cache:
    service: redis
    ports:
      - type: assigned
        port: 6379
        visibility: protected

  debug:
    service: api
    ports:
      - type: assigned
        port: 9229
        visibility: protected
```

---

## application.yaml — Minimal Example

```yaml
# Minimal application.yaml
# Uses defaults where possible

exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
```
