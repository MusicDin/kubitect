package cmd

import (
	"cli/env"
	"cli/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

// Root command (cli name)
var rootCmd = &cobra.Command{
	Use:   "kubitect",
	Short: "Kubitect",
	Long: `Kubitect is a CLI tool that helps you manage multiple Kubernetes
clusters.`,
	Version: "2.0.2",

	// This is run when any command is run (also applies to all subcommands).
	// Flags are in this stage already resolved.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := setup()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets the flags
// accordingly.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&env.DebugMode, "debug", false, "enable debug messages")
}

// setup function is an entry point into CLI tool which sets various global
// variables.
func setup() error {

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if env.Local {
		// Set project and cluster path to the current working
		// directory if local flag is set to true.
		env.ProjectHomePath = workingDir
		env.ClusterPath = workingDir
		env.ClusterName = "local"

	} else {
		env.ProjectHomePath = utils.GetEnv(env.ConstProjectHomeEnvName, filepath.Join(userHomeDir, env.ConstProjectHomeDir))
		env.ClusterPath = filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir, env.ClusterName)
	}

	// Make sure ConfigPath exists
	if len(env.ConfigPath) > 0 {

		env.IsCustomConfig = true

		// Convert config filepath to absolute path
		env.ConfigPath, err = filepath.Abs(env.ConfigPath)
		if err != nil {
			return fmt.Errorf("Config filepath '%s' cannot be converted to absolute path: %w", env.ConfigPath, err)
		}

		// Verify that the provided config file exists
		_, err = os.Stat(env.ConfigPath)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Config file does not exist on path '%s'.", env.ConfigPath)

		} else if err != nil {
			panic(err)
		}

	} else {
		env.IsCustomConfig = false
		env.ConfigPath = filepath.Join(env.ClusterPath, env.DefaultClusterConfigPath)
	}

	if env.DebugMode {
		utils.PrintDebug("env.ConfigPath: %s", env.ConfigPath)
		utils.PrintDebug("env.ClusterPath: %s", env.ClusterPath)
		utils.PrintDebug("env.ProjectHomePath: %s", env.ProjectHomePath)
		utils.PrintDebug("env.Local: %s", strconv.FormatBool(env.Local))
	}

	return nil
}
