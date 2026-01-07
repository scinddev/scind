// scaffold-utility.go
// Utility commands scaffold for Scind CLI
// Create as: internal/cli/doctor.go, internal/cli/open.go, internal/cli/urls.go

// --- doctor.go ---

package cli

import (
    "github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check system health and dependencies",
    Long:  `Verify that Docker, Docker Compose, and other dependencies are properly configured.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: check Docker, network, proxy status
        return nil
    },
}

func init() {
    rootCmd.AddCommand(doctorCmd)

    doctorCmd.Flags().Bool("fix", false, "attempt to fix issues automatically")
}

// --- open.go ---

var openCmd = &cobra.Command{
    Use:   "open [service]",
    Short: "Open service URL in browser",
    Long:  `Open the URL for a proxied service in the default browser.`,
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: resolve URL, open browser
        return nil
    },
}

func init() {
    rootCmd.AddCommand(openCmd)

    openCmd.Flags().Bool("print", false, "print URL instead of opening browser")
}

// --- urls.go ---

var urlsCmd = &cobra.Command{
    Use:   "urls",
    Short: "List all service URLs",
    Long:  `Display URLs for all proxied services in the current context.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation: list all proxied URLs
        return nil
    },
}

func init() {
    rootCmd.AddCommand(urlsCmd)
}
