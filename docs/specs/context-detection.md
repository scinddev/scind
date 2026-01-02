# Context Detection Specification

**Version**: 1.0.0
**Date**: 2026-01-02
**Status**: Accepted

---

## Overview

Context detection allows Contrail commands to automatically infer the target workspace and application from the current working directory. This reduces the need for explicit `--workspace` and `--app` flags in most cases, making the CLI more ergonomic for daily development workflows.

The algorithm uses a **workspace boundary** approach to prevent accidental detection of configuration files in vendor packages or nested test fixtures.

**Related Documents**:
- [ADR-0011: Options-Based Targeting](../decisions/0011-options-based-targeting.md)
- [Architecture: Overview](../architecture/overview.md)
- [CLI Reference](../reference/cli.md)

---

## Behavior

### Detection Algorithm

1. **Find workspace root**: Walk up from current working directory looking for `workspace.yaml`. The first one found establishes the workspace root.

2. **Find application context**: Walk up from current working directory toward the workspace root looking for `application.yaml`. Only consider `application.yaml` files that are within the workspace directory tree.

3. **Never traverse above workspace root** for application detection—this prevents vendor packages or nested fixtures from hijacking context.

### Context Resolution Flow

```
Current Directory: ~/workspaces/dev/app-one/src/components/

Step 1: Walk up looking for workspace.yaml
  ~/workspaces/dev/app-one/src/components/ - not found
  ~/workspaces/dev/app-one/src/ - not found
  ~/workspaces/dev/app-one/ - not found
  ~/workspaces/dev/ - FOUND workspace.yaml

  Result: Workspace root = ~/workspaces/dev/

Step 2: Walk up from current directory toward workspace root looking for application.yaml
  ~/workspaces/dev/app-one/src/components/ - not found
  ~/workspaces/dev/app-one/src/ - not found
  ~/workspaces/dev/app-one/ - FOUND application.yaml (within workspace tree)

  Result: Application = app-one (from directory name)

Final Context:
  workspace = "dev" (from workspace.yaml:workspace.name)
  app = "app-one" (from directory containing application.yaml)
```

---

## Data Schema

### Marker Files

| Marker | Detection Result |
|--------|-----------------|
| `workspace.yaml` | Workspace root directory |
| `application.yaml` | Application root directory |
| `docker-compose.yaml` | Not used for detection (applications must have `application.yaml`) |

### Context Object

When context is successfully detected, it produces:

```yaml
workspace:
  name: "dev"                              # From workspace.yaml:workspace.name
  path: "/home/user/workspaces/dev"        # Absolute path to workspace root

application:                               # May be null if not in an app directory
  name: "app-one"                          # From directory name
  path: "/home/user/workspaces/dev/app-one"  # Absolute path to application root
```

---

## Examples

### Example 1: Full Context Detection

**Directory**: `~/workspaces/dev/app-one/src/components/`

**Directory Structure**:
```
~/workspaces/dev/
├── workspace.yaml              # workspace.name = "dev"
├── app-one/
│   ├── application.yaml
│   └── src/
│       └── components/         # Current directory
└── app-two/
    └── application.yaml
```

**Detected Context**:
- Workspace: `dev` (from `../../../workspace.yaml`)
- Application: `app-one` (from `../../application.yaml`)

**CLI Output**:
```bash
$ contrail app status
# Using workspace: dev (from ../../../workspace.yaml)
# Using app: app-one (from ../../application.yaml)

Status: running
Services: 3 running, 0 stopped
```

### Example 2: Workspace Root (No Application)

**Directory**: `~/workspaces/dev/`

**Detected Context**:
- Workspace: `dev`
- Application: none

**Behavior for app-specific commands**:
```bash
$ contrail app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

### Example 3: Single-App Workspace

**Directory**: `~/my-project/src/`

**Directory Structure**:
```
~/my-project/
├── workspace.yaml              # workspace.name = "dev", applications.myapp.path = "."
├── application.yaml
└── src/                        # Current directory
```

**Detected Context**:
- Workspace: `dev`
- Application: `myapp`

Both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

### Example 4: Nested Vendor Package (Ignored)

**Directory**: `~/workspaces/dev/app-one/vendor/some-package/`

**Directory Structure**:
```
~/workspaces/dev/
├── workspace.yaml
└── app-one/
    ├── application.yaml
    └── vendor/
        └── some-package/
            └── application.yaml  # This package has its own application.yaml
```

**Detected Context**:
- Workspace: `dev`
- Application: `app-one` (NOT `some-package`)

**Rationale**: The vendor package's `application.yaml` is ignored because we walk **toward** the workspace root, finding `app-one/application.yaml` first.

### Example 5: Nested Workspace (Test Fixtures)

**Directory**: `~/workspaces/dev/app-one/tests/fixtures/workspace-test/`

**Directory Structure**:
```
~/workspaces/dev/
├── workspace.yaml              # Parent workspace
└── app-one/
    ├── application.yaml
    └── tests/
        └── fixtures/
            └── workspace-test/
                └── workspace.yaml  # Test fixture workspace
```

**Detected Context**:
- Workspace: from `workspace-test/workspace.yaml` (the closest/innermost wins)
- Application: none (no application.yaml found in test fixture)

**Rationale**: Test fixtures that create their own workspace should be isolated. The closest `workspace.yaml` is used.

---

## Edge Cases

### No Workspace Found, Application.yaml Present

**Scenario**: User is in a directory containing `application.yaml` but no `workspace.yaml` anywhere up the tree.

**Behavior**:
```bash
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

### Neither Workspace Nor Application Found

**Scenario**: User is in a directory with no Contrail configuration files.

**Behavior**:
```bash
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: contrail app status --workspace=NAME --app=NAME

Available workspaces: contrail workspace list
```

### Explicit Flags Override Detection

**Scenario**: User provides `--app` flag while in a different application directory.

**Behavior**: Explicit flags completely replace context detection.

```bash
# From within app-one directory (context would detect app-one)
$ cd ~/workspaces/dev/app-one
$ contrail up -a app-two
# Starting: app-two
# (app-one from context is ignored)

# To start multiple apps, list them all explicitly
$ contrail up -a app-one -a app-two
# Starting: app-one, app-two
```

This "explicit replaces context" behavior ensures predictable results—when you specify apps explicitly, you get exactly what you asked for.

### Symlinked Directories

**Scenario**: Application directory is a symlink.

**Behavior**: Symlinks are resolved before detection. The resolved (real) path is used for the application name.

**Rationale**: This ensures consistent naming regardless of how the directory is accessed.

---

## Error Handling

| Error Condition | Error Code | Message | Recovery |
|-----------------|------------|---------|----------|
| No workspace found | 5 | `No workspace found (workspace.yaml) in current directory or any parent directories` | Navigate to workspace or use explicit flags |
| No app context for app-specific command | 5 | `No application context detected` | Navigate to app directory or use `--app` flag |
| Application.yaml outside workspace | 5 | `Found application.yaml but no workspace.yaml above it` | Create workspace or move application.yaml |

### Exit Code 5

Context detection failures use exit code 5 specifically to distinguish from general errors (code 1). This allows scripts to detect when explicit flags are needed.

---

## Context Feedback

When context is detected, commands indicate what was found:

```bash
$ cd ~/workspaces/dev/app-one
$ contrail app status
# Using workspace: dev (from ../workspace.yaml)
# Using app: app-one (from ./application.yaml)

Status: running
Services: 3 running, 0 stopped
```

Use `--quiet` or `-q` to suppress context indicators:

```bash
$ contrail app status -q
running
```

---

## Global Commands

Some commands do not use directory context:

| Command | Context Used |
|---------|--------------|
| `workspace` | Yes (workspace only) |
| `app` | Yes (workspace + application) |
| `flavor` | Yes (workspace + application) |
| `port` | No (operates on global state) |
| `proxy` | No (operates on global proxy) |
| `config` | No (operates on global config) |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-02 | Initial specification extracted from CLI reference |
