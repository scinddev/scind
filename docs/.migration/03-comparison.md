# Migration Step: Layer 2 - Comparison

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 file, approximately 80 lines

---

## Overview

Create `product/comparison.md` consolidating comparison content from `specs/contrail-prd.md`.

**Source Sections**:
- "Why Existing Solutions Fall Short" table (lines 34-42)
- "Comparison with Related Tools" appendix (lines 632-642)

---

## Output File: `product/comparison.md`

**Source**: `specs/contrail-prd.md:34-42, 632-642`

### Content

```markdown
# Contrail Comparison

How Contrail compares to existing solutions for Docker-based development environments.

---

## Why Existing Solutions Fall Short

| Solution | Limitation |
|----------|------------|
| Docker Compose alone | No built-in multi-instance orchestration; manual project naming |
| Docker Compose `include` | Merges into single application model; doesn't handle parallel instances |
| DDEV / Lando / Docksal | Single-application focused (one Drupal site, not multi-app stacks) |
| Skaffold / Tilt / Garden | Kubernetes-focused, not Docker Compose |
| Manual scripts | Error-prone, hard to maintain, no conventions |

---

## Feature Comparison Matrix

| Feature | Contrail | Docker `include` | DDEV/Lando | Tilt/Garden |
|---------|----------|------------------|------------|-------------|
| Multi-app orchestration | Yes | Yes (merged) | No | Yes |
| Parallel workspace instances | Yes | No | No | No |
| Apps remain agnostic | Yes | N/A | N/A | No |
| Docker Compose native | Yes | Yes | Yes | No (K8s) |
| Generated integration | Yes | No | Yes | Yes |
| Service discovery | Yes | Manual | Limited | Yes |

---

## When to Use Contrail

**Use Contrail when you need:**
- Multiple isolated copies of a multi-application stack running simultaneously
- Applications that remain portable and framework-agnostic
- Generated configuration rather than manual Docker Compose wiring
- Docker Compose (not Kubernetes) as your runtime

**Consider alternatives when:**
- You only have a single application (DDEV, Lando may be simpler)
- You're targeting Kubernetes in production (consider Tilt, Garden, Skaffold)
- You need CI/CD pipeline integration (Contrail focuses on local development)

<!-- Migrated from specs/contrail-prd.md:34-42, 632-642 -->
```

---

## Completion Checklist

- [ ] `product/comparison.md` created
- [ ] All comparison content preserved
- [ ] Source attribution present
