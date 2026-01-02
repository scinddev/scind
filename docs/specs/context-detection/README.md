# Context Detection Specification

**Version**: 0.5.0
**Date**: December 2024
**Status**: Accepted

---

## Overview

Contrail automatically detects workspace and application context from the current directory, reducing the need for explicit `--workspace` and `--app` flags.

---

## Detection Algorithm

Context detection uses a **workspace boundary** approach to prevent accidental detection of config files in vendor packages or nested test fixtures.

### Steps

1. **Find workspace root**: Walk up from current directory looking for `workspace.yaml`
   - First one found establishes the workspace root
   - The `workspace.name` value becomes implicit `--workspace`

2. **Find application context**: Walk up from current directory toward workspace root looking for `application.yaml`
   - Only considers files **within the workspace directory tree**
   - Never traverses above workspace root
   - Directory name containing `application.yaml` becomes implicit `--app`

3. **Both can be detected simultaneously**:
   ```
   ~/workspaces/dev/app-one/src/components/
                   │        │
                   │        └── application.yaml → app = "app-one"
                   │
                   └── workspace.yaml → workspace = "dev"
   ```

---

## Flag Override Behavior

Explicit flags **replace** (not add to) context detection:

```bash
# From within app-one directory
$ cd ~/workspaces/dev/app-one

# This starts ONLY app-two (app-one from context is ignored)
$ contrail up -a app-two

# To start multiple apps, list them explicitly
$ contrail up -a app-one -a app-two
```

---

## Edge Cases

### Nested Vendor Packages

If working in `app-one/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored. The workspace's `app-one/application.yaml` is found first when walking toward the workspace root.

### Workspace Within Workspace

If a test fixture has its own `workspace.yaml` nested inside a workspace, the closest (innermost) `workspace.yaml` wins—this is the test fixture, which is the expected behavior.

### Single-App Workspaces

With `path: .`, both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

---

## Context Feedback

Commands indicate detected context:

```bash
$ cd ~/workspaces/dev/app-one
$ contrail app status
# Using workspace: dev (from ../workspace.yaml)
# Using app: app-one (from ./application.yaml)

Status: running
Services: 3 running, 0 stopped
```

Use `--quiet` or `-q` to suppress context indicators.

---

## Error Messages

### No workspace found, but application.yaml exists

```
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /path/to/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

### Neither found

```
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: contrail app status --workspace=NAME --app=NAME

Available workspaces: contrail workspace list
```

### Workspace found but no application context

```
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

---

## Global Commands

The following commands **ignore** directory context:

- `port` commands
- `proxy` commands
- `config` commands
- `doctor`
- `completion`
- `init-shell`

---

## Related Documentation

- [ADR-0011: Options-Based Targeting](../../decisions/0011-options-based-targeting/README.md)
- [CLI Reference](../../reference/cli/README.md)
