<!-- Migrated from specs/scind-technical-spec.md:874-954 -->
<!-- Extraction ID: spec-docker-labels -->

## Docker Labels

Scind uses Docker labels for workspace discovery, service routing, and external tool integration. All labels use the `scind.` namespace prefix with kebab-case for multi-word segments.

### Context Labels

Applied to all application containers for workspace discovery and registry reconstruction:

| Label | Description | Example |
|-------|-------------|---------|
| `scind.workspace.name` | Workspace identifier | `dev` |
| `scind.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `scind.app.name` | Application identifier | `app-one` |
| `scind.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/app-one` |

### Export Labels

Applied to containers with exported services. Labels are keyed by export name for consistency, supporting multiple exports per container.

**Proxied exports** (HTTP/HTTPS through Traefik):
```
scind.export.{name}.host={hostname}
scind.export.{name}.proxy.http.visibility={public|protected|private}
scind.export.{name}.proxy.http.url={url}
scind.export.{name}.proxy.https.visibility={public|protected|private}
scind.export.{name}.proxy.https.url={url}
```

**Assigned port exports** (direct port mapping):
```
scind.export.{name}.host={hostname}
scind.export.{name}.port.{internal-port}.visibility={public|protected|private}
scind.export.{name}.port.{internal-port}.assigned={external-port}
```

**Example** — web service with proxied HTTP/HTTPS and debug port:
```yaml
labels:
  # Context
  - "scind.workspace.name=dev"
  - "scind.workspace.path=/Users/beau/workspaces/dev"
  - "scind.app.name=app-one"
  - "scind.app.path=/Users/beau/workspaces/dev/app-one"
  # Proxied export: web
  - "scind.export.web.host=dev-app-one-web.scind.test"
  - "scind.export.web.proxy.http.visibility=public"
  - "scind.export.web.proxy.http.url=http://dev-app-one-web.scind.test"
  - "scind.export.web.proxy.https.visibility=public"
  - "scind.export.web.proxy.https.url=https://dev-app-one-web.scind.test"
  # Assigned port export: debug (uses internal alias, not proxied hostname)
  - "scind.export.debug.host=app-one-debug"
  - "scind.export.debug.port.9000.visibility=protected"
  - "scind.export.debug.port.9000.assigned=9003"
```

### Proxy Container Labels

Applied to the Scind-managed Traefik proxy container:

| Label | Description | Example |
|-------|-------------|---------|
| `scind.managed` | Indicates Scind manages this container | `true` |
| `scind.component` | Component type | `proxy` |

### External Tool Integration

External tools can discover Scind workspaces and services by querying Docker labels:

```bash
# Find all Scind-managed containers
docker ps --filter "label=scind.workspace.name" --format "{{.Names}}"

# Find all containers for a specific workspace
docker ps --filter "label=scind.workspace.name=dev" --format "{{.Names}}"

# Get workspace paths for registry reconstruction
docker inspect --format '{{index .Config.Labels "scind.workspace.path"}}' $(docker ps -q --filter "label=scind.workspace.name")
```
