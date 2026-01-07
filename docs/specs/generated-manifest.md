<!-- Migrated from specs/scind-technical-spec.md:708-765 -->
<!-- Extraction ID: spec-generated-manifest -->

## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

The manifest is a computed, read-only view of the workspace's current state. It captures all resolved values (hostnames, aliases, environment variables) derived from configuration and state. This serves several purposes:

- **Discoverability**: Humans and tools can inspect one file to understand the workspace topology
- **Tool integration**: Dashboards, DNS updaters, or service discovery tools can consume this structured data
- **Debugging**: Inspect computed hostnames and environment variables without reconstructing from templates
- **Caching**: Scind can compare the manifest against configuration to determine if regeneration is needed

```yaml
# AUTO-GENERATED - Computed from configuration and state
# Generated: 2024-12-27T10:30:00Z

workspace:
  name: dev
  network: dev-internal

proxy:
  domain: scind.test

applications:
  app-one:
    flavor: default
    project: dev-app-one
    exported_services:
      web:
        service: web                      # Underlying Compose service
        alias: app-one-web                # Internal network alias
        ports:
          - type: proxied
            protocol: https
            container_port: 443
            visibility: public
            hostname: dev-app-one-web.scind.test
        environment:
          SCIND_APP_ONE_WEB_HOST: dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_PORT: 443
          SCIND_APP_ONE_WEB_SCHEME: https
          SCIND_APP_ONE_WEB_URL: https://dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_HTTPS_HOST: dev-app-one-web.scind.test
          SCIND_APP_ONE_WEB_HTTPS_PORT: 443
          SCIND_APP_ONE_WEB_HTTPS_URL: https://dev-app-one-web.scind.test

      db:
        service: postgres
        alias: app-one-db
        ports:
          - type: assigned
            container_port: 5432
            host_port: 5432
            visibility: protected
        environment:
          SCIND_APP_ONE_DB_HOST: app-one-db
          SCIND_APP_ONE_DB_PORT: 5432
```
