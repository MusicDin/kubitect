package cluster

import (
	"path/filepath"

	"github.com/MusicDin/kubitect/cli/app"
	"github.com/MusicDin/kubitect/cli/cluster/executors"
	"github.com/MusicDin/kubitect/cli/cluster/provisioner"
	"github.com/MusicDin/kubitect/cli/cluster/provisioner/terraform"
	"github.com/MusicDin/kubitect/cli/utils/file"
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
	app.AppContext

	Name  string
	Path  string
	Local bool

	exec executors.Executor
	prov provisioner.Provisioner
}

// func NewClusterMeta(ctx app.AppContext, clusterPath string) (ClusterMeta, error) {
// 	cpStat, err := os.Stat(clusterPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("cluster meta: %v", err)
// 	}

// 	meta := clusterMeta{
// 		AppContext: ctx,
// 		Name:       cpStat.Name(),
// 		Path:       filepath.Join(clusterPath, cpStat.Name()),
// 		Local:      ctx.Local(),
// 	}

// 	if !cpStat.IsDir() || !meta.ContainsAppliedConfig() {
// 		return nil, fmt.Errorf("cluster meta: %s is not a cluster directory", err)
// 	}

// 	return &meta, nil
// }

func (c ClusterMeta) ConfigDir() string {
	return filepath.Join(c.Path, DefaultConfigDir)
}

func (c ClusterMeta) AppliedConfigPath() string {
	return filepath.Join(c.ConfigDir(), DefaultAppliedConfigFilename)
}

func (c ClusterMeta) InfrastructureConfigPath() string {
	return filepath.Join(c.ConfigDir(), DefaultInfraConfigFilename)
}

func (c ClusterMeta) TfStatePath() string {
	return filepath.Join(c.Path, DefaultTerraformDir, DefaultTerraformStateFilename)
}

func (c ClusterMeta) KubeconfigPath() string {
	return filepath.Join(c.ConfigDir(), DefaultKubeconfigFilename)
}

func (c ClusterMeta) PrivateSshKeyPath() string {
	return filepath.Join(c.ConfigDir(), ".ssh", "id_rsa")
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

	c.prov = terraform.NewTerraformProvisioner(
		c.Path,
		c.ShareDir(),
		c.ShowTerraformPlan(),
		nil,
	)

	return c.prov
}
