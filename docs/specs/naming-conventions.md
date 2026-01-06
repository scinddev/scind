# Naming Conventions Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail uses convention-based naming for predictability and discoverability. All naming patterns are deterministic, allowing users and tools to compute hostnames, aliases, and project names without querying Contrail.

**Related Documents**:
- [ADR-0001: Docker Compose Project Name Isolation](../decisions/0001-docker-compose-project-name-isolation.md)
- [ADR-0004: Convention-Based Naming](../decisions/0004-convention-based-naming.md)
- [Configuration Schemas](./configuration-schemas.md)
- [Environment Variables](./environment-variables.md)

---

## Behavior

### Naming Pattern Resolution

Names are resolved at generation time using template patterns. Default patterns are built-in but can be customized in `workspace.yaml`.

### Default Templates

```yaml
workspace:
  templates:
    hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
    alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
    project-name: "%WORKSPACE_NAME%-%APPLICATION_NAME%"
```

---

## Data Schema

### Naming Patterns

| Entity | Pattern | Example |
|--------|---------|---------|
| **Workspace names** | Lowercase alphanumeric with hyphens | `dev`, `feature-x`, `pr-123` |
| **Application names** | Lowercase alphanumeric with hyphens, inferred from directory name | `app-one`, `my-api` |
| **Exported service names** | The key in `exported_services`, may differ from the underlying Compose service name | `web`, `db`, `api` |
| **Proxied hostnames** | `{workspace}-{application}-{exported_service}.{domain}` | `dev-app-one-web.contrail.test` |
| **Internal aliases** | `{application}-{exported_service}` | `app-one-web`, `app-one-db` |
| **Environment variables** | `CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}` in SCREAMING_SNAKE_CASE | `CONTRAIL_APP_ONE_WEB_URL` |
| **Traefik router names** | `{workspace}-{application}-{exported_service}-{protocol}` | `dev-app-one-web-https` |
| **Docker Compose project names** | `{workspace}-{application}` | `dev-app-one` |

### Project Names

Docker Compose project names uniquely identify an application instance.

| Attribute | Value |
|-----------|-------|
| Pattern | `{workspace}-{app}` |
| Example | `dev-frontend` |
| Characters | Lowercase alphanumeric, hyphens |
| Max Length | 63 characters (Docker limit) |

### Hostnames (Public)

Fully qualified domain names for proxied services.

| Attribute | Value |
|-----------|-------|
| Pattern | `{workspace}-{app}-{exported_service}.{domain}` |
| Example | `dev-frontend-web.contrail.test` |
| Characters | Lowercase alphanumeric, hyphens, dots |
| Max Length | 253 characters (DNS limit) |

### Network Aliases (Internal)

Short names for inter-application communication within a workspace.

| Attribute | Value |
|-----------|-------|
| Pattern | `{app}-{exported_service}` |
| Example | `frontend-web` |
| Characters | Lowercase alphanumeric, hyphens |

### Network Names

| Network | Pattern | Example |
|---------|---------|---------|
| Proxy | `contrail-proxy` | `contrail-proxy` |
| Workspace Internal | `{workspace}-internal` | `dev-internal` |
| Application Default | Managed by Docker Compose | `dev-frontend_default` |

---

## Application Requirements

For an application to work well within a workspace, it should:

1. **Include an `application.yaml`**: This file defines the service contract - which services the application exports to the workspace. This is owned and maintained by the application developers.
2. **Use environment variables for external service URLs**: Don't hardcode hostnames for dependencies. Use the injected `CONTRAIL_{APP}_{EXPORTED_SERVICE}_*` variables.
3. **Expose ports without host bindings**: Use `ports: ["8080"]` not `ports: ["8080:8080"]` to avoid conflicts.
4. **Use relative volume paths**: Ensure builds and mounts work regardless of absolute path.

---

## The Service Contract (`application.yaml`)

The `application.yaml` file is the interface between the application and the workspace system:

- **Owned by**: Application developers (committed to the application's git repository)
- **Application name**: Inferred from the directory name (no explicit `name:` field required)
- **Purpose**:
  - Declares which services the application exports and how they should be exposed
  - Defines flavors (different ways to run the application with different compose file combinations)
- **Consumed by**: Workspace tooling (to generate override files)

Application developers should update `application.yaml` when:
- Adding a new service that other applications need to access
- Changing the type (proxied, assigned) or protocol (http, https) of an exported service
- Renaming services that are exposed to the workspace
- Adding a new flavor (e.g., a "lite" mode that excludes optional services)
- Changing which compose files are needed for a flavor

---

## Examples

### Example 1: Multi-Application Workspace

**Configuration**:
- Workspace: `dev`
- Domain: `contrail.test`
- Applications: `frontend`, `backend`, `api`

**Generated Names**:

| Resource | frontend/web | backend/api | api/graphql |
|----------|-------------|-------------|-------------|
| Project | `dev-frontend` | `dev-backend` | `dev-api` |
| Hostname | `dev-frontend-web.contrail.test` | `dev-backend-api.contrail.test` | `dev-api-graphql.contrail.test` |
| Alias | `frontend-web` | `backend-api` | `api-graphql` |
| Router | `dev-frontend-web-https` | `dev-backend-api-https` | `dev-api-graphql-https` |
| Env Var | `CONTRAIL_FRONTEND_WEB_*` | `CONTRAIL_BACKEND_API_*` | `CONTRAIL_API_GRAPHQL_*` |

### Example 2: Custom Templates

**Configuration**:
```yaml
workspace:
  name: staging
  templates:
    hostname: "%APPLICATION_NAME%-%EXPORTED_SERVICE%.staging.example.com"
    alias: "%APPLICATION_NAME%.%EXPORTED_SERVICE%"
```

**Result** for `frontend/web`:
- Hostname: `frontend-web.staging.example.com`
- Alias: `frontend.web`

---

## Validation Rules

### Workspace Names

- Lowercase alphanumeric with hyphens
- Must start with a letter
- Max 50 characters
- Pattern: `^[a-z][a-z0-9-]*$`

### Application Names

- Lowercase alphanumeric with hyphens
- Must start with a letter
- Max 50 characters (to allow room in combined names)
- Pattern: `^[a-z][a-z0-9-]*$`

### Exported Service Names

- Lowercase alphanumeric with hyphens
- Must start with a letter
- Max 30 characters
- Pattern: `^[a-z][a-z0-9-]*$`

### Domain Names

- Valid DNS domain
- Must end with a TLD
- Recommended: Use `.test` TLD (RFC 2606 reserved)

---

## Edge Cases

### Collision Warning

Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Creative naming that produces identical project names could cause conflicts:

- **Traefik routers**: Conflicting router names cause routing failures
- **Docker volumes**: Conflicting volume names (e.g., `dev-app-one_postgres_data`) could cause data to be shared unexpectedly or overwritten
- **Docker networks**: Conflicting network names could connect unrelated services

### Name Collision Example

Workspace `dev-app` with app `one` and workspace `dev` with app `app-one` both produce project name `dev-app-one`.

**Behavior**: Error during workspace initialization:
```
Error: Project name collision detected.
  Workspace "dev" app "app-one" produces project name "dev-app-one"
  This conflicts with workspace "dev-app" app "one".

  Choose a different workspace or application name.
```

### Long Names

**Scenario**: Combined name exceeds Docker's 63-character limit.

**Behavior**: Error during generation:
```
Error: Project name too long.
  Generated name: "very-long-workspace-name-extremely-long-application-name" (58 chars)
  Maximum allowed: 63 characters

  Shorten workspace or application name.
```

### Best Practice

Follow the lowercase-alphanumeric-with-hyphens convention and avoid names that could produce ambiguous concatenations.

---

## Template Variables Reference

| Variable | Scope | Description | Example |
|----------|-------|-------------|---------|
| `%PROXY_DOMAIN%` | Proxy | Domain from `proxy.yaml` | `contrail.test` |
| `%WORKSPACE_NAME%` | Workspace | Workspace name | `dev` |
| `%WORKSPACE_NETWORK%` | Workspace | Internal network name | `dev-internal` |
| `%APPLICATION_NAME%` | Application | Application identifier | `app-one` |
| `%APPLICATION_FLAVOR%` | Application | Resolved flavor | `default` |
| `%EXPORTED_SERVICE%` | Export | Key from `exported_services` | `web-debug` |
| `%SERVICE_NAME%` | Export | Underlying Compose service | `web` |
| `%SERVICE_PORT%` | Export | Container port number | `8080` |
| `%SERVICE_PROTOCOL%` | Export | Protocol (for proxied) | `https` |

---

## Git Strategy

**Workspace repository** (optional - can be version controlled):
- `workspace.yaml` - workspace definition
- `overrides/` - manual overrides

**Generated files** (gitignored):
- `.generated/` - generated override files

**Application directories** (cloned repositories):
- Each application is its own git repository
- `application.yaml` lives in the application repo and defines its service contract
- Can be managed as submodules, or simply cloned separately

### Example Workspace `.gitignore`

```
.generated/
app-*/
```

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Invalid name characters | `Name contains invalid characters` | Use allowed characters only |
| Name too long | `Name exceeds maximum length` | Shorten name |
| Name collision | `Name collision detected` | Choose unique names |
| Missing required variable | `Template variable not available in this context` | Use valid variables for scope |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
