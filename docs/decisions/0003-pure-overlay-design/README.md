# ADR-0003: Pure Overlay Design

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Applications could embed workspace configuration directly in their docker-compose.yaml files, or workspace integration could be applied externally. The choice affects application portability and the separation between application and workspace concerns.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files that are merged at runtime.

## Consequences

### Positive

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions required in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems
- Applications remain portable across different environments

### Negative

- Generated files add to the directory structure
- Users must understand Docker Compose override file merging behavior

### Neutral

- Application developers define a service contract (`application.yaml`) that declares exports

---

## Notes

The service contract (`application.yaml`) is the only Contrail-specific file in an application repository. It declares what the application exports but doesn't modify core compose configuration.
