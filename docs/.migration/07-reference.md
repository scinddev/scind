# Migration Step: Layer 5 — Reference

**Prerequisites**: Read `common-instructions.md`, complete Layer 4 steps
**Estimated Size**: 2 files + appendices, approximately 1,700 lines total

---

## Overview

Extract reference documentation for CLI commands and configuration options. Reference docs are lookup-oriented, not tutorial-oriented.

**Source documents**:
- `specs/contrail-cli-reference.md` (CLI reference)
- `specs/contrail-technical-spec.md` (configuration reference)

---

## File 1: `reference/cli.md`

**Source**: `specs/contrail-cli-reference.md:1-800`

### Summary

Complete CLI reference with all commands, options, and flags.

### Content Structure

```markdown
# CLI Reference

Quick reference for all Contrail commands.

---

## Global Options

| Option | Short | Description |
|--------|-------|-------------|
| `--workspace` | `-w` | Target workspace |
| `--app` | `-a` | Target application |
| `--verbose` | `-v` | Verbose output |
| `--help` | `-h` | Show help |

---

## Commands

### contrail up

Start workspace or application.

**Usage**: `contrail up [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--build` | Build images before starting | false |
| `--detach` | Run in background | true |
| `--flavor` | Use specific flavor | default |

**Examples**:
```bash
contrail up                      # Start based on context
contrail up -w main              # Start entire workspace
contrail up -w main -a frontend  # Start specific app
```

### contrail down

Stop workspace or application.

**Usage**: `contrail down [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--volumes` | Remove volumes | false |
| `--timeout` | Shutdown timeout (seconds) | 10 |

**Examples**:
```bash
contrail down                    # Stop based on context
contrail down -w main --volumes  # Stop and remove volumes
```

### contrail status

Show workspace/application status.

**Usage**: `contrail status [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--format` | Output format (table, json, yaml) | table |

### contrail logs

View container logs.

**Usage**: `contrail logs [service] [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--follow` | Follow log output | false |
| `--tail` | Number of lines | all |
| `--timestamps` | Show timestamps | false |

### contrail exec

Execute command in container.

**Usage**: `contrail exec <service> <command>`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--tty` | Allocate TTY | true |
| `--interactive` | Keep STDIN open | true |

### contrail generate

Generate override files without starting.

**Usage**: `contrail generate [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--dry-run` | Show what would be generated | false |

### contrail validate

Validate configuration files.

**Usage**: `contrail validate [options]`

**Options**:
| Option | Description | Default |
|--------|-------------|---------|
| `--strict` | Treat warnings as errors | false |

### contrail init

Initialize new workspace or application.

**Usage**: `contrail init <type> [name]`

**Types**:
- `workspace` — Create workspace.yaml
- `app` — Create application.yaml

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Docker error |

---

## Related Documents

- [Configuration Reference](./configuration.md)
- [Context Detection Spec](../specs/context-detection.md)

<!-- See appendices/cli/ for detailed examples and error messages -->
```

### Appendix Content

Create `reference/appendices/cli/`:
- `detailed-examples.md` — Extended examples for each command
- `error-messages.md` — Complete error message catalog

---

## File 2: `reference/configuration.md`

**Source**: `specs/contrail-technical-spec.md:200-400` + `specs/contrail-cli-reference.md:900-1200`

### Summary

Complete configuration reference for all YAML files.

### Content Structure

```markdown
# Configuration Reference

Quick reference for all configuration files and options.

---

## File Locations

| File | Location | Purpose |
|------|----------|---------|
| proxy.yaml | `~/.config/contrail/proxy.yaml` | Global proxy settings |
| workspace.yaml | `{workspace}/workspace.yaml` | Workspace definition |
| application.yaml | `{app}/application.yaml` | Application settings |

---

## proxy.yaml

Global Traefik and TLS configuration.

```yaml
domain: test                    # Base domain for hostnames
tls:
  mode: auto                    # auto, custom, disabled
  cert_file: ""                 # Path for custom mode
  key_file: ""                  # Path for custom mode
traefik:
  dashboard: true               # Enable Traefik dashboard
  dashboard_port: 8080          # Dashboard port
  log_level: ERROR              # Traefik log level
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `domain` | string | No | `test` | Base domain for public hostnames |
| `tls.mode` | enum | No | `auto` | TLS mode: auto, custom, disabled |
| `tls.cert_file` | string | If custom | - | Path to certificate file |
| `tls.key_file` | string | If custom | - | Path to private key file |
| `traefik.dashboard` | bool | No | `true` | Enable Traefik dashboard |
| `traefik.dashboard_port` | int | No | `8080` | Dashboard port |
| `traefik.log_level` | string | No | `ERROR` | Log level |

---

## workspace.yaml

Workspace definition and application references.

```yaml
name: main                      # Workspace name
applications:
  - path: ../frontend           # Relative path to app
    name: frontend              # Optional name override
  - path: ../backend
  - path: ../shared-db
    name: db
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | - | Workspace name (used in hostnames) |
| `applications` | array | Yes | - | List of applications |
| `applications[].path` | string | Yes | - | Path to application directory |
| `applications[].name` | string | No | dir name | Override application name |

---

## application.yaml

Application-specific configuration.

```yaml
name: frontend                  # Application name
flavors:
  default:                      # Default flavor
    profiles: []                # Docker Compose profiles
  full:
    profiles: [backend, db]
  minimal:
    profiles: []
exported_services:
  web:                          # Service name for discovery
    service: nginx              # Docker Compose service name
    ports:
      - type: proxied
        protocol: https
        visibility: public
        container_port: 80
      - type: proxied
        protocol: http
        visibility: protected
        container_port: 80
  api:
    service: node
    ports:
      - type: proxied
        protocol: http
        visibility: protected
        container_port: 3000
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | No | dir name | Application name |
| `flavors` | object | No | `{default: {}}` | Named configurations |
| `flavors.<name>.profiles` | array | No | `[]` | Docker Compose profiles |
| `exported_services` | object | No | `{}` | Services to expose |
| `exported_services.<name>.service` | string | No | key name | Compose service name |
| `exported_services.<name>.ports` | array | Yes | - | Port configurations |
| `exported_services.<name>.ports[].type` | enum | Yes | - | `proxied` or `assigned` |
| `exported_services.<name>.ports[].protocol` | string | If proxied | - | `http` or `https` |
| `exported_services.<name>.ports[].visibility` | enum | No | `protected` | `public` or `protected` |
| `exported_services.<name>.ports[].container_port` | int | Yes | - | Container port |
| `exported_services.<name>.ports[].port` | int | If assigned | - | Preferred host port |

---

## Related Documents

- [CLI Reference](./cli.md)
- [Configuration Schemas Spec](../specs/configuration-schemas.md)
- [ADR-0006: Three Configuration Schemas](../decisions/0006-three-configuration-schemas.md)

<!-- See appendices/configuration/ for complete examples and JSON schemas -->
```

### Appendix Content

Create `reference/appendices/configuration/`:
- `complete-examples.md` — Full working configuration examples
- `json-schemas/` — JSON Schema files for validation

---

## Completion Checklist

- [ ] `reference/cli.md` created
- [ ] `reference/configuration.md` created
- [ ] Appendix directories created
- [ ] Large content moved to appendices
- [ ] Cross-references added

