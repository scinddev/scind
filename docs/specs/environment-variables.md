# Environment Variables Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail injects environment variables into containers for service discovery. This enables applications to reference other services within the workspace without hardcoding hostnames. Environment variables use the `CONTRAIL_` prefix to avoid conflicts with application-defined variables.

**Related Documents**:
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [Port Types](./port-types.md)
- [Naming Conventions](./naming-conventions.md)
- [Generated Override Files](./generated-override-files.md)

---

## Behavior

### Variable Injection

Environment variables are generated during `workspace generate` and written to the generated override files. They are injected into containers when they start via Docker Compose's environment merge behavior.

### Name Transformation

- Hyphens in application and exported service names are converted to underscores
- All names are uppercased
- Examples:
  - `app-one` becomes `APP_ONE`
  - `web-debug` becomes `WEB_DEBUG`

---

## Data Schema

### Naming Convention

All Contrail environment variables use the `CONTRAIL_` prefix.

#### Base Variables

Generated for each exported service:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

#### Protocol-Specific Variables

Generated for each proxied protocol:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

#### Workspace Context Variable

```
CONTRAIL_WORKSPACE_NAME={workspace_name}
```

### Field Reference

| Variable Pattern | Type | Port Type | Description |
|------------------|------|-----------|-------------|
| `CONTRAIL_WORKSPACE_NAME` | string | All | Current workspace name |
| `*_HOST` | string | All | Hostname or alias for the service |
| `*_PORT` | integer | All | Port number (proxy port or assigned port) |
| `*_SCHEME` | string | Proxied only | `http` or `https` |
| `*_URL` | string | Proxied only | Complete URL |
| `*_{PROTOCOL}_HOST` | string | Proxied only | Protocol-specific hostname |
| `*_{PROTOCOL}_PORT` | integer | Proxied only | Protocol-specific port |
| `*_{PROTOCOL}_URL` | string | Proxied only | Protocol-specific URL |

### Variable Generation Rules

#### For `proxied` Type Ports

| Variable | Value |
|----------|-------|
| `*_HOST` | Fully qualified proxied hostname (e.g., `dev-app-one-web.contrail.test`) |
| `*_PORT` | Proxy port (443 for HTTPS, 80 for HTTP) — **not** the container port |
| `*_SCHEME` | `https` or `http` |
| `*_URL` | Complete URL (e.g., `https://dev-app-one-web.contrail.test`) |
| `*_{PROTOCOL}_HOST` | Same as `*_HOST` |
| `*_{PROTOCOL}_PORT` | Protocol-specific port (443 for HTTPS, 80 for HTTP) |
| `*_{PROTOCOL}_URL` | Protocol-specific URL |

#### For `assigned` Type Ports

| Variable | Value |
|----------|-------|
| `*_HOST` | Internal network alias (e.g., `app-one-db`) |
| `*_PORT` | Assigned host port (may differ from requested port) |
| `*_SCHEME` | Not generated |
| `*_URL` | Not generated |
| Protocol-specific | Not generated |

### Summary Table

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | Yes | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | Yes | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | Yes | Both |
| `assigned` | - | Internal alias | Assigned port | No | No | No |

### HTTPS-Default Rationale

When both HTTP and HTTPS are configured for an exported service, base variables (`*_PORT`, `*_SCHEME`, `*_URL`) default to HTTPS (port 443) following security-by-default principles. Applications should prefer HTTPS for service-to-service communication. Use protocol-specific variables (`*_HTTP_*`) when HTTP is explicitly required.

---

## Examples

### Example 1: Proxied Web Service (HTTPS only)

**Configuration**:
```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
```

**Generated Variables**:
```bash
CONTRAIL_WORKSPACE_NAME=dev
CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_PORT=443
CONTRAIL_APP_ONE_WEB_SCHEME=https
CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTPS_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTPS_PORT=443
CONTRAIL_APP_ONE_WEB_HTTPS_URL=https://dev-app-one-web.contrail.test
```

### Example 2: Proxied Web Service (HTTP and HTTPS)

**Configuration**:
```yaml
exported_services:
  web:
    ports:
      - type: proxied
        protocol: https
        visibility: public
      - type: proxied
        protocol: http
        visibility: protected
```

**Generated Variables**:
```bash
CONTRAIL_WORKSPACE_NAME=dev
# Base variables (default to HTTPS)
CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_PORT=443
CONTRAIL_APP_ONE_WEB_SCHEME=https
CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test
# HTTPS-specific
CONTRAIL_APP_ONE_WEB_HTTPS_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTPS_PORT=443
CONTRAIL_APP_ONE_WEB_HTTPS_URL=https://dev-app-one-web.contrail.test
# HTTP-specific
CONTRAIL_APP_ONE_WEB_HTTP_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTP_PORT=80
CONTRAIL_APP_ONE_WEB_HTTP_URL=http://dev-app-one-web.contrail.test
```

### Example 3: Assigned Database Port

**Configuration**:
```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: assigned
        port: 5432
        visibility: protected
```

**Generated Variables** (port 5432 was available):
```bash
CONTRAIL_WORKSPACE_NAME=dev
CONTRAIL_APP_ONE_DB_HOST=app-one-db
CONTRAIL_APP_ONE_DB_PORT=5432
```

**Generated Variables** (port 5432 was taken, assigned 5433):
```bash
CONTRAIL_WORKSPACE_NAME=dev
CONTRAIL_APP_ONE_DB_HOST=app-one-db
CONTRAIL_APP_ONE_DB_PORT=5433
```

---

## Usage in Applications

### PHP Example

```php
// Using URL directly for proxied services
$apiUrl = getenv('CONTRAIL_APP_TWO_API_URL') ?: 'https://app-two-api.contrail.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// Building connection for assigned port services
$dbHost = getenv('CONTRAIL_APP_ONE_DB_HOST') ?: 'app-one-db';
$dbPort = getenv('CONTRAIL_APP_ONE_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

### Node.js Example

```javascript
// Using URL directly
const apiUrl = process.env.CONTRAIL_APP_TWO_API_URL || 'https://app-two-api.contrail.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Building connection manually
const dbHost = process.env.CONTRAIL_APP_ONE_DB_HOST || 'app-one-db';
const dbPort = process.env.CONTRAIL_APP_ONE_DB_PORT || '5432';
```

### Python Example

```python
import os

# Using URL directly
api_url = os.environ.get('CONTRAIL_APP_TWO_API_URL', 'https://app-two-api.contrail.test')
response = requests.get(f'{api_url}/endpoint')

# Building connection
db_host = os.environ.get('CONTRAIL_APP_ONE_DB_HOST', 'app-one-db')
db_port = os.environ.get('CONTRAIL_APP_ONE_DB_PORT', '5432')
```

---

## Edge Cases

### Hyphenated Names

**Scenario**: Application name contains hyphens (e.g., `my-api-service`).

**Behavior**: Hyphens are converted to underscores.

```bash
CONTRAIL_MY_API_SERVICE_WEB_URL=https://dev-my-api-service-web.contrail.test
```

### Override Behavior

**Scenario**: Application's compose file already defines an environment variable with the same name.

**Behavior**: The application's explicitly defined environment variable takes precedence. Contrail's generated variables are added but can be overridden.

### Multiple Assigned Ports

**Scenario**: Exported service has multiple assigned ports.

**Behavior**: Only the first assigned port generates the base variables. Each port is documented in labels but not in environment variables.

**Rationale**: Environment variable names would conflict. Use labels for discovering additional ports.

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| Name collision | `Environment variable name collision after transformation` | Rename exported service |
| Invalid characters | `Exported service name contains invalid characters` | Use alphanumeric and hyphens only |

---

## Validation Rules

- Application names must be lowercase alphanumeric with hyphens
- Exported service names must be lowercase alphanumeric with hyphens
- Generated variable names must be valid shell variable names (uppercase alphanumeric with underscores)

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from technical spec |
