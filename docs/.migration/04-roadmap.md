# Migration Step: Layer 2 - Roadmap

**Prerequisites**: Read `common-instructions.md`
**Estimated Size**: 1 file, approximately 150 lines

---

## Overview

Create `product/roadmap.md` consolidating future considerations from both PRD and technical spec.

**Source Sections**:
- `specs/contrail-prd.md:514-558` - Future Considerations
- `specs/contrail-technical-spec.md:1318-1372` - Future Enhancements

---

## Output File: `product/roadmap.md`

**Source**: `specs/contrail-prd.md:514-558`, `specs/contrail-technical-spec.md:1318-1372`

### Content

```markdown
# Contrail Roadmap

Future considerations and planned enhancements for Contrail.

---

## Port Type Plugins

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

---

## Application Dependencies

**Context**: Some applications may need others to be running first.

**Consideration**: Add dependency ordering:

```yaml
workspace:
  applications:
    app-two:
      depends_on:
        - app-one
```

---

## Shared Volumes

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

---

## Health Checks

**Context**: Starting applications in order isn't sufficient if they need warm-up time.

**Consideration**: Integration with Docker health checks to wait for readiness.

---

## Per-Workspace TLS Overrides

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

<!-- Migrated from specs/contrail-prd.md:514-558, specs/contrail-technical-spec.md:1318-1372 -->
```

---

## Completion Checklist

- [ ] `product/roadmap.md` created
- [ ] All future considerations preserved
- [ ] Implementation priority included
- [ ] Source attribution present
