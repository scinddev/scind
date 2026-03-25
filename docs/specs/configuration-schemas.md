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
- `workspace up`: Automatically runs `proxy up` if proxy is not running

### Recovery

If a user manually edits the proxy configuration and breaks it, `proxy init --force` regenerates the default configuration.

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

### Port Type Constraints

- Each exported service may have at most **one `http`** and **one `https`** proxied port
- Each exported service may have **multiple `assigned`** ports
- If an exported service needs more than one http or https proxy mapping, create separate exported services

### Primary Export Designation

An exported service can be marked as the application's primary export by adding `primary: true`:

- **Type**: Boolean, optional, defaults to `false`
- **At most one** exported service per application may be marked `primary: true`
- **Implicit primary**: When an application has exactly one exported service, it is implicitly primary (no annotation needed)
- **No primary**: When an application has multiple exported services and none is marked `primary: true`, no apex URL is generated
- **Validation error**: If more than one exported service is marked `primary: true`, Scind emits a validation error at generation time

The primary export receives:

- An **apex hostname** (proxied types only): `{workspace}-{application}.{domain}`
- An **apex internal alias** (all types): `{application}`
- **Apex environment variables** (proxied types only): `SCIND_{APPLICATION}_APEX_*`
- **Apex Docker labels** (proxied types only): `scind.apex.*`

For assigned-port primary exports, only the apex internal alias is created (no hostname, no apex environment variables, no apex labels).

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
