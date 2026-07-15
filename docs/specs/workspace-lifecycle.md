# Workspace Lifecycle Specification

**Version**: 1.0.0
**Date**: 2026-01-07
**Status**: Accepted

---

## Workspace States

Workspaces transition through defined states during their lifecycle. State is determined at runtime by inspecting Docker container status and the presence of generated files.

### State Definitions

| State | Description | Entry Condition |
|-------|-------------|-----------------|
| `dormant` | Defined in `workspace.yaml` but never started | Initial state after `workspace init` |
| `starting` | `workspace up` initiated, containers launching | `workspace up` command executed |
| `running` | All containers up and healthy | Startup complete, health checks passed |
| `stopping` | `workspace down` initiated | `workspace down` command executed |
| `stopped` | Containers stopped, data preserved | Shutdown complete |
| `error` | Operation failed, requires intervention | Any failure during transition |

### State Transitions

```
dormant ──[workspace up]──> starting ──[containers ready]──> running
                                │                              │
                                └──[failure]──> error <────────┘
                                                  │
running ──[workspace down]──> stopping ──> stopped
                                              │
stopped ──[workspace up]──> starting          │
                                              │
stopped ──[workspace destroy]──> (removed) <──┘
```

### Error States and Recovery

- **Starting failure**: Check container logs with `scind-compose logs` or `workspace logs`
- **Runtime failure**: Container exited unexpectedly; check exit codes and logs
- **Recovery**: Run `workspace down` then `workspace up` to reset state

### State Persistence

Current state is determined by:
1. Docker container status (running/stopped/missing)
2. Presence of `.generated/` directory
3. Entry in `~/.config/scind/workspaces.yaml`

State is not explicitly stored; it is inferred from the environment at runtime.

The state table above is **descriptive**, not a stored artifact. The Xcind proof-of-concept confirmed that no persisted lifecycle-state record is needed to drive these transitions — status is inferred at runtime from Docker container status and the presence of generated files.

---

## Operations

### Startup Sequence (`workspace up`)

1. Ensure proxy is running
2. Create workspace network if it doesn't exist
3. Check if override files are stale; regenerate if needed (see Staleness Detection)
4. For each application (or specified apps if `-a` flag used):
   - Resolve active flavor
   - Execute `docker compose up -d --remove-orphans` with compose files + override

**Notes**:
- The `--remove-orphans` flag is always passed to `docker compose up`. This ensures that containers from services removed by flavor changes, manual compose file edits, or renamed services are automatically stopped and removed.
- The workspace network is always created if it doesn't exist, even when starting a single application with `-a`. This ensures cross-application communication is available when other apps are started later.

**Network and proxy creation failures**: Steps 1–2 create resources (the proxy, the workspace-internal network) that the generated overrides declare as `external: true` (see [Generated Override Files](./generated-override-files.md)). Because the override marks these networks external, Docker Compose will not create them itself and instead fails opaquely at container start if they are absent. Scind must therefore treat a failure to ensure the proxy or the workspace network as a first-class error: emit a diagnostic that names the specific resource (proxy, `{workspace}-internal` network) and includes the underlying Docker error, and **fail `up`** before invoking `docker compose` rather than letting compose surface an unattributable "network not found" error. Validated by the Xcind proof-of-concept.

### Staleness Detection

Scind uses **mtime comparison** to determine if generated override files need to be regenerated. Override files are considered stale if any of the following source files have a newer modification time than the generated override:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`, `docker-compose.worker.yaml`)
- Global assigned-ports state (`~/.config/scind/state.yaml`) — an assigned host port may have changed since the last generation (see [Port Assignment Rules](./state-management.md#port-assignment-rules))
- The generator/schema version — upgrading the Scind binary to a version with a different override/manifest schema must force regeneration even when no source file's mtime changed

**Additional staleness inputs**:
- **Assigned-port revalidation**: before treating overrides as up-to-date, revalidate the assigned host ports recorded in global state. If a recorded port is no longer valid/available, the assigned content is stale regardless of source mtimes.
- **Completeness**: an incomplete generation (missing files, missing completeness marker — see below) is treated as stale and regenerated.

**Behavior**:
- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Touch a file accidentally? Use `--force` to ensure clean state

**Note**: mtime comparison is simple and fast but may trigger unnecessary regeneration if files are touched without content changes. The `--force` flag provides explicit control when needed.

#### Config-Derived vs. Live-State-Derived Content

Staleness distinguishes two classes of generated content:

- **Config-derived (cacheable)**: hostnames, aliases, Traefik labels, and other values that are a pure function of `workspace.yaml`, `application.yaml`, and flavor state. The mtime check above fully governs these — if no source changed, they need not be rewritten.
- **Live-state-derived (non-cacheable, "always-regenerate")**: values that depend on machine-local runtime state rather than config alone — chiefly **assigned host ports** and the **discovery environment variables** that embed them. These cannot be validated by a config-based staleness check, because the config that produced them is unchanged while the underlying state (port availability, prior assignments in global state) may have moved. Live-state-derived content must be refreshed — re-resolved against current global state — on every generation, even when the config-based staleness check reports "up-to-date." See [Port Types](./port-types.md) (assigned values are live-state-derived) and [State Management](./state-management.md#global-state).

#### Atomicity and Completeness

Generation must be **atomic**: write the full set of generated artifacts into a temporary directory, then rename it into place as a single step, so a reader never observes a half-written `.generated/`. A **completeness marker** (written last) records that the generation finished successfully and captures the generator/schema version it was produced with. Output lacking a valid, current-version marker is treated as stale and regenerated. A failed resolve step (the `docker compose config` equivalent — see Generation Logic) must **fail generation** and leave the previous good output in place, rather than persist a truncated or partial artifact. Validated by the Xcind proof-of-concept.

### Generation Logic (`workspace generate`)

1. **Resolve flavor** for each application (CLI → state → default_flavor → "default")
2. **Get compose files** from resolved flavor's `compose_files` list
3. **Validate compose files exist** on disk; if any are missing, report error with available alternatives:
   ```
   Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
     Application: backend
     Available compose files: docker-compose.yaml, docker-compose.dev.yaml
   ```
4. **Validate service references** in `exported_services` point to actual Compose services:
   ```
   Error: Exported service "api" references non-existent Compose service: backend
     Application: my-app
     Available services in docker-compose.yaml: web, db, redis
   ```
5. **Infer port values** for any exported services with omitted `port:` field (see Port Configuration)
6. **Default service names** for any exported services with omitted `service:` field
7. **Validate port values**: every resolved or inferred port must be an integer in the range 1–65535; protocol suffixes (e.g. `443/tcp`) are parsed and handled explicitly rather than passed through as opaque strings. An invalid value fails generation with an error that names the offending exported service:
   ```
   Error: Exported service "web" has an invalid port value: "https"
     Application: frontend
     Expected an integer in 1–65535 (optionally with a protocol suffix)
   ```
8. **Collect all exported services** across all applications in workspace
9. **Allocate and validate assigned host ports**: for every `assigned`-type export, resolve the sticky host port from global state or allocate a new one (see [Port Assignment Rules](./state-management.md#port-assignment-rules)), and validate that each resolved host port is currently usable. This step **must precede** the override write and the manifest write (steps 10 and 12) so that both artifacts embed the same, freshly-allocated host ports and discovery environment variables. Because these are live-state-derived values, they are re-resolved on every generation (see Config-Derived vs. Live-State-Derived Content). Validated by the Xcind proof-of-concept.
10. **Generate override file** with networks, aliases, labels, and environment variables (using the host ports allocated in step 9)
11. **Update state file** with resolved flavors
12. **Update manifest** with computed values (reflecting the post-allocation host ports and discovery variables from step 9)

### Initialization and Config Writes (`workspace init`)

Commands that write the primary config files — `workspace init` and any command
that edits `workspace.yaml` or an application's `application.yaml` — perform
**targeted field updates**, not whole-file rewrites. Re-running such a command
adds or updates only the fields it owns and **preserves unrelated,
user-authored fields** (comments, custom keys, hand-tuned values) already
present in the file.

This extends the preservation guarantee that already applies to the
`overrides/` directory (see [Generated Override Files](./generated-override-files.md))
to the primary config files themselves: re-initialization is idempotent for the
fields Scind manages and non-destructive for everything else. Validated by the
Xcind proof-of-concept.

### Shutdown Sequence (`workspace down`)

1. For each application (or specified apps if `-a` flag used):
   - Execute `docker compose down`
2. If full workspace teardown (no `-a` flag): remove workspace network
3. If `--volumes` specified, remove associated volumes

**Network removal timing**: The workspace network (`{workspace}-internal`) is only removed during a full workspace teardown (i.e., `workspace down` without the `-a` flag). When stopping individual applications with `-a`, the network is preserved to allow other running applications to continue communicating.

### Destroy Sequence (`workspace destroy`)

Completely removes a workspace and optionally its application directories:

1. Run `workspace down --volumes` to stop all containers and remove volumes
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force` or `--keep-apps`)
4. Remove `workspace.yaml`
5. Release any assigned ports from global state
6. Remove workspace from registry (`~/.config/scind/workspaces.yaml`)

**Flags**:
- `--force`: Skip confirmation prompts and remove application directories
- `--keep-apps`: Preserve application directories without prompting

### Viewing Logs

Using `scind-compose` (recommended):
```bash
# All logs for an application (context-aware)
scind-compose logs -f

# Specific service
scind-compose logs -f web

# Different app from workspace root
scind-compose -a backend logs -f
```

Using raw Docker Compose:
```bash
# All logs for an application
docker compose -p dev-backend logs -f

# Specific service
docker compose -p dev-backend logs -f web

# All containers in a workspace (using labels)
docker logs $(docker ps -q --filter "label=scind.workspace.name=dev")
```

### Listing Workspace Status

Using `scind-compose`:
```bash
scind-compose ps
scind-compose -a backend ps
```

Using raw Docker Compose:
```bash
# All containers in a workspace
docker ps --filter "label=scind.workspace.name=dev"

# All containers for an application
docker ps --filter "label=scind.app.name=backend"
```
