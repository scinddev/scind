# Generated Override Files Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail generates Docker Compose override files to integrate applications into workspaces. These files add network configuration, labels, and environment variables without modifying the application's original compose files. This enables the "pure overlay" design where applications remain unaware of the workspace system.

**Related Documents**:
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [Docker Labels](./docker-labels.md)
- [Environment Variables](./environment-variables.md)
- [Naming Conventions](./naming-conventions.md)

**Appendices**:
- [Complete Override Example](./appendices/generated-override-files/complete-override-example.yaml)

---

## Behavior

### File Location

Override files are generated in the workspace's `.generated/` directory:

```
{workspace}/.generated/{application-name}.override.yaml
```

Example:
```
~/workspaces/dev/.generated/app-one.override.yaml
~/workspaces/dev/.generated/app-two.override.yaml
```

### Generation Trigger

Override files are generated:
- During `workspace generate`
- Automatically during `workspace up` when source files have changed (staleness detection)
- When `contrail flavor set` changes the active flavor

### Gitignore

The `.generated/` directory should be gitignored. Override files are transient and regenerated as needed.

---

## Data Schema

### File Structure

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + {app}/application.yaml
# Flavor: {flavor}
# Generated: {timestamp}

name: {workspace}-{app}                    # Explicit project name

services:
  {service}:
    networks:
      {workspace}-internal:
        aliases:
          - {app}-{service}
      contrail-proxy: {}                   # For proxied services only
    labels:
      # Traefik labels (for proxied services)
      # Contrail context labels
      # Contrail export labels
    environment:
      # Contrail environment variables
    ports:                                 # For assigned type only
      - "{host}:{container}"

networks:
  {workspace}-internal:
    external: true
  contrail-proxy:
    external: true                         # Only if proxied services exist
```

### Header Comments

Each generated file includes a header comment with:
- Warning not to edit directly
- Source files used for generation
- Active flavor name
- Generation timestamp

---

## Examples

### Example 1: Complete Override File

For application `app-one` in workspace `dev` with:
- Proxied web service (HTTPS)
- Assigned database port

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + app-one/application.yaml
# Flavor: default
# Generated: 2024-12-27T10:30:00Z

name: dev-app-one                         # Explicit project name to prevent conflicts

services:
  web:
    networks:
      dev-internal:
        aliases:
          - app-one-web
      contrail-proxy: {}                   # Connected to proxy for Traefik routing
    labels:
      # Traefik HTTPS router
      - "traefik.enable=true"
      - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.contrail.test`)"
      - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
      - "traefik.http.routers.dev-app-one-web-https.tls=true"
      - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"
      # Contrail context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/home/user/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/home/user/workspaces/dev/app-one"
      # Contrail export labels
      - "contrail.export.web.host=dev-app-one-web.contrail.test"
      - "contrail.export.web.proxy.https.visibility=public"
      - "contrail.export.web.proxy.https.url=https://dev-app-one-web.contrail.test"
    environment:
      - CONTRAIL_WORKSPACE_NAME=dev
      - CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
      - CONTRAIL_APP_ONE_WEB_PORT=443
      - CONTRAIL_APP_ONE_WEB_SCHEME=https
      - CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test

  postgres:
    ports:
      - "5432:5432"                        # host:container - assigned port mapping
    networks:
      dev-internal:
        aliases:
          - app-one-db
    labels:
      # Contrail context labels
      - "contrail.workspace.name=dev"
      - "contrail.workspace.path=/home/user/workspaces/dev"
      - "contrail.app.name=app-one"
      - "contrail.app.path=/home/user/workspaces/dev/app-one"
      # Contrail export labels
      - "contrail.export.db.host=app-one-db"
      - "contrail.export.db.port.5432.visibility=protected"
      - "contrail.export.db.port.5432.assigned=5432"
    environment:
      - CONTRAIL_WORKSPACE_NAME=dev
      - CONTRAIL_APP_ONE_DB_HOST=app-one-db
      - CONTRAIL_APP_ONE_DB_PORT=5432

networks:
  dev-internal:
    external: true
  contrail-proxy:
    external: true
```

---

## Generation Process

### Step-by-Step Algorithm

1. **Resolve flavor** for the application (CLI → state → default_flavor → "default")

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

7. **Collect all exported services** for the application

8. **Generate override content**:
   - Add project name
   - For each service with exported services:
     - Add networks with aliases
     - Add Traefik labels (for proxied)
     - Add Contrail labels
     - Add environment variables
     - Add port mappings (for assigned)
   - Add external network declarations

9. **Write override file** to `.generated/{app}.override.yaml`

10. **Update state file** with resolved flavor

11. **Update manifest** with computed values

---

## Merge Behavior

Docker Compose files are merged in this order:

```bash
docker compose \
  -f base.yaml \
  -f .generated/app.override.yaml \
  -f overrides/app.yaml  # Optional manual override
```

### Merge Semantics

| Element | Behavior |
|---------|----------|
| `services.{name}.networks` | Merged (new networks added) |
| `services.{name}.labels` | Appended |
| `services.{name}.environment` | Merged (later values override) |
| `services.{name}.ports` | Merged (new ports added) |
| `networks` | Merged |
| `name` | Last value wins (override sets project name) |

---

## Manual Override Files

**Location**: `{workspace}/overrides/{application-name}.yaml`

Optional. If present, merged after the generated override file. Useful for workspace-specific customizations.

### Preservation Guarantee

Files in `{workspace}/overrides/` are **never modified by Contrail**. They persist across all regeneration operations, including `workspace generate --force`.

### Example Manual Override

```yaml
# Manual overrides for app-one in dev workspace
services:
  web:
    environment:
      - DEBUG=true
      - FEATURE_FLAG_X=enabled
    labels:
      - "traefik.http.routers.dev-app-one-web-https.middlewares=dev-auth@docker"

  postgres:
    volumes:
      - ./local-db-init:/docker-entrypoint-initdb.d:ro
```

---

## Edge Cases

### Assigned Port Conflict

**Scenario**: Previously assigned port is now unavailable at startup.

**Behavior**: Error with guidance:
```
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

### --force Flag Behavior

**Scenario**: User runs `workspace generate --force`.

**Behavior**:
- Regenerates all override files regardless of staleness
- **Preserves existing port assignments** (does not reassign ports)
- Does not touch files in `overrides/` directory

### Flavor Change While Running

**Scenario**: User changes flavor while application is running.

**Behavior**:
- Override file is immediately regenerated
- Warning displayed:
  ```
  Warning: Application "app-name" is currently running.
  The new flavor has been applied to the configuration, but running
  containers still use the previous flavor.

  To apply the flavor change:
    contrail app restart -a app-name
  ```

---

## Staleness Detection

Override files are considered stale if any source file has a newer modification time:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files (e.g., `docker-compose.yaml`)

### Behavior

- `workspace up` automatically regenerates stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Use `--no-generate` to skip regeneration entirely

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Missing compose file | FILE_NOT_FOUND | `Flavor "{name}" references non-existent file: {file}` | Fix compose_files list |
| Invalid service reference | VALIDATION | `Exported service references non-existent Compose service` | Fix service reference |
| Write permission denied | FILE_SYSTEM | `Cannot write to .generated/: permission denied` | Fix directory permissions |
| YAML syntax error | YAML_PARSE | `Invalid YAML in {file}: {error}` | Fix YAML syntax |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
