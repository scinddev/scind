## Proxy Layer

### Architecture Overview

Traefik serves as the reverse proxy, routing external requests to workspace services.

```
[External Request] → [Traefik:443/80] → [scind-proxy network] → [Service Container]
```

#### Components

- **Traefik container**: Single instance managing all workspace routing
- **Shared proxy network**: Host-level Docker network connecting Traefik to services. Its name comes from `proxy.network` (default `scind-proxy`); under [External Proxy Mode](#external-proxy-mode) it names the network a foreign proxy already lives on.
- **Dynamic configuration**: Label-based routing rules on service containers

See [ADR-0008: Traefik for Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md).

### Entry Points

| Entrypoint | Port | Purpose |
|------------|------|---------|
| `proxy.http_entrypoint` (default `web`) | `proxy.http_port` (default 80) | HTTP traffic |
| `proxy.https_entrypoint` (default `websecure`) | `proxy.https_port` (default 443) | HTTPS traffic (only when `tls.mode != disabled`) |

Entry points are configured in the Traefik static configuration. Entrypoint **names** are configurable, not just their ports, so generated router labels can target an external proxy's entrypoints (see [External Proxy Mode](#external-proxy-mode)); the managed defaults are `web`/`websecure`.

HTTP traffic is **not** blanket-redirected to HTTPS. Redirection happens per export, only for exports declaring `tls: require` (see [HTTP→HTTPS Redirect](#httphttps-redirect)); `tls: auto` and `tls: disable` exports are served over HTTP on this entrypoint.

The dashboard has **no entrypoint of its own**. When `proxy.dashboard.enabled` is true, Scind adds `--api.dashboard=true`, an `api: {dashboard: true, insecure: true}` block to the static configuration, and a host publish of `proxy.dashboard.port` — the dashboard rides Traefik's built-in insecure-API listener. All dashboard settings are inert under `mode: external`, where Scind generates no static configuration and starts no container.

### Dynamic Routing

Routing rules are defined via Docker labels on service containers. Traefik watches for container changes and updates routing automatically.

See [Docker Labels - Traefik Routing Labels](docker-labels.md#traefik-routing-labels) for label documentation.

### HTTP→HTTPS Redirect

When an export requires TLS (per-export `tls: require`, see [Port Types — Per-Export TLS](port-types.md#per-export-tls)), HTTP requests are redirected to HTTPS. The redirect is realized with a Traefik `redirectscheme` middleware attached to a redirect-only HTTP router:

- The redirect-only HTTP router listens on the HTTP entrypoint (`proxy.http_entrypoint`, default `web`) and matches the export's host.
- It applies a shared `redirectscheme` middleware (`scheme: https`, `permanent: true`) that rewrites the request to HTTPS before it reaches any backend. The redirect is **permanent (301)**, not temporary: `tls: require` is a stable policy declared in configuration, not a transient runtime state, so caching it in clients is correct rather than a hazard.
- Because the request is redirected before a backend is ever selected, the redirect-only router may point its service at Traefik's built-in `noop@internal` backend — no real loadbalancer service need be resolved.

**Implementation constraint — emit the middleware definition on every service block.** Traefik's Docker provider loads labels only from **running** containers. If the shared `redirectscheme` middleware is *defined* on only one "first" service, its definition disappears whenever that container is down, and every router referencing it dangles (redirects break workspace-wide). Therefore the middleware **definition** labels must be emitted on **every** rendered service block that references it, not just one. Repeated identical definitions are idempotent — Traefik reconciles them to a single middleware — so there is no cost to the duplication.

*Validated by the Xcind proof-of-concept (Xcind ADR-0009).*

### Traefik Configuration

**proxy/docker-compose.yaml:**

See [Configuration Schemas - Proxy Behavior](./configuration-schemas.md#proxy-behavior) for behavioral rules. For complete Traefik configuration examples, see the [Proxy Infrastructure Appendix](appendices/proxy-infrastructure/). The dashboard port (`8080:8080`) is only included when `proxy.dashboard.enabled` is true (default). The `--api.dashboard` flag is set based on the same configuration.

The generated static configuration follows the same `proxy.network`, `proxy.http_entrypoint`, and `proxy.https_entrypoint` values that the generated labels use, so managed and external mode share a single label-generation path with no per-mode branching; the managed defaults (`scind-proxy`, `web`, `websecure`) reproduce Scind's historical output unchanged. Neither file is generated under `mode: external`.

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

**Which routers exist vs. who provisions certificates.** These are two separate questions, and the configuration answers them with different keys. `tls.mode` and the per-export `tls:` attribute decide **which routers exist** — whether an HTTPS router is generated for an export at all, and whether its HTTP sibling redirects. `proxy.certresolver` decides **who provisions the certificate** those routers serve: when set, every HTTPS router (per-export and apex) carries `traefik.http.routers.{name}.tls.certresolver={value}` and delegates certificate issuance to the proxy's own ACME machinery instead of to Scind's cascade below. The two are independent — a managed Traefik with an ACME resolver uses both — and under `mode: external` Scind provisions nothing at all.

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

#### Division of TLS Labor Under External Mode

Under `mode: external` the external proxy terminates TLS, so Scind's certificate machinery stands down entirely:

- the auto-mode cascade (user-provided cert → cached cert → mkcert → OpenSSL self-signed) is **skipped in full**,
- no `dynamic/tls.yaml` is written,
- no certificates are mounted anywhere — there is no Scind-managed container to mount them into.

**Router emission rules are unchanged.** `tls.mode` and the per-export `tls: auto|require|disable` attribute continue to decide which routers exist and which redirect, exactly as in managed mode; only provisioning is removed.

`tls.mode: custom` together with `mode: external` **MUST be rejected** at config validation (and again at `scind proxy init`) with a clear error. `custom` names cert and key files, and in external mode there is nowhere those files would ever be installed — accepting the combination would silently ignore an explicit certificate choice.

### External Proxy Mode

See [ADR-0017: External Proxy Mode](../decisions/0017-external-proxy-mode.md). Under `proxy.mode: external`, a proxy Scind does not manage — a Coolify host's Traefik, for example — serves the applications. Scind only attaches application containers to the configured shared network (`proxy.network`, pointed at the external proxy's network) and emits router labels using the configured entrypoint names (`proxy.http_entrypoint`, `proxy.https_entrypoint`) and optional `proxy.certresolver`.

#### Generation and Initialization

Under `mode: external`, `scind proxy init` skips every generation step — no compose file, no static configuration, no `dynamic/tls.yaml`, no certificate provisioning. Instead it:

1. Writes the configuration as usual.
2. Removes the stale managed static configuration and `dynamic/tls.yaml`.
3. Leaves `certs/` in place. Removing it is not required (nothing reads it in external mode) and keeping it preserves a user-supplied wildcard certificate across mode flip-flops.
4. Removes the generated compose file **only when no container labeled `scind.component=proxy` is running**. While such a container exists the compose file is the only handle for stopping it, so `init` retains the file and warns that a Scind-managed proxy is still running.
5. Verifies the configured shared network, creating it if absent.

#### Network Topology

Two wirings work.

**Recommended** — point `proxy.network` at the network the external proxy already lives on (e.g. `coolify`). Generated overrides declare that network `external: true`, so applications simply join it and the host proxy container is never touched. This is critical on platform-managed hosts, which recreate their proxy container on upgrade or redeploy.

**Alternate** — keep `proxy.network: scind-proxy` and run `docker network connect scind-proxy {host-proxy-container}`. This needs no knowledge of the host proxy's own configuration, but the connection is lost whenever that container is recreated and **must be re-applied** each time.

When the configured network is absent, Scind still **creates** it — the generated overrides declare it `external: true` and would fail to start otherwise — but **warns loudly** in external mode, because a freshly created network has no proxy attached to it and nothing will route until one is.

#### Command Behavior Under `mode: external`

The governing principle: **refusal beats a silent no-op.** These verbs promise actions — stop the proxy, recreate its network, tail its logs — that Scind cannot perform on a proxy it does not own. Exiting 0 having done nothing would misreport success for a request that was never carried out.

| Command | Behavior | Exit |
|---------|----------|------|
| `proxy init` | Writes config, removes stale managed artifacts (retaining the compose file and warning when a managed proxy still runs), verifies/creates the network | 0 |
| `proxy up` | Verifies the network; reports a detected proxy-like container; starts nothing | 0 |
| `proxy up --force` / `--recreate` | Refuses — would remove a network shared with the external proxy | 1 |
| `proxy down` | Refuses, except the managed→external migration escape hatch below | 1 / 0 |
| `proxy restart` | Refuses — there is no Scind-owned proxy to restart | 1 |
| `proxy logs` | Refuses, with a `docker logs {container}` hint naming the detected proxy container | 1 |
| `proxy status` | Renders an external-specific report (see [CLI Reference](../reference/cli.md#scind-proxy-status)) | 0 |
| `proxy destroy [--purge] [--force]` | Runs: stops a retained Scind-managed proxy via the previously generated compose file (leaving state intact for a retry if that teardown fails) and removes generated state, but **never** removes the configured shared network | 0 |

**Migration escape hatch for `proxy down`.** `down` stops a leftover Scind-managed proxy and exits 0 when **both** predicates hold:

1. a container labeled `scind.component=proxy` is running, **and**
2. the previously generated managed compose file still exists.

Either predicate alone leaves the refusal in place. The label alone could match a proxy the user ran by hand, which Scind has no mandate to stop; the compose file alone could be stale, describing a proxy that is long gone. Together they identify a proxy Scind itself started and still knows how to stop — exactly the managed→external migration case.

**Why `destroy` still runs.** `destroy` targets Scind's *own* generated state, which exists in external mode too. It must not remove the configured shared network: that network is owned by the external proxy and may carry unrelated workloads.

#### Known Caveats

Accepted and documented rather than worked around:

- **Host-level HTTP→HTTPS redirect.** An external proxy with an entrypoint-level redirect (Coolify's "Force HTTPS", for example) shadows Scind's HTTP routers *before* routing. This is harmless — requests still land on the HTTPS routers — but a per-export `tls: require` redirect will appear to come from the host proxy rather than from Scind's `scind-redirect-https` middleware, which can be confusing when diagnosing a redirect.
- **Docker-provider constraints.** A host proxy using `providers.docker.constraints` will not discover Scind containers at all unless the constrained label is added to the application's labels. Handling this automatically is out of scope; it is documented so operators recognize the symptom (containers on the shared network, no routers in the proxy).

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
- [ADR-0017: External Proxy Mode](../decisions/0017-external-proxy-mode.md)
