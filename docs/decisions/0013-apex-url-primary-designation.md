# Apex URL Primary Designation

**Status**: Accepted

## Context

Scind applications can export multiple services, each receiving a hostname like `{workspace}-{application}-{exported_service}.{domain}`. A common need is a shorter "apex" hostname at the application level: `{workspace}-{application}.{domain}`.

To generate an apex hostname, Scind must know which exported service is the primary one. Xcind solves this by treating the first entry in an ordered structure as primary — a **positional** rule. That positional-only approach was partly a consequence of Xcind's Bash configuration model, not a deliberate rejection of explicit designation; on reconciliation it does not earn a divergence from Scind. Scind should therefore offer **both** an explicit marker and a positional fallback, rather than choosing one.

Two refinements emerged from the Xcind proof-of-concept work:

1. **Apex applies only to proxied exports.** An apex hostname is a proxy concept; assigned exports have no hostname and can never receive an apex hostname. Any notion of "the single export ⇒ implicitly primary" must therefore be scoped to *proxied* exports, so that a mixed application like `web` (proxied) + `db` (assigned) still gets an apex with zero annotation instead of being treated as "multiple exports, ambiguous."
2. **Positional selection needs a defined order.** A positional fallback requires `exported_services` to have a well-defined declaration order. Scind's original framing modeled exported services as an unordered YAML map, for which "first" is undefined. Enabling positional fallback therefore requires either modeling `exported_services` as a sequence, or adopting a documented first-declared rule so ordering is well-defined.

## Decision

Add a `primary: true` boolean field to exported service definitions in `application.yaml`, and select the apex-bearing export with a **hybrid** rule scoped to proxied exports:

- **Explicit wins**: an export marked `primary: true` always receives the apex.
- **Positional fallback**: when multiple **proxied** exports exist and none is marked primary, the **first-declared proxied export** becomes primary and an apex is still generated. (This supersedes the earlier "none marked ⇒ no apex" behavior.)
- **Implicit single**: when there is exactly **one proxied** export, it is implicitly primary with no annotation. Assigned exports present alongside it do not affect this determination.
- **Validation error**: more than one export marked `primary: true`.

Apex eligibility is **scoped to proxied exports**. Assigned exports never count toward implicit-primary eligibility and never receive an apex hostname; they may still be marked `primary: true` explicitly to receive the apex internal alias.

Enabling the positional fallback requires giving `exported_services` a **defined declaration order** — model it as a sequence, or adopt a documented first-declared rule — so that "first proxied export" is well-defined.

All export types can be marked primary, but what they receive differs:

- **Proxied primary exports** receive: apex hostname, apex internal alias, apex Traefik routing, apex Docker labels, and apex environment variables
- **Assigned-port primary exports** receive: apex internal alias only (no hostname, routing, labels, or environment variables since there is no proxy routing)

This is a deliberate product decision to provide **both** mechanisms. Explicit `primary: true` gives multi-export applications unambiguous control; positional fallback means most applications — including single-proxied-export and mixed proxied/assigned applications — get a correct apex with zero annotation. Xcind's positional-only behavior is thus subsumed rather than diverged from.

## Consequences

### Positive

- Explicit `primary: true` gives multi-export applications unambiguous, order-independent control
- Positional fallback means single-proxied-export and mixed proxied/assigned applications get a correct apex with zero annotation
- Scoping apex eligibility to proxied exports makes `web` (proxied) + `db` (assigned) apps "just work"
- All export types can participate as explicit primary, enabling apex aliases for non-proxied services
- Subsumes Xcind's positional behavior rather than diverging from it — both mechanisms are available

### Negative

- Requires `exported_services` to carry a defined declaration order to make positional fallback well-defined (a sequence, or a documented first-declared rule)
- Two selection paths (explicit and positional) to implement and document instead of one
- One additional field to validate

## Related Documents

- [ADR-0004: Convention-Based Naming](0004-convention-based-naming.md) — Naming patterns extended by apex
- [ADR-0007: Port Type System](0007-port-type-system.md) — Port types that affect apex behavior
- [Naming Conventions](../specs/naming-conventions.md) — Apex naming patterns
- [Configuration Schemas](../specs/configuration-schemas.md) — Primary designation validation rules
