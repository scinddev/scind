# ADR-0004: Convention-Based Naming

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Hostnames, network aliases, and project names could be explicitly configured per-service, or derived from conventions.

## Decision

Derive names from conventions:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Protected alias: `{app}-{service}`
- Network name: `{workspace}-internal`
- Project name: `{workspace}-{app}`

## Consequences

### Positive

- Reduces configuration burden
- Ensures consistency across workspaces
- Makes the system predictable and debuggable
- Given workspace and app names, hostnames are deterministic

### Negative

- Less flexibility for unusual naming requirements
- Explicit overrides were considered but removed to keep the schema simple

### Neutral

- Templates can be customized at the workspace level for advanced use cases

---

## Related Documents

- [Naming Conventions Spec](../specs/naming-conventions.md) - Implements the naming conventions defined here

---

<!-- Migrated from specs/contrail-prd.md:194-202 -->
