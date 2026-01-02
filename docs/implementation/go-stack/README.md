# Contrail Go Stack

**Version**: 0.1.5
**Date**: December 2024
**Status**: Design Phase

---

## Stack Overview

Contrail's Go stack is intentionally aligned with the Weftlo project where patterns overlap.

---

## Core Dependencies

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

---

## Dependency Rationale

### Carried from Weftlo

| Package | Purpose in Contrail |
|---------|---------------------|
| **Cobra** | CLI framework—resource-first commands, shell completion |
| **Viper** | Configuration loading and merging |
| **Afero** | Filesystem abstraction for testing |
| **Sprig** | Template functions for hostname generation |
| **go-playground/validator** | Schema validation |
| **testify** | Testing assertions and mocks |
| **yaml.v3** | YAML parsing with comments |

### Contrail-Specific

| Package | Purpose |
|---------|---------|
| **docker/docker** | Official Docker SDK for network/container operations |

### Future Dependencies

```go
// Plugin system - add when implementing protocol plugins
github.com/hashicorp/go-plugin v1.6.2

// Compose file parsing - add only if programmatic parsing needed
github.com/compose-spec/compose-go/v2 v2.4.0
```

---

## Project Structure

```
contrail/
├── cmd/contrail/
│   └── main.go                 # Entry point
│
├── internal/
│   ├── cli/                    # Cobra command definitions
│   │   ├── root.go
│   │   ├── workspace.go
│   │   ├── app.go
│   │   ├── flavor.go
│   │   ├── port.go
│   │   ├── proxy.go
│   │   ├── config.go
│   │   ├── aliases.go
│   │   ├── compose_prefix.go
│   │   ├── init_shell.go
│   │   └── ...
│   │
│   ├── config/                 # Configuration types and loading
│   │   ├── workspace.go
│   │   ├── application.go
│   │   ├── proxy.go
│   │   ├── state.go
│   │   └── loader.go
│   │
│   ├── context/                # Context detection
│   │   ├── detector.go
│   │   └── resolver.go
│   │
│   ├── generator/              # Override file generation
│   │   ├── generator.go
│   │   ├── templates.go
│   │   ├── override.go
│   │   └── traefik.go
│   │
│   ├── docker/                 # Docker interaction
│   │   ├── compose.go
│   │   ├── network.go
│   │   ├── container.go
│   │   └── port.go
│   │
│   ├── workspace/              # Workspace operations
│   │   ├── lifecycle.go
│   │   ├── status.go
│   │   └── clone.go
│   │
│   ├── output/                 # Output formatting
│   │   ├── table.go
│   │   ├── json.go
│   │   ├── yaml.go
│   │   └── quiet.go
│   │
│   └── scripts/                # Embedded shell scripts
│       ├── bash.sh
│       ├── zsh.zsh
│       └── fish.fish
│
├── pkg/plugin/                 # Public API for future plugins
│   ├── protocol.go
│   └── types.go
│
├── testdata/                   # Test fixtures
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Architecture Patterns

### Context Detection

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

### Configuration Layering

Viper's merge capabilities support:
1. Global config (`~/.config/contrail/proxy.yaml`)
2. Workspace config (`workspace.yaml`)
3. Application config (`application.yaml`)
4. Manual overrides (`overrides/{app}.yaml`)
5. Environment variables (`CONTRAIL_*`)
6. Command-line flags

### Docker Interaction Strategy

**Primary approach**: Shell out to `docker compose` via `exec.Command`.

**Error handling**:
```go
func runCompose(appName string, args ...string) error {
    cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to %s %s:\n", args[0], appName)
        fmt.Fprint(os.Stderr, string(output))
        return err
    }
    return nil
}
```

**Docker SDK usage**: For operations that don't map to compose commands (network creation, container listing with labels, port checking).

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

---

## Testing Strategy

### Unit Tests with Afero

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

---

## Related Documentation

- [CLI Reference](../../reference/cli/README.md)
- [Configuration Schemas Spec](../../specs/configuration-schemas/README.md)
- [Shell Integration Spec](../../specs/shell-integration/README.md)
