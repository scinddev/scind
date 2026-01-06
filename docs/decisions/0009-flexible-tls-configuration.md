# ADR-0009: Flexible TLS Configuration

## Status

Accepted

## Context

HTTPS support for local development requires TLS certificates. Different environments have different constraints (personal dev machines, enterprise networks with managed CAs).

## Decision

Support three TLS modes via `proxy.yaml`:

| Mode | Use Case |
|------|----------|
| `auto` | Personal development - uses mkcert if available, falls back to self-signed |
| `custom` | Enterprise environments - user provides cert/key signed by enterprise CA |
| `disabled` | HTTP-only development (not recommended) |

## Consequences

### Positive

- `auto` provides zero-config HTTPS for most users with mkcert installed
- `custom` supports enterprise environments where developers already have CA-signed certs
- Avoids mandating a specific certificate tool while still enabling secure-by-default development

### Negative

- Multiple modes add complexity to documentation
- Users must understand which mode fits their environment

### Neutral

- Default mode is `auto` for simplicity

## Related Documents

- [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md) - Implements TLS configuration options

<!-- Migrated from specs/contrail-prd.md:272-289 -->
