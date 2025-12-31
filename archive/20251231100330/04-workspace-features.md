# Issue Group 4: Workspace Features

**Documents Affected**: PRD + Technical Spec + CLI Reference  
**Suggested Order**: 4 of 10 (high-severity core functionality gaps)  
**Estimated Effort**: Large

---

## Overview

These are significant gaps in workspace-level functionality. Two are high-severity issues that need resolution before implementation can proceed confidently.

---

## Issues

### M-1: No Workspace Discovery Mechanism Specified

**Severity**: High

**Issue**: `contrail workspace list` is documented in CLI Reference (lines 127-141) but there's no specification for how Contrail discovers workspaces.

**Current output implies discovery works**:
```
NAME     APPS  STATUS   PATH
dev      3     running  ~/workspaces/dev
review   3     stopped  ~/workspaces/review
control  2     running  ~/workspaces/control
```

But how does Contrail know about these workspaces?

**Options**:

**A) Global Registry File**
```yaml
# ~/.config/contrail/workspaces.yaml
workspaces:
  - path: ~/workspaces/dev
  - path: ~/workspaces/review
  - path: ~/projects/my-app
```
- Pro: Explicit, fast lookup
- Con: Can get out of sync, requires registration step

**B) Scan from Current Directory**
- Walk up from CWD looking for workspace.yaml files
- Pro: No registration needed
- Con: Only finds workspaces in current tree, `list` from ~ would find nothing

**C) Scan Known Locations**
- Check `~/workspaces/*/workspace.yaml`, configurable search paths
- Pro: Convention-based, no registration
- Con: Misses workspaces outside conventions

**D) Hybrid: Registry with Auto-Registration**
- `workspace init` automatically adds to registry
- `workspace list` reads registry
- `workspace prune` removes stale entries
- Pro: Best of both worlds
- Con: More complexity

**Questions**:
1. Which approach fits Contrail's philosophy best?
2. Should `workspace init` auto-register?
3. Should there be a `workspace register` command for existing workspaces?

**Response**:
> Global Registry file with a fallback / ability to scan for workspaces using Docker labels to bootstrap / update the global registry. For example, if someone removes `~/.config/contrail/workspaces.yaml`, we should be able to reconstruct most of it by looking at Docker labels from running containers.

---

### M-6: Single-App Workspace Clone Behavior Undefined

**Severity**: Medium

**Issue**: When `path: .` is used (single-app workspace), what does `contrail workspace clone` do?

**Scenario**:
```yaml
# workspace.yaml
workspace:
  name: dev
  applications:
    myapp:
      path: .
      repository: git@github.com:org/myapp.git  # Is this valid?
```

**Questions**:
1. Should `path: .` and `repository:` be mutually exclusive?
2. If both are present, should clone skip it? Error? Warn?
3. Should clone skip all apps with `path: .`?
4. What if `path: .` but no repository—is that the expected case?

**Suggested Resolution**:
- `path: .` implies the application is already present (it's the workspace directory)
- `repository:` is only meaningful for `path: ./subdir` or `path: ./app-name`
- `clone` should skip apps with `path: .` and log: "Skipping myapp: application is workspace root"
- Validation should warn if both `path: .` and `repository:` are set

**Response**:
> Clone should skip apps with `path: .` and log: "Skipping myapp: application is workspace root"

---

### M-7: No Collision Detection for Duplicate Workspace Names

**Severity**: High

**Issue**: No error handling is documented for when two directories have `workspace.yaml` files with the same `name:` value. Global port assignments in `~/.config/contrail/state.yaml` would collide silently.

**Example collision**:
```yaml
# ~/project-a/workspace.yaml
workspace:
  name: dev

# ~/project-b/workspace.yaml  
workspace:
  name: dev   # Same name!
```

Both would write to `assigned_ports.dev.*` in global state.

**Options**:

**A) Enforce Uniqueness via Registry** (pairs with M-1 Option D)
- `workspace init` checks registry for name conflicts
- Error: "Workspace 'dev' already exists at ~/project-a"

**B) Include Path Hash in State Keys**
- `assigned_ports.dev-a1b2c3.*` where hash is from workspace path
- Pro: No collisions possible
- Con: State file becomes less human-readable

**C) Namespace by Path**
- `assigned_ports["/home/user/project-a"].dev.*`
- Pro: Clear what's what
- Con: Verbose state file

**D) Detect at Runtime, Warn/Error**
- When loading state, check if multiple paths claim same workspace name
- Pro: No registry needed
- Con: Reactive rather than preventive

**Questions**:
1. Is workspace name uniqueness a hard requirement?
2. Should the solution tie into M-1 (workspace discovery)?
3. How important is human-readable state files?

**Response**:
> 1. Yes, workspace names should be unique since they are intended to be a unique namespace across an entire host.
> 2. Yes, tie into M-1—any attempt to access the "list of workspaces" should take the global registry into account *and also check Docker labels* to see if any new workspaces are active that aren't in the registry for some reason.
> 3. Human-readable state files are not super important. (The question about hash suffixes in Option B was about making state keys like `assigned_ports.dev-a1b2c3.*` instead of just `assigned_ports.dev.*`, but this is moot since we're enforcing uniqueness via the registry + Docker label fallback.)

---

## Checklist

- [x] Design and document workspace discovery mechanism (registry + Docker label fallback)
- [x] Add discovery details to Tech Spec (new "Workspace Registry" section)
- [x] Add `workspace prune` command to CLI Reference
- [x] Update `workspace list` with `--validate` and `--rebuild` flags
- [x] Document single-app workspace clone behavior (`path: .` skipped with message)
- [x] Document workspace name uniqueness enforcement in `workspace init`
- [x] Add `workspace.path` label to generated override examples
- [x] Add `workspacePruneCmd` to Go Stack command scaffolding

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
