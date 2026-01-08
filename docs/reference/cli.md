# Scind CLI Reference

**Date**: January 2025
**Status**: Design Phase

This document is the authoritative reference for Scind's command-line interface. It defines command structure, arguments, flags, and behaviors.

---

## Command Structure

Scind follows a **resource-first** pattern inspired by Docker and Mutagen:

```
scind [resource] [action] [options...]
```

All targeting uses **options** rather than positional arguments:

```bash
# Good: Options-based
scind app status --workspace=dev --app=frontend

# With context detection (from current directory)
scind app status
```

---

## Context Detection

Scind automatically detects workspace and application context from the current directory, reducing the need for explicit `--workspace` and `--app` flags.

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
   ~/workspaces/dev/frontend/src/components/
                   │        │
                   │        └── application.yaml → app = "frontend"
                   │
                   └── workspace.yaml → workspace = "dev"
   ```

4. **Explicit flags override detection**: `--workspace` and `--app` always take precedence over context detection
   - When any `--app` flag is specified, context-detected application is **completely ignored**
   - This applies even when multiple `-a` flags are used
   - To start multiple specific apps: `scind up -a frontend -a backend`

5. **Global commands ignore context**: `port`, `proxy`, and `config` commands don't use directory context

### Flag Override Behavior

When explicit flags are provided, they **replace** (not add to) context detection:

```bash
# From within frontend directory (context would detect frontend)
$ cd ~/workspaces/dev/frontend

# This starts ONLY backend, not both frontend and backend
$ scind up -a backend
# Starting: backend
# (frontend from context is ignored)

# To start multiple apps, list them all explicitly
$ scind up -a frontend -a backend
# Starting: frontend, backend
```

This "explicit replaces context" behavior ensures predictable results—when you specify apps explicitly, you get exactly what you asked for.

### Edge Cases

**Nested vendor packages**: If working in `frontend/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored. The workspace's `frontend/application.yaml` is found first when walking toward the workspace root.

**Workspace within workspace**: If a test fixture has its own `workspace.yaml` nested inside a workspace (e.g., for integration tests), the closest `workspace.yaml` wins—this is the test fixture, which is the expected behavior.

### Context Feedback

When context is detected, commands indicate what was found:

```bash
$ cd ~/workspaces/dev/frontend
$ scind app status
# Using workspace: dev (from ../workspace.yaml)
# Using app: frontend (from ./application.yaml)

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
$ scind app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: scind workspace init --workspace=NAME
```

**Neither workspace nor application found**:
```bash
$ cd ~
$ scind app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: scind app status --workspace=NAME --app=NAME

Available workspaces: scind workspace list
```

**Workspace found but no application context** (for app-specific commands):
```bash
$ cd ~/workspaces/dev
$ scind app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: scind app status --app=NAME

Available apps in 'dev': frontend, backend, shared-db
```

---

## Resources

Scind manages these resource types:

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
| `--version` | | Show Scind version |

### Version Information

Both `scind --version` and `scind version` display version information. The `--version` flag is available on all commands.

```bash
$ scind --version
scind version 0.2.3

$ scind version
scind version 0.2.3
```

### Output Behavior

**Progress output**: Multi-application operations show per-application progress by default:
```
Starting frontend... done
Starting backend... done
Starting shared-db... done
```

**`--quiet` behavior**:
- Suppresses context indicators ("Using workspace: dev")
- Suppresses progress messages ("Starting frontend... done")
- Status commands output just the value: `running`
- List commands output names only, one per line
- Errors are always shown regardless of `--quiet`

```bash
# Normal output
$ scind workspace list
NAME     APPS  STATUS   PATH
dev      3     running  ~/workspaces/dev
staging  2     stopped  ~/workspaces/staging

# Quiet output (machine-readable)
$ scind workspace list -q
dev
staging

# Status with quiet
$ scind app status -q
running
```

---

## Workspace Commands

Manage workspace lifecycle and orchestration.

### `scind workspace list`

List all registered workspaces.

```bash
scind workspace list [flags]
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

**Discovery mechanism**: Reads from the workspace registry (`~/.config/scind/workspaces.yaml`). If the registry is missing, automatically attempts to rebuild from Docker container labels.

---

### `scind workspace prune`

Remove stale entries from the workspace registry.

```bash
scind workspace prune [flags]
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
scind workspace prune
# Removed: old-project (path /home/user/old-project no longer exists)
# Registry: 3 workspaces remaining
```

---

### `scind workspace show`

Show detailed information about a workspace, including the computed manifest.

```bash
scind workspace show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**: Workspace configuration, applications, computed hostnames, port assignments.

**Example**:
```bash
scind workspace show --workspace=dev
```

---

### `scind workspace init`

Initialize a new workspace.

```bash
scind workspace init [flags]
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
- **Registers the workspace** in `~/.config/scind/workspaces.yaml`
- **Fails if name is already registered** to a different path (enforces workspace name uniqueness)

**Example** (new workspace):
```bash
scind workspace init --workspace=dev
# Creates ./dev/workspace.yaml
# Registers "dev" -> ./dev in workspace registry
```

**Example** (current directory):
```bash
cd ~/my-project
scind workspace init --workspace=dev
# Creates ./workspace.yaml with name: dev
# Registers "dev" -> ~/my-project in workspace registry
```

**Example** (name collision):
```bash
scind workspace init --workspace=dev
# Error: Workspace "dev" already registered at ~/workspaces/dev
# Use a different name, or run `scind workspace prune` if that path no longer exists
```

---

### `scind workspace clone`

Clone all application repositories defined in the workspace.

```bash
scind workspace clone [flags]
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
scind workspace clone
# Skipping myapp: application is workspace root (path: .)
# Cloned: backend -> ./backend
```

---

### `scind workspace generate`

Generate or regenerate Docker Compose override files.

```bash
scind workspace generate [flags]
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

### `scind workspace up`

Bring up a workspace (generate overrides if needed, create networks, start containers).

```bash
scind workspace up [flags]
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
4. Ensure `scind-proxy` network exists and proxy is running
5. For each application (or specified apps):
   - Resolve active flavor
   - Execute `docker compose up -d` with appropriate files
6. Report status

**Example**:
```bash
scind workspace up --workspace=dev
scind workspace up -a frontend -a backend  # With context
scind up  # Alias, with context
```

---

### `scind workspace down`

Tear down a workspace (stop containers, remove networks).

```bash
scind workspace down [flags]
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

### `scind workspace restart`

Restart a workspace or specific applications.

```bash
scind workspace restart [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Restart specific app(s) only (repeatable) |

**Behavior**: Equivalent to `down` followed by `up`. Volumes are always preserved (the internal `down` does not use `--volumes`).

---

### `scind workspace status`

Show the running status of a workspace.

```bash
scind workspace status [flags]
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

APPLICATION  FLAVOR   STATUS   SERVICES         URL
frontend    full     running  3/3 running      https://dev-frontend-web.scind.test
backend     full     running  5/5 running      https://dev-backend-api.scind.test
shared-db   default  stopped  0/2 running      —
```

---

### `scind workspace destroy`

Completely remove a workspace.

```bash
scind workspace destroy [flags]
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
6. Remove workspace from registry (`~/.config/scind/workspaces.yaml`)

**Warning**: This is destructive. Without `--force` or `--keep-apps`, prompts for confirmation showing what will be removed.

---

## Application Commands

Manage applications within workspaces.

### `scind app list`

List applications in a workspace.

```bash
scind app list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**:
```
NAME       FLAVOR   STATUS   SERVICES  PATH
frontend   full     running  3/3       ./frontend
backend    full     running  5/5       ./backend
shared-db  default  stopped  0/2       ./shared-db
```

---

### `scind app show`

Show detailed information about an application.

```bash
scind app show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Output**: Application configuration, exported services, current flavor, computed hostnames and ports.

**Example**:
```bash
scind app show --app=backend
```

---

### `scind app init`

Initialize an application configuration in the current directory.

```bash
scind app init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-a, --app` | Application name (default: current directory name) |

**Behavior**:
- Creates `application.yaml` in current directory
- Scans for existing `docker-compose.yaml` to suggest exported services
- Sets up default flavor pointing to existing compose file(s)

**Use case**: Promoting an existing Docker Compose project to a Scind application.

```bash
cd ~/my-project
scind app init --app=myapp
# Creates ./application.yaml
```

---

### `scind app add`

Add an application to a workspace.

```bash
scind app add [flags]
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
scind app add --app=api --repo=git@github.com:org/api.git
scind app add --app=main --path=.  # Use current directory
```

---

### `scind app remove`

Remove an application from a workspace.

```bash
scind app remove [flags]
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

### `scind app up`

Bring up a single application.

```bash
scind app up [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

Equivalent to `scind workspace up --app=NAME`.

---

### `scind app down`

Tear down a single application.

```bash
scind app down [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |
| `--volumes` | Also remove volumes |

Equivalent to `scind workspace down --app=NAME`.

---

### `scind app restart`

Restart a single application.

```bash
scind app restart [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Behavior**: Equivalent to `app down` followed by `app up`. Volumes are always preserved.

---

### `scind app status`

Show status of a single application.

```bash
scind app status [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

**Example**:
```bash
scind app status --app=frontend
```

---

## Flavor Commands

Manage application flavors (named configurations).

### `scind flavor list`

List available flavors for an application.

```bash
scind flavor list [flags]
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

### `scind flavor show`

Show the current active flavor for an application.

```bash
scind flavor show [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

---

### `scind flavor set`

Set the active flavor for an application.

```bash
scind flavor set <flavor> [flags]
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
| Flavor adds/removes services | `scind up` (starts new services, stops orphaned services) |
| Flavor changes environment or config | `scind app restart -a APP` |

**Example**:
```bash
# Change flavor (regenerates config, warns if running)
scind flavor set full --app=backend

# Apply to running application
scind app restart --app=backend
```

---

## Port Commands

Manage host port assignments for `assigned`-type services.

> **Port Types**: Services can expose `proxied` ports (through Traefik reverse proxy) or `assigned` ports (direct host binding). Port commands manage assigned ports only; proxied ports are handled automatically through the proxy. See [ADR-0007: Port Type System](../decisions/0007-port-type-system.md).

### `scind port list`

List all assigned ports.

```bash
scind port list [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Filter to specific workspace |
| `-v, --verbose` | Include bind status check |

**Output**:
```
PORT   WORKSPACE  APP        SERVICE  STATUS
5432   dev        shared-db  db       assigned
6379   dev        shared-db  cache    assigned
8080   dev        frontend   web      assigned
3000   dev        backend    api      assigned
```

With `--verbose`:
```
PORT   WORKSPACE  APP        SERVICE  STATUS    BOUND
5432   dev        shared-db  db       assigned  yes
6379   dev        shared-db  cache    assigned  yes
8080   dev        frontend   web      assigned  yes
3000   dev        backend    api      assigned  no
```

---

### `scind port show`

Show details about a specific port assignment.

```bash
scind port show <port>
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `port` | Port number to inspect |

---

### `scind port release`

Manually release a port assignment.

```bash
scind port release <port> [flags]
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

### `scind port assign`

Manually assign a port (advanced usage).

```bash
scind port assign <port> <workspace>/<app>/<service>
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `port` | Port number to assign |
| `target` | Target in format `workspace/app/service` |

---

### `scind port gc`

Garbage collect released and unbound ports.

```bash
scind port gc [flags]
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

### `scind port scan`

Re-scan and update port availability inventory.

```bash
scind port scan
```

**Behavior**:
- Checks all tracked ports for current bind status
- Updates `port_inventory` in global state
- Useful after external changes (other processes releasing ports)

---

## Proxy Commands

Manage the Traefik reverse proxy. The proxy is a shared instance that serves all workspaces on the host.

### `scind proxy init`

Bootstrap the proxy configuration. This creates the Traefik Docker Compose project at `~/.config/scind/proxy/`.

```bash
scind proxy init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing configuration (useful for recovery) |
| `--domain` | Set proxy domain (default: `scind.test`) |
| `--path` | Directory to create proxy in (default: `~/.config/scind/proxy/`) |

**Behavior**:
1. Check if proxy configuration already exists
   - If exists and no `--force`: error with message
   - If exists and `--force`: backup existing config and overwrite
2. Create proxy directory structure (`docker-compose.yaml`, `traefik.yaml`, `dynamic/`, `certs/`)
3. Create `scind-proxy` Docker network if it doesn't exist
4. Output next steps (DNS setup, starting proxy)

**Example**:
```bash
$ scind proxy init
Created proxy configuration at ~/.config/scind/proxy/

Next steps:
  1. Configure DNS for *.scind.test → 127.0.0.1
     (See: scind doctor for DNS verification)
  2. Start the proxy:
     scind proxy up

$ scind proxy init
Error: Proxy configuration already exists at ~/.config/scind/proxy/
Use --force to overwrite, or --path to create elsewhere.

$ scind proxy init --force
Backed up existing configuration to ~/.config/scind/proxy.backup.20241230/
Created proxy configuration at ~/.config/scind/proxy/
```

**Custom domain example**:
```bash
$ scind proxy init --domain mydev.local
Created proxy configuration at ~/.config/scind/proxy/
Domain set to: mydev.local
```

---

### `scind proxy up`

Start the Traefik proxy.

```bash
scind proxy up [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--recreate` | Recreate the proxy network even if it exists |

**Behavior**:
- Creates `scind-proxy` network if it doesn't exist
- Validates existing network configuration matches expected settings
- Starts Traefik container from proxy configuration
- If proxy configuration doesn't exist, runs `proxy init` first

**Network conflict handling**:
If the `scind-proxy` network exists but was created by a different tool or has incompatible settings, `proxy up` will warn:
```
Warning: Network 'scind-proxy' exists but may not have been created by Scind.
  Driver: bridge (expected: bridge) ✓
  Labels: scind.managed not found ⚠

Use 'scind proxy up --recreate' to recreate the network.
```

**Note**: Users rarely need to call this directly. `workspace up` automatically starts the proxy if it's not running.

---

### `scind proxy down`

Stop the Traefik proxy.

```bash
scind proxy down
```

---

### `scind proxy restart`

Restart the Traefik proxy.

```bash
scind proxy restart
```

---

### `scind proxy status`

Show proxy status.

```bash
scind proxy status
```

**Output** (dashboard enabled):
```
Proxy: running
Network: scind-proxy (created)
Dashboard: http://localhost:8080
Entrypoints:
  - web: :80
  - websecure: :443
```

**Output** (dashboard disabled):
```
Proxy: running
Network: scind-proxy (created)
Dashboard: disabled
Entrypoints:
  - web: :80
  - websecure: :443
```

**Note**: The dashboard URL reflects the configured port from `proxy.yaml` (`dashboard.port`, default 8080). If `dashboard.enabled` is false, the dashboard line shows "disabled".

---

## Config Commands

Manage Scind configuration.

### `scind config show`

Show all configuration values.

```bash
scind config show
```

**Output**:
```yaml
proxy:
  domain: scind.test
paths:
  global_config: ~/.config/scind/proxy.yaml
  global_state: ~/.config/scind/state.yaml
```

---

### `scind config get`

Get a specific configuration value.

```bash
scind config get <key>
```

**Example**:
```bash
scind config get proxy.domain
# scind.test
```

---

### `scind config set`

Set a configuration value.

```bash
scind config set <key> <value>
```

**Example**:
```bash
scind config set proxy.domain local.test
```

---

### `scind config path`

Show configuration file locations.

```bash
scind config path
```

**Output**:
```
Global config: ~/.config/scind/proxy.yaml
Global state:  ~/.config/scind/state.yaml
```

---

### `scind config edit`

Open the global configuration file in your default editor.

```bash
scind config edit
```

**Behavior**:
- Opens `~/.config/scind/proxy.yaml` in `$EDITOR` (or `$VISUAL`, or falls back to `vi`)
- Creates the config file with defaults if it doesn't exist

---

## Docker Compose Integration

Scind provides a `scind-compose` shell function for direct Docker Compose interaction with automatic context awareness. This function is installed via `scind init-shell` and delegates to Docker Compose with the correct project name and compose files.

For full documentation on `scind-compose` and shell integration, see the [Shell Integration Specification](../specs/shell-integration.md).

### `scind compose-prefix`

Output the Docker Compose command prefix for the current context. This is primarily used by the `scind-compose` shell function.

```bash
scind compose-prefix [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |

> **Note**: There is no `--flavor` flag for `compose-prefix`. Flavor changes require regeneration and can impact running applications. Use `scind flavor set` to change the active flavor before running `scind-compose`.

**Output**:
```bash
$ cd ~/workspaces/dev/frontend
$ scind compose-prefix
docker compose -p dev-frontend -f '/home/user/workspaces/dev/frontend/docker-compose.yaml' -f '/home/user/workspaces/dev/.generated/frontend.override.yaml'
```

**Error Handling**:

If context cannot be resolved, exits with code 5:
```bash
$ cd ~
$ scind compose-prefix
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly with -w and -a flags

Available workspaces: scind workspace list
```

---

### `scind init-shell`

Output shell integration script for the specified shell. This script provides the `scind-compose` function and its completion.

```bash
scind init-shell {bash|zsh|fish}
```

**Installation**:
```bash
# Bash
scind init-shell bash >> ~/.bashrc

# Zsh
scind init-shell zsh >> ~/.zshrc

# Fish
scind init-shell fish >> ~/.config/fish/conf.d/scind.fish
```

**Provides**:
- `scind-compose` function with context-aware Docker Compose passthrough
- Tab completion for `scind-compose` that delegates to Docker's completion
- Automatic resolution of workspace, app, and compose files

**Example Usage** (after installation):
```bash
# From within an application directory
$ cd ~/workspaces/dev/frontend
$ scind-compose exec node bash
$ scind-compose logs -f

# From workspace root with explicit app
$ cd ~/workspaces/dev
$ scind-compose -a backend ps

# From anywhere with explicit workspace and app
$ scind-compose -w dev -a frontend up -d
```

---

## Top-Level Aliases

For common operations, these aliases are provided:

| Alias | Equivalent |
|-------|------------|
| `scind up` | `scind workspace up` |
| `scind down` | `scind workspace down` |
| `scind ps` | `scind workspace status` |
| `scind generate` | `scind workspace generate` |

All aliases accept the same flags as their full forms and support context detection.

---

## Utility Commands

### `scind validate`

Validate configuration files.

```bash
scind validate [flags]
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

**Example**:
```bash
scind validate --workspace=dev
scind validate --workspace=dev --app=frontend
```

---

### `scind doctor`

Check system health and requirements.

```bash
scind doctor
```

**Output**:
```
Checking Scind environment...

✓ Docker: running (version 24.0.7)
✓ Docker Compose: available (version 2.23.0)
✓ Proxy network: created
✓ Traefik: running
✓ Config directory: ~/.config/scind
✓ Domain resolution: scind.test → 127.0.0.1
✓ Workspace domains:
  - dev-frontend-web.scind.test → 127.0.0.1
  - dev-backend-api.scind.test → 127.0.0.1

All checks passed.
```

**DNS checking behavior**:
- Uses the **system DNS resolver** (respects `/etc/hosts` and `/etc/resolv.conf`)
- Timeout: **5 seconds** per query
- Checks base domain (`scind.test`) resolution
- If workspaces exist, checks all public proxied hostnames from workspace manifests
- If no workspaces exist, checks a test subdomain (`check-{timestamp}.scind.test`) to verify wildcard configuration

**Offline/air-gapped environments**: DNS checks may fail in environments without network access. Use `/etc/hosts` entries or a local dnsmasq configuration for offline development.

**Wildcard DNS warning** (when base resolves but subdomains don't):
```
⚠ Wildcard DNS not configured. Individual hostnames may not resolve.
  Configure dnsmasq: address=/scind.test/127.0.0.1
```

---

### `scind open`

Open a service URL in the default browser.

```bash
scind open [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |
| `-a, --app` | Target application (or use context) |
| `--service` | Specific exported service (default: first web service) |

---

### `scind urls`

List all accessible URLs for a workspace.

```bash
scind urls [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `-w, --workspace` | Target workspace (or use context) |

**Output**:
```
APP        SERVICE  URL
frontend   web      https://dev-frontend-web.scind.test
backend    web      https://dev-backend-web.scind.test
backend    api      https://dev-backend-api.scind.test
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
$ scind workspace list --json
[
  {"name": "dev", "apps": 3, "status": "running"},
  {"name": "review", "apps": 3, "status": "stopped"}
]

$ scind workspace list --quiet
dev
review
```

---

### `scind completion`

Generate shell completion scripts.

```bash
scind completion {bash|zsh|fish}
```

**Subcommands**:

| Command | Description |
|---------|-------------|
| `completion bash` | Generate Bash completion script |
| `completion zsh` | Generate Zsh completion script |
| `completion fish` | Generate Fish completion script |

**Installation**:

```bash
# Bash
scind completion bash > /etc/bash_completion.d/scind

# Zsh
scind completion zsh > "${fpath[1]}/_scind"

# Fish
scind completion fish > ~/.config/fish/completions/scind.fish
```

---

## Shell Completion

Scind provides two types of shell integration:

1. **Standard CLI completion**: Completions for the `scind` command itself
2. **Docker Compose passthrough**: The `scind-compose` function with delegated completion

For complete documentation on shell integration, including the `scind-compose` function and its completion delegation, see the [Shell Integration Specification](../specs/shell-integration.md).

### Standard CLI Completion

Install completion scripts for the `scind` command:

```bash
# Bash
scind completion bash > /etc/bash_completion.d/scind

# Zsh
scind completion zsh > "${fpath[1]}/_scind"

# Fish
scind completion fish > ~/.config/fish/completions/scind.fish
```

**Completion Features**:
- Command and subcommand completion
- Flag completion with descriptions
- Dynamic completion for `--workspace` (from `scind workspace list -q`)
- Dynamic completion for `--app` (from `scind app list -q`)
- Flavor name completion
- Port number completion

### Full Shell Integration

For both CLI completion and `scind-compose` support:

```bash
# Bash - includes scind completion + scind-compose function
scind init-shell bash >> ~/.bashrc

# Zsh
scind init-shell zsh >> ~/.zshrc

# Fish
scind init-shell fish >> ~/.config/fish/conf.d/scind.fish
```

This provides:
- All standard `scind` CLI completions
- The `scind-compose` shell function
- Tab completion for `scind-compose` that delegates to Docker's completion

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SCIND_CONFIG` | Path to global config file | `~/.config/scind/proxy.yaml` |
| `SCIND_STATE_DIR` | State file directory | `~/.config/scind` |
| `SCIND_WORKSPACE` | Default workspace (overrides context detection) | — |
| `NO_COLOR` | Disable colored output | — |
| `SCIND_DEBUG` | Enable debug logging | — |

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
Run 'scind doctor' for setup guidance.
```

Exit code: 4

---

## Examples

### New workspace from scratch

```bash
mkdir ~/workspaces && cd ~/workspaces
scind workspace init --workspace=dev
cd dev
scind app add --app=frontend --repo=git@github.com:org/frontend.git
scind app add --app=backend --repo=git@github.com:org/backend.git
scind up
scind urls
```

### Promote existing project to workspace

```bash
cd ~/my-docker-project
scind workspace init --workspace=dev
scind app init --app=myapp
# Edit application.yaml to define exported_services
scind up
```

### Daily development workflow

```bash
cd ~/workspaces/dev/frontend

# Start your day
scind up                    # Brings up entire dev workspace

# Work on frontend
scind-compose logs -f       # Tail frontend logs (context detected)
scind app restart           # Restart after changes

# Direct Docker Compose interaction (uses scind-compose)
scind-compose exec node bash         # Shell into container
scind-compose exec node npm run dev  # Run npm command

# Check on another app
scind app status -a backend
scind-compose -a backend logs --tail=50
scind-compose -a backend exec node npm test

# End of day
scind down
```

### Direct Docker Compose operations

```bash
# scind-compose provides context-aware docker compose access
cd ~/workspaces/dev/frontend

# These are equivalent:
scind-compose exec node bash
# ...to running:
docker compose -p dev-frontend \
  -f ~/workspaces/dev/frontend/docker-compose.yaml \
  -f ~/workspaces/dev/.generated/frontend.override.yaml \
  exec node bash

# Target different app from workspace root
cd ~/workspaces/dev
scind-compose -a backend logs -f api

# Full docker compose functionality with tab completion
scind-compose build --no-cache php
scind-compose run --rm php composer install
```

### Switch application flavor

```bash
scind flavor list -a backend
scind flavor set full -a backend
scind app restart -a backend
```

### Manage ports

```bash
scind port list
scind port list --verbose   # Check bind status
scind port gc               # Clean up stale assignments
```
