package cmd

import (
	"cli/env"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List initialized clusters",
	Long:  `Command 'list clusters' lists initialized clusters that are located in TKK_HOME/clusters directory.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return listClusters()
	},
}

func init() {
	listCmd.AddCommand(clustersCmd)
}

// listClusters lists all clusters located in the project clusters directory.
func listClusters() error {

	clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)

	files, err := ioutil.ReadDir(clustersPath)
	if err != nil {
		return fmt.Errorf("Failed reading a cluster directory: %w", err)
	}

	if len(files) > 0 {

		fmt.Println("Clusters:")
		for _, f := range files {
			fmt.Println("  - " + f.Name())
		}

		fmt.Println("\nSelect specific cluster with '--cluster <CLUSTER>' flag.")

	} else {
		fmt.Println("No clusters initialized yet. Run 'tkk apply' to create the default cluster.")
	}

	return nil
}
