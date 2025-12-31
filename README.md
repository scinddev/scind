# Contrail Specification Review

This directory contains the Contrail specification documents and the issues identified during review.

## Directory Structure

```
contrail-spec-review/
├── README.md
├── specs/                          # Original specification documents
│   ├── contrail-prd.md
│   ├── contrail-technical-spec.md
│   ├── contrail-cli-reference.md
│   ├── contrail-shell-integration.md
│   └── contrail-go-stack.md
└── issues/                         # Review findings organized by group
    ├── 00-index.md                 # Start here - overview and recommended order
    ├── 01-go-stack-cleanup.md
    ├── 02-schema-validation.md
    ├── 03-conceptual-foundations.md
    ├── 04-workspace-features.md
    ├── 05-operations-generation.md
    ├── 06-context-detection.md
    ├── 07-cli-commands.md
    ├── 08-shell-integration.md
    ├── 09-proxy-init.md
    ├── 10-documentation-polish.md
    ├── 11-docker-labels.md
    ├── 12-flavor-set-running-app.md
    ├── 13-cli-tech-spec-alignment.md
    ├── 14-defaults-assumptions.md
    ├── 15-error-handling.md
    ├── 16-security-platform.md
    ├── 17-dns-networking.md
    ├── 18-docker-label-consistency.md
    ├── 19-proxy-network-naming.md
    ├── 20-workspace-destroy-registry.md
    ├── 21-go-stack-missing-app-exec.md
    ├── 22-multiple-app-flags.md
    ├── 23-shell-integration-version.md
    ├── 24-traefik-dashboard-port.md
    ├── 25-workspace-prune-go-stack.md
    ├── 26-workspace-list-flags.md
    ├── 27-keep-apps-flag.md
    ├── 28-app-commands-go-stack.md
    ├── 29-port-commands-go-stack.md
    └── 30-flavor-commands-go-stack.md
```

## Getting Started

1. Read `issues/00-index.md` for an overview of all findings
2. Work through issue groups in order (1-30)
3. Each issue file has response sections for your decisions
4. Update specs as you resolve issues

## Summary

- **59 total issues** identified across 30 groups
- **7 high severity** (should resolve before implementation)
- **17 medium severity** (resolve during implementation)
- **35 low severity** (nice-to-haves)

## Recommended Order

| # | File | Focus |
|---|------|-------|
| 1 | Go Stack Cleanup | Quick wins - typos and missing mappings |
| 2 | Schema Validation | Foundational - struct validation alignment |
| 3 | Conceptual Foundations | Core PRD/Tech Spec decisions |
| 4 | Workspace Features | High-severity functionality gaps |
| 5 | Operations & Generation | Lifecycle details |
| 6 | Context Detection | Edge case handling |
| 7 | CLI Commands | Command naming and behavior |
| 8 | Shell Integration | Shell function concerns |
| 9 | Proxy Init | New feature addition |
| 10 | Documentation Polish | Final cleanup |
| 11 | Docker Labels | Docker label schema formalization |
| 12 | Flavor Set Running App | Flavor set behavior with running apps |
| 13 | CLI & Tech Spec Alignment | Command/flag documentation sync |
| 14 | Defaults & Assumptions | Unspecified default values |
| 15 | Error Handling | Edge cases and error scenarios |
| 16 | Security & Platform | Dashboard security, platform scope |
| 17 | DNS & Networking | DNS resolution, naming collisions |
| 18 | Docker Label Consistency | Minor label prefix fix |
| 19 | Proxy Network Naming | Architecture diagram inconsistency |
| 20 | Workspace Destroy Registry | Missing registry step in CLI docs |
| 21 | Go Stack Missing App Exec | Clarify app exec design decision |
| 22 | Multiple App Flags | Repeatable --app flag implementation |
| 23 | Shell Integration Version | Version header mismatch |
| 24 | Traefik Dashboard Port | Dashboard configuration inconsistency |
| 25 | Workspace Prune Go Stack | Missing --dry-run flag |
| 26 | Workspace List Flags | Missing flags in Go Stack |
| 27 | Keep Apps Flag | Consider adding --keep-apps |
| 28 | App Commands Go Stack | Missing app command scaffolding |
| 29 | Port Commands Go Stack | Missing port commands in mapping |
| 30 | Flavor Commands Go Stack | Missing flavor command scaffolding |
