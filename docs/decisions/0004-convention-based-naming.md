# Convention-Based Naming

**Status**: Accepted

## Context

Hostnames and aliases could be explicitly configured or derived from conventions.

## Decision

Derive names from conventions:
- Public hostname: `{workspace}-{application}-{service}.{domain}`
- Internal alias: `{application}-{service}`
- Network name: `{workspace}-internal`

## Consequences

Conventions reduce configuration, ensure consistency, and make the system predictable. Explicit overrides were considered but removed to keep the schema simple.
