# Naming Conventions Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail uses convention-based naming to derive hostnames, aliases, project names, and other identifiers from workspace and application names.

---

## Naming Patterns

| Element | Pattern | Example |
|---------|---------|---------|
| Public hostname | `{workspace}-{app}-{export}.{domain}` | `dev-app-one-web.contrail.test` |
| Internal alias | `{app}-{export}` | `app-one-web` |
| Network name | `{workspace}-internal` | `dev-internal` |
| Project name | `{workspace}-{app}` | `dev-app-one` |
| Traefik router | `{workspace}-{app}-{export}-{protocol}` | `dev-app-one-web-https` |

---

## Name Requirements

### Workspace Names

- Lowercase alphanumeric with hyphens
- Must be unique across all registered workspaces
- Examples: `dev`, `feature-x`, `pr-123`

### Application Names

- Lowercase alphanumeric with hyphens
- Inferred from directory name
- Examples: `app-one`, `frontend`, `api-gateway`

### Exported Service Names

- Key in `exported_services` map
- May differ from underlying Compose service name
- Used in hostname and alias generation

---

## Environment Variable Naming

Variables use `CONTRAIL_` prefix with underscore conversion:

- Hyphens → underscores
- Names → SCREAMING_SNAKE_CASE

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}
```

Example: `app-one` → `APP_ONE`

```bash
CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_PORT=443
CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test
```

---

## Collision Warning

Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns. Creative naming could produce conflicts:

**Example collision**:
- Workspace `dev-app` with app `one` → project `dev-app-one`
- Workspace `dev` with app `app-one` → project `dev-app-one`

**Potential impacts**:
- Traefik routers: Conflicting names cause routing failures
- Docker volumes: Conflicting names could cause data sharing/overwriting
- Docker networks: Conflicting names could connect unrelated services

**Prevention**: Follow lowercase-alphanumeric-with-hyphens convention and avoid ambiguous concatenations.

---

## Related Documentation

- [ADR-0004: Convention-Based Naming](../../decisions/0004-convention-based-naming/README.md)
- [Configuration Schemas Spec](../configuration-schemas/README.md)
- [Environment Variables Spec](../environment-variables/README.md)
