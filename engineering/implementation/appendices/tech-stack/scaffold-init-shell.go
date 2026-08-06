// scaffold-init-shell.go
// init-shell command scaffold for Scind CLI
// Create as: internal/cli/init_shell.go

package cli

import (
    _ "embed"
    "fmt"

    "github.com/spf13/cobra"
)

//go:embed scripts/bash.sh
var bashScript string

//go:embed scripts/zsh.zsh
var zshScript string

//go:embed scripts/fish.fish
var fishScript string

var initShellCmd = &cobra.Command{
    Use:   "init-shell <shell>",
    Short: "Output shell integration script",
    Long: `Output shell integration script for the specified shell.

Supported shells: bash, zsh, fish

Add to your shell configuration:
  # Bash
  scind init-shell bash >> ~/.bashrc

  # Zsh
  scind init-shell zsh >> ~/.zshrc

  # Fish
  scind init-shell fish >> ~/.config/fish/conf.d/scind.fish`,
    Args:      cobra.ExactArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish"},
    RunE: func(cmd *cobra.Command, args []string) error {
        shell := args[0]
        switch shell {
        case "bash":
            fmt.Print(bashScript)
        case "zsh":
            fmt.Print(zshScript)
        case "fish":
            fmt.Print(fishScript)
        default:
            return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initShellCmd)
}
