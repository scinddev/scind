# Issue Group 9: New Feature — `proxy init`

**Documents Affected**: PRD + Technical Spec + CLI Reference + Go Stack  
**Suggested Order**: 9 of 10 (new feature after existing features are clarified)  
**Estimated Effort**: Medium

---

## Overview

The proxy layer is a prerequisite for Contrail's proxied services, but there's no documented way to bootstrap it. This group addresses adding a `proxy init` command.

---

## Issues

### M-2: No `contrail proxy init` Command

**Severity**: Medium

**Issue**: The proxy layer assumes `proxy/docker-compose.yaml` exists with Traefik configuration (Technical Spec lines 742-771), but there's no command to bootstrap this infrastructure.

**Current State**:
- Tech Spec shows the expected Traefik docker-compose.yaml
- PRD mentions proxy layer exists
- CLI Reference has `proxy up`, `proxy down`, `proxy status` but no `proxy init`
- Users must manually create the proxy configuration

**Questions to Resolve**:

**1. Where should proxy configuration live?**

| Option | Location | Pro | Con |
|--------|----------|-----|-----|
| A | Per-project `proxy/` | Each project self-contained | Duplicated across projects |
| B | Global `~/.config/contrail/proxy/` | Single source of truth | Shared across all projects |
| C | Configurable | Flexibility | Complexity |

**2. What should `proxy init` create?**

Minimum viable:
```
proxy/
├── docker-compose.yaml    # Traefik service definition
└── .env                   # Optional environment variables
```

Extended:
```
proxy/
├── docker-compose.yaml
├── traefik.yaml          # Static configuration
├── dynamic/              # Dynamic configuration directory
└── certs/                # TLS certificates (future)
```

**3. Should proxy be auto-started on first `workspace up`?**

| Option | Behavior | Pro | Con |
|--------|----------|-----|-----|
| A | Manual | Explicit control | Extra step for users |
| B | Auto-start | Seamless | Surprises users, permissions issues? |
| C | Prompt | Best of both | Interrupts automation |

**4. What if proxy already exists?**

- Error and exit?
- Skip with message?
- Offer to overwrite with `--force`?

---

## Proposed Design

### Command: `contrail proxy init`

```bash
contrail proxy init [flags]
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--path` | Directory to create proxy in (default: `~/.config/contrail/proxy/`) |
| `--force` | Overwrite existing configuration |
| `--domain` | Set proxy domain (default: `contrail.test`) |

**Behavior**:
1. Check if proxy configuration already exists
   - If exists and no `--force`: error with message
   - If exists and `--force`: backup and overwrite
2. Create proxy directory structure
3. Generate `docker-compose.yaml` with Traefik configuration
4. Create `proxy` Docker network if it doesn't exist
5. Output next steps

**Example Output**:
```
$ contrail proxy init
Created proxy configuration at ~/.config/contrail/proxy/

Next steps:
  1. Configure DNS for *.contrail.test → 127.0.0.1
     (See: https://docs.contrail.dev/dns-setup)
  2. Start the proxy:
     contrail proxy up

$ contrail proxy init
Error: Proxy configuration already exists at ~/.config/contrail/proxy/
Use --force to overwrite, or --path to create elsewhere.
```

---

## Documents to Update

### PRD
- Add `proxy init` to CLI quick reference
- Mention bootstrapping in Architecture section

### Technical Spec
- Add `proxy init` to Operations section
- Document proxy configuration location decision
- Add generated file templates

### CLI Reference
- Add full `proxy init` command documentation
- Add to proxy commands section (before `proxy up`)

### Go Stack
- Add `proxyInitCmd` to command mapping
- Add to implementation priority (Phase 1 or 2?)

---

## Implementation Checklist

- [x] Decide on proxy configuration location (per-project vs global) → Global at `~/.config/contrail/proxy/`
- [x] Decide on auto-start behavior → Auto-start on `workspace up`
- [x] Design `proxy init` command interface → `--force`, `--domain`, `--path` flags
- [x] Document in PRD (Quick Reference, Architecture section)
- [x] Document in Technical Spec with templates (Proxy Infrastructure section with docker-compose.yaml and traefik.yaml)
- [x] Document in CLI Reference (full `proxy init` command documentation)
- [x] Add to Go Stack command mapping and scaffolding (proxyInitCmd with full proxy.go)
- [x] Create embedded template for docker-compose.yaml (in Tech Spec)

---

## Response

**Location Decision**:
> Option B (Global `~/.config/contrail/proxy/`). The proxy should be a single shared instance across all workspaces. It could be managed as a Docker Compose project or directly via Docker commands. Using Docker Compose means users could manually edit the configuration and break things, but that's acceptable if well-documented—and `--force` on `proxy init` provides a recovery path.

**Auto-Start Decision**:
> Option B (Auto-start). The proxy should be brought up automatically when someone brings up a workspace for the first time. This was partially addressed in Group 7 (C-4) where we noted that workspace/application commands should automatically bring up the proxy if needed.

**Additional Notes**:
> - `proxy init` should create the complete Docker Compose project including core Traefik configuration
> - The `dynamic/` directory is for Traefik dynamic configuration and `certs/` for TLS certificates—include both for a complete setup
> - The proxy should be completely self-contained and managed by Contrail; users should be mostly unaware of the proxy internals
> - All wiring should happen via Docker labels on workspaces and applications
> - The `--domain` flag should be included with `contrail.test` as default, allowing users to specify a custom domain (e.g., for wildcard DNS they want to share externally)
> - If proxy config already exists, error by default; allow `--force` to overwrite (useful as a recovery mechanism if manual edits break things)

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
