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

## Detailed Comparisons

### Tilt

**Similarities**:
- Focused on local development experience
- Watches files and rebuilds
- Provides a dashboard

**Differences**:
- Tilt targets Kubernetes; Contrail targets Docker Compose
- Tilt has its own configuration language (Tiltfile); Contrail uses YAML
- Tilt handles builds; Contrail delegates to Compose

**When to use Tilt**: Your production is Kubernetes and you want local dev to mirror that.

**When to use Contrail**: Your local dev uses Docker Compose and you want to enhance it without adopting Kubernetes.

### DevSpace

**Similarities**:
- Simplifies development workflows
- Provides CLI for common operations

**Differences**:
- DevSpace abstracts Kubernetes; Contrail enhances Compose
- DevSpace can sync files into running containers; Contrail relies on volume mounts
- DevSpace manages deployments; Contrail generates overrides

**When to use DevSpace**: You're deploying to Kubernetes and want development shortcuts.

**When to use Contrail**: You're staying with Docker Compose and want workspace isolation.

### Garden

**Similarities**:
- Declarative environment definition
- Supports multiple "stacks" or environments

**Differences**:
- Garden is a full platform with cloud features; Contrail is a local tool
- Garden has its own execution model; Contrail delegates to Compose
- Garden requires more configuration; Contrail emphasizes conventions

**When to use Garden**: You need a comprehensive platform spanning local to cloud.

**When to use Contrail**: You want a focused local tool that enhances Compose.

### Skaffold

**Similarities**:
- CLI-driven workflow
- Handles multiple applications

**Differences**:
- Skaffold focuses on build/deploy pipeline; Contrail focuses on runtime orchestration
- Skaffold targets Kubernetes; Contrail targets Compose
- Skaffold integrates with CI/CD; Contrail is development-only

**When to use Skaffold**: You need a build/deploy pipeline for Kubernetes.

**When to use Contrail**: You need workspace isolation for local Compose development.

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

---

## Summary

Contrail occupies a specific niche: **enhancing Docker Compose for teams who want workspace isolation, automatic routing, and service discovery without adopting Kubernetes**.

If your workflow is Docker Compose locally and you want to improve the developer experience without changing your deployment target, Contrail is designed for you.

If you're moving to Kubernetes or need cloud-native features, consider Tilt, DevSpace, Garden, or Skaffold based on your specific needs.

---

## Related Documentation

- [Product Vision](./vision.md) - What Contrail aims to achieve
- [Roadmap](./roadmap.md) - Planned enhancements

<!-- Migrated from specs/contrail-prd.md:34-42, 317-393, 632-642 -->
