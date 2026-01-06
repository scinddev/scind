# Migration Step: Layer 7 - Implementation

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 file + appendices, approximately 1,600 lines

---

## Overview

Create implementation documentation from `specs/contrail-go-stack.md`.

**Target Files**:
1. `implementation/tech-stack.md` - Go technology stack and scaffolding

---

## Implementation: `implementation/tech-stack.md`

**Source**: `specs/contrail-go-stack.md` (entire file ~1,616 lines)

This file contains significant code blocks that exceed the appendix threshold. Create with appendices:

**Main file content**:
- Stack Overview (lines 9-55)
- Dependency Rationale (lines 57-88)
- Architecture Patterns (lines 88-182)
- Project Structure (lines 184-267)
- CLI to Cobra Command Mapping table (lines 1479-1529)
- Testing Strategy overview (lines 1530-1561)
- Implementation Priority (lines 1563-1594)

**Appendices** (code blocks >= 50 lines):
- `implementation/appendices/tech-stack/scaffold-main.go` - Entry point scaffold
- `implementation/appendices/tech-stack/scaffold-app.go` - App commands scaffold
- `implementation/appendices/tech-stack/scaffold-cmd-root.go` - Root command scaffold
- `implementation/appendices/tech-stack/scaffold-workspace.go` - Workspace commands scaffold
- `implementation/appendices/tech-stack/scaffold-aliases.go` - Aliases scaffold
- `implementation/appendices/tech-stack/scaffold-config.go` - Config types scaffold
- `implementation/appendices/tech-stack/scaffold-context.go` - Context detection scaffold
- `implementation/appendices/tech-stack/scaffold-generator.go` - Generator scaffold
- `implementation/appendices/tech-stack/goreleaser.yaml` - GoReleaser config
- `implementation/appendices/tech-stack/makefile` - Makefile

---

## Code Block Mapping

| Source Lines | Content | Target |
|--------------|---------|--------|
| 17-43 | go.mod dependencies | Main file |
| 47-54 | Future dependencies | Main file |
| 94-105 | Context detection pattern | Main file (small) |
| 123-128 | Shell completion pattern | Main file (small) |
| 134-143 | Embed shell scripts | Main file (small) |
| 155-167 | Docker compose pattern | Main file (small) |
| 175-181 | Port check pattern | Main file (small) |
| 306-427 | root.go scaffold | Appendix: scaffold-cmd-root.go |
| 432-572 | workspace.go scaffold | Appendix: scaffold-workspace.go |
| 578-700 | app.go scaffold | Appendix: scaffold-app.go |
| 706-757 | flavor.go scaffold | Appendix: scaffold-flavor.go |
| 763-838 | aliases.go scaffold | Appendix: scaffold-aliases.go |
| 844-898 | compose_prefix.go scaffold | Appendix: scaffold-compose-prefix.go |
| 904-1011 | proxy.go scaffold | Appendix: scaffold-proxy.go |
| 1017-1115 | port.go scaffold | Appendix: scaffold-port.go |
| 1120-1196 | config.go scaffold | Appendix: scaffold-config.go |
| 1200-1249 | validate.go scaffold | Appendix: scaffold-validate.go |
| 1254-1301 | urls.go + open.go scaffold | Appendix: scaffold-utility.go |
| 1305-1363 | init_shell.go scaffold | Appendix: scaffold-init-shell.go |
| 1369-1425 | config/workspace.go types | Appendix: scaffold-config-types.go |
| 1455-1467 | main.go entry point | Appendix: scaffold-main.go |

---

## Also Create: `implementation/README.md`

```markdown
# Implementation Documentation

Implementation guides and developer documentation for Contrail.

## Contents

| Document | Description |
|----------|-------------|
| [tech-stack.md](./tech-stack.md) | Go technology stack and project scaffolding |

## Appendices

Code scaffolds are in `appendices/tech-stack/`:
- Go source file scaffolds for all commands
- Configuration type definitions
- Build configuration files (Makefile, GoReleaser)
```

---

## Appendix Directories to Create

```
implementation/appendices/
  tech-stack/
    scaffold-main.go
    scaffold-cmd-root.go
    scaffold-workspace.go
    scaffold-app.go
    scaffold-flavor.go
    scaffold-aliases.go
    scaffold-compose-prefix.go
    scaffold-proxy.go
    scaffold-port.go
    scaffold-config.go
    scaffold-validate.go
    scaffold-utility.go
    scaffold-init-shell.go
    scaffold-config-types.go
    goreleaser.yaml
    makefile
```

---

## Notes for Migration

1. **Preserve all Go code exactly**: The scaffolding code is meant to be copy-pasted, so preserve formatting and comments
2. **Include package declarations**: Each Go file appendix should include the package declaration
3. **CLI mapping table**: The CLI to Cobra Command Mapping table is valuable reference material - preserve in main file
4. **Testing section**: Include testing strategy and Afero example in main file

---

## Completion Checklist

- [ ] `implementation/tech-stack.md` created
- [ ] `implementation/README.md` created
- [ ] All scaffold appendix files created
- [ ] Complete Go code preserved in appendices
- [ ] CLI mapping table preserved
- [ ] Source attributions present
