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
	Use:   "tkk",
	Short: "Terraform KVM Kubespray - TKK",
	Long: `tkk is a CLI tool that helps you manage multiple Kubernetes
clusters running on KVM.`,
	Version: "0.0.1",

	// This is run when any command is run (also applies to all subcommands).
	// Flags are in this stage already resolved.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := setup()
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&env.DebugMode, "debug", false, "enable debug messages")
}

// setup is an entry point into CLI tool and it sets
// various global variables.
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

		// Convert config filepath to absoultue
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
		env.ConfigPath = filepath.Join(env.ClusterPath, env.DefaultClusterConfigPath)
	}

	if env.DebugMode {
		fmt.Println("env.ConfigPath: " + env.ConfigPath)
		fmt.Println("env.ClusterPath: " + env.ClusterPath)
		fmt.Println("env.ProjectHomePath: " + env.ProjectHomePath)
		fmt.Println("env.Local: " + strconv.FormatBool(env.Local))
		fmt.Println("")
	}

	return nil
}
