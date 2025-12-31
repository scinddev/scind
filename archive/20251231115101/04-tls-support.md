# Issue Group 4: TLS/HTTPS Configuration Gap

**Documents Affected**: PRD, Technical Spec
**Suggested Order**: 4 of 6 (future feature clarification)
**Estimated Effort**: Medium

---

## Overview

The specs describe HTTPS routing through Traefik but lack clarity on the current TLS implementation status and how TLS termination actually works.

---

## Issues

### A-2: Unclear TLS Certificate Handling for HTTPS Proxied Services

**Severity**: Medium

**Issue**: The PRD (lines 90-95) and Technical Spec (lines 100-118) describe `proxied` type with `protocol: https` that routes through Traefik's `websecure` entrypoint with "TLS termination." However:

1. The Technical Spec's Traefik static config (lines 294-309) doesn't show any TLS certificate configuration
2. The proxy directory structure (lines 252-258) mentions a `certs/` directory but with "(future)" annotation
3. The Technical Spec's Future Enhancements section (lines 1297-1308) lists "HTTPS Support" as a future feature with `mkcert` integration

This creates ambiguity: Can HTTPS work today? Does Traefik auto-generate self-signed certs? Is this a blocking gap?

**Questions**:
1. Should the current implementation support HTTPS without explicit certificate configuration (using self-signed or auto-generated certs)?
2. Or should the specs clarify that HTTPS requires manual certificate setup until the mkcert integration is implemented?

**Suggested Resolution**: Add a note in the Technical Spec's Proxy Layer section clarifying the current TLS status—either documenting that Traefik will use its default self-signed certificate, or noting that HTTPS requires manual cert configuration until the mkcert feature is implemented.

**Response**:
> Implemented flexible TLS configuration with three modes:
> - `auto`: Uses mkcert if available, falls back to Traefik's self-signed certificate
> - `custom`: User-provided cert/key files (supports enterprise CA certificates)
> - `disabled`: HTTP only
>
> Added TLS schema to proxy.yaml, updated proxy directory structure, and added decision record to PRD. Enterprise CA certificates are fully supported via `custom` mode by providing cert/key paths.

---

## Checklist

- [x] Clarify current TLS/HTTPS status in Technical Spec Proxy Layer section
- [x] Add TLS configuration schema to proxy.yaml in Technical Spec
- [x] Update proxy directory structure to include certs/ and dynamic/tls.yaml
- [x] Move "HTTPS Support" from Future Enhancements to implemented (now "Per-Workspace TLS Overrides" as future)
- [x] Add TLS decision record to PRD Architecture Decisions section

---

## Archived

This issue was archived on 2024-12-31 at 11:51:01.
