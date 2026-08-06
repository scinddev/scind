## Environment Variable Injection

All exported services receive environment variables for service discovery. This enables applications to reference other services without hardcoding hostnames.

> **Schema validated by Xcind.** The `SCIND_{APP}_{EXPORT}_{SUFFIX}` service-discovery injection schema described here — env-safing rules, HTTPS-default base variables, protocol-specific `_HTTP_*`/`_HTTPS_*` variables, and apex variables on the primary export — is build-validated by the Xcind Bash proof-of-concept, which generates and injects these variables end-to-end. The one exception is the `_HOST`/`_PORT`/`_HOST_PORT` coherence contract for assigned exports, which is a correction surfaced by that same work (see [Assigned-Export Discovery Contract](#assigned-export-discovery-contract-three-values) below).

### Naming Convention

Environment variables use a `SCIND_` prefix to avoid conflicts with application-defined variables.

**Name transformation**: Hyphens in application and exported service names are converted to underscores, and names are uppercased (e.g., `shared-db` becomes `SHARED_DB`, `web-debug` becomes `WEB_DEBUG`).

**Base variables** (always generated for each exported service):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_HOST_PORT={host_port}  # Only for assigned types
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

**Protocol-specific variables** (generated for each proxied protocol):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

**Apex variables** (generated only for proxied primary exports):
```
SCIND_{APPLICATION}_APEX_HOST={apex_hostname}
SCIND_{APPLICATION}_APEX_PORT={port}
SCIND_{APPLICATION}_APEX_SCHEME={scheme}
SCIND_{APPLICATION}_APEX_URL={url}
```

Apex environment variables follow the pattern `SCIND_{APPLICATION}_APEX_{SUFFIX}` — the exported service name segment is omitted. These are only injected when the application has a primary export with proxied ports. Assigned-port primary exports do not generate apex environment variables (they have no hostname). See [ADR-0013](../decisions/0013-apex-url-primary-designation.md) for primary designation rules.

### Variable Generation Rules

**For `proxied` type ports**:
- `*_HOST` contains the fully qualified proxied hostname (e.g., `dev-frontend-web.scind.test`)
- `*_PORT` contains the proxy port (443 for HTTPS, 80 for HTTP)—**not** the container port
- `*_SCHEME` and `*_URL` are generated
- Protocol-specific variables (`*_HTTPS_*`, `*_HTTP_*`) are also generated

#### Assigned-Export Discovery Contract (Three Values)

**For `assigned` type ports**, discovery requires **three** values, not two:

- `SCIND_{APP}_{EXPORT}_HOST` — the **in-network alias** (e.g., `shared-db-db`), used for container-to-container access
- `SCIND_{APP}_{EXPORT}_PORT` — the **container port**, the port the service listens on inside the network; pairs with `*_HOST`
- `SCIND_{APP}_{EXPORT}_HOST_PORT` — the **allocated host-published port**, used to reach the service from the host or via `host.docker.internal`
- No `*_SCHEME` or `*_URL` variables
- No protocol-specific variables

**Why three values.** Scind previously paired an assigned export's in-network `_HOST` with a single `_PORT`. That pairing is self-contradictory whenever the host-published port differs from the container port — the normal outcome after a conflict reassignment (e.g., a Postgres container listening on `5432` but published on `5433` because `5432` was taken). With a single `_PORT`:

- If `_PORT` held the container port, `${_HOST}:${_PORT}` resolved correctly in-network but the host-published mapping was undiscoverable.
- If `_PORT` held the host port, `${_HOST}:${_PORT}` resolved to the *wrong* port in-network (`shared-db-db:5433`, where nothing listens), and the alias was still unreachable from the host.

Splitting the container port (`_PORT`, paired with the in-network `_HOST`) from the host-published port (`_HOST_PORT`, paired with `127.0.0.1` / `host.docker.internal`) makes both access paths correct simultaneously. This coherence fix was surfaced and validated by the Xcind proof-of-concept — see [Xcind ADR-0018].

| Type | Protocol | `*_HOST` | `*_PORT` | `*_HOST_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|---------------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | — | `https` | ✓ | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | — | `http` | ✓ | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | — | `https` | ✓ | Both |
| `assigned` | - | Internal alias | Container port | Host-published port | ✗ | ✗ | ✗ |
| `proxied` (primary) | `https` | Apex hostname | 443 | — | `https` | ✓ | `*_APEX_*` |
| `assigned` (primary) | - | *(no apex env vars)* | - | - | - | - | *(alias only)* |

[Xcind ADR-0018]: Xcind proof-of-concept ADR validating the three-value assigned-export discovery contract (`_HOST` in-network alias, `_PORT` container port, `_HOST_PORT` host-published port).

**HTTPS-default rationale**: When both HTTP and HTTPS are configured for an exported service, base variables (`*_PORT`, `*_SCHEME`, `*_URL`) default to HTTPS (port 443) following security-by-default principles. Applications should prefer HTTPS for service-to-service communication. Use protocol-specific variables (`*_HTTP_*`) when HTTP is explicitly required.

### Usage in Applications

```php
// PHP example - using URL directly for proxied services
$apiUrl = getenv('SCIND_BACKEND_API_URL') ?: 'https://backend-api.scind.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// PHP example - building connection for assigned port services
$dbHost = getenv('SCIND_SHARED_DB_DB_HOST') ?: 'shared-db-db';
$dbPort = getenv('SCIND_SHARED_DB_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

```javascript
// Node.js example - using URL directly
const apiUrl = process.env.SCIND_BACKEND_API_URL || 'https://backend-api.scind.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Node.js example - building connection manually
const dbHost = process.env.SCIND_SHARED_DB_DB_HOST || 'shared-db-db';
const dbPort = process.env.SCIND_SHARED_DB_DB_PORT || '5432';
```

---

## Host-View Environment

The service-discovery variables above are **container-flavored**: they describe how one container reaches another (in-network alias + container port) or how a browser reaches a proxied service (routable hostname). But developers routinely run **host** processes that must reach the same services — mise/direnv-loaded shells, test runners, language servers, debuggers, and one-off scripts. From the host, an assigned export does not live at its in-network `host:container_port`; it lives at `127.0.0.1:<host_port>`. This is exactly the [three-value assigned-export contract](#assigned-export-discovery-contract-three-values): the container view uses `_HOST` + `_PORT`, the host view uses `127.0.0.1` + `_HOST_PORT`.

The consequence is that a single committed connection string cannot resolve correctly in both places. A value like:

```
DATABASE_URL=postgres://app@${SCIND_SHARED_DB_DB_HOST}:${SCIND_SHARED_DB_DB_PORT}/app
```

is correct inside a container (`shared-db-db:5432`) and wrong on the host (`shared-db-db` is not resolvable there, and the host-published port may differ).

**Host-view env file.** Scind offers an **opt-in** host-view environment file that renders the *same* discovery computation with **host-flavored** values:

- **Assigned exports** resolve to `127.0.0.1` + the allocated host-published port (`_HOST_PORT`).
- **Proxied and apex exports** keep their routable hostname and URL — these already resolve identically from host and container.

Because the host-view file is produced by the same computation that drives container injection — not a separate hand-maintained file — the two views cannot drift. The container injection and the host-view file are two renderings of one resolved topology (the manifest).

### Layering Guarantees

The host-view file is designed to compose safely with both host shells and Compose, relying on two language-agnostic layering facts:

1. **Dotenv loaders never override a real OS environment variable.** Well-behaved dotenv implementations (mise, direnv, and every mainstream language's dotenv library) treat a variable already present in the OS environment as authoritative and skip it. A developer who has intentionally exported `DATABASE_URL` in their shell keeps that value; the host-view file only fills in what is unset.
2. **Compose `environment:` beats `env_file:`.** When the generated override injects a variable via `environment:` (the container-flavored value), it wins over any value the same key might carry through an `env_file:` entry. This means a host-view file sourced on the host and a container-flavored value injected into a container do not collide even if both define the same key.

### Opt-In and Write Mode

The host-view file is **never** generated implicitly and is **never** inferred from a filename. Enabling it is an explicit choice, and Scind must be told how it may write:

- **own** — Scind owns the file end-to-end and regenerates it wholesale.
- **block** — Scind manages only a delimited block within an otherwise developer-owned file, leaving surrounding content untouched.

Making the mode explicit prevents Scind from clobbering a developer's own dotenv file merely because it shares a conventional name (e.g., `.env`).

### Scoping Note

Host-side process access sits at the edge of Scind's runtime model, which is otherwise Compose/container-centric. The capability is recorded here as an accepted design because host processes are a first-class part of the development loop the tool targets; whether every Scind runtime enables it by default is a separate product decision. See [ADR-0015: Host-View Environment Symmetry](../decisions/0015-host-view-environment-symmetry.md) for the full rationale. The Xcind proof-of-concept implements this pattern (Xcind ADR-0020) and validates both layering guarantees end-to-end.

---

## Configuration Environment Variables

These environment variables configure Scind itself (distinct from service discovery variables injected into containers).

| Variable | Description | Default |
|----------|-------------|---------|
| `TRAEFIK_IMAGE` | Traefik Docker image to use | `traefik:v3.2.3` |
| `SCIND_CONFIG_DIR` | Configuration directory (declarative structure) | `$XDG_CONFIG_HOME/scind` (fallback `~/.config/scind`) |
| `SCIND_STATE_DIR` | State file directory (ephemeral runtime state) | `$XDG_STATE_HOME/scind` (fallback `~/.local/state/scind`) |
| `SCIND_PROXY_HTTP_PORT` | Host port for the shared proxy's HTTP entrypoint | `80` |
| `SCIND_PROXY_HTTPS_PORT` | Host port for the shared proxy's HTTPS entrypoint | `443` |
| `SCIND_PROXY_DASHBOARD_PORT` | Host port for the Traefik dashboard | `8080` |

These can also be set in `proxy.yaml` configuration (see [Proxy Configuration](../reference/configuration.md#proxy-configuration)).

### XDG Base Directory Alignment

`SCIND_CONFIG_DIR` and `SCIND_STATE_DIR` follow the [XDG Base Directory specification](https://specifications.freedesktop.org/basedir-spec/latest/), which separates **config** (declarative, version-controllable structure) from **state** (ephemeral, machine-specific runtime data). Defaulting both under `~/.config/scind` would collapse ephemeral state into config-home and contradict the structure-vs-state separation established in [ADR-0005](../decisions/0005-structure-vs-state-separation.md). State therefore defaults to `$XDG_STATE_HOME/scind` (`~/.local/state/scind`) while config remains under `$XDG_CONFIG_HOME/scind` (`~/.config/scind`).

### Configurable Proxy Host Ports

The shared proxy binds host ports `80`, `443`, and `8080` by default. On developer machines these are frequently already occupied (another proxy, a system service, a previously-running stack). `SCIND_PROXY_HTTP_PORT`, `SCIND_PROXY_HTTPS_PORT`, and `SCIND_PROXY_DASHBOARD_PORT` (equivalently `proxy.http_port`, `proxy.https_port`, `proxy.dashboard.port` in `proxy.yaml`) let the proxy run on alternate host ports **without editing generated files**.

When an HTTP or HTTPS entrypoint runs on a **non-default** port, that port is reflected in service-discovery values: the proxied `_PORT` carries the non-default port and `_URL` includes it (e.g., `https://dev-frontend-web.scind.test:8443`). When the entrypoint runs on the standard port (`80`/`443`), the port is omitted from `_URL` as usual.

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md) - Configuration file formats
- [Generated Override Files](generated-override-files.md) - How variables are injected
- [ADR-0013: Apex URL Primary Designation](../decisions/0013-apex-url-primary-designation.md) - Primary designation design
- [ADR-0015: Host-View Environment Symmetry](../decisions/0015-host-view-environment-symmetry.md) - Host-flavored discovery variables
