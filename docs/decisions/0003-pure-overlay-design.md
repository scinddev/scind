<!-- Migrated from specs/contrail-prd.md:181-192 -->
<!-- Extraction ID: adr-0003-pure-overlay-design -->

# Pure Overlay Design (Applications Remain Workspace-Agnostic)

**Status**: Accepted

## Context

Applications could embed workspace configuration, or it could be applied externally.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

## Consequences

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems
