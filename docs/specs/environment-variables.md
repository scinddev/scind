# Specification: Environment Variable Injection

**Version**: 0.5.0
**Date**: December 2024

<!-- Migrated from specs/contrail-technical-spec.md -->

---

## Overview

All exported services receive environment variables for service discovery. This enables applications to reference other services without hardcoding hostnames.

---

## Naming Convention

Environment variables use a `CONTRAIL_` prefix to avoid conflicts with application-defined variables.

**Name transformation**: Hyphens in application and exported service names are converted to underscores, and names are uppercased (e.g., `app-one` becomes `APP_ONE`, `web-debug` becomes `WEB_DEBUG`).

---

## Base Variables

Always generated for each exported service:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

---

## Protocol-Specific Variables

Generated for each proxied protocol:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

---

## Variable Generation Rules

### For `proxied` Type Ports

- `*_HOST` contains the fully qualified proxied hostname (e.g., `dev-app-one-web.contrail.test`)
- `*_PORT` contains the proxy port (443 for HTTPS, 80 for HTTP)—**not** the container port
- `*_SCHEME` and `*_URL` are generated
- Protocol-specific variables (`*_HTTPS_*`, `*_HTTP_*`) are also generated

### For `assigned` Type Ports

- `*_HOST` contains the internal alias (e.g., `app-one-db`)
- `*_PORT` contains the assigned host port (which may differ from the requested port)
- No `*_SCHEME` or `*_URL` variables
- No protocol-specific variables

---

## Summary Table

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | Yes | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | Yes | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | Yes | Both |
| `assigned` | - | Internal alias | Assigned port | No | No | No |

**HTTPS-default rationale**: When both HTTP and HTTPS are configured for an exported service, base variables (`*_PORT`, `*_SCHEME`, `*_URL`) default to HTTPS (port 443) following security-by-default principles.

---

## Usage Examples

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

---

## Related Documents

- [Port Types and Proxying](port-types.md)
- [Generated Override Files](generated-override-files.md)
