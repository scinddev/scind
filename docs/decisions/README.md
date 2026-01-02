# Decisions (ADRs)

This directory contains Architectural Decision Records (ADRs) for the Contrail project. ADRs capture the *why* behind significant technical and product decisions.

## Index

| Number | Title | Status | Date |
|--------|-------|--------|------|
| [0000](./0000-template/README.md) | ADR Template | Template | - |
| [0001](./0001-docker-compose-project-name-isolation/README.md) | Docker Compose Project Name Isolation | Accepted | Dec 2024 |
| [0002](./0002-two-layer-networking/README.md) | Two-Layer Networking | Accepted | Dec 2024 |
| [0003](./0003-pure-overlay-design/README.md) | Pure Overlay Design | Accepted | Dec 2024 |
| [0004](./0004-convention-based-naming/README.md) | Convention-Based Naming | Accepted | Dec 2024 |
| [0005](./0005-structure-vs-state-separation/README.md) | Structure vs State Separation | Accepted | Dec 2024 |
| [0006](./0006-three-configuration-schemas/README.md) | Three Configuration Schemas | Accepted | Dec 2024 |
| [0007](./0007-port-type-system/README.md) | Port Type System for Exported Services | Accepted | Dec 2024 |
| [0008](./0008-traefik-reverse-proxy/README.md) | Traefik for Reverse Proxy | Accepted | Dec 2024 |
| [0009](./0009-flexible-tls-configuration/README.md) | Flexible TLS Configuration | Accepted | Dec 2024 |
| [0010](./0010-up-down-command-semantics/README.md) | `up`/`down` Command Semantics | Accepted | Dec 2024 |
| [0011](./0011-options-based-targeting/README.md) | Options-Based Targeting with Context Detection | Accepted | Dec 2024 |
| [0012](./0012-layered-documentation-system/README.md) | Layered Documentation System | Accepted | Dec 2024 |

## About ADRs

ADRs are immutable once accepted. If a decision needs to change, a new ADR supersedes the old one.

### When to Create an ADR

- Choosing between multiple viable technical approaches
- Making a decision that would be expensive to reverse
- Establishing a pattern that will be followed throughout the codebase
- Deviating from common practice or industry norms

### ADR Lifecycle

```
Draft → Proposed → Accepted → [Superseded by NNNN]
                            → [Deprecated]
```
