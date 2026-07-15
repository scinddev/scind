# Flexible TLS Configuration

**Status**: Accepted

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints (personal dev machines, enterprise networks with managed CAs).

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development—uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments—user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

### Wildcard Proxy Domain Constraint

Under any non-`disabled` TLS mode, the proxy domain used for wildcard TLS **MUST** contain at least one dot (that is, be at least two labels, e.g. `scind.test`, not `scind`). A single-label domain cannot be secured with a wildcard certificate in a portable way:

- Strict RFC 6125 TLS stacks — macOS system curl/LibreSSL, Apple Secure Transport/Safari, and Go's `crypto/tls` — reject a `*.singlelabel` wildcard even when the corresponding SAN is present, and a wildcard matches exactly one label deep.
- Lenient stacks (OpenSSL, Chrome) accept the same certificate, which masks the defect during casual testing and makes it appear environment-specific.

Scind's default proxy domain, `scind.test`, already satisfies this constraint, which is why the limitation stayed latent in the design. Implementations **SHOULD warn** (not fail) when a configured proxy domain has a single label, so users on lenient stacks are not surprised when a teammate on a strict stack cannot connect.

*Note on TLD choice:* domains under HSTS-preloaded TLDs (`.dev`, `.app`) force HTTPS in browsers regardless of Scind's TLS mode; `.io` is **not** preloaded. This does not affect the wildcard constraint but is worth considering when picking a proxy domain.

*Validated by the Xcind proof-of-concept (Xcind ADR-0016).*

### Per-Export TLS Attribute

TLS enforcement is expressed per export via a `tls` attribute, so a single application can mix HTTP-only and HTTPS-required exports:

| Value | Behavior |
|-------|----------|
| `auto` (default) | Serve HTTPS when the TLS mode makes certificates available; otherwise fall back to HTTP. |
| `require` | Always serve HTTPS and redirect HTTP → HTTPS. Incompatible with `tls.mode: disabled`. |
| `disable` | Serve HTTP only for this export, even when the proxy has TLS enabled. |

The per-export `tls` attribute is the export-scoped complement to the workspace-level `tls.mode` above: `tls.mode` decides whether the proxy *can* terminate TLS at all, while the export's `tls` decides what each individual export *does* with that capability. See [Port Types — Per-Export TLS](../specs/port-types.md#per-export-tls) for the attribute definition and [Proxy Infrastructure — HTTP→HTTPS Redirect](../specs/proxy-infrastructure.md#httphttps-redirect) for how `require` is realized.

*Validated by the Xcind proof-of-concept (Xcind ADR-0009).*

## Consequences

- `auto` provides zero-config HTTPS for most users with mkcert installed
- `custom` supports enterprise environments where developers already have CA-signed certs
- Avoids mandating a specific certificate tool while still enabling secure-by-default development

## Related Decisions

- [ADR-0008: Traefik for Reverse Proxy](0008-traefik-reverse-proxy.md) - Traefik performs TLS termination
