# Workspace Lifecycle Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

This specification defines the operations for workspace startup, shutdown, generation, and destruction.

---

## Startup Sequence (`workspace up`)

1. Ensure proxy is running
2. Create workspace network if it doesn't exist
3. Check if override files are stale; regenerate if needed
4. For each application (or specified apps if `-a` flag used):
   - Resolve active flavor
   - Execute `docker compose up -d --remove-orphans` with compose files + override

**Notes**:
- `--remove-orphans` always passed to handle removed services
- Workspace network created even when starting single app with `-a`

---

## Staleness Detection

Override files are considered stale if any source file has newer mtime:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (flavor may have changed)
- Active flavor's compose files

**Commands**:
- `workspace up` and `workspace generate` auto-regenerate stale overrides
- Use `--force` to regenerate regardless of staleness

---

## Generation Logic (`workspace generate`)

1. **Resolve flavor** for each application (CLI → state → default_flavor → "default")
2. **Get compose files** from resolved flavor's `compose_files` list
3. **Validate compose files exist**:
   ```
   Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
     Application: app-two
     Available compose files: docker-compose.yaml, docker-compose.dev.yaml
   ```
4. **Validate service references** in `exported_services`:
   ```
   Error: Exported service "api" references non-existent Compose service: backend
     Application: my-app
     Available services: web, db, redis
   ```
5. **Infer port values** for omitted `port:` fields
6. **Default service names** for omitted `service:` fields
7. **Collect all exported services** across workspace
8. **Generate override file** with networks, aliases, labels, environment
9. **Update state file** with resolved flavors
10. **Update manifest** with computed values

---

## Shutdown Sequence (`workspace down`)

1. For each application (or specified apps if `-a` flag used):
   - Execute `docker compose down`
2. If full workspace teardown (no `-a` flag): remove workspace network
3. If `--volumes` specified, remove associated volumes

**Network removal**: Only on full teardown (preserves cross-app communication when stopping individual apps)

---

## Destroy Sequence (`workspace destroy`)

1. Run `workspace down --volumes`
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force` or `--keep-apps`)
4. Remove `workspace.yaml`
5. Release any assigned ports from global state
6. Remove workspace from registry

**Flags**:
- `--force`: Skip prompts, remove application directories
- `--keep-apps`: Preserve application directories without prompting

---

## Restart Behavior

`workspace restart` equivalent to `down` followed by `up`:
- Volumes always preserved (internal `down` does not use `--volumes`)

---

## Flavor Changes

When `contrail flavor set FLAVOR` is executed:

1. Updates `.generated/state.yaml` with new flavor
2. Immediately regenerates affected application's override file
3. If application is running, displays warning:
   ```
   Warning: Application "app-name" is currently running.
   The new flavor has been applied to the configuration, but running
   containers still use the previous flavor.

   To apply the flavor change:
     contrail app restart -a app-name
   ```

### Flavor Change Effects

| Scenario | Effect | Resolution |
|----------|--------|------------|
| Adds services | New services defined but not running | `contrail up` |
| Removes services | Services still running but not in override | `contrail up` (with --remove-orphans) |
| Changes environment | Running containers have old values | `contrail app restart` |

---

## Related Documentation

- [ADR-0010: up/down Command Semantics](../../decisions/0010-up-down-command-semantics/README.md)
- [Configuration Schemas Spec](../configuration-schemas/README.md)
- [Context Detection Spec](../context-detection/README.md)
