# ADR-0006: Three Configuration Schemas

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Design Team

<!-- Migrated from specs/contrail-prd.md -->

## Context

Configuration could be organized as:

1. **Single monolithic file**: All settings in one place
2. **Separated by concern**: Different files for different scopes and owners

A monolithic approach is simpler initially but doesn't scale well when different stakeholders own different parts of the configuration.

## Decision

Three schema types that can be combined:

| Schema | Location | Owner | Purpose |
|--------|----------|-------|---------|
| `proxy` | `~/.config/contrail/proxy.yaml` | User/global | Domain, TLS, Traefik settings |
| `workspace` | `{workspace}/workspace.yaml` | Workspace maintainer | Name, application list |
| `application` | `{app}/application.yaml` | Application team | Flavors, exported services |

## Consequences

### Positive

- Separation of concerns—proxy config rarely changes, workspace defines the environment, application config is owned by the app team
- Each schema has clear ownership
- Application config can be version-controlled with the application
- Workspace config can be shared or kept local

### Negative

- Multiple files to manage
- Configuration loading must merge multiple sources

### Neutral

- Environment variables can override settings from any schema
- CLI flags take highest precedence
