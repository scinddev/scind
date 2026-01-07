// Scaffold: internal/cli/root.go
// Migrated from specs/scind-go-stack.md

package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	workspace string
	apps      []string
	quiet     bool
	verbose   bool
	jsonOut   bool
	yamlOut   bool
	colorMode string // auto, always, never
)

var rootCmd = &cobra.Command{
	Use:   "scind",
	Short: "Workspace orchestration for Docker Compose",
	Long: `Scind is a workspace orchestration system for Docker Compose that enables
developers to run multiple isolated instances of multi-application stacks
simultaneously on a single host.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip context detection for global commands
		if isGlobalCommand(cmd) {
			return nil
		}
		return detectAndSetContext(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/scind/proxy.yaml)")
	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "specify workspace (overrides context detection)")
	rootCmd.PersistentFlags().StringSliceVarP(&apps, "app", "a", nil, "specify application(s) (repeatable, overrides context detection)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output, suppress context indicators and progress")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&yamlOut, "yaml", false, "output in YAML format")
	rootCmd.PersistentFlags().StringVar(&colorMode, "color", "auto", "color output: auto, always, or never")

	// Register flag completion
	rootCmd.RegisterFlagCompletionFunc("workspace", completeWorkspace)
	rootCmd.RegisterFlagCompletionFunc("app", completeApp)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/.config/scind")
		viper.SetConfigType("yaml")
		viper.SetConfigName("proxy")
	}

	viper.SetEnvPrefix("SCIND")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

// isGlobalCommand returns true for commands that don't use context detection
func isGlobalCommand(cmd *cobra.Command) bool {
	name := cmd.Name()
	return name == "port" || name == "proxy" || name == "config" ||
		name == "doctor" || name == "completion" || name == "help" ||
		name == "init-shell"
}

func detectAndSetContext(cmd *cobra.Command) error {
	// Implementation: walk up directory tree using workspace boundary approach
	// See internal/context/detector.go for full implementation
	//
	// Algorithm (workspace boundary):
	// 1. Walk up from CWD looking for workspace.yaml → establishes workspace root
	// 2. Walk up from CWD toward workspace root looking for application.yaml
	//    - Only consider application.yaml files WITHIN the workspace directory tree
	//    - Never traverse above workspace root (prevents vendor hijacking)
	// 3. Set results in viper:
	//    - viper.Set("detected.workspace", workspaceName)
	//    - viper.Set("detected.workspace_path", workspaceRoot)
	//    - viper.Set("detected.app", appName)
	//    - viper.Set("detected.app_path", appPath)
	//
	// Error messages for debugging:
	// - If application.yaml found but no workspace.yaml:
	//   "No workspace found (workspace.yaml) in current directory or any parent directories,
	//    but found an application (application.yaml) at: {path}"
	// - If neither found:
	//   "No workspace found (workspace.yaml) in current directory or any parent directories,
	//    and no application (application.yaml) found either"
	return nil
}

func completeWorkspace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// List available workspaces
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func completeApp(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// List available apps in current workspace
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}
