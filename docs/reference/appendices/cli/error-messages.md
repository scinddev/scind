# CLI Error Messages

Complete catalog of Scind CLI error messages and their meanings.

<!-- Source: specs/scind-cli-reference.md -->

---

## Docker Errors (Exit Code 4)

### Docker Not Available

Commands that require Docker check for availability upfront:

```bash
$ scind up
Error: Docker is not installed or not running.
Run 'scind doctor' for setup guidance.
```

---

## Context Detection Errors (Exit Code 5)

### No Workspace Found

**Neither workspace nor application found**:

```bash
$ cd ~
$ scind app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: scind app status --workspace=NAME --app=NAME

Available workspaces: scind workspace list
```

### Application Found Without Workspace

Helps identify misplaced applications not in a workspace:

```bash
$ cd ~/random-project
$ scind app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: scind workspace init --workspace=NAME
```

### No Application Context

For app-specific commands when you're in the workspace root:

```bash
$ cd ~/workspaces/dev
$ scind app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: scind app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

---

## Configuration Errors (Exit Code 3)

### Workspace Name Already Registered

```bash
$ scind workspace init --workspace=dev
Error: Workspace "dev" already registered at ~/workspaces/dev
Use a different name, or run `scind workspace prune` if that path no longer exists
```

### Invalid Configuration

When configuration files have syntax or schema errors:

```bash
$ scind validate
Error: Invalid configuration in ./application.yaml

  Line 12: unknown field "export_services" (did you mean "exported_services"?)
  Line 15: missing required field "type" in ports configuration
```

### Flavor References Non-Existent File

```bash
$ scind up
Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
  Application: app-two
  Available compose files: docker-compose.yaml, docker-compose.dev.yaml
```

### Exported Service References Non-Existent Compose Service

```bash
$ scind up
Error: Exported service "api" references non-existent Compose service: backend
  Application: my-app
  Available services in docker-compose.yaml: web, db, redis
```

---

## Port Errors (Exit Code 4)

### Port Conflict at Startup

```bash
$ scind up
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  scind port scan       # Check which ports are conflicting
  scind port release 5432   # Release the conflicting assignment
  scind generate --force    # Regenerate with new port assignment
```

---

## Proxy Errors

### Proxy Configuration Already Exists

```bash
$ scind proxy init
Error: Proxy configuration already exists at ~/.config/scind/proxy/
Use --force to overwrite, or --path to create elsewhere.
```

### Network Conflict Warning

Not an error, but a warning when the proxy network may not have been created by Scind:

```
Warning: Network 'scind-proxy' exists but may not have been created by Scind.
  Driver: bridge (expected: bridge) ✓
  Labels: scind.managed not found ⚠

Use 'scind proxy up --recreate' to recreate the network.
```

---

## DNS Warnings

### Wildcard DNS Not Configured

```
⚠ Wildcard DNS not configured. Individual hostnames may not resolve.
  Configure dnsmasq: address=/scind.test/127.0.0.1
```

This warning appears when `scind doctor` detects that the base domain resolves but subdomains do not, indicating wildcard DNS is not properly configured.

---

## compose-prefix Errors (Exit Code 5)

### Context Cannot Be Resolved

```bash
$ cd ~
$ scind compose-prefix
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly with -w and -a flags

Available workspaces: scind workspace list
```

---

## Exit Codes Reference

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments or flags |
| 3 | Configuration error |
| 4 | Docker/Compose error |
| 5 | Context detection failed (workspace/app not found) |

---

## Related Documents

- [CLI Reference](../../cli.md)
- [Detailed Examples](./detailed-examples.md)
