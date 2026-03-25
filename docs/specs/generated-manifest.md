## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

The manifest is a computed, read-only view of the workspace's current state. It captures all resolved values (hostnames, aliases, environment variables) derived from configuration and state. This serves several purposes:

- **Discoverability**: Humans and tools can inspect one file to understand the workspace topology
- **Tool integration**: Dashboards, DNS updaters, or service discovery tools can consume this structured data
- **Debugging**: Inspect computed hostnames and environment variables without reconstructing from templates
- **Caching**: Scind can compare the manifest against configuration to determine if regeneration is needed

```yaml
# AUTO-GENERATED - Computed from configuration and state
workspace:
  name: dev
  network: dev-internal
proxy:
  domain: scind.test

applications:
  frontend:
    flavor: full
    project: dev-frontend
    exported_services:
      web:
        service: web
        alias: frontend-web
        primary: true                     # Implicit (single export)
        apex_alias: frontend              # Apex internal alias
        ports:
          - type: proxied
            protocol: https
            container_port: 443
            visibility: public
            hostname: dev-frontend-web.scind.test
            apex_hostname: dev-frontend.scind.test
        environment:
          SCIND_FRONTEND_WEB_URL: https://dev-frontend-web.scind.test
          # ... additional export variables ...
          SCIND_FRONTEND_APEX_HOST: dev-frontend.scind.test
          SCIND_FRONTEND_APEX_PORT: 443
          SCIND_FRONTEND_APEX_SCHEME: https
          SCIND_FRONTEND_APEX_URL: https://dev-frontend.scind.test

  shared-db:                              # Multi-export, no primary — no apex
    flavor: default
    project: dev-shared-db
    exported_services:
      db:
        service: postgres
        alias: shared-db-db
        ports:
          - type: assigned
            container_port: 5432
            host_port: 5432
            visibility: protected
        environment:
          SCIND_SHARED_DB_DB_HOST: shared-db-db
          SCIND_SHARED_DB_DB_PORT: 5432
```
