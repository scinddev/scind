# Port Types Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Each exported service port has a `type` that determines how traffic reaches the container. The two port types—`proxied` and `assigned`—serve different use cases and have distinct behaviors for routing, environment variables, and discovery.

**Related Documents**:
- [ADR-0007: Port Type System](../decisions/0007-port-type-system.md)
- [Environment Variables](./environment-variables.md)
- [Proxy Infrastructure](./proxy-infrastructure.md)

---

## Behavior

### Port Type Summary

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `*_HOST`, `*_PORT` |

---

## Data Schema

### Proxied Type

Traffic is routed through Traefik. The exported service gets a hostname and Traefik labels are generated.

```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https           # Required: http, https, or future SNI types
        port: 8080                # Optional: container port (can be inferred)
        visibility: public        # Optional: public or protected (default: protected)
```

**Environment variables contain proxy values** (hostname and proxy port 80/443), not the container port.

#### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `proxied` |
| `protocol` | string | Yes | - | `http`, `https`, or future SNI types |
| `port` | integer | No | (inferred) | Container port |
| `visibility` | string | No | `protected` | `public` or `protected` |

#### Protocol Options

| Protocol | Entrypoint | Port | TLS |
|----------|------------|------|-----|
| `https` | `websecure` | 443 | Yes (termination at Traefik) |
| `http` | `web` | 80 | No |
| `tcp`, `postgresql`, `mysql` | Future | Custom | SNI-based routing |

### Assigned Type

The port is bound directly to the host. If the specified port is unavailable, Contrail increments until an available port is found.

```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432                # Preferred port (may be reassigned)
        visibility: protected
```

**Environment variables point to the internal alias and assigned host port.**

#### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | Must be `assigned` |
| `port` | integer | Yes | - | Preferred host port |
| `visibility` | string | No | `protected` | `public` or `protected` |

---

## Examples

### Example 1: Simple Proxied Web Service

**Configuration**:
```yaml
# docker-compose.yaml
services:
  web:
    ports:
      - "8080"
```

```yaml
# application.yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
```

**Result**:
- Hostname: `dev-app-one-web.contrail.test`
- Traefik routes HTTPS to container port 8080
- Environment: `CONTRAIL_APP_ONE_WEB_PORT=443` (proxy port, not container port)

### Example 2: Database with Assigned Port

**Configuration**:
```yaml
exported_services:
  db:
    service: mysql
    ports:
      - type: assigned
        port: 3306
        visibility: protected
```

**Result** (port 3306 available):
- Host port binding: `3306:3306`
- Environment: `CONTRAIL_APP_ONE_DB_PORT=3306`

**Result** (port 3306 taken by another workspace):
- Host port binding: `3307:3306`
- Environment: `CONTRAIL_APP_ONE_DB_PORT=3307`
- Assignment recorded in global state for persistence

### Example 3: Service with Both Proxy and Direct Ports

**Configuration**:
```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        port: 443
        visibility: public
      - type: proxied
        protocol: http
        port: 80
        visibility: protected
      - type: assigned
        port: 9229
        visibility: protected           # Node.js debug port
```

**Result**:
- HTTPS via Traefik to port 443
- HTTP via Traefik to port 80
- Debug port 9229 bound directly to host

---

## Validation Rules

### Port Type Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

### Port Inference Rules

- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

---

## Visibility

| Visibility | Description |
|------------|-------------|
| `public` | This port is intended for external/production use |
| `protected` (default) | This port exists for development/debugging |

Visibility is primarily **documentation** to communicate intent to collaborators. It does not change Contrail's core behavior:
- All exported services receive internal network aliases
- All exported services receive environment variables
- Both public and protected proxied services route through Traefik

Visibility is exposed in Docker labels (`contrail.export.*.visibility`) enabling external tools (dashboards, service discovery) to filter or display services differently.

---

## Private Services

Services not listed in `exported_services` remain private. They are only accessible within the application's own compose network, following standard Docker Compose behavior.

---

## Edge Cases

### Port Assignment Algorithm

**Scenario**: Preferred port is unavailable.

**Behavior**:
1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in global state (`~/.config/contrail/state.yaml`)
4. Subsequent runs use the recorded port (sticky assignment)

### Port Conflict at Startup

**Scenario**: Previously assigned port has become unavailable when `workspace up` runs.

**Behavior**: Error with resolution steps:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

### Port Status Transitions

| From | To | Trigger |
|------|-----|---------|
| `unavailable` | `assigned` | Port became free, Contrail claimed it |
| `assigned` | `released` | Workspace/app removed, port freed |
| `unavailable` | `released` | External process stopped, `contrail port gc` cleaned it up |

### Port Availability Checking

`contrail port scan` and `contrail port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state.

This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

### Force Flag Behavior

`--force` on `workspace generate` regenerates override files but **preserves existing port assignments**. This prevents accidental port reassignment that could break external connections.

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Port conflict at startup | `Port {port} is assigned but no longer available` | Release port and regenerate |
| Invalid port type | `Port type must be "proxied" or "assigned"` | Fix configuration |
| Missing protocol | `Protocol is required for proxied ports` | Add protocol field |
| Multiple HTTP protocols | `Each exported service may have at most one http and one https port` | Split into separate exports |
| Port out of range | `Port must be between 1 and 65535` | Use valid port number |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
