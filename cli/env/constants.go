// The env package provides constants for all other packages to consume,
// without creating import cycles.
//
// This package should not import any other packages.
package env

// Project related constants
const (
	ConstProjectUrl       = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion   = "v2.2.0"
	ConstKubesprayVersion = "v2.20.0"
	ConstTerraformVersion = "1.3.7"
)

// Defines applications that Kubitect depends on.
var ProjectRequiredApps = []string{
	"virtualenv",
	"python3",
	"git",
}

// Defines required files/directories that are copied from tmp git project.
var ProjectRequiredFiles = []string{
	"ansible/",
	"resources/",
	"terraform/modules/",
	"terraform/templates/",
	"terraform/scripts/",
	"terraform/main.tf.tpl",
	"terraform/defaults.yaml",
	"terraform/output.tf",
	"terraform/variables.tf",
	"terraform/versions.tf",
	// "LICENSE",
}

// Defines options for "apply --action" command.
var ProjectApplyActions = [...]string{
	"create",
	"upgrade",
	"scale",
}
