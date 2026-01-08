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

### Staleness Detection

Scind uses **mtime comparison** to determine if generated override files need to be regenerated. Override files are considered stale if any of the following source files have a newer modification time than the generated override:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`, `docker-compose.worker.yaml`)

**Behavior**:
- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Touch a file accidentally? Use `--force` to ensure clean state

**Note**: mtime comparison is simple and fast but may trigger unnecessary regeneration if files are touched without content changes. The `--force` flag provides explicit control when needed.

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
7. **Collect all exported services** across all applications in workspace
8. **Generate override file** with networks, aliases, labels, and environment variables
9. **Update state file** with resolved flavors
10. **Update manifest** with computed values

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
