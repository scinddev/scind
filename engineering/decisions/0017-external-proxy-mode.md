# External Proxy Mode

**Status**: Accepted

## Context

[ADR-0008](0008-traefik-reverse-proxy.md) chose Traefik with the Docker provider as Scind's reverse proxy, and Scind runs that Traefik itself, bound to host ports `80`/`443`. On a host that already owns those ports with its own Docker-provider Traefik — the motivating case is a Coolify server, whose proxy owns `80`/`443`, lives on its own Docker network, names its entrypoints `http`/`https`, and terminates TLS with an ACME certresolver named `letsencrypt` — the Scind-managed proxy cannot bind its ports at all. Moving it to alternate ports works but produces second-class URLs and two proxy stacks doing the same job.

`proxy.auto_start: false` (see [Configuration Schemas — Proxy Auto-Start](../specs/configuration-schemas.md#proxy-auto-start)) names this case but only solves half of it: it suppresses the container start while leaving every routing value hardcoded. The network name (`scind-proxy`), the entrypoint names (`web`/`websecure`), and the absence of any certresolver label mean the labels Scind emits are unconsumable by a foreign Traefik — the proxy is not started, and nothing routes either.

Because Scind's routing is *already* expressed as Docker-provider labels, an existing host Traefik can serve Scind applications directly, provided the labels name **its** entrypoints, the application containers share a network with it, and HTTPS routers name its certresolver. Those three hardcoded values are the entire obstacle.

## Decision

Add a `proxy.mode` key to `proxy.yaml` with two values:

- **`managed`** (default) — exactly today's behavior: Scind generates, starts, and stops its own Traefik.
- **`external`** — Scind never starts, stops, or configures a proxy. Applications join a configured shared network and emit labels tuned for the host proxy; that proxy's Docker provider discovers them.

Four further keys make the emitted label surface configurable. All are meaningful in **both** modes, and all are set the ordinary way — typed YAML in `proxy.yaml`, or `scind config set`:

| Field (`proxy.yaml`) | Type | Default | Coolify-style value |
|----------------------|------|---------|---------------------|
| `proxy.mode` | string | `managed` | `external` |
| `proxy.network` | string | `scind-proxy` | `coolify` |
| `proxy.http_entrypoint` | string | `web` | `http` |
| `proxy.https_entrypoint` | string | `websecure` | `https` |
| `proxy.certresolver` | string | (empty) | `letsencrypt` |

When `proxy.certresolver` is non-empty, every generated HTTPS router — per-export and apex — gains a `traefik.http.routers.{name}.tls.certresolver={value}` label. This is mode-independent: a managed Traefik configured with an ACME resolver benefits identically.

Managed mode's generated `traefik.yaml` follows the same configured network and entrypoint names, so **label generation stays a single code path with no per-mode branching**. The managed defaults reproduce the current output byte-for-byte.

### Network topology

Two wirings work.

**Recommended**: point `proxy.network` at the network the host proxy already lives on (e.g. `coolify`). Generated overrides declare that network `external: true`, so applications simply join it and the host proxy container is **never touched**. This matters on platform-managed hosts, which recreate their proxy container on upgrade or redeploy.

**Alternate**: keep `proxy.network: scind-proxy` and run `docker network connect scind-proxy {host-proxy-container}`. This needs no awareness of the host proxy's own configuration, but carries exactly that recreation caveat — the connect must be re-applied whenever the host proxy container is recreated.

When the configured network does not exist, Scind still creates it (the overrides declare it `external: true` and would otherwise fail to start) but **warns loudly** in external mode, because a freshly created network by definition has no proxy attached to it.

### Division of TLS labor

`tls.mode` keeps its meaning of *which routers exist*, and the per-export `tls: auto|require|disable` attribute ([ADR-0009](0009-flexible-tls-configuration.md)) is unchanged. What external mode removes is certificate **provisioning**: the auto-mode cascade is skipped in full, no `dynamic/tls.yaml` is written, and no certificates are mounted — the external proxy terminates TLS. Consequently `tls.mode: custom` together with `mode: external` is a **hard error**: it names certificate files Scind would never install anywhere.

### Command semantics

`scind proxy` verbs promise actions Scind cannot perform on a proxy it does not own. Under `mode: external` they therefore **refuse** rather than silently no-op, with the exception of the verbs that can still be satisfied honestly (`init`, `up`, `status`, `destroy`). The full command table lives in [Proxy Infrastructure — External Proxy Mode](../specs/proxy-infrastructure.md#external-proxy-mode).

## Consequences

### Positive

- A Coolify (or any existing-Traefik) host serves Scind applications on real ports `80`/`443` with ACME certificates, one proxy hop, and zero modification of the host proxy.
- Managed-mode behavior and generated output are unchanged; the new keys take their defaults when absent.
- `proxy.certresolver` benefits managed mode too — an ACME-configured managed Traefik gains per-router resolvers without any external-mode involvement.
- Label generation remains one code path, so managed and external mode cannot drift apart.

### Negative

- The label vocabulary Scind emits becomes a semi-public contract with foreign proxies; entrypoint names and certresolver are the configuration surface for that contract.
- Switching `external` → `managed` requires setting the managed values back explicitly (e.g. `proxy.network: scind-proxy`), since the keys are persisted independently of the mode.
- Several `scind proxy` verbs now have a mode-dependent outcome, which is more surface to document and test.

### Neutral

- `proxy.auto_start` and `proxy.mode` both describe "Scind does not start a proxy", but they answer different questions and do not conflict: `mode: external` short-circuits before `auto_start` is consulted.

## Related Documents

- [ADR-0002: Two-Layer Networking](0002-two-layer-networking.md) — the shared proxy network layer this mode re-points at an existing network
- [ADR-0008: Traefik for Reverse Proxy](0008-traefik-reverse-proxy.md) — Docker-provider labels are what make a foreign Traefik a drop-in consumer
- [ADR-0009: Flexible TLS Configuration](0009-flexible-tls-configuration.md) — TLS mode semantics this mode builds on
- [Proxy Infrastructure](../specs/proxy-infrastructure.md) — external-mode network topology, caveats, and command table

*Credit to the Xcind proof-of-concept (Xcind ADR-0022) for validating the mode, the configurable label surface, and the refuse-don't-no-op command semantics (Xcind PR #83).*
