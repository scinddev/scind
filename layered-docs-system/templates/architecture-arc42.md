# [System Name]: Architecture Documentation

**Version**: 0.1.0
**Date**: YYYY-MM-DD
**Status**: Draft | Active | Deprecated

Based on the [arc42](https://arc42.org/) template.

---

## 1. Introduction and Goals

### Requirements Overview

[Brief description of the key functional requirements.]

### Quality Goals

| Priority | Quality Goal | Scenario |
|----------|--------------|----------|
| 1 | [e.g., Performance] | [Concrete scenario] |
| 2 | [e.g., Maintainability] | [Concrete scenario] |
| 3 | [e.g., Security] | [Concrete scenario] |

### Stakeholders

| Role | Expectations |
|------|--------------|
| [Developer] | [What they need from the architecture] |
| [Operator] | [What they need from the architecture] |
| [End User] | [What they need from the architecture] |

---

## 2. Constraints

### Technical Constraints

| Constraint | Background |
|------------|------------|
| [e.g., Must run on Docker] | [Why this is required] |
| [e.g., Go 1.21+] | [Why this version] |

### Organizational Constraints

| Constraint | Background |
|------------|------------|
| [e.g., Small team] | [Impact on decisions] |
| [e.g., Open source] | [Licensing considerations] |

### Conventions

| Convention | Description |
|------------|-------------|
| [Naming] | [Standards used] |
| [Coding style] | [Linters, formatters] |

---

## 3. Context and Scope

### Business Context

```
[Diagram showing the system and its business partners/users]
```

| Partner/User | Inputs | Outputs |
|--------------|--------|---------|
| [User 1] | [What they provide] | [What they receive] |
| [External System] | [What it provides] | [What it receives] |

### Technical Context

```
[Diagram showing technical interfaces and protocols]
```

| Interface | Protocol | Purpose |
|-----------|----------|---------|
| [Interface 1] | [HTTP/gRPC] | [What it's used for] |

---

## 4. Solution Strategy

[Summary of fundamental design decisions and solution strategies. Reference ADRs for detailed rationale.]

| Approach | Rationale | Reference |
|----------|-----------|-----------|
| [Strategy 1] | [Why chosen] | [ADR-NNNN](../decisions/NNNN.md) |
| [Strategy 2] | [Why chosen] | [ADR-NNNN](../decisions/NNNN.md) |

---

## 5. Building Block View

### Level 1: System Overview

```
[High-level diagram of major building blocks]
```

| Building Block | Purpose |
|----------------|---------|
| [Block 1] | [What it does] |
| [Block 2] | [What it does] |

### Level 2: [Container/Subsystem Name]

[Only include for complex subsystems that need decomposition.]

```
[Diagram of internal components]
```

---

## 6. Runtime View

### [Scenario 1 Name]

[Sequence or activity diagram showing runtime behavior.]

```
[Participant 1] -> [Participant 2]: [Message]
[Participant 2] -> [Participant 3]: [Message]
[Participant 3] --> [Participant 2]: [Response]
```

**Description**: [What this scenario demonstrates.]

### [Scenario 2 Name]

[Another key runtime scenario.]

---

## 7. Deployment View

### Infrastructure

```
[Diagram showing deployment topology]
```

| Node | Description | Software |
|------|-------------|----------|
| [Host/Container] | [What it is] | [What runs there] |

### Deployment Artifacts

| Artifact | Description |
|----------|-------------|
| [Binary/Image] | [What it is, how it's built] |

---

## 8. Cross-cutting Concepts

### [Concept 1: e.g., Error Handling]

[Pattern or approach used across the system.]

### [Concept 2: e.g., Logging]

[Pattern or approach used across the system.]

### [Concept 3: e.g., Configuration]

[Pattern or approach used across the system.]

---

## 9. Architecture Decisions

See the [decisions directory](../decisions/) for all ADRs.

Key decisions:

- [ADR-0001: Title](../decisions/0001-title.md) - [Brief summary]
- [ADR-0002: Title](../decisions/0002-title.md) - [Brief summary]

---

## 10. Quality Requirements

### Quality Tree

```
Quality
├── Performance
│   ├── Response Time
│   └── Throughput
├── Reliability
│   ├── Availability
│   └── Fault Tolerance
└── Security
    ├── Authentication
    └── Authorization
```

### Quality Scenarios

| Quality | Scenario | Measure |
|---------|----------|---------|
| [Performance] | [Under X load...] | [< Y ms] |
| [Reliability] | [When component fails...] | [System continues] |

---

## 11. Risks and Technical Debt

| Type | Description | Priority | Mitigation |
|------|-------------|----------|------------|
| Risk | [What could go wrong] | High/Med/Low | [Plan] |
| Debt | [Current shortcut] | High/Med/Low | [Plan] |

---

## 12. Glossary

| Term | Definition |
|------|------------|
| [Term 1] | [Definition] |
| [Term 2] | [Definition] |

---

## Related Documents

- [Vision](../product/vision.md)
- [Specifications](../specs/)
- [Reference](../reference/)
