// scaffold-proxy.go
// Proxy commands scaffold for Scind CLI
// Create as: internal/cli/proxy.go

package cli

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
    Use:   "proxy",
    Short: "Manage the Traefik reverse proxy",
}

var proxyInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Bootstrap proxy configuration",
    Long:  `Creates the Traefik Docker Compose project at ~/.config/scind/proxy/`,
    RunE: func(cmd *cobra.Command, args []string) error {
        force, _ := cmd.Flags().GetBool("force")
        domain, _ := cmd.Flags().GetString("domain")
        path, _ := cmd.Flags().GetString("path")

        // Check if proxy config exists
        if _, err := os.Stat(filepath.Join(path, "docker-compose.yaml")); err == nil {
            if !force {
                return fmt.Errorf("proxy configuration already exists at %s\nUse --force to overwrite", path)
            }
            // Backup existing config
            // ...
        }

        // Create directory structure
        // Create docker-compose.yaml, traefik.yaml, dynamic/, certs/
        // Create scind-proxy network if needed
        // Output next steps

        _ = domain // TODO: use domain in generated config
        return nil
    },
}

var proxyUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Start the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        recreate, _ := cmd.Flags().GetBool("recreate")

        // Run proxy init if config doesn't exist
        // Create scind-proxy network if needed
        // If recreate flag is set, remove and recreate the network
        // Validate existing network configuration
        // Start containers via docker compose
        _ = recreate // TODO: implement
        return nil
    },
}

var proxyDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Stop the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Stop containers via docker compose
        return nil
    },
}

var proxyRestartCmd = &cobra.Command{
    Use:   "restart",
    Short: "Restart the Traefik proxy",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Restart via docker compose
        return nil
    },
}

var proxyStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show proxy status",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Check container status, network, entrypoints
        return nil
    },
}

func init() {
    rootCmd.AddCommand(proxyCmd)
    proxyCmd.AddCommand(proxyInitCmd)
    proxyCmd.AddCommand(proxyUpCmd)
    proxyCmd.AddCommand(proxyDownCmd)
    proxyCmd.AddCommand(proxyRestartCmd)
    proxyCmd.AddCommand(proxyStatusCmd)

    // proxy init flags
    proxyInitCmd.Flags().Bool("force", false, "overwrite existing configuration")
    proxyInitCmd.Flags().String("domain", "scind.test", "proxy domain for generated hostnames")
    proxyInitCmd.Flags().String("path", defaultProxyPath(), "directory to create proxy in")

    // proxy up flags
    proxyUpCmd.Flags().Bool("recreate", false, "recreate the proxy network even if it exists")
}

func defaultProxyPath() string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "scind", "proxy")
}
