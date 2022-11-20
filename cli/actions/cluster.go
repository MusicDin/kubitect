package actions

import (
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/file"
	"cli/tools/virtualenv"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	DefaultConfigDir    = "config"
	DefaultShareDir     = "share"
	DefaultTerraformDir = DefaultConfigDir + "/terraform"

	DefaultNewConfigFilename     = "kubitect.yaml"
	DefaultAppliedConfigFilename = "kubitect-applied.yaml"
	DefaultInfraConfigFilename   = "infrastructure.yaml"

	DefaultTerraformStateFilename = "terraform.tfstate"
	DefaultKubeconfigFilename     = "admin.conf"
)

type ClusterMeta struct {
	Name  string
	Path  string
	Local bool
}

func (c ClusterMeta) AppliedConfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultAppliedConfigFilename)
}

func (c ClusterMeta) InfrastructureConfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultInfraConfigFilename)
}

func (c ClusterMeta) KubeconfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultKubeconfigFilename)
}

func (c ClusterMeta) TfStatePath() string {
	return filepath.Join(c.Path, DefaultTerraformDir, DefaultTerraformStateFilename)
}

func (c ClusterMeta) ContainsKubeconfig() bool {
	return file.Exists(c.KubeconfigPath())
}

func (c ClusterMeta) ContainsAppliedConfig() bool {
	return file.Exists(c.TfStatePath())
}

func (c ClusterMeta) Valid() bool {
	return verifyClusterDir(c.Path) == nil
}

type Cluster struct {
	ClusterMeta

	NewConfigPath string

	// Configuration files
	NewConfig     *modelconfig.Config
	AppliedConfig *modelconfig.Config
	InfraConfig   *modelinfra.Config

	Ctx *env.Context
}

// NewCluster returns new Cluster instance with populated general fields.
// Cluster name and path are extracted from the provided configuration file.
// Previously applied configuration is also read, if cluster already exists.
func NewCluster(ctx *env.Context, userCfgPath string) (Cluster, error) {
	var err error
	var c Cluster

	c.Ctx = ctx
	c.NewConfigPath = userCfgPath
	c.NewConfig, err = readConfig(c.NewConfigPath, modelconfig.Config{})

	if err != nil {
		return c, err
	}

	if err := validateConfig(c.NewConfig); err != nil {
		return c, fmt.Errorf("Provided configuration file is not valid.\n%v", err)
	}

	c.Local = c.Ctx.Local()
	c.Name = string(*c.NewConfig.Cluster.Name)
	c.Path = filepath.Join(c.Ctx.ClustersDir(), c.Name)

	c.SetVirtualEnvironments()

	return c, c.Sync()
}

// Sync ensures that cluster properties are up to data.
func (c *Cluster) Sync() error {
	var err error

	c.AppliedConfig, err = readConfigIfExists(c.AppliedConfigPath(), modelconfig.Config{})

	if err != nil {
		return fmt.Errorf("failed to read previously applied configuration file: %v", err)
	}

	c.InfraConfig, err = readConfigIfExists(c.InfrastructureConfigPath(), modelinfra.Config{})

	if err != nil {
		return fmt.Errorf("failed to read infrastructure file: %v", err)
	}

	if c.InfraConfig != nil {
		if err := validateConfig(c.NewConfig); err != nil {
			return fmt.Errorf("Infrastructure file (produced by Terraform) is invalid.\n%v", err)
		}
	}

	return nil
}

func (c *Cluster) SetVirtualEnvironments() {

	t := virtualenv.MAIN

	virtualenv.Set(t, &virtualenv.VirtualEnv{
		Name:             string(t),
		RequirementsPath: "ansible/kubitect/requirements.txt",
		Path:             filepath.Join(c.Ctx.ShareDir(), "venv", string(t), c.KubitectVersion()),
		ClusterPath:      c.Path,
	})

	t = virtualenv.KUBESPRAY

	virtualenv.Set(t, &virtualenv.VirtualEnv{
		Name:             string(t),
		RequirementsPath: "ansible/kubespray/requirements.txt",
		Path:             filepath.Join(c.Ctx.ShareDir(), "venv", string(t), c.KubesprayVersion()),
		ClusterPath:      c.Path,
	})
}

// StoreNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func (c *Cluster) StoreNewConfig() error {
	src := c.NewConfigPath
	dst := filepath.Join(c.Path, DefaultConfigDir, DefaultNewConfigFilename)

	c.NewConfigPath = dst

	return file.ForceCopy(src, dst)
}

// ApplyNewConfig replaces currently applied config with new one.
func (c *Cluster) ApplyNewConfig() error {
	return file.Move(c.NewConfigPath, c.AppliedConfigPath())
}

func (c *Cluster) KubitectURL() string {
	if c.NewConfig.Kubitect.Url != nil {
		return string(*c.NewConfig.Kubitect.Url)
	}

	return env.ConstProjectUrl
}

func (c *Cluster) KubitectVersion() string {
	if c.NewConfig.Kubitect.Version != nil {
		return string(*c.NewConfig.Kubitect.Version)
	}

	return env.ConstProjectVersion
}

func (c *Cluster) KubesprayVersion() string {
	if c.NewConfig.Kubernetes.Kubespray.Version != nil {
		return string(*c.NewConfig.Kubernetes.Kubespray.Version)
	}

	return env.ConstKubesprayVersion
}

type ClustersMeta []ClusterMeta

func (cs ClustersMeta) Names() []string {
	var names []string
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}

func (cs ClustersMeta) Find(name string) *ClusterMeta {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}

	return nil
}

// Clusters returns cluster meta list of clusters from both project and local
// directory (if working directory is the Kubitect project).
func Clusters(ctx *env.Context) (ClustersMeta, error) {
	cs, err := clusters(ctx.ClustersDir(), false)

	if err != nil {
		return nil, err
	}

	lcs, err := clusters(ctx.LocalClustersDir(), true)

	if err == nil {
		cs = append(cs, lcs...)
	}

	return cs, nil
}

// clusters returns cluster meta list of clusters located in the given directory.
func clusters(clustersPath string, local bool) (ClustersMeta, error) {
	files, err := ioutil.ReadDir(clustersPath)

	if err != nil {
		return nil, fmt.Errorf("failed reading cluster directory: %v", err)
	}

	var clusters ClustersMeta

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()

			clusters = append(clusters, ClusterMeta{
				Name:  name,
				Path:  filepath.Join(clustersPath, name),
				Local: local,
			})
		}
	}

	return clusters, nil
}
