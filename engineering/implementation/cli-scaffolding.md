# CLI Scaffolding

> **Note**: `yourorg` is used as a placeholder for the GitHub organization or username. Replace with your actual organization name when initializing the project. See [Go Technology Stack](tech-stack.md) for the authoritative dependency list.

This document provides scaffolding instructions for the Cobra-based CLI structure.

## Scaffolding Instructions

### Step 1: Initialize Module

```bash
mkdir scind && cd scind
go mod init github.com/yourorg/scind
```

### Step 2: Install Dependencies

```bash
go get github.com/spf13/cobra@v1.10.2
go get github.com/spf13/viper@v1.21.0
go get github.com/spf13/pflag@v1.0.10
go get github.com/spf13/afero@v1.15.0
go get github.com/Masterminds/sprig/v3@v3.3.0
go get github.com/go-playground/validator/v10@v10.29.0
go get github.com/docker/docker@v27.4.0+incompatible
go get github.com/stretchr/testify@v1.11.1
go get gopkg.in/yaml.v3@v3.0.1
```

### Step 3: Create Directory Structure

```bash
mkdir -p cmd/scind
mkdir -p internal/{cli,config,context,generator,docker,workspace,output,scripts}
mkdir -p pkg/plugin
mkdir -p testdata/{workspaces,applications}
```

### Step 4: Create Entry Point

Create `cmd/scind/main.go`:

```go
package main

import "github.com/yourorg/scind/internal/cli"

func main() {
    cli.Execute()
}
```

### Step 5: Build and Verify

```bash
go build -o scind ./cmd/scind
./scind --help
./scind workspace --help
```

## CLI to Cobra Command Mapping

| CLI Command | Cobra Location | Notes |
|-------------|----------------|-------|
| `scind workspace list` | `workspaceCmd` → `workspaceListCmd` | |
| `scind workspace show` | `workspaceCmd` → `workspaceShowCmd` | |
| `scind workspace init` | `workspaceCmd` → `workspaceInitCmd` | |
| `scind workspace clone` | `workspaceCmd` → `workspaceCloneCmd` | |
| `scind workspace generate` | `workspaceCmd` → `workspaceGenerateCmd` | |
| `scind workspace prune` | `workspaceCmd` → `workspacePruneCmd` | |
| `scind workspace up` | `workspaceCmd` → `workspaceUpCmd` | |
| `scind workspace down` | `workspaceCmd` → `workspaceDownCmd` | |
| `scind workspace restart` | `workspaceCmd` → `workspaceRestartCmd` | |
| `scind workspace status` | `workspaceCmd` → `workspaceStatusCmd` | |
| `scind workspace destroy` | `workspaceCmd` → `workspaceDestroyCmd` | Removes workspace, containers, networks, optionally volumes |
| `scind app list` | `appCmd` → `appListCmd` | |
| `scind app show` | `appCmd` → `appShowCmd` | |
| `scind app init` | `appCmd` → `appInitCmd` | |
| `scind app add` | `appCmd` → `appAddCmd` | |
| `scind app remove` | `appCmd` → `appRemoveCmd` | |
| `scind app up` | `appCmd` → `appUpCmd` | |
| `scind app down` | `appCmd` → `appDownCmd` | |
| `scind app restart` | `appCmd` → `appRestartCmd` | |
| `scind app status` | `appCmd` → `appStatusCmd` | |
| `scind flavor list` | `flavorCmd` → `flavorListCmd` | |
| `scind flavor show` | `flavorCmd` → `flavorShowCmd` | |
| `scind flavor set` | `flavorCmd` → `flavorSetCmd` | |
| `scind port list` | `portCmd` → `portListCmd` | Global (no context) |
| `scind port show` | `portCmd` → `portShowCmd` | Global (no context) |
| `scind port release` | `portCmd` → `portReleaseCmd` | Global (no context) |
| `scind port assign` | `portCmd` → `portAssignCmd` | Global (no context) |
| `scind port gc` | `portCmd` → `portGcCmd` | Global (no context) |
| `scind port scan` | `portCmd` → `portScanCmd` | Global (no context) |
| `scind proxy init` | `proxyCmd` → `proxyInitCmd` | Global (no context) |
| `scind proxy status` | `proxyCmd` → `proxyStatusCmd` | Global (no context) |
| `scind proxy up` | `proxyCmd` → `proxyUpCmd` | Global (no context) |
| `scind proxy down` | `proxyCmd` → `proxyDownCmd` | Global (no context) |
| `scind proxy restart` | `proxyCmd` → `proxyRestartCmd` | Global (no context) |
| `scind config show` | `configCmd` → `configShowCmd` | Global (no context) |
| `scind config get` | `configCmd` → `configGetCmd` | Global (no context) |
| `scind config set` | `configCmd` → `configSetCmd` | Global (no context) |
| `scind config path` | `configCmd` → `configPathCmd` | Global (no context) |
| `scind config edit` | `configCmd` → `configEditCmd` | Global (no context) |
| `scind up` | `upCmd` | Alias for `workspace up` |
| `scind down` | `downCmd` | Alias for `workspace down` |
| `scind ps` | `psCmd` | Alias for `workspace status` |
| `scind generate` | `generateCmd` | Alias for `workspace generate` |
| `scind compose-prefix` | `composePrefixCmd` | Hidden, for shell integration |
| `scind init-shell` | `initShellCmd` | Outputs shell scripts |
| `scind completion` | `completionCmd` | Cobra built-in pattern |
| `scind validate` | `validateCmd` | |
| `scind doctor` | `doctorCmd` | |
| `scind open` | `openCmd` | |
| `scind urls` | `urlsCmd` | |

## Implementation Priority

### Phase 1: Core CLI Structure
1. Root command with context detection
2. Workspace commands (init, up, down, destroy, status)
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

## Related Documents

- [Go Technology Stack](tech-stack.md) — Authoritative dependency list and architecture patterns
- [Project Layout](project-layout.md) — Directory structure and file organization
- [CLI Reference](../reference/cli.md) — Command documentation
- [Shell Integration](../specs/shell-integration.md) — Shell function and completion specifications
