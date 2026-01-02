# Technology Stack

**Version**: 1.0.0
**Date**: 2024-12

---

## Overview

Contrail is implemented in Go, chosen for:
- Single binary distribution (no runtime dependencies)
- Excellent CLI library ecosystem
- Strong Docker SDK support
- Fast compilation and execution

The Go stack is intentionally aligned with the Weftlo project where patterns overlap, enabling knowledge transfer and consistent tooling across projects.

---

## Philosophy

### Minimal Dependencies

Only include dependencies that provide significant value. Prefer standard library where reasonable.

### No Frameworks

Use libraries, not frameworks. Keep control of the application lifecycle.

### Explicit Over Magic

Configuration and behavior should be explicit and traceable.

---

## Core Dependencies

### CLI: Cobra + Viper

**Cobra** provides command structure, **Viper** handles configuration.

```go
import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)
```

**Why**: Industry standard for Go CLIs. Used by kubectl, docker, gh.

**Patterns**:
- One command per file in `cmd/`
- Viper binds flags to config automatically
- Environment variables via `CONTRAIL_` prefix

### YAML: gopkg.in/yaml.v3

```go
import "gopkg.in/yaml.v3"
```

**Why**: Full YAML 1.2 support, better than encoding/json for config files.

**Patterns**:
- Strict unmarshaling to catch typos
- Custom unmarshalers for complex types

### Docker: docker/docker

```go
import (
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
)
```

**Why**: Official Docker SDK, required for container and network operations.

**Patterns**:
- Create client with `client.NewClientWithOpts(client.FromEnv)`
- Always close response bodies
- Use context for cancellation

### Filesystem: afero

```go
import "github.com/spf13/afero"
```

**Why**: Filesystem abstraction critical for testing override file generation without touching disk.

**Patterns**:
- In-memory filesystem for unit tests
- Real filesystem for production

### Templates: sprig

```go
import "github.com/Masterminds/sprig/v3"
```

**Why**: Template functions for hostname templates, alias templates, generated override file content.

### Validation: go-playground/validator

```go
import "github.com/go-playground/validator/v10"
```

**Why**: Declarative validation via struct tags.

**Patterns**:
- Validate at load time, not use time
- Custom validators for domain rules

### Testing: testify

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)
```

**Why**: Cleaner assertions, mock generation.

**Patterns**:
- `assert` for non-fatal checks
- `require` for fatal preconditions
- Mocks generated with mockery

---

## Complete Dependency List

```go
// go.mod
module github.com/yourorg/contrail

go 1.23

require (
    // CLI Framework
    github.com/spf13/cobra v1.10.2
    github.com/spf13/viper v1.21.0
    github.com/spf13/pflag v1.0.10

    // Filesystem & Templates
    github.com/spf13/afero v1.15.0
    github.com/Masterminds/sprig/v3 v3.3.0

    // Validation & Serialization
    github.com/go-playground/validator/v10 v10.29.0
    gopkg.in/yaml.v3 v3.0.1

    // Docker Integration
    github.com/docker/docker v27.4.0+incompatible

    // Testing
    github.com/stretchr/testify v1.11.1
)
```

### Future Dependencies (Add When Needed)

```go
// Plugin system - add when implementing protocol plugins
github.com/hashicorp/go-plugin v1.6.2

// Compose file parsing - add only if programmatic parsing needed
github.com/compose-spec/compose-go/v2 v2.4.0
```

### Intentionally Excluded

| Package | Reason |
|---------|--------|
| **compose-go** | Not needed initially—Contrail shells out to `docker compose` rather than parsing compose files programmatically. Add later if direct parsing becomes necessary. |
| **Alternative CLI frameworks** (urfave/cli, Kong) | Cobra has superior ecosystem (completion, doc generation) and precedent (docker, kubectl, gh, terraform). |
| **Alternative plugin systems** | go-plugin is the clear choice for Contrail's requirements. |

---

## Project Structure

```
contrail/
├── cmd/
│   └── contrail/
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

**Key principle**: `internal/` for all application code. Only expose `pkg/` if building a library.

---

## Patterns

### Error Handling

Wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("loading workspace config: %w", err)
}
```

Use error types for programmatic handling:

```go
type ConfigNotFoundError struct {
    Path string
}

func (e *ConfigNotFoundError) Error() string {
    return fmt.Sprintf("config not found: %s", e.Path)
}
```

### Configuration Loading

```go
type Config struct {
    Proxy     ProxyConfig
    Workspace WorkspaceConfig
    App       AppConfig
}

func LoadConfig(workspacePath string) (*Config, error) {
    // Load in order: proxy (global) -> workspace -> app
}
```

### Configuration Layering

Viper's merge capabilities support Contrail's configuration hierarchy:

1. Global config (`~/.config/contrail/proxy.yaml`)
2. Workspace config (`workspace.yaml`)
3. Application config (`application.yaml`)
4. Manual overrides (`overrides/{app}.yaml`)
5. Environment variables (`CONTRAIL_*`)
6. Command-line flags

### Context Detection

Use Cobra's `PersistentPreRunE` on the root command to implement directory-walking context detection:

```go
var rootCmd = &cobra.Command{
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if isGlobalCommand(cmd) {
            return nil
        }
        return detectContext(cmd)
    },
}
```

Algorithm (workspace boundary):
1. Walk up from CWD looking for `workspace.yaml` → establishes workspace root
2. Walk up from CWD toward workspace root looking for `application.yaml`
   - Only consider `application.yaml` files WITHIN the workspace directory tree
   - Never traverse above workspace root (prevents vendor hijacking)
3. Set results in viper for access by all subcommands

### Shell Completion

Cobra's `ValidArgsFunction` enables dynamic completion:

```go
cmd.RegisterFlagCompletionFunc("workspace",
    func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        workspaces, _ := listWorkspaces()
        return workspaces, cobra.ShellCompDirectiveNoFileComp
    })
```

### Shell Integration Scripts

Use `//go:embed` for script content:

```go
//go:embed scripts/bash.sh
var bashScript string

//go:embed scripts/zsh.zsh
var zshScript string

//go:embed scripts/fish.fish
var fishScript string
```

### Docker Interaction Strategy

**Primary approach**: Shell out to `docker compose` via `exec.Command`. This is simpler, maintains full Docker Compose compatibility.

**Error handling**: Use pass-through with context prefix:

```go
func runCompose(appName string, args ...string) error {
    cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to %s %s:\n", args[0], appName)
        fmt.Fprint(os.Stderr, string(output))
        return err
    }
    fmt.Print(string(output))
    return nil
}
```

**Docker SDK usage**: For operations that don't map to compose commands:
- Network creation/inspection (`{workspace}-internal`)
- Container listing with label filters (`workspace.name=dev`)
- Port availability checking

---

## Testing Strategy

### Unit Tests

Test pure functions and logic in isolation.

```bash
go test ./internal/...
```

### Unit Tests with Afero

Use Afero's in-memory filesystem for testing file operations:

```go
func TestGenerateOverrideFile(t *testing.T) {
    fs := afero.NewMemMapFs()

    afero.WriteFile(fs, "/workspace/workspace.yaml", []byte(`...`), 0644)
    afero.WriteFile(fs, "/workspace/app-one/application.yaml", []byte(`...`), 0644)

    gen := generator.New(fs)
    err := gen.Generate("/workspace")

    content, _ := afero.ReadFile(fs, "/workspace/.generated/app-one.override.yaml")
    assert.Contains(t, string(content), "dev-internal")
}
```

### Integration Tests

Test Docker interactions with real containers.

```bash
go test -tags=integration ./...
```

### End-to-End Tests

Test full CLI workflows.

```bash
go test -tags=e2e ./...
```

---

## Build & Release

### Local Build

```bash
go build -o contrail ./cmd/contrail
```

### Release Build

```bash
goreleaser release --snapshot --clean
```

### Supported Platforms

- linux/amd64
- linux/arm64
- darwin/amd64
- darwin/arm64
- windows/amd64

---

## Implementation Priority

### Phase 1: Core CLI Structure
1. Root command with context detection
2. Workspace commands (init, up, down, status)
3. App commands (init, up, down, status)
4. Basic configuration loading

### Phase 2: Override Generation
1. Template system for hostnames/aliases
2. Override file generation
3. Manifest generation
4. Traefik label generation

### Phase 3: Shell Integration
1. `compose-prefix` command
2. `init-shell` command with embedded scripts
3. Shell completion for flags

### Phase 4: Polish
1. Port management commands
2. Proxy management commands
3. `doctor` command
4. `validate` command
5. `open` and `urls` commands

### Future: Plugins
1. Define plugin interface in `pkg/plugin`
2. Integrate go-plugin
3. Extract protocol handlers to plugins

---

## Related Documents

- [Architecture Overview](../architecture/overview.md)
- [CLI Reference](../reference/cli.md)
- [Configuration Reference](../reference/configuration.md)

<!-- See appendices/tech-stack/ for complete code examples -->
