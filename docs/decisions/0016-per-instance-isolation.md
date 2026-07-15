# Per-Instance Isolation for Multiple Working Copies

**Status**: Accepted

## Context

[ADR-0001](./0001-docker-compose-project-name-isolation.md) isolates workspaces
and applications by deriving a Docker Compose project name of the form
`{workspace}-{application}` and a workspace-internal network named
`{workspace}-internal`. This works because the derivation is deterministic: the
same names always resolve to the same Docker resources.

That determinism carries a hidden assumption: **one working copy per project.**
The derived names are a pure function of the workspace and application
identities, with no term that varies between two simultaneously-checked-out
copies of the same source.

The assumption breaks with linked git worktrees. When a developer checks the
same repository out into multiple worktrees (for example, to run a feature
branch and `main` side by side), every worktree resolves to the **same**
workspace name, therefore the **same** compose project name and the **same**
`{workspace}-internal` network. Instead of two isolated stacks, the second `up`
attaches to, restarts, or collides with the first. The mechanism that exists to
provide isolation actively removes it in this case.

This is a defect in the design assumption, not in any particular
implementation. Any implementation that derives resource names purely from
workspace/application identity — the Go implementation included — inherits it.
The Xcind Bash proof-of-concept encountered and validated the fix
(see Xcind ADR-0019).

## Decision

Introduce a **per-instance disambiguation token** that participates in the
derivation of the compose project name and the workspace-internal network name.

The token **extends** the existing names; it does not replace them:

- When the token is **empty**, names are exactly as ADR-0001 specifies today
  (`{workspace}-{application}`, `{workspace}-internal`). This is the default,
  and it is fully backward compatible — existing workspaces keep their current
  resource names, and there is no migration.
- When the token is **non-empty**, it is folded into the derived names (for
  example `{workspace}-{instance}-internal` for the internal network, and the
  corresponding extension of the project name). Two worktrees with distinct
  tokens therefore derive distinct project names and distinct networks, and are
  isolated as intended.

**Resolution order** for the token (first match wins):

1. **Explicit value** — a value configured/passed by the user.
2. **Opt-out** — an explicit request for the empty token (preserve legacy
   names) even when auto-detection would otherwise produce one.
3. **Auto-detect** — derive a stable token from a linked git worktree when one
   is detected.

Auto-detection keys on a **stable worktree-directory identity**, not on the
branch. Branches move (checkout, rebase, rename) while the resource identity
must stay pinned to the working copy for the life of that copy; a branch-based
token would rename networks and orphan containers on every branch switch.

## Consequences

### Positive

- Multiple linked worktrees of the same repository run as genuinely isolated
  stacks, delivering the isolation ADR-0001 promised.
- Fully backward compatible: the empty-token default reproduces today's names
  byte-for-byte, so existing workspaces need no migration.
- The fix lives in name derivation, a single well-defined seam, rather than
  being scattered across lifecycle operations.

### Negative

- Name derivation gains a conditional branch (token present vs. absent),
  slightly complicating the naming rules and the collision analysis.
- Auto-detection introduces an environment-sensitive input (worktree layout),
  so identical config in different checkouts can now legitimately produce
  different resource names.

### Neutral

- The token is a general disambiguation seam; multiple-worktree isolation is its
  first consumer, but any future need for parallel same-source instances can
  reuse it.

## Related Documents

- [ADR-0001: Docker Compose Project Name Isolation](./0001-docker-compose-project-name-isolation.md) - The base isolation mechanism this ADR extends
- [ADR-0004: Convention-Based Naming](./0004-convention-based-naming.md) - Naming derivation this token participates in
- [Naming Conventions](../specs/naming-conventions.md) - Where the instance token folds into the project-name and internal-network patterns
- Xcind ADR-0019 (proof-of-concept) - Validating implementation of per-instance isolation
