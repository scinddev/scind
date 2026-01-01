# Specification: Workspace Lifecycle

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

This document describes the operations for managing workspace lifecycle: startup, shutdown, generation, and destruction.

---

## Startup Sequence (`workspace up`)

1. Ensure proxy is running
2. Create workspace network if it doesn't exist
3. Check if override files are stale; regenerate if needed
4. For each application (or specified apps if `-a` flag used):
   - Resolve active flavor
   - Execute `docker compose up -d --remove-orphans` with compose files + override

**Notes**:
- The `--remove-orphans` flag is always passed to `docker compose up`. This ensures that containers from services removed by flavor changes, manual compose file edits, or renamed services are automatically stopped and removed.
- The workspace network is always created if it doesn't exist, even when starting a single application with `-a`. This ensures cross-application communication is available when other apps are started later.

---

## Generation Logic (`workspace generate`)

1. **Resolve flavor** for each application (CLI → state → default_flavor → "default")
2. **Get compose files** from resolved flavor's `compose_files` list
3. **Validate compose files exist** on disk; if any are missing, report error:
   ```
   Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
     Application: app-two
     Available compose files: docker-compose.yaml, docker-compose.dev.yaml
   ```
4. **Validate service references** in `exported_services` point to actual Compose services:
   ```
   Error: Exported service "api" references non-existent Compose service: backend
     Application: my-app
     Available services in docker-compose.yaml: web, db, redis
   ```
5. **Infer port values** for any exported services with omitted `port:` field
6. **Default service names** for any exported services with omitted `service:` field
7. **Collect all exported services** across all applications in workspace
8. **Generate override file** with networks, aliases, labels, and environment variables
9. **Update state file** with resolved flavors
10. **Update manifest** with computed values

---

## Staleness Detection

Contrail uses **mtime comparison** to determine if generated override files need to be regenerated. Override files are considered stale if any of the following source files have a newer modification time:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`, `docker-compose.worker.yaml`)

**Behavior**:
- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness

---

## Flavor Changes and Running Applications

When `contrail flavor set FLAVOR` is executed:
1. Updates `.generated/state.yaml` with the new flavor
2. Immediately regenerates the affected application's override file
3. If the application is currently running, displays a warning with restart guidance

Flavor changes affect running applications in different ways:

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Flavor adds services | New services defined in override but not running | Run `contrail up` to start new services |
| Flavor removes services | Services still running but not in override | Run `contrail up` to stop orphaned services |
| Flavor changes environment | Running containers have old values | Run `contrail app restart` to pick up changes |

**Orphaned service handling**: When `contrail up` is run after a flavor change that removes services, Contrail passes `--remove-orphans` to Docker Compose to stop and remove containers for services no longer defined in the active configuration.

---

## Shutdown Sequence (`workspace down`)

1. For each application (or specified apps if `-a` flag used):
   - Execute `docker compose down`
2. If full workspace teardown (no `-a` flag): remove workspace network
3. If `--volumes` specified, remove associated volumes

**Network removal timing**: The workspace network (`{workspace}-internal`) is only removed during a full workspace teardown. When stopping individual applications with `-a`, the network is preserved.

---

## Destroy Sequence (`workspace destroy`)

Completely removes a workspace and optionally its application directories:

1. Run `workspace down --volumes` to stop all containers and remove volumes
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force` or `--keep-apps`)
4. Remove `workspace.yaml`
5. Release any assigned ports from global state
6. Remove workspace from registry (`~/.config/contrail/workspaces.yaml`)

**Flags**:
- `--force`: Skip confirmation prompts and remove application directories
- `--keep-apps`: Preserve application directories without prompting

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

# Specific service
docker compose -p dev-app-two logs -f web
```

---

## Listing Workspace Status

Using `contrail-compose`:
```bash
contrail-compose ps
contrail-compose -a app-two ps
```

Using raw Docker Compose:
```bash
# All containers in a workspace
docker ps --filter "label=contrail.workspace.name=dev"

# All containers for an application
docker ps --filter "label=contrail.app.name=app-two"
```

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md)
- [Generated Override Files](generated-override-files.md)
- [Architecture Overview](../architecture/overview.md)
