<!-- Migrated from specs/contrail-prd.md:266-271 -->
<!-- Extraction ID: adr-0008-traefik-reverse-proxy -->

# Traefik for Reverse Proxy

**Status**: Accepted

## Context

Need a reverse proxy that can dynamically route to containers.

## Decision

Use Traefik with Docker provider, reading labels from containers.

## Consequences

Traefik's Docker integration allows dynamic routing without config file changes. Labels on containers (added via generated overrides) define routing rules.
