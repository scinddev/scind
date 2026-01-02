# Generated Override Files Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail generates Docker Compose override files that add workspace integration to applications without modifying their original compose files.

---

## Location

**Generated overrides**: `{workspace}/.generated/{application-name}.override.yaml`

**Manual overrides**: `{workspace}/overrides/{application-name}.yaml`

---

## Merge Order

Docker Compose files are merged in this order:

```
docker compose -f base.yaml -f .generated/app.override.yaml -f overrides/app.yaml
```

This allows:
1. Application's base compose file(s) define core services
2. Generated override adds workspace integration
3. Manual override provides workspace-specific customizations

---

## Generated Override Structure

```yaml
# AUTO-GENERATED - Do not edit directly
# Source: workspace.yaml + app-one/application.yaml
# Flavor: default
# Generated: 2024-12-27T10:30:00Z

name: dev-app-one                        # Explicit project name

services:
  web:
    networks:
      dev-internal:
        aliases:
          - app-one-web
      contrail-proxy: {}
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
      - "5432:5432"                       # Assigned port mapping
    networks:
      dev-internal:
        aliases:
          - app-one-db
    labels:
      # Contrail context labels
      - "contrail.workspace.name=dev"
      - "contrail.app.name=app-one"
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

## What Gets Added

### Networks

- **Workspace internal network** (`{workspace}-internal`): Added to all exported services for internal communication
- **Proxy network** (`contrail-proxy`): Added to proxied services for Traefik routing

### Network Aliases

Each exported service gets an alias on the internal network:
- Pattern: `{app}-{export}` (e.g., `app-one-web`, `app-one-db`)

### Labels

- **Traefik routing labels**: For proxied services
- **Contrail context labels**: Workspace and app identification
- **Contrail export labels**: Export metadata for external tools

### Environment Variables

All `CONTRAIL_*` environment variables for service discovery.

### Port Mappings

For `assigned` type ports, host:container port mappings are added.

---

## Manual Override File

**Location**: `{workspace}/overrides/{application-name}.yaml`

**Preservation guarantee**: Never modified by Contrail; persists across all regeneration.

**Example**:
```yaml
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

## Related Documentation

- [ADR-0003: Pure Overlay Design](../../decisions/0003-pure-overlay-design/README.md)
- [Configuration Schemas Spec](../configuration-schemas/README.md)
- [Docker Labels Spec](../docker-labels/README.md)
