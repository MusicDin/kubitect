package cmd

import (
	"cli/env"
	"cli/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

type ClusterFilter uint8

const (
	IsActive           ClusterFilter = iota
	IsInactive         ClusterFilter = iota
	ContainsKubeconfig ClusterFilter = iota
	ContainsConfig     ClusterFilter = iota
)

// clustersCmd represents the clusters command
var listClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List initialized clusters",
	Long: `Lists clusters that have been initialized.
Local clusters (applied with --local flag) are not listed.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return listClusters()
	},
}

func init() {
	listCmd.AddCommand(listClustersCmd)
}

// listClusters lists all clusters located in the project clusters directory
// that contain terraform state file.
func listClusters() error {

	clusters, err := GetClusters(nil)
	if err != nil {
		return err
	}

	activeClusters, err := GetClusters([]ClusterFilter{IsActive})
	if err != nil {
		return err
	}

	// Tag active clusters (intersection of active and all clusters).
	for i, cluster := range clusters {

		if utils.StrArrayContains(activeClusters, cluster) {
			clusters[i] = cluster + " (active)"
		}
	}

	// Print clusters.
	if len(clusters) > 0 {

		fmt.Println("Clusters:")

		for _, clusterName := range clusters {
			fmt.Println("  - " + clusterName)
		}

	} else {

		fmt.Println("No clusters initialized yet. Run 'kubitect apply' to create the cluster.")
	}

	return nil
}

// GetClusters returns clusters located in the project clusters directory.
// If filters are provided, clusters will be appropriately filtered.
func GetClusters(filters []ClusterFilter) ([]string, error) {

	clusters, err := getAllClusters()
	if err != nil {
		return nil, err
	}

	filteredClusters, err := filterClusters(clusters, filters)
	if err != nil {
		return nil, err
	}

	return filteredClusters, nil
}

// getAllClusters returns clusters located in the project clusters directory.
func getAllClusters() ([]string, error) {

	clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)

	files, err := ioutil.ReadDir(clustersPath)
	if err != nil {
		return nil, fmt.Errorf("Failed reading a cluster directory: %w", err)
	}

	clusterNames := []string{}

	for _, file := range files {

		// Only list directories.
		if file.IsDir() {
			clusterNames = append(clusterNames, file.Name())
		}
	}

	return clusterNames, nil
}

// filterClusters returns clusters that pass provided filters.
func filterClusters(clusters []string, filters []ClusterFilter) ([]string, error) {

	clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)

	// Return all clusters if filters are not provided
	if filters == nil || len(filters) == 0 {
		return clusters, nil
	}

	filteredClusters := []string{}

	for _, cluster := range clusters {

		clusterPath := filepath.Join(clustersPath, cluster)
		passedFilters := true

		// Iterate over provided filters
		for _, filter := range filters {

			switch filter {

			case IsActive:

				// To determine if cluster is active, check if terraform state
				// file exists.
				if !utils.Exists(filepath.Join(clusterPath, env.ConstTerraformStatePath)) {
					passedFilters = false
					break
				}

			case IsInactive:

				if utils.Exists(filepath.Join(clusterPath, env.ConstTerraformStatePath)) {
					passedFilters = false
					break
				}

			case ContainsKubeconfig:

				if !utils.Exists(filepath.Join(clusterPath, env.ConstKubeconfigPath)) {
					passedFilters = false
					break
				}

			case ContainsConfig:

				if !utils.Exists(filepath.Join(clusterPath, env.DefaultClusterConfigPath)) {
					passedFilters = false
					break
				}
			}
		}

		if passedFilters {
			filteredClusters = append(filteredClusters, cluster)
		}
	}

	return filteredClusters, nil
}
