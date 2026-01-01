# Specification: Context Detection

**Version**: 0.2.0
**Date**: December 2024

<!-- Migrated from specs/contrail-cli-reference.md and specs/contrail-technical-spec.md -->

---

## Overview

Contrail automatically detects workspace and application context from the current directory, reducing the need for explicit `--workspace` and `--app` flags.

---

## Detection Algorithm

Context detection uses a **workspace boundary** approach to prevent accidental detection of config files in vendor packages or nested test fixtures.

### Step 1: Find Workspace Root

Walk up the directory tree from current working directory looking for `workspace.yaml`. The first one found establishes the **workspace root**.

### Step 2: Find Application Context (Bounded by Workspace)

Walk up from current directory toward the workspace root looking for `application.yaml`:
- Only considers `application.yaml` files **within the workspace directory tree**
- If found, the directory name containing it becomes the implicit `--app` value
- **Never traverses above the workspace root**—this prevents vendor packages from hijacking context

### Step 3: Apply Context

- The `workspace.name` value becomes the implicit `--workspace` value
- The directory containing `application.yaml` becomes the implicit `--app` value
- Explicit flags always override detected context

---

## Flag Override Behavior

When explicit flags are provided, they **replace** (not add to) context detection:

```bash
# From within app-one directory (context would detect app-one)
$ cd ~/workspaces/dev/app-one

# This starts ONLY app-two, not both app-one and app-two
$ contrail up -a app-two
# Starting: app-two
# (app-one from context is ignored)

# To start multiple apps, list them all explicitly
$ contrail up -a app-one -a app-two
# Starting: app-one, app-two
```

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

Use `--quiet` or `-q` to suppress context indicators.

---

## Edge Cases

### Nested Vendor Packages

If working in `app-one/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored. The workspace's `app-one/application.yaml` is found first when walking toward the workspace root.

### Workspace Within Workspace

If a test fixture has its own `workspace.yaml` nested inside a workspace (e.g., for integration tests), the closest `workspace.yaml` wins—this is the test fixture, which is the expected behavior.

### Single-App Workspaces

With `path: .`, both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

---

## Error Handling

### No Workspace Found, But Application.yaml Exists

```bash
$ cd ~/random-project
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

### Neither Workspace Nor Application Found

```bash
$ cd ~
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: contrail app status --workspace=NAME --app=NAME

Available workspaces: contrail workspace list
```

### Workspace Found But No Application Context

```bash
$ cd ~/workspaces/dev
$ contrail app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

---

## Global Commands

These commands ignore directory context entirely:

- `port` commands
- `proxy` commands
- `config` commands
- `doctor`
- `completion`
- `init-shell`

---

## Related Documents

- [CLI Reference](../reference/cli.md)
- [ADR-0011: Options-Based Targeting](../decisions/0011-options-based-targeting.md)
