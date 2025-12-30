# Issue Group 11: Docker Labels Specification

**Documents Affected**: Technical Spec + Go Stack
**Suggested Order**: 11 of 11 (after existing issues are resolved)
**Estimated Effort**: Medium

---

## Overview

Docker labels are used throughout Contrail for workspace metadata, service discovery, and tool integration, but they are not formally specified in one place. This issue group consolidates the Docker label requirements.

---

## Issues

### L-1: Docker Labels Not Formally Specified

**Severity**: Medium

**Issue**: Docker labels are referenced in multiple places but there is no authoritative specification for:
- The complete list of labels Contrail generates
- Label naming conventions and namespacing
- Which labels are required vs optional
- How external tools (like Servlo) should consume these labels
- How the workspace discovery fallback (M-1) uses labels to reconstruct the registry

**Current references to labels**:
- Technical Spec generated override example shows `workspace.name`, `workspace.application`, `workspace.exported_service`, `workspace.visibility`
- M-1 and M-7 responses indicate labels should be used to bootstrap/update the workspace registry

**Questions**:
1. Should labels use a `contrail.` prefix for namespacing (e.g., `contrail.workspace.name`)?
2. What additional labels might be needed for workspace discovery?
   - `contrail.workspace.path` (to reconstruct registry)?
   - `contrail.workspace.config_hash` (to detect stale state)?
3. Should there be a label schema version for future compatibility?
4. Should labels be documented in a dedicated section or as part of the generated override file documentation?

**Suggested Resolution**:
Add a "Docker Labels" section to the Technical Spec that formally specifies:
- All generated labels with their purpose
- Naming convention (recommend `contrail.*` namespace)
- Labels required for workspace discovery fallback
- Guidance for external tool integration

**Response**:
> _[Your response here]_

---

## Checklist

- [ ] Define complete Docker label schema
- [ ] Document label naming conventions
- [ ] Add labels required for workspace discovery (M-1 fallback)
- [ ] Add Docker Labels section to Technical Spec
- [ ] Update generated override example with any new labels
- [ ] Add label documentation for external tool integration
