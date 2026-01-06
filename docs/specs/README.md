# Specifications

This directory contains detailed feature specifications. Each specification describes behavior, data schemas, examples, and edge cases for specific Contrail features.

## Specifications

| Specification | Status | Description |
|---------------|--------|-------------|
| [Configuration Schemas](./configuration-schemas.md) | Accepted | Configuration file schemas (proxy, workspace, application) |
| [Naming Conventions](./naming-conventions.md) | Accepted | Naming patterns for hostnames, aliases, and projects |
| [Port Types](./port-types.md) | Accepted | Port type system (proxied vs assigned) |
| [Proxy Infrastructure](./proxy-infrastructure.md) | Accepted | Traefik proxy configuration and network |
| [Context Detection](./context-detection.md) | Accepted | Directory context detection algorithm |
| [Workspace Lifecycle](./workspace-lifecycle.md) | Accepted | Workspace operations (up/down/generate) |
| [Generated Override Files](./generated-override-files.md) | Accepted | Override file generation specification |
| [Docker Labels](./docker-labels.md) | Accepted | Docker label conventions |
| [Environment Variables](./environment-variables.md) | Accepted | Environment variable injection |
| [Shell Integration](./shell-integration.md) | Accepted | Shell integration and contrail-compose |

## Specification Status

- **Draft**: Initial design, subject to significant changes
- **Review**: Under review, feedback requested
- **Accepted**: Approved for implementation
- **Implemented**: Fully implemented in code

## Creating New Specifications

Use `_template.md` as a starting point for new specifications.

## Appendices

Large content (code blocks >50 lines, tables >20 rows) is stored in:
- `appendices/{spec-name}/` directories

Available appendices:
- [Configuration Schemas Examples](./appendices/configuration-schemas/complete-examples.md)
- [Generated Override Example](./appendices/generated-override-files/complete-override-example.yaml)
- [Traefik Compose](./appendices/proxy-infrastructure/traefik-compose.yaml)
- [Traefik Config](./appendices/proxy-infrastructure/traefik-config.yaml)
- [Bash Setup](./appendices/shell-integration/bash-setup.sh)
- [Zsh Setup](./appendices/shell-integration/zsh-setup.zsh)
- [Fish Setup](./appendices/shell-integration/fish-setup.fish)

## Related Documents

- [Architecture](../architecture/README.md) — System architecture
- [Decisions](../decisions/) — Decision rationale
- [Reference](../reference/) — CLI and configuration reference
