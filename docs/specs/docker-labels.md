## Docker Labels

Scind uses Docker labels for workspace discovery, service routing, and external tool integration. All labels use the `scind.` namespace prefix with kebab-case for multi-word segments.

### Context Labels

Applied to all application containers for workspace discovery and registry reconstruction:

| Label | Description | Example |
|-------|-------------|---------|
| `scind.workspace.name` | Workspace identifier | `dev` |
| `scind.workspace.path` | Absolute path to workspace directory | `/Users/beau/workspaces/dev` |
| `scind.app.name` | Application identifier | `frontend` |
| `scind.app.path` | Absolute path to application directory | `/Users/beau/workspaces/dev/frontend` |

### Export Labels

Applied to containers with exported services. Labels are keyed by export name for consistency, supporting multiple exports per container.

**Proxied exports** (HTTP/HTTPS through Traefik):
```
scind.export.{name}.host={hostname}
scind.export.{name}.url={preferred-url}
scind.export.{name}.proxy.http.visibility={public|protected}
scind.export.{name}.proxy.http.url={url}
scind.export.{name}.proxy.https.visibility={public|protected}
scind.export.{name}.proxy.https.url={url}
```

The `scind.export.{name}.url` label is the **canonical preferred-scheme URL** for the export: `https://…` when an HTTPS router serves the export, otherwise `http://…`. It is additive to the per-protocol `scind.export.{name}.proxy.{http,https}.url` labels above — those still advertise each individual scheme — so consumers get "the URL to use" directly, without re-deriving scheme preference from which per-protocol labels happen to be present.

*Validated by the Xcind proof-of-concept.*

**Assigned port exports** (direct port mapping):
```
scind.export.{name}.host={hostname}
scind.export.{name}.port.{internal-port}.visibility={public|protected}
scind.export.{name}.port.{internal-port}.assigned={external-port}
```

**Example** — web service with proxied HTTP/HTTPS and debug port:
```yaml
labels:
  # Context
  - "scind.workspace.name=dev"
  - "scind.workspace.path=/Users/beau/workspaces/dev"
  - "scind.app.name=frontend"
  - "scind.app.path=/Users/beau/workspaces/dev/frontend"
  # Proxied export: web
  - "scind.export.web.host=dev-frontend-web.scind.test"
  - "scind.export.web.proxy.http.visibility=public"
  - "scind.export.web.proxy.http.url=http://dev-frontend-web.scind.test"
  - "scind.export.web.proxy.https.visibility=public"
  - "scind.export.web.proxy.https.url=https://dev-frontend-web.scind.test"
  # Assigned port export: debug (uses internal alias, not proxied hostname)
  - "scind.export.debug.host=frontend-debug"
  - "scind.export.debug.port.9000.visibility=protected"
  - "scind.export.debug.port.9000.assigned=9003"
```

### Apex Labels

Applied to the container running the primary exported service, when that export is proxied. These labels provide apex (application-level) hostname discovery:

```
scind.apex.host={apex_hostname}
scind.apex.url={preferred_apex_url}
scind.apex.proxy.{protocol}.url={apex_url}
```

As with exports, `scind.apex.url` is the canonical preferred-scheme apex URL (HTTPS when an HTTPS apex router exists, else HTTP), additive to the per-protocol `scind.apex.proxy.{protocol}.url` labels.

**Example** — frontend with web as implicit primary:
```yaml
labels:
  - "scind.apex.host=dev-frontend.scind.test"
  - "scind.apex.url=https://dev-frontend.scind.test"
  - "scind.apex.proxy.https.url=https://dev-frontend.scind.test"
```

These labels only appear on proxied primary exports. Assigned-port primary exports receive the apex internal alias but no apex Docker labels (since there is no hostname to advertise). An application that sets `apex: false` emits no `scind.apex.*` labels at all.

See [ADR-0013](../decisions/0013-apex-url-primary-designation.md) for primary designation and opt-out rules.

### Proxy Container Labels

Applied to the Scind-managed Traefik proxy container:

| Label | Description | Example |
|-------|-------------|---------|
| `scind.managed` | Indicates Scind manages this container | `true` |
| `scind.component` | Component type | `proxy` |

### Traefik Routing Labels

These labels are added to containers to configure Traefik routing. They are generated automatically by Scind based on exported service configuration.

| Label | Description | Example |
|-------|-------------|---------|
| `traefik.enable` | Exposes container to Traefik | `true` |
| `traefik.docker.network` | Which network Traefik routes on | `scind-proxy` |
| `traefik.http.routers.{name}.rule` | Routing rule (Host matcher) | `Host(\`dev-app-web.scind.test\`)` |
| `traefik.http.routers.{name}.entrypoints` | Entry points to use | `websecure` |
| `traefik.http.routers.{name}.tls` | Enable TLS | `true` |
| `traefik.http.routers.{name}.service` | Bind this router to its own service | `dev-app-web-https` |
| `traefik.http.services.{name}.loadbalancer.server.port` | Container port | `8080` |

#### Network Selection and Router→Service Binding

Two labels are load-bearing for correctness under Scind's two-layer networking:

- **`traefik.docker.network={proxy-network}`** (e.g. `scind-proxy`) is **required**. Under [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md), each exported container is attached to more than one network (its application-internal network plus the shared proxy network). If the network label is omitted, Traefik cannot unambiguously determine which of the container's IPs to route to and may pick the wrong one, causing intermittent mis-routing. The label pins routing to the shared proxy network.
- **`traefik.http.routers.{name}.service={name}`** explicitly binds each router to its own loadbalancer service. With multiple routers/services on one container, relying on Traefik's implicit service selection is fragile; the explicit binding is deterministic.

*Validated by the Xcind proof-of-concept.*

#### Router Naming Convention

Router names follow the pattern: `{workspace}-{application}-{service}-{protocol}`

Example: `dev-frontend-web-https`

**Apex router names** omit the service segment: `{workspace}-{application}-{protocol}`

Example: `dev-frontend-https`

#### Example Labels

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.dev-frontend-web-https.rule=Host(`dev-frontend-web.scind.test`)"
  - "traefik.http.routers.dev-frontend-web-https.entrypoints=websecure"
  - "traefik.http.routers.dev-frontend-web-https.tls=true"
  - "traefik.http.services.dev-frontend-web-https.loadbalancer.server.port=3000"
```

#### HTTP→HTTPS Redirect Labels

When an export requires TLS (`tls: require`), the HTTP→HTTPS redirect is expressed with a shared `redirectscheme` middleware and a redirect-only HTTP router:

```yaml
labels:
  # Shared redirect middleware DEFINITION — emitted on EVERY rendered service block (see note)
  - "traefik.http.middlewares.scind-redirect-https.redirectscheme.scheme=https"
  # Redirect-only HTTP router: matches the host on the web entrypoint, applies the middleware
  - "traefik.http.routers.dev-frontend-web-http.rule=Host(`dev-frontend-web.scind.test`)"
  - "traefik.http.routers.dev-frontend-web-http.entrypoints=web"
  - "traefik.http.routers.dev-frontend-web-http.middlewares=scind-redirect-https"
  - "traefik.http.routers.dev-frontend-web-http.service=noop@internal"
```

**Emit the middleware definition on every service block.** Traefik's Docker provider loads labels only from **running** containers, so a shared middleware defined on only one "first" container vanishes whenever that container is down, leaving every referencing router dangling. The `...middlewares.scind-redirect-https.redirectscheme.scheme` definition label must therefore be repeated on **every** rendered service block that references it; identical repeated definitions are idempotent. The redirect-only router may target `noop@internal` because the request is redirected before any backend is resolved. See [Proxy Infrastructure — HTTP→HTTPS Redirect](proxy-infrastructure.md#httphttps-redirect).

See [Proxy Infrastructure - Dynamic Routing](proxy-infrastructure.md#dynamic-routing) for more details.

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

---

## Related Documents

- [Proxy Infrastructure](proxy-infrastructure.md) - Traefik configuration and routing
- [Generated Override Files](generated-override-files.md) - How labels are generated
- [ADR-0013: Apex URL Primary Designation](../decisions/0013-apex-url-primary-designation.md) - Primary designation design
