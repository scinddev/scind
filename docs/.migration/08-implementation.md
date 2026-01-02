# Migration Step: Layer 7 — Implementation

**Prerequisites**: Read `common-instructions.md`, complete Layers 4-5
**Estimated Size**: 1 file + appendices, approximately 1,700 lines total

---

## Overview

Extract implementation documentation from the Go stack specification. This includes dependency decisions, code patterns, and scaffolding.

**Source document**: `specs/contrail-go-stack.md` (entire file, 1,615 lines)

---

## File: `implementation/tech-stack.md`

**Sources**:
- `specs/contrail-go-stack.md:1-58` (Overview, Philosophy)
- `specs/contrail-go-stack.md:86-400` (Core Dependencies)
- `specs/contrail-go-stack.md:400-800` (Patterns)
- `specs/contrail-go-stack.md:800-1200` (Code Structure)

### Content Structure

```markdown
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

## Project Structure

```
contrail/
├── cmd/                    # CLI commands
│   ├── root.go
│   ├── up.go
│   ├── down.go
│   └── ...
├── internal/               # Private packages
│   ├── config/             # Configuration loading
│   ├── context/            # Context detection
│   ├── generator/          # Override generation
│   ├── docker/             # Docker client wrapper
│   └── workspace/          # Workspace operations
├── pkg/                    # Public packages (if any)
├── main.go                 # Entry point
├── go.mod
└── go.sum
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

### Context Detection

```go
type Context struct {
    WorkspacePath string
    WorkspaceName string
    AppPath       string
    AppName       string
}

func DetectContext(startDir string) (*Context, error) {
    // Walk up directories looking for markers
}
```

---

## Testing Strategy

### Unit Tests

Test pure functions and logic in isolation.

```bash
go test ./internal/...
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
go build -o contrail .
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

## Related Documents

- [Architecture Overview](../architecture/overview.md)
- [CLI Reference](../reference/cli.md)

<!-- See appendices/tech-stack/ for complete code examples -->
```

### Appendix Content

Create `implementation/appendices/tech-stack/`:
- `scaffold-main.go` — Complete main.go scaffold
- `scaffold-cmd-root.go` — Root command scaffold
- `scaffold-config.go` — Configuration loading scaffold
- `scaffold-context.go` — Context detection scaffold
- `scaffold-generator.go` — Override generator scaffold
- `makefile` — Complete Makefile
- `goreleaser.yaml` — GoReleaser configuration

**Note**: These are large code blocks (>50 lines each) that belong in appendices.

---

## Completion Checklist

- [ ] `implementation/tech-stack.md` created
- [ ] Appendix directory created
- [ ] Scaffold files extracted to appendices
- [ ] Code blocks in main doc are <50 lines
- [ ] Cross-references added

