# ADR-0009: Flexible TLS Configuration

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints:

- **Personal dev machines**: Can install local CA via mkcert
- **Enterprise networks**: May have managed CAs with pre-issued certificates
- **CI/Docker-in-Docker**: May need HTTP-only or self-signed

Mandating a specific certificate tool would limit adoption in enterprise environments.

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development—uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments—user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

```yaml
proxy:
  tls:
    mode: auto  # or: custom, disabled
    # For mode: custom
    cert_file: ~/.config/contrail/certs/wildcard.crt
    key_file: ~/.config/contrail/certs/wildcard.key
```

## Consequences

### Positive

- Zero-config HTTPS for most users with mkcert installed
- Enterprise environments can use existing CA-signed certificates
- Flexibility for edge cases without mandating specific tools

### Negative

- Multiple modes to document and test
- Self-signed fallback results in browser warnings

### Neutral

- Certificate files for custom mode are user-managed
- mkcert must be separately installed and configured (`mkcert -install`)
