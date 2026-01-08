## Context Detection Algorithm

Context detection uses a **workspace boundary** approach to prevent accidental detection of config files in vendor packages or nested test fixtures.

**Detection steps**:
1. **Find workspace root**: Walk up from current working directory looking for `workspace.yaml`. The first one found establishes the workspace root.
2. **Find application context**: Walk up from current working directory toward the workspace root looking for `application.yaml`. Only consider `application.yaml` files that are within the workspace directory tree.
3. **Never traverse above workspace root** for application detection—this prevents vendor packages or nested fixtures from hijacking context.

**Error handling** (when no workspace found):
- If `application.yaml` found but no `workspace.yaml`:
  ```
  No workspace found (workspace.yaml) in current directory or any parent directories,
  but found an application (application.yaml) at: /path/to/application.yaml
  ```
- If neither found:
  ```
  No workspace found (workspace.yaml) in current directory or any parent directories,
  and no application (application.yaml) found either
  ```

**Edge cases**:
- **Nested vendor packages**: If working in `frontend/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored because the workspace's `frontend/application.yaml` is found first when walking toward the workspace root.
- **Workspace within workspace**: If a test fixture has its own `workspace.yaml` nested inside a workspace, the closest (innermost) `workspace.yaml` wins—this is typically the test fixture, which is the expected behavior.
- **Single-app workspaces**: With `path: .`, both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

## Error Cases

### No Configuration Found

If no `workspace.yaml` or `application.yaml` is found in the current directory or parents:
- Error: "Not in a workspace or application directory"
- Exit code: 1

### Ambiguous Context

If both `workspace.yaml` and `application.yaml` exist in the same directory:
- Behavior: Treat as workspace context (single-app workspace pattern)
- The application is resolved using the workspace's `applications` section

### Quick Reference

```bash
# Workspace operations
scind workspace init --workspace=dev
scind workspace up [-w NAME]
scind workspace down [-w NAME]
scind workspace status [-w NAME]

# Application operations
scind app add --app=NAME --repo=URL
scind app up [-a NAME]

# Flavor management
scind flavor set FLAVOR [-a NAME]

# Port management
scind port list
scind port gc

# Top-level aliases (with context detection)
scind up
scind down

# Docker Compose passthrough (shell function)
scind-compose exec php bash
scind-compose logs -f
scind-compose -a backend ps
```

## Related Documents

- [Shell Integration](./shell-integration.md) - How context detection integrates with shell functions
- [Directory Structure](./directory-structure.md) - File locations for workspace.yaml and application.yaml
- [ADR-0011: Options-Based Targeting](../decisions/0011-options-based-targeting.md) - CLI targeting strategy
