# Issue Group 6: Context Detection

**Documents Affected**: Technical Spec + CLI Reference + Go Stack  
**Suggested Order**: 6 of 10 (depends on workspace structure understanding from earlier groups)  
**Estimated Effort**: Small

---

## Overview

A single focused issue about how context detection handles edge cases with nested or conflicting configuration files.

---

## Issues

### A-3: Nested/Conflicting Config File Precedence Undefined

**Severity**: Medium

**Issue**: The context detection documentation describes walking up the directory tree to find `workspace.yaml` and `application.yaml`, but doesn't address edge cases.

**Location**: CLI Reference lines 36-52

**Scenario 1: Nested vendor package**
```
~/workspaces/dev/app-one/
├── application.yaml           # app-one's config
├── vendor/
│   └── some-package/
│       └── application.yaml   # vendor package has its own!
└── src/
```

If working in `~/workspaces/dev/app-one/vendor/some-package/`, which `application.yaml` wins?

**Scenario 2: Workspace within workspace**
```
~/workspaces/dev/
├── workspace.yaml             # dev workspace
├── app-one/
│   ├── application.yaml
│   └── test-fixtures/
│       └── workspace.yaml     # test fixture workspace!
```

**Scenario 3: Application outside workspace**
```
~/random-project/
├── application.yaml           # Has app config but no workspace
└── src/
```

**Questions**:
1. Should "closest to CWD" always win (current implied behavior)?
2. Should there be a boundary—only configs within a workspace are valid?
3. Should orphan `application.yaml` (no parent workspace) be an error or work standalone?
4. Should there be a max depth or marker file to stop traversal?

**Options**:

**A) Closest Wins (Current Implied)**
- First `application.yaml` found walking up is used
- First `workspace.yaml` found walking up is used
- Pro: Simple, predictable
- Con: Vendor packages could hijack context

**B) Workspace Boundary**
- Find `workspace.yaml` first
- Only look for `application.yaml` within that workspace's directory tree
- Pro: Clear boundaries, no vendor hijacking
- Con: More complex detection logic

**C) Explicit Ignore Patterns**
- Skip `vendor/`, `node_modules/`, `.git/` when walking up
- Pro: Handles common cases
- Con: Incomplete, different ecosystems have different patterns

**D) Marker File**
- Stop traversal at `.contrail-root` or when hitting `workspace.yaml`
- Pro: Explicit control
- Con: Another file to manage

**Suggested Resolution**: Option B (Workspace Boundary) with fallback:
1. Walk up to find `workspace.yaml` → establishes workspace root
2. Walk up from CWD to workspace root looking for `application.yaml`
3. If no `workspace.yaml` found, `application.yaml` alone is valid (single-app standalone mode)
4. Never traverse above workspace root for application detection

**Response**:
> Use Option B (Workspace Boundary) with required `workspace.yaml`. Since `workspace.yaml` is required in all cases (even single-app workspaces use `path: .`), context detection should:
> 1. Walk up to find `workspace.yaml` → establishes workspace root
> 2. Walk up from CWD to workspace root looking for `application.yaml`
> 3. Never traverse above workspace root for application detection
> 4. If no `workspace.yaml` found, error with helpful debugging info:
>    - If `application.yaml` found: "No workspace found (workspace.yaml) in current directory or any parent directories, but found an application (application.yaml) at: /full/path/to/application.yaml"
>    - If neither found: "No workspace found (workspace.yaml) in current directory or any parent directories, and no application (application.yaml) found either"

---

## Documents to Update

### Technical Spec
Add to "Context Detection" or "CLI Interface" section:
- Detection algorithm with workspace boundary behavior
- Edge case handling (nested configs, orphan apps)

### CLI Reference
Update "Detection Rules" (lines 36-52):
- Add boundary behavior
- Add edge case examples
- Document standalone application mode

### Go Stack
Update `detectAndSetContext` implementation notes:
- Algorithm description
- Boundary detection logic

---

## Checklist

- [x] Define context detection algorithm with workspace boundaries
- [x] Document edge cases in CLI Reference
- [x] Add algorithm details to Tech Spec
- [x] Update Go Stack implementation notes
- [ ] Add test cases for edge scenarios to testdata/ (deferred to implementation)
