# Automatic `host.docker.internal` Normalization

**Status**: Accepted

## Context

Docker's `host.docker.internal` DNS name is intended to resolve to the developer's workstation from inside a container. It is an official Docker convention, introduced in Docker Desktop 18.03 (March 2018), and Docker's documentation explicitly states it is "for development purpose and will not work in a production environment outside of Docker Desktop."

However, its availability varies dramatically across platforms:

| Platform | `host.docker.internal` works? | Resolves to |
|---|---|---|
| Docker Desktop (macOS) | Yes, automatic | Mac host via vpnkit proxy |
| Docker Desktop (Windows/WSL2) | Yes, automatic | Windows host via vpnkit proxy |
| Docker Desktop (Linux) | Yes, automatic | Linux host via Desktop's managed DNS |
| Docker Engine on native Linux | No | Must be manually configured |
| Docker Engine in WSL2 (no Desktop) | No, and `host-gateway` points to wrong host | WSL2 VM, not Windows host |

This inconsistency breaks development tools that need to connect back to the developer's workstation — debuggers (Xdebug), webhook receivers, and IDE integrations. The common advice to "add `host.docker.internal:host-gateway`" is written for native Linux users and actively misleads macOS and Windows developers: on Docker Desktop, `host.docker.internal` already works without it, and on WSL2 without Docker Desktop, `host-gateway` resolves to the WSL2 VM rather than the Windows host.

Docker Engine's `host-gateway` feature (Docker 20.10+, moby/moby PR #40007) substitutes the daemon's configured gateway IP — typically `172.17.0.1` (the `docker0` bridge). On Docker Desktop, both `host.docker.internal` and `host-gateway` reach the same destination because Desktop configures `--host-gateway-ip` internally. On standalone Docker Engine in WSL2, they diverge.

## Decision

Scind implementations SHOULD automatically ensure `host.docker.internal` resolves to the developer's workstation in all containers, for all services, unless the user has already configured it. This is achieved through a generated compose overlay file following the [Pure Overlay Design](0003-pure-overlay-design.md).

The normalization:

1. Detects the current platform and Docker runtime configuration.
2. Determines the correct value that will reach the developer's workstation (an IP address or the literal string `host-gateway`).
3. Generates a compose overlay that adds `host.docker.internal` to `extra_hosts` for every service that does not already define it.
4. Provides opt-out and override mechanisms for users with custom configurations.

## Consequences

### Positive

- Development tools (Xdebug, debuggers, webhook receivers) work out of the box across all platforms without per-project configuration.
- Users on Docker Desktop see no behavioral change — `host.docker.internal` already works; the generated overlay is additive and harmless.
- Users on native Linux get `host.docker.internal` without manual `extra_hosts` configuration.
- Users on WSL2 without Docker Desktop get the correct Windows host IP instead of the WSL2 VM IP.
- Follows the overlay pattern — application compose files remain workspace-agnostic.

### Negative

- Adds platform detection complexity to implementations.
- WSL2 detection heuristics may need updates as WSL evolves (e.g., new networking modes).
- Users unfamiliar with the feature may be surprised by the extra `extra_hosts` entry in their containers.

### Neutral

- The overlay does not modify containers where the user has already defined `host.docker.internal` in `extra_hosts`, preserving explicit configuration.
- The detection algorithm is inherently platform-specific and must be maintained as Docker and WSL2 evolve.

## Related Documents

- [Host Gateway Resolution Spec](../specs/host-gateway-resolution.md) - Detection algorithm, overlay contract, and opt-out mechanisms
- [ADR-0003: Pure Overlay Design](0003-pure-overlay-design.md) - This hook follows the overlay pattern
- [Generated Override Files](../specs/generated-override-files.md) - Override file generation conventions
