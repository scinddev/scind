# Comparison with Alternatives

This document compares Contrail to other tools in the container orchestration and development environment space.

---

## Comparison Matrix

| Aspect | Contrail | Tilt | DevSpace | Garden | Skaffold |
|--------|----------|------|----------|--------|----------|
| **Target** | Docker Compose | Kubernetes | Kubernetes | Kubernetes | Kubernetes |
| **Philosophy** | Enhance existing workflow | Replace build/deploy | Abstract k8s | Declarative stack | Build/deploy pipeline |
| **Learning curve** | Low (if you know Compose) | Medium | Medium | High | Medium |
| **Isolation** | Per-workspace | Per-namespace | Per-namespace | Per-environment | Per-namespace |
| **Local focus** | Yes | Yes | Yes | Mixed | Mixed |

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

## Summary

Contrail occupies a specific niche: **enhancing Docker Compose for teams who want workspace isolation, automatic routing, and service discovery without adopting Kubernetes**.

If your workflow is Docker Compose locally and you want to improve the developer experience without changing your deployment target, Contrail is designed for you.

If you're moving to Kubernetes or need cloud-native features, consider Tilt, DevSpace, Garden, or Skaffold based on your specific needs.

---

## Related Documents

- [Vision](./vision.md) — Product vision and concepts
- [Architecture Overview](../architecture/overview.md) — How Contrail works

<!-- Migrated from specs/contrail-prd.md:317-393 -->
