# Shell Integration Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Contrail provides shell integration that enables context-aware Docker Compose interaction through the `contrail-compose` shell function. This function automatically injects the correct project name and compose files based on the current directory, providing a seamless development experience.

**Related Documents**:
- [Context Detection](./context-detection.md)
- [CLI Reference](../reference/cli.md)

**Appendices**:
- [Bash Setup Script](./appendices/shell-integration/bash-setup.sh)
- [Zsh Setup Script](./appendices/shell-integration/zsh-setup.zsh)
- [Fish Setup Script](./appendices/shell-integration/fish-setup.fish)

---

## Behavior

### Architecture

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

### Two Interfaces

1. **`contrail` CLI**: Primary command for managing workspaces, applications, ports, and configuration
2. **`contrail-compose` shell function**: Context-aware passthrough to `docker compose`

This separation exists because:
- Docker Compose has ~20 subcommands with hundreds of flags
- Reimplementing completion would be significant work with no added value
- Users already know `docker compose` — Contrail just makes it workspace-aware
- Shell function + completion delegation achieves full docker compose completion "for free"

---

## Data Schema

### compose-prefix Command

The `contrail compose-prefix` command outputs a docker compose command prefix that includes all necessary flags.

**Flags**:
| Flag | Short | Description |
|------|-------|-------------|
| `--workspace` | `-w` | Specify workspace (overrides context detection) |
| `--app` | `-a` | Specify application (overrides context detection) |

> **Note**: There is no `--flavor` flag. Flavor changes require regeneration and can impact running applications. Use `contrail flavor set` to change the active flavor.

**Output Format**:
- All paths are single-quoted to handle spaces and special characters
- Compose files are included in the correct order (base files first, then generated override)
- Project name follows the `{workspace}-{app}` convention
- Output is a single line suitable for `eval`

**Example Output**:
```bash
$ cd ~/workspaces/dev/app-one
$ contrail compose-prefix
docker compose -p dev-app-one -f '/home/user/workspaces/dev/app-one/docker-compose.yaml' -f '/home/user/workspaces/dev/.generated/app-one.override.yaml'
```

**Exit Codes**:
| Code | Meaning |
|------|---------|
| 0 | Success - prefix output to stdout |
| 1 | General error (configuration issues, file not found, etc.) |
| 5 | Context detection failed (no workspace/app found) |

---

## Examples

### Example 1: Basic Usage

```bash
# From within an application directory (context detected)
$ cd ~/workspaces/dev/app-one
$ contrail-compose exec php bash
$ contrail-compose logs -f
$ contrail-compose up -d

# From workspace root with explicit app
$ cd ~/workspaces/dev
$ contrail-compose -a app-two exec php bash
$ contrail-compose -a app-one logs -f php

# From anywhere with explicit workspace and app
$ contrail-compose -w dev -a app-one ps
```

### Example 2: What contrail-compose Executes

```bash
# These are equivalent:
$ contrail-compose exec php bash

# ...to running:
$ docker compose -p dev-app-one \
    -f ~/workspaces/dev/app-one/docker-compose.yaml \
    -f ~/workspaces/dev/.generated/app-one.override.yaml \
    exec php bash
```

### Example 3: Tab Completion

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
--index     --no-tty        --privileged ...

# Complete with explicit app (contrail flags first)
$ contrail-compose -a app-[TAB]
app-one     app-two     app-three

$ contrail-compose -a app-two exec [TAB]
php     nginx   postgres    worker
```

### Example 4: Daily Development Workflow

```bash
cd ~/workspaces/dev/frontend

# Start your day
contrail up                    # Brings up entire dev workspace

# Work on frontend
contrail-compose logs -f       # Tail frontend logs (context detected)
contrail app restart           # Restart after changes

# Direct Docker Compose interaction
contrail-compose exec php bash          # Shell into container
contrail-compose exec php php artisan   # Run artisan command

# Check on another app
contrail app status -a backend
contrail-compose -a backend logs --tail=50
contrail-compose -a backend exec node npm test

# End of day
contrail down
```

---

## Installation

### Setup Command

```bash
# Bash
contrail init-shell bash >> ~/.bashrc

# Zsh
contrail init-shell zsh >> ~/.zshrc

# Fish
contrail init-shell fish >> ~/.config/fish/conf.d/contrail.fish
```

### What init-shell Provides

- `contrail-compose` function with context-aware Docker Compose passthrough
- Tab completion for `contrail-compose` that delegates to Docker's completion
- Automatic resolution of workspace, app, and compose files

---

## Edge Cases

### Context Resolution Failure

**Scenario**: Running `contrail-compose` outside any workspace.

**Behavior**:
```bash
$ cd ~
$ contrail-compose ps
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail-compose -a APP [command]

Available workspaces: contrail workspace list
```

### Missing Docker Completion

**Scenario**: Docker's shell completion isn't installed.

**Behavior**: Falls back to basic subcommand completion:
```bash
$ contrail-compose [TAB]
build   config  create  down    events  exec    images  kill    logs ...
```

Service names and advanced flags won't complete, but basic operation works.

### Paths with Spaces

**Scenario**: Workspace or application path contains spaces.

**Behavior**: Paths are single-quoted in `compose-prefix` output:
```bash
$ contrail compose-prefix
docker compose -p dev-my-app -f '/Users/name/My Documents/workspaces/dev/my-app/docker-compose.yaml' ...
```

---

## Known Limitations

### All Shells

1. **Docker completion must be installed**: Without docker's shell completion, delegation falls back to basic subcommand completion only

2. **First invocation latency**: First tab completion calls `contrail compose-prefix`, which involves context detection. Subsequent completions in the same directory use shell caching

3. **Path spaces**: Supported but require careful quoting in `compose-prefix` output

### Bash-Specific

1. **Requires bash-completion package**: The `_init_completion` helper function comes from bash-completion

2. **COMP_LINE reconstruction**: Some docker completions use `COMP_LINE` for substring matching. Reconstruction from words may lose original spacing

### Zsh-Specific

1. **1-indexed arrays**: Completion function must account for zsh's 1-indexed arrays

2. **Completion function detection**: Docker's completion might be `_docker`, `_docker-compose`, or loaded lazily

### Fish-Specific

1. **`complete -C` overhead**: Fish's completion delegation spawns a subshell, adding latency

2. **Completion caching**: Fish aggressively caches completions. If compose files change, completions may be stale until shell restarts

---

## Error Handling

| Error Condition | Error Code | Message | Recovery |
|-----------------|------------|---------|----------|
| No workspace found | 5 | `No workspace found in current directory or parents` | Navigate to workspace or use explicit flags |
| No app context | 5 | `No application context detected` | Navigate to app directory or use `-a` flag |
| Docker not available | 4 | `Docker is not installed or not running` | Start Docker |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from shell integration spec |
