# CLI Error Messages

Complete catalog of Contrail CLI error messages and their meanings.

<!-- Source: specs/contrail-cli-reference.md -->

---

## Docker Errors (Exit Code 4)

### Docker Not Available

Commands that require Docker check for availability upfront:

```bash
$ contrail up
Error: Docker is not installed or not running.
Run 'contrail doctor' for setup guidance.
```

---

## Context Detection Errors (Exit Code 5)

### No Workspace Found

**Neither workspace nor application found**:

```bash
$ cd ~
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
and no application (application.yaml) found either.

Either:
  1. Run from within a workspace directory
  2. Specify explicitly: contrail app status --workspace=NAME --app=NAME

Available workspaces: contrail workspace list
```

### Application Found Without Workspace

Helps identify misplaced applications not in a workspace:

```bash
$ cd ~/random-project
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

### No Application Context

For app-specific commands when you're in the workspace root:

```bash
$ cd ~/workspaces/dev
$ contrail app status
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly: contrail app status --app=NAME

Available apps in 'dev': app-one, app-two, app-three
```

---

## Configuration Errors (Exit Code 3)

### Workspace Name Already Registered

```bash
$ contrail workspace init --workspace=dev
Error: Workspace "dev" already registered at ~/workspaces/dev
Use a different name, or run `contrail workspace prune` if that path no longer exists
```

### Invalid Configuration

When configuration files have syntax or schema errors:

```bash
$ contrail validate
Error: Invalid configuration in ./application.yaml

  Line 12: unknown field "export_services" (did you mean "exported_services"?)
  Line 15: missing required field "type" in ports configuration
```

### Flavor References Non-Existent File

```bash
$ contrail up
Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
  Application: app-two
  Available compose files: docker-compose.yaml, docker-compose.dev.yaml
```

### Exported Service References Non-Existent Compose Service

```bash
$ contrail up
Error: Exported service "api" references non-existent Compose service: backend
  Application: my-app
  Available services in docker-compose.yaml: web, db, redis
```

---

## Port Errors (Exit Code 4)

### Port Conflict at Startup

```bash
$ contrail up
Error: Port conflict detected for app-one

Port 5432 is assigned to app-one/postgres but is no longer available.
Another process may be using this port.

To resolve:
  contrail port scan       # Check which ports are conflicting
  contrail port release 5432   # Release the conflicting assignment
  contrail generate --force    # Regenerate with new port assignment
```

---

## Proxy Errors

### Proxy Configuration Already Exists

```bash
$ contrail proxy init
Error: Proxy configuration already exists at ~/.config/contrail/proxy/
Use --force to overwrite, or --path to create elsewhere.
```

### Network Conflict Warning

Not an error, but a warning when the proxy network may not have been created by Contrail:

```
Warning: Network 'contrail-proxy' exists but may not have been created by Contrail.
  Driver: bridge (expected: bridge) ✓
  Labels: contrail.managed not found ⚠

Use 'contrail proxy up --recreate' to recreate the network.
```

---

## DNS Warnings

### Wildcard DNS Not Configured

```
⚠ Wildcard DNS not configured. Individual hostnames may not resolve.
  Configure dnsmasq: address=/contrail.test/127.0.0.1
```

This warning appears when `contrail doctor` detects that the base domain resolves but subdomains do not, indicating wildcard DNS is not properly configured.

---

## compose-prefix Errors (Exit Code 5)

### Context Cannot Be Resolved

```bash
$ cd ~
$ contrail compose-prefix
Error: No application context detected.

Either:
  1. Run from within an application directory
  2. Specify explicitly with -w and -a flags

Available workspaces: contrail workspace list
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
