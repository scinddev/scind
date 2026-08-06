## Generated Manifest

**Location**: `{workspace}/.generated/manifest.yaml` (gitignored)

The manifest is a computed, read-only view of the workspace's current state. It captures all resolved values (hostnames, aliases, environment variables) derived from configuration and state. This serves several purposes:

- **Discoverability**: Humans and tools can inspect one file to understand the workspace topology
- **Tool integration**: Dashboards, DNS updaters, or service discovery tools can consume this structured data
- **Debugging**: Inspect computed hostnames and environment variables without reconstructing from templates
- **At-rest integration**: The manifest is the one topology surface readable **without Docker running** and without invoking Scind — a file another process can watch or read directly, which Docker labels (running containers only) and `--json` (invoke the binary in context) cannot provide

**The manifest is an output, never a staleness input.** Regeneration is decided by [mtime comparison](./workspace-lifecycle.md#staleness-detection) against the source files, plus the revalidation and completeness checks recorded there. Scind does not compare the manifest against configuration to decide whether to regenerate; the manifest records what the last generation produced.

**Cacheable vs. non-cacheable content**: Staleness/regeneration distinguishes **config-derived content** (hostnames, aliases, labels — a pure function of configuration and flavor state, and therefore cacheable) from **live-state-derived content** (assigned `host_port` values and the discovery environment variables that embed them — a function of machine-local runtime state, and therefore *not* cacheable). Live-state-derived fields are re-resolved on every generation even when the config-based staleness check reports up-to-date. See [Workspace Lifecycle: Config-Derived vs. Live-State-Derived Content](./workspace-lifecycle.md#config-derived-vs-live-state-derived-content) and [Port Types](./port-types.md) (assigned values are live-state-derived).

**Freshness contract**: The manifest is consumed **directly** by external tools (dashboards, DNS updaters, service discovery). It must therefore be written **after** assigned-port allocation completes and must reflect the **post-allocation** `host_port` and discovery-variable values — never a pre-allocation placeholder. See [Generation Logic](./workspace-lifecycle.md#generation-logic-workspace-generate), where the allocation step precedes the manifest write.

**Separate host and container ports**: For an `assigned` export, the manifest represents `container_port` (the in-network port the service listens on) and `host_port` (the allocated, host-published port) as **separate** fields. They are distinct values — the container port comes from configuration, the host port is allocated from live state — and consumers must not conflate them.

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
        primary: true                     # Implicit (single proxied export)
        apex_alias: frontend              # Apex internal alias
        ports:
          - type: proxied
            protocol: https
            container_port: 80
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

  shared-db:                              # No proxied export — not apex-eligible, no apex
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

### Portable Resolved-Configuration Snapshot

Beyond the manifest, Scind can emit a **fully-resolved, flattened Compose
configuration** on demand — the merge of base compose files, the generated
override, and any manual overrides, with all variables resolved (the
`docker compose config` equivalent). This can be written to stdout or to an
arbitrary file, producing a portable, self-contained snapshot for external
consumers such as devcontainers or CI that need the effective configuration
without access to Scind or its template inputs. Validated by the Xcind
proof-of-concept. (The specific CLI surface for requesting this snapshot is
defined in the CLI specification.)
