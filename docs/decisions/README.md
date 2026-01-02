# Architectural Decision Records

This directory contains Architectural Decision Records (ADRs) documenting significant technical and product decisions.

## What is an ADR?

An ADR captures the **why** behind important decisions. Each ADR is a simple file describing:
- The context that motivated the decision
- The decision that was made
- The consequences (positive, negative, neutral)

## ADR Index

| Number | Title | Status | Date |
|--------|-------|--------|------|
| [0000](./0000-template.md) | Template | Template | - |
| [0001](./0001-docker-compose-project-name-isolation.md) | Docker Compose Project Name Isolation | Accepted | 2024-12 |
| [0002](./0002-two-layer-networking.md) | Two-Layer Networking | Accepted | 2024-12 |
| [0003](./0003-pure-overlay-design.md) | Pure Overlay Design | Accepted | 2024-12 |
| [0004](./0004-convention-based-naming.md) | Convention-Based Naming | Accepted | 2024-12 |
| [0005](./0005-structure-vs-state-separation.md) | Structure vs State Separation | Accepted | 2024-12 |
| [0006](./0006-three-configuration-schemas.md) | Three Configuration Schemas | Accepted | 2024-12 |
| [0007](./0007-port-type-system.md) | Port Type System for Exported Services | Accepted | 2024-12 |
| [0008](./0008-traefik-reverse-proxy.md) | Traefik for Reverse Proxy | Accepted | 2024-12 |
| [0009](./0009-flexible-tls-configuration.md) | Flexible TLS Configuration | Accepted | 2024-12 |
| [0010](./0010-up-down-command-semantics.md) | up/down Command Semantics | Accepted | 2024-12 |
| [0011](./0011-options-based-targeting.md) | Options-Based Targeting with Context Detection | Accepted | 2024-12 |

*ADRs are numbered sequentially. Use `0000-template.md` as a starting point.*

## Creating a New ADR

1. Copy `0000-template.md` to `NNNN-{short-title}.md`
2. Fill in the context, decision, and consequences
3. Update this index

## Related Documents

- [Vision](../product/vision.md) — Product vision and goals
- [Architecture](../architecture/overview.md) — System architecture overview
