// scaffold-compose-prefix.go
// compose-prefix command scaffold for Contrail CLI
// Create as: internal/cli/compose_prefix.go

package cli

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var composePrefixCmd = &cobra.Command{
    Use:    "compose-prefix",
    Short:  "Output docker compose command prefix",
    Long:   `Outputs a docker compose command prefix with project name and compose files for the current context.`,
    Hidden: true, // Internal command for shell integration
    RunE: func(cmd *cobra.Command, args []string) error {
        ws := viper.GetString("resolved.workspace")
        app := viper.GetString("resolved.app")

        if ws == "" || app == "" {
            return fmt.Errorf("no application context detected")
        }

        projectName := fmt.Sprintf("%s-%s", ws, app)

        // Get compose files from resolved flavor
        composeFiles, err := getComposeFilesForApp(ws, app)
        if err != nil {
            return err
        }

        // Build output
        fmt.Printf("docker compose -p %s", projectName)
        for _, f := range composeFiles {
            fmt.Printf(" -f '%s'", f)
        }
        fmt.Println()

        return nil
    },
}

func init() {
    rootCmd.AddCommand(composePrefixCmd)
    // Note: No --flavor flag. Flavor changes require regeneration and can impact
    // running applications. Users must use `contrail flavor set` instead.
}

func getComposeFilesForApp(workspace, app string) ([]string, error) {
    // Implementation: read application.yaml, resolve flavor, return file list
    return []string{}, nil
}
