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

### HTTP→HTTPS Redirect

When an export requires TLS (per-export `tls: require`, see [Port Types — Per-Export TLS](port-types.md#per-export-tls)), HTTP requests are redirected to HTTPS. The redirect is realized with a Traefik `redirectscheme` middleware attached to a redirect-only HTTP router:

- The redirect-only HTTP router listens on the `web` entrypoint and matches the export's host.
- It applies a shared `redirectscheme` middleware (`scheme: https`) that rewrites the request to HTTPS before it reaches any backend.
- Because the request is redirected before a backend is ever selected, the redirect-only router may point its service at Traefik's built-in `noop@internal` backend — no real loadbalancer service need be resolved.

**Implementation constraint — emit the middleware definition on every service block.** Traefik's Docker provider loads labels only from **running** containers. If the shared `redirectscheme` middleware is *defined* on only one "first" service, its definition disappears whenever that container is down, and every router referencing it dangles (redirects break workspace-wide). Therefore the middleware **definition** labels must be emitted on **every** rendered service block that references it, not just one. Repeated identical definitions are idempotent — Traefik reconciles them to a single middleware — so there is no cost to the duplication.

*Validated by the Xcind proof-of-concept (Xcind ADR-0009).*

### Traefik Configuration

**proxy/docker-compose.yaml:**

See [Configuration Schemas - Proxy Behavior](./configuration-schemas.md#proxy-behavior) for behavioral rules. For complete Traefik configuration examples, see the [Proxy Infrastructure Appendix](appendices/proxy-infrastructure/). The dashboard port (`8080:8080`) is only included when `proxy.dashboard.enabled` is true (default). The `--api.dashboard` flag is set based on the same configuration.

#### TLS-Mode-Conditional Emission

The HTTPS-related pieces of both the static configuration and the compose file are emitted **only when `tls.mode != disabled`**, mirroring how the dashboard is already conditional on `proxy.dashboard.enabled`. Specifically, the following are omitted under `tls.mode: disabled`:

- the `websecure` (`:443`) entrypoint in the static configuration,
- the `file` provider watching the dynamic directory (there are no certificates or TLS middlewares to load),
- the `443:443` port mapping in the compose file,
- the `./certs` and `./dynamic` bind-mounts in the compose file.

This avoids a proxy that advertises HTTPS on `:443` and watches a certificate directory that is never populated. When TLS is enabled, all four are emitted. See the [Proxy Infrastructure Appendix](appendices/proxy-infrastructure/) for the annotated static config and compose file.

#### Traefik Options

| Option | Value | Description |
|--------|-------|-------------|
| `api.insecure` | `true` | Dashboard accessible without authentication (local development) |
| `providers.docker.watch` | `true` | Live container discovery for dynamic routing updates |
| `accessLog` | `{}` | Access logging enabled for debugging and monitoring |

### TLS Certificate Management

Three modes are supported (see [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)):

#### Wildcard Proxy Domain Constraint

Under any non-`disabled` TLS mode, the wildcard certificate is minted for the proxy domain (`*.{domain}`), and that domain **MUST** be at least two labels (contain at least one dot). Strict RFC 6125 stacks (macOS system curl/LibreSSL, Apple Secure Transport/Safari, Go `crypto/tls`) reject a `*.singlelabel` wildcard even when the SAN is present, while lenient stacks (OpenSSL, Chrome) accept it and mask the problem. Implementations SHOULD warn (not fail) on a single-label domain. The default `scind.test` satisfies the constraint. See [ADR-0009 — Wildcard Proxy Domain Constraint](../decisions/0009-flexible-tls-configuration.md#wildcard-proxy-domain-constraint).

#### Auto Mode

Auto mode resolves a wildcard certificate for the proxy domain through the following cascade, using the first source that yields a usable cert/key pair:

1. **User-provided cert** — an existing `*.{domain}` cert/key pair placed in `~/.config/scind/certs/` (or discovered in the working directory / mkcert default location).
2. **Cached cert** — a previously provisioned cert cached under `~/.config/scind/certs/`, provided it still matches the configured domain (see domain-change regeneration below).
3. **mkcert** — if `mkcert` is installed (with its local CA trusted via `mkcert -install`), generate `mkcert "*.{domain}"`.
4. **OpenSSL self-signed fallback** — if mkcert is unavailable, generate a self-signed wildcard cert with OpenSSL so HTTPS still terminates (browsers will warn until the CA is trusted).

The resolved certificates are mounted into the Traefik container.

**Domain-change regeneration.** Scind records the domain a cached wildcard certificate was minted for (e.g. in a `domain` marker alongside the cert). When the configured proxy domain changes, the recorded marker no longer matches, and auto mode **re-mints** the wildcard for the new domain instead of serving a stale-CN certificate. This prevents a domain change from silently continuing to present a certificate whose CN/SAN no longer covers the active hostnames.

*Validated by the Xcind proof-of-concept.*

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
