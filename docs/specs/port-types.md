<!-- Migrated from specs/scind-technical-spec.md:96-135 -->
<!-- Extraction ID: spec-port-types -->

## Port Types and Proxying

Exported services declare ports with a `type` that determines how the port is routed, and optionally a `protocol` for proxied services:

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `*_HOST`, `*_PORT`, `*_SCHEME`, `*_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `*_HOST`, `*_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `*_HOST`, `*_PORT` |

### Port Type Descriptions

- **proxied**: Traffic is routed through Traefik. The exported service gets a hostname (`{workspace}-{app}-{export}.{domain}`) and Traefik labels are generated. Environment variables contain the **proxy values** (hostname and proxy port 80/443), not the container port.
- **assigned**: The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Scind increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

### Port Type Constraints

Each exported service may have:
- At most **one `http`** proxied port
- At most **one `https`** proxied port
- **Multiple `assigned`** ports

If an exported service needs more than one HTTP or HTTPS proxy mapping, create separate exported services.

### Protocol (for proxied type)

When `type: proxied`, the `protocol` field is **required** and specifies how Traefik routes the traffic:

- **https**: Routes through Traefik's `websecure` entrypoint (port 443) with TLS termination
- **http**: Routes through Traefik's `web` entrypoint (port 80)
- **tcp**, **postgresql**, **mysql**, etc. (future): SNI-based TCP routing for database connections. Plugins will handle generating appropriate Traefik TCP router configuration.

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent to collaborators:

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Scind's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

**Docker label exposure**: Visibility is included in the generated Docker labels (`workspace.visibility=public` or `workspace.visibility=protected`), enabling external tools (such as Servlo) to distinguish between public and protected services for display or filtering purposes.

### Private Services

Services not listed in `exported_services` remain private (standard Docker Compose behavior—only accessible within the application's own compose network).
