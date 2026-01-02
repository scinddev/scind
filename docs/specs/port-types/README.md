# Port Types Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Exported services declare ports with a `type` that determines how the port is routed, and optionally a `protocol` for proxied services.

---

## Port Types

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql` | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding | No | `*_HOST`, `*_PORT` |

---

## Port Type Descriptions

### Proxied

Traffic is routed through Traefik:
- Exported service gets a hostname (`{workspace}-{app}-{export}.{domain}`)
- Traefik labels are generated
- Environment variables contain **proxy values** (ports 80/443), not container port

**Protocols for proxied type**:
- `https`: Routes through Traefik's `websecure` entrypoint (port 443) with TLS
- `http`: Routes through Traefik's `web` entrypoint (port 80)
- `tcp`, `postgresql`, `mysql` (future): SNI-based TCP routing

### Assigned

Port is bound directly to the host:
- If specified port is unavailable, Contrail increments until available
- Assignment recorded in global state
- Environment variables point to internal alias and assigned host port

---

## Visibility

Each port can have `visibility` of `public` or `protected` (default):

- **public**: Intended for external/production use
- **protected**: Development/debugging, not for production dependency

Visibility is **documentation only**—does not change routing behavior. Exposed via Docker labels for external tools.

---

## Configuration Examples

### Basic Proxied Service

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
```

### Database with Assigned Port

```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

### Multiple Protocols

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
      - type: proxied
        protocol: http
        visibility: protected
      - type: assigned
        port: 9229                       # Debug port
        visibility: protected
```

---

## Port Inference

If `port:` is omitted:
- If Compose service has exactly one port → use that port
- If multiple ports → error with clear message

```
Error: Port must be specified for exported service "web"
  Application: app-one
  Compose service "web" has multiple ports: 80, 443, 9229
```

---

## Port Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- For multiple HTTP/HTTPS proxies, create separate exported services

---

## Port Assignment Behavior

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use recorded port (sticky assignment)

**Port conflict at startup**:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

---

## Related Documentation

- [ADR-0007: Port Type System](../../decisions/0007-port-type-system/README.md)
- [Environment Variables Spec](../environment-variables/README.md)
- [Configuration Schemas Spec](../configuration-schemas/README.md)
