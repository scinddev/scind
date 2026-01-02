# ADR-0006: Three Configuration Schemas

**Status**: Accepted
**Date**: December 2024
**Decision-Makers**: Contrail Core Team

---

## Context

Configuration could be in one monolithic file or separated by concern and ownership. Different aspects of the system change at different rates and are owned by different roles.

## Decision

Implement three schema types with clear separation of concerns:

| Schema | Location | Owner | Purpose |
|--------|----------|-------|---------|
| `proxy.yaml` | `~/.config/contrail/` | User | Global proxy settings (domain, TLS, dashboard) |
| `workspace.yaml` | Workspace root | Workspace owner | Workspace definition and application list |
| `application.yaml` | Application directory | Application team | Flavors and exported services |

## Consequences

### Positive

- Proxy config rarely changes; can be set once per machine
- Workspace config defines the environment structure
- Application config is owned by the application team and travels with the app
- Clear ownership boundaries reduce coordination overhead
- Each file can have its own change frequency

### Negative

- Users must understand three different configuration files
- Some settings require knowing which file to modify

### Neutral

- Configuration files can be combined when using single-app workspaces

---

## Notes

The three-schema approach follows the principle of separating concerns by rate of change and ownership. Proxy config changes when migrating machines, workspace config changes when restructuring environments, and application config changes when the application's exported services change.
