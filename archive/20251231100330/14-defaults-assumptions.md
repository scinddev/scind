# Issue Group 14: Defaults & Assumptions

**Documents Affected**: PRD + Technical Spec (+ CLI Reference for naming)
**Suggested Order**: 14 of 18 (clarifies fundamentals, low urgency)
**Estimated Effort**: Small

---

## Overview

Several configuration options lack explicit defaults. Documenting these prevents implementation ambiguity.

---

## Issues

### N-4: Default Visibility Unspecified

**Severity**: Low

**Issue**: The `visibility` field is documented as optional (`public` or `protected`) but no default is specified.

**Locations**:
- PRD line 101: "Each port can have a `visibility` of `public` or `protected`"
- Technical Spec line 122: Same description
- Go Stack line 809: `validate:"omitempty,oneof=public protected"`

**Questions**:
1. What is the default visibility when not specified?
2. Should the default be `protected` (more conservative) or `public`?

**Suggested Resolution**: Default to `protected` since it's more conservative—developers must explicitly opt into `public` visibility.

**Your Response**:
> Default to `protected`.

---

### N-5: Default Protocol for Proxied Ports Unspecified

**Severity**: Low

**Issue**: For `type: proxied` ports, `protocol` is described as optional in Go Stack (line 807):

```go
Protocol string `yaml:"protocol,omitempty" validate:"omitempty,oneof=http https tcp postgresql mysql"`
```

But what happens if a proxied port has no protocol specified? The Technical Spec doesn't address this case.

**Questions**:
1. Should `protocol` be required for proxied ports?
2. Or should it default to `https` (secure by default)?
3. Or should it default to `http` (simpler)?

**Suggested Resolution**: Make `protocol` required for `type: proxied` ports with a clear validation error if omitted.

**Your Response**:
> Make `protocol` required for proxied ports.

---

### N-6: Traefik Version Pinning

**Severity**: Low

**Issue**: The Technical Spec shows `image: traefik:v3.0` in examples (lines 263, 996) but doesn't specify if this is a pinned version or just an example.

**Questions**:
1. Should Contrail pin to a specific Traefik version (e.g., `v3.0.4`)?
2. Or use a floating tag like `v3` or `v3.0`?
3. Should this be configurable in `proxy.yaml`?

**Your Response**:
> Make it configurable in `proxy.yaml` but default to a pinned version.

---

### N-7: Single-App Workspace Naming Ambiguity

**Severity**: Medium

**Issue**: In a single-app workspace (where `path: .`), the relationship between workspace name and app name is unclear across examples.

**PRD line 86** shows:
```bash
contrail app init --app=myapp
```

**Technical Spec lines 206-213** shows:
```yaml
workspace:
  name: dev
  applications:
    my-project:      # <-- Different name from --app flag example
      path: .
```

**Questions**:
1. In a single-app workspace, how is the app name determined?
   - From the `--app` flag during `app init`?
   - From the directory name?
   - From the key in `workspace.yaml` applications map?
2. Are all three expected to match?
3. What happens if they don't match?

**Your Response**:
> The app name comes from the `--app` flag during `app init`. The workspace name can be different from the app name—if you create three separate single-application workspaces for the same app, you'd still want to distinguish between them by the workspace name even if the app name is the same. Since they don't all need to match, this may not be an issue that needs documenting beyond clarifying the `--app` flag behavior.

---

## Checklist

- [x] Document default visibility value — defaults to `protected` in PRD and Tech Spec
- [x] Clarify protocol requirement/default for proxied ports — now required for `type: proxied` in Tech Spec and Go Stack
- [x] Clarify Traefik version pinning strategy — configurable via `proxy.yaml` with pinned default `v3.2.3`
- [x] Clarify single-app workspace naming rules — added note to CLI Reference about independent naming

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
