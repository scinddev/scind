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

### Project Names

Docker Compose project names uniquely identify an application instance.

| Pattern | `{workspace}-{app}` |
|---------|---------------------|
| Example | `dev-frontend` |
| Characters | Lowercase alphanumeric, hyphens |
| Max Length | 63 characters (Docker limit) |

**Purpose**: Prevents container name collisions when running multiple workspaces or applications.

### Hostnames (Public)

Fully qualified domain names for proxied services.

| Pattern | `{workspace}-{app}-{exported_service}.{domain}` |
|---------|------------------------------------------------|
| Example | `dev-frontend-web.contrail.test` |
| Characters | Lowercase alphanumeric, hyphens, dots |
| Max Length | 253 characters (DNS limit) |

**Purpose**: Routes external HTTP/HTTPS requests through Traefik to the correct container.

### Network Aliases (Internal)

Short names for inter-application communication within a workspace.

| Pattern | `{app}-{exported_service}` |
|---------|---------------------------|
| Example | `frontend-web` |
| Characters | Lowercase alphanumeric, hyphens |

**Purpose**: Provides stable, workspace-agnostic names for service discovery. Applications can reference `frontend-web` regardless of which workspace they're in.

### Network Names

| Network | Pattern | Example |
|---------|---------|---------|
| Proxy | `contrail-proxy` | `contrail-proxy` |
| Workspace Internal | `{workspace}-internal` | `dev-internal` |
| Application Default | Managed by Docker Compose | `dev-frontend_default` |

### Environment Variables

| Pattern | `CONTRAIL_{APP}_{EXPORTED_SERVICE}_{SUFFIX}` |
|---------|---------------------------------------------|
| Example | `CONTRAIL_FRONTEND_WEB_URL` |
| Transformation | Uppercase, hyphens to underscores |

### Traefik Router Names

| Pattern | `{workspace}-{app}-{exported_service}-{protocol}` |
|---------|--------------------------------------------------|
| Example | `dev-frontend-web-https` |

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

### Example 2: Single-Application Workspace

**Configuration**:
- Workspace: `dev`
- Application: `myapp` (at workspace root with `path: .`)
- Exported service: `web`

**Generated Names**:
- Project: `dev-myapp`
- Hostname: `dev-myapp-web.contrail.test`
- Alias: `myapp-web`
- Network: `dev-internal`

### Example 3: Custom Templates

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

### Name Collisions

**Scenario**: Two different naming combinations produce the same project name.

**Example**:
- Workspace `dev-app` with app `one` → `dev-app-one`
- Workspace `dev` with app `app-one` → `dev-app-one`

**Behavior**: Error during workspace initialization:
```
Error: Project name collision detected.
  Workspace "dev" app "app-one" produces project name "dev-app-one"
  This conflicts with workspace "dev-app" app "one".

  Choose a different workspace or application name.
```

**Collision Warning**: Docker Compose project names, Traefik router names, volume names, and network names are all derived from naming patterns. Collisions could cause:
- Traefik routers: Routing failures
- Docker volumes: Unexpected data sharing or overwriting
- Docker networks: Unrelated services connected

### Long Names

**Scenario**: Combined name exceeds Docker's 63-character limit.

**Behavior**: Error during generation:
```
Error: Project name too long.
  Generated name: "very-long-workspace-name-extremely-long-application-name" (58 chars)
  Maximum allowed: 63 characters

  Shorten workspace or application name.
```

### Special Characters

**Scenario**: User attempts to use underscores or other special characters.

**Behavior**: Validation error on configuration:
```
Error: Invalid workspace name: "my_workspace"
  Workspace names must be lowercase alphanumeric with hyphens only.
  Pattern: ^[a-z][a-z0-9-]*$
```

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

**Note on `%SERVICE_PORT%`**: This variable provides the container's internal port number. While not used in the default templates, it's available for advanced customization such as adding debugging labels.

---

## Error Handling

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| Invalid name characters | VALIDATION | `Name contains invalid characters` | Use allowed characters only |
| Name too long | VALIDATION | `Name exceeds maximum length` | Shorten name |
| Name collision | VALIDATION | `Name collision detected` | Choose unique names |
| Missing required variable | TEMPLATE | `Template variable not available in this context` | Use valid variables for scope |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
