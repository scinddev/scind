# Architecture Overview

**Version**: 1.0.0
**Date**: 2024-12

---

## System Overview

Contrail is an orchestration layer that sits between the developer and Docker Compose. It does not replace Docker Compose—it generates configuration that Compose consumes.

```
┌─────────────────────────────────────────────────────────────────┐
│                        Developer                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Contrail CLI                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Context   │  │   Config    │  │      Generator          │  │
│  │  Detection  │  │   Loader    │  │  (Override Files)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Docker Compose                                │
│         (with generated override files)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Docker Engine                               │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    Containers                                ││
│  └─────────────────────────────────────────────────────────────┘│
│  ┌───────────────────────┐  ┌─────────────────────────────────┐ │
│  │   contrail-proxy      │  │   {workspace}-internal          │ │
│  │   (network)           │  │   (network per workspace)       │ │
│  └───────────────────────┘  └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

---

## Core Components

### CLI (`contrail`)

The command-line interface is the primary user interaction point.

**Responsibilities**:
- Parse commands and options
- Detect context from current directory
- Load and validate configuration
- Delegate to internal components
- Invoke Docker Compose with generated overrides

**Key design**: The CLI is stateless. All state lives in configuration files, Docker, and git.

### Context Detector

Determines the current workspace and application based on directory location.

**Algorithm**:
1. Walk up from current directory
2. Look for `workspace.yaml` (workspace root)
3. Look for `application.yaml` or `docker-compose.yaml` (application root)
4. Return detected context or empty

**See**: [Context Detection Spec](../specs/context-detection.md)

### Configuration Loader

Loads and validates the three configuration schemas.

**Schemas**:
- `proxy.yaml`: Global Traefik and TLS settings
- `workspace.yaml`: Workspace definition and application list
- `application.yaml`: Application-specific settings (flavors, services)

**See**: [Configuration Schemas Spec](../specs/configuration-schemas.md)

### Override Generator

Generates Docker Compose override files that integrate applications into workspaces.

**Generates**:
- Network attachments (proxy and internal)
- Traefik labels for routing
- Environment variables for service discovery
- Project name settings

**Output location**: `.generated/` directory (gitignored)

**See**: [Generated Override Files Spec](../specs/generated-override-files.md)

---

## Network Architecture

Contrail creates a two-layer network topology:

```
┌─────────────────────────────────────────────────────────────────┐
│                         Host                                     │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   contrail-proxy                           │  │
│  │                   (host-wide network)                      │  │
│  │                                                            │  │
│  │   ┌─────────┐                                              │  │
│  │   │ Traefik │◄──── HTTP/HTTPS from browser                 │  │
│  │   └────┬────┘                                              │  │
│  │        │                                                   │  │
│  │        ▼                                                   │  │
│  │   Routes to containers with `visibility: public`          │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              main-internal (workspace network)             │  │
│  │                                                            │  │
│  │   ┌─────────┐    ┌─────────┐    ┌─────────┐               │  │
│  │   │ app-one │◄──►│ app-two │◄──►│ app-db  │               │  │
│  │   └─────────┘    └─────────┘    └─────────┘               │  │
│  │                                                            │  │
│  │   All containers can reach each other via aliases          │  │
│  │   e.g., app-one-web, app-two-web, app-db-postgres          │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │            feature-internal (another workspace)            │  │
│  │                                                            │  │
│  │   ┌─────────┐    ┌─────────┐                               │  │
│  │   │ app-one │    │ app-two │   (isolated from main)        │  │
│  │   └─────────┘    └─────────┘                               │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

**Key points**:
- `contrail-proxy`: Single network connecting Traefik to all public services
- `{workspace}-internal`: Per-workspace network for internal communication
- Workspaces are isolated—containers in different workspaces cannot communicate

**See**: [Two-Layer Networking ADR](../decisions/0002-two-layer-networking.md)

---

## Data Flow

### `contrail up` Flow

```
1. Parse command options
2. Detect context (workspace, app) if not specified
3. Load configuration (proxy.yaml, workspace.yaml, application.yaml)
4. For each application:
   a. Generate override file with:
      - Network attachments
      - Traefik labels
      - Environment variables
   b. Write to .generated/
5. Invoke: docker compose -f docker-compose.yaml -f .generated/override.yaml up
6. Report status
```

### Request Flow (HTTP)

```
1. Browser requests https://main-app-one-web.test
2. DNS resolves to 127.0.0.1 (via /etc/hosts or local DNS)
3. Traefik receives request on port 443
4. Traefik matches Host header to routing rule
5. Traefik forwards to container on internal port
6. Container responds
7. Response flows back through Traefik to browser
```

---

## Related Documents

- [ADR-0002: Two-Layer Networking](../decisions/0002-two-layer-networking.md)
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md)
- [ADR-0008: Traefik Reverse Proxy](../decisions/0008-traefik-reverse-proxy.md)
- [Proxy Infrastructure Spec](../specs/proxy-infrastructure.md)
- [Generated Override Files Spec](../specs/generated-override-files.md)

<!-- Migrated from specs/contrail-technical-spec.md:1-200 -->
