# [Component] Configuration Reference

**Version**: 0.1.0

---

## Overview

[Brief description of this configuration.]

---

## Location

**File**: `path/to/config.yaml`

---

## Schema

```yaml
section:
  key: value
  nested:
    subkey: subvalue
```

---

## Options

| Key | Type | Required | Default | Description |
|-----|------|----------|---------|-------------|
| `section.key` | string | Yes | | Description |
| `section.nested.subkey` | bool | No | `true` | Description |

---

## Examples

### Minimal Configuration

```yaml
section:
  key: value
```

### Full Configuration

```yaml
section:
  key: value
  nested:
    subkey: true
```

---

## Related Documents

- [Related Reference](link.md)
