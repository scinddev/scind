# Configuration Schemas

This specification defines the **behavioral rules** for Scind configuration. For schema definitions and field references, see [Configuration Reference](../reference/configuration.md).

---

## Design Rationale: Structure vs State

The system uses three schema types, separating structure (configuration) from state (runtime):

| Aspect | Structure (config) | State (runtime) |
|--------|-------------------|-----------------|
| Proxy settings | `proxy.yaml` | - |
| Port assignments | - | `~/.config/scind/state.yaml` |
| What apps exist | `workspace.yaml` | - |
| Available flavors | `application.yaml` | - |
| Active flavor | - | `.generated/state.yaml` or CLI |
| Active branch | - | git working directory |
| Running containers | - | Docker |

This separation ensures configuration files are declarative and version-controllable while runtime state remains ephemeral and machine-specific.

---

## Proxy Behavior

> For the complete proxy schema and field definitions, see [Configuration Reference - Proxy Configuration](../reference/configuration.md#proxy-configuration).

### Lifecycle

- `proxy init`: Creates the directory structure and configuration files
- `proxy up`: Starts the Traefik container (creates `scind-proxy` network if needed)
- `proxy down`: Stops the Traefik container
- `workspace up`: Ensures the `scind-proxy` network always, and runs `proxy up` when the applications being brought up declare at least one proxied export (see Auto-Start below)

### Proxy Auto-Start

| Field (`proxy.yaml`) | Type | Default |
|----------------------|------|---------|
| `proxy.auto_start` | boolean | `true` |

`workspace up` (and therefore `app up`) starts the shared proxy only when it is needed:

- The `scind-proxy` **network** is ensured **unconditionally** — generated overrides declare it `external: true`, so applications need it whether or not Traefik runs.
- The proxy **container** starts only when at least one application in this invocation declares a **proxied** export. An assigned-ports-only `up` starts no proxy.
- `auto_start: false` disables the side-effect start entirely, for users who run their own proxy. The network is still ensured. If proxied exports are present and the proxy is not running, Scind **warns** and continues rather than failing `up` — with `auto_start: false` the user owns proxy lifecycle. This is the sole exception to the fail-fast proxy/network gate.
- `auto_start` lives in `proxy.yaml`, a global (per-user) file, so it applies uniformly across every workspace on the machine — there is no per-workspace override.

See [Workspace Lifecycle — Proxy Start Conditions](./workspace-lifecycle.md#proxy-start-conditions). Both refinements come from the Xcind proof-of-concept.

### Recovery

If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

### Proxy Host Ports

The shared proxy binds host ports for its HTTP, HTTPS, and dashboard entrypoints. These are configurable so the proxy can run on machines where `80`/`443` are already occupied — a common developer situation — without editing generated files:

| Field (`proxy.yaml`) | Environment variable | Default |
|----------------------|----------------------|---------|
| `proxy.http_port` | `SCIND_PROXY_HTTP_PORT` | `80` |
| `proxy.https_port` | `SCIND_PROXY_HTTPS_PORT` | `443` |
| `proxy.dashboard.port` | `SCIND_PROXY_DASHBOARD_PORT` | `8080` |

**Discovery interaction**: When an HTTP or HTTPS entrypoint runs on a **non-default** port, service-discovery values must reflect it — the proxied `_PORT` carries the non-default port and `_URL` includes it (e.g., `https://dev-frontend-web.scind.test:8443`). When the entrypoint runs on the standard port (`80`/`443`), the port is omitted from `_URL`. See [Environment Variables — Configurable Proxy Host Ports](./environment-variables.md#configurable-proxy-host-ports). The Traefik static-configuration side of this change is covered by the proxy infrastructure spec.

### TLS Mode Behavior

| Mode | Behavior |
|------|----------|
| `auto` | Uses mkcert if available to generate locally-trusted certificates; falls back to Traefik's default self-signed certificate (browser warnings) |
| `custom` | Uses user-provided certificate and key files (for enterprise CA or manually generated certs) |
| `disabled` | HTTP only, no HTTPS entrypoint (not recommended for production-like testing) |

**Certificate Setup by Mode**:

- **auto with mkcert**: Run `mkcert -install` once per machine to add the local CA to your trust store, then `mkcert "*.scind.test"` to generate a wildcard certificate. Scind will detect and use these automatically.
- **custom (enterprise CA)**: Obtain a wildcard certificate signed by your enterprise CA for `*.scind.test` (or your configured domain). Place the cert and key files at the configured paths.
- **auto without mkcert**: Traefik serves its default self-signed certificate. Browsers will show security warnings.

---

## Workspace Registry Behavior

> For the registry schema, see [Configuration Reference - Workspace Registry](../reference/configuration.md#workspace-registry).

### Registration Rules

- `workspace init` automatically registers the workspace, failing if the name is already registered to a different path
- `workspace list` reads the registry and optionally validates entries still exist
- `workspace prune` removes stale entries (paths that no longer contain `workspace.yaml`)

### Fallback Recovery

If the registry file is missing or corrupted, Scind can reconstruct it by querying Docker for containers with `workspace.name` and `workspace.path` labels. This provides resilience against accidental deletion of `~/.config/scind/workspaces.yaml`.

---

## Port Assignment Behavior

> For the global state schema, see [Configuration Reference - Global State](../reference/configuration.md#global-state).

### Assignment Algorithm

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

### Port Conflict at Startup

If a previously assigned port has become unavailable (e.g., taken by an external process) when `workspace up` runs, Scind fails with a clear error:

```
Error: Port conflict detected for frontend

Port 5432 is assigned to frontend/postgres but is no longer available.
Another process may be using this port.

To resolve:
  scind port scan       # Check which ports are conflicting
  scind port release 5432   # Release the conflicting assignment
  scind generate --force    # Regenerate with new port assignment
```

### Port Status Transitions

- `unavailable` → `assigned`: Port became free, Scind claimed it
- `assigned` → `released`: Workspace/app removed, port freed
- `unavailable` → `released`: External process stopped, `scind port gc` cleaned it up

### Availability Checking

`scind port scan` and `scind port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

---

## Workspace Configuration Behavior

> For the workspace schema, see [Configuration Reference - Workspace Configuration](../reference/configuration.md#workspace-configuration).

### Conventions

- Application name = directory path (e.g., `frontend` → `./{workspace}/frontend/`)
- Network name defaults to `{workspace.name}-internal`

**Note**: Branch (`ref`) and flavor are runtime state, not configuration.

---

## Template Resolution Behavior

> For template variable definitions, see [Configuration Reference - Template Variables](../reference/configuration.md#template-variables).
>
> Apex templates (`hostname-apex`, `alias-apex`) are resolved using the same mechanism but intentionally exclude the `%EXPORTED_SERVICE%` variable. See [Configuration Reference - Default Templates](../reference/configuration.md#default-templates).

### Resolution Timing

Template variables are resolved at **generation time** (when `workspace generate` or `workspace up` runs). The resolved values are written into the generated override files.

### Flavor Change Handling

When `scind flavor set FLAVOR` is executed, it:
1. Updates `.generated/state.yaml` with the new flavor
2. Immediately regenerates the affected application's override file
3. If the application is currently running, displays a warning:
   ```
   Warning: Application "app-name" is currently running.
   The new flavor has been applied to the configuration, but running
   containers still use the previous flavor.

   To apply the flavor change:
     scind app restart -a app-name
   ```

This ensures override files always reflect the current flavor without requiring a separate `generate` step.

### Running Application Considerations

Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `scind up` to start new services |
| Flavor removes services | Services still running but not in override | Run `scind up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `scind app restart` to pick up changes |

### Orphaned Service Handling

When `scind up` is run after a flavor change that removes services, Scind passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

---

## Application Configuration Behavior

> For the application schema, see [Configuration Reference - Application Configuration](../reference/configuration.md#application-configuration-service-contract).

### Exported Service Mapping

Each key in `exported_services` is the "exported service name" used for hostname generation, network aliases, and environment variables. By default, this key maps to a Compose service of the same name.

Use the `service:` property when the exported name differs from the Compose service name.

### Application Env File Injection

Applications may declare dotenv files at the application level through two fields whose behaviors map to Docker Compose's two distinct env-file scopes. The distinction is a common source of confusion and is therefore modeled as two explicit fields rather than one:

| Field | Compose mechanism | When applied | Reaches container processes? |
|-------|-------------------|--------------|------------------------------|
| `compose_env_files` | `--env-file` | While Compose parses the YAML | No — interpolation only |
| `app_env_files` | `env_file:` on every generated service | Container runtime | Yes |

- **`compose_env_files`** are passed to Compose as `--env-file` and drive `${VAR}` **interpolation of the Compose YAML**. They influence how the file is parsed; they do **not** inject variables into running containers.
- **`app_env_files`** cause Scind to add an `env_file:` entry referencing each listed file to **every service** in the generated override, so the variables are present **inside running containers**.

A value needed for both YAML interpolation and container runtime must appear in both lists. For field definitions see [Configuration Reference — Application Env Files](../reference/configuration.md#application-env-files). This two-field model was validated by the Xcind proof-of-concept.

### Default Compose File Resolution

`flavors` and `compose_files` are optional. When an application declares neither and the active flavor is the literal `default` (i.e., `default_flavor` is unset or set to `default`), that flavor resolves to the existence-filtered candidate list `compose.yaml`, `compose.yml`, `docker-compose.yaml`, `docker-compose.yml`, and the env-file default `.env` (applied to `compose_env_files`, the interpolation scope). This mirrors Docker Compose's own default file discovery, so an application behaves under Scind the way it already behaves under plain `docker compose`.

- The default applies **only when nothing is declared and the active flavor is `default`**. Any `flavors` or `compose_files` declaration governs completely; defaults are never merged into a declared list. A non-`default` active flavor with no declaration (for example, `default_flavor: full` with no `flavors:` block) **fails** rather than falling back to convention defaults, with its own error naming the flavor and where it came from — not the missing-file error, since an undeclared flavor has no file list. See [Generation Logic](./workspace-lifecycle.md#generation-logic-workspace-generate) step 3.
- If nothing is declared and **no** candidate file exists, generation fails with an error naming the four candidates and the directory searched.
- Reporting commands (`app show`, `app diagnose`) render the resolved list and mark it convention-derived, so the application still describes itself.

See [Configuration Reference — Default Compose File Resolution](../reference/configuration.md#default-compose-file-resolution). Adopted from the Xcind proof-of-concept.

### Port Type Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

### Primary Export Designation

An exported service can be marked as the application's primary export by adding `primary: true`:

- **Type**: Boolean, optional, defaults to `false`
- **At most one** exported service per application may be marked `primary: true`

**Apex eligibility is scoped to proxied exports.** Apex hostnames only ever apply to proxied exports, so only proxied exports are candidates for implicit primary. Assigned exports do not count toward implicit-primary eligibility (they have no hostname):

- **Implicit primary**: When an application has exactly **one proxied** exported service, it is implicitly primary and apex-eligible (no annotation needed). Assigned exports present alongside it do not disqualify this — a `web` (proxied) + `db` (assigned) application receives an apex with zero annotation.
- **Hybrid selection when several proxied exports exist**: An explicit `primary: true` always wins. When multiple proxied exports exist and **none** is marked primary, Scind falls back to **positional** selection — the first-declared proxied export becomes primary and an apex is still generated. (Prior behavior generated no apex in this case; the hybrid rule always emits one.) Explicit `primary: true` is therefore required only to override the positional default when 2+ proxied exports compete.
- **Assigned primary**: An assigned export may still be marked `primary: true` explicitly to receive the apex internal alias, but it is never implicitly primary and never yields an apex hostname.
- **Validation error**: If more than one exported service is marked `primary: true`, Scind emits a validation error at generation time.

**Ordering requirement**: Positional fallback requires `exported_services` to have a defined declaration order. If `exported_services` is modeled as an unordered YAML map, adopt a documented first-declared rule (or model it as a sequence) so that "first proxied export" is well-defined. See [ADR-0013](../decisions/0013-apex-url-primary-designation.md).

The primary export receives:

- An **apex hostname** (proxied types only): `{workspace}-{application}.{domain}`
- An **apex internal alias** (all types): `{application}`
- **Apex environment variables** (proxied types only): `SCIND_{APPLICATION}_APEX_*`
- **Apex Docker labels** (proxied types only): `scind.apex.*`

For assigned-port primary exports, only the apex internal alias is created (no hostname, no apex environment variables, no apex labels).

**Apex opt-out**: An application-level `apex` boolean in `application.yaml` (default `true`) declines the apex entirely. With `apex: false`, Scind generates no apex hostname, no apex internal alias, no apex Traefik router, no `scind.apex.*` labels, and no `SCIND_{APPLICATION}_APEX_*` environment variables; `apex` and `apex_host` report `null` in the [JSON introspection contract](../reference/cli.md#json-introspection-contract). The opt-out is application-level because there is exactly one apex per application. It does **not** conflict with `primary: true` — an export may still be marked primary, and the designation simply yields nothing apex-related — so this combination is not a validation error.

See [ADR-0013](../decisions/0013-apex-url-primary-designation.md) for the design rationale.

### Port Inference Rules

- If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
- If the Compose service has multiple ports, `port:` must be explicitly specified
- For Compose port mappings like `"80:8080"`, the container port (`8080`) is used

---

## Flavor Resolution Order

1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Related Documentation

- [Configuration Reference](../reference/configuration.md) - Schema definitions and field reference
- [Port Types Spec](./port-types.md) - Detailed port type behaviors
- [Generated Override Files Spec](./generated-override-files.md) - Override file generation rules
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md) - Design rationale
- [ADR-0013: Apex URL Primary Designation](../decisions/0013-apex-url-primary-designation.md) - Primary designation design
