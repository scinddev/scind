# Issue Group 2: Schema & Validation

**Documents Affected**: Technical Spec + Go Stack  
**Suggested Order**: 2 of 10 (foundational—other implementation depends on correct schemas)  
**Estimated Effort**: Medium

---

## Overview

These issues concern the alignment between what the Technical Spec says about configuration schemas and how the Go Stack implements validation. Getting these right early prevents implementation headaches.

---

## Issues

### C-2: Go Struct Marks `port` as Required, Spec Says Optional

**Severity**: High

**Issue**: The Technical Spec (lines 436-439) states that `port:` is optional with inference rules:
> - If the Compose service has exactly one port in its `ports:` configuration, that port is used as the default
> - If the Compose service has multiple ports, `port:` must be explicitly specified

However, the Go Stack struct definition (lines 698-700) marks it as required:
```go
Port int `yaml:"port" validate:"required,min=1,max=65535"`
```

**Questions**:
1. Should the Go struct match the spec (make port optional)?
2. Where should port inference logic live (loader? generator?)?
3. Should inference failure produce a clear error message?

**Suggested Resolution**: 
- Change Go struct to `validate:"omitempty,min=1,max=65535"`
- Add inference logic in config loader
- Document the inference in both Tech Spec and Go Stack

**Response**:  
> Yes, if your changes to "omitempty" to the Go struct result in making the port optional (per the technical spec), go ahead with your plan.

---

### C-3: Go Struct Marks `service` as Required, Spec Says Optional

**Severity**: High

**Issue**: The Technical Spec (lines 405-417) says the exported service key defaults to the Compose service name when `service:` is omitted:

```yaml
exported_services:
  web:           # No service: field, defaults to Compose service "web"
    ports: [...]
  db:
    service: postgres   # Explicit: maps to Compose service "postgres"
```

The Go struct has `validate:"required"`:
```go
type ExportedService struct {
    Service string `yaml:"service" validate:"required"`
    // ...
}
```

**Questions**:
1. Should the Go struct match the spec (make service optional)?
2. Should the default be set during YAML unmarshaling or in a post-load step?

**Suggested Resolution**:
- Change to `validate:"omitempty"` 
- In loader: if `Service` is empty, set it to the map key

**Response**:  
> Your suggestions are sound; we should have the Go struct match the spec, the service should be optional, and should default to the map key if it is not specified.

---

### A-9: Traefik Router Name Collision Potential

**Severity**: Low

**Issue**: Generated Traefik labels use patterns like `traefik.http.routers.dev-app-one-web-https.rule=...`. With creative naming, collisions are theoretically possible. No collision detection is documented.

**Location**: Technical Spec lines 604-614

**Questions**:
1. Is this a real concern given naming conventions (lowercase alphanumeric + hyphens)?
2. Should router names include a hash suffix for safety?
3. Should generation validate uniqueness across the workspace?

**Possible Resolutions**:
- A) Document that naming conventions prevent collisions (do nothing)
- B) Add collision detection during generation
- C) Add a short hash suffix to router names

**Response**:  
> Document the naming conventions could cause collisions if people are not aware of the issue.

---

### A-10: Compose File Existence Not Validated

**Severity**: Medium

**Issue**: A flavor's `compose_files` list is validated for structure (`min=1`) but nothing validates that referenced files actually exist on disk.

**Location**: Go Stack lines 703-705

**Questions**:
1. When should validation occur?
   - At config load time (early failure)?
   - At `generate` time (with clear context)?
   - At `up` time (let docker compose fail with its own error)?
2. Should missing file errors include suggestions (e.g., "Did you mean docker-compose.yml?")?

**Suggested Resolution**: Validate at `generate` time with a clear error:
```
Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
  Application: app-two
  Available compose files: docker-compose.yaml, docker-compose.dev.yaml
```

**Response**:  
> Your suggestion is sound

---

## Checklist

- [x] Update `Port` struct field validation in Go Stack (changed to `omitempty`)
- [x] Document port inference logic in Tech Spec and Go Stack
- [x] Update `Service` struct field validation in Go Stack (changed to `omitempty`)
- [x] Document service name defaulting behavior (in Go Stack)
- [x] Document Traefik router collision warning in Tech Spec Naming Conventions
- [x] Document compose file existence validation timing (at `generate` time)
- [x] Update Tech Spec Generation Logic with validation steps

---

## Archived

This issue was archived on 2025-12-31 at 10:03:30.
