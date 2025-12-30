# Issue Group 3: Conceptual Foundations

**Documents Affected**: PRD + Technical Spec  
**Suggested Order**: 3 of 10 (establishes core concepts before feature details)  
**Estimated Effort**: Small

---

## Overview

These issues involve aligning the PRD and Technical Spec on fundamental concepts. Mostly "document the rationale" or "pick one approach" decisions that establish the foundation for everything else.

---

## Issues

### A-2: Workspace Network Creation Timing Inconsistent

**Severity**: Medium

**Issue**: The documents disagree on when the workspace network is created:

- **PRD** (line 85): Network is "Created by: Workspace initialization"
- **Technical Spec** (line 843): Step 2 of `workspace up` is "Create workspace network if it doesn't exist"

**Questions**:
1. Should the network be created eagerly at `workspace init` time?
2. Should it be created lazily at first `workspace up`?
3. Should both be valid (idempotent creation at both points)?

**Considerations**:
- Eager (at init): Network exists immediately, but workspace might never be started
- Lazy (at up): Network only exists when needed, but adds complexity to up sequence
- Idempotent: Most flexible, but need to document clearly

**Suggested Resolution**: Lazy creation at `up` time with idempotent behavior. Update PRD to match Tech Spec.

**Response**:  
> _[Your response here]_

---

### A-5: Base Environment Variable Protocol Selection Rationale

**Severity**: Low

**Issue**: When an exported service has both HTTP and HTTPS proxied ports, the "base" variables (`CONTRAIL_*_PORT`, `CONTRAIL_*_SCHEME`) default to HTTPS (port 443). The behavior is documented but the rationale is missing.

**Location**: Technical Spec lines 707-712

**Current behavior**:
```
| `proxied` | both | Proxied hostname | 443 | `https` | ✓ | Both |
```

**Questions**:
1. Is the HTTPS preference intentional (security-by-default)?
2. Should it be configurable per-service (e.g., `default_protocol: http`)?
3. Should we just document the rationale explicitly?

**Suggested Resolution**: Document the rationale in Tech Spec:
> When both HTTP and HTTPS are configured, base variables default to HTTPS (port 443) following security-by-default principles. Applications should prefer HTTPS for service-to-service communication. Use protocol-specific variables (`*_HTTP_*`) when HTTP is explicitly required.

**Response**:  
> _[Your response here]_

---

### A-7: Visibility Field Purpose and Future Intent

**Severity**: Low

**Issue**: The `visibility` field (`public`/`protected`) is documented as "documentation only" with no behavioral difference. If purely informational, why formalize it in the schema?

**Location**: 
- PRD lines 99-101
- Technical Spec lines 120-127

**Current documentation** (PRD):
> Each port can have a `visibility` of `public` or `protected`. This is primarily **documentation** to communicate intent to collaborators—it does not change Contrail's behavior.

**Questions**:
1. Is this scaffolding for future features (auth middleware, access controls)?
2. Should we document the intended future use?
3. Should it be removed if truly unused?
4. Could tools like Servlo use this for display purposes?

**Options**:
- A) Keep as-is, document it's for human communication and future extensibility
- B) Remove from schema entirely (YAGNI)
- C) Add concrete behavior now (e.g., protected services get a warning banner in Servlo)

**Response**:  
> _[Your response here]_

---

## Checklist

- [ ] Align PRD and Tech Spec on network creation timing
- [ ] Add rationale for HTTPS-default in base environment variables
- [ ] Decide on visibility field's purpose and document intent
