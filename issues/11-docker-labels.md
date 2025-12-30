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
> Use `contrail.` prefix for all labels, with kebab-case for multi-word segments (following Traefik convention). No schema versioning needed for now.
>
> **Label schema:**
>
> Context labels (on all app containers):
> ```
> contrail.workspace.name=review
> contrail.workspace.path=/Users/beau/workspaces/review
> contrail.app.name=app-one
> contrail.app.path=/Users/beau/workspaces/review/app-one
> ```
>
> Export labels (keyed by export name for consistency):
> ```
> # Proxied export with http + https
> contrail.export.web.host=review-app-one-web.contrail.test
> contrail.export.web.proxy.http.visibility=public
> contrail.export.web.proxy.http.url=http://review-app-one-web.contrail.test
> contrail.export.web.proxy.https.visibility=public
> contrail.export.web.proxy.https.url=https://review-app-one-web.contrail.test
>
> # Assigned port export
> contrail.export.debug.host=review-app-one-debug.contrail.test
> contrail.export.debug.port.9000.visibility=protected
> contrail.export.debug.port.9000.assigned=9003
> ```
>
> Proxy container labels:
> ```
> contrail.managed=true
> contrail.component=proxy
> ```

---

## Checklist

- [x] Define complete Docker label schema
- [x] Document label naming conventions (`contrail.` prefix, kebab-case)
- [x] Add labels required for workspace discovery (M-1 fallback) — `contrail.workspace.path`, `contrail.app.path`
- [x] Add Docker Labels section to Technical Spec
- [x] Update generated override example with new labels
- [x] Add label documentation for external tool integration
