## Directory Structure

Applications are **derived, not tracked**: the set of applications in a workspace
is resolved from the workspace directory's contents (the application
subdirectories present on disk), with no separate application registry to keep in
sync. The Xcind proof-of-concept validated this model end-to-end.

### Standard Multi-Application Workspace

```
~/.config/scind/
├── proxy.yaml                        # Global proxy configuration
├── state.yaml                        # Global port assignments and inventory
└── workspaces.yaml                   # Workspace registry

workspaces/
├── proxy/
│   ├── docker-compose.yaml           # Traefik service definition
│   ├── traefik.yaml                  # Traefik static configuration
│   └── .env                          # Proxy-level environment variables
│
├── dev/                              # Workspace root
│   ├── workspace.yaml                # Workspace configuration (structure)
│   │
│   ├── overrides/                    # Manual overrides (optional, workspace-specific)
│   │   └── backend.yaml              # Merged after generated config
│   │
│   ├── .generated/                   # Generated files (gitignored)
│   │   ├── state.yaml                # Runtime state (active flavors)
│   │   ├── manifest.yaml             # Computed values (read-only)
│   │   ├── frontend.override.yaml    # Generated compose override
│   │   ├── backend.override.yaml
│   │   └── shared-db.override.yaml
│   │
│   ├── frontend/                     # Cloned application repository
│   │   ├── docker-compose.yaml       # Application's compose file (app-owned)
│   │   ├── application.yaml          # Service contract + flavors (app-owned)
│   │   └── ...
│   ├── backend/
│   │   ├── docker-compose.yaml
│   │   ├── docker-compose.worker.yaml
│   │   ├── docker-compose.extras.yaml
│   │   ├── application.yaml
│   │   └── ...
│   └── shared-db/
│       ├── docker-compose.yaml
│       ├── application.yaml
│       └── ...
│
├── review/                           # Another workspace (same structure)
│   └── ...
│
└── control/                          # Another workspace
    └── ...
```

### Single-Application Workspace

When promoting an existing Docker Compose project:

```
~/my-project/                         # Workspace AND application directory
├── workspace.yaml                    # workspace.name = "dev"
├── application.yaml                  # Application service contract
├── docker-compose.yaml               # Existing compose file (unchanged)
├── docker-compose.worker.yaml        # Optional additional compose files
├── .generated/                       # Generated files (gitignored)
│   ├── state.yaml
│   ├── manifest.yaml
│   └── my-project.override.yaml
├── overrides/                        # Manual overrides (optional)
└── src/                              # Application source code
```

In this configuration, `workspace.yaml` references the application with `path: .`:

```yaml
workspace:
  name: dev
  applications:
    my-project:
      path: .                         # Application is in workspace root
```

### External Applications

An application directory does not have to sit inside the workspace tree. `applications.{name}.path` may be **absolute, or relative and pointing outside the workspace**, so a repository that already lives elsewhere on disk can join a workspace without moving:

```yaml
workspace:
  name: dev
  applications:
    frontend: {}                      # ./frontend — conventional, inside the tree
    shared-db:
      path: ~/src/shared-db           # Outside the workspace tree
```

This does not invert the ownership direction. The workspace still enumerates its applications, so the roster stays authoritative and enumerable for `workspace up`, the generated manifest, and validation — the application says nothing about its own membership. What changes is only where the directory sits.

**Detection does not follow the path outward.** [Context detection](./context-detection.md) never traverses above the workspace root, so running Scind from inside an external application's directory detects neither the workspace nor the application. Target external applications explicitly with `-w`/`-a` ([ADR-0011](../decisions/0011-options-based-targeting.md)):

```bash
scind app up --workspace=dev --app=shared-db
```

Relative paths resolve from the workspace root and `~` expands to the user's home directory. A path that leaves the workspace tree is machine-specific, which makes that `workspace.yaml` less portable than the conventional `./{app}` form — a real cost, and the reason the conventional form stays the default.
