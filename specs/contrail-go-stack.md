# Contrail Go Stack Specification

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

Create `internal/cli/root.go`:

```go
package cli

import (
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile   string
    workspace string
    apps      []string
    quiet     bool
    verbose   bool
    jsonOut   bool
    yamlOut   bool
    colorMode string  // auto, always, never
)

var rootCmd = &cobra.Command{
    Use:   "contrail",
    Short: "Workspace orchestration for Docker Compose",
    Long: `Contrail is a workspace orchestration system for Docker Compose that enables
developers to run multiple isolated instances of multi-application stacks
simultaneously on a single host.`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Skip context detection for global commands
        if isGlobalCommand(cmd) {
            return nil
        }
        return detectAndSetContext(cmd)
    },
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Global flags
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/contrail/proxy.yaml)")
    rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "specify workspace (overrides context detection)")
    rootCmd.PersistentFlags().StringSliceVarP(&apps, "app", "a", nil, "specify application(s) (repeatable, overrides context detection)")
    rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output, suppress context indicators and progress")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output")
    rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output in JSON format")
    rootCmd.PersistentFlags().BoolVar(&yamlOut, "yaml", false, "output in YAML format")
    rootCmd.PersistentFlags().StringVar(&colorMode, "color", "auto", "color output: auto, always, or never")

    // Register flag completion
    rootCmd.RegisterFlagCompletionFunc("workspace", completeWorkspace)
    rootCmd.RegisterFlagCompletionFunc("app", completeApp)
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        viper.AddConfigPath(home + "/.config/contrail")
        viper.SetConfigType("yaml")
        viper.SetConfigName("proxy")
    }

    viper.SetEnvPrefix("CONTRAIL")
    viper.AutomaticEnv()

    viper.ReadInConfig()
}

// isGlobalCommand returns true for commands that don't use context detection
func isGlobalCommand(cmd *cobra.Command) bool {
    name := cmd.Name()
    return name == "port" || name == "proxy" || name == "config" ||
           name == "doctor" || name == "completion" || name == "help" ||
           name == "init-shell"
}

func detectAndSetContext(cmd *cobra.Command) error {
    // Implementation: walk up directory tree using workspace boundary approach
    // See internal/context/detector.go for full implementation
    //
    // Algorithm (workspace boundary):
    // 1. Walk up from CWD looking for workspace.yaml → establishes workspace root
    // 2. Walk up from CWD toward workspace root looking for application.yaml
    //    - Only consider application.yaml files WITHIN the workspace directory tree
    //    - Never traverse above workspace root (prevents vendor hijacking)
    // 3. Set results in viper:
    //    - viper.Set("detected.workspace", workspaceName)
    //    - viper.Set("detected.workspace_path", workspaceRoot)
    //    - viper.Set("detected.app", appName)
    //    - viper.Set("detected.app_path", appPath)
    //
    // Error messages for debugging:
    // - If application.yaml found but no workspace.yaml:
    //   "No workspace found (workspace.yaml) in current directory or any parent directories,
    //    but found an application (application.yaml) at: {path}"
    // - If neither found:
    //   "No workspace found (workspace.yaml) in current directory or any parent directories,
    //    and no application (application.yaml) found either"
    return nil
}

func completeWorkspace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    // List available workspaces
    return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func completeApp(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    // List available apps in current workspace
    return []string{}, cobra.ShellCompDirectiveNoFileComp
}
```

### Step 5: Scaffold Workspace Commands

Create `internal/cli/workspace.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
    Use:   "workspace",
    Short: "Manage workspaces",
    Long:  `Create, configure, and manage workspace environments.`,
}

var workspaceListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all workspaces",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show workspace details",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceCloneCmd = &cobra.Command{
    Use:   "clone",
    Short: "Clone application repositories",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceGenerateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate override files",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspacePruneCmd = &cobra.Command{
    Use:   "prune",
    Short: "Remove stale workspace registry entries",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Bring up a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Tear down a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceRestartCmd = &cobra.Command{
    Use:   "restart",
    Short: "Restart a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var workspaceStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show workspace status",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(workspaceCmd)

    workspaceCmd.AddCommand(workspaceListCmd)
    workspaceCmd.AddCommand(workspaceShowCmd)
    workspaceCmd.AddCommand(workspaceInitCmd)
    workspaceCmd.AddCommand(workspaceCloneCmd)
    workspaceCmd.AddCommand(workspaceGenerateCmd)
    workspaceCmd.AddCommand(workspacePruneCmd)
    workspaceCmd.AddCommand(workspaceUpCmd)
    workspaceCmd.AddCommand(workspaceDownCmd)
    workspaceCmd.AddCommand(workspaceRestartCmd)
    workspaceCmd.AddCommand(workspaceStatusCmd)

    // workspace list flags
    workspaceListCmd.Flags().Bool("validate", false, "check that registered paths still contain workspace.yaml")
    workspaceListCmd.Flags().Bool("rebuild", false, "rebuild registry from Docker labels")

    // workspace init flags
    workspaceInitCmd.Flags().String("path", "", "directory to create workspace in")

    // workspace up flags
    workspaceUpCmd.Flags().Bool("no-generate", false, "skip automatic regeneration")
    workspaceUpCmd.Flags().Bool("force-generate", false, "force regeneration even if up-to-date")
    workspaceUpCmd.Flags().BoolP("detach", "d", true, "run in background")

    // workspace down flags
    workspaceDownCmd.Flags().Bool("volumes", false, "also remove volumes")
    workspaceDownCmd.Flags().Bool("force", false, "skip confirmation")

    // workspace generate flags
    workspaceGenerateCmd.Flags().Bool("force", false, "regenerate even if up-to-date")

    // workspace prune flags
    workspacePruneCmd.Flags().Bool("dry-run", false, "show what would be removed without making changes")
}
```

### Step 6: Scaffold App Commands

Create `internal/cli/app.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
    Use:   "app",
    Short: "Manage applications",
    Long:  `Manage applications within workspaces.`,
}

var appListCmd = &cobra.Command{
    Use:   "list",
    Short: "List applications in a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show application details",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize an application configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appAddCmd = &cobra.Command{
    Use:   "add",
    Short: "Add an application to a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appRemoveCmd = &cobra.Command{
    Use:   "remove",
    Short: "Remove an application from a workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Bring up an application",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Tear down an application",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appRestartCmd = &cobra.Command{
    Use:   "restart",
    Short: "Restart an application",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var appStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show application status",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(appCmd)

    appCmd.AddCommand(appListCmd)
    appCmd.AddCommand(appShowCmd)
    appCmd.AddCommand(appInitCmd)
    appCmd.AddCommand(appAddCmd)
    appCmd.AddCommand(appRemoveCmd)
    appCmd.AddCommand(appUpCmd)
    appCmd.AddCommand(appDownCmd)
    appCmd.AddCommand(appRestartCmd)
    appCmd.AddCommand(appStatusCmd)

    // app init flags
    appInitCmd.Flags().StringP("app", "a", "", "application name (default: current directory name)")

    // app add flags
    appAddCmd.Flags().StringP("app", "a", "", "application name (required)")
    appAddCmd.Flags().String("repo", "", "git repository URL to clone")
    appAddCmd.Flags().String("path", "", "custom path relative to workspace")
    appAddCmd.MarkFlagRequired("app")

    // app remove flags
    appRemoveCmd.Flags().Bool("force", false, "skip confirmation, also remove directory")

    // app down flags
    appDownCmd.Flags().Bool("volumes", false, "also remove volumes")
}
```

### Step 6b: Scaffold Flavor Commands

Create `internal/cli/flavor.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var flavorCmd = &cobra.Command{
    Use:   "flavor",
    Short: "Manage application flavors",
    Long:  `Manage application flavors (named configurations).`,
}

var flavorListCmd = &cobra.Command{
    Use:   "list",
    Short: "List available flavors",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var flavorShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show current active flavor",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var flavorSetCmd = &cobra.Command{
    Use:   "set <flavor>",
    Short: "Set the active flavor",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // args[0] contains the flavor name
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(flavorCmd)

    flavorCmd.AddCommand(flavorListCmd)
    flavorCmd.AddCommand(flavorShowCmd)
    flavorCmd.AddCommand(flavorSetCmd)

    // Note: -w/--workspace and -a/--app flags are inherited from root command
}
```

### Step 7: Scaffold Top-Level Aliases

Create `internal/cli/aliases.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

// Top-level aliases for common operations
// These call shared implementation functions to ensure proper context handling
var upCmd = &cobra.Command{
    Use:   "up",
    Short: "Alias for 'workspace up'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceUp(cmd, args)
    },
}

var downCmd = &cobra.Command{
    Use:   "down",
    Short: "Alias for 'workspace down'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceDown(cmd, args)
    },
}

var psCmd = &cobra.Command{
    Use:   "ps",
    Short: "Alias for 'workspace status'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceStatus(cmd, args)
    },
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Alias for 'workspace generate'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceGenerate(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(upCmd)
    rootCmd.AddCommand(downCmd)
    rootCmd.AddCommand(psCmd)
    rootCmd.AddCommand(generateCmd)

    // Copy flags from workspace commands to aliases
    upCmd.Flags().AddFlagSet(workspaceUpCmd.Flags())
    downCmd.Flags().AddFlagSet(workspaceDownCmd.Flags())
    generateCmd.Flags().AddFlagSet(workspaceGenerateCmd.Flags())
}

// Shared implementation functions (called by both workspace commands and aliases)
// These are defined here but the actual implementation lives in workspace.go

func runWorkspaceUp(cmd *cobra.Command, args []string) error {
    // Implementation: generate overrides, start containers
    return nil
}

func runWorkspaceDown(cmd *cobra.Command, args []string) error {
    // Implementation: stop containers, optionally remove volumes
    return nil
}

func runWorkspaceStatus(cmd *cobra.Command, args []string) error {
    // Implementation: aggregate status from all apps
    return nil
}

func runWorkspaceGenerate(cmd *cobra.Command, args []string) error {
    // Implementation: generate override files
    return nil
}
```

Note: The workspace commands in `workspace.go` should call these same `runWorkspace*` functions to ensure consistent behavior between the full command and alias forms.

### Step 8: Scaffold compose-prefix Command

Create `internal/cli/compose_prefix.go`:

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var composePrefixCmd = &cobra.Command{
    Use:    "compose-prefix",
    Short:  "Output docker compose command prefix",
    Long:   `Outputs a docker compose command prefix with project name and compose files for the current context.`,
    Hidden: true, // Internal command for shell integration
    RunE: func(cmd *cobra.Command, args []string) error {
        ws := viper.GetString("resolved.workspace")
        app := viper.GetString("resolved.app")

        if ws == "" || app == "" {
            return fmt.Errorf("no application context detected")
        }

        projectName := fmt.Sprintf("%s-%s", ws, app)
        
        // Get compose files from resolved flavor
        composeFiles, err := getComposeFilesForApp(ws, app)
        if err != nil {
            return err
        }

        // Build output
        fmt.Printf("docker compose -p %s", projectName)
        for _, f := range composeFiles {
            fmt.Printf(" -f '%s'", f)
        }
        fmt.Println()

        return nil
    },
}

func init() {
    rootCmd.AddCommand(composePrefixCmd)
    // Note: No --flavor flag. Flavor changes require regeneration and can impact
    // running applications. Users must use `contrail flavor set` instead.
}

func getComposeFilesForApp(workspace, app string) ([]string, error) {
    // Implementation: read application.yaml, resolve flavor, return file list
    return []string{}, nil
}
```

### Step 8b: Scaffold Proxy Commands

Create `internal/cli/proxy.go`:

```go
package cli

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
    Use:   "proxy",
    Short: "Manage the Traefik reverse proxy",
}

var proxyInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Bootstrap proxy configuration",
    Long:  `Creates the Traefik Docker Compose project at ~/.config/contrail/proxy/`,
    RunE: func(cmd *cobra.Command, args []string) error {
        force, _ := cmd.Flags().GetBool("force")
        domain, _ := cmd.Flags().GetString("domain")
        path, _ := cmd.Flags().GetString("path")

        // Check if proxy config exists
        if _, err := os.Stat(filepath.Join(path, "docker-compose.yaml")); err == nil {
            if !force {
                return fmt.Errorf("proxy configuration already exists at %s\nUse --force to overwrite", path)
            }
            // Backup existing config
            // ...
        }

        // Create directory structure
        // Create docker-compose.yaml, traefik.yaml, dynamic/, certs/
        // Create contrail-proxy network if needed
        // Output next steps

        return nil
    },
}

var proxyUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Start the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        recreate, _ := cmd.Flags().GetBool("recreate")

        // Run proxy init if config doesn't exist
        // Create contrail-proxy network if needed
        // If recreate flag is set, remove and recreate the network
        // Validate existing network configuration
        // Start containers via docker compose
        _ = recreate // TODO: implement
        return nil
    },
}

var proxyDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Stop the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Stop containers via docker compose
        return nil
    },
}

var proxyRestartCmd = &cobra.Command{
    Use:   "restart",
    Short: "Restart the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Restart via docker compose
        return nil
    },
}

var proxyStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show proxy status",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Check container status, network, entrypoints
        return nil
    },
}

func init() {
    rootCmd.AddCommand(proxyCmd)
    proxyCmd.AddCommand(proxyInitCmd)
    proxyCmd.AddCommand(proxyUpCmd)
    proxyCmd.AddCommand(proxyDownCmd)
    proxyCmd.AddCommand(proxyRestartCmd)
    proxyCmd.AddCommand(proxyStatusCmd)

    // proxy init flags
    proxyInitCmd.Flags().Bool("force", false, "overwrite existing configuration")
    proxyInitCmd.Flags().String("domain", "contrail.test", "proxy domain for generated hostnames")
    proxyInitCmd.Flags().String("path", defaultProxyPath(), "directory to create proxy in")

    // proxy up flags
    proxyUpCmd.Flags().Bool("recreate", false, "recreate the proxy network even if it exists")
}

func defaultProxyPath() string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "contrail", "proxy")
}
```

### Step 8c: Scaffold Port Commands

Create `internal/cli/port.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
    Use:   "port",
    Short: "Manage port assignments",
    Long:  `View and manage port assignments across all workspaces.`,
}

var portListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all port assignments",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var portShowCmd = &cobra.Command{
    Use:   "show <port>",
    Short: "Show details for a specific port",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // args[0] contains the port number
        return nil
    },
}

var portReleaseCmd = &cobra.Command{
    Use:   "release <port>",
    Short: "Release a port assignment",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var portAssignCmd = &cobra.Command{
    Use:   "assign",
    Short: "Manually assign a port",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var portGcCmd = &cobra.Command{
    Use:   "gc",
    Short: "Garbage collect stale port assignments",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var portScanCmd = &cobra.Command{
    Use:   "scan",
    Short: "Scan for port conflicts",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(portCmd)

    portCmd.AddCommand(portListCmd)
    portCmd.AddCommand(portShowCmd)
    portCmd.AddCommand(portReleaseCmd)
    portCmd.AddCommand(portAssignCmd)
    portCmd.AddCommand(portGcCmd)
    portCmd.AddCommand(portScanCmd)

    // port release flags
    portReleaseCmd.Flags().Bool("force", false, "release even if container is running")

    // port assign flags
    portAssignCmd.Flags().Int("port", 0, "specific port to assign (required)")
    portAssignCmd.Flags().StringP("workspace", "w", "", "workspace name (required)")
    portAssignCmd.Flags().StringP("app", "a", "", "application name (required)")
    portAssignCmd.Flags().String("service", "", "service name (required)")
    portAssignCmd.MarkFlagRequired("port")
    portAssignCmd.MarkFlagRequired("workspace")
    portAssignCmd.MarkFlagRequired("app")
    portAssignCmd.MarkFlagRequired("service")

    // port gc flags
    portGcCmd.Flags().Bool("dry-run", false, "show what would be released without making changes")

    // port scan flags
    portScanCmd.Flags().Bool("fix", false, "attempt to resolve conflicts automatically")
}
```

### Step 8d: Scaffold Config Commands

Create `internal/cli/config.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage Contrail configuration",
    Long:  `View and modify Contrail configuration settings.`,
}

var configShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show current configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var configGetCmd = &cobra.Command{
    Use:   "get <key>",
    Short: "Get a configuration value",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // args[0] contains the key
        return nil
    },
}

var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Set a configuration value",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        // args[0] = key, args[1] = value
        return nil
    },
}

var configPathCmd = &cobra.Command{
    Use:   "path",
    Short: "Show configuration file paths",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

var configEditCmd = &cobra.Command{
    Use:   "edit",
    Short: "Open configuration in editor",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: open $EDITOR or default editor
        return nil
    },
}

func init() {
    rootCmd.AddCommand(configCmd)

    configCmd.AddCommand(configShowCmd)
    configCmd.AddCommand(configGetCmd)
    configCmd.AddCommand(configSetCmd)
    configCmd.AddCommand(configPathCmd)
    configCmd.AddCommand(configEditCmd)

    // config show flags
    configShowCmd.Flags().Bool("resolved", false, "show fully resolved configuration with all defaults")

    // config edit flags
    configEditCmd.Flags().String("file", "proxy", "which config to edit: proxy, workspace, or application")
}
```

### Step 8e: Scaffold Utility Commands

Create `internal/cli/validate.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate configuration files",
    Long:  `Validate workspace.yaml and application.yaml files for correctness.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: validate schemas, check references
        return nil
    },
}

func init() {
    rootCmd.AddCommand(validateCmd)

    validateCmd.Flags().Bool("strict", false, "treat warnings as errors")
}
```

Create `internal/cli/doctor.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check system health and dependencies",
    Long:  `Verify that Docker, Docker Compose, and other dependencies are properly configured.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: check Docker, network, proxy status
        return nil
    },
}

func init() {
    rootCmd.AddCommand(doctorCmd)

    doctorCmd.Flags().Bool("fix", false, "attempt to fix issues automatically")
}
```

Create `internal/cli/open.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
    Use:   "open [service]",
    Short: "Open service URL in browser",
    Long:  `Open the URL for a proxied service in the default browser.`,
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: resolve URL, open browser
        return nil
    },
}

func init() {
    rootCmd.AddCommand(openCmd)

    openCmd.Flags().Bool("print", false, "print URL instead of opening browser")
}
```

Create `internal/cli/urls.go`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var urlsCmd = &cobra.Command{
    Use:   "urls",
    Short: "List all service URLs",
    Long:  `Display URLs for all proxied services in the current context.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: list all proxied URLs
        return nil
    },
}

func init() {
    rootCmd.AddCommand(urlsCmd)
}
```

### Step 8f: Scaffold init-shell Command

Create `internal/cli/init_shell.go`:

```go
package cli

import (
    _ "embed"
    "fmt"

    "github.com/spf13/cobra"
)

//go:embed scripts/bash.sh
var bashScript string

//go:embed scripts/zsh.zsh
var zshScript string

//go:embed scripts/fish.fish
var fishScript string

var initShellCmd = &cobra.Command{
    Use:   "init-shell <shell>",
    Short: "Output shell integration script",
    Long: `Output shell integration script for the specified shell.

Supported shells: bash, zsh, fish

Add to your shell configuration:
  # Bash
  contrail init-shell bash >> ~/.bashrc

  # Zsh
  contrail init-shell zsh >> ~/.zshrc

  # Fish
  contrail init-shell fish >> ~/.config/fish/conf.d/contrail.fish`,
    Args:      cobra.ExactArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish"},
    RunE: func(cmd *cobra.Command, args []string) error {
        shell := args[0]
        switch shell {
        case "bash":
            fmt.Print(bashScript)
        case "zsh":
            fmt.Print(zshScript)
        case "fish":
            fmt.Print(fishScript)
        default:
            return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initShellCmd)
}
```

Note: The embedded shell scripts (`scripts/bash.sh`, `scripts/zsh.zsh`, `scripts/fish.fish`) contain the shell integration code from the [Shell Integration Specification](./contrail-shell-integration.md).

### Step 9: Scaffold Config Types

Create `internal/config/workspace.go`:

```go
package config

// WorkspaceConfig represents workspace.yaml
type WorkspaceConfig struct {
    Workspace WorkspaceSpec `yaml:"workspace" validate:"required"`
}

type WorkspaceSpec struct {
    Name         string                      `yaml:"name" validate:"required,dns_label"`
    Domain       string                      `yaml:"domain,omitempty"`
    Templates    *TemplateOverrides          `yaml:"templates,omitempty"`
    Applications map[string]ApplicationRef   `yaml:"applications,omitempty"`
}

type TemplateOverrides struct {
    Hostname    string `yaml:"hostname,omitempty"`
    Alias       string `yaml:"alias,omitempty"`
    ProjectName string `yaml:"project_name,omitempty"`
}

type ApplicationRef struct {
    Path       string `yaml:"path,omitempty"`
    Repository string `yaml:"repository,omitempty"`
}
```

Create `internal/config/application.go`:

```go
package config

// ApplicationConfig represents application.yaml
type ApplicationConfig struct {
    ExportedServices map[string]ExportedService `yaml:"exported_services" validate:"required"`
    Flavors          map[string]Flavor          `yaml:"flavors,omitempty"`
    DefaultFlavor    string                     `yaml:"default_flavor,omitempty"`
}

type ExportedService struct {
    Service string `yaml:"service,omitempty" validate:"omitempty"`  // Optional: defaults to map key
    Ports   []Port `yaml:"ports" validate:"required,dive"`
}

type Port struct {
    Type       string `yaml:"type" validate:"required,oneof=proxied assigned"`
    Protocol   string `yaml:"protocol,omitempty" validate:"required_if=Type proxied,omitempty,oneof=http https tcp postgresql mysql"`
    Port       int    `yaml:"port,omitempty" validate:"omitempty,min=1,max=65535"`  // Optional: inferred from Compose service
    Visibility string `yaml:"visibility,omitempty" validate:"omitempty,oneof=public protected"`  // Defaults to "protected"
}

type Flavor struct {
    ComposeFiles []string `yaml:"compose_files" validate:"required,min=1"`
}
```

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

Create `cmd/contrail/main.go`:

```go
package main

import "github.com/yourorg/contrail/internal/cli"

func main() {
    cli.Execute()
}
```

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

- [Contrail PRD](./contrail-prd.md) — Product requirements and concepts
- [Contrail Technical Specification](./contrail-technical-spec.md) — Architecture and configuration schemas
- [Contrail CLI Reference](./contrail-cli-reference.md) — Complete CLI documentation
- [Contrail Shell Integration](./contrail-shell-integration.md) — Shell functions and completion

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
