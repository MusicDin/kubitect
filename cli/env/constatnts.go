package env

const (
	// Global constants
	ConstTerraformVersion   = "1.1.4"
	ConstProjectHomeEnvName = "TKK_HOME"
	ConstProjectHomeDir     = ".tkk"
	ConstProjectClustersDir = "clusters"

	// default values
	DefaultClusterName       = "default"
	DefaultClusterAction     = "create"
	DefaultClusterConfigPath = "config/tkk.yaml"
	DefaultGitProjectUrl     = "https://github.com/MusicDin/terraform-kvm-kubespray"
	DefaultGitProjectVersion = "feature/multiple-servers"
)

var (
	// Defines required files/directories that are copied from tmp git project.
	ProjectRequiredFiles = [...]string{
		"ansible/",
		"terraform/",
		"templates/",
		"examples/",
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
