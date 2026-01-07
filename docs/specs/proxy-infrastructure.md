<!-- Migrated from specs/contrail-technical-spec.md:1030-1050 -->
<!-- Extraction ID: spec-proxy-infrastructure -->

## Proxy Layer

### Traefik Configuration

**proxy/docker-compose.yaml:**

See [Configuration Schemas](./configuration-schemas.md#proxy-infrastructure) for the generated docker-compose.yaml structure. The dashboard port (`8080:8080`) is only included when `proxy.dashboard.enabled` is true (default). The `--api.dashboard` flag is set based on the same configuration.

### DNS Configuration

For local development, configure DNS resolution for the workspace domains. Options include:

1. **dnsmasq**: Route all `*.contrail.test` to `127.0.0.1`
   ```
   address=/contrail.test/127.0.0.1
   ```
2. **/etc/hosts**: Manual entries for each hostname
3. **Local DNS server**: More complex but flexible

**Note**: The `.test` TLD is reserved by RFC 2606 for testing purposes and will not conflict with real domains or mDNS (unlike `.local`).
