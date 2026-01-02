# Workspace Lifecycle Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Workspace lifecycle commands manage the state of workspaces and applications: starting containers, stopping them, and handling the transitions between states. This specification defines the behavior of `up`, `down`, `restart`, and `destroy` operations.

**Related Documents**:
- [ADR-0010: Up/Down Command Semantics](../decisions/0010-up-down-command-semantics.md)
- [Generated Override Files](./generated-override-files.md)
- [Context Detection](./context-detection.md)

---

## Behavior

### State Model

Workspaces and applications can be in one of these states:

| State | Description |
|-------|-------------|
| `stopped` | No containers running |
| `running` | All containers running |
| `partial` | Some containers running, some stopped |

State is determined by querying Docker, not stored in files.

---

## Data Schema

### Workspace State

Tracked via Docker labels on running containers:
- `contrail.workspace.name` identifies workspace membership
- `contrail.app.name` identifies application membership

### Runtime State File

**Location**: `{workspace}/.generated/state.yaml`

Stores explicit user choices (like active flavor), not container state:

```yaml
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite
```

---

## Operations

### up Command

**Command**: `contrail workspace up` / `contrail up`

Brings up a workspace: generates overrides if needed, creates networks, starts containers.

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Start specific app(s) only (repeatable) |
| `--no-generate` | Skip automatic regeneration |
| `--force-generate` | Force regeneration even if up-to-date |
| `-d, --detach` | Run in background (default: true) |

**Sequence**:
1. Ensure proxy is running
2. Create workspace network if it doesn't exist
3. Check if override files are stale; regenerate if needed (see Staleness Detection)
4. For each application (or specified apps if `-a` flag used):
   - Resolve active flavor
   - Execute `docker compose up -d --remove-orphans` with compose files + override

**Notes**:
- The `--remove-orphans` flag is **always** passed to `docker compose up`. This ensures containers from services removed by flavor changes, manual compose file edits, or renamed services are automatically stopped and removed.
- The workspace network is **always** created if it doesn't exist, even when starting a single application with `-a`. This ensures cross-application communication is available when other apps start later.

**Example**:
```bash
contrail workspace up --workspace=dev
contrail workspace up -a app-one -a app-two  # With context
contrail up  # Alias, with context
```

**Output**:
```
Starting app-one... done
Starting app-two... done
Starting app-three... done

Workspace 'dev' is running.
```

### down Command

**Command**: `contrail workspace down` / `contrail down`

Tears down a workspace: stops containers, optionally removes volumes.

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Stop specific app(s) only (repeatable) |
| `--volumes` | Also remove volumes |
| `--force` | Skip confirmation for destructive actions |

**Sequence**:
1. For each application (or specified apps if `-a` flag used):
   - Execute `docker compose down`
2. If full workspace teardown (no `-a` flag): remove workspace network
3. If `--volumes` specified, remove associated volumes

**Network removal timing**: The workspace network (`{workspace}-internal`) is only removed during a full workspace teardown (i.e., `workspace down` without the `-a` flag). When stopping individual applications with `-a`, the network is preserved to allow other running applications to continue communicating.

**Example**:
```bash
contrail workspace down
contrail down -a app-one          # Stop just app-one
contrail down --volumes --force   # Full teardown with volumes
```

### restart Command

**Command**: `contrail workspace restart`

Restarts a workspace or specific applications.

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Restart specific app(s) only (repeatable) |

**Behavior**: Equivalent to `down` followed by `up`. Volumes are **always preserved** (the internal `down` does not use `--volumes`).

### destroy Command

**Command**: `contrail workspace destroy`

Completely removes a workspace.

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `--force` | Skip confirmation prompts and remove application directories |
| `--keep-apps` | Preserve application directories without prompting |

**Sequence**:
1. Run `workspace down --volumes` to stop all containers and remove volumes
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force` or `--keep-apps`)
4. Remove `workspace.yaml`
5. Release any assigned ports from global state
6. Remove workspace from registry (`~/.config/contrail/workspaces.yaml`)

**Warning**: This is destructive. Without `--force` or `--keep-apps`, prompts for confirmation showing what will be removed.

---

## Examples

### Example 1: Complete Startup Sequence

```bash
$ contrail workspace up --workspace=dev
Proxy not running, starting...
Proxy: started
Creating network dev-internal...
Regenerating app-one.override.yaml (stale)...
Starting app-one... done
Starting app-two... done
Starting app-three... done

Workspace 'dev' is running.
URLs:
  - https://dev-app-one-web.contrail.test
  - https://dev-app-two-api.contrail.test
```

### Example 2: Partial Operations

```bash
# Stop just one app (network preserved)
$ contrail down -a app-one
Stopping app-one... done

# Other apps still running
$ contrail workspace status
Workspace: dev
Network: dev-internal (created)

APPLICATION  FLAVOR   STATUS   SERVICES
app-one      default  stopped  0/3 running
app-two      full     running  5/5 running
app-three    lite     running  2/2 running

# Start it back up
$ contrail up -a app-one
Starting app-one... done
```

### Example 3: Restart After Configuration Change

```bash
$ contrail flavor set full -a app-two
Warning: Application "app-two" is currently running.
The new flavor has been applied to the configuration, but running
containers still use the previous flavor.

To apply the flavor change:
  contrail app restart -a app-two

$ contrail app restart -a app-two
Stopping app-two... done
Starting app-two... done
```

### Example 4: Full Workspace Destruction

```bash
$ contrail workspace destroy
This will:
  - Stop all containers in workspace 'dev'
  - Remove all volumes
  - Delete .generated/ directory
  - Remove workspace from registry

Application directories:
  ./app-one/
  ./app-two/
  ./app-three/

Delete application directories? [y/N/keep]
> keep

Stopping containers...
Removing volumes...
Removing .generated/...
Releasing assigned ports...
Removing from registry...

Workspace 'dev' destroyed.
Application directories preserved.
```

---

## Staleness Detection

Override files are considered stale if any source file has a newer modification time:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`, `docker-compose.worker.yaml`)

**Behavior**:
- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Use `--no-generate` to skip regeneration entirely

**Note**: mtime comparison is simple and fast but may trigger unnecessary regeneration if files are touched without content changes. The `--force` flag provides explicit control when needed.

---

## Edge Cases

### Orphaned Containers

**Scenario**: Flavor change removes a service that was previously running.

**Behavior**: `--remove-orphans` is always passed to `docker compose up`, automatically stopping and removing containers for services no longer defined.

**Example**:
```
Flavor "full" → "lite" removes worker service

$ contrail up
Warning: Removing orphaned containers: dev-app-one-worker-1
Starting app-one... done
```

### Port Conflict at Startup

**Scenario**: Previously assigned port is now unavailable.

**Behavior**:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

### Docker Compose Error

**Scenario**: `docker compose up` fails for an application.

**Behavior**: Error is reported, other applications continue:
```
Starting app-one... done
Starting app-two... ERROR
  Error: Service 'web' failed to build: Dockerfile not found

Starting app-three... done

Workspace 'dev' partially running (1 error).
```

Exit code is non-zero if any application failed.

### Proxy Not Running

**Scenario**: User runs `workspace up` but proxy is down.

**Behavior**: Contrail automatically starts the proxy:
```
$ contrail workspace up
Proxy not running, starting...
Proxy: started

Starting app-one... done
```

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Docker not running | DOCKER | `Docker is not installed or not running` | Start Docker |
| Compose file not found | FILE_NOT_FOUND | `Compose file not found: {path}` | Fix configuration |
| Build failure | COMPOSE | `Service '{name}' failed to build` | Fix Dockerfile |
| Port conflict | PORT_CONFLICT | `Port {port} is assigned but no longer available` | Release and regenerate |
| Network creation failed | NETWORK | `Cannot create network: {error}` | Check Docker daemon |

---

## Viewing Logs

Using `contrail-compose` (recommended):
```bash
# All logs for an application (context-aware)
contrail-compose logs -f

# Specific service
contrail-compose logs -f web

# Different app from workspace root
contrail-compose -a app-two logs -f
```

Using raw Docker Compose:
```bash
# All logs for an application
docker compose -p dev-app-two logs -f

# All containers in a workspace (using labels)
docker logs $(docker ps -q --filter "label=contrail.workspace.name=dev")
```

---

## Listing Workspace Status

Using `contrail`:
```bash
$ contrail workspace status
Workspace: dev
Network: dev-internal (created)
Proxy: running

APPLICATION  FLAVOR   STATUS   SERVICES         URL
app-one      default  running  3/3 running      https://dev-app-one-web.contrail.test
app-two      full     running  5/5 running      https://dev-app-two-web.contrail.test
app-three    lite     stopped  0/2 running      —
```

Using Docker filters:
```bash
# All containers in a workspace
docker ps --filter "label=contrail.workspace.name=dev"

# All containers for an application
docker ps --filter "label=contrail.app.name=app-two"
```

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
