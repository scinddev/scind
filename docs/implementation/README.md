# Implementation Documentation

This directory contains implementation guides and technology stack documentation.

## Documents

| Document | Description |
|----------|-------------|
| [Technology Stack](./tech-stack.md) | Go dependencies, patterns, and scaffolding |

## Appendices

Large content is stored in:
- `appendices/tech-stack/` — Scaffold scripts, full code examples
  - `scaffold-main.go` — Entry point scaffold
  - `scaffold-cmd-root.go` — Root command with context detection
  - `scaffold-config.go` — Configuration type definitions
  - `scaffold-context.go` — Context detection logic
  - `scaffold-generator.go` — Override file generation
  - `scaffold-workspace.go` — Workspace subcommands
  - `scaffold-app.go` — App subcommands
  - `scaffold-aliases.go` — Top-level command aliases
  - `makefile` — Complete Makefile for build/test
  - `goreleaser.yaml` — GoReleaser configuration

## Related Documents

- [Architecture](../architecture/overview.md) — System architecture
- [Specifications](../specs/) — Feature specifications
