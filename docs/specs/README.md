# Specifications

This directory contains precise behavioral definitions for Contrail features.

## Index

| Specification | Description |
|--------------|-------------|
| [Configuration Schemas](./configuration-schemas/README.md) | Proxy, workspace, and application configuration |
| [Context Detection](./context-detection/README.md) | Automatic workspace/app detection algorithm |
| [Docker Labels](./docker-labels/README.md) | Label conventions for container metadata |
| [Environment Variables](./environment-variables/README.md) | Service discovery via CONTRAIL_* variables |
| [Generated Override Files](./generated-override-files/README.md) | Docker Compose override generation |
| [Naming Conventions](./naming-conventions/README.md) | Hostname, alias, and project name patterns |
| [Port Types](./port-types/README.md) | Proxied vs assigned port routing |
| [Proxy Infrastructure](./proxy-infrastructure/README.md) | Traefik proxy setup and TLS |
| [Shell Integration](./shell-integration/README.md) | contrail-compose and shell functions |
| [Workspace Lifecycle](./workspace-lifecycle/README.md) | Up, down, generate, destroy operations |

## About Specifications

Specifications define **how** features work in precise detail. They are the source of truth for implementation.

### When to Create a Spec

- Defining behavior that spans multiple components
- Documenting protocols, formats, or conventions
- Specifying validation rules or error handling
- Capturing complex algorithms or state machines

### Spec Structure

Each specification should include:
- Clear problem/context statement
- Detailed behavioral definition
- Examples and edge cases
- Related ADRs and other specs
