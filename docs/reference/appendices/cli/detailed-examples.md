# CLI Detailed Examples

Extended examples for Contrail CLI commands.

<!-- Source: specs/contrail-cli-reference.md -->

---

## Workspace Lifecycle Examples

### New Workspace from Scratch

```bash
# Create workspace directory
mkdir ~/workspaces && cd ~/workspaces

# Initialize a new workspace
contrail workspace init --workspace=dev

# Navigate into the workspace
cd dev

# Add applications with git repositories
contrail app add --app=frontend --repo=git@github.com:org/frontend.git
contrail app add --app=backend --repo=git@github.com:org/backend.git

# Start all applications
contrail up

# View all accessible URLs
contrail urls
```

### Promote Existing Project to Workspace

```bash
# Navigate to your existing Docker Compose project
cd ~/my-docker-project

# Initialize as a workspace
contrail workspace init --workspace=dev

# Initialize the application configuration
contrail app init --app=myapp

# Edit application.yaml to define exported_services
# (manual step)

# Start the workspace
contrail up
```

---

## Daily Development Workflow

### Starting Your Day

```bash
cd ~/workspaces/dev/frontend

# Start the entire dev workspace
contrail up

# Check status
contrail workspace status
```

### Working on an Application

```bash
# Tail logs for the current application (context detected)
contrail-compose logs -f

# Restart after code changes
contrail app restart

# Shell into a container
contrail-compose exec php bash

# Run a command in a container
contrail-compose exec php php artisan migrate
```

### Checking Another Application

```bash
# Check status of a different app (from anywhere in workspace)
contrail app status -a backend

# View logs for a different app
contrail-compose -a backend logs --tail=50

# Run tests in a different app
contrail-compose -a backend exec node npm test
```

### End of Day

```bash
# Stop all applications
contrail down
```

---

## Direct Docker Compose Operations

The `contrail-compose` function provides context-aware Docker Compose access:

```bash
cd ~/workspaces/dev/app-one

# These are equivalent:
contrail-compose exec php bash

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
contrail-compose -a app-two logs -f php

# Build without cache
contrail-compose build --no-cache php

# Run one-off commands
contrail-compose run --rm php composer install
```

---

## Flavor Management Examples

### Listing and Switching Flavors

```bash
# List available flavors for an application
contrail flavor list -a backend

# Output:
# NAME     COMPOSE FILES                                    ACTIVE
# lite     docker-compose.yaml
# full     docker-compose.yaml, docker-compose.worker.yaml  ✓
# debug    docker-compose.yaml, docker-compose.debug.yaml

# Switch to a different flavor
contrail flavor set lite -a backend

# Apply the change (if app is running)
contrail app restart -a backend
```

### Flavor Change Scenarios

| Scenario | Command |
|----------|---------|
| Flavor adds/removes services | `contrail up` (starts new services, stops orphaned services) |
| Flavor changes environment or config | `contrail app restart -a APP` |

---

## Port Management Examples

### Viewing Port Assignments

```bash
# List all assigned ports
contrail port list

# Output:
# PORT   WORKSPACE  APP      SERVICE  STATUS
# 5432   dev        app-one  db       assigned
# 5433   dev        app-two  db       assigned
# 5434   review     app-one  db       assigned
# 6379   dev        app-one  cache    assigned

# With bind status check
contrail port list --verbose

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
contrail port gc --dry-run

# Actually release stale ports
contrail port gc
```

### Manual Port Operations

```bash
# View details about a specific port
contrail port show 5432

# Manually release a port
contrail port release 5432

# Force release even if in use
contrail port release 5432 --force

# Manually assign a port (advanced)
contrail port assign 5432 dev/app-one/db
```

---

## Proxy Management Examples

### Initial Setup

```bash
# Bootstrap the proxy
contrail proxy init

# Output:
# Created proxy configuration at ~/.config/contrail/proxy/
#
# Next steps:
#   1. Configure DNS for *.contrail.test -> 127.0.0.1
#      (See: contrail doctor for DNS verification)
#   2. Start the proxy:
#      contrail proxy up

# Start the proxy
contrail proxy up
```

### Custom Domain Setup

```bash
# Use a custom domain
contrail proxy init --domain mydev.local

# Output:
# Created proxy configuration at ~/.config/contrail/proxy/
# Domain set to: mydev.local
```

### Recovery from Broken Configuration

```bash
# Force regeneration of proxy config
contrail proxy init --force

# Output:
# Backed up existing configuration to ~/.config/contrail/proxy.backup.20241230/
# Created proxy configuration at ~/.config/contrail/proxy/
```

### Checking Proxy Status

```bash
contrail proxy status

# Output (dashboard enabled):
# Proxy: running
# Network: contrail-proxy (created)
# Dashboard: http://localhost:8080
# Entrypoints:
#   - web: :80
#   - websecure: :443

# Output (dashboard disabled):
# Proxy: running
# Network: contrail-proxy (created)
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
contrail config show

# Output:
# proxy:
#   domain: contrail.test
# paths:
#   global_config: ~/.config/contrail/proxy.yaml
#   global_state: ~/.config/contrail/state.yaml

# Get a specific value
contrail config get proxy.domain
# contrail.test

# Show file locations
contrail config path
# Global config: ~/.config/contrail/proxy.yaml
# Global state:  ~/.config/contrail/state.yaml
```

### Modifying Configuration

```bash
# Set a configuration value
contrail config set proxy.domain local.test

# Edit configuration in your default editor
contrail config edit
```

---

## Output Format Examples

### JSON Output

```bash
contrail workspace list --json

# Output:
# [
#   {"name": "dev", "apps": 3, "status": "running"},
#   {"name": "review", "apps": 3, "status": "stopped"}
# ]
```

### YAML Output

```bash
contrail workspace list --yaml

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
contrail workspace list --quiet

# Output:
# dev
# review

# Useful in scripts
for ws in $(contrail workspace list -q); do
  echo "Processing workspace: $ws"
  contrail workspace status -w "$ws"
done
```

---

## Validation and Diagnostics

### Validating Configuration

```bash
contrail validate

# Or for a specific workspace/app
contrail validate -w dev -a frontend
```

### System Health Check

```bash
contrail doctor

# Output:
# Checking Contrail environment...
#
# ✓ Docker: running (version 24.0.7)
# ✓ Docker Compose: available (version 2.23.0)
# ✓ Proxy network: created
# ✓ Traefik: running
# ✓ Config directory: ~/.config/contrail
# ✓ Domain resolution: contrail.test -> 127.0.0.1
# ✓ Workspace domains:
#   - dev-app-one-web.contrail.test -> 127.0.0.1
#   - dev-app-two-api.contrail.test -> 127.0.0.1
#
# All checks passed.
```

---

## Browser and URL Commands

### Opening Services in Browser

```bash
# Open the first web service for current context
contrail open

# Open a specific service
contrail open --service api

# Open for a specific app
contrail open -a backend --service api
```

### Listing All URLs

```bash
contrail urls

# Output:
# APP        SERVICE  URL
# app-one    web      https://dev-app-one-web.contrail.test
# app-two    web      https://dev-app-two-web.contrail.test
# app-two    api      https://dev-app-two-api.contrail.test
```

---

## Related Documents

- [CLI Reference](../../cli.md)
- [Error Messages](./error-messages.md)
