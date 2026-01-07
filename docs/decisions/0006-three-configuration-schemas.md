<!-- Migrated from specs/contrail-prd.md:224-233 -->
<!-- Extraction ID: adr-0006-three-configuration-schemas -->

# Three Configuration Schemas

**Status**: Accepted

## Context

Configuration could be in one monolithic file or separated by concern.

## Decision

Three schema types that can be combined:
- `proxy`: Global/per-user (domain, Traefik settings)
- `workspace`: Per-workspace (name, application list)
- `application`: Per-application (flavors, exported services)

## Consequences

Separation of concerns—proxy config rarely changes, workspace config defines the environment, application config is owned by the application team.
