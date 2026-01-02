# Docker Labels Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail uses Docker labels for workspace discovery, service routing through Traefik, and external tool integration. All labels use the `contrail.` namespace prefix with kebab-case for multi-word segments.

Labels serve three primary purposes:
1. **Workspace discovery**: Enable reconstruction of workspace registry from running containers
2. **Traefik routing**: Configure reverse proxy routing for proxied services
3. **External tool integration**: Allow dashboards and monitoring tools to discover and categorize services

**Related Documents**:
- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [Naming Conventions](./naming-conventions.md)
- [Generated Override Files](./generated-override-files.md)

---

## Behavior

### Label Generation

Labels are generated during `workspace generate` and written to the generated override files. They are applied when containers start via Docker Compose.

### Label Namespace

All Contrail labels use the `contrail.` prefix to avoid conflicts with other tools. Segments within label names use kebab-case (e.g., `contrail.workspace.name`, not `contrail.workspace_name`).

---

## Data Schema

### Context Labels

Applied to **all application containers** for workspace discovery and registry reconstruction.

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.workspace.name` | Workspace identifier | `dev` |
| `contrail.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `contrail.app.name` | Application identifier | `app-one` |
| `contrail.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

**Example**:
```yaml
labels:
  - "contrail.workspace.name=dev"
  - "contrail.workspace.path=/Users/beau/workspaces/dev"
  - "contrail.app.name=app-one"
  - "contrail.app.path=/Users/beau/workspaces/dev/app-one"
```

### Export Labels

Applied to containers with exported services. Labels are keyed by export name for consistency, supporting multiple exports per container.

#### Proxied Exports (HTTP/HTTPS through Traefik)

Pattern:
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.proxy.http.visibility={public|protected|private}
contrail.export.{name}.proxy.http.url={url}
contrail.export.{name}.proxy.https.visibility={public|protected|private}
contrail.export.{name}.proxy.https.url={url}
```

**Example** — web service with HTTP and HTTPS:
```yaml
labels:
  - "contrail.export.web.host=dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.http.visibility=protected"
  - "contrail.export.web.proxy.http.url=http://dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.https.visibility=public"
  - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"
```

#### Assigned Port Exports (Direct Port Mapping)

Pattern:
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.port.{internal-port}.visibility={public|protected|private}
contrail.export.{name}.port.{internal-port}.assigned={external-port}
```

**Example** — database with assigned port:
```yaml
labels:
  - "contrail.export.db.host=dev-app-one-db.contrail.test"
  - "contrail.export.db.port.5432.visibility=protected"
  - "contrail.export.db.port.5432.assigned=5432"
```

### Traefik Labels

Applied to containers with proxied exported services to configure Traefik routing.

#### Standard Traefik Labels

| Label | Description |
|-------|-------------|
| `traefik.enable=true` | Enable Traefik routing for this container |
| `traefik.http.routers.{name}.rule` | Routing rule (typically Host matcher) |
| `traefik.http.routers.{name}.entrypoints` | Entrypoint (web or websecure) |
| `traefik.http.routers.{name}.tls` | Enable TLS (for websecure) |
| `traefik.http.services.{name}.loadbalancer.server.port` | Container port to route to |

**Example** — HTTPS router for web service:
```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.contrail.test`)"
  - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
  - "traefik.http.routers.dev-app-one-web-https.tls=true"
  - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"
```

#### Router Naming Convention

Router names follow the pattern: `{workspace}-{application}-{exported_service}-{protocol}`

Examples:
- `dev-app-one-web-https`
- `dev-app-one-web-http`
- `dev-app-two-api-https`

### Proxy Container Labels

Applied to the Contrail-managed Traefik proxy container:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.managed` | Indicates Contrail manages this container | `true` |
| `contrail.component` | Component type | `proxy` |

---

## Examples

### Example 1: Complete Label Set for Web Service

A web service with both HTTP and HTTPS proxying and a debug port:

```yaml
services:
  web:
    labels:
      # Context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/Users/beau/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/Users/beau/workspaces/dev/app-one"

      # Traefik HTTPS router
      - "traefik.enable=true"
      - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.contrail.test`)"
      - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
      - "traefik.http.routers.dev-app-one-web-https.tls=true"
      - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"

      # Traefik HTTP router
      - "traefik.http.routers.dev-app-one-web-http.rule=Host(`dev-app-one-web.contrail.test`)"
      - "traefik.http.routers.dev-app-one-web-http.entrypoints=web"
      - "traefik.http.services.dev-app-one-web-http.loadbalancer.server.port=80"

      # Proxied export: web
      - "contrail.export.web.host=dev-app-one-web.contrail.test"
      - "contrail.export.web.proxy.http.visibility=protected"
      - "contrail.export.web.proxy.http.url=http://dev-app-one-web.contrail.test"
      - "contrail.export.web.proxy.https.visibility=public"
      - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"

      # Assigned port export: debug
      - "contrail.export.debug.host=dev-app-one-debug.contrail.test"
      - "contrail.export.debug.port.9229.visibility=protected"
      - "contrail.export.debug.port.9229.assigned=9229"
```

### Example 2: Database with Assigned Port

```yaml
services:
  postgres:
    labels:
      # Context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/Users/beau/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/Users/beau/workspaces/dev/app-one"

      # Assigned port export: db (mapped from postgres service)
      - "contrail.export.db.host=dev-app-one-db.contrail.test"
      - "contrail.export.db.port.5432.visibility=protected"
      - "contrail.export.db.port.5432.assigned=5433"  # Assigned 5433 because 5432 was taken
```

---

## Edge Cases

### Multiple Exports from Same Container

**Scenario**: A container exposes multiple exported services (e.g., web and debug).

**Behavior**: Multiple export label groups are added to the same container.

```yaml
labels:
  # First export
  - "contrail.export.web.host=dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.https.visibility=public"
  - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"

  # Second export
  - "contrail.export.debug.host=dev-app-one-debug.contrail.test"
  - "contrail.export.debug.port.9229.visibility=protected"
  - "contrail.export.debug.port.9229.assigned=9229"
```

### Port Reassignment

**Scenario**: Preferred port was unavailable, assigned a different port.

**Behavior**: The `.assigned` label reflects the actual assigned port, not the preferred port.

**Example**: Database wanted 5432, got 5433:
```yaml
labels:
  - "contrail.export.db.port.5432.assigned=5433"
```

### Registry Reconstruction

**Scenario**: Workspace registry file is missing or corrupted.

**Behavior**: Contrail can reconstruct the registry by querying Docker for containers with `contrail.workspace.name` and `contrail.workspace.path` labels.

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Label collision | GENERATION | `Router name collision: {name} already defined` | Use unique exported service names |
| Invalid label value | GENERATION | `Invalid label value: {label}` | Check for special characters |

---

## External Tool Integration

External tools can discover Contrail workspaces and services by querying Docker labels.

### Discovery Commands

```bash
# Find all Contrail-managed containers
docker ps --filter "label=contrail.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=contrail.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "contrail.workspace.path"}}' \
  $(docker ps -q --filter "label=contrail.workspace.name")

# Find all public proxied services
docker ps --filter "label=contrail.export.web.proxy.https.visibility=public" \
  --format "{{.Names}}"
```

### Label Queries for Dashboards

A dashboard tool could:
1. Query all containers with `contrail.workspace.name` label
2. Group by workspace name
3. For each container, inspect labels matching `contrail.export.*` to discover services
4. Display service URLs from `contrail.export.*.proxy.*.url` labels
5. Filter by visibility using `*.visibility` labels

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
