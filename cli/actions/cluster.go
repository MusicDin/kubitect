package actions

import (
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/tools/virtualenv"
	"cli/utils"
	"cli/validation"
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
	Name         string
	Path         string
	Local        bool
	Invalid      bool
	NewCfg       *modelconfig.Config
	NewCfgPath   string
	OldCfg       *modelconfig.Config
	OldCfgPath   string
	InfraCfg     *modelinfra.Config
	InfraCfgPath string
}

// NewCluster returns new Cluster instance with populated general fields.
// Cluster name and path are extracted from the provided configuration file.
// Previously applied configuration is also read, if cluster already exists.
func NewCluster(userCfgPath string) (Cluster, error) {
	var err error
	var c Cluster

	c.NewCfg, err = readConfig(userCfgPath, modelconfig.Config{})

	if err != nil {
		return c, err
	}

	if err := validateConfig(c.NewCfg); err != nil {
		e := fmt.Errorf("Provided configuration file is not valid:")
		return c, utils.NewErrors(e, err)
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

	if err != nil {
		return fmt.Errorf("Error reading previously applied configuration file: %v", err)
	}

	c.InfraCfg, err = readConfigIfExists(c.InfraCfgPath, modelinfra.Config{})

	if err != nil {
		return fmt.Errorf("Error reading infrastructure file: %v", err)
	}

	if c.InfraCfg != nil {
		if err := validateConfig(c.NewCfg); err != nil {
			e := fmt.Errorf("Infrastructure file (produced by Terraform) is invalid.")
			return utils.NewErrors(e, err)
		}
	}

	return nil
}

// readConfig reads configuration file on the given path and converts it into
// the provided model.
func readConfig[T validation.Validatable](path string, model T) (*T, error) {
	if !utils.Exists(path) {
		return nil, fmt.Errorf("file '%s' does not exist", path)
	}

	return utils.ReadYaml(path, model)
}

// readConfig reads configuration file on the given path and converts it into
// the provided model. If file on the provided path does not exist, neither error
// nor model is returned.
func readConfigIfExists[T validation.Validatable](path string, model T) (*T, error) {
	if !utils.Exists(path) {
		return nil, nil
	}

	return utils.ReadYaml(path, model)
}

// validateConfig validates provided configuration file.
func validateConfig[T validation.Validatable](config T) error {
	var errs utils.Errors

	err := config.Validate()

	for _, e := range err.(validation.ValidationErrors) {
		errs = append(errs, NewValidationError(e.Error(), e.Namespace))
	}

	return errs
}

// StoreNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func (c *Cluster) StoreNewConfig() error {
	src := c.NewCfgPath
	dst := path.Join(c.Path, env.ConstClusterConfigDir, newCfgFileName)

	c.NewCfgPath = dst

	return utils.CopyFile(src, dst)
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
