<!-- Migrated from specs/scind-prd.md:600-616 -->
<!-- Extraction ID: vision-glossary -->

## Glossary

| Term | Definition |
|------|------------|
| **Alias** | A DNS name on the workspace-internal network |
| **Application** | A Docker Compose-based service that participates in workspaces |
| **Context Detection** | Automatic discovery of workspace/app from current directory |
| **Exported Service** | A named export in `application.yaml`, may map to a different Compose service |
| **Flavor** | A named configuration specifying which compose files to use |
| **Internal Network** | Per-workspace Docker network (`{workspace}-internal`) enabling communication between applications within a workspace. See [ADR-0002](../decisions/0002-two-layer-networking.md). |
| **Manifest** | Generated read-only file showing computed hostnames, ports, and environment variables |
| **Overlay** | The architectural approach where workspace integration is achieved entirely through generated Docker Compose override files, without modifying application source files. See [ADR-0003](../decisions/0003-pure-overlay-design.md). |
| **Override File** | Generated Docker Compose file that adds workspace integration |
| **Port Type** | How a port is routed: `proxied` (through Traefik) or `assigned` (direct port binding) |
| **Project** | Docker Compose project name, formatted as `{workspace}-{application}`, providing namespace isolation for containers. See [ADR-0001](../decisions/0001-docker-compose-project-name-isolation.md). |
| **Protocol** | For proxied types, the traffic protocol: `http`, `https`, or future SNI types |
| **Proxy Network** | Host-level Docker network (`scind-proxy`) connecting the Traefik reverse proxy to services requiring external access. See [ADR-0002](../decisions/0002-two-layer-networking.md). |
| **scind-compose** | Shell function that provides context-aware passthrough to Docker Compose with full tab completion. Automatically injects the appropriate project name and configuration files based on the current workspace and application context. |
| **Service** | A container defined in a Docker Compose file. Distinct from "Exported Service" which is a Scind abstraction for services exposed beyond their application's network. |
| **Service Contract** | The `application.yaml` file defining what an application exports |
| **Single-Application Workspace** | A workspace containing only one application, useful for isolated development of a single project. Allows promoting existing Docker Compose projects to Scind without restructuring. |
| **Visibility** | Flag (`public`/`protected`) indicating intended use; exposed via Docker labels for external tools |
| **Workspace** | An isolated environment containing multiple applications |
