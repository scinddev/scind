## Port Types and Proxying

Exported services declare ports with a `type` that determines how the port is routed, and optionally a `protocol` for proxied services:

| Type | Protocol | Behavior | Traefik | Environment Variables |
|------|----------|----------|---------|----------------------|
| `proxied` | `https` | HTTPS proxy via Traefik | Yes (HTTPS router) | `SCIND_{APP}_{SERVICE}_HOST`, `_PORT`, `_SCHEME`, `_URL` |
| `proxied` | `http` | HTTP proxy via Traefik | Yes (HTTP router) | `SCIND_{APP}_{SERVICE}_HOST`, `_PORT`, `_SCHEME`, `_URL` |
| `proxied` | `tcp`, `postgresql`, etc. | SNI-based TCP proxy (future) | Yes (TCP router) | `SCIND_{APP}_{SERVICE}_HOST`, `_PORT` |
| `assigned` | - | Direct port binding, auto-assigned if unavailable | No | `SCIND_{APP}_{SERVICE}_HOST`, `_PORT` |

### Port Type Descriptions

- **proxied**: Traffic is routed through Traefik. The exported service gets a hostname (`{workspace}-{application}-{export}.{domain}`) and Traefik labels are generated. Environment variables contain the **proxy values** (hostname and proxy port 80/443), not the container port.
- **assigned**: The port is bound directly to the host. If the specified port is unavailable (used by another workspace or external process), Scind increments until an available port is found and records the assignment in global state. Environment variables point to the internal alias and assigned host port.

  **Assigned values are live-state-derived, not a pure function of config.** The external port an `assigned` export resolves to depends on the **live global port state** at generation time, not on the application config alone. Another workspace releasing a port (or a garbage-collection pass reclaiming stale reservations) can change what this export resolves to **with no change to this application's config**. Consumers must therefore treat assigned values as state-dependent and re-read them after regeneration rather than caching them across workspace lifecycle events. See [Workspace Lifecycle — Staleness Detection](./workspace-lifecycle.md#staleness-detection) for the generation/staleness contract that governs when these values are recomputed. *Validated by the Xcind proof-of-concept.*

### Port Type Constraints

Each exported service may have:
- At most **one `http`** proxied port
- At most **one `https`** proxied port
- **Multiple `assigned`** ports

If an exported service needs more than one HTTP or HTTPS proxy mapping, create separate exported services.

### Protocol (for proxied type)

When `type: proxied`, the `protocol` field is **required** and specifies how Traefik routes the traffic:

- **https**: Routes through Traefik's HTTPS entrypoint (`proxy.https_entrypoint`, default `websecure`, port 443) with TLS termination
- **http**: Routes through Traefik's HTTP entrypoint (`proxy.http_entrypoint`, default `web`, port 80)
- **tcp**, **postgresql**, **mysql**, etc. (future): SNI-based TCP routing for database connections. Plugins will handle generating appropriate Traefik TCP router configuration.

### Per-Export TLS

Each proxied export may declare a `tls` attribute that controls TLS enforcement **for that export independently**, letting a single application mix HTTP-only and HTTPS-required exports:

| Value | Behavior |
|-------|----------|
| `auto` (default) | Serve HTTPS when the proxy's TLS mode makes certificates available; otherwise fall back to HTTP. |
| `require` | Always serve HTTPS and redirect HTTP → HTTPS. Incompatible with `tls.mode: disabled`. |
| `disable` | Serve HTTP only for this export, even when the proxy has TLS enabled. |

The per-export `tls` attribute is scoped below the workspace-level `tls.mode` (see [ADR-0009: Flexible TLS Configuration](../decisions/0009-flexible-tls-configuration.md)): `tls.mode` decides whether the proxy *can* terminate TLS at all, while each export's `tls` decides what it *does* with that capability. The `require` value's HTTP→HTTPS redirect is realized as described in [Proxy Infrastructure — HTTP→HTTPS Redirect](./proxy-infrastructure.md#httphttps-redirect).

*Validated by the Xcind proof-of-concept (Xcind ADR-0009).*

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). Visibility is **advisory display metadata** — it communicates intent to collaborators and to tools that render service lists. It is **not** access control: Scind enforces nothing on the basis of it, and neither value restricts who can reach the port.

- **public**: This port is intended for external/production use
- **protected** (default): This port exists for development/debugging but should not be depended on in production

Visibility does not change Scind's core behavior—all exported services receive internal network aliases and environment variables regardless of visibility. Both public and protected proxied services route through Traefik.

`private` is **not** a visibility value. A service is private by being absent from `exported_services` (see below), which is a different mechanism at a different level; there is nothing to label, because no labels are generated for it.

**Docker label exposure**: Visibility is emitted on the generated Docker labels — `scind.export.{name}.proxy.{protocol}.visibility` for proxied exports and `scind.export.{name}.port.{internal-port}.visibility` for assigned exports (see [Docker Labels](./docker-labels.md#export-labels)). External tools such as Servlo read these to distinguish public from protected services when displaying or filtering them.

### Private Services

Services not listed in `exported_services` remain private (standard Docker Compose behavior—only accessible within the application's own compose network).

## Port Value Validation

Resolved and inferred port values are validated **at generation time**, not left to fail downstream in Traefik or Docker:

- A resolved port value MUST be an integer in the range **1–65535**. Protocol suffixes (e.g. a `postgresql`/`tcp` designation) are handled explicitly and stripped before the numeric check, rather than being parsed as part of the integer.
- An out-of-range or non-integer port value fails generation with an **export-named error** identifying the offending export, so the problem is reported where it can be fixed rather than surfacing later as an opaque proxy/runtime failure.

*Validated by the Xcind proof-of-concept.*

## Unified Export Descriptor

Although exports resolve through two distinct port types (`proxied` and `assigned`), Scind exposes a **single unified per-export descriptor** keyed by export name, so a consumer sees one consistent view regardless of type. The descriptor carries:

| Field | Description |
|-------|-------------|
| `name` | Export name (the key). |
| `type` | `proxied` or `assigned`. |
| `service` | The Compose service backing the export. |
| `container_port` | The in-container port the service listens on. |
| `host_port_or_url` | For `assigned`: the resolved host port. For `proxied`: the preferred-scheme URL. |
| `tls` | The per-export TLS setting (`auto`/`require`/`disable`; proxied exports only). |
| `apex` | Whether this export is the application's apex (primary) export. |

Merging both port types into one descriptor means downstream consumers do not special-case proxied vs. assigned exports to answer "what/where is this export?". The descriptor is surfaced through the CLI via `app show` / `exports` (the CLI surface is specified elsewhere; this section defines the descriptor shape only).

*Validated by the Xcind proof-of-concept.*

## Related Documents

- [Naming Conventions](./naming-conventions.md) - Hostname and environment variable patterns
- [Proxy Infrastructure](./proxy-infrastructure.md) - Traefik configuration for proxied ports
- [Generated Override Files](./generated-override-files.md) - How port types translate to Docker labels
- [Configuration Schemas](./configuration-schemas.md) - Port configuration in application.yaml
