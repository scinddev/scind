# Contrail Roadmap

Future considerations and planned enhancements for Contrail.

---

## Potential Future Features

### Port Type Plugins

**Status**: Under consideration

**Context**: Different services need different proxying strategies - HTTP/HTTPS is handled by Traefik's HTTP routers, but databases need TCP routing with SNI.

**Consideration**: Plugin system where protocols can register handlers:

```yaml
exported_services:
  db:
    service: postgres
    ports:
      - type: proxied
        protocol: postgresql          # Plugin handles this protocol
        port: 5432
        visibility: public
```

Plugins would generate appropriate Traefik configuration (TCP routers, SNI rules) for their protocol.

**Why not now**: Core HTTP/HTTPS use cases must be solid before expanding protocol support.

---

### Application Dependencies

**Status**: Under consideration

**Context**: Some applications may need others to be running first.

**Consideration**: Add dependency ordering:

```yaml
workspace:
  applications:
    app-two:
      depends_on:
        - app-one
```

**Why not now**: Most use cases can work around this with manual ordering. Adding dependency graphs increases complexity significantly.

---

### Shared Volumes

**Status**: Under consideration

**Context**: Applications might need to share files (uploads, assets).

**Consideration**: Workspace-level volume definitions that can be mounted into multiple applications:

```yaml
# workspace.yaml
volumes:
  shared-uploads:
    driver: local

applications:
  app-one:
    volumes:
      - shared-uploads:/app/uploads
  app-two:
    volumes:
      - shared-uploads:/app/uploads:ro
```

**Why not now**: This can be achieved with Docker named volumes today. A Contrail-native solution would add convenience but isn't blocking any use cases.

---

### Health Checks

**Status**: Under consideration

**Context**: Starting applications in order isn't sufficient if they need warm-up time.

**Consideration**: Integration with Docker health checks to wait for readiness.

**Why not now**: Docker Compose already supports `depends_on` with `condition: service_healthy`. Contrail can document this pattern without needing native support.

---

### Per-Workspace TLS Overrides

**Status**: Under consideration

**Context**: Allow workspaces to override the global TLS configuration for testing with different certificates.

**Consideration**:

```yaml
# workspace.yaml
workspace:
  name: dev
  tls:
    mode: custom
    cert_file: ./certs/dev.crt
    key_file: ./certs/dev.key
```

**Why not now**: The global TLS configuration handles most development needs. Custom certificates per workspace is an edge case.

---

### Remote Workspaces

**Status**: Under consideration

Allow workspaces to span multiple machines, enabling:

- Shared development environments
- Running heavy services on a remote machine
- Distributed team workflows

**Why not now**: Local development is the primary use case. Remote adds significant complexity around networking, authentication, and state synchronization.

---

### GUI Dashboard

**Status**: Under consideration

A web-based dashboard for:

- Viewing workspace status
- Starting/stopping applications
- Viewing logs across services
- Managing flavors

**Why not now**: CLI-first approach ensures scriptability and works in any environment. GUI can be added as a complement once the CLI is stable.

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

### Cloud Integration

Contrail is intentionally a local development tool. Cloud deployment is handled by other tools (Kubernetes, cloud-native Compose, etc.).

---

## Implementation Priority

### Phase 1: Core CLI Structure
1. Root command with context detection
2. Workspace commands (init, up, down, status)
3. App commands (init, up, down, status)
4. Basic configuration loading

### Phase 2: Override Generation
1. Template system for hostnames/aliases
2. Override file generation
3. Manifest generation
4. Traefik label generation

### Phase 3: Shell Integration
1. `compose-prefix` command
2. `init-shell` command with embedded scripts
3. Shell completion for flags

### Phase 4: Polish
1. Port management commands
2. Proxy management commands
3. `doctor` command
4. `validate` command
5. `open` and `urls` commands

### Future: Plugins
1. Define plugin interface in `pkg/plugin`
2. Integrate go-plugin
3. Extract protocol handlers to plugins

---

## Feedback Channels

If you have feature requests or suggestions:

1. Open an issue on the project repository
2. Describe the use case, not just the solution
3. Explain how current tools fall short

---

## Related Documents

- [Vision](./vision.md) — Current product vision
- [Comparison](./comparison.md) — How Contrail compares to alternatives

<!-- Migrated from specs/contrail-prd.md:514-569, specs/contrail-technical-spec.md:1318-1372 -->
