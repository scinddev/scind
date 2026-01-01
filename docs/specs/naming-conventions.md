# Specification: Naming Conventions

**Version**: 0.1.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

Contrail uses consistent naming conventions to generate predictable hostnames, aliases, project names, and environment variables. Following these conventions prevents collisions and ensures discoverability.

---

## Naming Patterns

| Entity | Pattern | Example |
|--------|---------|---------|
| Workspace name | Lowercase alphanumeric with hyphens | `dev`, `feature-x`, `pr-123` |
| Application name | Lowercase alphanumeric with hyphens (inferred from directory) | `app-one`, `frontend` |
| Exported service name | Key in `exported_services` | `web`, `api`, `db` |
| Docker Compose project | `{workspace}-{application}` | `dev-app-one` |
| Proxied hostname | `{workspace}-{application}-{export}.{domain}` | `dev-app-one-web.contrail.test` |
| Internal alias | `{application}-{export}` | `app-one-web`, `app-one-db` |
| Traefik router | `{workspace}-{application}-{export}-{protocol}` | `dev-app-one-web-https` |
| Environment variable | `CONTRAIL_{APP}_{EXPORT}_{SUFFIX}` | `CONTRAIL_APP_ONE_WEB_URL` |

---

## Environment Variable Name Transformation

Hyphens in application and exported service names are converted to underscores, and names are uppercased:

| Original | Transformed |
|----------|-------------|
| `app-one` | `APP_ONE` |
| `web-debug` | `WEB_DEBUG` |
| `my-app` | `MY_APP` |

---

## Collision Warnings

Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Creative naming that produces identical concatenations can cause serious conflicts.

### Collision Example

Consider two workspaces:
- Workspace `dev-app` with application `one`
- Workspace `dev` with application `app-one`

Both produce project name `dev-app-one`, causing:

| Resource | Collision Impact |
|----------|-----------------|
| **Docker Compose project** | Containers from different workspaces share the same project; `down` removes wrong containers |
| **Traefik routers** | Conflicting router names cause routing failures; requests go to wrong service |
| **Docker volumes** | Volume `dev-app-one_postgres_data` shared unexpectedly; data corruption or leakage |
| **Docker networks** | Conflicting network names connect unrelated services |

### Prevention

1. **Use clear namespace boundaries**: Avoid names that could produce ambiguous concatenations
2. **Avoid hyphens that mimic separators**: `dev-app` + `one` = `dev` + `app-one` = `dev-app-one`
3. **Prefer simple, distinct workspace names**: `dev`, `review`, `staging` rather than `dev-feature-x`

### Detection

Contrail does **not** currently detect potential collisions at registration time. This is a known limitation. Users should follow the conventions above to avoid issues.

---

## Template Customization

Templates can be customized at the workspace level to change the generated patterns:

```yaml
workspace:
  name: dev
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

The above shows the default templates. Custom templates use the same `%VARIABLE_NAME%` syntax.

**Note**: Template customization allows flexibility but increases collision risk. If you change the templates, ensure the resulting patterns remain unique across your workspaces.

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md)
- [ADR-0004: Convention-Based Naming](../decisions/0004-convention-based-naming.md)
