# Contrail CLI Reference

**Version**: 0.2.0
**Date**: December 2024

<!-- Migrated from specs/contrail-cli-reference.md -->

---

## Overview

Contrail uses a **resource-first** command structure (noun-verb) for most operations, similar to `docker`, `kubectl`, and `gh`.

---

## Global Flags

These flags can be used with any command:

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Specify workspace (overrides context detection) |
| `--app` | `-a` | Specify application(s) (repeatable, overrides context detection) |
| `--quiet` | `-q` | Minimal output, suppress context indicators and progress |
| `--verbose` | `-v` | Detailed output |
| `--json` | | Output in JSON format |
| `--yaml` | | Output in YAML format |
| `--color` | | Color output: `auto`, `always`, or `never` |
| `--config` | | Path to config file |
| `--help` | `-h` | Show help |

---

## Workspace Commands

### `contrail workspace list`

List all workspaces.

| Flag | Description |
|------|-------------|
| `--validate` | Check that registered paths still contain `workspace.yaml` |
| `--rebuild` | Rebuild registry from Docker labels |

### `contrail workspace show`

Show details for the current/specified workspace.

### `contrail workspace init`

Initialize a new workspace in the current or specified directory.

| Flag | Description |
|------|-------------|
| `--path` | Directory to create workspace in |

### `contrail workspace clone`

Clone application repositories defined in `workspace.yaml`.

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing application directories |

### `contrail workspace generate`

Generate override files for the workspace.

| Flag | Description |
|------|-------------|
| `--force` | Regenerate even if up-to-date |

### `contrail workspace prune`

Remove stale workspace registry entries.

| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be removed without making changes |

### `contrail workspace up`

Bring up a workspace (generate overrides and start containers).

| Flag | Description |
|------|-------------|
| `--no-generate` | Skip automatic regeneration |
| `--force-generate` | Force regeneration even if up-to-date |
| `-d, --detach` | Run in background (default: true) |

### `contrail workspace down`

Tear down a workspace.

| Flag | Description |
|------|-------------|
| `--volumes` | Also remove volumes |
| `--force` | Skip confirmation |

### `contrail workspace destroy`

Completely remove a workspace and optionally its application directories.

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompts and remove application directories |
| `--keep-apps` | Preserve application directories without prompting |

### `contrail workspace restart`

Restart a workspace (down + up).

### `contrail workspace status`

Show workspace status (container states, URLs).

---

## Application Commands

### `contrail app list`

List applications in the current/specified workspace.

### `contrail app show`

Show details for the current/specified application.

### `contrail app init`

Initialize an `application.yaml` in the current directory.

| Flag | Description |
|------|-------------|
| `-a, --app` | Application name (default: directory name) |

### `contrail app add`

Add an application to a workspace.

| Flag | Description |
|------|-------------|
| `-a, --app` | Application name (required) |
| `--repo` | Git repository URL to clone |
| `--path` | Custom path relative to workspace |

### `contrail app remove`

Remove an application from a workspace.

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation, also remove directory |

### `contrail app up`

Bring up a specific application.

### `contrail app down`

Tear down a specific application.

| Flag | Description |
|------|-------------|
| `--volumes` | Also remove volumes |

### `contrail app restart`

Restart a specific application.

### `contrail app status`

Show status for a specific application.

---

## Flavor Commands

### `contrail flavor list`

List available flavors for the current/specified application.

### `contrail flavor show`

Show the currently active flavor.

### `contrail flavor set <flavor>`

Set the active flavor for an application.

---

## Port Commands

Global commands (ignore directory context).

### `contrail port list`

List all port assignments across workspaces.

### `contrail port show <port>`

Show details for a specific port assignment.

### `contrail port release <port>`

Release a port assignment.

| Flag | Description |
|------|-------------|
| `--force` | Release even if container is running |

### `contrail port assign`

Manually assign a port.

| Flag | Description |
|------|-------------|
| `--port` | Port number (required) |
| `-w, --workspace` | Workspace name (required) |
| `-a, --app` | Application name (required) |
| `--service` | Service name (required) |

### `contrail port gc`

Garbage collect stale port assignments.

| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be released without making changes |

### `contrail port scan`

Scan for port conflicts between assignments and system.

| Flag | Description |
|------|-------------|
| `--fix` | Attempt to resolve conflicts automatically |

---

## Proxy Commands

Global commands (ignore directory context).

### `contrail proxy init`

Bootstrap proxy configuration.

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing configuration |
| `--domain` | Proxy domain (default: `contrail.test`) |
| `--path` | Directory to create proxy in |

### `contrail proxy up`

Start the Traefik proxy.

| Flag | Description |
|------|-------------|
| `--recreate` | Recreate the proxy network even if it exists |

### `contrail proxy down`

Stop the Traefik proxy.

### `contrail proxy restart`

Restart the Traefik proxy.

### `contrail proxy status`

Show proxy status.

---

## Config Commands

Global commands (ignore directory context).

### `contrail config show`

Show current configuration.

| Flag | Description |
|------|-------------|
| `--resolved` | Show fully resolved configuration with all defaults |

### `contrail config get <key>`

Get a configuration value.

### `contrail config set <key> <value>`

Set a configuration value.

### `contrail config path`

Show configuration file paths.

### `contrail config edit`

Open configuration in editor.

| Flag | Description |
|------|-------------|
| `--file` | Which config to edit: `proxy`, `workspace`, or `application` |

---

## Utility Commands

### `contrail validate`

Validate configuration files.

| Flag | Description |
|------|-------------|
| `--strict` | Treat warnings as errors |

### `contrail doctor`

Check system health and dependencies.

| Flag | Description |
|------|-------------|
| `--fix` | Attempt to fix issues automatically |

### `contrail open [service]`

Open service URL in browser.

| Flag | Description |
|------|-------------|
| `--print` | Print URL instead of opening browser |

### `contrail urls`

List all service URLs for current context.

### `contrail completion <shell>`

Generate shell completion scripts.

Supported shells: `bash`, `zsh`, `fish`, `powershell`

### `contrail init-shell <shell>`

Output shell integration script for `contrail-compose`.

Supported shells: `bash`, `zsh`, `fish`

---

## Top-Level Aliases

Convenience aliases for common operations:

| Alias | Equivalent |
|-------|------------|
| `contrail up` | `contrail workspace up` |
| `contrail down` | `contrail workspace down` |
| `contrail ps` | `contrail workspace status` |
| `contrail generate` | `contrail workspace generate` |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Docker/Compose error |
| 4 | Configuration error |
| 5 | Context detection failed |

---

## Output Behavior

### `--quiet` Mode

When `--quiet` or `-q` is specified:
- Context indicators are suppressed (no "Using workspace: ..." messages)
- Progress output is minimized
- Only essential output is shown (e.g., URLs, error messages)

### Context Feedback

When context is detected from the current directory, commands indicate what was found:

```bash
$ cd ~/workspaces/dev/app-one
$ contrail app status
# Using workspace: dev (from ../workspace.yaml)
# Using app: app-one (from ./application.yaml)

Status: running
Services: 3 running, 0 stopped
```

### Docker Not Available

If Docker is not running or not installed:

```
Error: Docker is not available.

Contrail requires Docker Engine and Docker Compose.
Run 'contrail doctor' for more details.
```

---

## Related Documents

- [Context Detection](../specs/context-detection.md)
- [Shell Integration](../specs/shell-integration.md)
- [Configuration Schemas](../specs/configuration-schemas.md)
