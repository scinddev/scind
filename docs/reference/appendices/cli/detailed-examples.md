# CLI Detailed Examples

Extended examples for Scind CLI commands.

<!-- Source: specs/scind-cli-reference.md -->

---

## Workspace Lifecycle Examples

### New Workspace from Scratch

```bash
# Create workspace directory
mkdir ~/workspaces && cd ~/workspaces

# Initialize a new workspace
scind workspace init --workspace=dev

# Navigate into the workspace
cd dev

# Add applications with git repositories
scind app add --app=frontend --repo=git@github.com:org/frontend.git
scind app add --app=backend --repo=git@github.com:org/backend.git

# Start all applications
scind up

# View all accessible URLs
scind urls
```

### Promote Existing Project to Workspace

```bash
# Navigate to your existing Docker Compose project
cd ~/my-docker-project

# Initialize as a workspace
scind workspace init --workspace=dev

# Initialize the application configuration
scind app init --app=myapp

# Edit application.yaml to define exported_services
# (manual step)

# Start the workspace
scind up
```

---

## Daily Development Workflow

### Starting Your Day

```bash
cd ~/workspaces/dev/frontend

# Start the entire dev workspace
scind up

# Check status
scind workspace status
```

### Working on an Application

```bash
# Tail logs for the current application (context detected)
scind-compose logs -f

# Restart after code changes
scind app restart

# Shell into a container
scind-compose exec php bash

# Run a command in a container
scind-compose exec php php artisan migrate
```

### Checking Another Application

```bash
# Check status of a different app (from anywhere in workspace)
scind app status -a backend

# View logs for a different app
scind-compose -a backend logs --tail=50

# Run tests in a different app
scind-compose -a backend exec node npm test
```

### End of Day

```bash
# Stop all applications
scind down
```

---

## Direct Docker Compose Operations

The `scind-compose` function provides context-aware Docker Compose access:

```bash
cd ~/workspaces/dev/app-one

# These are equivalent:
scind-compose exec php bash

# ...to running:
docker compose -p dev-app-one \
  -f ~/workspaces/dev/app-one/docker-compose.yaml \
  -f ~/workspaces/dev/.generated/app-one.override.yaml \
  exec php bash
```

### Targeting Different Apps

```bash
# From workspace root
cd ~/workspaces/dev

# Target specific app with -a flag
scind-compose -a app-two logs -f php

# Build without cache
scind-compose build --no-cache php

# Run one-off commands
scind-compose run --rm php composer install
```

---

## Flavor Management Examples

### Listing and Switching Flavors

```bash
# List available flavors for an application
scind flavor list -a backend

# Output:
# NAME     COMPOSE FILES                                    ACTIVE
# lite     docker-compose.yaml
# full     docker-compose.yaml, docker-compose.worker.yaml  ✓
# debug    docker-compose.yaml, docker-compose.debug.yaml

# Switch to a different flavor
scind flavor set lite -a backend

# Apply the change (if app is running)
scind app restart -a backend
```

### Flavor Change Scenarios

| Scenario | Command |
|----------|---------|
| Flavor adds/removes services | `scind up` (starts new services, stops orphaned services) |
| Flavor changes environment or config | `scind app restart -a APP` |

---

## Port Management Examples

### Viewing Port Assignments

```bash
# List all assigned ports
scind port list

# Output:
# PORT   WORKSPACE  APP      SERVICE  STATUS
# 5432   dev        app-one  db       assigned
# 5433   dev        app-two  db       assigned
# 5434   review     app-one  db       assigned
# 6379   dev        app-one  cache    assigned

# With bind status check
scind port list --verbose

# Output:
# PORT   WORKSPACE  APP      SERVICE  STATUS    BOUND
# 5432   dev        app-one  db       assigned  yes
# 5433   dev        app-two  db       assigned  yes
# 5434   review     app-one  db       assigned  no
# 6379   dev        app-one  cache    assigned  yes
```

### Cleaning Up Stale Ports

```bash
# Check which ports would be released
scind port gc --dry-run

# Actually release stale ports
scind port gc
```

### Manual Port Operations

```bash
# View details about a specific port
scind port show 5432

# Manually release a port
scind port release 5432

# Force release even if in use
scind port release 5432 --force

# Manually assign a port (advanced)
scind port assign 5432 dev/app-one/db
```

---

## Proxy Management Examples

### Initial Setup

```bash
# Bootstrap the proxy
scind proxy init

# Output:
# Created proxy configuration at ~/.config/scind/proxy/
#
# Next steps:
#   1. Configure DNS for *.scind.test -> 127.0.0.1
#      (See: scind doctor for DNS verification)
#   2. Start the proxy:
#      scind proxy up

# Start the proxy
scind proxy up
```

### Custom Domain Setup

```bash
# Use a custom domain
scind proxy init --domain mydev.local

# Output:
# Created proxy configuration at ~/.config/scind/proxy/
# Domain set to: mydev.local
```

### Recovery from Broken Configuration

```bash
# Force regeneration of proxy config
scind proxy init --force

# Output:
# Backed up existing configuration to ~/.config/scind/proxy.backup.20241230/
# Created proxy configuration at ~/.config/scind/proxy/
```

### Checking Proxy Status

```bash
scind proxy status

# Output (dashboard enabled):
# Proxy: running
# Network: scind-proxy (created)
# Dashboard: http://localhost:8080
# Entrypoints:
#   - web: :80
#   - websecure: :443

# Output (dashboard disabled):
# Proxy: running
# Network: scind-proxy (created)
# Dashboard: disabled
# Entrypoints:
#   - web: :80
#   - websecure: :443
```

---

## Configuration Management Examples

### Viewing Configuration

```bash
# Show all configuration values
scind config show

# Output:
# proxy:
#   domain: scind.test
# paths:
#   global_config: ~/.config/scind/proxy.yaml
#   global_state: ~/.config/scind/state.yaml

# Get a specific value
scind config get proxy.domain
# scind.test

# Show file locations
scind config path
# Global config: ~/.config/scind/proxy.yaml
# Global state:  ~/.config/scind/state.yaml
```

### Modifying Configuration

```bash
# Set a configuration value
scind config set proxy.domain local.test

# Edit configuration in your default editor
scind config edit
```

---

## Output Format Examples

### JSON Output

```bash
scind workspace list --json

# Output:
# [
#   {"name": "dev", "apps": 3, "status": "running"},
#   {"name": "review", "apps": 3, "status": "stopped"}
# ]
```

### YAML Output

```bash
scind workspace list --yaml

# Output:
# - name: dev
#   apps: 3
#   status: running
# - name: review
#   apps: 3
#   status: stopped
```

### Quiet Output (Scripting)

```bash
# Just names, one per line
scind workspace list --quiet

# Output:
# dev
# review

# Useful in scripts
for ws in $(scind workspace list -q); do
  echo "Processing workspace: $ws"
  scind workspace status -w "$ws"
done
```

---

## Validation and Diagnostics

### Validating Configuration

```bash
scind validate

# Or for a specific workspace/app
scind validate -w dev -a frontend
```

### System Health Check

```bash
scind doctor

# Output:
# Checking Scind environment...
#
# ✓ Docker: running (version 24.0.7)
# ✓ Docker Compose: available (version 2.23.0)
# ✓ Proxy network: created
# ✓ Traefik: running
# ✓ Config directory: ~/.config/scind
# ✓ Domain resolution: scind.test -> 127.0.0.1
# ✓ Workspace domains:
#   - dev-app-one-web.scind.test -> 127.0.0.1
#   - dev-app-two-api.scind.test -> 127.0.0.1
#
# All checks passed.
```

---

## Browser and URL Commands

### Opening Services in Browser

```bash
# Open the first web service for current context
scind open

# Open a specific service
scind open --service api

# Open for a specific app
scind open -a backend --service api
```

### Listing All URLs

```bash
scind urls

# Output:
# APP        SERVICE  URL
# app-one    web      https://dev-app-one-web.scind.test
# app-two    web      https://dev-app-two-web.scind.test
# app-two    api      https://dev-app-two-api.scind.test
```

---

## Related Documents

- [CLI Reference](../../cli.md)
- [Error Messages](./error-messages.md)
