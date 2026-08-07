# State Management Specification

**Version**: 1.0.0
**Date**: 2026-01-06
**Status**: Accepted

---

## Overview

Scind separates structure (configuration) from state (runtime). State represents explicit choices made by the user and system-managed assignments, not computed values. This specification covers both workspace-level state and global state.

**State is load-bearing.** The Xcind proof-of-concept attempted a leaner,
stateless design and could not avoid either a **workspace registry** or a
**sticky assigned-port inventory** — both proved necessary in practice. This
confirms the structure-vs-state separation of [ADR-0005](../decisions/0005-structure-vs-state-separation.md):
these are genuine state, not values that can be recomputed from configuration.

**Cacheable vs. non-cacheable content.** For the purposes of staleness and
regeneration, generated content divides into **config-derived** content
(hostnames, aliases, labels — a pure function of configuration and flavor
state, and therefore cacheable) and **live-state-derived** content (assigned
host ports and the discovery environment variables that embed them — a function
of the machine-local state described here, and therefore *not* cacheable and
re-resolved on every generation). See
[Workspace Lifecycle: Config-Derived vs. Live-State-Derived Content](./workspace-lifecycle.md#config-derived-vs-live-state-derived-content)
and [Port Types](./port-types.md).

**Machine-local state belongs under `$XDG_STATE_HOME`.** The ephemeral,
machine-local artifacts described here — the assigned-port inventory, the
workspace registry, and proxy runtime state — are *state*, not *configuration*,
and per the [XDG Base Directory](https://specifications.freedesktop.org/basedir-spec/latest/)
specification they belong under `$XDG_STATE_HOME` rather than being collapsed
under the config home. This mirrors, at the storage layer, the structure-vs-state
separation of [ADR-0005](../decisions/0005-structure-vs-state-separation.md).
(The concrete environment-variable defaults are defined in the environment
variables specification; the `~/.config/scind/...` paths shown below are
illustrative.)

---

## Workspace State

**Location**: `{workspace}/.generated/state.yaml` (gitignored)

Runtime state is tracked separately from configuration. State represents explicit choices made by the user (e.g., which flavor to use), not computed values.

```yaml
# AUTO-GENERATED - Managed by workspace tooling
applications:
  frontend:
    flavor: full
  backend:
    flavor: lite                      # Overridden from default
```

### Flavor Resolution Order

1. CLI flag (`--flavor=X`)
2. State file (`.generated/state.yaml`)
3. Application's `default_flavor`
4. `"default"`

The resolved name must be a **declared** flavor — the literal `default` is the
only name that may resolve without a declaration, and an unresolvable active
flavor fails generation rather than falling back. A flavor recorded here can
become undeclared later, when an `application.yaml` edit removes it. See
[Configuration Schemas — Resolution Validity](./configuration-schemas.md#resolution-validity).

---

## Global State

**Location**: `~/.config/scind/state.yaml` (global/per-user)

This file tracks port assignments for `assigned` type ports across all workspaces, plus an inventory of port availability for garbage collection and debugging.

```yaml
# AUTO-GENERATED - Managed by Scind
# Records assigned ports and port availability inventory

assigned_ports:
  dev:                                  # Workspace name
    frontend:                           # Application name
      web: 8080                         # Exported service: assigned host port
    backend:
      api: 3000
      worker: 9000
    shared-db:
      db: 5432
      cache: 6379
  review:
    frontend:
      web: 8081                         # Different workspace, different port

port_inventory:
  5432:
    status: assigned                    # assigned | unavailable | released
    first_seen: 2025-12-28T17:53:55Z
    last_checked: 2025-12-29T13:01:33Z
    assignment:                         # Present only if status=assigned
      workspace: dev
      application: shared-db
      exported_service: db
  6379:
    status: assigned
    first_seen: 2025-12-28T17:53:58Z
    last_checked: 2025-12-29T13:01:37Z
    assignment:
      workspace: dev
      application: shared-db
      exported_service: cache
  5434:
    status: unavailable                 # External process using this port
    first_seen: 2025-12-28T17:54:00Z
    last_checked: 2025-12-29T13:01:40Z
    # No assignment - taken by external process
```

---

## Workspace Registry

**Location**: `~/.config/scind/workspaces.yaml` (global/per-user; see
`$XDG_STATE_HOME` note in the Overview)

The registry tracks known workspaces so that global operations (listing,
port reconciliation) can find them without scanning the whole filesystem.
Registry maintenance is not limited to init-time registration and
whole-registry pruning; it also supports:

- **Registering a cloned-but-never-run workspace**: a workspace that has been
  cloned/initialized but never started — and therefore carries **no Docker
  labels or containers yet** — can be added to the registry. Registration does
  not depend on the presence of running resources.
- **Dropping a single entry by path**: an individual workspace that has been
  moved or deleted can be removed from the registry by its recorded path,
  independently of a prune-all operation.

Validated by the Xcind proof-of-concept. (The specific CLI verbs are defined in
the CLI specification.)

---

## Port Assignment Rules

1. Try the port specified in `application.yaml`
2. If unavailable, increment and try again
3. Record the assignment in `assigned_ports` and `port_inventory`
4. Subsequent runs use the recorded port (sticky assignment)
5. The `--force` flag on `workspace generate` regenerates override files but preserves existing port assignments

---

## Port Conflict at Startup

**Attribute before declaring a conflict.** The startup conflict check must first
attribute any bound port to its owner via the port inventory, and must
**exclude ports bound by the workspace's own already-running containers** before
declaring a conflict. `workspace up` is idempotent — a re-up of a
still-running workspace will legitimately find its own assigned ports bound. A
bare bind-probe (attempting `net.Listen` on the port) cannot distinguish the
workspace's own running container from a foreign process; treating any bound
port as a conflict would raise **false conflicts** on the common re-up case.
Only a port bound by a **foreign** process — one not attributable to this
workspace's own containers in the inventory — is a real conflict. (Reverse
learning from the Xcind proof-of-concept, which surfaced this false-positive.)

If a previously assigned port has become unavailable because it is now held by a
**foreign** process (not one of the workspace's own containers) when
`workspace up` runs, Scind fails with a clear error:

```
Error: Port conflict detected for frontend

Port 5432 is assigned to frontend/postgres but is no longer available.
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
- `assigned` -> `released`: Proxy destroyed — `scind proxy destroy` removes the generated proxy state that holds the assigned-port inventory, so every workspace loses its assignments at once and receives new host ports on the next `up`
- `unavailable` -> `released`: External process stopped, `scind port gc` cleaned it up

---

## Port Availability Checking

`scind port scan` and `scind port gc` check port availability by attempting to bind to each tracked port using `net.Listen("tcp", ":PORT")`. Ports that can be bound are marked as available; ports that fail with "address already in use" remain in their current state. This method is reliable across platforms and doesn't require parsing system-specific files like `/proc/net`.

---

## Port Garbage Collection

The port-inventory maintenance commands (`scind port gc`, and any command that
reconciles the assigned-port state) apply a **path-existence** GC rule: an
assigned-port entry whose recorded application path no longer exists on disk is
dropped and its port released. A deleted application directory would otherwise
keep holding its port reservation indefinitely, silently shrinking the
available range. Because the application model is derived from directory
contents rather than a registry (see [Directory Structure](./directory-structure.md)),
the presence of the recorded path is the authoritative signal that an
assignment is still live. Validated by the Xcind proof-of-concept.

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
- [Workspace States](./workspace-lifecycle.md#workspace-states)
- [Workspace Operations](./workspace-lifecycle.md#operations)
