# Environment Variables Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

All exported services receive environment variables for service discovery, enabling applications to reference other services without hardcoding hostnames.

---

## Naming Convention

Variables use `CONTRAIL_` prefix with underscore conversion:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}
```

**Name transformation**: Hyphens → underscores, UPPERCASE

Example: `app-one` → `APP_ONE`, `web-debug` → `WEB_DEBUG`

---

## Variable Types

### Base Variables

Generated for each exported service:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # proxied only
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # proxied only
```

### Protocol-Specific Variables

Generated for each proxied protocol:

```
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

---

## Variables by Port Type

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | Yes | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | Yes | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | Yes | Both |
| `assigned` | - | Internal alias | Assigned port | No | No | No |

**HTTPS-default rationale**: When both HTTP and HTTPS are configured, base variables default to HTTPS (port 443) following security-by-default principles.

---

## Examples

### Proxied Service with HTTPS

```bash
# Base variables
CONTRAIL_APP_ONE_WEB_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_PORT=443
CONTRAIL_APP_ONE_WEB_SCHEME=https
CONTRAIL_APP_ONE_WEB_URL=https://dev-app-one-web.contrail.test

# Protocol-specific
CONTRAIL_APP_ONE_WEB_HTTPS_HOST=dev-app-one-web.contrail.test
CONTRAIL_APP_ONE_WEB_HTTPS_PORT=443
CONTRAIL_APP_ONE_WEB_HTTPS_URL=https://dev-app-one-web.contrail.test
```

### Assigned Port Service

```bash
CONTRAIL_APP_ONE_DB_HOST=app-one-db
CONTRAIL_APP_ONE_DB_PORT=5432
```

---

## Usage Examples

### PHP

```php
// Using URL directly for proxied services
$apiUrl = getenv('CONTRAIL_APP_TWO_API_URL') ?: 'https://app-two-api.contrail.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// Building connection for assigned port services
$dbHost = getenv('CONTRAIL_APP_ONE_DB_HOST') ?: 'app-one-db';
$dbPort = getenv('CONTRAIL_APP_ONE_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

### Node.js

```javascript
// Using URL directly
const apiUrl = process.env.CONTRAIL_APP_TWO_API_URL || 'https://app-two-api.contrail.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Building connection manually
const dbHost = process.env.CONTRAIL_APP_ONE_DB_HOST || 'app-one-db';
const dbPort = process.env.CONTRAIL_APP_ONE_DB_PORT || '5432';
```

---

## Related Documentation

- [Naming Conventions Spec](../naming-conventions/README.md)
- [Configuration Schemas Spec](../configuration-schemas/README.md)
- [ADR-0007: Port Type System](../../decisions/0007-port-type-system/README.md)
