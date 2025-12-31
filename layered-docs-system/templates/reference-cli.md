# [Tool Name] CLI Reference

**Version**: [Software version this documents]
**Generated**: YYYY-MM-DD (or "Hand-maintained")

---

## Synopsis

```
toolname [global-options] <command> [command-options] [arguments]
```

---

## Global Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | `~/.config/tool/config.yaml` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--help` | `-h` | Show help | - |
| `--version` | | Show version | - |

---

## Commands

### `command-one`

[Brief description of what this command does.]

**Usage**:
```
toolname command-one [options] <required-arg>
```

**Arguments**:

| Argument | Description |
|----------|-------------|
| `<required-arg>` | [What this argument is] |

**Options**:

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--flag-one` | `-f` | [What it does] | `false` |
| `--option-with-value` | `-o` | [What it does] | `default` |

**Examples**:

```bash
# Basic usage
toolname command-one myarg

# With options
toolname command-one --flag-one --option-with-value=custom myarg
```

---

### `command-two`

[Brief description of what this command does.]

**Usage**:
```
toolname command-two [options]
```

**Options**:

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--option` | `-o` | [What it does] | `default` |

**Examples**:

```bash
toolname command-two
toolname command-two --option=value
```

---

### `command-three subcommand`

[Commands can have subcommands. Document them as needed.]

**Usage**:
```
toolname command-three subcommand [options]
```

---

## Environment Variables

| Variable | Description | Default | Equivalent Option |
|----------|-------------|---------|-------------------|
| `TOOL_CONFIG` | Path to config file | `~/.config/tool/config.yaml` | `--config` |
| `TOOL_VERBOSE` | Enable verbose mode | `false` | `--verbose` |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Invalid arguments |
| `3` | Configuration error |

---

## See Also

- [Configuration Reference](./configuration.md)
- [Getting Started Guide](../tutorials/getting-started.md)

---

## Notes

[Optional: Any additional notes about CLI behavior, platform differences, etc.]
