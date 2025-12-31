# Refinement Issues — Index

**Created**: December 2024
**Total Issues**: 11
**Groups**: 6

---

## Recommended Order

| Order | File | Issues | Documents | Effort | Rationale |
|-------|------|--------|-----------|--------|-----------|
| 1 | 01-workspace-destroy.md | M-1 | PRD | Small | ✅ COMPLETED |
| 2 | 02-proxy-docker-compose.md | C-1 | Technical Spec | Small | ✅ COMPLETED |
| 3 | 03-environment-variable-pattern.md | A-1 | PRD, Technical Spec | Small | ✅ COMPLETED |
| 4 | 04-tls-support.md | A-2 | PRD, Technical Spec | Medium | ✅ COMPLETED |
| 5 | 05-go-stack-missing-commands.md | M-2, M-3, M-4, M-5, C-2 | CLI Reference, Go Stack | Medium | ✅ COMPLETED |
| 6 | 06-template-variables.md | A-3, A-4 | PRD, Technical Spec | Small | ✅ COMPLETED |

---

## By Severity

### High Severity (0 issues)
| Issue | Group | File |
|-------|-------|------|

### Medium Severity (4 issues)
| Issue | Group | File |
|-------|-------|------|
| C-1 | 2 | 02-proxy-docker-compose.md |
| A-1 | 3 | 03-environment-variable-pattern.md |
| A-2 | 4 | 04-tls-support.md |
| C-2 | 5 | 05-go-stack-missing-commands.md |

### Low Severity (7 issues)
| Issue | Group | File |
|-------|-------|------|
| M-1 | 1 | 01-workspace-destroy.md |
| M-2 | 5 | 05-go-stack-missing-commands.md |
| M-3 | 5 | 05-go-stack-missing-commands.md |
| M-4 | 5 | 05-go-stack-missing-commands.md |
| M-5 | 5 | 05-go-stack-missing-commands.md |
| A-3 | 6 | 06-template-variables.md |
| A-4 | 6 | 06-template-variables.md |

---

## By Document Impact

### PRD (contrail-prd.md)
- M-1: Workspace destroy missing from Quick Reference
- A-1: Environment variable hyphen conversion rule
- A-2: TLS/HTTPS configuration clarity
- A-3: Template syntax mismatch with Technical Spec

### Technical Spec (contrail-technical-spec.md)
- C-1: Proxy network name inconsistency in docker-compose example
- A-1: Environment variable hyphen conversion rule
- A-2: TLS/HTTPS configuration clarity
- A-3: Template syntax documentation
- A-4: SERVICE_PORT variable use case

### CLI Reference (contrail-cli-reference.md)
- M-2, M-3, M-4, M-5: Commands defined but not scaffolded in Go Stack

### Go Stack (contrail-go-stack.md)
- M-2: Missing port commands scaffold
- M-3: Missing config commands scaffold
- M-4: Missing utility commands scaffolds
- M-5: Missing init-shell command scaffold
- C-2: Alias commands reference pattern issue

### Shell Integration (contrail-shell-integration.md)
_No issues identified._

---

## Quick Start

1. Run refinement research to identify new issues
2. Work through issues in recommended order
3. Archive when complete

---

## Files in This Directory

```
issues/
├── 00-index.md                    # This file
├── 01-workspace-destroy.md        # PRD Quick Reference completeness
├── 02-proxy-docker-compose.md     # Technical Spec internal fix
├── 03-environment-variable-pattern.md  # Cross-doc naming convention
├── 04-tls-support.md              # TLS/HTTPS implementation status
├── 05-go-stack-missing-commands.md     # Go Stack scaffolding
└── 06-template-variables.md       # Template documentation
```

---

## Archived

This issue was archived on 2024-12-31 at 11:51:01.
