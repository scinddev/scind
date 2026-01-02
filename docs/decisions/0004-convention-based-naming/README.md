# ADR-0004: Convention-Based Naming

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Hostnames, aliases, project names, and network names could either be explicitly configured for each service or derived from conventions. Explicit configuration offers flexibility but increases configuration burden.

## Decision

Derive names from conventions:

- **Public hostname**: `{workspace}-{app}-{service}.{domain}` (e.g., `dev-app-one-web.contrail.test`)
- **Internal alias**: `{app}-{service}` (e.g., `app-one-web`)
- **Network name**: `{workspace}-internal` (e.g., `dev-internal`)
- **Project name**: `{workspace}-{app}` (e.g., `dev-app-one`)

## Consequences

### Positive

- Reduces configuration burden significantly
- Ensures consistency across all workspaces and applications
- Makes the system predictable and debuggable
- Given workspace and app names, all derived names are deterministic

### Negative

- Less flexibility for edge cases requiring custom naming
- Potential for naming collisions with creative workspace/app name combinations

### Neutral

- Template customization is available for advanced users who need to modify patterns

---

## Notes

Explicit override options were considered but removed to keep the schema simple. Template variables (`%WORKSPACE_NAME%`, etc.) provide escape hatch for advanced customization.
