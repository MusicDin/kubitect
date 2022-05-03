package env

const (
	// Global constants
	ConstTerraformVersion   = "1.1.4"
	ConstProjectHomeEnvName = "KUBITECT_HOME"
	ConstProjectHomeDir     = ".kubitect"
	ConstProjectClustersDir = "clusters"
	ConstProjectUrl         = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion     = "v2.0.4"

	// default values
	DefaultClusterName       = "default"
	DefaultClusterAction     = "create"
	DefaultClusterConfigPath = "config/kubitect.yaml"
)

var (
	// Defines required files/directories that are copied from tmp git project.
	ProjectRequiredFiles = [...]string{
		"ansible/",
		"terraform/",
		"templates/",
		"resources/",
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
