# Migration Step: Layer 5 - Reference

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 2 files + appendices, approximately 1,800 lines

---

## Overview

Create reference documentation from `specs/contrail-cli-reference.md` and configuration reference from `specs/contrail-technical-spec.md`.

**Target Files**:
1. `reference/cli.md` - CLI command reference
2. `reference/configuration.md` - Configuration file reference

---

## Reference 1: `reference/cli.md`

**Source**: `specs/contrail-cli-reference.md` (entire file ~1,599 lines)

This is the main CLI reference document. Create with appendices for:
- `reference/appendices/cli/detailed-examples.md` - Example workflows (lines 1491-1578)
- `reference/appendices/cli/error-messages.md` - Error message catalog (lines 1475-1489)

The main file should contain:
- Command Structure (lines 11-27)
- Context Detection (moved to specs/context-detection.md, link to it)
- Resources table (lines 150-163)
- Global Flags (lines 165-212)
- All command documentation:
  - Workspace Commands (lines 216-525)
  - Application Commands (lines 527-721)
  - Flavor Commands (lines 723-810)
  - Port Commands (lines 812-934)
  - Proxy Commands (lines 936-1075)
  - Config Commands (lines 1077-1162)
  - Docker Compose Integration (lines 1164-1247)
  - Top-Level Aliases (lines 1250-1262)
  - Utility Commands (lines 1264-1369)
- Output Formats (lines 1371-1394)
- Shell Completion (lines 1396-1443)
- Environment Variables (lines 1450-1461)
- Exit Codes (lines 1463-1473)

---

## Reference 2: `reference/configuration.md`

**Source**: `specs/contrail-technical-spec.md:216-765` (Configuration Schemas section)

Create a configuration reference that shows:
- All configuration file types and their locations
- Complete schema reference with field descriptions
- Examples for common configurations

Create appendix:
- `reference/appendices/configuration/complete-examples.md` - Full configuration file examples

---

## Also Create: `reference/README.md`

```markdown
# Reference Documentation

API and configuration reference for Contrail.

## Contents

| Document | Description |
|----------|-------------|
| [cli.md](./cli.md) | CLI command reference |
| [configuration.md](./configuration.md) | Configuration file reference |

## Appendices

- `appendices/cli/detailed-examples.md` - Complete workflow examples
- `appendices/cli/error-messages.md` - Error message catalog
- `appendices/configuration/complete-examples.md` - Full configuration examples
```

---

## Appendix Directories to Create

```
reference/appendices/
  cli/
    detailed-examples.md
    error-messages.md
  configuration/
    complete-examples.md
```

---

## Notes for Migration

1. **Context Detection**: The Context Detection section (lines 30-147 of cli-reference.md) should be linked to `specs/context-detection.md` rather than duplicated
2. **Shell Integration**: Shell integration details (lines 1164-1247) should link to `specs/shell-integration.md`
3. **Preserve all tables**: Command flags, options, and output format tables must be preserved exactly
4. **Exit codes**: Include complete exit code table

---

## Completion Checklist

- [ ] `reference/cli.md` created
- [ ] `reference/configuration.md` created
- [ ] `reference/README.md` created
- [ ] All appendix files created
- [ ] Complete command documentation preserved
- [ ] All tables preserved exactly
- [ ] Source attributions present
