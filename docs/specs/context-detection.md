<!-- Migrated from specs/contrail-technical-spec.md:1059-1122 -->
<!-- Extraction ID: spec-context-detection -->

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
- **Nested vendor packages**: If working in `app-one/vendor/some-package/` where the vendor package has its own `application.yaml`, it is ignored because the workspace's `app-one/application.yaml` is found first when walking toward the workspace root.
- **Workspace within workspace**: If a test fixture has its own `workspace.yaml` nested inside a workspace, the closest (innermost) `workspace.yaml` wins—this is typically the test fixture, which is the expected behavior.
- **Single-app workspaces**: With `path: .`, both `workspace.yaml` and `application.yaml` are in the same directory, so detection finds both immediately.

### Quick Reference

```bash
# Workspace operations
contrail workspace init --workspace=dev
contrail workspace up [-w NAME]
contrail workspace down [-w NAME]
contrail workspace status [-w NAME]

# Application operations
contrail app add --app=NAME --repo=URL
contrail app up [-a NAME]

# Flavor management
contrail flavor set FLAVOR [-a NAME]

# Port management
contrail port list
contrail port gc

# Top-level aliases (with context detection)
contrail up
contrail down

# Docker Compose passthrough (shell function)
contrail-compose exec php bash
contrail-compose logs -f
contrail-compose -a app-two ps
```
