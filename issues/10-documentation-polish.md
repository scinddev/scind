# Issue Group 10: Documentation Polish

**Documents Affected**: Various (see individual issues)  
**Suggested Order**: 10 of 10 (final cleanup after substantive changes)  
**Estimated Effort**: Small

---

## Overview

Final polish items—version synchronization and minor documentation gaps. Best tackled last after all substantive changes are complete.

---

## Issues

### C-1: Version Number Desynchronization

**Severity**: Low

**Issue**: The five documents have inconsistent version numbers:

| Document | Version |
|----------|---------|
| PRD | 0.5.0-draft |
| Technical Spec | 0.5.0-draft |
| CLI Reference | 0.2.0-draft |
| Shell Integration | 0.1.0-draft |
| Go Stack | 0.1.0-draft |

This makes it unclear which document is authoritative when they conflict, and suggests the newer docs (Shell Integration, Go Stack) may not have been reviewed against the updated PRD/Tech Spec.

**Options**:

**A) Synchronize All Versions**
- All documents share a single version number
- When any doc changes, all versions increment
- Pro: Clear that docs are a coordinated set
- Con: Version churn, hard to track which doc actually changed

**B) Independent Versions with Compatibility Matrix**
- Each doc versions independently
- Add compatibility note: "CLI Reference 0.2.0 is compatible with PRD 0.5.0"
- Pro: Granular tracking
- Con: Complex to maintain

**C) Primary + Satellite Model**
- PRD is the "primary" version (0.5.0)
- Other docs note which PRD version they implement
- "This CLI Reference implements PRD v0.5.0"
- Pro: Clear hierarchy
- Con: Still need to track alignment

**D) Date-Based Versioning**
- Use dates instead: "December 2024 Draft"
- All docs in sync by definition
- Pro: Simple, always current
- Con: Loses semantic versioning benefits

**Suggested Resolution**: Option A (synchronized versions) for draft phase. Once stable, can revisit.

After completing all issue groups, update all documents to the same version (suggest 0.6.0-draft to indicate significant revision).

**Response**:  
> _[Your response here]_

---

### X-2: `--since` Flag for Logs Not Demonstrated

**Severity**: Low

**Issue**: CLI Reference documents `--since` for `workspace logs` (line 364) but no examples demonstrate its use.

**Current documentation**:
```
| `--since` | Show logs since timestamp or duration (e.g., "10m", "2024-01-01") |
```

**Action**: Add example to CLI Reference in the logs section or examples section:

```bash
# View logs from the last 30 minutes
contrail logs --since=30m

# View logs since a specific time
contrail logs --since="2024-12-28T10:00:00"

# Combine with follow and app filter
contrail logs -a app-one --since=1h -f
```

**Response**:  
> _[Your response here]_

---

## Post-Completion Checklist

After all other issue groups are resolved:

- [ ] Update all document versions to synchronized number
- [ ] Update revision history in each document
- [ ] Add `--since` examples to CLI Reference
- [ ] Do a final cross-reference check for any new inconsistencies
- [ ] Update the "Related Documentation" links if any docs were renamed/restructured

---

## Version Update Template

When updating versions, use this template for revision history:

```markdown
| 0.6.0-draft | Dec 2024 | Specification review: [brief description of changes] |
```

Example entries:
- PRD: "Specification review: clarified network creation timing, added workspace discovery"
- Tech Spec: "Specification review: added destroy sequence, staleness detection, port scanning"
- CLI Reference: "Specification review: aligned proxy commands, added --since examples"
- Shell Integration: "Specification review: removed -f flag collision, added exit codes"
- Go Stack: "Specification review: fixed validation rules, added missing commands"
