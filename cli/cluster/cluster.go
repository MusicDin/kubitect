package cluster

import (
	"cli/app"
	"cli/cluster/executors"
	"cli/cluster/executors/kubespray"
	"cli/cluster/provisioner"
	"cli/cluster/provisioner/terraform"
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/tools/virtualenv"
	"cli/ui"
	"cli/utils/defaults"
	"cli/utils/file"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
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
func NewCluster(ctx app.AppContext, configPath string) (*Cluster, error) {
	newCfg, err := readConfig(configPath, modelconfig.Config{})
	if err != nil {
		return nil, err
	}

	c := &Cluster{
		ClusterMeta: ClusterMeta{
			AppContext: ctx,
			Local:      ctx.Local(),
		},
		NewConfig:     newCfg,
		NewConfigPath: configPath,
	}

	if err := defaults.Set(c.NewConfig); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %v", err)
	}

	if err := validateConfig(c.NewConfig); err != nil {
		ui.PrintBlockE(err...)
		return nil, fmt.Errorf("invalid configuration file")
	}

	// Throw an error if cluster name contains prefix local.
	// Prepend prefix "local" to the cluster name, if cluster is local.
	if strings.HasPrefix(c.NewConfig.Cluster.Name, "local") {
		return nil, fmt.Errorf("Cluster name cannot contain a prefix 'local'. This prefix is reserved for clusters created with --local flag.")
	} else if ctx.Local() {
		c.NewConfig.Cluster.Name = "local-" + c.NewConfig.Cluster.Name
	}

	c.Name = c.NewConfig.Cluster.Name
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
		if err := validateConfig(c.InfraConfig); err != nil {
			ui.PrintBlockE(err...)
			return fmt.Errorf("infrastructure file (produced by Terraform) is invalid")
		}
	}

	return nil
}

func (c *Cluster) Executor() executors.Executor {
	if c.exec != nil {
		return c.exec
	}

	veReqPath := "ansible/kubespray/requirements.txt"
	vePath := path.Join(c.ShareDir(), "venv", "kubespray", env.ConstKubesprayVersion)
	ve := virtualenv.NewVirtualEnv(vePath, c.Path, veReqPath)

	c.exec = kubespray.NewKubesprayExecutor(
		c.Name,
		c.Path,
		c.PrivateSshKeyPath(),
		c.ConfigDir(),
		c.NewConfig,
		c.InfraConfig,
		ve,
	)

	return c.exec
}

func (c *Cluster) Provisioner() provisioner.Provisioner {
	if c.prov != nil {
		return c.prov
	}

	c.prov = terraform.NewTerraformProvisioner(
		c.Path,
		c.ShareDir(),
		c.ShowTerraformPlan(),
		c.NewConfig,
	)

	return c.prov
}

// ApplyNewConfig replaces currently applied config with new one.
func (c *Cluster) ApplyNewConfig() error {
	err := os.MkdirAll(path.Dir(c.AppliedConfigPath()), 0744)
	if err != nil {
		return err
	}

	return file.WriteYaml(c.NewConfig, c.AppliedConfigPath(), 0644)
}

// StoreNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func (c *Cluster) StoreNewConfig() error {
	src := c.NewConfigPath
	dst := filepath.Join(c.Path, DefaultConfigDir, DefaultNewConfigFilename)

	c.NewConfigPath = dst

	return file.ForceCopy(src, dst)
}
