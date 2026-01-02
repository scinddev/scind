# ADR-0003: Pure Overlay Design (Applications Remain Workspace-Agnostic)

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Applications could embed workspace configuration (labels, network definitions, environment variables), or this integration could be applied externally.

## Decision

Applications' own `docker-compose.yaml` files have no knowledge of workspaces. All workspace integration is achieved through generated Docker Compose override files.

## Consequences

### Positive

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions in application code
- Workspace concerns are cleanly separated from application concerns
- Same application can participate in multiple workspace systems

### Negative

- Requires a generation step before running
- Generated files must be kept in sync with source configuration

### Neutral

- Override files are stored in `.generated/` and gitignored

---

## Related Documents

- [Generated Override Files Spec](../specs/generated-override-files.md) - Implements the overlay generation mechanism
- [Environment Variables Spec](../specs/environment-variables.md) - Defines the environment injection approach

---

<!-- Migrated from specs/contrail-prd.md:181-192 -->
