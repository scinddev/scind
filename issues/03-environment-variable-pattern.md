# Issue Group 3: Environment Variable Pattern Inconsistency

**Documents Affected**: PRD, Technical Spec
**Suggested Order**: 3 of 6 (affects generated output)
**Estimated Effort**: Small

---

## Overview

The environment variable naming pattern for service discovery has subtle inconsistencies between documents regarding underscores vs hyphens in generated names.

---

## Issues

### A-1: Ambiguity in Environment Variable Name Generation for Hyphenated Names

**Severity**: Medium

**Issue**: Both PRD (line 415-416) and Technical Spec (lines 930-935) define the environment variable pattern as:
- `CONTRAIL_{APPLICATION}_{EXPORTED_SERVICE}_{SUFFIX}`

However, application names and exported service names can contain hyphens (e.g., `app-one`, `web-debug`), but environment variable names conventionally use underscores.

The examples in PRD (lines 429-443) show `CONTRAIL_APP_ONE_WEB_HOST` for application `app-one` and export `web`, implying hyphens are converted to underscores. But this transformation rule is not explicitly documented.

**Questions**:
1. Should the specs explicitly state that hyphens in application and exported service names are converted to underscores in environment variable names?

**Suggested Resolution**: Add explicit documentation in both PRD and Technical Spec stating: "Hyphens in application and exported service names are converted to underscores in environment variable names."

**Response**:
> Approved. Added explicit name transformation rule to both documents.

---

## Checklist

- [x] Add hyphen-to-underscore conversion rule to PRD Service Discovery section
- [x] Add hyphen-to-underscore conversion rule to Technical Spec Environment Variable Injection section
