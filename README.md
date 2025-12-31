# Contrail Specification Review

This directory contains the Contrail specification documents and the issues identified during review.

## Directory Structure

```
contrail-spec-review/
├── README.md
├── specs/                          # Original specifications and documents
│   ├── contrail-prd.md
│   ├── contrail-technical-spec.md
│   ├── contrail-cli-reference.md
│   ├── contrail-shell-integration.md
│   └── contrail-go-stack.md
├── issues/                         # Review findings organized by group
│   ├── 00-index.md                 # Start here - overview and recommended order
│   ├── 01-workspace-destroy.md     # Missing workspace destroy in PRD
│   ├── 02-proxy-docker-compose.md  # Network name inconsistency
│   ├── 03-environment-variable-pattern.md  # Hyphen conversion rule
│   ├── 04-tls-support.md           # TLS/HTTPS configuration gap
│   ├── 05-go-stack-missing-commands.md    # Missing command scaffolds
│   └── 06-template-variables.md    # Template syntax and variables
└── archive/                        # Archived issue sets
    └── 20251231100330/             # Previous review cycle (59 issues, 30 groups)
```

## Getting Started

1. Run "research refinement" to identify new issues
2. Work through issues in recommended order
3. Archive when complete

## Summary

- **11 total issues** currently tracked
- **6 groups** pending review

## Recommended Order

| # | File | Focus |
|---|------|-------|
| 1 | 01-workspace-destroy.md | PRD Quick Reference completeness |
| 2 | 02-proxy-docker-compose.md | Technical Spec internal consistency |
| 3 | 03-environment-variable-pattern.md | Cross-doc naming convention clarity |
| 4 | 04-tls-support.md | TLS/HTTPS implementation status |
| 5 | 05-go-stack-missing-commands.md | Go Stack scaffolding completeness |
| 6 | 06-template-variables.md | Template documentation clarity |
