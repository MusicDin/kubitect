package env

import "fmt"

const (
	ConstClusterConfigDir = "config"
	ConstSharedDir        = "share"

	ConstProjectHomeEnvName = "KUBITECT_HOME"
	ConstProjectHomeDir     = ".kubitect"
	ConstProjectClustersDir = "clusters"

	ConstProjectUrl       = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion   = "v2.2.0"
	ConstKubesprayVersion = "v2.20.0"
	ConstTerraformVersion = "1.2.4"

	ConstTerraformStatePath = ConstClusterConfigDir + "/terraform/terraform.tfstate"
	ConstKubeconfigPath     = ConstClusterConfigDir + "/admin.conf"
	ConstClusterConfigPath  = ConstClusterConfigDir + "/kubitect.yaml"

	// default values
	// DefaultClusterName       = "default"
	// DefaultClusterAction     = "create"
)

var (
	// Defines required files/directories that are copied from tmp git project.
	ProjectRequiredFiles = [...]string{
		"ansible/",
		"resources/",
		"terraform/modules/",
		"terraform/templates/",
		"terraform/scripts/",
		"terraform/defaults.yaml",
		"terraform/output.tf",
		"terraform/variables.tf",
		"terraform/versions.tf",
		// "LICENSE",
	}

	// Defines options for "apply --action" command.
	ProjectApplyActions = [...]string{
		"create",
		"upgrade",
		"scale",
	}
)

type ApplyAction string

const (
	UNKNOWN ApplyAction = "unknown"
	CREATE  ApplyAction = "create"
	UPGRADE ApplyAction = "upgrade"
	SCALE   ApplyAction = "scale"
)

func (a ApplyAction) String() string {
	return string(a)
}

func ToApplyAction(a string) (ApplyAction, error) {
	switch a {
	case CREATE.String(), "":
		return CREATE, nil
	case UPGRADE.String():
		return UPGRADE, nil
	case SCALE.String():
		return SCALE, nil
	default:
		return UNKNOWN, fmt.Errorf("Unknown cluster action: %s", a)
	}
}
