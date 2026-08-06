# Apex URL Primary Designation

**Status**: Accepted

## Context

Scind applications can export multiple services, each receiving a hostname like `{workspace}-{application}-{exported_service}.{domain}`. A common need is a shorter "apex" hostname at the application level: `{workspace}-{application}.{domain}`.

To generate an apex hostname, Scind must know which exported service is the primary one. Xcind solves this by treating the first entry in an ordered structure as primary — a **positional** rule. That positional-only approach was partly a consequence of Xcind's Bash configuration model, not a deliberate rejection of explicit designation; on reconciliation it does not earn a divergence from Scind. Scind should therefore offer **both** an explicit marker and a positional fallback, rather than choosing one.

Two refinements emerged from the Xcind proof-of-concept work:

1. **Apex applies only to proxied exports.** An apex hostname is a proxy concept; assigned exports have no hostname and can never receive an apex hostname. Any notion of "the single export ⇒ implicitly primary" must therefore be scoped to *proxied* exports, so that a mixed application like `web` (proxied) + `db` (assigned) still gets an apex with zero annotation instead of being treated as "multiple exports, ambiguous."
2. **Positional selection needs a defined order.** A positional fallback requires `exported_services` to have a well-defined declaration order. Scind's original framing modeled exported services as an unordered YAML map, for which "first" is undefined. Enabling positional fallback therefore requires either modeling `exported_services` as a sequence, or adopting a documented first-declared rule so ordering is well-defined.
3. **Some applications must be able to decline the apex.** The hybrid rule makes Scind strictly more apex-eager: every application with at least one proxied export now receives an apex hostname, where the earlier "none marked ⇒ no apex" behavior gave a way out. Xcind field experience (Xcind ADR-0017) shows that some applications need no apex — the short `{workspace}-{application}.{domain}` name is unwanted, or it collides with a hostname the developer serves by other means. Xcind expresses this by setting the apex hostname template to an empty string, which disables apex generation for that application. Scind therefore needs an explicit opt-out of its own.

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

### Apex Opt-Out

Add an application-level `apex` boolean to `application.yaml`, defaulting to `true`. Setting `apex: false` declares that the application produces **no apex at all**:

- no apex hostname and no apex internal alias,
- no apex Traefik router,
- no `scind.apex.*` Docker labels,
- no `SCIND_{APPLICATION}_APEX_*` environment variables,
- `apex` and `apex_host` report `null` in the JSON introspection contract, and `scind urls` shows no apex URL for the application.

The opt-out is **application-level, not export-level**, because the apex is an application-level name: there is exactly one apex per application, so declining it is one decision, not one per export. Per-export designation (`primary: true`) is unaffected — an export may still be marked primary under `apex: false`, and that designation simply produces nothing apex-related. Marking an export `primary: true` while the application sets `apex: false` is **not** an error; the opt-out wins, so a shared `application.yaml` can carry the designation for environments that do want the apex.

Xcind validates the capability with an equivalent mechanism (an empty apex hostname template disables apex generation, Xcind ADR-0017); Scind expresses it as a typed field rather than an empty template value, because an empty string is an easy accident and a boolean is not.

## Consequences

### Positive

- Explicit `primary: true` gives multi-export applications unambiguous, order-independent control
- Positional fallback means single-proxied-export and mixed proxied/assigned applications get a correct apex with zero annotation
- Scoping apex eligibility to proxied exports makes `web` (proxied) + `db` (assigned) apps "just work"
- All export types can participate as explicit primary, enabling apex aliases for non-proxied services
- Subsumes Xcind's positional behavior rather than diverging from it — both mechanisms are available
- `apex: false` gives applications an explicit way to decline the apex, which the hybrid rule would otherwise force on every proxied application

### Negative

- Requires `exported_services` to carry a defined declaration order to make positional fallback well-defined (a sequence, or a documented first-declared rule)
- Two selection paths (explicit and positional) to implement and document instead of one
- Two additional fields to validate (`primary` per export, `apex` per application)
- The apex can now be absent for two different reasons — no proxied export, or an explicit opt-out — so reporting surfaces must render `null` apex values in both cases

## Related Documents

- [ADR-0004: Convention-Based Naming](0004-convention-based-naming.md) — Naming patterns extended by apex
- [ADR-0007: Port Type System](0007-port-type-system.md) — Port types that affect apex behavior
- [Naming Conventions](../specs/naming-conventions.md) — Apex naming patterns
- [Configuration Schemas](../specs/configuration-schemas.md) — Primary designation validation rules
