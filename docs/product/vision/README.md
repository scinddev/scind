# Contrail Product Vision

**Version**: 0.5.0
**Date**: December 2024
**Status**: Design Phase

---

## Executive Summary

Contrail is a workspace orchestration system for Docker Compose that enables developers to run multiple isolated instances of multi-application stacks simultaneously on a single host. It solves the problem of needing complete, independent environments for development, code review, and testing without requiring Kubernetes or cloud infrastructure.

The name "Contrail" evokes the trails left by aircraft—parallel paths that don't intersect, much like the isolated workspaces the system creates.

---

## Problem Statement

### The Scenario

A developer works on a system composed of multiple applications (e.g., `app-one`, `app-two`, `app-three`) that communicate with each other. They need to run multiple complete copies of this stack simultaneously:

- **dev**: Active development with local changes
- **review**: Clean checkout of a PR branch for code review
- **control**: Stable baseline for comparison

Each environment needs:
1. Internal communication between applications (app-one can reach app-two's API)
2. External access via unique hostnames (dev-app-one-web.contrail.test, review-app-one-web.contrail.test)
3. Complete isolation (dev's database is separate from review's database)

### Why Existing Solutions Fall Short

| Solution | Limitation |
|----------|------------|
| Docker Compose alone | No built-in multi-instance orchestration; manual project naming |
| Docker Compose `include` | Merges into single application model; doesn't handle parallel instances |
| DDEV / Lando / Docksal | Single-application focused (one Drupal site, not multi-app stacks) |
| Skaffold / Tilt / Garden | Kubernetes-focused, not Docker Compose |
| Manual scripts | Error-prone, hard to maintain, no conventions |

### The Gap

There is no existing tool that orchestrates **multiple isolated instances of multi-application Docker Compose stacks** with **generated integration configuration** while keeping **applications workspace-agnostic**.

---

## Vision

Contrail provides a thin coordination layer over Docker Compose that:

1. **Preserves application independence**: Applications don't know they're in a workspace
2. **Uses pure overlay**: All integration happens via generated Docker Compose override files
3. **Follows conventions**: Predictable naming for hostnames, aliases, and networks
4. **Separates structure from state**: Configuration describes what exists; runtime state describes what's active
5. **Enables direct Docker Compose access**: The `contrail-compose` shell function provides context-aware passthrough to Docker Compose with full tab completion

---

## Core Concepts

### Workspace

A logical grouping of applications that run together, sharing an internal network for communication while isolated from other workspaces.

```
workspace: dev
├── app-one (project: dev-app-one)
├── app-two (project: dev-app-two)
└── app-three (project: dev-app-three)
    └── all connected via: dev-internal network
```

### Application

A Docker Compose-based service that can participate in workspaces. Applications define a **service contract** (`application.yaml`) that declares what they export to the workspace, but their core `docker-compose.yaml` remains workspace-agnostic.

### Single-Application Workspace

Contrail supports workspaces where the workspace directory is also the application directory. This allows promoting existing Docker Compose projects without restructuring.

### Flavor

A named configuration that specifies which Docker Compose files to use when running an application. Enables "lite" vs "full" modes without duplicating configuration.

---

## Success Criteria

1. **Zero application changes**: Existing Docker Compose applications work without modification
2. **Parallel environments**: Can run dev, review, and control simultaneously without conflicts
3. **Predictable naming**: Given workspace and app names, hostnames and aliases are deterministic
4. **Fast iteration**: Switching flavors or regenerating config takes seconds
5. **Debuggable**: Generated files are human-readable; easy to understand what's happening

---

## Non-Goals

1. **Kubernetes support**: Contrail is specifically for Docker Compose environments
2. **Production deployment**: Focused on local development and testing
3. **CI/CD integration**: May come later, but not initial focus
4. **GUI**: CLI-first; GUI could be added later
5. **Image building**: Uses existing images; doesn't manage builds
6. **Secret management**: Uses Docker Compose's existing mechanisms
7. **Windows native support**: Initial release targets macOS and Linux. Windows users should use WSL2.

---

## Known Limitations (v1)

1. **Concurrent operations**: Running multiple Contrail commands simultaneously may cause race conditions. Use one terminal per workspace for operations.
2. **Port garbage collection**: Ports from workspaces deleted via filesystem (not `workspace destroy`) require manual cleanup with `contrail port gc`.

---

## Future Considerations

- **Port Type Plugins**: Protocol handlers for SNI-based TCP routing (postgresql, mysql)
- **Application Dependencies**: Startup ordering between applications
- **Shared Volumes**: Workspace-level volumes mountable into multiple applications
- **Health Checks**: Integration with Docker health checks for readiness waiting

---

## Related Documentation

- [Architecture Overview](../../architecture/overview/README.md)
- [CLI Reference](../../reference/cli/README.md)
- [Configuration Reference](../../reference/configuration/README.md)
