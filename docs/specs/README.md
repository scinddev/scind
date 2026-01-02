# Specifications

This directory contains detailed feature specifications. Each specification describes behavior, data schemas, examples, and edge cases for specific Contrail features.

## Specifications

| Specification | Status | Description |
|---------------|--------|-------------|
| [Configuration Schemas](./configuration-schemas.md) | Accepted | workspace.yaml, application.yaml, proxy.yaml schemas |
| [Context Detection](./context-detection.md) | Accepted | Directory walking and context resolution |
| [Docker Labels](./docker-labels.md) | Accepted | Traefik labels and workspace labels |
| [Environment Variables](./environment-variables.md) | Accepted | Environment variable mapping and overrides |
| [Generated Override Files](./generated-override-files.md) | Accepted | Override file generation specification |
| [Naming Conventions](./naming-conventions.md) | Accepted | Project names, hostnames, aliases |
| [Port Types](./port-types.md) | Accepted | Proxied and assigned port handling |
| [Proxy Infrastructure](./proxy-infrastructure.md) | Accepted | Traefik configuration and network |
| [Shell Integration](./shell-integration.md) | Accepted | contrail-compose shell function |
| [Workspace Lifecycle](./workspace-lifecycle.md) | Accepted | up/down/restart semantics |

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

- [Architecture](../architecture/overview.md) — System architecture
- [Decisions](../decisions/) — Decision rationale
- [Reference](../reference/) — CLI and configuration reference
