# Specification: Shell Integration

**Version**: 0.1.0
**Date**: December 2024

<!-- Migrated from specs/contrail-shell-integration.md -->

---

## Overview

Contrail provides two complementary shell interfaces:

1. **`contrail` CLI**: The primary command for managing workspaces, applications, ports, and configuration
2. **`contrail-compose` shell function**: A context-aware passthrough to `docker compose` that automatically injects the correct project name and compose files

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
│         │                    │ (outputs docker compose  │           │
│         │                    │  prefix with -p, -f, etc)│           │
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
│                                                                     │
│   Completion:                Completion:                            │
│   Standard CLI framework     Delegated to docker's completion       │
│   (Cobra, Click, etc.)       with transformed context               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## The `contrail compose-prefix` Command

This command outputs a docker compose command prefix that includes all necessary flags.

### Behavior

```bash
$ cd ~/workspaces/dev/app-one
$ contrail compose-prefix
docker compose -p dev-app-one -f '/home/user/workspaces/dev/app-one/docker-compose.yaml' -f '/home/user/workspaces/dev/.generated/app-one.override.yaml'
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Specify workspace (overrides context detection) |
| `--app` | `-a` | Specify application (overrides context detection) |

**Note**: There is no `--flavor` flag. Flavor changes require regeneration via `contrail flavor set`.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success - prefix output to stdout |
| 1 | General error |
| 5 | Context detection failed |

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

This provides:
- The `contrail-compose` shell function
- Tab completion for `contrail-compose` that delegates to Docker's completion
- Standard `contrail` CLI completions

---

## Usage Examples

### Basic Usage

```bash
# From within an application directory (context detected)
$ cd ~/workspaces/dev/app-one
$ contrail-compose exec php bash
$ contrail-compose logs -f
$ contrail-compose up -d

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
build   config  create  down    events  exec    images  kill    logs ...

# Complete service names (from the resolved compose file)
$ contrail-compose exec [TAB]
php     nginx   mysql   redis

# Complete docker compose flags
$ contrail-compose exec -[TAB]
-d          --detach        -e          --env

# Complete with explicit app
$ contrail-compose -a app-[TAB]
app-one     app-two     app-three
```

---

## Known Limitations

### All Shells

1. **Docker completion must be installed**: Falls back to basic subcommand completion without it
2. **First invocation latency**: First tab completion calls `contrail compose-prefix`
3. **Path spaces**: Paths with spaces are supported but require careful quoting

### Bash-Specific

- Requires bash-completion package

### Zsh-Specific

- 1-indexed arrays affect completion position calculation
- Docker's completion might be `_docker`, `_docker-compose`, or loaded lazily

### Fish-Specific

- `complete -C` overhead from subshell spawning
- Aggressive completion caching may show stale results

---

## Related Documents

- [Context Detection](context-detection.md)
- [CLI Reference](../reference/cli.md)
