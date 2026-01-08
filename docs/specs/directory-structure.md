## Directory Structure

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
