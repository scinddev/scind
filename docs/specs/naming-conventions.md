## Conventions and Best Practices

### Application Requirements

For an application to work well within a workspace, it should:

1. **Include an `application.yaml`**: This file defines the service contract - which services the application exports to the workspace. This is owned and maintained by the application developers.
2. **Use environment variables for external service URLs**: Don't hardcode hostnames for dependencies. Use the injected `SCIND_{APP}_{EXPORTED_SERVICE}_*` variables.
3. **Expose ports without host bindings**: Use `ports: ["8080"]` not `ports: ["8080:8080"]` to avoid conflicts.
4. **Use relative volume paths**: Ensure builds and mounts work regardless of absolute path.

### The Service Contract (`application.yaml`)

The `application.yaml` file is the interface between the application and the workspace system:

- **Owned by**: Application developers (committed to the application's git repository)
- **Application name**: Inferred from the directory name (no explicit `name:` field required)
- **Purpose**:
  - Declares which services the application exports and how they should be exposed
  - Defines flavors (different ways to run the application with different compose file combinations)
- **Consumed by**: Workspace tooling (to generate override files)

Application developers should update `application.yaml` when:
- Adding a new service that other applications need to access
- Changing the type (proxied, assigned) or protocol (http, https) of an exported service
- Renaming services that are exposed to the workspace
- Adding a new flavor (e.g., a "lite" mode that excludes optional services)
- Changing which compose files are needed for a flavor

### Naming Conventions

- **Workspace names**: Lowercase alphanumeric with hyphens (e.g., `dev`, `feature-x`, `pr-123`)
- **Application names**: Lowercase alphanumeric with hyphens, inferred from directory name
- **Compose project names**: `{workspace}-{application}` (e.g., `dev-frontend`), or `{workspace}-{instance}-{application}` when the per-instance token is non-empty
- **Workspace-internal network**: `{workspace}-internal` (e.g., `dev-internal`), or `{workspace}-{instance}-internal` when the per-instance token is non-empty
- **Exported service names**: The key in `exported_services`, may differ from the underlying Compose service name
- **Proxied hostnames** (proxied type): `{workspace}-{application}-{exported_service}.{domain}` (e.g., `dev-frontend-web.scind.test`)
- **Internal aliases** (all types): `{application}-{exported_service}` (e.g., `frontend-web`, `shared-db-db`)
- **Environment variables**: `SCIND_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}` in SCREAMING_SNAKE_CASE
- **Traefik router names**: `{workspace}-{application}-{exported_service}-{protocol}` (e.g., `dev-frontend-web-https`)
- **Apex hostnames** (proxied primary export): `{workspace}-{application}.{domain}` (e.g., `dev-frontend.scind.test`)
- **Apex internal aliases** (all primary exports): `{application}` (e.g., `frontend`)
- **Apex Traefik router names** (proxied primary export): `{workspace}-{application}-{protocol}` (e.g., `dev-frontend-https`)
- **Apex environment variables** (proxied primary export): `SCIND_{APPLICATION}_APEX_{SUFFIX}` (e.g., `SCIND_FRONTEND_APEX_URL`)

**Implicit primary**: If an application has exactly one exported service, it is implicitly primary — no annotation needed. Apex patterns are only generated for the primary exported service. See [ADR-0013](../decisions/0013-apex-url-primary-designation.md) for the design rationale.

### Per-Instance Token

The **per-instance disambiguation token** ([ADR-0016](../decisions/0016-per-instance-isolation.md)) folds into the compose project name and the workspace-internal network name so that multiple linked git worktrees of the same repository derive distinct Docker resources instead of colliding:

- **Empty token** (default): names are exactly the base patterns above (`{workspace}-{application}`, `{workspace}-internal`) — fully backward compatible, no migration.
- **Non-empty token**: the token is inserted, producing `{workspace}-{instance}-{application}` and `{workspace}-{instance}-internal`.

**Resolution order** (first match wins):

1. **Explicit value** — configured/passed by the user.
2. **Opt-out** — an explicit request for the empty token (preserve legacy names).
3. **Auto-detect** — a stable token derived from a linked git worktree when one is detected. Auto-detection keys on the worktree-directory identity, **not** the branch (branches move; the resource identity must stay pinned to the working copy).

**Collision warning**: Docker Compose project names, Traefik router names, volume names, and network names are derived from the naming patterns above. Creative naming that produces identical project names could cause conflicts:
- **Traefik routers**: Conflicting router names cause routing failures
- **Docker volumes**: Conflicting volume names (e.g., `dev-frontend_postgres_data`) could cause data to be shared unexpectedly or overwritten
- **Docker networks**: Conflicting network names could connect unrelated services

Example collision: workspace `dev-front` with app `end` and workspace `dev` with app `frontend` both produce project name `dev-frontend`.

**Multi-worktree collision**: The same repository checked out into multiple linked git worktrees derives the **same** workspace name — and therefore the same project name and `{workspace}-internal` network — in every worktree. Without a per-instance token this is a collision, not isolation: the second `up` attaches to or disrupts the first. Supplying a non-empty token per worktree (see resolution order above) restores isolation. See [ADR-0016](../decisions/0016-per-instance-isolation.md).

Follow the lowercase-alphanumeric-with-hyphens convention and avoid names that could produce ambiguous concatenations.

### Git Strategy

**Workspace repository** (optional - can be version controlled):
- `workspace.yaml` - workspace definition
- `overrides/` - manual overrides

**Generated files** (gitignored):
- `.generated/` - generated override files

**Application directories** (cloned repositories):
- Each application is its own git repository
- `application.yaml` lives in the application repo and defines its service contract
- Can be managed as submodules, or simply cloned separately

Example workspace `.gitignore`:
```
.generated/
app-*/
```

## Related Documents

- [ADR-0004: Convention-Based Naming](../decisions/0004-convention-based-naming.md) - Rationale for naming conventions
- [Port Types](./port-types.md) - How port types affect naming
- [Generated Override Files](./generated-override-files.md) - Naming in generated Docker labels
- [Docker Labels](./docker-labels.md) - Label naming conventions
- [ADR-0013: Apex URL Primary Designation](../decisions/0013-apex-url-primary-designation.md) - Why `primary: true` field
- [ADR-0016: Per-Instance Isolation](../decisions/0016-per-instance-isolation.md) - Per-instance token in project and network names
