# Apex URL Primary Designation

**Status**: Accepted

## Context

Scind applications can export multiple services, each receiving a hostname like `{workspace}-{application}-{exported_service}.{domain}`. A common need is a shorter "apex" hostname at the application level: `{workspace}-{application}.{domain}`.

To generate an apex hostname, Scind must know which exported service is the primary one. xcind solves this by treating the first entry in an ordered bash array as primary. Scind uses YAML maps for exported services, which are unordered by specification — key ordering cannot be relied upon for semantics.

## Decision

Add a `primary: true` boolean field to exported service definitions in `application.yaml`:

- **Single exported service**: implicitly primary (no annotation needed)
- **Multiple exports, one `primary: true`**: that export gets the apex URL
- **Multiple exports, none marked primary**: no apex URL generated
- **Multiple exports, more than one primary**: validation error

All export types (proxied and assigned) can be marked primary:

- **Proxied primary exports** receive: apex hostname, apex internal alias, apex Traefik routing, apex Docker labels, and apex environment variables
- **Assigned-port primary exports** receive: apex internal alias only (no hostname, routing, labels, or environment variables since there is no proxy routing)

## Consequences

### Positive

- Explicit opt-in avoids surprising behavior for multi-export applications
- Works correctly with unordered YAML maps (unlike xcind's position-based approach)
- Single-export applications get apex automatically with zero configuration
- All export types can participate as primary, enabling apex aliases for non-proxied services

### Negative

- Multi-export applications must explicitly annotate their primary (xcind gets this implicitly from ordering)
- One additional field to validate

## Related Documents

- [ADR-0004: Convention-Based Naming](0004-convention-based-naming.md) — Naming patterns extended by apex
- [ADR-0007: Port Type System](0007-port-type-system.md) — Port types that affect apex behavior
- [Naming Conventions](../specs/naming-conventions.md) — Apex naming patterns
- [Configuration Schemas](../specs/configuration-schemas.md) — Primary designation validation rules
