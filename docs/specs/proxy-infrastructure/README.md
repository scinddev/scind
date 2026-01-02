# Proxy Infrastructure Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail uses Traefik as a reverse proxy to route external requests to workspace applications. The proxy is managed as a Docker Compose project by Contrail.

---

## Directory Structure

**Location**: `~/.config/contrail/proxy/`

```
proxy/
├── docker-compose.yaml    # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration
└── certs/                # TLS certificates
```

---

## Docker Compose Configuration

```yaml
name: contrail-proxy

services:
  traefik:
    image: ${TRAEFIK_IMAGE:-traefik:v3.2.3}
    command:
      - "--configFile=/etc/traefik/traefik.yaml"
      - "--api.dashboard=true"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"                      # Dashboard (if enabled)
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

---

## Traefik Static Configuration

```yaml
api:
  dashboard: true                        # Based on proxy.yaml

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

## Lifecycle Commands

| Command | Behavior |
|---------|----------|
| `proxy init` | Creates directory structure and configuration files |
| `proxy up` | Starts Traefik container (creates network if needed) |
| `proxy down` | Stops Traefik container |
| `workspace up` | Auto-starts proxy if not running |

---

## TLS Configuration

**Mode: auto** (default)
- Uses mkcert if available to generate locally-trusted certificates
- Falls back to Traefik's default self-signed certificate

**Mode: custom**
- Uses user-provided certificate and key files
- For enterprise CA or manually generated certs

**Mode: disabled**
- HTTP only, no HTTPS entrypoint

### Certificate Setup (auto with mkcert)

```bash
# Install mkcert (one-time)
mkcert -install

# Generate wildcard certificate
mkcert "*.contrail.test"
```

---

## DNS Configuration

Configure DNS resolution for workspace domains:

**Option 1: dnsmasq**
```
address=/contrail.test/127.0.0.1
```

**Option 2: /etc/hosts**
```
127.0.0.1 dev-app-one-web.contrail.test
127.0.0.1 dev-app-two-api.contrail.test
```

**Note**: The `.test` TLD is reserved by RFC 2606 for testing and won't conflict with real domains.

---

## Recovery

If proxy configuration is corrupted:
```bash
contrail proxy init --force
```

This regenerates default configuration with backup of existing files.

---

## Related Documentation

- [ADR-0008: Traefik for Reverse Proxy](../../decisions/0008-traefik-reverse-proxy/README.md)
- [ADR-0009: Flexible TLS Configuration](../../decisions/0009-flexible-tls-configuration/README.md)
- [Port Types Spec](../port-types/README.md)
