## Proxy Layer

### Architecture Overview

Traefik serves as the reverse proxy, routing external requests to workspace services.

```
[External Request] → [Traefik:443/80] → [scind-proxy network] → [Service Container]
```

#### Components

- **Traefik container**: Single instance managing all workspace routing
- **scind-proxy network**: Host-level Docker network connecting Traefik to services
- **Dynamic configuration**: Label-based routing rules on service containers

See [ADR-0008: Traefik for Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md).

### Entry Points

| Entrypoint | Port | Purpose |
|------------|------|---------|
| web | 80 | HTTP traffic (redirects to HTTPS) |
| websecure | 443 | HTTPS traffic |
| dashboard | 8080 | Traefik dashboard (local access) |

Entry points are configured in the Traefik static configuration.

### Dynamic Routing

Routing rules are defined via Docker labels on service containers. Traefik watches for container changes and updates routing automatically.

See [Docker Labels - Traefik Routing Labels](docker-labels.md#traefik-routing-labels) for label documentation.

### Traefik Configuration

**proxy/docker-compose.yaml:**

See [Configuration Schemas - Proxy Behavior](./configuration-schemas.md#proxy-behavior) for behavioral rules. For complete Traefik configuration examples, see the [Proxy Infrastructure Appendix](appendices/proxy-infrastructure/). The dashboard port (`8080:8080`) is only included when `proxy.dashboard.enabled` is true (default). The `--api.dashboard` flag is set based on the same configuration.

#### Traefik Options

| Option | Value | Description |
|--------|-------|-------------|
| `api.insecure` | `true` | Dashboard accessible without authentication (local development) |
| `providers.docker.watch` | `true` | Live container discovery for dynamic routing updates |
| `accessLog` | `{}` | Access logging enabled for debugging and monitoring |

### TLS Certificate Management

Three modes are supported (see [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)):

#### Auto Mode (mkcert)

1. User runs `mkcert -install` to add local CA
2. User generates wildcard cert: `mkcert "*.scind.test"`
3. Scind discovers certificates in:
   - `~/.config/scind/certs/`
   - Current working directory
   - mkcert default location
4. Certificates are mounted into Traefik container

#### Custom Mode

Place certificates at:
- `~/.config/scind/certs/{domain}.crt`
- `~/.config/scind/certs/{domain}.key`

#### Disabled Mode

No TLS termination at proxy. Services must handle their own TLS.

### DNS Configuration

For local development, configure DNS resolution for the workspace domains. Options include:

1. **dnsmasq**: Route all `*.scind.test` to `127.0.0.1`
   ```
   address=/scind.test/127.0.0.1
   ```
2. **/etc/hosts**: Manual entries for each hostname
3. **Local DNS server**: More complex but flexible

**Note**: The `.test` TLD is reserved by RFC 2606 for testing purposes and will not conflict with real domains or mDNS (unlike `.local`).

### Health Checks and Monitoring

Traefik provides several endpoints for health monitoring:

| Endpoint | Port | Description |
|----------|------|-------------|
| `/ping` | 8080 | Health check endpoint |
| `/api/overview` | 8080 | Dashboard overview API |
| `/dashboard/` | 8080 | Web-based dashboard UI |

The dashboard is only accessible when `proxy.dashboard.enabled` is true.

### Related Decisions

- [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md)
- [ADR-0008: Traefik for Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)
