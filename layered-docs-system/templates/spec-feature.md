# [Feature Name] Specification

**Version**: 0.1.0
**Date**: YYYY-MM-DD
**Status**: Draft | Review | Accepted | Implemented

---

> **File Location**: `docs/specs/{feature-name}.md`
>
> Create appendices directory at `docs/specs/appendices/{feature-name}/` if large content is needed.

---

## Overview

[One paragraph describing what this feature does and why it exists.]

**Related Documents**:
- [ADR-NNNN: Relevant Decision](../decisions/NNNN-{title}.md)
- [Architecture: Relevant Section](../architecture/overview.md#section)

**Appendices** (if applicable):
- [Detailed Examples](./appendices/{feature-name}/examples.md)
- [Error Catalog](./appendices/{feature-name}/errors.md)

---

## Behavior

### Normal Flow

[Describe the typical behavior step by step.]

1. [Step 1]
2. [Step 2]
3. [Step 3]

### State Machine

[If the feature has distinct states, document them.]

```
┌─────────┐     trigger      ┌─────────┐
│ State A │─────────────────►│ State B │
└─────────┘                  └────┬────┘
                                  │
                             trigger
                                  │
                                  ▼
                            ┌─────────┐
                            │ State C │
                            └─────────┘
```

| State | Description | Transitions |
|-------|-------------|-------------|
| State A | [What this state means] | → State B (on trigger) |
| State B | [What this state means] | → State C (on trigger) |
| State C | [What this state means] | Terminal |

---

## Data Schema

### [Configuration/Input Name]

```yaml
# Example configuration
field_one: value
field_two:
  nested_field: value
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `field_one` | string | Yes | - | [What it does] |
| `field_two.nested_field` | string | No | `default` | [What it does] |

### Validation Rules

- [Rule 1: e.g., `field_one` must be non-empty]
- [Rule 2: e.g., `field_two.nested_field` must match pattern X]

---

## Examples

[Brief examples here. For complete workflow examples, see appendix.]

### Example 1: [Scenario Name]

**Input**:
```yaml
# Configuration
field_one: example
```

**Behavior**: [What happens]

**Result**: [Expected outcome]

### Example 2: [Another Scenario]

**Input**:
```yaml
# Configuration
field_one: other_example
field_two:
  nested_field: custom
```

**Behavior**: [What happens]

**Result**: [Expected outcome]

> **Note**: If code blocks exceed 50 lines or you need complete file examples, move to `appendices/{feature-name}/examples.md`.

---

## Edge Cases

### [Edge Case 1]

**Scenario**: [Description of unusual situation]

**Behavior**: [What the system does]

**Rationale**: [Why this behavior was chosen. Link to ADR if applicable.]

### [Edge Case 2]

**Scenario**: [Description of unusual situation]

**Behavior**: [What the system does]

---

## Error Handling

[Brief error handling overview. For the complete error catalog, see appendix.]

| Error Condition | Error Code/Type | Message | Recovery |
|-----------------|-----------------|---------|----------|
| [Condition 1] | [E001] | [User-facing message] | [What user can do] |
| [Condition 2] | [E002] | [User-facing message] | [What user can do] |

> **Note**: If this table exceeds 20 rows, move the full catalog to `appendices/{feature-name}/errors.md`.

---

## Integration Points

### With [Other Feature/Component]

[How this feature interacts with other parts of the system.]

**Data Flow**:
```
[This Feature] ──[data type]──► [Other Component]
                                      │
                                      ▼
                               [Result/Effect]
```

---

## Implementation Notes

[Optional: Guidance for implementers that isn't part of the spec itself.]

- [Note 1]
- [Note 2]

---

## Open Questions

- [ ] [Question that needs resolution before implementation]
- [ ] [Question that needs resolution]

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0 | YYYY-MM-DD | Initial draft |
