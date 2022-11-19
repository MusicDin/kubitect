package actions

import (
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/tools/virtualenv"
	"cli/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	newCfgFileName   = "kubitect.yaml"
	oldCfgFileName   = "kubitect-applied.yaml"
	infraCfgFileName = "infrastructure.yaml"
)

type Cluster struct {
	// Metadata
	Name  string
	Path  string
	Local bool

	// Configuration files
	NewCfg   *modelconfig.Config
	OldCfg   *modelconfig.Config
	InfraCfg *modelinfra.Config

	// Actual paths to the configuration files
	NewCfgPath   string
	OldCfgPath   string
	InfraCfgPath string
}

// NewCluster returns new Cluster instance with populated general fields.
// Cluster name and path are extracted from the provided configuration file.
// Previously applied configuration is also read, if cluster already exists.
func NewCluster(userCfgPath string) (Cluster, error) {
	var err error
	var c Cluster

	c.NewCfgPath = userCfgPath
	c.NewCfg, err = readConfig(c.NewCfgPath, modelconfig.Config{})

	if err != nil {
		return c, err
	}

	if err := validateConfig(c.NewCfg); err != nil {
		return c, fmt.Errorf("Provided configuration file is not valid.\n%v", err)
	}

	c.Name = string(*c.NewCfg.Cluster.Name)
	c.Path = env.ClusterPath(c.Name)
	c.Local = env.Local

	return c, c.Sync()
}

// Sync ensures that cluster properties are up to data.
func (c *Cluster) Sync() error {
	var err error

	c.OldCfgPath = path.Join(c.Path, env.ConstClusterConfigDir, oldCfgFileName)
	c.InfraCfgPath = path.Join(c.Path, env.ConstClusterConfigDir, infraCfgFileName)

	c.OldCfg, err = readConfigIfExists(c.OldCfgPath, modelconfig.Config{})

	if err != nil {
		return fmt.Errorf("failed to read previously applied configuration file: %v", err)
	}

	c.InfraCfg, err = readConfigIfExists(c.InfraCfgPath, modelinfra.Config{})

	if err != nil {
		return fmt.Errorf("failed to read infrastructure file: %v", err)
	}

	if c.InfraCfg != nil {
		if err := validateConfig(c.NewCfg); err != nil {
			return fmt.Errorf("Infrastructure file (produced by Terraform) is invalid.\n%v", err)
		}
	}

	return nil
}

// StoreNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func (c *Cluster) StoreNewConfig() error {
	src := c.NewCfgPath
	dst := path.Join(c.Path, env.ConstClusterConfigDir, newCfgFileName)

	c.NewCfgPath = dst

	return utils.ForceCopy(src, dst)
}

// ApplyNewConfig replaces currently applied config with new one.
func (c Cluster) ApplyNewConfig() error {
	return utils.ForceMove(c.NewCfgPath, c.OldCfgPath)
}

// setupMainVE creates main (Kubitect) virtual environment.
func (c Cluster) SetupMainVE() error {
	ktVer := env.ConstProjectVersion

	if c.NewCfg.Kubitect.Version != nil {
		ktVer = string(*c.NewCfg.Kubitect.Version)
	}

	return virtualenv.Env.Main.Setup(c.Path, ktVer)
}

// setupKubesprayVE creates Kubespray virtual environment.
func (c Cluster) SetupKubesprayVE() error {
	ksVer := string(*c.NewCfg.Kubernetes.Kubespray.Version)
	return virtualenv.Env.Kubespray.Setup(c.Path, ksVer)
}

// IsActive returns true if cluster contains terraform state file.
func (c Cluster) Active() bool {
	return utils.Exists(filepath.Join(c.Path, env.ConstTerraformStatePath))
}

func (c Cluster) Valid() bool {
	return verifyClusterDir(c.Path) == nil
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

// GetClusters returns clusters located in the project clusters directory.
// If current (working) directory is a Kubitect project, it also returns
// all local clusters. Provided filters are used for filtering found clusters.
func GetClusters() (Clusters, error) {
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

			clusters = append(clusters, Cluster{
				Name:  name,
				Path:  filepath.Join(path, name),
				Local: local,
			})
		}
	}

	return clusters, nil
}
