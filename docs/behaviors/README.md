# Behaviors

Executable behavior specifications for Scind.

## Contents

This directory contains Gherkin feature files that verify expected system behaviors. These serve as **living documentation** — if the tests pass, the documentation is accurate.

*No behavior specifications have been added yet.*

## Directory Structure

Feature files are organized by domain:

```
behaviors/
├── README.md
├── {domain}/                    # e.g., workspace/, proxy/
│   └── {feature}.feature
└── support/                     # optional: step definitions
    └── step_definitions/
```

Example paths:
- `behaviors/workspace/workspace-lifecycle.feature`
- `behaviors/proxy/proxy-routing.feature`

## When to Create a Behavior File

Use executable specs for:
- **Behaviors that have historically broken** — prevent regressions
- **Complex multi-step workflows** — document the expected sequence
- **Integration points between components** — verify contracts
- **Critical user journeys** — ensure key paths always work

## Template

```gherkin
# This feature verifies behaviors from:
# See: ../specs/{feature}.md

Feature: [Feature Name]
  As a [role]
  I want [capability]
  So that [benefit]

  Background:
    Given [common precondition]

  Scenario: [Scenario Name]
    Given [initial context]
    When [action is taken]
    Then [expected outcome]

  Scenario: [Edge Case Name]
    Given [edge case context]
    When [action is taken]
    Then [expected outcome]
```

## Linking to Specifications

Every behavior file should reference the specification it verifies. Add a comment at the top:

```gherkin
# This feature verifies behaviors from:
# See: ../specs/workspace-lifecycle.md
```

This creates traceability between the executable test and the specification it validates.

## Running Tests

```bash
# Run all behavior tests
cucumber-js docs/behaviors/

# Run a specific domain
cucumber-js docs/behaviors/workspace/
```

## Related Documents

- [DOCUMENTATION-GUIDE.md](../DOCUMENTATION-GUIDE.md) — Full LDS reference with classification heuristics
- [Layer 6: Behaviors](../DOCUMENTATION-GUIDE.md#layer-6-behaviors-detailed-guidance) — Detailed guidance
