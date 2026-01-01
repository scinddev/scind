# ADR-0004: Convention-Based Naming

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Hostnames, network aliases, and project names could be:

1. **Explicitly configured**: Users specify every name in configuration files
2. **Convention-derived**: Names follow predictable patterns based on workspace and application names

Explicit configuration offers flexibility but requires more work and can lead to inconsistencies. Convention-based naming is predictable and reduces configuration overhead.

## Decision

Derive names from conventions:

- **Public hostname**: `{workspace}-{app}-{service}.{domain}` (e.g., `dev-app-one-web.contrail.test`)
- **Internal alias**: `{app}-{service}` (e.g., `app-one-web`)
- **Network name**: `{workspace}-internal` (e.g., `dev-internal`)
- **Project name**: `{workspace}-{app}` (e.g., `dev-app-one`)

Template customization is available at the workspace level for advanced users, but defaults work for most cases.

## Consequences

### Positive

- Reduced configuration—names are predictable without explicit specification
- Consistency across all workspaces and applications
- Given workspace and app names, all derived names are deterministic
- Easy to communicate and document

### Negative

- Less flexibility for edge cases (though templates can override)
- Creative naming that produces ambiguous concatenations could cause conflicts

### Neutral

- Explicit overrides were considered but removed to keep the schema simple
- Collision warnings are documented in the technical specification
