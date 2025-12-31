# Issue Group 24: Traefik Dashboard Port Configuration

**Documents**: Technical Spec, CLI Reference
**Effort**: Small

---

## Issues

### Issue N-24: Traefik Dashboard Port Binding in Proxy Configuration

**Severity**: Medium
**Category**: Security/Configuration

**Finding**: The Technical Spec's proxy docker-compose.yaml example (lines 1009-1013) shows:
```yaml
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Dashboard
```

And the CLI Reference's `proxy status` output (lines 1016-1024) shows:
```
Dashboard: http://localhost:8080
```

However, the Technical Spec's "Proxy Infrastructure" generated docker-compose.yaml (lines 260-287) does NOT include port 8080:
```yaml
    ports:
      - "80:80"
      - "443:443"
```

These two sections of the Technical Spec are inconsistent about whether port 8080 is exposed.

**Additional concern**: Exposing the Traefik dashboard on port 8080 without authentication is a security concern. The Technical Spec (line 292) shows:
```yaml
api:
  dashboard: false
```

This suggests the dashboard is disabled by default in the generated traefik.yaml, but then the later section shows port 8080 exposed and `--api.dashboard=true` in command args.

**Documents**:
- `contrail-technical-spec.md` (lines 260-287) - no port 8080
- `contrail-technical-spec.md` (lines 1009-1013) - has port 8080
- `contrail-technical-spec.md` (line 292) - `dashboard: false`
- `contrail-technical-spec.md` (line 1002-1004) - `--api.dashboard=true`

**Recommendation**:
1. Decide if the dashboard should be enabled by default or opt-in
2. Make both sections consistent
3. If enabled, consider security implications (addressed in Issue 16)

---

## Tasks

### Task 24.1: Align Dashboard Configuration

Decide on the default dashboard behavior and update both proxy configuration sections in the Technical Spec to be consistent.

**Your Decision**:

> Add `dashboard` configuration to `proxy.yaml` with `enabled` (default: true) and `port` (default: 8080) settings. When `dashboard.enabled: true`, expose the port and use `--api.dashboard=true`. When `dashboard.enabled: false`, do NOT expose the port and use `--api.dashboard=false`. Update `proxy status` to show `Dashboard: http://localhost:{port}` when enabled or `Dashboard: disabled` when not.

---

## Summary

| Issue | Severity | Decision |
|-------|----------|----------|
| N-24: Dashboard port configuration inconsistency | Medium | ✅ Fixed |

