package cmd

import (
	"cli/env"
	"cli/utils"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

// Root command (cli name)
var rootCmd = &cobra.Command{
	Use:   "kubitect",
	Short: "Kubitect",
	Long: `
Kubitect is a CLI tool that helps you manage multiple Kubernetes clusters.`,
	Version: env.ConstProjectVersion,

	// This is run when any command is run (also applies to all subcommands).
	// Flags are in this stage already resolved.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setup()
	},
}

// Execute adds all child commands to the root command and sets the flags
// accordingly.
func Execute() error {
	err := rootCmd.Execute()

	if err == nil {
		return nil
	}

	if _, ok := err.(utils.Error); ok {
		return err
	}

	if _, ok := err.(utils.Errors); ok {
		return err
	}

	return utils.NewError(err)
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SuggestionsMinimumDistance = 3

	rootCmd.AddGroup(
		&cobra.Group{
			Title: "Cluster Management Commands:",
			ID:    "mgmt",
		},
		&cobra.Group{
			Title: "Support Commands:",
			ID:    "support",
		},
		&cobra.Group{
			Title: "Other Commands:",
			ID:    "other",
		},
	)

	rootCmd.SetCompletionCommandGroupID("other")
	rootCmd.SetHelpCommandGroupID("other")

	rootCmd.PersistentFlags().BoolVar(&env.DebugMode, "debug", false, "enable debug messages")
}

// setup function is an entry point into CLI tool which sets global variables.
func setup() error {

	if env.Local {
		workingDir, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		env.ProjectHomePath = filepath.Join(workingDir, env.ConstProjectHomeDir)

	} else {
		userHomeDir, err := os.UserHomeDir()

		if err != nil {
			panic(err)
		}

		def := filepath.Join(userHomeDir, env.ConstProjectHomeDir)
		env.ProjectHomePath = utils.GetEnv(env.ConstProjectHomeEnvName, def)
	}

	if env.DebugMode {
		utils.PrintDebug("env.Local: %s", strconv.FormatBool(env.Local))
		utils.PrintDebug("env.ProjectHomePath: %s", env.ProjectHomePath)
	}

	return nil
}
