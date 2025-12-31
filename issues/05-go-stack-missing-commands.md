# Issue Group 5: Go Stack Missing Command Scaffolds

**Documents Affected**: CLI Reference, Go Stack
**Suggested Order**: 5 of 6 (implementation scaffolding)
**Estimated Effort**: Medium

---

## Overview

The Go Stack specification provides scaffolding for most commands but is missing scaffolds for several commands defined in the CLI Reference. Additionally, one scaffolded command references a variable that doesn't exist.

---

## Issues

### M-2: Missing Port Commands Scaffold

**Severity**: Low

**Issue**: The CLI Reference defines six port commands (lines 779-898):
- `port list`, `port show`, `port release`, `port assign`, `port gc`, `port scan`

The Go Stack's CLI to Cobra mapping table (lines 1094-1099) references these, but no `internal/cli/port.go` scaffold file content is provided in the scaffolding instructions (Steps 4-9).

**Questions**: None needed

**Suggested Resolution**: Add a Step 6c with `internal/cli/port.go` scaffold containing all port commands with their flags.

**Response**:
> _[Your response here]_

---

### M-3: Missing Config Commands Scaffold

**Severity**: Low

**Issue**: The CLI Reference defines config commands (lines 1043-1122):
- `config show`, `config get`, `config set`, `config path`, `config edit`

The Go Stack mapping table (lines 1105-1107) references these, but no `internal/cli/config.go` scaffold content is provided.

**Questions**: None needed

**Suggested Resolution**: Add a Step 6d with `internal/cli/config.go` scaffold containing all config commands.

**Response**:
> _[Your response here]_

---

### M-4: Missing Utility Commands Scaffolds

**Severity**: Low

**Issue**: The CLI Reference defines utility commands (lines 1229-1331):
- `validate`, `doctor`, `open`, `urls`

The Go Stack project structure (lines 186-189) lists these files but no scaffold content is provided.

**Questions**: None needed

**Suggested Resolution**: Add scaffolds for `validate.go`, `doctor.go`, `open.go`, and `urls.go` in a new Step 10b.

**Response**:
> _[Your response here]_

---

### M-5: Missing init-shell Command Scaffold

**Severity**: Low

**Issue**: The CLI Reference documents `init-shell` (lines 1172-1196) and the Go Stack project structure lists `init_shell.go` (line 184), but no scaffold content is provided. The Shell Integration spec (lines 132-146) describes the installation process but the Go implementation details are missing.

**Questions**: None needed

**Suggested Resolution**: Add a scaffold for `internal/cli/init_shell.go` that uses `//go:embed` for the shell scripts.

**Response**:
> _[Your response here]_

---

### C-2: Alias Commands Reference Undefined RunE

**Severity**: Medium

**Issue**: In Go Stack's `internal/cli/aliases.go` scaffold (lines 750-783), the alias commands reference `workspaceUpCmd.RunE`, `workspaceDownCmd.RunE`, etc. However, looking at the `workspace.go` scaffold (lines 409-549), these commands define `RunE` as inline functions, not as named variables that can be shared.

When Cobra evaluates `upCmd.RunE = workspaceUpCmd.RunE`, it works, but the pattern shown won't properly inherit context detection since `PersistentPreRunE` runs on the command being executed, not the delegated command.

**Questions**:
1. Should aliases share the implementation via a common helper function instead of sharing RunE?
2. Or should aliases be implemented differently (e.g., as hidden commands that call the actual commands)?

**Suggested Resolution**: Refactor the alias pattern to call a shared implementation function, e.g.:
```go
var upCmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceUp(cmd, args) // shared function
    },
}
```

**Response**:
> _[Your response here]_

---

## Checklist

- [ ] Add port.go scaffold to Go Stack (Step 6c)
- [ ] Add config.go scaffold to Go Stack (Step 6d)
- [ ] Add utility command scaffolds (validate, doctor, open, urls)
- [ ] Add init_shell.go scaffold to Go Stack
- [ ] Fix alias command pattern to use shared implementation functions
