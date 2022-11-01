package env

import (
	"path"
)

var (
	// Global shared variables (should always have a valid value)
	// See cmd/root:setup
	ProjectHomePath string

	// Local options (flags)
	Local       bool
	AutoApprove bool

	// Global options (flags)
	DebugMode bool
)

// clusterPath returns path of the current cluster.
func ClusterPath(clusterName string) string {
	return path.Join(ProjectHomePath, ConstProjectClustersDir, clusterName)
}

func BinDirPath(binName string, version string) string {
	return path.Join(ProjectHomePath, ConstSharedDir, "bin", binName, version)
}

func VenvDirPath(venvName string, version string) string {
	return path.Join(ProjectHomePath, ConstSharedDir, "venv", venvName, version)
}
