# Issue Group 17: DNS & Networking Details

**Documents Affected**: Technical Spec + CLI Reference
**Suggested Order**: 17 of 18 (implementation details, can defer)
**Estimated Effort**: Small

---

## Overview

These issues concern DNS resolution behavior and Docker networking edge cases.

---

## Issues

### N-15: DNS Resolution Behavior in `doctor` Command

**Severity**: Low

**Issue**: The CLI Reference (lines 1243-1252) documents DNS checking in `doctor` but leaves several details unspecified:

```
✓ Domain resolution: contrail.test → 127.0.0.1
✓ Workspace domains:
  - dev-app-one-web.contrail.test → 127.0.0.1
```

**Questions**:
1. What DNS resolver is used? System default? A specific library?
2. Does it respect `/etc/hosts` entries?
3. How long is the timeout for DNS queries?
4. What happens in air-gapped/offline environments?

**Suggested Resolution**: Document that:
- System DNS resolver is used (respects `/etc/hosts` and `/etc/resolv.conf`)
- Timeout is 5 seconds per query
- Offline environments should use `/etc/hosts` or dnsmasq with local config

**Your Response**:
> Use suggested resolution: document system DNS resolver, 5 second timeout, and offline environment guidance.

---

### N-16: Volume Naming Collision Potential

**Severity**: Low

**Issue**: Docker Compose creates volumes with the project name prefix (e.g., `dev-app-one_postgres_data`). The collision warning in Technical Spec (lines 1256-1257) only mentions Traefik router names:

> "Creative naming that produces identical router names... could cause routing conflicts."

The same issue applies to volumes and networks. Example collision:
- Workspace `dev-app` with app `one` → project name `dev-app-one`
- Workspace `dev` with app `app-one` → project name `dev-app-one`

Both would create `dev-app-one_postgres_data` volume.

**Questions**:
1. Should the collision warning in Tech Spec be expanded to cover volumes and networks?
2. Should Contrail validate for naming collisions at `workspace init` time?
3. Is this edge case rare enough to just document?

**Your Response**:
> Option 1: Expand the collision warning in Tech Spec to cover volumes and networks.

---

### N-17: Existing Proxy Network Conflict

**Severity**: Low

**Issue**: Technical Spec line 307 states:

> "proxy up: Starts the Traefik container (creates `proxy` network if needed)"

But what if the `proxy` network already exists and was created by a different tool with different settings (e.g., different subnet, different driver)?

**Questions**:
1. Should `proxy up` validate the existing network's configuration?
2. Should it warn if the network wasn't created by Contrail?
3. Should it have a `--recreate` flag for recovery?

**Your Response**:
> Add `--recreate` flag to `proxy up` for recovery. Also rename the network from `proxy` to `contrail-proxy` for uniqueness and to avoid conflicts with other tools.

---

## Checklist

- [x] Document DNS resolution behavior in doctor command — added to CLI Reference
- [x] Expand naming collision warning to cover volumes/networks — added to Tech Spec Naming Conventions
- [x] Document existing proxy network handling — added `--recreate` flag, renamed network to `contrail-proxy`
