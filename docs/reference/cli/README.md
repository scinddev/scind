# Contrail CLI Reference

**Version**: 0.2.4
**Date**: December 2024

---

## Command Structure

```
contrail [resource] [action] [options...]
```

All targeting uses **options** rather than positional arguments.

---

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Specify workspace (overrides context) |
| `--app` | `-a` | Specify application (overrides context) |
| `--quiet` | `-q` | Minimal output, suppress progress |
| `--verbose` | `-v` | Detailed output |
| `--json` | | Output in JSON format |
| `--yaml` | | Output in YAML format |
| `--color` | | Color output: `auto`, `always`, `never` |
| `--help` | `-h` | Show help |
| `--version` | | Show version |

---

## Workspace Commands

### `contrail workspace list`
List all registered workspaces.

**Flags**: `--validate`, `--rebuild`

### `contrail workspace init`
Initialize a new workspace.

**Flags**: `-w, --workspace` (name), `--path`

### `contrail workspace up`
Bring up a workspace.

**Flags**: `-w`, `-a` (repeatable), `--no-generate`, `--force-generate`, `-d, --detach`

### `contrail workspace down`
Tear down a workspace.

**Flags**: `-w`, `-a` (repeatable), `--volumes`, `--force`

### `contrail workspace status`
Show workspace status.

**Flags**: `-w`

### `contrail workspace restart`
Restart a workspace.

**Flags**: `-w`, `-a` (repeatable)

### `contrail workspace generate`
Generate override files.

**Flags**: `-w`, `-a` (repeatable), `--force`

### `contrail workspace clone`
Clone application repositories.

**Flags**: `-w`, `-a` (repeatable)

### `contrail workspace prune`
Remove stale registry entries.

**Flags**: `--dry-run`

### `contrail workspace destroy`
Completely remove a workspace.

**Flags**: `-w`, `--force`, `--keep-apps`

### `contrail workspace show`
Show detailed workspace information.

**Flags**: `-w`

---

## Application Commands

### `contrail app list`
List applications in a workspace.

**Flags**: `-w`

### `contrail app show`
Show application details.

**Flags**: `-w`, `-a`

### `contrail app init`
Initialize application configuration in current directory.

**Flags**: `-a` (name, default: directory name)

### `contrail app add`
Add application to workspace.

**Flags**: `-w`, `-a` (required), `--repo`, `--path`

### `contrail app remove`
Remove application from workspace.

**Flags**: `-w`, `-a`, `--force`

### `contrail app up`
Bring up an application.

**Flags**: `-w`, `-a`

### `contrail app down`
Tear down an application.

**Flags**: `-w`, `-a`, `--volumes`

### `contrail app restart`
Restart an application.

**Flags**: `-w`, `-a`

### `contrail app status`
Show application status.

**Flags**: `-w`, `-a`

---

## Flavor Commands

### `contrail flavor list`
List available flavors.

**Flags**: `-w`, `-a`

### `contrail flavor show`
Show current active flavor.

**Flags**: `-w`, `-a`

### `contrail flavor set <flavor>`
Set active flavor.

**Flags**: `-w`, `-a`

---

## Port Commands

### `contrail port list`
List all port assignments.

**Flags**: `-w` (filter), `-v, --verbose`

### `contrail port show <port>`
Show port details.

### `contrail port release <port>`
Release a port assignment.

**Flags**: `--force`

### `contrail port assign <port> <workspace/app/service>`
Manually assign a port.

### `contrail port gc`
Garbage collect stale ports.

**Flags**: `--dry-run`

### `contrail port scan`
Scan and update port inventory.

---

## Proxy Commands

### `contrail proxy init`
Bootstrap proxy configuration.

**Flags**: `--force`, `--domain`, `--path`

### `contrail proxy up`
Start the Traefik proxy.

**Flags**: `--recreate`

### `contrail proxy down`
Stop the proxy.

### `contrail proxy restart`
Restart the proxy.

### `contrail proxy status`
Show proxy status.

---

## Config Commands

### `contrail config show`
Show all configuration.

### `contrail config get <key>`
Get a configuration value.

### `contrail config set <key> <value>`
Set a configuration value.

### `contrail config path`
Show configuration file locations.

### `contrail config edit`
Open configuration in editor.

---

## Utility Commands

### `contrail validate`
Validate configuration files.

**Flags**: `-w`, `-a`

### `contrail doctor`
Check system health and dependencies.

### `contrail open`
Open service URL in browser.

**Flags**: `-w`, `-a`, `--service`

### `contrail urls`
List all service URLs.

**Flags**: `-w`

---

## Shell Integration Commands

### `contrail init-shell <shell>`
Output shell integration script.

**Arguments**: `bash`, `zsh`, or `fish`

### `contrail compose-prefix`
Output docker compose command prefix.

**Flags**: `-w`, `-a`

### `contrail completion <shell>`
Generate shell completion script.

---

## Top-Level Aliases

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
| 3 | Configuration error |
| 4 | Docker/Compose error |
| 5 | Context detection failed |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CONTRAIL_CONFIG` | Path to global config |
| `CONTRAIL_WORKSPACE` | Default workspace |
| `NO_COLOR` | Disable colored output |
| `CONTRAIL_DEBUG` | Enable debug logging |

---

## Related Documentation

- [Context Detection Spec](../../specs/context-detection/README.md)
- [Shell Integration Spec](../../specs/shell-integration/README.md)
- [Implementation: Go Stack](../../implementation/go-stack/README.md)
