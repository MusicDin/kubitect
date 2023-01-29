package cluster

import (
	"cli/cluster/executors"
	"cli/cluster/provisioner"
	"cli/cluster/provisioner/terraform"
	"cli/file"
	"cli/ui"
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

type ClusterContext interface {
	WorkingDir() string
	HomeDir() string
	ShareDir() string
	ClustersDir() string
	LocalClustersDir() string

	Local() bool
	ShowTerraformPlan() bool

	Ui() *ui.Ui
}

type ClusterMeta struct {
	ClusterContext

	Name  string
	Path  string
	Local bool

	exec executors.Executor
	prov provisioner.Provisioner
}

func (c ClusterMeta) AppliedConfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultAppliedConfigFilename)
}

func (c ClusterMeta) InfrastructureConfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultInfraConfigFilename)
}

func (c ClusterMeta) TfStatePath() string {
	return filepath.Join(c.Path, DefaultTerraformDir, DefaultTerraformStateFilename)
}

func (c ClusterMeta) KubeconfigPath() string {
	return filepath.Join(c.Path, DefaultConfigDir, DefaultKubeconfigFilename)
}

func (c ClusterMeta) ContainsAppliedConfig() bool {
	return file.Exists(c.AppliedConfigPath())
}

func (c ClusterMeta) ContainsTfStateConfig() bool {
	return file.Exists(c.TfStatePath())
}

func (c ClusterMeta) ContainsKubeconfig() bool {
	return file.Exists(c.KubeconfigPath())
}

func (c *ClusterMeta) Provisioner() provisioner.Provisioner {
	if c.prov != nil {
		return c.prov
	}

	c.prov, _ = terraform.NewTerraformProvisioner(
		c.Path,
		c.ShareDir(),
		true,
		nil,
		// c.Ui(),
	)

	return c.prov
}
