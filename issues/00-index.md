# Contrail Specification Issues — Index

**Created**: December 2024  
**Total Issues**: 27  
**Groups**: 10

---

## Recommended Order

| Order | File | Issues | Documents | Effort | Rationale |
|-------|------|--------|-----------|--------|-----------|
| 1 | `01-go-stack-cleanup.md` | 3 | Go Stack | Small | ✅ COMPLETED |
| 2 | `02-schema-validation.md` | 4 | Tech Spec, Go Stack | Medium | ✅ COMPLETED |
| 3 | `03-conceptual-foundations.md` | 3 | PRD, Tech Spec | Small | Establishes core concepts before feature details |
| 4 | `04-workspace-features.md` | 3 | PRD, Tech Spec, CLI Ref | Large | High-severity gaps in core functionality |
| 5 | `05-operations-generation.md` | 5 | Tech Spec, CLI Ref | Medium | Generate/up/down lifecycle details |
| 6 | `06-context-detection.md` | 1 | Tech Spec, CLI Ref, Go Stack | Small | Single focused issue, needs workspace concepts |
| 7 | `07-cli-commands.md` | 3 | CLI Ref, Go Stack | Small | Command naming, depends on concepts being settled |
| 8 | `08-shell-integration.md` | 2 | CLI Ref, Shell Integration, Go Stack | Small | Isolated shell concerns, depends on CLI |
| 9 | `09-proxy-init.md` | 1 | PRD, Tech Spec, CLI Ref, Go Stack | Medium | New feature, best after existing features clarified |
| 10 | `10-documentation-polish.md` | 2 | All | Small | Final cleanup after substantive changes |

---

## By Severity

### High Severity (7 issues)
| Issue | Group | File |
|-------|-------|------|
| C-2: Port validation mismatch | 2 | `02-schema-validation.md` |
| C-3: Service validation mismatch | 2 | `02-schema-validation.md` |
| C-5: `-f` flag collision | 8 | `08-shell-integration.md` |
| M-1: Workspace discovery missing | 4 | `04-workspace-features.md` |
| M-7: Workspace name collisions | 4 | `04-workspace-features.md` |
| A-1: Staleness detection undefined | 5 | `05-operations-generation.md` |

### Medium Severity (10 issues)
| Issue | Group | File |
|-------|-------|------|
| C-4: Proxy command naming | 7 | `07-cli-commands.md` |
| M-2: No proxy init command | 9 | `09-proxy-init.md` |
| M-6: Single-app clone behavior | 4 | `04-workspace-features.md` |
| M-8: Port release detection | 5 | `05-operations-generation.md` |
| M-9: Manual override behavior | 5 | `05-operations-generation.md` |
| A-2: Network creation timing | 3 | `03-conceptual-foundations.md` |
| A-3: Nested config precedence | 6 | `06-context-detection.md` |
| A-4: Template resolution timing | 5 | `05-operations-generation.md` |
| A-6: Repeatable --app behavior | 7 | `07-cli-commands.md` |
| A-10: Compose file validation | 2 | `02-schema-validation.md` |

### Low Severity (10 issues)
| Issue | Group | File |
|-------|-------|------|
| C-1: Version desync | 10 | `10-documentation-polish.md` |
| M-3: destroy in Tech Spec | 5 | `05-operations-generation.md` |
| M-4: app exec mapping | 1 | `01-go-stack-cleanup.md` |
| M-5: flavor reset mapping | 1 | `01-go-stack-cleanup.md` |
| A-5: Base env var rationale | 3 | `03-conceptual-foundations.md` |
| A-7: Visibility field purpose | 3 | `03-conceptual-foundations.md` |
| A-8: DNS validation | 7 | `07-cli-commands.md` |
| A-9: Router name collisions | 2 | `02-schema-validation.md` |
| X-1: BoolVarP typo | 1 | `01-go-stack-cleanup.md` |
| X-2: --since examples | 10 | `10-documentation-polish.md` |
| X-3: Exit code reference | 8 | `08-shell-integration.md` |

---

## By Document Impact

### Go Stack Only
- `01-go-stack-cleanup.md` (3 issues)

### Technical Spec + Go Stack
- `02-schema-validation.md` (4 issues)

### PRD + Technical Spec
- `03-conceptual-foundations.md` (3 issues)

### PRD + Technical Spec + CLI Reference
- `04-workspace-features.md` (3 issues)

### Technical Spec + CLI Reference
- `05-operations-generation.md` (5 issues)

### Technical Spec + CLI Reference + Go Stack
- `06-context-detection.md` (1 issue)

### CLI Reference + Go Stack
- `07-cli-commands.md` (3 issues)

### CLI Reference + Shell Integration + Go Stack
- `08-shell-integration.md` (2 issues)

### All Four Main Docs
- `09-proxy-init.md` (1 issue)

### Various
- `10-documentation-polish.md` (2 issues)

---

## Quick Start

1. **Start with Group 1** (`01-go-stack-cleanup.md`) — 3 quick fixes, builds momentum
2. **Then Group 2** (`02-schema-validation.md`) — foundational decisions
3. **Groups 3-5** form the conceptual core — work through in order
4. **Groups 6-8** are implementation details — can parallelize if needed
5. **Group 9** is additive (new feature) — can defer if time-constrained
6. **Group 10** is final polish — do last

---

## Files in This Directory

```
contrail-issues/
├── 00-index.md                    # This file
├── 01-go-stack-cleanup.md         # Quick wins
├── 02-schema-validation.md        # Foundational schemas
├── 03-conceptual-foundations.md   # Core concepts
├── 04-workspace-features.md       # Major functionality gaps
├── 05-operations-generation.md    # Lifecycle operations
├── 06-context-detection.md        # Edge case handling
├── 07-cli-commands.md             # Command details
├── 08-shell-integration.md        # Shell concerns
├── 09-proxy-init.md               # New feature
└── 10-documentation-polish.md     # Final cleanup
```
