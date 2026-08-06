## Host Gateway Resolution

This spec defines how Scind implementations automatically normalize `host.docker.internal` so that it resolves to the developer's workstation across all platforms.

See [ADR-0014](../decisions/0014-host-docker-internal-normalization.md) for the decision rationale.

---

### The Workstation Host Concept

Docker provides no single abstraction for "the developer's actual workstation." `host.docker.internal` was designed for this purpose but only works automatically on Docker Desktop. `host-gateway` resolves to the Docker daemon's bridge gateway, which is not always the workstation (notably in WSL2 without Docker Desktop).

Scind defines the **workstation host** as: the machine where the developer's IDE, browser, and development tools run — regardless of whether Docker runs natively on that machine, inside a Linux VM managed by Docker Desktop, or inside a WSL2 distribution.

`host.docker.internal` is the DNS name containers use to reach the workstation host. This spec ensures that name resolves correctly in all supported configurations.

---

### Detection Algorithm

Implementations MUST determine the correct resolution value using the following precedence:

```
1. User override → use value directly
2. Docker Desktop detected → skip (already works)
3. WSL2 detected:
   a. Detect networking mode
   b. Mirrored mode + docker-ce inside WSL2 → "host-gateway"
   c. Mirrored mode + Docker Desktop → LAN IP from hostname -I
   d. NAT mode → gateway IP from default route
4. Native Linux → "host-gateway"
5. Fallback → "host-gateway"
```

**Step 1 — User override**: If the implementation's override variable (e.g., `XCIND_HOST_GATEWAY`) is set, use its value directly. No detection is performed.

**Step 2 — Docker Desktop**: If Docker Desktop is detected (see [Platform Detection Methods](#platform-detection-methods)), `host.docker.internal` already resolves correctly via Docker Desktop's internal DNS. The implementation SHOULD still generate the overlay (the entry is harmless and ensures consistency) but MAY skip generation as an optimization.

**Step 3 — WSL2**: If `/proc/version` contains `microsoft` (case-insensitive), the environment is WSL2. Determine the networking mode:

- **Mirrored mode**: Detected via `wslinfo --networking-mode` returning `mirrored`, or by the presence of the `loopback0` network interface (`ip link show loopback0`).
- **NAT mode**: The default; assumed when mirrored mode is not detected.

In mirrored mode, the resolution depends on the Docker runtime:
- **docker-ce inside WSL2** (no Docker Desktop): Use the literal string `host-gateway`. In mirrored mode, the WSL2 network stack shares the Windows host's interfaces, so `host-gateway` resolves to the correct destination.
- **Docker Desktop in WSL2**: Use the first LAN IP from `hostname -I`. Docker Desktop runs its own networking layer, so `host-gateway` may not resolve to the Windows host.

In NAT mode, extract the Windows host IP from the default route: `ip -4 route show default | awk '{print $3}'`. This is the NAT gateway, which is the Windows host.

**Step 4 — Native Linux**: Use the literal string `host-gateway`. Docker Engine resolves this to the `docker0` bridge gateway (typically `172.17.0.1`), which is the host machine.

**Step 5 — Fallback**: Use `host-gateway`. This is the safest default for unknown environments.

> **Note on `host-gateway`**: When the resolved value is the literal string `host-gateway`, Docker Engine substitutes it at container creation time with the daemon's configured gateway IP. This is not a DNS name — it is a Docker-specific token recognized only in `extra_hosts` entries.

---

### Generated Overlay Contract

The implementation generates a compose overlay file (e.g., `compose.host-gateway.yaml`) that adds `host.docker.internal` to `extra_hosts` for every service in the resolved compose configuration.

```yaml
services:
  web:
    extra_hosts:
      - "host.docker.internal:<resolved_value>"
  worker:
    extra_hosts:
      - "host.docker.internal:<resolved_value>"
  redis:
    extra_hosts:
      - "host.docker.internal:<resolved_value>"
```

**Inclusion rule**: Every service in the resolved compose config receives the entry, **unless** that service already defines `host.docker.internal` in its `extra_hosts` (in any of the user's compose files, including manual overrides). This preserves user-provided values.

**Resolved value**: Either a literal IP address (e.g., `192.168.1.100`) or the string `host-gateway`, depending on the detection algorithm's result.

**File location and merge order**: The overlay file is implementation-defined. It MUST be merged after the application's base compose file(s) and before any user-provided manual overrides, so that user overrides take final precedence.

---

### Environment Variable Exposure

The resolved workstation host value SHOULD also be exposed as an environment variable inside containers:

```yaml
services:
  web:
    environment:
      - SCIND_HOST_GATEWAY=<resolved_value>
```

This enables applications to reference the workstation host programmatically without relying on DNS resolution of `host.docker.internal`. A common use case is configuring Xdebug's `client_host`:

```ini
xdebug.client_host=${SCIND_HOST_GATEWAY}
```

The variable name is implementation-defined (e.g., `XCIND_HOST_GATEWAY` for Xcind) but SHOULD follow the `SCIND_` prefix convention when feasible.

---

### Opt-Out and Override

#### Per-Service Opt-Out

If a service already defines `host.docker.internal` in its `extra_hosts` — in the application's base compose file, the generated override, or a manual override — the implementation MUST NOT override it. The user's explicit configuration takes precedence.

#### Global Opt-Out

Setting the implementation's enable variable to `0` (e.g., `XCIND_HOST_GATEWAY_ENABLED=0`) disables host gateway normalization entirely. No overlay is generated. This is useful for users who manage `host.docker.internal` through other means (e.g., system-level DNS configuration).

#### Manual Override

Setting the implementation's override variable (e.g., `XCIND_HOST_GATEWAY=192.168.1.50`) forces the resolved value, bypassing all platform detection. This is the first check in the detection algorithm.

---

### Platform Detection Methods

#### Docker Desktop Detection

Check in order; stop at the first match:

| Method | What it detects | Reliability |
|---|---|---|
| `docker info` output contains `Docker Desktop` or `desktop-linux` context | Docker Desktop is the active runtime | High |
| `/mnt/wsl/docker-desktop` exists | Docker Desktop integration in WSL2 | High |
| Docker socket path points to Docker Desktop's socket | Desktop-managed daemon | Medium |

#### WSL2 Detection

| Method | What it detects | Reliability | Caveats |
|---|---|---|---|
| `grep -qi microsoft /proc/version` | WSL2 environment | High | Also matches WSL1 (rare in practice) |

#### WSL2 Networking Mode Detection

| Method | What it detects | Reliability | Caveats |
|---|---|---|---|
| `wslinfo --networking-mode` | Mirrored vs NAT mode | Authoritative | Not available on older WSL builds |
| `ip link show loopback0` | Mirrored mode (heuristic) | High | `loopback0` is mirrored-mode-only; fallback when `wslinfo` is unavailable |

#### Host IP Resolution

| Method | What it resolves | Reliability | Caveats |
|---|---|---|---|
| `ip -4 route show default \| awk '{print $3}'` | NAT gateway (= Windows host in WSL2 NAT mode) | High | Returns WSL2 VM gateway in mirrored mode — do not use in mirrored mode |
| `hostname -I \| awk '{print $1}'` | First LAN IP (shared with Windows in mirrored mode) | Medium | May return a VPN or virtual adapter IP first |

---

### Prior Art

**DDEV** is the most mature reference implementation of this pattern. Its `GetHostDockerInternalIP()` function (in Go) performs equivalent detection:

- Runs on every `ddev start`, making the behavior always-on.
- Handles native Linux, Docker Desktop (macOS/Windows), Colima, OrbStack, WSL2 with Docker Desktop, WSL2 with standalone Docker Engine, Gitpod, and Codespaces.
- For WSL2 NAT mode, extracts the Windows host IP from `ip -4 route show default`.
- For WSL2 mirrored mode (fixed in DDEV v1.24.5), uses `wslinfo --networking-mode` to detect, then selects `host-gateway` for docker-ce-inside-WSL2 or the LAN IP for Docker Desktop.
- Exposes `HOST_DOCKER_INTERNAL_IP` as an environment variable inside containers.

**Where Scind aligns**: The detection algorithm, the always-on default, the per-service opt-out via existing `extra_hosts`, and the environment variable exposure all follow DDEV's proven approach.

**Where Scind differs**: Scind uses a compose overlay file rather than direct Docker API manipulation, consistent with the [Pure Overlay Design](../decisions/0003-pure-overlay-design.md). Scind also provides a formal opt-out mechanism (`*_HOST_GATEWAY_ENABLED=0`) and a user override variable (`*_HOST_GATEWAY=<value>`), which DDEV does not expose as first-class configuration.

---

## Related Documents

- [ADR-0014: Automatic `host.docker.internal` Normalization](../decisions/0014-host-docker-internal-normalization.md) - Decision rationale
- [ADR-0003: Pure Overlay Design](../decisions/0003-pure-overlay-design.md) - Overlay pattern this follows
- [Generated Override Files](./generated-override-files.md) - Override file conventions
- [Environment Variables](./environment-variables.md) - Environment variable conventions
