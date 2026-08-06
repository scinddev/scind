<!-- Migrated from specs/scind-prd.md:1-108 -->
<!-- Additional content from specs/scind-prd.md:562-571 (Non-Goals) -->
<!-- Additional content from specs/scind-prd.md:574-578 (Known Limitations) -->
<!-- Additional content from specs/scind-prd.md:580-588 (Success Criteria) -->
<!-- Extraction ID: vision-main, vision-non-goals, vision-known-limitations, vision-success-criteria -->

# Scind: Product Requirements Document

**Version**: 0.6.0
**Date**: January 2025
**Status**: Design Phase

---

## Executive Summary

Scind is a workspace orchestration system for Docker Compose that enables developers to run multiple isolated instances of multi-application stacks simultaneously on a single host. It solves the problem of needing complete, independent environments for development, code review, and testing without requiring Kubernetes or cloud infrastructure.

The name "Scind" (from Latin *scindere*, to cut/split) reflects the system's core purpose: cleanly separating isolated workspaces that don't intersect.

---

## Problem Statement

### The Scenario

A developer works on a system composed of multiple applications (e.g., `frontend`, `backend`, `shared-db`) that communicate with each other. They need to run multiple complete copies of this stack simultaneously:

- **dev**: Active development with local changes
- **review**: Clean checkout of a PR branch for code review
- **control**: Stable baseline for comparison

Each environment needs:
1. Internal communication between applications (frontend can reach backend's API)
2. External access via unique hostnames (dev-frontend-web.scind.test, review-frontend-web.scind.test)
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

## Product Vision

Scind provides a thin coordination layer over Docker Compose that:

1. **Preserves application independence**: Applications don't know they're in a workspace
2. **Uses pure overlay**: All integration happens via generated Docker Compose override files
3. **Follows conventions**: Predictable naming for hostnames, aliases, and networks
4. **Separates structure from state**: Configuration describes what exists; runtime state describes what's active
5. **Enables direct Docker Compose access**: The `scind-compose` shell function provides context-aware passthrough to Docker Compose with full tab completion

---

## Target Audience

Scind is designed for:

- **Developers working on multi-application projects**: Teams building systems composed of multiple cooperating services (frontend, backend, databases)
- **Teams needing isolated local development environments**: Reviewers and developers who need to run multiple complete copies of a system simultaneously
- **Projects using Docker Compose for local development**: Teams already invested in Docker Compose who want workspace orchestration without migrating to Kubernetes

---

## Core Concepts

### Workspace

A logical grouping of applications that run together, sharing an internal network for communication while isolated from other workspaces.

```
workspace: dev
├── frontend (project: dev-frontend)
├── backend (project: dev-backend)
└── shared-db (project: dev-shared-db)
    └── all connected via: dev-internal network
```

*For the isolation mechanism, see [ADR-0001: Docker Compose Project Name Isolation](../decisions/0001-docker-compose-project-name-isolation.md).*

### Application

A Docker Compose-based service that can participate in workspaces. Applications define a **service contract** (`application.yaml`) that declares what they export to the workspace, but their core `docker-compose.yaml` remains workspace-agnostic.

### Single-Application Workspace

Scind supports workspaces where the workspace directory is also the application directory. This allows promoting existing Docker Compose projects without restructuring:

```bash
cd ~/my-existing-project
scind workspace init --workspace=dev
scind app init --app=myapp
```

### Port Types and Proxying

| Type | Protocol | Behavior | External Access | Internal Access | Use Case |
|------|----------|----------|-----------------|-----------------|----------|
| **proxied** | `https` | Traefik HTTPS proxy | Via hostname (port 443) | Via alias | Web frontends, public APIs |
| **proxied** | `http` | Traefik HTTP proxy | Via hostname (port 80) | Via alias | HTTP services, redirects |
| **assigned** | - | Direct port binding | Via assigned host port | Via alias | Databases, caches, debug ports |
| **proxied** | `tcp`, etc. | SNI TCP proxy (future) | Via hostname | Via alias | Database GUIs, external tools |

Services not listed in `exported_services` remain **private**—only accessible within the application's own Docker Compose network.

*For the technical decision, see [ADR-0007: Port Type System](../decisions/0007-port-type-system.md).*

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent to collaborators—it does not change Scind's core behavior. Both public and protected proxied services route through Traefik.

Visibility is exposed via Docker labels (`workspace.visibility`), enabling external tools like Servlo to distinguish between public and protected services for display or filtering purposes.

### Flavor

A named configuration that specifies which Docker Compose files to use when running an application. Enables "lite" vs "full" modes without duplicating configuration.

---

## Non-Goals

1. **Kubernetes support**: Scind is specifically for Docker Compose environments
2. **Production deployment**: Focused on local development and testing
3. **CI/CD integration**: May come later, but not initial focus
4. **GUI**: CLI-first; GUI could be added later
5. **Image building**: Uses existing images; doesn't manage builds
6. **Secret management**: Uses Docker Compose's existing mechanisms
7. **Windows native support**: Initial release targets macOS and Linux. Windows users should use WSL2.

---

## Known Limitations (v1)

1. **Concurrent operations**: Running multiple Scind commands simultaneously (e.g., two terminals running `workspace up`) may cause race conditions. Use one terminal per workspace for operations.
2. **Port garbage collection**: Ports from workspaces deleted via filesystem (not `workspace destroy`) require manual cleanup with `scind port gc`.

---

## Success Criteria

1. **Zero application changes**: Existing Docker Compose applications work without modification
2. **Parallel environments**: Can run dev, review, and control simultaneously without conflicts
3. **Predictable naming**: Given workspace and app names, hostnames and aliases are deterministic
4. **Fast iteration**: Switching flavors or regenerating config takes seconds
5. **Debuggable**: Generated files are human-readable; easy to understand what's happening
