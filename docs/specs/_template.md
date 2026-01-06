# [Feature Name] Specification

**Version**: 0.1.0
**Date**: YYYY-MM-DD
**Status**: Draft | Review | Accepted | Implemented

---

## Overview

[One paragraph describing what this feature does and why it exists.]

**Related Documents**:
- [ADR-NNNN: Relevant Decision](../decisions/NNNN-{title}.md)
- [Architecture: Relevant Section](../architecture/README.md#section)

**Appendices** (if applicable):
- [Detailed Examples](./appendices/{feature-name}/examples.md)

---

## Behavior

### Normal Flow

[Describe the typical behavior step by step.]

1. [Step 1]
2. [Step 2]
3. [Step 3]

---

## Data Schema

### [Configuration/Input Name]

```yaml
# Example configuration
field_one: value
field_two:
  nested_field: value
```

### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `field_one` | string | Yes | - | [What it does] |
| `field_two.nested_field` | string | No | `default` | [What it does] |

### Validation Rules

- [Rule 1: e.g., `field_one` must be non-empty]
- [Rule 2: e.g., `field_two.nested_field` must match pattern X]

---

## Examples

### Example 1: [Scenario Name]

**Input**:
```yaml
field_one: example
```

**Behavior**: [What happens]

**Result**: [Expected outcome]

---

## Edge Cases

### [Edge Case 1]

**Scenario**: [Description of unusual situation]

**Behavior**: [What the system does]

**Rationale**: [Why this behavior was chosen. Link to ADR if applicable.]

---

## Error Handling

| Error Condition | Message | Recovery |
|-----------------|---------|----------|
| [Condition 1] | [User-facing message] | [What user can do] |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0 | YYYY-MM-DD | Initial draft |
