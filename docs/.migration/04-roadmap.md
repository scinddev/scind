# Migration Step: Layer 2 — Roadmap

**Prerequisites**: Read `common-instructions.md`, complete `02-vision.md`
**Estimated Size**: 1 file, approximately 100 lines total

---

## Overview

Extract future considerations and roadmap content from source documents.

**Source documents**:
- `specs/contrail-prd.md` lines 395-412 (Future Considerations section)
- `specs/contrail-go-stack.md` lines 74-85 (Intentionally Excluded section)

---

## File: `product/roadmap.md`

**Sources**:
- `specs/contrail-prd.md:395-412`
- `specs/contrail-go-stack.md:74-85`

### Content

```markdown
# Roadmap & Future Considerations

This document outlines potential future directions for Contrail. Items here are not committed—they represent areas being considered based on user feedback and use cases.

---

## Potential Future Features

### Plugin System for Protocols

**Status**: Under consideration

Currently, Contrail supports `proxied` (HTTP/HTTPS via Traefik) and `assigned` (direct port binding) port types. A plugin system could enable:

- PostgreSQL SNI routing
- MySQL SNI routing
- gRPC-specific handling
- Custom protocol handlers

**Why not now**: Core HTTP/HTTPS use cases must be solid before expanding protocol support.

### Remote Workspaces

**Status**: Under consideration

Allow workspaces to span multiple machines, enabling:

- Shared development environments
- Running heavy services on a remote machine
- Distributed team workflows

**Why not now**: Local development is the primary use case. Remote adds significant complexity.

### GUI Dashboard

**Status**: Under consideration

A web-based dashboard for:

- Viewing workspace status
- Starting/stopping applications
- Viewing logs across services
- Managing flavors

**Why not now**: CLI-first approach ensures scriptability. GUI can be added as a complement.

### Cloud Integration

**Status**: Not planned

Contrail is intentionally a local development tool. Cloud deployment is handled by other tools (Kubernetes, cloud-native Compose, etc.).

---

## Intentionally Excluded

These items were considered and explicitly excluded from scope:

### Kubernetes Support

Contrail targets Docker Compose exclusively. For Kubernetes development workflows, see tools like Tilt, DevSpace, or Skaffold.

### Build System

Contrail does not handle building images. It delegates to `docker compose build` or expects pre-built images. For advanced build workflows, integrate with your existing CI or tools like Buildx.

### Secret Management

Contrail does not manage secrets. Use Docker secrets, environment files, or external secret managers (Vault, AWS Secrets Manager, etc.).

### Production Deployment

Contrail is a development tool. Production deployment should use appropriate production-grade tooling.

---

## Feedback Channels

If you have feature requests or suggestions:

1. Open an issue on the project repository
2. Describe the use case, not just the solution
3. Explain how current tools fall short

---

## Related Documents

- [Vision](./vision.md) — Current product vision
- [Architecture Overview](../architecture/overview.md) — Current architecture

<!-- Migrated from specs/contrail-prd.md:395-412 and specs/contrail-go-stack.md:74-85 -->
```

---

## Completion Checklist

- [ ] `product/roadmap.md` created
- [ ] Update `product/README.md` to list this document

