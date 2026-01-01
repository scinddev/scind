# ADR-0003: Pure Overlay Design (Applications Remain Workspace-Agnostic)

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Applications could be designed with two approaches:

1. **Embedded workspace awareness**: Applications include workspace-specific configuration in their own `docker-compose.yaml` files
2. **Pure overlay**: Applications have no knowledge of workspaces; all integration is applied externally through generated override files

The choice affects portability, maintainability, and vendor lock-in.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

The generated override files add:
- Network attachments (`{workspace}-internal`, `contrail-proxy`)
- Network aliases for service discovery
- Traefik labels for routing
- Environment variables for service discovery (`CONTRAIL_*`)
- Project name for isolation

## Consequences

### Positive

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions required in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems
- Application developers don't need to know about Contrail internals

### Negative

- Generated override files must be kept in sync with application changes
- Some complexity is hidden in the generation layer
- Debugging may require inspecting generated files

### Neutral

- Applications define a service contract (`application.yaml`) describing what they export
- Generated files are gitignored and recreated on demand
