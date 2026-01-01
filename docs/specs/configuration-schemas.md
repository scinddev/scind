# Specification: Configuration Schemas

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

The system uses three schema types, separating structure (configuration) from state (runtime):

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/contrail/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

---

## Proxy Configuration

**Location**: `~/.config/contrail/proxy.yaml` (global/per-user)

```yaml
proxy:
  domain: contrail.test                  # TLD for generated hostnames
  traefik_image: traefik:v3.2.3          # Traefik Docker image
  dashboard:
    enabled: true                        # Enable/disable Traefik dashboard (default: true)
    port: 8080                           # Dashboard port (default: 8080)
  tls:
    mode: auto                           # auto | custom | disabled
    # For mode: custom (e.g., enterprise CA certificates)
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

### TLS Modes

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available; falls back to Traefik's self-signed certificate |
| `custom` | Uses user-provided certificate and key files |
| `disabled` | HTTP only, no HTTPS entrypoint |

### Certificate Setup by Mode

- **auto with mkcert**: Run `mkcert -install` once per machine, then `mkcert "*.contrail.test"` to generate a wildcard certificate
- **custom (enterprise CA)**: Obtain a wildcard certificate signed by your enterprise CA
- **auto without mkcert**: Traefik serves its default self-signed certificate (browser warnings)

---

## DNS Configuration

For local development, configure DNS resolution for the workspace domains. Options include:

1. **dnsmasq**: Route all `*.contrail.test` to `127.0.0.1`
   ```
   address=/contrail.test/127.0.0.1
   ```
2. **/etc/hosts**: Manual entries for each hostname
3. **Local DNS server**: More complex but flexible

**Note**: The `.test` TLD is reserved by RFC 2606 for testing purposes and will not conflict with real domains or mDNS (unlike `.local`).

---

## Workspace Configuration

**Location**: `{workspace}/workspace.yaml`

```yaml
workspace:
  name: dev                           # Required. Used as prefix for project names
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

### Conventions

- Application name = directory path (e.g., `app-one` → `./{workspace}/app-one/`)
- Network name defaults to `{workspace.name}-internal`

---

## Application Configuration (Service Contract)

**Location**: `{application}/application.yaml` (lives in the application's git repository)

This file defines the application's exported services and flavors.

```yaml
default_flavor: full                    # Optional. Defaults to "default"

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

### Mapping to a Different Compose Service

Use the `service:` property when the exported name differs from the Compose service name:

```yaml
exported_services:
  db:
    service: postgres                   # Maps to Compose service "postgres", exported as "db"
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

---

## Template Variables

Contrail uses template variables to generate hostnames, aliases, and other computed values.

### Default Templates

```yaml
workspace:
  name: dev
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
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
| `%SERVICE_PORT%` | Export | Container port number | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

---

## Workspace Registry

**Location**: `~/.config/contrail/workspaces.yaml`

Tracks all known workspaces across the system:

```yaml
workspaces:
  dev:
    path: /home/user/workspaces/dev
    registered_at: 2024-12-28T10:30:00Z
    last_seen: 2024-12-29T14:22:00Z
  review:
    path: /home/user/workspaces/review
    registered_at: 2024-12-28T11:00:00Z
    last_seen: 2024-12-29T13:15:00Z
```

**Registry behavior**:
- `workspace init` automatically registers the workspace
- `workspace list` reads the registry and optionally validates entries
- `workspace prune` removes stale entries
- **Docker label fallback**: If registry is missing, reconstruct from Docker container labels

---

## State Management

### Workspace State

**Location**: `{workspace}/.generated/state.yaml` (gitignored)

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite                      # Overridden from default
```

### Flavor Resolution Order

1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

A computed, read-only view of the workspace's current state:

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

---

## Related Documents

- [Port Types and Proxying](port-types.md)
- [Generated Override Files](generated-override-files.md)
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md)
