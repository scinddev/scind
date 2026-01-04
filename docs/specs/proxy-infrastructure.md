# Proxy Infrastructure Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail uses Traefik as a reverse proxy to route HTTP/HTTPS traffic to workspace services. The proxy runs as a Docker Compose project managed by Contrail and is shared across all workspaces on the host.

**Related Documents**:
- [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md)
- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)
- [Port Types](./port-types.md)
- [Configuration Schemas](./configuration-schemas.md)

**Appendices**:
- [Traefik docker-compose.yaml](./appendices/proxy-infrastructure/traefik-compose.yaml)
- [Traefik static configuration](./appendices/proxy-infrastructure/traefik-config.yaml)

---

## Behavior

### Proxy as a Docker Compose Project

The proxy is implemented as a Docker Compose project at `~/.config/contrail/proxy/`. This approach:
- Uses familiar Docker Compose patterns
- Enables standard `docker compose` debugging
- Keeps proxy configuration in a known location

### Automatic Startup

`workspace up` automatically starts the proxy if it's not running. Users rarely need to manage the proxy directly.

---

## Data Schema

### Directory Structure

**Location**: `~/.config/contrail/proxy/` (created by `contrail proxy init`)

```
~/.config/contrail/proxy/
├── docker-compose.yaml   # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

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

### traefik.yaml (Static Configuration)

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

## Networks

### contrail-proxy Network

| Property | Value |
|----------|-------|
| Name | `contrail-proxy` |
| Scope | Host-level, shared across all workspaces |
| Purpose | Connects Traefik to services that need external access |
| Created by | `contrail proxy up` or `contrail proxy init` |
| Driver | bridge |

Services with proxied exports are attached to this network, allowing Traefik to route traffic to them.

### Network Labels

The network is labeled for identification:
```
contrail.managed=true
```

---

## Examples

### Example 1: Basic Proxy Initialization

```bash
$ contrail proxy init
Created proxy configuration at ~/.config/contrail/proxy/

Next steps:
  1. Configure DNS for *.contrail.test → 127.0.0.1
     (See: contrail doctor for DNS verification)
  2. Start the proxy:
     contrail proxy up
```

### Example 2: Custom Domain

```bash
$ contrail proxy init --domain mydev.local
Created proxy configuration at ~/.config/contrail/proxy/
Domain set to: mydev.local
```

### Example 3: Proxy Status

```bash
$ contrail proxy status
Proxy: running
Network: contrail-proxy (created)
Dashboard: http://localhost:8080
Entrypoints:
  - web: :80
  - websecure: :443
```

Dashboard disabled:
```bash
$ contrail proxy status
Proxy: running
Network: contrail-proxy (created)
Dashboard: disabled
Entrypoints:
  - web: :80
  - websecure: :443
```

---

## Lifecycle

### Commands

| Command | Description |
|---------|-------------|
| `proxy init` | Creates directory structure and configuration files |
| `proxy up` | Starts the Traefik container (creates network if needed) |
| `proxy down` | Stops the Traefik container |
| `proxy restart` | Restarts the Traefik container |
| `proxy status` | Shows proxy status, network, and entrypoints |

### proxy init

Creates the proxy directory structure and configuration files.

**Flags**:
- `--force`: Overwrite existing configuration (backups old config first)
- `--domain`: Set proxy domain (default: `contrail.test`)
- `--path`: Directory to create proxy in (default: `~/.config/contrail/proxy/`)

**Behavior**:
1. Check if proxy configuration already exists
   - If exists and no `--force`: error with message
   - If exists and `--force`: backup existing config and overwrite
2. Create proxy directory structure
3. Create `contrail-proxy` Docker network if it doesn't exist
4. Output next steps (DNS setup, starting proxy)

### proxy up

**Flags**:
- `--recreate`: Recreate the proxy network even if it exists

**Behavior**:
- Creates `contrail-proxy` network if it doesn't exist
- Validates existing network configuration matches expected settings
- Starts Traefik container from proxy configuration
- If proxy configuration doesn't exist, runs `proxy init` first

**Network conflict handling**:
```
Warning: Network 'contrail-proxy' exists but may not have been created by Contrail.
  Driver: bridge (expected: bridge) ✓
  Labels: contrail.managed not found ⚠

Use 'contrail proxy up --recreate' to recreate the network.
```

### Recovery

If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration:

```bash
$ contrail proxy init --force
Backed up existing configuration to ~/.config/contrail/proxy.backup.20241230/
Created proxy configuration at ~/.config/contrail/proxy/
```

---

## TLS Configuration

### TLS Modes

Configured in `~/.config/contrail/proxy.yaml`:

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available; falls back to Traefik's default self-signed certificate |
| `custom` | Uses user-provided certificate and key files |
| `disabled` | HTTP only, no HTTPS entrypoint |

### Certificate Setup

**auto with mkcert**:
1. Run `mkcert -install` once per machine to add local CA to trust store
2. Run `mkcert "*.contrail.test"` to generate wildcard certificate
3. Contrail detects and uses these automatically

**custom (enterprise CA)**:
1. Obtain wildcard certificate from your CA for `*.contrail.test`
2. Configure paths in `proxy.yaml`:
   ```yaml
   proxy:
     tls:
       mode: custom
       cert_file: ~/.config/contrail/certs/wildcard.crt
       key_file: ~/.config/contrail/certs/wildcard.key
   ```

**auto without mkcert**:
Traefik serves its default self-signed certificate. Browsers will show security warnings.

---

## DNS Configuration

For local development, configure DNS resolution for workspace domains.

### Options

**dnsmasq** (recommended):
```
address=/contrail.test/127.0.0.1
```

**/etc/hosts** (manual entries):
```
127.0.0.1 dev-app-one-web.contrail.test
127.0.0.1 dev-app-two-api.contrail.test
```

**Local DNS server**: More complex but flexible for team settings.

### .test TLD

The `.test` TLD is reserved by RFC 2606 for testing purposes:
- Will not conflict with real domains
- Will not conflict with mDNS (unlike `.local`)
- Safe for local development

---

## Edge Cases

### Dashboard Port Conflict

**Scenario**: Port 8080 is already in use.

**Behavior**: Error on proxy startup:
```
Error: Cannot bind to port 8080 for Traefik dashboard.
  Another process is using this port.

To resolve:
  1. Stop the other process, or
  2. Change the dashboard port in ~/.config/contrail/proxy.yaml:
     proxy:
       dashboard:
         port: 8081
```

### Network Already Exists (External)

**Scenario**: Network `contrail-proxy` was created by another tool.

**Behavior**: Warning on `proxy up`:
```
Warning: Network 'contrail-proxy' exists but may not have been created by Contrail.
  Driver: bridge (expected: bridge) ✓
  Labels: contrail.managed not found ⚠

Use 'contrail proxy up --recreate' to recreate the network.
```

### Proxy Not Running During workspace up

**Scenario**: User runs `workspace up` but proxy is down.

**Behavior**: Contrail automatically starts the proxy:
```
$ contrail workspace up
Proxy not running, starting...
Proxy: started

Starting app-one... done
Starting app-two... done
```

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Docker not running | DOCKER | `Docker is not installed or not running` | Start Docker |
| Port 80/443 in use | NETWORK | `Cannot bind to port 80: address already in use` | Stop conflicting process |
| Network conflict | NETWORK | `Network exists with incompatible settings` | Use `--recreate` flag |
| Missing cert files | FILE_NOT_FOUND | `Certificate file not found` | Create certs or change TLS mode |
| Config corruption | PARSE | `Invalid proxy configuration` | Use `proxy init --force` |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
