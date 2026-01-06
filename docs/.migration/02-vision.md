# Migration Step: Layer 2 - Vision

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 main file, approximately 400 lines total

---

## Overview

Create `product/vision.md` using the Lean PRD template, consolidating content from `specs/contrail-prd.md`.

**Source Sections**:
- Executive Summary (lines 9-14)
- Problem Statement (lines 16-45)
- Product Vision (lines 47-57)
- Core Concepts (lines 59-108)
- Success Criteria (lines 580-588)
- Non-Goals (lines 561-571)
- Known Limitations (lines 573-578)
- Terminology (lines 600-616)

---

## Output File: `product/vision.md`

**Source**: `specs/contrail-prd.md` (multiple sections)

### Content

```markdown
# Contrail Product Vision

**Version**: 0.5.2
**Status**: Design Phase

---

## Executive Summary

Contrail is a workspace orchestration system for Docker Compose that enables developers to run multiple isolated instances of multi-application stacks simultaneously on a single host. It solves the problem of needing complete, independent environments for development, code review, and testing without requiring Kubernetes or cloud infrastructure.

The name "Contrail" evokes the trails left by aircraft - parallel paths that don't intersect, much like the isolated workspaces the system creates.

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

### The Gap

There is no existing tool that orchestrates **multiple isolated instances of multi-application Docker Compose stacks** with **generated integration configuration** while keeping **applications workspace-agnostic**.

---

## Product Vision

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
  app-one (project: dev-app-one)
  app-two (project: dev-app-two)
  app-three (project: dev-app-three)
    all connected via: dev-internal network
```

### Application

A Docker Compose-based service that can participate in workspaces. Applications define a **service contract** (`application.yaml`) that declares what they export to the workspace, but their core `docker-compose.yaml` remains workspace-agnostic.

### Single-Application Workspace

Contrail supports workspaces where the workspace directory is also the application directory. This allows promoting existing Docker Compose projects without restructuring:

```bash
cd ~/my-existing-project
contrail workspace init --workspace=dev
contrail app init --app=myapp
```

### Port Types and Proxying

| Type | Protocol | Behavior | External Access | Internal Access | Use Case |
|------|----------|----------|-----------------|-----------------|----------|
| **proxied** | `https` | Traefik HTTPS proxy | Via hostname (port 443) | Via alias | Web frontends, public APIs |
| **proxied** | `http` | Traefik HTTP proxy | Via hostname (port 80) | Via alias | HTTP services, redirects |
| **assigned** | - | Direct port binding | Via assigned host port | Via alias | Databases, caches, debug ports |
| **proxied** | `tcp`, etc. | SNI TCP proxy (future) | Via hostname | Via alias | Database GUIs, external tools |

Services not listed in `exported_services` remain **private** - only accessible within the application's own Docker Compose network.

### Visibility

Each port can have a `visibility` of `public` or `protected` (defaults to `protected` if not specified). This is primarily **documentation** to communicate intent to collaborators - it does not change Contrail's core behavior.

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

1. **Concurrent operations**: Running multiple Contrail commands simultaneously (e.g., two terminals running `workspace up`) may cause race conditions. Use one terminal per workspace for operations.
2. **Port garbage collection**: Ports from workspaces deleted via filesystem (not `workspace destroy`) require manual cleanup with `contrail port gc`.

---

## Terminology

| Term | Definition |
|------|------------|
| **Workspace** | An isolated environment containing multiple applications |
| **Application** | A Docker Compose-based service that participates in workspaces |
| **Flavor** | A named configuration specifying which compose files to use |
| **Service Contract** | The `application.yaml` file defining what an application exports |
| **Exported Service** | A named export in `application.yaml`, may map to a different Compose service |
| **Override File** | Generated Docker Compose file that adds workspace integration |
| **Manifest** | Generated read-only file showing computed hostnames, ports, and environment variables |
| **Port Type** | How a port is routed: `proxied` (through Traefik) or `assigned` (direct port binding) |
| **Protocol** | For proxied types, the traffic protocol: `http`, `https`, or future SNI types |
| **Visibility** | Flag (`public`/`protected`) indicating intended use; exposed via Docker labels for external tools |
| **Alias** | A DNS name on the workspace-internal network |
| **Context Detection** | Automatic discovery of workspace/app from current directory |

<!-- Migrated from specs/contrail-prd.md -->
```

---

## Also Create: `product/README.md`

```markdown
# Product Documentation

This directory contains product-level documentation for Contrail.

## Contents

| Document | Description |
|----------|-------------|
| [vision.md](./vision.md) | Product vision, problem statement, and core concepts |
| [comparison.md](./comparison.md) | How Contrail compares to alternative solutions |
| [roadmap.md](./roadmap.md) | Future considerations and planned enhancements |
```

---

## Completion Checklist

- [ ] `product/vision.md` created
- [ ] `product/README.md` created
- [ ] All content preserved (not summarized)
- [ ] Source attribution present
