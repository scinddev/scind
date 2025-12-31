# Issue Group 18: Docker Label Consistency

**Documents Affected**: Technical Spec only
**Suggested Order**: 18 of 18 (minor fix, quick)
**Estimated Effort**: Small

---

## Overview

Minor inconsistency in Docker label prefix usage within the Technical Spec.

---

## Issues

### N-18: Docker Label Prefix Inconsistency

**Severity**: Low

**Issue**: The Technical Spec uses inconsistent Docker label prefixes in different sections.

**Docker Labels section (lines 836-911)** correctly uses `contrail.` prefix:
```yaml
labels:
  - "contrail.workspace.name=dev"
  - "contrail.workspace.path=/Users/beau/workspaces/dev"
  - "contrail.app.name=app-one"
```

**Operations section (line 1197)** uses unprefixed labels:
```bash
docker logs $(docker ps -q --filter "label=workspace.name=dev")
```

This should be:
```bash
docker logs $(docker ps -q --filter "label=contrail.workspace.name=dev")
```

**Fix**: Update Operations section examples to use `contrail.` prefix consistently.

**Your Response**:
> Fix the inconsistency — update Operations section examples to use `contrail.` prefix.

---

## Checklist

- [x] Fix Docker label prefix in Operations section examples — corrected to `contrail.workspace.name` and `contrail.app.name`
