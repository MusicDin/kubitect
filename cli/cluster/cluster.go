package cluster

import (
	"cli/cluster/executors"
	"cli/cluster/executors/kubespray"
	"cli/cluster/provisioner"
	"cli/cluster/provisioner/terraform"
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/file"
	"cli/tools/virtualenv"
	"cli/ui"
	"fmt"
	"path/filepath"
)

type Cluster struct {
	ClusterMeta

	NewConfigPath string

	// Configuration files
	NewConfig     *modelconfig.Config
	AppliedConfig *modelconfig.Config
	InfraConfig   *modelinfra.Config
}

// NewCluster returns new Cluster instance with populated general fields.
// Cluster name and path are extracted from the provided configuration file.
// Previously applied configuration is also read, if cluster already exists.
func NewCluster(ctx ClusterContext, configPath string) (*Cluster, error) {
	newCfg, err := readConfig(configPath, modelconfig.Config{})

	if err != nil {
		return nil, err
	}

	c := &Cluster{
		ClusterMeta: ClusterMeta{
			ClusterContext: ctx,
			Local:          ctx.Local(),
		},
		NewConfig:     newCfg,
		NewConfigPath: configPath,
	}

	if err := validateConfig(c.NewConfig); err != nil {
		ui.PrintBlockE(err...)
		return c, fmt.Errorf("Provided configuration file is not valid.")
	}

	c.Name = string(*c.NewConfig.Cluster.Name)
	c.Path = filepath.Join(c.ClustersDir(), c.Name)

	return c, c.Sync()
}

// Sync ensures that cluster configuration files are up to data.
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
			ui.PrintBlockE(err...)
			return fmt.Errorf("Infrastructure file (produced by Terraform) is invalid.")
		}
	}

	return nil
}

func (c *Cluster) Executor() executors.Executor {

	if c.exec != nil {
		return c.exec
	}

	veName := "venv"
	vePath := filepath.Join(c.ShareDir(), "venv", "kubespray", c.KubesprayVersion())
	veReqPath := "ansible/kubespray/requirements.txt"
	veWorkingDir := c.Path

	ve := virtualenv.NewVirtualEnv(veName, vePath, veWorkingDir, veReqPath)

	c.exec, _ = kubespray.NewKubesprayExecutor(
		c.Name,
		c.Path,
		string(*c.NewConfig.Kubernetes.Version),
		string(*c.InfraConfig.Cluster.NodeTemplate.User),
		string(*c.InfraConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath),
		ve,
	)

	return c.exec
}

func (c *Cluster) Provisioner() provisioner.Provisioner {
	if c.prov != nil {
		return c.prov
	}

	c.prov, _ = terraform.NewTerraformProvisioner(
		c.Path,
		c.ShareDir(),
		true,
		c.NewConfig.Hosts,
	)

	return c.prov
}

// ApplyNewConfig replaces currently applied config with new one.
func (c *Cluster) ApplyNewConfig() error {
	return file.ForceCopy(c.NewConfigPath, c.AppliedConfigPath())
}

// StoreNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func (c *Cluster) StoreNewConfig() error {
	src := c.NewConfigPath
	dst := filepath.Join(c.Path, DefaultConfigDir, DefaultNewConfigFilename)

	c.NewConfigPath = dst

	return file.ForceCopy(src, dst)
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
