<!-- Migrated from specs/scind-technical-spec.md:686-706 -->
<!-- Extraction ID: spec-state-management -->

# State Management Specification

**Version**: 1.0.0
**Date**: 2026-01-06
**Status**: Accepted

---

## Overview

Scind separates structure (configuration) from state (runtime). State represents explicit choices made by the user and system-managed assignments, not computed values. This specification covers both workspace-level state and global state.

---

## Workspace State

**Location**: `{workspace}/.generated/state.yaml` (gitignored)

Runtime state is tracked separately from configuration. State represents explicit choices made by the user (e.g., which flavor to use), not computed values.

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  app-one:
    flavor: default
  app-two:
    flavor: lite                      # Overridden from default
```

### Flavor Resolution Order

1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

---

## Global State

**Location**: `~/.config/scind/state.yaml` (global/per-user)

This file tracks port assignments for `assigned` type ports across all workspaces, plus an inventory of port availability for garbage collection and debugging.

```yaml
# AUTO-GENERATED - Managed by Scind
# Records assigned ports and port availability inventory

assigned_ports:
  dev:                                  # Workspace name
    app-one:                            # Application name
      db: 5432                          # Exported service: assigned host port
    app-two:
      db: 5433                          # Incremented because 5432 was taken
      cache: 6379
  review:
    app-one:
      db: 5434                          # Different workspace, different port

port_inventory:
  5432:
    status: assigned                    # assigned | unavailable | released
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:                         # Present only if status=assigned
      workspace: dev
      application: app-one
      exported_service: db
  5433:
    status: assigned
    first_seen: 2025-12-28T17:53:58Z
    last_checked: 2025-12-29T13:01:37Z
    assignment:
      workspace: dev
      application: app-two
      exported_service: db
  5434:
    status: unavailable                 # External process using this port
    first_seen: 2025-12-28T17:54:00Z
    last_checked: 2025-12-29T13:01:40Z
    # No assignment - taken by external process
```

---

## Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

---

## Port Conflict at Startup

If a previously assigned port has become unavailable (e.g., taken by an external process) when `workspace up` runs, Scind fails with a clear error:

```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  scind port scan       # Check which ports are conflicting
  scind port release 5432   # Release the conflicting assignment
  scind generate --force    # Regenerate with new port assignment
```

---

## Port Status Values

| Status | Description |
|--------|-------------|
| `assigned` | Port is assigned to a Scind workspace/application |
| `unavailable` | Port is in use by an external process (not managed by Scind) |
| `released` | Port was previously tracked but has been freed |

---

## Port Status Transitions

- `unavailable` -> `assigned`: Port became free, Scind claimed it
- `assigned` -> `released`: Workspace/app removed, port freed
- `unavailable` -> `released`: External process stopped, `scind port gc` cleaned it up

---

## Port Availability Checking

`scind port scan` and `scind port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

---

## Port Inventory Fields

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | One of: `assigned`, `unavailable`, `released` |
| `first_seen` | timestamp | When this port was first tracked by Scind |
| `last_checked` | timestamp | When port availability was last verified |
| `assignment` | object | Present only when `status=assigned` |
| `assignment.workspace` | string | Workspace name that owns this port |
| `assignment.application` | string | Application name within the workspace |
| `assignment.exported_service` | string | Exported service name from application.yaml |

---

## Related Documents

- [Port Types Specification](./port-types.md)
- [Configuration Schemas](./configuration-schemas.md)
- [Workspace Lifecycle](./workspace-lifecycle.md)
