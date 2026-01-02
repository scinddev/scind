# CLI Error Messages

Complete catalog of Contrail CLI error messages and their meanings.

<!-- Source: specs/contrail-cli-reference.md -->

---

## Context Detection Errors

### No Workspace Found (Exit Code 5)

**No workspace found, but application.yaml exists**:

Helps identify misplaced applications not in a workspace.

```bash
$ cd ~/random-project
$ contrail app status
Error: No workspace found (workspace.yaml) in current directory or any parent directories,
but found an application (application.yaml) at: /home/user/random-project/application.yaml

Create a workspace with: contrail workspace init --workspace=NAME
```

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

**Workspace found but no application context**:

For app-specific commands when you're in the workspace root.

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

## Docker Errors (Exit Code 4)

### Docker Not Available

```bash
$ contrail up
Error: Docker is not installed or not running.
Run 'contrail doctor' for setup guidance.
```

---

## Configuration Errors (Exit Code 3)

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

## Port Errors

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

## Workspace Errors

### Workspace Name Already Registered

```bash
$ contrail workspace init --workspace=dev
Error: Workspace "dev" already registered at ~/workspaces/dev
Use a different name, or run `contrail workspace prune` if that path no longer exists
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

---

## compose-prefix Errors

### Context Cannot Be Resolved (Exit Code 5)

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
