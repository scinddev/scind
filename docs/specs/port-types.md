# Specification: Port Types and Proxying

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

Exported services declare ports with a `type` that determines how the port is routed, and optionally a `protocol` for proxied services.

---

## Port Types

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `*_HOST`, `*_PORT` |

---

## Port Type Descriptions

### Proxied

Traffic is routed through Traefik. The exported service gets a hostname (`{workspace}-{app}-{export}.{domain}`) and Traefik labels are generated. Environment variables contain the **proxy values** (hostname and proxy port 80/443), not the container port.

### Assigned

The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Contrail increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

---

## Protocol (for proxied type)

When `type: proxied`, the `protocol` field is **required** and specifies how Traefik routes the traffic:

- **https**: Routes through Traefik's `websecure` entrypoint (port 443) with TLS termination
- **http**: Routes through Traefik's `web` entrypoint (port 80)
- **tcp**, **postgresql**, **mysql**, etc. (future): SNI-based TCP routing for database connections

---

## Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent to collaborators:

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Contrail's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

**Docker label exposure**: Visibility is included in the generated Docker labels (`workspace.visibility=public` or `workspace.visibility=protected`), enabling external tools (such as Servlo) to distinguish between public and protected services for display or filtering purposes.

---

## Private Services

Services not listed in `exported_services` remain private (standard Docker Compose behavior—only accessible within the application's own compose network).

---

## Port Configuration

Each exported service declares one or more ports:

```yaml
ports:
  - type: proxied                       # Required: proxied or assigned
    protocol: https                     # Required for proxied: http, https, or future SNI types
    port: 8080                          # Optional: container port (see inference rules below)
    visibility: public                  # Optional: public or protected (documentation only)
```

**Port type constraints**:
- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

**Port inference rules**:
- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

---

## Port Assignment

### Global State

**Location**: `~/.config/contrail/state.yaml`

This file tracks port assignments for `assigned` type ports across all workspaces, plus an inventory of port availability.

```yaml
assigned_ports:
  dev:
    app-one:
      db: 5432
    app-two:
      db: 5433                          # Incremented because 5432 was taken
      cache: 6379
  review:
    app-one:
      db: 5434                          # Different workspace, different port

port_inventory:
  5432:
    status: assigned
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:
      workspace: dev
      application: app-one
      exported_service: db
```

### Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

### Port Conflict at Startup

If a previously assigned port has become unavailable when `workspace up` runs:

```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

### Port Availability Checking

`contrail port scan` and `contrail port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state.

This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

### Port Status Transitions

The port inventory tracks status changes over time:

| Transition | When It Occurs |
|------------|----------------|
| `unavailable` → `assigned` | Port became free, Contrail claimed it |
| `assigned` → `released` | Workspace/app removed, port freed |
| `unavailable` → `released` | External process stopped, `contrail port gc` cleaned it up |

---

## Examples

### Simple web service (single port in Compose, inferred)

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

Result: HTTPS proxy to container port 8080. Environment variables will use proxy port 443.

### Database with direct port (assigned port, auto-assigned if unavailable)

```yaml
# application.yaml
exported_services:
  db:
    service: mysql
    ports:
      - type: assigned
        port: 3306
        visibility: protected
```

### Service with both proxy and direct ports

```yaml
# application.yaml
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

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md)
- [Generated Override Files](generated-override-files.md)
- [ADR-0007: Port Type System](../decisions/0007-port-type-system.md)
