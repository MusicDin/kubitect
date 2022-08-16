package env

const (
	// Global constants
	ConstProjectHomeEnvName = "KUBITECT_HOME"
	ConstProjectHomeDir     = ".kubitect"
	ConstProjectClustersDir = "clusters"
	ConstProjectUrl         = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion     = "v2.1.0"
	ConstTerraformVersion   = "1.2.4"
	ConstTerraformStatePath = "config/terraform/terraform.tfstate"
	ConstKubeconfigPath     = "config/admin.conf"
	ConstVenvBinDir         = "bin/venvs"

	// default values
	DefaultClusterName       = "default"
	DefaultClusterAction     = "create"
	DefaultClusterConfigPath = "config/kubitect.yaml"
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
		"LICENSE",
	}

	// Defines options for "apply --action" command.
	ProjectApplyActions = [...]string{
		"create",
		"upgrade",
		"scale",
	}
)
