<!-- Migrated from specs/contrail-cli-reference.md:1-1599 -->
<!-- Extraction ID: reference-cli -->

# Contrail CLI Reference

**Version**: 0.2.3-draft
**Date**: December 2024
**Status**: Design Phase

This document is the authoritative reference for Contrail's command-line interface. It defines command structure, arguments, flags, and behaviors.

---

## Command Structure

Contrail follows a **resource-first** pattern inspired by Docker and Mutagen:

```
contrail [resource] [action] [options...]
```

All targeting uses **options** rather than positional arguments:

```bash
# Good: Options-based
contrail app status --workspace=dev --app=app-one

# With context detection (from current directory)
contrail app status
```

---

## Context Detection

Contrail automatically detects workspace and application context from the current directory, reducing the need for explicit `--workspace` and `--app` flags.

### Detection Rules

Context detection uses a **workspace boundary** approach to prevent accidental detection of config files in vendor packages or nested test fixtures.

1. **Workspace context** (found first): Walk up the directory tree looking for `workspace.yaml`
   - If found, this establishes the **workspace root**
   - The `workspace.name` value becomes the implicit `--workspace` value

2. **Application context** (bounded by workspace): Walk up from current directory toward the workspace root looking for `application.yaml`
   - Only considers `application.yaml` files **within the workspace directory tree**
   - If found, the directory name containing it becomes the implicit `--app` value
   - **Never traverses above the workspace root**—this prevents vendor packages from hijacking context

3. **Both can be detected simultaneously**:
   ```
   ~/workspaces/dev/app-one/src/components/
                   │        │
                   │        └── application.yaml → app = "app-one"
                   │
                   └── workspace.yaml → workspace = "dev"
   ```

4. **Explicit flags override detection**: `--workspace` and `--app` always take precedence over context detection
   - When any `--app` flag is specified, context-detected application is **completely ignored**
   - This applies even when multiple `-a` flags are used
   - To start multiple specific apps: `contrail up -a app-one -a app-two`

5. **Global commands ignore context**: `port`, `proxy`, and `config` commands don't use directory context

### Flag Override Behavior

When explicit flags are provided, they **replace** (not add to) context detection:

```bash
# From within app-one directory (context would detect app-one)
$ cd ~/workspaces/dev/app-one

# This starts ONLY app-two, not both app-one and app-two
$ contrail up -a app-two
# Starting: app-two
# (app-one from context is ignored)

# To start multiple apps, list them all explicitly
$ contrail up -a app-one -a app-two
# Starting: app-one, app-two
```

This "explicit replaces context" behavior ensures predictable results—when you specify apps explicitly, you get exactly what you asked for.

### Edge Cases

**Nested vendor packages**: If working in `app-one/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored. The workspace's `app-one/application.yaml` is found first when walking toward the workspace root.

**Workspace within workspace**: If a test fixture has its own `workspace.yaml` nested inside a workspace (e.g., for integration tests), the closest `workspace.yaml` wins—this is the test fixture, which is the expected behavior.

### Context Feedback

When context is detected, commands indicate what was found:

```bash
$ cd ~/workspaces/dev/app-one
$ contrail app status
# Using workspace: dev (from ../workspace.yaml)
# Using app: app-one (from ./application.yaml)

Status: running
Services: 3 running, 0 stopped
...
```

Use `--quiet` or `-q` to suppress context indicators.

### Error Handling

When context is required but not detected, error messages provide debugging hints:

**No workspace found, but application.yaml exists** (helps identify misplaced apps):
```bash
$ cd ~/random-project
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

**Neither workspace nor application found**:
```bash
$ cd ~
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: contrail app status --workspace=NAME --app=NAME

Available workspaces: contrail workspace list
```

**Workspace found but no application context** (for app-specific commands):
```bash
$ cd ~/workspaces/dev
$ contrail app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

---

## Resources

Contrail manages these resource types:

| Resource | Description | Context-Aware |
|----------|-------------|---------------|
| `workspace` | Isolated environment containing applications | Yes |
| `app` | A Docker Compose application within a workspace | Yes |
| `flavor` | Named configuration for running an application | Yes |
| `port` | Host port assignments for assigned-type services | No |
| `proxy` | The Traefik reverse proxy layer | No |
| `config` | Global and user configuration | No |

---

## Global Flags

These flags are available on all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Specify workspace (overrides context detection) |
| `--app` | `-a` | Specify application (overrides context detection) |
| `--quiet` | `-q` | Minimal output, suppress context indicators and progress |
| `--verbose` | `-v` | Detailed output |
| `--json` | | Output in JSON format |
| `--yaml` | | Output in YAML format |
| `--color` | | Color output: `auto` (default), `always`, or `never` |
| `--help` | `-h` | Show help for command |
| `--version` | | Show Contrail version |

### Output Behavior

**Progress output**: Multi-application operations show per-application progress by default:
```
Starting app-one... done
Starting app-two... done
Starting app-three... done
```

**`--quiet` behavior**:
- Suppresses context indicators ("Using workspace: dev")
- Suppresses progress messages ("Starting app-one... done")
- Status commands output just the value: `running`
- List commands output names only, one per line
- Errors are always shown regardless of `--quiet`

```bash
# Normal output
$ contrail workspace list
NAME    PATH                      APPS   STATUS
dev     ~/workspaces/dev          3      running
staging ~/workspaces/staging      2      stopped

# Quiet output (machine-readable)
$ contrail workspace list -q
dev
staging

# Status with quiet
$ contrail app status -q
running
```

---

## Workspace Commands

Manage workspace lifecycle and orchestration.

### `contrail workspace list`

List all registered workspaces.

```bash
contrail workspace list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--validate` | Check that registered paths still contain `workspace.yaml` |
| `--rebuild` | Rebuild registry from Docker labels (useful if registry is missing) |

**Output** (default table):
```
NAME     APPS  STATUS   PATH
dev      3     running  ~/workspaces/dev
review   3     stopped  ~/workspaces/review
control  2     running  ~/workspaces/control
```

**Discovery mechanism**: Reads from the workspace registry (`~/.config/contrail/workspaces.yaml`). If the registry is missing, automatically attempts to rebuild from Docker container labels.

---

### `contrail workspace prune`

Remove stale entries from the workspace registry.

```bash
contrail workspace prune [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be removed without making changes |

**Behavior**:
- Checks each registered workspace path for a valid `workspace.yaml`
- Removes entries where the path no longer exists or lacks a valid config
- Reports removed entries

**Example**:
```bash
contrail workspace prune
# Removed: old-project (path /home/user/old-project no longer exists)
# Registry: 3 workspaces remaining
```

---

### `contrail workspace show`

Show detailed information about a workspace, including the computed manifest.

```bash
contrail workspace show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**: Workspace configuration, applications, computed hostnames, port assignments.

---

### `contrail workspace init`

Initialize a new workspace.

```bash
contrail workspace init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Workspace name (required if not in existing workspace dir) |
| `--path` | Directory to create workspace in (default: `./{name}` or current dir) |

**Behavior**:
- If run in an empty directory without `--workspace`, prompts for name
- If run with `--workspace=NAME`, creates `./NAME/workspace.yaml` or `./workspace.yaml` if `--path=.`
- Creates initial directory structure
- **Registers the workspace** in `~/.config/contrail/workspaces.yaml`
- **Fails if name is already registered** to a different path (enforces workspace name uniqueness)

**Example** (new workspace):
```bash
contrail workspace init --workspace=dev
# Creates ./dev/workspace.yaml
# Registers "dev" -> ./dev in workspace registry
```

**Example** (current directory):
```bash
cd ~/my-project
contrail workspace init --workspace=dev
# Creates ./workspace.yaml with name: dev
# Registers "dev" -> ~/my-project in workspace registry
```

**Example** (name collision):
```bash
contrail workspace init --workspace=dev
# Error: Workspace "dev" already registered at ~/workspaces/dev
# Use a different name, or run `contrail workspace prune` if that path no longer exists
```

---

### `contrail workspace clone`

Clone all application repositories defined in the workspace.

```bash
contrail workspace clone [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Clone specific app(s) only (repeatable) |

**Behavior**:
- Reads `repository` URLs from `workspace.yaml`
- Clones each application to its configured path
- Skips applications that already exist locally
- **Skips applications with `path: .`** (single-app workspaces where the application is the workspace root)

**Error handling**:
- On clone failure (network error, auth failure, invalid URL), aborts immediately
- Partial clones are cleaned up on failure
- Exit code 1 on any failure

**Single-app workspace handling**:
```bash
contrail workspace clone
# Skipping myapp: application is workspace root (path: .)
# Cloned: app-two -> ./app-two
```

---

### `contrail workspace generate`

Generate or regenerate Docker Compose override files.

```bash
contrail workspace generate [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Generate for specific app(s) only (repeatable) |
| `--force` | Regenerate even if files appear up-to-date |

**Behavior**:
- Reads `workspace.yaml` and each `application.yaml`
- Generates `.generated/*.override.yaml` files
- Updates `.generated/manifest.yaml`
- Updates `.generated/state.yaml` if needed

**Note**: Generation happens automatically during `workspace up` when source files have changed. Use this command for explicit control.

---

### `contrail workspace up`

Bring up a workspace (generate overrides if needed, create networks, start containers).

```bash
contrail workspace up [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Start specific app(s) only (repeatable) |
| `--no-generate` | Skip automatic regeneration |
| `--force-generate` | Force regeneration even if up-to-date |
| `-d, --detach` | Run in background (default: true) |

**Behavior**:
1. Detect or require workspace context
2. Check if override files are stale; regenerate if needed
3. Ensure workspace network (`{workspace}-internal`) exists
4. Ensure `contrail-proxy` network exists and proxy is running
5. For each application (or specified apps):
   - Resolve active flavor
   - Execute `docker compose up -d` with appropriate files
6. Report status

**Example**:
```bash
contrail workspace up --workspace=dev
contrail workspace up -a app-one -a app-two  # With context
contrail up  # Alias, with context
```

---

### `contrail workspace down`

Tear down a workspace (stop containers, remove networks).

```bash
contrail workspace down [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Stop specific app(s) only (repeatable) |
| `--volumes` | Also remove volumes |
| `--force` | Skip confirmation for destructive actions |

**Behavior**:
1. For each application (or specified apps):
   - Execute `docker compose down`
2. If all apps stopped, optionally remove workspace network
3. If `--volumes` specified, remove associated volumes
4. Prompt before removing volumes that may contain data (unless `--force`)

---

### `contrail workspace restart`

Restart a workspace or specific applications.

```bash
contrail workspace restart [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Restart specific app(s) only (repeatable) |

**Behavior**: Equivalent to `down` followed by `up`. Volumes are always preserved (the internal `down` does not use `--volumes`).

---

### `contrail workspace status`

Show the running status of a workspace.

```bash
contrail workspace status [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**:
```
Workspace: dev
Network: dev-internal (created)
Proxy: running

APPLICATON  FLAVOR   STATUS   SERVICES         URL
app-one     default  running  3/3 running      https://dev-app-one-web.contrail.test
app-two     full     running  5/5 running      https://dev-app-two-web.contrail.test
app-three   lite     stopped  0/2 running      —
```

---

### `contrail workspace destroy`

Completely remove a workspace.

```bash
contrail workspace destroy [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `--force` | Skip confirmation prompts and remove application directories |
| `--keep-apps` | Preserve application directories without prompting |

**Behavior**:
1. Run `workspace down --volumes`
2. Remove `.generated/` directory
3. Prompt before removing application directories (unless `--force` or `--keep-apps`)
4. Remove `workspace.yaml`
5. Release any assigned ports
6. Remove workspace from registry (`~/.config/contrail/workspaces.yaml`)

**Warning**: This is destructive. Without `--force` or `--keep-apps`, prompts for confirmation showing what will be removed.

---

## Application Commands

Manage applications within workspaces.

### `contrail app list`

List applications in a workspace.

```bash
contrail app list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**:
```
NAME       FLAVOR   STATUS   SERVICES  PATH
app-one    default  running  3/3       ./app-one
app-two    full     running  5/5       ./app-two
app-three  lite     stopped  0/2       ./app-three
```

---

### `contrail app show`

Show detailed information about an application.

```bash
contrail app show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Output**: Application configuration, exported services, current flavor, computed hostnames and ports.

---

### `contrail app init`

Initialize an application configuration in the current directory.

```bash
contrail app init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-a, --app` | Application name (default: current directory name) |

**Behavior**:
- Creates `application.yaml` in current directory
- Scans for existing `docker-compose.yaml` to suggest exported services
- Sets up default flavor pointing to existing compose file(s)

**Use case**: Promoting an existing Docker Compose project to a Contrail application.

```bash
cd ~/my-project
contrail app init --app=myapp
# Creates ./application.yaml
```

---

### `contrail app add`

Add an application to a workspace.

```bash
contrail app add [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Application name (required) |
| `--repo` | Git repository URL to clone |
| `--path` | Custom path relative to workspace (default: `./{app}`) |

**Behavior**:
- Adds application entry to `workspace.yaml`
- If `--repo` provided, clones repository to path
- If `--path=.`, configures app to use workspace root (single-app workspace)

**Single-app workspace naming**: The workspace name and app name are independent. You can have workspace `dev` with app `myapi`, or create multiple workspaces (`dev`, `review`, `feature-x`) each containing the same app name (`myapi`). The app name comes from `--app`; the workspace name comes from `workspace init --workspace`.

**Example**:
```bash
contrail app add --app=api --repo=git@github.com:org/api.git
contrail app add --app=main --path=.  # Use current directory
```

---

### `contrail app remove`

Remove an application from a workspace.

```bash
contrail app remove [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |
| `--force` | Skip confirmation, also remove directory |

**Behavior**:
- Stops application if running
- Removes from `workspace.yaml`
- Prompts before removing application directory (unless `--force`)
- Releases any assigned ports

---

### `contrail app up`

Bring up a single application.

```bash
contrail app up [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

Equivalent to `contrail workspace up --app=NAME`.

---

### `contrail app down`

Tear down a single application.

```bash
contrail app down [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |
| `--volumes` | Also remove volumes |

Equivalent to `contrail workspace down --app=NAME`.

---

### `contrail app restart`

Restart a single application.

```bash
contrail app restart [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Behavior**: Equivalent to `app down` followed by `app up`. Volumes are always preserved.

---

### `contrail app status`

Show status of a single application.

```bash
contrail app status [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

---

## Flavor Commands

Manage application flavors (named configurations).

### `contrail flavor list`

List available flavors for an application.

```bash
contrail flavor list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Output**:
```
NAME     COMPOSE FILES                                    ACTIVE
lite     docker-compose.yaml
full     docker-compose.yaml, docker-compose.worker.yaml  ✓
debug    docker-compose.yaml, docker-compose.debug.yaml
```

---

### `contrail flavor show`

Show the current active flavor for an application.

```bash
contrail flavor show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

---

### `contrail flavor set`

Set the active flavor for an application.

```bash
contrail flavor set <flavor> [flags]
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `flavor` | Name of flavor to activate |

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Behavior**:
- Updates `.generated/state.yaml` with the new flavor
- Immediately regenerates the affected application's override file
- If the application is running, displays a warning with restart guidance
- Does NOT automatically restart containers

> **Warning when application is running**: If the target application is currently running, `flavor set` displays a warning explaining that running containers still use the previous flavor and provides the command to apply changes.

**After changing flavors**:
| Scenario | Command |
|----------|---------|
| Flavor adds/removes services | `contrail up` (starts new services, stops orphaned services) |
| Flavor changes environment or config | `contrail app restart -a APP` |

**Example**:
```bash
# Change flavor (regenerates config, warns if running)
contrail flavor set full --app=app-two

# Apply to running application
contrail app restart --app=app-two
```

---

## Port Commands

Manage host port assignments for `assigned`-type services.

### `contrail port list`

List all assigned ports.

```bash
contrail port list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Filter to specific workspace |
| `-v, --verbose` | Include bind status check |

**Output**:
```
PORT   WORKSPACE  APP      SERVICE  STATUS
5432   dev        app-one  db       assigned
5433   dev        app-two  db       assigned
5434   review     app-one  db       assigned
6379   dev        app-one  cache    assigned
```

With `--verbose`:
```
PORT   WORKSPACE  APP      SERVICE  STATUS    BOUND
5432   dev        app-one  db       assigned  yes
5433   dev        app-two  db       assigned  yes
5434   review     app-one  db       assigned  no
6379   dev        app-one  cache    assigned  yes
```

---

### `contrail port show`

Show details about a specific port assignment.

```bash
contrail port show <port>
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `port` | Port number to inspect |

---

### `contrail port release`

Manually release a port assignment.

```bash
contrail port release <port> [flags]
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `port` | Port number to release |

**Flags**:
| Flag | Description |
|------|-------------|
| `--force` | Release even if port appears in use |

---

### `contrail port assign`

Manually assign a port (advanced usage).

```bash
contrail port assign <port> <workspace>/<app>/<service>
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `port` | Port number to assign |
| `target` | Target in format `workspace/app/service` |

---

### `contrail port gc`

Garbage collect released and unbound ports.

```bash
contrail port gc [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be released without doing it |

**Behavior**:
- Scans port inventory for ports marked as `assigned` but not actually bound
- Releases ports whose workspaces/apps no longer exist
- Reports actions taken

---

### `contrail port scan`

Re-scan and update port availability inventory.

```bash
contrail port scan
```

**Behavior**:
- Checks all tracked ports for current bind status
- Updates `port_inventory` in global state
- Useful after external changes (other processes releasing ports)

---

## Proxy Commands

Manage the Traefik reverse proxy. The proxy is a shared instance that serves all workspaces on the host.

### `contrail proxy init`

Bootstrap the proxy configuration. This creates the Traefik Docker Compose project at `~/.config/contrail/proxy/`.

```bash
contrail proxy init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing configuration (useful for recovery) |
| `--domain` | Set proxy domain (default: `contrail.test`) |
| `--path` | Directory to create proxy in (default: `~/.config/contrail/proxy/`) |

**Behavior**:
1. Check if proxy configuration already exists
   - If exists and no `--force`: error with message
   - If exists and `--force`: backup existing config and overwrite
2. Create proxy directory structure (`docker-compose.yaml`, `traefik.yaml`, `dynamic/`, `certs/`)
3. Create `contrail-proxy` Docker network if it doesn't exist
4. Output next steps (DNS setup, starting proxy)

**Example**:
```bash
$ contrail proxy init
Created proxy configuration at ~/.config/contrail/proxy/

Next steps:
  1. Configure DNS for *.contrail.test → 127.0.0.1
     (See: contrail doctor for DNS verification)
  2. Start the proxy:
     contrail proxy up

$ contrail proxy init
Error: Proxy configuration already exists at ~/.config/contrail/proxy/
Use --force to overwrite, or --path to create elsewhere.

$ contrail proxy init --force
Backed up existing configuration to ~/.config/contrail/proxy.backup.20241230/
Created proxy configuration at ~/.config/contrail/proxy/
```

**Custom domain example**:
```bash
$ contrail proxy init --domain mydev.local
Created proxy configuration at ~/.config/contrail/proxy/
Domain set to: mydev.local
```

---

### `contrail proxy up`

Start the Traefik proxy.

```bash
contrail proxy up [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--recreate` | Recreate the proxy network even if it exists |

**Behavior**:
- Creates `contrail-proxy` network if it doesn't exist
- Validates existing network configuration matches expected settings
- Starts Traefik container from proxy configuration
- If proxy configuration doesn't exist, runs `proxy init` first

**Network conflict handling**:
If the `contrail-proxy` network exists but was created by a different tool or has incompatible settings, `proxy up` will warn:
```
Warning: Network 'contrail-proxy' exists but may not have been created by Contrail.
  Driver: bridge (expected: bridge) ✓
  Labels: contrail.managed not found ⚠

Use 'contrail proxy up --recreate' to recreate the network.
```

**Note**: Users rarely need to call this directly. `workspace up` automatically starts the proxy if it's not running.

---

### `contrail proxy down`

Stop the Traefik proxy.

```bash
contrail proxy down
```

---

### `contrail proxy restart`

Restart the Traefik proxy.

```bash
contrail proxy restart
```

---

### `contrail proxy status`

Show proxy status.

```bash
contrail proxy status
```

**Output** (dashboard enabled):
```
Proxy: running
Network: contrail-proxy (created)
Dashboard: http://localhost:8080
Entrypoints:
  - web: :80
  - websecure: :443
```

**Output** (dashboard disabled):
```
Proxy: running
Network: contrail-proxy (created)
Dashboard: disabled
Entrypoints:
  - web: :80
  - websecure: :443
```

**Note**: The dashboard URL reflects the configured port from `proxy.yaml` (`dashboard.port`, default 8080). If `dashboard.enabled` is false, the dashboard line shows "disabled".

---

## Config Commands

Manage Contrail configuration.

### `contrail config show`

Show all configuration values.

```bash
contrail config show
```

**Output**:
```yaml
proxy:
  domain: contrail.test
paths:
  global_config: ~/.config/contrail/proxy.yaml
  global_state: ~/.config/contrail/state.yaml
```

---

### `contrail config get`

Get a specific configuration value.

```bash
contrail config get <key>
```

**Example**:
```bash
contrail config get proxy.domain
# contrail.test
```

---

### `contrail config set`

Set a configuration value.

```bash
contrail config set <key> <value>
```

**Example**:
```bash
contrail config set proxy.domain local.test
```

---

### `contrail config path`

Show configuration file locations.

```bash
contrail config path
```

**Output**:
```
Global config: ~/.config/contrail/proxy.yaml
Global state:  ~/.config/contrail/state.yaml
```

---

### `contrail config edit`

Open the global configuration file in your default editor.

```bash
contrail config edit
```

**Behavior**:
- Opens `~/.config/contrail/proxy.yaml` in `$EDITOR` (or `$VISUAL`, or falls back to `vi`)
- Creates the config file with defaults if it doesn't exist

---

## Docker Compose Integration

Contrail provides a `contrail-compose` shell function for direct Docker Compose interaction with automatic context awareness. This function is installed via `contrail init-shell` and delegates to Docker Compose with the correct project name and compose files.

For full documentation on `contrail-compose` and shell integration, see the [Shell Integration Specification](../specs/shell-integration.md).

### `contrail compose-prefix`

Output the Docker Compose command prefix for the current context. This is primarily used by the `contrail-compose` shell function.

```bash
contrail compose-prefix [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

> **Note**: There is no `--flavor` flag for `compose-prefix`. Flavor changes require regeneration and can impact running applications. Use `contrail flavor set` to change the active flavor before running `contrail-compose`.

**Output**:
```bash
$ cd ~/workspaces/dev/app-one
$ contrail compose-prefix
docker compose -p dev-app-one -f '/home/user/workspaces/dev/app-one/docker-compose.yaml' -f '/home/user/workspaces/dev/.generated/app-one.override.yaml'
```

**Error Handling**:

If context cannot be resolved, exits with code 5:
```bash
$ cd ~
$ contrail compose-prefix
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly with -w and -a flags

Available workspaces: contrail workspace list
```

---

### `contrail init-shell`

Output shell integration script for the specified shell. This script provides the `contrail-compose` function and its completion.

```bash
contrail init-shell {bash|zsh|fish}
```

**Installation**:
```bash
# Bash
contrail init-shell bash >> ~/.bashrc

# Zsh
contrail init-shell zsh >> ~/.zshrc

# Fish
contrail init-shell fish >> ~/.config/fish/conf.d/contrail.fish
```

**Provides**:
- `contrail-compose` function with context-aware Docker Compose passthrough
- Tab completion for `contrail-compose` that delegates to Docker's completion
- Automatic resolution of workspace, app, and compose files

**Example Usage** (after installation):
```bash
# From within an application directory
$ cd ~/workspaces/dev/app-one
$ contrail-compose exec php bash
$ contrail-compose logs -f

# From workspace root with explicit app
$ cd ~/workspaces/dev
$ contrail-compose -a app-two ps

# From anywhere with explicit workspace and app
$ contrail-compose -w dev -a app-one up -d
```

---

## Top-Level Aliases

For common operations, these aliases are provided:

| Alias | Equivalent |
|-------|------------|
| `contrail up` | `contrail workspace up` |
| `contrail down` | `contrail workspace down` |
| `contrail ps` | `contrail workspace status` |
| `contrail generate` | `contrail workspace generate` |

All aliases accept the same flags as their full forms and support context detection.

---

## Utility Commands

### `contrail validate`

Validate configuration files.

```bash
contrail validate [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Behavior**:
- Validates YAML syntax
- Checks schema compliance
- Verifies referenced files exist
- Reports errors with file locations

---

### `contrail doctor`

Check system health and requirements.

```bash
contrail doctor
```

**Output**:
```
Checking Contrail environment...

✓ Docker: running (version 24.0.7)
✓ Docker Compose: available (version 2.23.0)
✓ Proxy network: created
✓ Traefik: running
✓ Config directory: ~/.config/contrail
✓ Domain resolution: contrail.test → 127.0.0.1
✓ Workspace domains:
  - dev-app-one-web.contrail.test → 127.0.0.1
  - dev-app-two-api.contrail.test → 127.0.0.1

All checks passed.
```

**DNS checking behavior**:
- Uses the **system DNS resolver** (respects `/etc/hosts` and `/etc/resolv.conf`)
- Timeout: **5 seconds** per query
- Checks base domain (`contrail.test`) resolution
- If workspaces exist, checks all public proxied hostnames from workspace manifests
- If no workspaces exist, checks a test subdomain (`check-{timestamp}.contrail.test`) to verify wildcard configuration

**Offline/air-gapped environments**: DNS checks may fail in environments without network access. Use `/etc/hosts` entries or a local dnsmasq configuration for offline development.

**Wildcard DNS warning** (when base resolves but subdomains don't):
```
⚠ Wildcard DNS not configured. Individual hostnames may not resolve.
  Configure dnsmasq: address=/contrail.test/127.0.0.1
```

---

### `contrail open`

Open a service URL in the default browser.

```bash
contrail open [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |
| `--service` | Specific exported service (default: first web service) |

---

### `contrail urls`

List all accessible URLs for a workspace.

```bash
contrail urls [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**:
```
APP        SERVICE  URL
app-one    web      https://dev-app-one-web.contrail.test
app-two    web      https://dev-app-two-web.contrail.test
app-two    api      https://dev-app-two-api.contrail.test
```

---

## Output Formats

All list and show commands support multiple output formats:

| Flag | Format | Description |
|------|--------|-------------|
| (default) | Table | Human-readable aligned columns |
| `--json` | JSON | Machine-readable JSON |
| `--yaml` | YAML | Machine-readable YAML |
| `--quiet` | Names | Just names/IDs, one per line |

**Example**:
```bash
$ contrail workspace list --json
[
  {"name": "dev", "apps": 3, "status": "running"},
  {"name": "review", "apps": 3, "status": "stopped"}
]

$ contrail workspace list --quiet
dev
review
```

---

## Shell Completion

Contrail provides two types of shell integration:

1. **Standard CLI completion**: Completions for the `contrail` command itself
2. **Docker Compose passthrough**: The `contrail-compose` function with delegated completion

For complete documentation on shell integration, including the `contrail-compose` function and its completion delegation, see the [Shell Integration Specification](../specs/shell-integration.md).

### Standard CLI Completion

Install completion scripts for the `contrail` command:

```bash
# Bash
contrail completion bash > /etc/bash_completion.d/contrail

# Zsh
contrail completion zsh > "${fpath[1]}/_contrail"

# Fish
contrail completion fish > ~/.config/fish/completions/contrail.fish
```

**Completion Features**:
- Command and subcommand completion
- Flag completion with descriptions
- Dynamic completion for `--workspace` (from `contrail workspace list -q`)
- Dynamic completion for `--app` (from `contrail app list -q`)
- Flavor name completion
- Port number completion

### Full Shell Integration

For both CLI completion and `contrail-compose` support:

```bash
# Bash - includes contrail completion + contrail-compose function
contrail init-shell bash >> ~/.bashrc

# Zsh
contrail init-shell zsh >> ~/.zshrc

# Fish
contrail init-shell fish >> ~/.config/fish/conf.d/contrail.fish
```

This provides:
- All standard `contrail` CLI completions
- The `contrail-compose` shell function
- Tab completion for `contrail-compose` that delegates to Docker's completion

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CONTRAIL_CONFIG` | Path to global config file | `~/.config/contrail/proxy.yaml` |
| `CONTRAIL_WORKSPACE` | Default workspace (overrides context detection) | — |
| `NO_COLOR` | Disable colored output | — |
| `CONTRAIL_DEBUG` | Enable debug logging | — |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments or flags |
| 3 | Configuration error |
| 4 | Docker/Compose error |
| 5 | Context detection failed (workspace/app not found) |

---

## Error Messages

### Docker Not Available

Commands that require Docker check for availability upfront. If Docker is not installed or not running:

```
Error: Docker is not installed or not running.
Run 'contrail doctor' for setup guidance.
```

Exit code: 4

---

## Examples

### New workspace from scratch

```bash
mkdir ~/workspaces && cd ~/workspaces
contrail workspace init --workspace=dev
cd dev
contrail app add --app=frontend --repo=git@github.com:org/frontend.git
contrail app add --app=backend --repo=git@github.com:org/backend.git
contrail up
contrail urls
```

### Promote existing project to workspace

```bash
cd ~/my-docker-project
contrail workspace init --workspace=dev
contrail app init --app=myapp
# Edit application.yaml to define exported_services
contrail up
```

### Daily development workflow

```bash
cd ~/workspaces/dev/frontend

# Start your day
contrail up                    # Brings up entire dev workspace

# Work on frontend
contrail-compose logs -f       # Tail frontend logs (context detected)
contrail app restart           # Restart after changes

# Direct Docker Compose interaction (uses contrail-compose)
contrail-compose exec php bash          # Shell into container
contrail-compose exec php php artisan   # Run artisan command

# Check on another app
contrail app status -a backend
contrail-compose -a backend logs --tail=50
contrail-compose -a backend exec node npm test

# End of day
contrail down
```

### Direct Docker Compose operations

```bash
# contrail-compose provides context-aware docker compose access
cd ~/workspaces/dev/app-one

# These are equivalent:
contrail-compose exec php bash
# ...to running:
docker compose -p dev-app-one \
  -f ~/workspaces/dev/app-one/docker-compose.yaml \
  -f ~/workspaces/dev/.generated/app-one.override.yaml \
  exec php bash

# Target different app from workspace root
cd ~/workspaces/dev
contrail-compose -a app-two logs -f php

# Full docker compose functionality with tab completion
contrail-compose build --no-cache php
contrail-compose run --rm php composer install
```

### Switch application flavor

```bash
contrail flavor list -a backend
contrail flavor set full -a backend
contrail app restart -a backend
```

### Manage ports

```bash
contrail port list
contrail port list --verbose   # Check bind status
contrail port gc               # Clean up stale assignments
```
