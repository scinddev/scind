# Host-View Environment Symmetry

**Status**: Accepted

## Context

Scind injects service-discovery environment variables (`SCIND_{APP}_{EXPORT}_*`) into application containers so services can find one another without hardcoded hostnames. These variables are inherently **container-flavored**: an assigned export is described by its in-network alias plus the container port (`_HOST` + `_PORT`), and a proxied export is described by its routable proxy hostname and URL.

Container-flavored values are correct for container-to-container traffic, but the development loop is not confined to containers. Developers routinely run **host** processes that must reach the very same services:

- language runtimes launched via mise/direnv-loaded shells
- test suites and REPLs run directly on the host
- language servers, type checkers, and linters
- debuggers attaching to a host-run process
- ad-hoc scripts and database clients

From the host, these processes cannot use a container's view. An assigned export does not live at `{alias}:{container_port}` on the host; the alias is not resolvable outside the Compose network, and the port that is actually reachable is the **host-published** port, which frequently differs from the container port after a conflict reassignment. On the host the same export lives at `127.0.0.1:<host_port>`. This is precisely the distinction captured by the [three-value assigned-export discovery contract](../specs/environment-variables.md#assigned-export-discovery-contract-three-values): container view is `_HOST` + `_PORT`, host view is `127.0.0.1` + `_HOST_PORT`.

The practical failure is that **one committed connection string cannot resolve correctly in both places.** A value such as:

```
DATABASE_URL=postgres://app@${SCIND_SHARED_DB_DB_HOST}:${SCIND_SHARED_DB_DB_PORT}/app
```

is right inside a container (`shared-db-db:5432`) and wrong on the host (`shared-db-db` does not resolve, and `5432` may not be the published port). Asking developers to hand-maintain a parallel host-flavored file invites drift the moment a port is reassigned.

Proxied and apex exports do not have this problem: their hostnames route identically from host and container, so only assigned exports need a host-flavored rendering.

## Decision

Scind offers an **opt-in host-view environment file** that renders the **same discovery computation** as container injection, but with **host-flavored values**:

- **Assigned exports** resolve to `127.0.0.1` plus the allocated host-published port (`_HOST_PORT`).
- **Proxied and apex exports** keep their routable hostname, scheme, and URL unchanged.

The host-view file is produced from the same resolved topology (the manifest) that drives container injection — it is a second rendering of one computation, not a separately maintained artifact. Because both renderings share a source, the container view and the host view **cannot drift**.

The design leans on two language-agnostic layering guarantees so the host-view file composes safely:

1. **Dotenv loaders never override a real OS environment variable.** Mainstream dotenv implementations (mise, direnv, and each language's dotenv library) treat an already-set OS variable as authoritative. A developer who exports a value intentionally keeps it; the host-view file only fills what is unset.
2. **Compose `environment:` beats `env_file:`.** A variable injected by the generated override via `environment:` (the container-flavored value) wins over any `env_file:`-supplied value with the same key. A host-view file sourced on the host and a container-flavored value injected into a container therefore never collide.

Enabling the file is **explicit and never inferred from a filename.** The user selects a write mode:

- **own** — Scind owns the whole file and regenerates it wholesale.
- **block** — Scind manages only a delimited block inside a developer-owned file, leaving surrounding content untouched.

Requiring an explicit mode prevents Scind from clobbering a developer's own dotenv file simply because it carries a conventional name such as `.env`.

This capability **depends on** the three-value assigned-export contract (`_HOST`, `_PORT`, `_HOST_PORT`): without a distinct host-published port value, host-flavored rendering has nothing correct to emit.

### Scoping Note

Host-side process access sits at the edge of Scind's otherwise container-centric runtime model. It is recorded here as an accepted design because host processes are a first-class part of the development loop Scind targets, and dropping the capability would leave the container/host asymmetry unaddressed. Whether a given Scind runtime enables the host-view file by default (versus purely on request) is a separate product decision; this ADR establishes the capability and its contract, not a mandatory default.

## Consequences

### Positive

- Host processes (test runners, language servers, debuggers, mise/direnv shells) reach the same services as containers using one coherent set of values.
- A single committed connection string can be made to resolve in both host and container contexts without hand-maintaining parallel files.
- Shared computation makes container and host views drift-proof by construction.
- Opt-in, explicit write modes prevent Scind from overwriting developer-owned dotenv files.

### Negative

- Adds a second rendering path and a write-mode concept to the configuration surface.
- Relies on dotenv-loader and Compose layering semantics that, while widely honored, are conventions a badly-behaved loader could violate.
- Broadens Scind's footprint from container orchestration toward host-process environment management.

### Neutral

- Requires the three-value assigned-export contract as a prerequisite; the two changes ship together.
- Proxied and apex exports render identically in both views, so the host view differs from the container view only for assigned exports.

## Related Documents

- [Environment Variables Spec — Host-View Environment](../specs/environment-variables.md#host-view-environment) — Behavior, layering guarantees, and write modes
- [Environment Variables Spec — Assigned-Export Discovery Contract](../specs/environment-variables.md#assigned-export-discovery-contract-three-values) — The three-value contract this depends on
- [ADR-0005: Structure vs State Separation](0005-structure-vs-state-separation.md) — Config/state boundary the host-view file respects
- [ADR-0014: host.docker.internal Normalization](0014-host-docker-internal-normalization.md) — Related host-reachability concern from the container side

<!-- Validated by the Xcind proof-of-concept (Xcind ADR-0020), which implements the opt-in host-view env file and exercises both layering guarantees end-to-end. -->
