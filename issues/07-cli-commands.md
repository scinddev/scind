# Issue Group 7: CLI Command Implementation

**Documents Affected**: CLI Reference + Go Stack  
**Suggested Order**: 7 of 10 (finalize command details after concepts are settled)  
**Estimated Effort**: Small

---

## Overview

These issues concern CLI command naming consistency and behavior details. Best addressed after conceptual decisions are made.

---

## Issues

### C-4: Proxy Commands Named Differently in CLI vs Go Stack

**Severity**: Medium

**Issue**: The CLI Reference documents one naming convention, the Go Stack uses another:

| CLI Reference | Go Stack |
|---------------|----------|
| `contrail proxy up` | `proxyStartCmd` |
| `contrail proxy down` | `proxyStopCmd` |
| `contrail proxy restart` | `proxyRestartCmd` |

**Location**:
- CLI Reference lines 849-881
- Go Stack lines 763-765

**Options**:

**A) Use `up`/`down` (match workspace commands)**
- Consistent with `workspace up`, `workspace down`
- Semantic: "bring up" / "tear down"
- Update Go Stack to `proxyUpCmd`, `proxyDownCmd`

**B) Use `start`/`stop` (match Docker)**
- More intuitive for a single service
- Docker uses `docker start`, `docker stop`
- Update CLI Reference to `proxy start`, `proxy stop`

**C) Support Both as Aliases**
- `proxy up` = `proxy start`
- `proxy down` = `proxy stop`
- Pro: Maximum flexibility
- Con: More surface area to document/test

**Suggested Resolution**: Option A (`up`/`down`) for consistency with the rest of Contrail's vocabulary.

**Response**:  
> _[Your response here]_

---

### A-6: Repeatable `--app` Flag Interaction with Context Detection

**Severity**: Medium

**Issue**: CLI Reference (line 254) says `--app` is "repeatable" for `workspace up`:

```bash
contrail workspace up -a app-one -a app-two
```

But how does this interact with context detection?

**Scenario**:
```bash
cd ~/workspaces/dev/app-one  # Context detects: workspace=dev, app=app-one
contrail workspace up -a app-two
```

**Questions**:
1. Does `-a app-two` replace the context-detected app?
2. Does it add to the context-detected app (start both app-one and app-two)?
3. Should context detection be skipped entirely when any `-a` is provided?

**Options**:

**A) Explicit Replaces Context**
- Any `-a` flag completely overrides context detection
- `contrail up -a app-two` starts only app-two
- Pro: Explicit is explicit, no surprises
- Con: Loses convenience of "also start these"

**B) Explicit Adds to Context**
- `-a` flags add to context-detected app
- From app-one dir: `contrail up -a app-two` starts app-one AND app-two
- Pro: Convenient for "start me plus dependencies"
- Con: Confusing—how to start ONLY app-two from app-one dir?

**C) Context Only When No Flags**
- If any `-a` provided, ignore context entirely
- If no `-a` provided and in app dir, use context
- If no `-a` provided and in workspace root, start all apps
- Pro: Clear mental model
- Con: Less flexible

**Suggested Resolution**: Option A (Explicit Replaces Context). Document:
> When `--app` is specified, context-detected application is ignored. To start multiple specific apps, use multiple `-a` flags: `contrail up -a app-one -a app-two`

**Response**:  
> _[Your response here]_

---

### A-8: DNS Validation in `contrail doctor` May Be Insufficient

**Severity**: Low

**Issue**: `contrail doctor` checks domain resolution (CLI Reference lines 1117-1137):

```
✓ Domain resolution: contrail.test → 127.0.0.1
```

But if using wildcard DNS via dnsmasq (`address=/contrail.test/127.0.0.1`), checking just the base domain works—but checking a random subdomain would be more thorough validation.

**Current check**: `contrail.test → 127.0.0.1`
**Better check**: `test-check-{random}.contrail.test → 127.0.0.1`

**Questions**:
1. Is wildcard DNS the expected setup?
2. Should doctor check both base and a random subdomain?
3. Should it check an actual expected hostname like `dev-app-one-web.contrail.test`?

**Suggested Resolution**: Check both:
```
✓ Domain resolution: contrail.test → 127.0.0.1
✓ Wildcard resolution: check-1735567890.contrail.test → 127.0.0.1
```

If base resolves but wildcard doesn't, warn:
```
⚠ Wildcard DNS not configured. Individual hostnames may not resolve.
  Configure dnsmasq: address=/contrail.test/127.0.0.1
```

**Response**:  
> _[Your response here]_

---

## Checklist

- [ ] Align proxy command naming between CLI Reference and Go Stack
- [ ] Document `--app` flag behavior with context detection
- [ ] Add explicit documentation about explicit flags overriding context
- [ ] Enhance `doctor` DNS check to verify wildcard resolution
