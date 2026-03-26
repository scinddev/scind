## Environment Variable Injection

All exported services receive environment variables for service discovery. This enables applications to reference other services without hardcoding hostnames.

### Naming Convention

Environment variables use a `SCIND_` prefix to avoid conflicts with application-defined variables.

**Name transformation**: Hyphens in application and exported service names are converted to underscores, and names are uppercased (e.g., `shared-db` becomes `SHARED_DB`, `web-debug` becomes `WEB_DEBUG`).

**Base variables** (always generated for each exported service):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_HOST={hostname_or_alias}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_SCHEME={scheme}    # Only for proxied types
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_URL={url}          # Only for proxied types
```

**Protocol-specific variables** (generated for each proxied protocol):
```
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_HOST={hostname}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_PORT={port}
SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{PROTOCOL}_URL={url}
```

**Apex variables** (generated only for proxied primary exports):
```
SCIND_{APPLICATION}_APEX_HOST={apex_hostname}
SCIND_{APPLICATION}_APEX_PORT={port}
SCIND_{APPLICATION}_APEX_SCHEME={scheme}
SCIND_{APPLICATION}_APEX_URL={url}
```

Apex environment variables follow the pattern `SCIND_{APPLICATION}_APEX_{SUFFIX}` — the exported service name segment is omitted. These are only injected when the application has a primary export with proxied ports. Assigned-port primary exports do not generate apex environment variables (they have no hostname). See [ADR-0013](../decisions/0013-apex-url-primary-designation.md) for primary designation rules.

### Variable Generation Rules

**For `proxied` type ports**:
- `*_HOST` contains the fully qualified proxied hostname (e.g., `dev-frontend-web.scind.test`)
- `*_PORT` contains the proxy port (443 for HTTPS, 80 for HTTP)—**not** the container port
- `*_SCHEME` and `*_URL` are generated
- Protocol-specific variables (`*_HTTPS_*`, `*_HTTP_*`) are also generated

**For `assigned` type ports**:
- `*_HOST` contains the internal alias (e.g., `shared-db-db`)
- `*_PORT` contains the assigned host port (which may differ from the requested port)
- No `*_SCHEME` or `*_URL` variables
- No protocol-specific variables

| Type | Protocol | `*_HOST` | `*_PORT` | `*_SCHEME` | `*_URL` | Protocol Vars |
|------|----------|----------|----------|------------|---------|---------------|
| `proxied` | `https` | Proxied hostname | 443 | `https` | ✓ | `*_HTTPS_*` |
| `proxied` | `http` | Proxied hostname | 80 | `http` | ✓ | `*_HTTP_*` |
| `proxied` | both | Proxied hostname | 443 | `https` | ✓ | Both |
| `assigned` | - | Internal alias | Assigned port | ✗ | ✗ | ✗ |
| `proxied` (primary) | `https` | Apex hostname | 443 | `https` | ✓ | `*_APEX_*` |
| `assigned` (primary) | - | *(no apex env vars)* | - | - | - | *(alias only)* |

**HTTPS-default rationale**: When both HTTP and HTTPS are configured for an exported service, base variables (`*_PORT`, `*_SCHEME`, `*_URL`) default to HTTPS (port 443) following security-by-default principles. Applications should prefer HTTPS for service-to-service communication. Use protocol-specific variables (`*_HTTP_*`) when HTTP is explicitly required.

### Usage in Applications

```php
// PHP example - using URL directly for proxied services
$apiUrl = getenv('SCIND_BACKEND_API_URL') ?: 'https://backend-api.scind.test';
$response = $httpClient->get("{$apiUrl}/endpoint");

// PHP example - building connection for assigned port services
$dbHost = getenv('SCIND_SHARED_DB_DB_HOST') ?: 'shared-db-db';
$dbPort = getenv('SCIND_SHARED_DB_DB_PORT') ?: '5432';
$dsn = "pgsql:host={$dbHost};port={$dbPort};dbname=app";
```

```javascript
// Node.js example - using URL directly
const apiUrl = process.env.SCIND_BACKEND_API_URL || 'https://backend-api.scind.test';
const response = await fetch(`${apiUrl}/endpoint`);

// Node.js example - building connection manually
const dbHost = process.env.SCIND_SHARED_DB_DB_HOST || 'shared-db-db';
const dbPort = process.env.SCIND_SHARED_DB_DB_PORT || '5432';
```

---

## Configuration Environment Variables

These environment variables configure Scind itself (distinct from service discovery variables injected into containers).

| Variable | Description | Default |
|----------|-------------|---------|
| `TRAEFIK_IMAGE` | Traefik Docker image to use | `traefik:v3.2.3` |
| `SCIND_CONFIG_DIR` | Configuration directory | `~/.config/scind` |
| `SCIND_STATE_DIR` | State file directory | `~/.config/scind` |

These can also be set in `proxy.yaml` configuration.

---

## Related Documents

- [Configuration Schemas](configuration-schemas.md) - Configuration file formats
- [Generated Override Files](generated-override-files.md) - How variables are injected
- [ADR-0013: Apex URL Primary Designation](../decisions/0013-apex-url-primary-designation.md) - Primary designation design
