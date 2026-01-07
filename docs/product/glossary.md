<!-- Migrated from specs/contrail-prd.md:600-616 -->
<!-- Extraction ID: vision-glossary -->

## Terminology

| Term | Definition |
|------|------------|
| **Workspace** | An isolated environment containing multiple applications |
| **Application** | A Docker Compose-based service that participates in workspaces |
| **Flavor** | A named configuration specifying which compose files to use |
| **Service Contract** | The `application.yaml` file defining what an application exports |
| **Exported Service** | A named export in `application.yaml`, may map to a different Compose service |
| **Override File** | Generated Docker Compose file that adds workspace integration |
| **Manifest** | Generated read-only file showing computed hostnames, ports, and environment variables |
| **Port Type** | How a port is routed: `proxied` (through Traefik) or `assigned` (direct port binding) |
| **Protocol** | For proxied types, the traffic protocol: `http`, `https`, or future SNI types |
| **Visibility** | Flag (`public`/`protected`) indicating intended use; exposed via Docker labels for external tools |
| **Alias** | A DNS name on the workspace-internal network |
| **Context Detection** | Automatic discovery of workspace/app from current directory |
