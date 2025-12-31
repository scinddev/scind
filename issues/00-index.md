# Contrail Specification Issues — Index

**Created**: December 2024
**Total Issues**: 47
**Groups**: 18

---

## Recommended Order

| Order | File | Issues | Documents | Effort | Rationale |
|-------|------|--------|-----------|--------|-----------|
| 1 | `01-go-stack-cleanup.md` | 3 | Go Stack | Small | ✅ COMPLETED |
| 2 | `02-schema-validation.md` | 4 | Tech Spec, Go Stack | Medium | ✅ COMPLETED |
| 3 | `03-conceptual-foundations.md` | 3 | PRD, Tech Spec | Small | ✅ COMPLETED |
| 4 | `04-workspace-features.md` | 3 | PRD, Tech Spec, CLI Ref | Large | ✅ COMPLETED |
| 5 | `05-operations-generation.md` | 5 | Tech Spec, CLI Ref | Medium | ✅ COMPLETED |
| 6 | `06-context-detection.md` | 1 | Tech Spec, CLI Ref, Go Stack | Small | ✅ COMPLETED |
| 7 | `07-cli-commands.md` | 3 | CLI Ref, Go Stack | Small | ✅ COMPLETED |
| 8 | `08-shell-integration.md` | 2 | CLI Ref, Shell Integration, Go Stack | Small | ✅ COMPLETED |
| 9 | `09-proxy-init.md` | 1 | PRD, Tech Spec, CLI Ref, Go Stack | Medium | ✅ COMPLETED |
| 10 | `10-documentation-polish.md` | 2 | All | Small | ✅ COMPLETED |
| 11 | `11-docker-labels.md` | 1 | Tech Spec, Go Stack | Medium | ✅ COMPLETED |
| 12 | `12-flavor-set-running-app.md` | 1 | Tech Spec, CLI Ref | Small | ✅ COMPLETED |
| 13 | `13-cli-tech-spec-alignment.md` | 3 | Tech Spec, CLI Ref | Small | ✅ COMPLETED |
| 14 | `14-defaults-assumptions.md` | 4 | PRD, Tech Spec | Small | ✅ COMPLETED |
| 15 | `15-error-handling.md` | 5 | Tech Spec, CLI Ref, Go Stack | Medium | ✅ COMPLETED |
| 16 | `16-security-platform.md` | 2 | PRD, Tech Spec | Small | ✅ COMPLETED |
| 17 | `17-dns-networking.md` | 3 | Tech Spec, CLI Ref | Small | ✅ COMPLETED |
| 18 | `18-docker-label-consistency.md` | 1 | Tech Spec | Small | ✅ COMPLETED |

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

### Medium Severity (15 issues)
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
| A-11: Flavor set running app | 12 | `12-flavor-set-running-app.md` |
| L-1: Docker labels not specified | 11 | `11-docker-labels.md` |
| N-7: Single-app naming ambiguity | 14 | `14-defaults-assumptions.md` |
| N-8: Concurrent operations undefined | 15 | `15-error-handling.md` |
| N-11: Non-existent Compose service | 15 | `15-error-handling.md` |
| N-13: Traefik dashboard security | 16 | `16-security-platform.md` |

### Low Severity (25 issues)
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
| N-1: --keep-apps missing from CLI | 13 | `13-cli-tech-spec-alignment.md` |
| N-2: config edit missing from CLI | 13 | `13-cli-tech-spec-alignment.md` |
| N-3: proxy logs status | 13 | `13-cli-tech-spec-alignment.md` |
| N-4: Default visibility unspecified | 14 | `14-defaults-assumptions.md` |
| N-5: Default protocol unspecified | 14 | `14-defaults-assumptions.md` |
| N-6: Traefik version pinning | 14 | `14-defaults-assumptions.md` |
| N-9: Git clone failure behavior | 15 | `15-error-handling.md` |
| N-10: Docker not installed | 15 | `15-error-handling.md` |
| N-12: Orphaned port cleanup timing | 15 | `15-error-handling.md` |
| N-14: Platform support scope | 16 | `16-security-platform.md` |
| N-15: DNS resolution in doctor | 17 | `17-dns-networking.md` |
| N-16: Volume naming collision | 17 | `17-dns-networking.md` |
| N-17: Existing proxy network | 17 | `17-dns-networking.md` |
| N-18: Docker label prefix inconsistency | 18 | `18-docker-label-consistency.md` |

---

## By Document Impact

### Go Stack Only
- `01-go-stack-cleanup.md` (3 issues)

### Technical Spec + Go Stack
- `02-schema-validation.md` (4 issues)
- `11-docker-labels.md` (1 issue)

### PRD + Technical Spec
- `03-conceptual-foundations.md` (3 issues)
- `14-defaults-assumptions.md` (4 issues)
- `16-security-platform.md` (2 issues)

### PRD + Technical Spec + CLI Reference
- `04-workspace-features.md` (3 issues)

### Technical Spec + CLI Reference
- `05-operations-generation.md` (5 issues)
- `12-flavor-set-running-app.md` (1 issue)
- `13-cli-tech-spec-alignment.md` (3 issues)
- `17-dns-networking.md` (3 issues)

### Technical Spec + CLI Reference + Go Stack
- `06-context-detection.md` (1 issue)
- `15-error-handling.md` (5 issues)

### CLI Reference + Go Stack
- `07-cli-commands.md` (3 issues)

### CLI Reference + Shell Integration + Go Stack
- `08-shell-integration.md` (2 issues)

### All Four Main Docs
- `09-proxy-init.md` (1 issue)

### Various
- `10-documentation-polish.md` (2 issues)

### Technical Spec Only
- `18-docker-label-consistency.md` (1 issue)

---

## Quick Start

1. **Start with Group 1** (`01-go-stack-cleanup.md`) — 3 quick fixes, builds momentum
2. **Then Group 2** (`02-schema-validation.md`) — foundational decisions
3. **Groups 3-5** form the conceptual core — work through in order
4. **Groups 6-8** are implementation details — can parallelize if needed
5. **Group 9** is additive (new feature) — can defer if time-constrained
6. **Group 10** is final polish — do last
7. **Group 11** formalizes Docker labels — added during review
8. **Group 12** documents flavor set behavior — added during Group 8
9. **Groups 13-18** are new issues from second review pass:
   - **Group 13**: Quick alignment fixes between CLI and Tech Spec
   - **Group 14**: Clarify default values and assumptions
   - **Group 15**: Document error handling and edge cases
   - **Group 16**: Security and platform scope decisions
   - **Group 17**: DNS and networking details
   - **Group 18**: Single fix for label prefix consistency

---

## Files in This Directory

```
issues/
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
├── 10-documentation-polish.md     # Final cleanup
├── 11-docker-labels.md            # Docker label schema
├── 12-flavor-set-running-app.md   # Flavor set behavior
├── 13-cli-tech-spec-alignment.md  # CLI/Tech Spec sync
├── 14-defaults-assumptions.md     # Default values
├── 15-error-handling.md           # Error scenarios
├── 16-security-platform.md        # Security & platform
├── 17-dns-networking.md           # DNS & networking
└── 18-docker-label-consistency.md # Label prefix fix
```
