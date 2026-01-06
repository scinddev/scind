# Proxy Infrastructure Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail uses Traefik as a reverse proxy to route HTTP/HTTPS traffic to workspace services. The proxy runs as a Docker Compose project managed by Contrail and is shared across all workspaces on the host.

**Location**: `~/.config/contrail/proxy/` (global/per-user)

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

Created by `contrail proxy init`:

```
~/.config/contrail/proxy/
├── docker-compose.yaml    # Traefik service definition
├── traefik.yaml          # Traefik static configuration
├── dynamic/              # Dynamic configuration (auto-discovered)
│   └── tls.yaml          # TLS certificate configuration (generated)
└── certs/                # TLS certificates (copied or generated here)
```

### Generated Configuration Files

#### docker-compose.yaml

See [appendices/proxy-infrastructure/traefik-compose.yaml](./appendices/proxy-infrastructure/traefik-compose.yaml) for the complete generated file.

Key features:
- Configurable Traefik image via `proxy.yaml`
- Ports 80, 443, and optional dashboard port (8080 by default)
- Volume mounts for Docker socket, config, and certificates
- Uses external `contrail-proxy` network
- Managed container labels for identification

#### traefik.yaml

See [appendices/proxy-infrastructure/traefik-config.yaml](./appendices/proxy-infrastructure/traefik-config.yaml) for the complete generated file.

Key features:
- Dashboard enabled by default (configurable)
- Web entrypoint on port 80
- Websecure entrypoint on port 443
- Docker provider with `exposedByDefault: false`
- File provider for dynamic configuration

---

## Lifecycle

| Command | Description |
|---------|-------------|
| `proxy init` | Creates directory structure and configuration files |
| `proxy up` | Starts the Traefik container (creates network if needed) |
| `proxy down` | Stops the Traefik container |
| `proxy restart` | Restarts the Traefik container |
| `proxy status` | Shows proxy status, network, and entrypoints |
| `workspace up` | Automatically runs `proxy up` if proxy is not running |

### proxy init Behavior

1. Check if proxy configuration already exists
   - If exists and no `--force`: error with message
   - If exists and `--force`: backup existing config and overwrite
2. Create proxy directory structure
3. Create `contrail-proxy` Docker network if it doesn't exist
4. Output next steps (DNS setup, starting proxy)

### Recovery

If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration:

```bash
$ contrail proxy init --force
Backed up existing configuration to ~/.config/contrail/proxy.backup.20241230/
Created proxy configuration at ~/.config/contrail/proxy/
```

---

## Network

The proxy uses the `contrail-proxy` network, which is external and shared across all workspaces. This network:

- Is created automatically by `proxy up` if it doesn't exist
- Connects Traefik to all proxied services
- Is referenced as `external: true` in the proxy's docker-compose.yaml
- Is labeled with `contrail.managed=true`

| Property | Value |
|----------|-------|
| Name | `contrail-proxy` |
| Scope | Host-level, shared across all workspaces |
| Purpose | Connects Traefik to services that need external access |
| Created by | `contrail proxy up` or `contrail proxy init` |
| Driver | bridge |

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

For local development, configure DNS resolution for the workspace domains. Options include:

### dnsmasq (recommended)

Route all `*.contrail.test` to `127.0.0.1`:
```
address=/contrail.test/127.0.0.1
```

### /etc/hosts (manual entries)

```
127.0.0.1 dev-app-one-web.contrail.test
127.0.0.1 dev-app-two-api.contrail.test
```

### Local DNS server

More complex but flexible for team settings.

**Note**: The `.test` TLD is reserved by RFC 2606 for testing purposes and will not conflict with real domains or mDNS (unlike `.local`).

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

### Example 2: Proxy Status

```bash
$ contrail proxy status
Proxy: running
Network: contrail-proxy (created)
Dashboard: http://localhost:8080
Entrypoints:
  - web: :80
  - websecure: :443
```

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

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Docker not running | `Docker is not installed or not running` | Start Docker |
| Port 80/443 in use | `Cannot bind to port 80: address already in use` | Stop conflicting process |
| Network conflict | `Network exists with incompatible settings` | Use `--recreate` flag |
| Missing cert files | `Certificate file not found` | Create certs or change TLS mode |
| Config corruption | `Invalid proxy configuration` | Use `proxy init --force` |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
