<!-- Migrated from specs/contrail-prd.md:194-202 -->
<!-- Extraction ID: adr-0004-convention-based-naming -->

# Convention-Based Naming

**Status**: Accepted

## Context

Hostnames and aliases could be explicitly configured or derived from conventions.

## Decision

Derive names from conventions:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Protected alias: `{app}-{service}`
- Network name: `{workspace}-internal`

## Consequences

Conventions reduce configuration, ensure consistency, and make the system predictable. Explicit overrides were considered but removed to keep the schema simple.
