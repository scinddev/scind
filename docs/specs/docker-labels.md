<!-- Migrated from specs/contrail-technical-spec.md:874-954 -->
<!-- Extraction ID: spec-docker-labels -->

## Docker Labels

Contrail uses Docker labels for workspace discovery, service routing, and external tool integration. All labels use the `contrail.` namespace prefix with kebab-case for multi-word segments.

### Context Labels

Applied to all application containers for workspace discovery and registry reconstruction:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.workspace.name` | Workspace identifier | `dev` |
| `contrail.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `contrail.app.name` | Application identifier | `app-one` |
| `contrail.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

### Export Labels

Applied to containers with exported services. Labels are keyed by export name for consistency, supporting multiple exports per container.

**Proxied exports** (HTTP/HTTPS through Traefik):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.proxy.http.visibility={public|protected|private}
contrail.export.{name}.proxy.http.url={url}
contrail.export.{name}.proxy.https.visibility={public|protected|private}
contrail.export.{name}.proxy.https.url={url}
```

**Assigned port exports** (direct port mapping):
```
contrail.export.{name}.host={hostname}
contrail.export.{name}.port.{internal-port}.visibility={public|protected|private}
contrail.export.{name}.port.{internal-port}.assigned={external-port}
```

**Example** — web service with proxied HTTP/HTTPS and debug port:
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
  # Assigned port export: debug (uses internal alias, not proxied hostname)
  - "contrail.export.debug.host=app-one-debug"
  - "contrail.export.debug.port.9000.visibility=protected"
  - "contrail.export.debug.port.9000.assigned=9003"
```

### Proxy Container Labels

Applied to the Contrail-managed Traefik proxy container:

| Label | Description | Example |
|-------|-------------|---------|
| `contrail.managed` | Indicates Contrail manages this container | `true` |
| `contrail.component` | Component type | `proxy` |

### External Tool Integration

External tools can discover Contrail workspaces and services by querying Docker labels:

```bash
# Find all Contrail-managed containers
docker ps --filter "label=contrail.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=contrail.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "contrail.workspace.path"}}' $(docker ps -q --filter "label=contrail.workspace.name")
```
