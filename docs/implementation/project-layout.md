<!-- Migrated from specs/scind-go-stack.md:183-267 -->
<!-- Extraction ID: impl-project-layout -->

# Project Structure

```
scind/
├── cmd/
│   └── scind/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── cli/                        # Cobra command definitions
│   │   ├── root.go                 # Root command, global flags, context detection
│   │   ├── workspace.go            # workspace subcommands
│   │   ├── app.go                  # app subcommands
│   │   ├── flavor.go               # flavor subcommands
│   │   ├── port.go                 # port subcommands
│   │   ├── proxy.go                # proxy subcommands
│   │   ├── config.go               # config subcommands
│   │   ├── aliases.go              # Top-level aliases (up, down, ps, generate)
│   │   ├── compose_prefix.go       # compose-prefix command
│   │   ├── init_shell.go           # init-shell command
│   │   ├── completion.go           # completion command
│   │   ├── validate.go             # validate command
│   │   ├── doctor.go               # doctor command
│   │   ├── open.go                 # open command
│   │   └── urls.go                 # urls command
│   │
│   ├── config/                     # Configuration types and loading
│   │   ├── workspace.go            # workspace.yaml schema
│   │   ├── application.go          # application.yaml schema
│   │   ├── proxy.go                # proxy.yaml schema
│   │   ├── state.go                # state.yaml (global and workspace)
│   │   ├── manifest.go             # manifest.yaml (computed values)
│   │   └── loader.go               # Viper-based config loading
│   │
│   ├── context/                    # Context detection
│   │   ├── detector.go             # Directory walking, file detection
│   │   └── resolver.go             # Resolve workspace/app from flags or detection
│   │
│   ├── generator/                  # Override file generation
│   │   ├── generator.go            # Main generation logic
│   │   ├── templates.go            # Sprig templates for hostnames, aliases
│   │   ├── override.go             # Docker Compose override file builder
│   │   ├── manifest.go             # Manifest file builder
│   │   └── traefik.go              # Traefik label generation
│   │
│   ├── docker/                     # Docker interaction
│   │   ├── compose.go              # exec.Command wrapper for docker compose
│   │   ├── network.go              # Network creation/inspection via SDK
│   │   ├── container.go            # Container listing/inspection
│   │   └── port.go                 # Port availability checking
│   │
│   ├── workspace/                  # Workspace operations
│   │   ├── lifecycle.go            # up, down, restart logic
│   │   ├── status.go               # Status aggregation
│   │   └── clone.go                # Git clone operations
│   │
│   ├── output/                     # Output formatting
│   │   ├── table.go                # Table formatter
│   │   ├── json.go                 # JSON output
│   │   ├── yaml.go                 # YAML output
│   │   └── quiet.go                # Quiet/names-only output
│   │
│   └── scripts/                    # Embedded shell scripts
│       ├── bash.sh
│       ├── zsh.zsh
│       └── fish.fish
│
├── pkg/                            # Public API (if needed for plugins later)
│   └── plugin/                     # Plugin interface definitions
│       ├── protocol.go             # Protocol handler interface
│       └── types.go                # Shared types for plugins
│
├── testdata/                       # Test fixtures
│   ├── workspaces/
│   │   ├── valid/
│   │   └── invalid/
│   └── applications/
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Package Responsibilities

### `internal/cli/`

Contains all Cobra command definitions. Each file corresponds to a resource type or command group:

- **root.go** - Root command setup, global flags (`-w`, `-a`, `-q`, `-v`, `--json`, `--yaml`, `--color`), context detection via `PersistentPreRunE`
- **workspace.go** - Workspace lifecycle and management commands
- **app.go** - Application lifecycle and management commands
- **flavor.go** - Flavor listing and switching
- **aliases.go** - Top-level convenience aliases (`up`, `down`, `ps`, `generate`)

### `internal/config/`

Configuration type definitions and loading logic:

- Struct definitions with validation tags
- Viper-based loading with layered merging
- Inference rules (service name defaulting, port inference)
- Schema validation

### `internal/context/`

Directory-walking context detection:

- Workspace boundary detection (find `workspace.yaml` walking up)
- Application detection (find `application.yaml` within workspace)
- Resolver for flag-based overrides

### `internal/generator/`

Override file generation:

- Hostname and alias template evaluation using Sprig
- Docker Compose override file assembly
- Manifest file generation
- Traefik label generation for proxied services

### `internal/docker/`

Docker interaction layer:

- `exec.Command` wrapper for `docker compose` calls
- Docker SDK usage for network/container operations
- Port availability checking

### `internal/output/`

Output formatting:

- Table formatting for `list` and `status` commands
- JSON/YAML output modes
- Quiet mode (names only, scriptable)
