# Issue Group 16: Security & Platform Concerns

**Documents Affected**: PRD + Technical Spec
**Suggested Order**: 16 of 18 (security merits attention, platform can defer)
**Estimated Effort**: Small

---

## Overview

These issues concern security posture and platform support scope.

---

## Issues

### N-13: Traefik Dashboard Security

**Severity**: Medium

**Issue**: The Traefik configuration exposes the dashboard on port 8080 without authentication:

**Technical Spec line 1007-1008**:
```yaml
ports:
  - "8080:8080"  # Dashboard
```

**CLI Reference line 997**:
```
Dashboard: http://localhost:8080
```

The Traefik dashboard can show routing rules, middleware, and service information. In shared development environments, this could expose information to other users on the same machine.

**Questions**:
1. Should the dashboard be disabled by default?
2. Should it require a password (add `--api.dashboard.basicAuth`)?
3. Should it bind to localhost only (`127.0.0.1:8080:8080`)?
4. Is this a non-issue since Contrail targets single-user development machines?

**Your Response**:
> Document future functionality: add ability to specify a password in `proxy.yaml` (possibly via environment variable). Also add `dashboard.enabled` config option that defaults to `true`, allowing users to disable if desired.

---

### N-14: Platform Support Scope

**Severity**: Low

**Issue**: All documentation assumes Unix-like systems:
- Shell scripts are bash/zsh/fish only (no PowerShell)
- Paths use forward slashes and `~/` home directory notation
- `~/.config/contrail/` assumes XDG-style config

Windows developers using Docker Desktop may expect Contrail to work.

**Questions**:
1. Is Windows support a goal for v1?
2. If not, should this be explicitly stated as a non-goal in the PRD?
3. If yes, what's the minimum scope (WSL only? Native Windows? PowerShell support?)

**Suggested Resolution**: Add to PRD Non-Goals section:
```markdown
7. **Windows native support**: Initial release targets macOS and Linux. Windows users should use WSL2.
```

**Your Response**:
> Add note to PRD Non-Goals section. (Already added during Group 15 implementation.)

---

## Checklist

- [x] Resolve Traefik dashboard security posture — added `dashboard.enabled` config with future password note
- [x] Document platform support scope in PRD — Windows non-goal already added

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
