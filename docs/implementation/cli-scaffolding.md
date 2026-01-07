<!-- Migrated from specs/contrail-go-stack.md:269-1525 -->
<!-- Extraction ID: impl-cobra-structure -->

# CLI Scaffolding

This document provides scaffolding instructions for the Cobra-based CLI structure.

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

### Step 4: Create Entry Point

Create `cmd/contrail/main.go`:

```go
package main

import "github.com/yourorg/contrail/internal/cli"

func main() {
    cli.Execute()
}
```

### Step 5: Build and Verify

```bash
go build -o contrail ./cmd/contrail
./contrail --help
./contrail workspace --help
```

## CLI to Cobra Command Mapping

| CLI Command | Cobra Location | Notes |
|-------------|----------------|-------|
| `contrail workspace list` | `workspaceCmd` → `workspaceListCmd` | |
| `contrail workspace show` | `workspaceCmd` → `workspaceShowCmd` | |
| `contrail workspace init` | `workspaceCmd` → `workspaceInitCmd` | |
| `contrail workspace clone` | `workspaceCmd` → `workspaceCloneCmd` | |
| `contrail workspace generate` | `workspaceCmd` → `workspaceGenerateCmd` | |
| `contrail workspace prune` | `workspaceCmd` → `workspacePruneCmd` | |
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
| `contrail config get` | `configCmd` → `configGetCmd` | Global (no context) |
| `contrail config set` | `configCmd` → `configSetCmd` | Global (no context) |
| `contrail config path` | `configCmd` → `configPathCmd` | Global (no context) |
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
