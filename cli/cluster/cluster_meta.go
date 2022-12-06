package cluster

import (
	"cli/env"
	"cli/file"
	"cli/tools/terraform"
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

	tf *terraform.Terraform
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

func (c *ClusterMeta) Terraform() *terraform.Terraform {
	if c.tf != nil {
		return c.tf
	}

	tfVer := env.ConstTerraformVersion

	return &terraform.Terraform{
		Version:    tfVer,
		WorkingDir: filepath.Join(c.Path, "terraform"),
		BinDir:     filepath.Join(c.ShareDir(), "terraform", tfVer),
		Ui:         c.Ui(),
	}
}
