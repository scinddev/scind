# [System Name]: Architecture Overview

**Version**: 0.1.0
**Date**: YYYY-MM-DD
**Status**: Draft | Active | Deprecated

---

## System Context

### Overview

[One paragraph describing what the system does and its primary users.]

### Context Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         ENVIRONMENT                              │
│                                                                  │
│   ┌──────────┐                              ┌──────────────┐    │
│   │  User 1  │──────────┐      ┌───────────│External Sys 1│    │
│   └──────────┘          │      │            └──────────────┘    │
│                         ▼      ▼                                 │
│   ┌──────────┐     ┌──────────────┐         ┌──────────────┐    │
│   │  User 2  │────►│   [SYSTEM]   │◄───────│External Sys 2│    │
│   └──────────┘     └──────────────┘         └──────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### External Dependencies

| System | Purpose | Protocol |
|--------|---------|----------|
| [External 1] | [What it provides] | [HTTP/gRPC/etc.] |
| [External 2] | [What it provides] | [Protocol] |

---

## Containers

### Container Diagram

```
┌───────────────────────────────────────────────────────────────┐
│                        [SYSTEM NAME]                           │
│                                                                │
│  ┌────────────────┐      ┌────────────────┐                   │
│  │  [Container 1] │◄────►│  [Container 2] │                   │
│  │  (technology)  │      │  (technology)  │                   │
│  └───────┬────────┘      └────────┬───────┘                   │
│          │                        │                            │
│          │       ┌────────────────┘                            │
│          │       │                                             │
│          ▼       ▼                                             │
│  ┌────────────────────┐                                        │
│  │   [Container 3]    │                                        │
│  │   (technology)     │                                        │
│  └────────────────────┘                                        │
│                                                                │
└───────────────────────────────────────────────────────────────┘
```

### Container Descriptions

| Container | Technology | Purpose |
|-----------|------------|---------|
| [Container 1] | [Language/framework] | [What it does] |
| [Container 2] | [Language/framework] | [What it does] |
| [Container 3] | [Database/cache/etc.] | [What it stores/provides] |

---

## Key Components

[Only include this section for containers complex enough to warrant it.]

### [Container 1] Components

| Component | Purpose |
|-----------|---------|
| [Component A] | [What it does] |
| [Component B] | [What it does] |

---

## Communication Patterns

### Internal Communication

| From | To | Method | Purpose |
|------|----|--------|---------|
| [Container 1] | [Container 2] | [HTTP/gRPC/queue] | [What data/commands] |

### External Communication

| Direction | Endpoint | Purpose |
|-----------|----------|---------|
| Inbound | [Entry point] | [What external systems/users connect here] |
| Outbound | [External API] | [What we send/request] |

---

## Cross-Cutting Concerns

### Security

[How authentication, authorization, and secrets are handled. Reference ADRs for decisions.]

### Logging & Observability

[Logging strategy, metrics, tracing approach.]

### Error Handling

[Error propagation patterns, retry strategies.]

---

## Quality Attributes

| Attribute | Requirement | How Achieved |
|-----------|-------------|--------------|
| [Performance] | [e.g., < 100ms response] | [Caching, indexing, etc.] |
| [Reliability] | [e.g., 99.9% uptime] | [Redundancy, failover] |
| [Scalability] | [e.g., 10k concurrent users] | [Horizontal scaling] |

---

## Known Risks & Technical Debt

| Risk/Debt | Impact | Mitigation |
|-----------|--------|------------|
| [Risk 1] | [What could go wrong] | [Plan to address] |
| [Debt 1] | [Current limitation] | [Plan to address] |

---

## Related Documents

- [Vision](../product/vision.md)
- [Decisions](../decisions/)
- [Specifications](../specs/)
