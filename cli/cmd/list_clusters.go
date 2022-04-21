package cmd

import (
	"cli/env"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var listClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List initialized clusters",
	Long: `Lists clusters that have been initialized in $KUBITECT_HOME.
Local clusters (applied with --local flag) are not listed.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return listClusters()
	},
}

func init() {
	listCmd.AddCommand(listClustersCmd)
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
		fmt.Println("No clusters initialized yet. Run 'kubitect apply' to create the default cluster.")
	}

	return nil
}
