# Docker Labels Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail uses Docker labels for workspace discovery, service routing, and external tool integration. All labels use the `contrail.` namespace prefix.

---

## Context Labels

Applied to all application containers:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.workspace.name` | Workspace identifier | `dev` |
| `contrail.workspace.path` | Absolute path to workspace | `/Users/beau/workspaces/dev` |
| `contrail.app.name` | Application identifier | `app-one` |
| `contrail.app.path` | Absolute path to application | `/Users/beau/workspaces/dev/app-one` |

---

## Export Labels

Applied to containers with exported services, keyed by export name.

### Proxied Exports (HTTP/HTTPS)

```
contrail.export.{name}.host={hostname}
contrail.export.{name}.proxy.http.visibility={public|protected}
contrail.export.{name}.proxy.http.url={url}
contrail.export.{name}.proxy.https.visibility={public|protected}
contrail.export.{name}.proxy.https.url={url}
```

### Assigned Port Exports

```
contrail.export.{name}.host={hostname}
contrail.export.{name}.port.{internal-port}.visibility={public|protected}
contrail.export.{name}.port.{internal-port}.assigned={external-port}
```

### Example

Web service with proxied HTTP/HTTPS and debug port:

```yaml
labels:
  # Context
  - "contrail.workspace.name=dev"
  - "contrail.workspace.path=/Users/beau/workspaces/dev"
  - "contrail.app.name=app-one"
  - "contrail.app.path=/Users/beau/workspaces/dev/app-one"
  # Proxied export: web
  - "contrail.export.web.host=dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.http.visibility=public"
  - "contrail.export.web.proxy.http.url=http://dev-app-one-web.contrail.test"
  - "contrail.export.web.proxy.https.visibility=public"
  - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"
  # Assigned port export: debug
  - "contrail.export.debug.host=dev-app-one-debug.contrail.test"
  - "contrail.export.debug.port.9000.visibility=protected"
  - "contrail.export.debug.port.9000.assigned=9003"
```

---

## Proxy Container Labels

Applied to the Contrail-managed Traefik proxy:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.managed` | Indicates Contrail management | `true` |
| `contrail.component` | Component type | `proxy` |

---

## External Tool Integration

External tools can discover Contrail workspaces and services:

```bash
# Find all Contrail-managed containers
docker ps --filter "label=contrail.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=contrail.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "contrail.workspace.path"}}' \
  $(docker ps -q --filter "label=contrail.workspace.name")
```

---

## Registry Reconstruction

If `~/.config/contrail/workspaces.yaml` is missing or corrupted, Contrail can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels.

---

## Related Documentation

- [Configuration Schemas Spec](../configuration-schemas/README.md)
- [Port Types Spec](../port-types/README.md)
- [Proxy Infrastructure Spec](../proxy-infrastructure/README.md)
