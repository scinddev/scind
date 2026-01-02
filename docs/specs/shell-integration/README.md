# Shell Integration Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail provides two complementary shell interfaces:

1. **`contrail` CLI**: Primary command for managing workspaces, applications, ports, configuration
2. **`contrail-compose` shell function**: Context-aware passthrough to `docker compose`

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User's Shell                                │
│                                                                     │
│   contrail workspace up          contrail-compose exec php bash     │
│         │                                    │                      │
│         ▼                                    ▼                      │
│   ┌───────────┐                    ┌─────────────────┐              │
│   │ contrail  │                    │ contrail-compose│              │
│   │  binary   │                    │ shell function  │              │
│   └───────────┘                    └────────┬────────┘              │
│         │                                   │                       │
│         │                                   ▼                       │
│         │                    ┌──────────────────────────┐           │
│         │                    │ contrail compose-prefix  │           │
│         │                    └──────────────┬───────────┘           │
│         │                                   │                       │
│         │                                   ▼                       │
│         │                    ┌──────────────────────────┐           │
│         │                    │     docker compose       │           │
│         │                    │  -p dev-app-one          │           │
│         │                    │  -f docker-compose.yaml  │           │
│         │                    │  -f override.yaml        │           │
│         │                    │  exec php bash           │           │
│         │                    └──────────────────────────┘           │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Installation

```bash
# Bash
contrail init-shell bash >> ~/.bashrc

# Zsh
contrail init-shell zsh >> ~/.zshrc

# Fish
contrail init-shell fish >> ~/.config/fish/conf.d/contrail.fish
```

---

## The `contrail compose-prefix` Command

Outputs a docker compose command prefix with project name and compose files:

```bash
$ cd ~/workspaces/dev/app-one
$ contrail compose-prefix
docker compose -p dev-app-one -f '/path/to/docker-compose.yaml' -f '/path/to/override.yaml'
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Override workspace context |
| `--app` | `-a` | Override application context |

**Note**: No `--flavor` flag. Flavor changes require regeneration via `contrail flavor set`.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 5 | Context detection failed |

---

## Usage Examples

### Basic Usage

```bash
# From application directory (context detected)
$ cd ~/workspaces/dev/app-one
$ contrail-compose exec php bash
$ contrail-compose logs -f

# From workspace root with explicit app
$ cd ~/workspaces/dev
$ contrail-compose -a app-two exec php bash

# From anywhere with explicit workspace and app
$ contrail-compose -w dev -a app-one ps
```

### Tab Completion

```bash
# Complete subcommands
$ contrail-compose [TAB]
build   config  create  down    events  exec ...

# Complete service names
$ contrail-compose exec [TAB]
php     nginx   mysql   redis

# Complete with explicit app
$ contrail-compose -a app-[TAB]
app-one     app-two     app-three
```

---

## Provided Functions

The shell integration provides:

- `contrail-compose` function with context-aware Docker Compose passthrough
- Tab completion delegated to Docker's own completion
- Automatic resolution of workspace, app, and compose files

---

## Known Limitations

1. **Docker completion must be installed** for full completion delegation
2. **First invocation latency** from context detection
3. **Paths with spaces** require careful quoting in `compose-prefix` output

---

## Related Documentation

- [CLI Reference](../../reference/cli/README.md)
- [Context Detection Spec](../context-detection/README.md)
- [Implementation: Go Stack](../../implementation/go-stack/README.md)
