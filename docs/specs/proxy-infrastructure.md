# Specification: Proxy Infrastructure

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

The proxy is implemented as a Docker Compose project managed by Contrail. It runs a Traefik instance that handles reverse proxying for all workspaces on the host.

---

## Directory Structure

**Location**: `~/.config/contrail/proxy/` (global/per-user)

Created by `contrail proxy init`:

```
~/.config/contrail/proxy/
├── docker-compose.yaml    # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

---

## Generated Files

### docker-compose.yaml

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
      - "8080:8080"                          # Only included if dashboard.enabled: true
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

**Notes**:
- The dashboard port (`8080:8080`) is only included when `proxy.dashboard.enabled` is `true` (the default)
- The `--api.dashboard` flag is set based on the same configuration
- The Traefik image can be customized via `proxy.yaml`

### traefik.yaml

```yaml
api:
  dashboard: true                          # Set based on proxy.yaml dashboard.enabled

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    exposedByDefault: false
    network: contrail-proxy
  file:
    directory: /etc/traefik/dynamic
    watch: true
```

---

## Network

The `contrail-proxy` network connects Traefik to all application services that need external access.

| Aspect | Value |
|--------|-------|
| Name | `contrail-proxy` |
| Scope | Host-wide, shared across all workspaces |
| Created by | `proxy init` or `proxy up` |
| Type | External (must exist before container start) |

---

## Lifecycle Commands

| Command | Behavior |
|---------|----------|
| `proxy init` | Creates the directory structure and configuration files |
| `proxy up` | Starts the Traefik container (creates `contrail-proxy` network if needed) |
| `proxy down` | Stops the Traefik container |
| `workspace up` | Automatically runs `proxy up` if proxy is not running |

---

## Network Conflict Handling

If the `contrail-proxy` network already exists but was not created by Contrail:
- `proxy up` will use the existing network
- Contrail does not validate network configuration compatibility
- Users are responsible for ensuring the network is properly configured

---

## Recovery

If a user manually edits the proxy configuration and breaks it:
- `proxy init --force` regenerates the default configuration
- Existing certificates in `certs/` are preserved
- The proxy container must be restarted after regeneration

---

## External Access Routing

For a service to be accessible externally through Traefik:

1. The service must have `type: proxied` in its `exported_services` configuration
2. The generated override file adds the service to the `contrail-proxy` network
3. Traefik labels are generated for hostname-based routing
4. Traefik automatically discovers the service via Docker labels

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md)
- [Generated Override Files](generated-override-files.md)
- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)
