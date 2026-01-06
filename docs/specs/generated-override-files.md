# Generated Override Files Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail generates Docker Compose override files to integrate applications into workspaces. These files add network configuration, labels, and environment variables without modifying the application's original compose files. This enables the "pure overlay" design where applications remain unaware of the workspace system.

**Location**: `{workspace}/.generated/{application-name}.override.yaml`

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

### Generated Override Structure

The generated override file includes:

- **Project name**: Explicit project name (`name: {workspace}-{app}`) to prevent conflicts
- **Service networks**: Connection to workspace internal network and proxy network
- **Network aliases**: Internal aliases for service discovery
- **Traefik labels**: Router configuration for proxied services
- **Contrail labels**: Context and export labels for discovery
- **Environment variables**: Service discovery variables for all exported services
- **Port mappings**: Host port bindings for assigned-type services

### File Header

Each generated file includes a header with generation metadata:

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + {app}/application.yaml
# Flavor: {flavor}
# Generated: {timestamp}
```

### Networks Section

The generated override declares both workspace networks as external:

```yaml
networks:
  {workspace}-internal:
    external: true
  contrail-proxy:
    external: true
```

### Service Configuration

#### Network Connections

Each exported service is connected to the workspace internal network with its alias, and to the proxy network for Traefik routing:

```yaml
services:
  {service}:
    networks:
      {workspace}-internal:
        aliases:
          - {app}-{export}
      contrail-proxy: {}
```

#### Traefik Labels

For proxied services, Traefik routing labels are generated:

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.{workspace}-{app}-{export}-{protocol}.rule=Host(`{hostname}`)"
  - "traefik.http.routers.{workspace}-{app}-{export}-{protocol}.entrypoints={entrypoint}"
  - "traefik.http.routers.{workspace}-{app}-{export}-{protocol}.tls=true"
  - "traefik.http.services.{workspace}-{app}-{export}-{protocol}.loadbalancer.server.port={port}"
```

#### Assigned Port Mappings

For assigned-type services, direct host port mappings are generated:

```yaml
services:
  {service}:
    ports:
      - "{host_port}:{container_port}"
```

---

## Examples

See [appendices/generated-override-files/complete-override-example.yaml](./appendices/generated-override-files/complete-override-example.yaml) for a complete example.

### Example: Basic Override Structure

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + app-one/application.yaml
# Flavor: default
# Generated: 2024-12-27T10:30:00Z

name: dev-app-one

services:
  web:
    networks:
      dev-internal:
        aliases:
          - app-one-web
      contrail-proxy: {}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.dev-app-one-web-https.rule=Host(`dev-app-one-web.contrail.test`)"
      - "traefik.http.routers.dev-app-one-web-https.entrypoints=websecure"
      - "traefik.http.routers.dev-app-one-web-https.tls=true"
      - "traefik.http.services.dev-app-one-web-https.loadbalancer.server.port=443"
      - "contrail.workspace.name=dev"
      - "contrail.app.name=app-one"
    environment:
      - CONTRAIL_WORKSPACE_NAME=dev
      - CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
      - CONTRAIL_APP_ONE_WEB_PORT=443

  postgres:
    ports:
      - "5432:5432"
    networks:
      dev-internal:
        aliases:
          - app-one-db
    labels:
      - "contrail.workspace.name=dev"
      - "contrail.app.name=app-one"
    environment:
      - CONTRAIL_APP_ONE_DB_HOST=app-one-db
      - CONTRAIL_APP_ONE_DB_PORT=5432

networks:
  dev-internal:
    external: true
  contrail-proxy:
    external: true
```

---

## Manual Override File

**Location**: `{workspace}/overrides/{application-name}.yaml`

Optional. If present, merged after the generated override file. Useful for workspace-specific customizations that can't be expressed in the application config.

### Preservation Guarantee

Files in `{workspace}/overrides/` are **never modified by Contrail**. They persist across all regeneration operations, including `workspace generate --force`. Only the `.generated/` directory contents are affected by regeneration.

### Merge Order

Docker Compose files are merged in this order:

```
docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
```

This allows workspace-specific customizations (extra environment variables, volume mounts, middleware) that persist across regeneration.

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

## Staleness Detection

Override files are considered stale if any source file has a newer modification time:

- `workspace.yaml`
- `{app}/application.yaml`
- `.generated/state.yaml` (active flavor may have changed)
- Active flavor's compose files

### Behavior

- `workspace up` and `workspace generate` automatically regenerate stale overrides
- Use `--force` flag to regenerate regardless of staleness
- Use `--no-generate` to skip regeneration entirely

---

## Edge Cases

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
- Warning displayed about restart requirement

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Missing compose file | `Flavor "{name}" references non-existent file: {file}` | Fix compose_files list |
| Invalid service reference | `Exported service references non-existent Compose service` | Fix service reference |
| Write permission denied | `Cannot write to .generated/: permission denied` | Fix directory permissions |
| YAML syntax error | `Invalid YAML in {file}: {error}` | Fix YAML syntax |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
