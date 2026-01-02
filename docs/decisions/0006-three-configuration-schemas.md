# ADR-0006: Three Configuration Schemas

**Status**: Accepted
**Date**: 2024-12
**Decision-Makers**: Contrail Design Team

---

## Context

Configuration could be in one monolithic file or separated by concern and ownership.

## Decision

Three schema types that can be combined:
- `proxy`: Global/per-user (domain, Traefik settings)
- `workspace`: Per-workspace (name, application list)
- `application`: Per-application (flavors, exported services)

## Consequences

### Positive

- Separation of concerns—proxy config rarely changes, workspace config defines the environment, application config is owned by the application team
- Application config (`application.yaml`) can live in the application's own repository
- Changes to one layer don't require touching others

### Negative

- Multiple files to understand and maintain
- Configuration hierarchy must be documented

### Neutral

- Proxy config: `~/.config/contrail/proxy.yaml`
- Workspace config: `{workspace}/workspace.yaml`
- Application config: `{app}/application.yaml`

---

## Related Documents

- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Full specification of each configuration schema
- [Configuration Reference](../reference/configuration.md) - Quick reference for configuration files

---

<!-- Migrated from specs/contrail-prd.md:224-233 -->
