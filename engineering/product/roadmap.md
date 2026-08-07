<!-- Migrated from specs/scind-prd.md:515-559 -->
<!-- Extraction ID: vision-roadmap -->

## Future Considerations

> **Note**: This roadmap describes planned features without specific version targets. Features are prioritized based on community feedback and project needs rather than fixed release schedules.

### Port Type Plugins (Future)

*Extends [ADR-0007: Port Type System](../decisions/0007-port-type-system.md)*

**Context**: Different services need different proxying strategies—HTTP/HTTPS is handled by Traefik's HTTP routers, but databases need TCP routing with SNI.

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

### Application Dependencies (Future)

*Related to [ADR-0010: Up/Down Command Semantics](../decisions/0010-up-down-command-semantics.md)*

**Context**: Some applications may need others to be running first.

**Consideration**: Add dependency ordering:
```yaml
workspace:
  applications:
    backend:
      depends_on:
        - shared-db
```

### Shared Volumes (Future)

**Context**: Applications might need to share files (uploads, assets).

**Consideration**: Workspace-level volume definitions that can be mounted into multiple applications.

### Health Checks (Future)

**Context**: Starting applications in order isn't sufficient if they need warm-up time.

**Consideration**: Integration with Docker health checks to wait for readiness.
