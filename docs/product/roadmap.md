<!-- Migrated from specs/contrail-prd.md:515-559 -->
<!-- Extraction ID: vision-roadmap -->

## Future Considerations

### Port Type Plugins (Future)

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

**Context**: Some applications may need others to be running first.

**Consideration**: Add dependency ordering:
```yaml
workspace:
  applications:
    app-two:
      depends_on:
        - app-one
```

### Shared Volumes (Future)

**Context**: Applications might need to share files (uploads, assets).

**Consideration**: Workspace-level volume definitions that can be mounted into multiple applications.

### Health Checks (Future)

**Context**: Starting applications in order isn't sufficient if they need warm-up time.

**Consideration**: Integration with Docker health checks to wait for readiness.
