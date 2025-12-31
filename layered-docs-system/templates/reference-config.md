# [Tool Name] Configuration Reference

**Version**: [Software version this documents]
**Generated**: YYYY-MM-DD (or "Hand-maintained")

---

## Configuration Files

| File | Purpose | Location |
|------|---------|----------|
| `config.yaml` | Main configuration | `~/.config/tool/config.yaml` |
| `local.yaml` | Local overrides | `./.tool/local.yaml` |

**Precedence** (highest to lowest):
1. Command-line flags
2. Environment variables
3. Local configuration file
4. User configuration file
5. Default values

---

## Configuration Schema

### Top-Level Structure

```yaml
# Full configuration example
section_one:
  option_a: value
  option_b: value

section_two:
  option_c: value
  nested:
    option_d: value
```

---

## `section_one`

[Description of what this section configures.]

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `option_a` | `string` | Yes | - | [What it does] |
| `option_b` | `integer` | No | `100` | [What it does] |

**Example**:
```yaml
section_one:
  option_a: "my-value"
  option_b: 200
```

---

## `section_two`

[Description of what this section configures.]

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `option_c` | `boolean` | No | `true` | [What it does] |
| `nested.option_d` | `string` | No | `"default"` | [What it does] |

**Example**:
```yaml
section_two:
  option_c: false
  nested:
    option_d: "custom"
```

---

## `section_two.nested`

[If nesting is deep, break out subsections.]

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `option_d` | `string` | No | `"default"` | [What it does] |
| `option_e` | `list[string]` | No | `[]` | [What it does] |

---

## Type Reference

### Custom Types

#### `type_name`

[Description of what values are valid.]

**Valid values**:
- `value_one` - [Meaning]
- `value_two` - [Meaning]
- `value_three` - [Meaning]

**Example**:
```yaml
field: value_one
```

---

## Environment Variable Mapping

Configuration values can be set via environment variables:

| Config Path | Environment Variable |
|-------------|---------------------|
| `section_one.option_a` | `TOOL_SECTION_ONE_OPTION_A` |
| `section_one.option_b` | `TOOL_SECTION_ONE_OPTION_B` |
| `section_two.option_c` | `TOOL_SECTION_TWO_OPTION_C` |

---

## Validation Rules

- [Rule 1: e.g., `option_a` must not contain spaces]
- [Rule 2: e.g., `option_b` must be between 1 and 1000]
- [Rule 3: e.g., If `option_c` is true, `option_d` is required]

---

## Complete Example

```yaml
# Full working configuration
section_one:
  option_a: "production"
  option_b: 500

section_two:
  option_c: true
  nested:
    option_d: "custom-value"
    option_e:
      - "item1"
      - "item2"
```

---

## See Also

- [CLI Reference](./cli.md)
- [Architecture Overview](../architecture/overview.md)
