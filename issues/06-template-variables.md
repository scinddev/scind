# Issue Group 6: Template Variables Documentation

**Documents Affected**: PRD, Technical Spec
**Suggested Order**: 6 of 6 (documentation enhancement)
**Estimated Effort**: Small

---

## Overview

Template variables are documented in the Technical Spec but could benefit from additional clarity about default values and the relationship between PRD's simplified examples and Technical Spec's detailed variables.

---

## Issues

### A-3: Template Variable Syntax Mismatch Between PRD and Technical Spec

**Severity**: Low

**Issue**: The PRD's Appendix on Template Customization (lines 599-606) shows:
```
- Hostname: `{workspace}-{app}-{export}.{domain}`
- Alias: `{app}-{export}`
```

While the Technical Spec (lines 450-461) shows:
```yaml
hostname: "%WORKSPACE_NAME%-%APPLICATION_NAME%-%EXPORTED_SERVICE%.%PROXY_DOMAIN%"
alias: "%APPLICATION_NAME%-%EXPORTED_SERVICE%"
```

The PRD uses `{placeholder}` syntax while Technical Spec uses `%PLACEHOLDER%` syntax. They also use different terminology (`app` vs `APPLICATION_NAME`, `export` vs `EXPORTED_SERVICE`).

**Questions**:
1. Should the PRD's simplified examples be updated to match the Technical Spec's actual variable names and syntax?
2. Or should there be an explicit note that the PRD shows conceptual patterns while Technical Spec shows actual implementation?

**Suggested Resolution**: Add a note to the PRD's Template Customization appendix clarifying that the simplified `{placeholder}` syntax is conceptual, and refer readers to the Technical Spec for the actual `%VARIABLE%` syntax.

**Response**:
> _[Your response here]_

---

### A-4: Missing SERVICE_PORT in Template Examples

**Severity**: Low

**Issue**: The Technical Spec's Template Variables table (lines 464-477) lists `%SERVICE_PORT%` (the container port number), but there's no example showing when this variable would be used in a template. The default templates don't use it.

**Questions**:
1. Is `%SERVICE_PORT%` intended for custom user templates? If so, what's a use case?
2. Should it be removed if there's no practical use case?

**Suggested Resolution**: Either add an example use case for `%SERVICE_PORT%` (e.g., for debugging labels) or add a note explaining it's available for advanced customization.

**Response**:
> _[Your response here]_

---

## Checklist

- [ ] Clarify template syntax relationship between PRD and Technical Spec
- [ ] Document `%SERVICE_PORT%` use case or rationale
