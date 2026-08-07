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

## Scope Note

This ADR covers **configuration schemas** that define structure:
- `proxy.yaml` - Global/per-user configuration
- `workspace.yaml` - Per-workspace configuration
- `application.yaml` - Per-application configuration

**State files** are not configuration schemas and are documented separately:
- `state.yaml` (global and per-workspace) - See [State Management Spec](../specs/state-management.md)
- `workspaces.yaml` - See [Workspace Registry Spec](../specs/state-management.md#workspace-registry)
- `manifest.yaml` - See [Generated Manifest Spec](../specs/generated-manifest.md)
