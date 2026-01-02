# ADR-0009: Flexible TLS Configuration

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints: personal dev machines can use locally-trusted CAs, enterprise networks may have managed CAs, and some environments may not support HTTPS at all.

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development - uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments - user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

```yaml
proxy:
  tls:
    mode: auto  # or custom, disabled
    # For mode: custom
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

## Consequences

### Positive

- `auto` provides zero-config HTTPS for most users with mkcert installed
- `custom` supports enterprise environments with existing CA-signed certs
- Avoids mandating a specific certificate tool
- Enables secure-by-default development

### Negative

- Users must install mkcert separately for browser-trusted certificates
- `disabled` mode exists but is discouraged

### Neutral

- Default self-signed certificates work but trigger browser warnings

---

## Notes

Certificate setup for `auto` mode: Run `mkcert -install` once per machine to add the local CA to your trust store, then `mkcert "*.contrail.test"` to generate a wildcard certificate.
