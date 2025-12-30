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

- [ ] Decide on proxy configuration location (per-project vs global)
- [ ] Decide on auto-start behavior
- [ ] Design `proxy init` command interface
- [ ] Document in PRD
- [ ] Document in Technical Spec with templates
- [ ] Document in CLI Reference
- [ ] Add to Go Stack command mapping
- [ ] Create embedded template for docker-compose.yaml

---

## Response

**Location Decision**:  
> _[Where should proxy config live? Why?]_

**Auto-Start Decision**:  
> _[Should proxy auto-start? Why?]_

**Additional Notes**:  
> _[Any other considerations?]_
