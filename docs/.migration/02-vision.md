# Migration Step: Layer 2 — Vision (PRD-Lite)

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 file, approximately 200 lines total

---

## Overview

Extract product vision and problem statement from the PRD. This becomes a condensed "PRD-Lite" focusing on the why and what, not the how.

**Source document**: `specs/contrail-prd.md` lines 1-159 (Introduction through Concepts sections)

---

## File: `product/vision.md`

**Sources**:
- `specs/contrail-prd.md:1-24` (Introduction, Problem Statement)
- `specs/contrail-prd.md:26-80` (Solution, Key Concepts)
- `specs/contrail-prd.md:82-159` (Detailed Concepts)

### Content

```markdown
# Contrail Vision

**Version**: 1.0.0
**Date**: 2024-12

---

## Problem Statement

Modern development often involves multiple applications that need to talk to each other—a frontend, a backend API, a database, maybe a message queue. Docker Compose handles running these services, but it doesn't solve:

1. **Routing**: How does traffic get from your browser to the right container?
2. **Naming**: How do services find each other consistently?
3. **Isolation**: How do you run multiple copies of the same stack without conflicts?
4. **Discoverability**: How does your IDE/shell know what's running and where?

Developers currently solve these problems ad-hoc with manual `/etc/hosts` entries, port number spreadsheets, and tribal knowledge about which ports are "claimed" by which project.

---

## Solution: Workspace-Oriented Development

Contrail is a thin orchestration layer over Docker Compose that provides:

- **Workspace isolation**: Run multiple copies of the same application simultaneously
- **Automatic routing**: Traefik-based reverse proxy with convention-based hostnames
- **Service discovery**: Environment variables and internal DNS for inter-app communication
- **Shell integration**: Context-aware commands that know your current workspace

---

## Core Concepts

### Workspace

A workspace is an isolated development environment containing one or more applications. Think of it as a "project" or "environment" that groups related work.

**Key properties**:
- Has a unique name (e.g., `main`, `feature-auth`, `hotfix-login`)
- Contains references to one or more applications
- Provides an isolated internal network for its applications
- Can coexist with other workspaces on the same machine

**Example**: You might have:
- `main` workspace tracking the main branch of all apps
- `feature-auth` workspace for a cross-app authentication feature
- `review-pr-123` workspace for reviewing a specific PR

### Application

An application is a Docker Compose project that Contrail manages. It's typically a git repository containing a `docker-compose.yaml`.

**Key properties**:
- Lives in its own directory with a `docker-compose.yaml`
- Optionally has an `application.yaml` defining exported services and flavors
- Can participate in multiple workspaces simultaneously
- Remains completely unaware of Contrail (pure overlay design)

### Exported Service

An exported service is a container port that Contrail exposes, either through the reverse proxy or via direct port binding.

**Port types**:
- `proxied`: Traffic routes through Traefik (HTTP/HTTPS)
- `assigned`: Direct host port binding (databases, non-HTTP protocols)

**Visibility**:
- `public`: Accessible via public hostname (Traefik routing)
- `protected`: Accessible only via internal network alias

### Flavor

A flavor is a named configuration variant for an application. It maps to Docker Compose profiles.

**Common patterns**:
- `full`: All services running locally
- `backend-only`: Just the API and database
- `external-db`: Backend using an external database
- `minimal`: Bare minimum for testing

---

## Design Principles

### Pure Overlay

Applications have no knowledge of workspaces. All Contrail integration happens through generated Docker Compose override files. This means:

- Applications can run standalone without Contrail
- No vendor lock-in or special conventions required
- Same application works with any workspace system

### Convention Over Configuration

Names are derived from conventions, not explicit configuration:
- Public hostname: `{workspace}-{app}-{service}.{domain}`
- Internal alias: `{app}-{service}`
- Project name: `{workspace}-{app}`

### Structure vs State

Configuration files describe what *exists* (structure), not what's *active* (state):
- What apps exist → `workspace.yaml`
- Which flavor is active → runtime state
- Which containers are running → Docker state

---

## Related Documents

- [Architecture Overview](../architecture/overview.md) — How components interact
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [ADR-0004: Convention-Based Naming](../decisions/0004-convention-based-naming.md)
- [ADR-0005: Structure vs State Separation](../decisions/0005-structure-vs-state-separation.md)

<!-- Migrated from specs/contrail-prd.md:1-159 -->
```

---

## Completion Checklist

- [ ] `product/vision.md` created
- [ ] Cross-references to ADRs verified
- [ ] Related documents links added

