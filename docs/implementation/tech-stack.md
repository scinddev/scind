# Contrail Go Stack Specification

<!-- migrated from: specs/contrail-go-stack.md -->

**Version**: 0.1.4-draft
**Date**: December 2024
**Status**: Design Phase

This document defines the Go technology stack for Contrail and provides scaffolding instructions for the initial project structure.

---

## Stack Overview

Contrail's Go stack is intentionally aligned with the Weftlo project where patterns overlap, enabling knowledge transfer and consistent tooling across projects.

### Core Dependencies

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

---

## Dependency Rationale

### Carried from Weftlo (Validated for Contrail)

| Package | Purpose in Contrail |
|---------|---------------------|
| **Cobra** | CLI framework—resource-first command structure, subcommands, persistent flags, shell completion generation |
| **Viper** | Configuration loading—`proxy.yaml`, `workspace.yaml`, `application.yaml`, environment variable overrides, config merging |
| **Afero** | Filesystem abstraction—critical for testing override file generation without touching disk, enables in-memory filesystem for unit tests |
| **Sprig** | Template functions—hostname templates (`{workspace}-{app}-{export}.{domain}`), alias templates, generated override file content |
| **go-playground/validator** | Struct validation—schema validation for workspace and application config files, custom validation rules |
| **testify** | Testing assertions and mocks |
| **yaml.v3** | YAML parsing with support for comments and anchors |

### Contrail-Specific

| Package | Purpose |
|---------|---------|
| **docker/docker** | Official Docker SDK—network creation/inspection, container listing with label filters, event streaming. This is the canonical SDK (same code as Docker CLI). |
| **go-plugin** (future) | HashiCorp's gRPC-based plugin system for protocol handlers. Provides process isolation, crash recovery, and language-agnostic plugins. Battle-tested in Terraform, Vault, Packer. |

### Intentionally Excluded

| Package | Reason |
|---------|--------|
| **compose-go** | Not needed initially—Contrail shells out to `docker compose` rather than parsing compose files programmatically. Add later if direct parsing becomes necessary. |
| **Alternative CLI frameworks** (urfave/cli, Kong) | Cobra has superior ecosystem (completion, doc generation) and precedent (docker, kubectl, gh, terraform). |
| **Alternative plugin systems** | go-plugin is the clear choice for Contrail's requirements. Alternatives solve different problems. |

---

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

Viper's merge capabilities support Contrail's configuration hierarchy:

1. Global config (`~/.config/contrail/proxy.yaml`)
2. Workspace config (`workspace.yaml`)
3. Application config (`application.yaml`)
4. Manual overrides (`overrides/{app}.yaml`)
5. Environment variables (`CONTRAIL_*`)
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

---

## Scaffolding Instructions

### Step 1: Initialize Module

```bash
mkdir contrail && cd contrail
go mod init github.com/yourorg/contrail
```

### Step 2: Install Dependencies

```bash
go get github.com/spf13/cobra@v1.10.2
go get github.com/spf13/viper@v1.21.0
go get github.com/spf13/afero@v1.15.0
go get github.com/Masterminds/sprig/v3@v3.3.0
go get github.com/go-playground/validator/v10@v10.29.0
go get github.com/docker/docker@v27.4.0+incompatible
go get github.com/stretchr/testify@v1.11.1
go get gopkg.in/yaml.v3@v3.0.1
```

### Step 3: Create Directory Structure

```bash
mkdir -p cmd/contrail
mkdir -p internal/{cli,config,context,generator,docker,workspace,output,scripts}
mkdir -p pkg/plugin
mkdir -p testdata/{workspaces,applications}
```

### Step 4: Scaffold Root Command

Create `internal/cli/root.go`: See [scaffold-cmd-root.go](./appendices/tech-stack/scaffold-cmd-root.go)

### Step 5: Scaffold Workspace Commands

Create `internal/cli/workspace.go`: See [scaffold-workspace.go](./appendices/tech-stack/scaffold-workspace.go)

### Step 6: Scaffold App Commands

Create `internal/cli/app.go`: See [scaffold-app.go](./appendices/tech-stack/scaffold-app.go)

### Step 6b: Scaffold Flavor Commands

Create `internal/cli/flavor.go`: See [scaffold-flavor.go](./appendices/tech-stack/scaffold-flavor.go)

### Step 7: Scaffold Top-Level Aliases

Create `internal/cli/aliases.go`: See [scaffold-aliases.go](./appendices/tech-stack/scaffold-aliases.go)

Note: The workspace commands in `workspace.go` should call these same `runWorkspace*` functions to ensure consistent behavior between the full command and alias forms.

### Step 8: Scaffold compose-prefix Command

Create `internal/cli/compose_prefix.go`: See [scaffold-compose-prefix.go](./appendices/tech-stack/scaffold-compose-prefix.go)

### Step 8b: Scaffold Proxy Commands

Create `internal/cli/proxy.go`: See [scaffold-proxy.go](./appendices/tech-stack/scaffold-proxy.go)

### Step 8c: Scaffold Port Commands

Create `internal/cli/port.go`: See [scaffold-port.go](./appendices/tech-stack/scaffold-port.go)

### Step 8d: Scaffold Config Commands

Create `internal/cli/config.go`: See [scaffold-config.go](./appendices/tech-stack/scaffold-config.go)

### Step 8e: Scaffold Utility Commands

Create utility commands:
- `internal/cli/validate.go`: See [scaffold-validate.go](./appendices/tech-stack/scaffold-validate.go)
- `internal/cli/doctor.go`, `internal/cli/open.go`, `internal/cli/urls.go`: See [scaffold-utility.go](./appendices/tech-stack/scaffold-utility.go)

### Step 8f: Scaffold init-shell Command

Create `internal/cli/init_shell.go`: See [scaffold-init-shell.go](./appendices/tech-stack/scaffold-init-shell.go)

Note: The embedded shell scripts (`scripts/bash.sh`, `scripts/zsh.zsh`, `scripts/fish.fish`) contain the shell integration code from the Shell Integration Specification.

### Step 9: Scaffold Config Types

Create configuration types:
- `internal/config/workspace.go`
- `internal/config/application.go`

See [scaffold-config-types.go](./appendices/tech-stack/scaffold-config-types.go)

#### Config Loading: Inference and Defaults

The config loader (`internal/config/loader.go`) applies these inference rules after unmarshaling:

**Service name defaulting** (C-3):
- If `ExportedService.Service` is empty, set it to the map key from `exported_services`
- Example: `exported_services.web` with no `service:` field defaults to Compose service `"web"`

**Port inference** (C-2):
- If `Port.Port` is zero (omitted), infer from the Compose service's `ports:` configuration
- If the Compose service has exactly one port, use that port
- If the Compose service has multiple ports, return a clear error:
  ```
  Error: Port must be specified for exported service "web"
    Application: app-one
    Compose service "web" has multiple ports: 80, 443, 9229
    Specify which port to use in application.yaml
  ```

**Compose file existence validation** (A-10):
- At `generate` time, validate that all files in `Flavor.ComposeFiles` exist on disk
- If a file is missing, return a clear error:
  ```
  Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
    Application: app-two
    Available compose files: docker-compose.yaml, docker-compose.dev.yaml
  ```

### Step 10: Create Entry Point

Create `cmd/contrail/main.go`: See [scaffold-main.go](./appendices/tech-stack/scaffold-main.go)

### Step 11: Build and Verify

```bash
go build -o contrail ./cmd/contrail
./contrail --help
./contrail workspace --help
```

---

## CLI to Cobra Command Mapping

| CLI Command | Cobra Location | Notes |
|-------------|----------------|-------|
| `contrail workspace list` | `workspaceCmd` → `workspaceListCmd` | |
| `contrail workspace show` | `workspaceCmd` → `workspaceShowCmd` | |
| `contrail workspace init` | `workspaceCmd` → `workspaceInitCmd` | |
| `contrail workspace clone` | `workspaceCmd` → `workspaceCloneCmd` | |
| `contrail workspace generate` | `workspaceCmd` → `workspaceGenerateCmd` | |
| `contrail workspace up` | `workspaceCmd` → `workspaceUpCmd` | |
| `contrail workspace down` | `workspaceCmd` → `workspaceDownCmd` | |
| `contrail workspace restart` | `workspaceCmd` → `workspaceRestartCmd` | |
| `contrail workspace status` | `workspaceCmd` → `workspaceStatusCmd` | |
| `contrail app list` | `appCmd` → `appListCmd` | |
| `contrail app show` | `appCmd` → `appShowCmd` | |
| `contrail app init` | `appCmd` → `appInitCmd` | |
| `contrail app add` | `appCmd` → `appAddCmd` | |
| `contrail app remove` | `appCmd` → `appRemoveCmd` | |
| `contrail app up` | `appCmd` → `appUpCmd` | |
| `contrail app down` | `appCmd` → `appDownCmd` | |
| `contrail app restart` | `appCmd` → `appRestartCmd` | |
| `contrail app status` | `appCmd` → `appStatusCmd` | |
| `contrail flavor list` | `flavorCmd` → `flavorListCmd` | |
| `contrail flavor show` | `flavorCmd` → `flavorShowCmd` | |
| `contrail flavor set` | `flavorCmd` → `flavorSetCmd` | |
| `contrail port list` | `portCmd` → `portListCmd` | Global (no context) |
| `contrail port show` | `portCmd` → `portShowCmd` | Global (no context) |
| `contrail port release` | `portCmd` → `portReleaseCmd` | Global (no context) |
| `contrail port assign` | `portCmd` → `portAssignCmd` | Global (no context) |
| `contrail port gc` | `portCmd` → `portGcCmd` | Global (no context) |
| `contrail port scan` | `portCmd` → `portScanCmd` | Global (no context) |
| `contrail proxy init` | `proxyCmd` → `proxyInitCmd` | Global (no context) |
| `contrail proxy status` | `proxyCmd` → `proxyStatusCmd` | Global (no context) |
| `contrail proxy up` | `proxyCmd` → `proxyUpCmd` | Global (no context) |
| `contrail proxy down` | `proxyCmd` → `proxyDownCmd` | Global (no context) |
| `contrail proxy restart` | `proxyCmd` → `proxyRestartCmd` | Global (no context) |
| `contrail config show` | `configCmd` → `configShowCmd` | Global (no context) |
| `contrail config set` | `configCmd` → `configSetCmd` | Global (no context) |
| `contrail config edit` | `configCmd` → `configEditCmd` | Global (no context) |
| `contrail up` | `upCmd` | Alias for `workspace up` |
| `contrail down` | `downCmd` | Alias for `workspace down` |
| `contrail ps` | `psCmd` | Alias for `workspace status` |
| `contrail generate` | `generateCmd` | Alias for `workspace generate` |
| `contrail compose-prefix` | `composePrefixCmd` | Hidden, for shell integration |
| `contrail init-shell` | `initShellCmd` | Outputs shell scripts |
| `contrail completion` | `completionCmd` | Cobra built-in pattern |
| `contrail validate` | `validateCmd` | |
| `contrail doctor` | `doctorCmd` | |
| `contrail open` | `openCmd` | |
| `contrail urls` | `urlsCmd` | |

---

## Testing Strategy

### Unit Tests with Afero

Use Afero's in-memory filesystem for testing file operations:

```go
func TestGenerateOverrideFile(t *testing.T) {
    fs := afero.NewMemMapFs()

    // Create test workspace.yaml
    afero.WriteFile(fs, "/workspace/workspace.yaml", []byte(`...`), 0644)

    // Create test application.yaml
    afero.WriteFile(fs, "/workspace/app-one/application.yaml", []byte(`...`), 0644)

    // Run generator
    gen := generator.New(fs)
    err := gen.Generate("/workspace")

    // Assert override file was created correctly
    content, _ := afero.ReadFile(fs, "/workspace/.generated/app-one.override.yaml")
    assert.Contains(t, string(content), "dev-internal")
}
```

### Integration Tests

Use Docker-in-Docker or testcontainers for integration tests that verify actual Docker Compose behavior.

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

## Related Documentation

- [CLI Reference](../reference/cli.md) - Complete CLI documentation
- [Shell Integration Spec](../specs/shell-integration.md) - Shell scripts to embed
- [Configuration Schemas Spec](../specs/configuration-schemas.md) - Config types to implement

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0-draft | Dec 2024 | Initial Go stack specification |
| 0.1.1-draft | Dec 2024 | Spec review: fixed validation rules, added missing commands, proxy commands, removed logs command |
| 0.1.2-draft | Dec 2024 | Spec review: added proxy up --recreate flag, renamed proxy network to contrail-proxy |
| 0.1.3-draft | Dec 2024 | Spec review: completed issues 19-30 (app exec clarification, repeatable flags, missing commands scaffolding) |
| 0.1.4-draft | Dec 2024 | Spec review: added port, config, utility, and init-shell command scaffolds; fixed alias pattern to use shared helper functions |
| 0.1.5-draft | Dec 2024 | Spec review: docker compose error handling pattern, --color flag, --quiet description, removed --available from port list |
