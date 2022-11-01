package actions

import (
	"cli/config/modelconfig"
	"cli/env"
	"cli/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Cluster struct {
	Name    string
	Path    string
	Local   bool
	Invalid bool
	NewCfg  *modelconfig.Config
	OldCfg  *modelconfig.Config
}

// NewCluster returns new Cluster instance with populated general fields.
// Cluster name and path are extracted from the provided configuration file.
// Previously applied configuration is also read, if cluster already exists.
func NewCluster(userCfgPath string) (Cluster, error) {
	var c Cluster

	newCfg, err := readNewConfig(userCfgPath)

	if err != nil {
		return c, err
	}

	if err := validateNewConfig(newCfg); err != nil {
		return Cluster{}, err
	}

	c.NewCfg = newCfg
	c.Name = string(*newCfg.Cluster.Name)
	c.Path = env.ClusterPath(c.Name)
	c.Local = env.Local

	oldCfg, err := readOldConfig(c.Name)

	if err != nil {
		return c, err
	}

	c.OldCfg = oldCfg

	return c, nil
}

// IsActive returns true if cluster contains terraform state file.
func (c Cluster) Active() bool {
	return utils.Exists(filepath.Join(c.Path, env.ConstTerraformStatePath))
}

// Returns true if cluster contains Kubeconfig file.
func (c Cluster) ContainsKubeconfig() bool {
	return utils.Exists(filepath.Join(c.Path, env.ConstKubeconfigPath))
}

// Returns true if cluster contains applied configuration file.
func (c Cluster) ContainsConfig() bool {
	return utils.Exists(filepath.Join(c.Path, env.ConstClusterConfigPath))
}

type Clusters []Cluster

func (cs Clusters) Names() []string {
	var names []string
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}

func (cs Clusters) Find(name string) *Cluster {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}

	return nil
}

// ReadClustersInfo returns clusters located in the project clusters directory.
// If current (working) directory is a Kubitect project, it also returns
// all local clusters. Provided filters are used for filtering found clusters.
func ReadClustersInfo() (Clusters, error) {
	path := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)

	clusters, err := clusters(path, false)

	if err != nil {
		return nil, err
	}

	localClusters, err := localClusters()

	if err == nil {
		clusters = append(clusters, localClusters...)
	}

	return clusters, nil
}

// localClusters returns a list of local clusters. Local clusters are
// searched for only if current directory is Kubitect project.
func localClusters() (Clusters, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	path := filepath.Join(wd, env.ConstProjectHomeDir, env.ConstProjectClustersDir)

	return clusters(path, true)
}

// getAllClusters returns clusters located in the project clusters directory.
func clusters(path string, local bool) (Clusters, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("Failed reading cluster directory: %v", err)
	}

	var clusters []Cluster

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()

			var invalid bool

			if err := verifyClusterDir(path); err != nil {
				invalid = true
			}

			clusters = append(clusters, Cluster{
				Name:    name,
				Path:    filepath.Join(path, name),
				Local:   local,
				Invalid: invalid,
			})
		}
	}

	return clusters, nil
}
