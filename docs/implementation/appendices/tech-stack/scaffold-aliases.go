// scaffold-aliases.go
// Top-level aliases scaffold for Contrail CLI
// Create as: internal/cli/aliases.go

package cli

import (
    "github.com/spf13/cobra"
)

// Top-level aliases for common operations
// These call shared implementation functions to ensure proper context handling
var upCmd = &cobra.Command{
    Use:   "up",
    Short: "Alias for 'workspace up'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceUp(cmd, args)
    },
}

var downCmd = &cobra.Command{
    Use:   "down",
    Short: "Alias for 'workspace down'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceDown(cmd, args)
    },
}

var psCmd = &cobra.Command{
    Use:   "ps",
    Short: "Alias for 'workspace status'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceStatus(cmd, args)
    },
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Alias for 'workspace generate'",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runWorkspaceGenerate(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(upCmd)
    rootCmd.AddCommand(downCmd)
    rootCmd.AddCommand(psCmd)
    rootCmd.AddCommand(generateCmd)

    // Copy flags from workspace commands to aliases
    upCmd.Flags().AddFlagSet(workspaceUpCmd.Flags())
    downCmd.Flags().AddFlagSet(workspaceDownCmd.Flags())
    generateCmd.Flags().AddFlagSet(workspaceGenerateCmd.Flags())
}

// Shared implementation functions (called by both workspace commands and aliases)
// These are defined here but the actual implementation lives in workspace.go

func runWorkspaceUp(cmd *cobra.Command, args []string) error {
    // Implementation: generate overrides, start containers
    return nil
}

func runWorkspaceDown(cmd *cobra.Command, args []string) error {
    // Implementation: stop containers, optionally remove volumes
    return nil
}

func runWorkspaceStatus(cmd *cobra.Command, args []string) error {
    // Implementation: aggregate status from all apps
    return nil
}

func runWorkspaceGenerate(cmd *cobra.Command, args []string) error {
    // Implementation: generate override files
    return nil
}
