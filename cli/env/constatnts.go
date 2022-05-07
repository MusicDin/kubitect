package env

const (
	// Global constants
	ConstProjectHomeEnvName = "KUBITECT_HOME"
	ConstProjectHomeDir     = ".kubitect"
	ConstProjectClustersDir = "clusters"
	ConstProjectUrl         = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion     = "v2.0.4"
	ConstTerraformVersion   = "1.1.4"
	ConstTerraformStatePath = "config/terraform/terraform.tfstate"
	ConstKubeconfigPath     = "config/admin.conf"

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
		"templates/",
		"terraform/modules/",
		"terraform/output.tf",
		"terraform/variables.tf",
		"terraform/versions.tf",
		"scripts/",
		"requirements.txt",
		"LICENSE",
	}

	// Defines options for "apply --action" command.
	ProjectApplyActions = [...]string{
		"create",
		"upgrade",
		"scale",
	}
)
