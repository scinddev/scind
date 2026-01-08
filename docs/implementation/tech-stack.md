# Go Technology Stack

> **Note**: Throughout this document and related implementation files, `yourorg` is used as a placeholder for the GitHub organization or username. Replace with your actual organization name when initializing the project (e.g., `github.com/mycompany/scind`).

## Stack Overview

Scind's Go stack is intentionally aligned with the Weftlo project where patterns overlap, enabling knowledge transfer and consistent tooling across projects.

### Core Dependencies

```go
// go.mod
module github.com/yourorg/scind

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

## Dependency Rationale

### Carried from Weftlo (Validated for Scind)

| Package | Purpose in Scind |
|---------|---------------------|
| **Cobra** | CLI framework—resource-first command structure, subcommands, persistent flags, shell completion generation |
| **Viper** | Configuration loading—`proxy.yaml`, `workspace.yaml`, `application.yaml`, environment variable overrides, config merging |
| **Afero** | Filesystem abstraction—critical for testing override file generation without touching disk, enables in-memory filesystem for unit tests |
| **Sprig** | Template functions—hostname templates (`{workspace}-{app}-{export}.{domain}`), alias templates, generated override file content |
| **go-playground/validator** | Struct validation—schema validation for workspace and application config files, custom validation rules |
| **testify** | Testing assertions and mocks |
| **yaml.v3** | YAML parsing with support for comments and anchors |

### Scind-Specific

| Package | Purpose |
|---------|---------|
| **docker/docker** | Official Docker SDK—network creation/inspection, container listing with label filters, event streaming. This is the canonical SDK (same code as Docker CLI). |
| **go-plugin** (future) | HashiCorp's gRPC-based plugin system for protocol handlers. Provides process isolation, crash recovery, and language-agnostic plugins. Battle-tested in Terraform, Vault, Packer. |

### External Dependencies

| Component | Version | Purpose | ADR |
|-----------|---------|---------|-----|
| **Docker** | >= 24.0 | Container runtime with Compose V2 | [ADR-0001](../decisions/0001-docker-compose-project-name-isolation.md) |
| **Docker Compose** | >= 2.20 | Multi-container orchestration | [ADR-0003](../decisions/0003-pure-overlay-design.md) |
| **Traefik** | v3.2.3 | Reverse proxy and TLS termination | [ADR-0008](../decisions/0008-traefik-reverse-proxy.md) |

### Intentionally Excluded

| Package | Reason |
|---------|--------|
| **compose-go** | Not needed initially—Scind shells out to `docker compose` rather than parsing compose files programmatically. Add later if direct parsing becomes necessary. |
| **Alternative CLI frameworks** (urfave/cli, Kong) | Cobra has superior ecosystem (completion, doc generation) and precedent (docker, kubectl, gh, terraform). |
| **Alternative plugin systems** | go-plugin is the clear choice for Scind's requirements. Alternatives solve different problems. |

## Architecture Patterns

### Context Detection

Use Cobra's `PersistentPreRunE` on the root command to implement directory-walking context detection:

```go
// Detect workspace.yaml and application.yaml by walking up the directory tree
// Store resolved context in Viper for access by all subcommands
var rootCmd = &cobra.Command{
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Skip context detection for global commands
        if isGlobalCommand(cmd) {
            return nil
        }
        return detectContext(cmd)
    },
}
```

### Configuration Layering

Viper's merge capabilities support Scind's configuration hierarchy:

1. Global config (`~/.config/scind/proxy.yaml`)
2. Workspace config (`workspace.yaml`)
3. Application config (`application.yaml`)
4. Manual overrides (`overrides/{app}.yaml`)
5. Environment variables (`SCIND_*`)
6. Command-line flags

### Shell Completion

Cobra's `ValidArgsFunction` enables dynamic completion for `--workspace` and `--app` flags:

```go
cmd.RegisterFlagCompletionFunc("workspace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    workspaces, _ := listWorkspaces()
    return workspaces, cobra.ShellCompDirectiveNoFileComp
})
```

### Shell Integration Scripts

The `init-shell` command outputs static shell scripts. Use `//go:embed` for the script content:

```go
//go:embed scripts/bash.sh
var bashScript string

//go:embed scripts/zsh.zsh
var zshScript string

//go:embed scripts/fish.fish
var fishScript string
```

### Docker Interaction Strategy

**Primary approach**: Shell out to `docker compose` via `exec.Command`. This is simpler, maintains full Docker Compose compatibility, and is what the specs describe.

**Error handling**: Use pass-through with context prefix. Capture both stdout and stderr from docker compose, and on non-zero exit:
1. Print a context line: `"Failed to start {app-name}:"` or `"Failed to stop {app-name}:"`
2. Print the full docker compose output (stdout + stderr) unmodified
3. Return an appropriate exit code

```go
// Example: Running docker compose with error handling
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

```go
// Example: Check if port is available
func isPortAvailable(port int) bool {
    cli, _ := client.NewClientWithOpts(client.FromEnv)
    containers, _ := cli.ContainerList(ctx, container.ListOptions{})
    // Check port bindings...
}
```

## Testing Strategy

### Unit Tests with Afero

Use Afero's in-memory filesystem for testing file operations:

```go
func TestGenerateOverrideFile(t *testing.T) {
    fs := afero.NewMemMapFs()

    // Create test workspace.yaml
    afero.WriteFile(fs, "/workspace/workspace.yaml", []byte(`...`), 0644)

    // Create test application.yaml
    afero.WriteFile(fs, "/workspace/frontend/application.yaml", []byte(`...`), 0644)

    // Run generator
    gen := generator.New(fs)
    err := gen.Generate("/workspace")

    // Assert override file was created correctly
    content, _ := afero.ReadFile(fs, "/workspace/.generated/frontend.override.yaml")
    assert.Contains(t, string(content), "dev-internal")
}
```

### Integration Tests

Use Docker-in-Docker or testcontainers for integration tests that verify actual Docker Compose behavior.

## Version Constraints

| Dependency | Minimum | Reason |
|------------|---------|--------|
| Go | 1.21 | Required for `log/slog` structured logging |
| Docker | 24.0 | Compose V2 integration as `docker compose` subcommand |
| Docker Compose | 2.20 | Required compose features (project isolation, label support) |
| Traefik | 3.0 | V3 API and configuration format |

## Version Handling

Version information is injected at build time via ldflags:

```bash
go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" ./cmd/scind
```

The goreleaser configuration handles this automatically for releases. See [goreleaser.yaml](appendices/tech-stack/goreleaser.yaml) for the complete release configuration.

## Related Documents

- [CLI Scaffolding](cli-scaffolding.md) — Cobra command structure and scaffolding instructions
- [Project Layout](project-layout.md) — Directory structure and file organization
- [CLI Reference](../reference/cli.md) — Command documentation
- [Configuration Reference](../reference/configuration.md) — Configuration options
