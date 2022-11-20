package env

/*
 * Global flags
 */

var (
	// Automatically approve all user requests.
	AutoApprove bool

	// Show debug messages
	Debug bool

	// Prevent colored output
	NoColor bool
)

/*
 * Globally accessible constants
 */

// Project related constants
const (
	ConstProjectUrl       = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion   = "v2.2.0"
	ConstKubesprayVersion = "v2.20.0"
	ConstTerraformVersion = "1.2.4"
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
	ProjectApplyActions = [3]string{
		"create",
		"upgrade",
		"scale",
	}
)
