// The env package provides constants for all other packages to consume,
// without creating import cycles.
//
// This package should not import any other packages.
package env

// Project related constants
const (
	ConstProjectUrl       = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion   = "v2.2.0"
	ConstKubesprayUrl     = "https://github.com/kubernetes-sigs/kubespray"
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

var ProjectOsPresets = map[string]string{
	"ubuntu":   "https://cloud-images.ubuntu.com/releases/jammy/release/ubuntu-22.04-server-cloudimg-amd64.img",
	"ubuntu22": "https://cloud-images.ubuntu.com/releases/jammy/release-20220712/ubuntu-22.04-server-cloudimg-amd64.img",
	"ubuntu20": "https://cloud-images.ubuntu.com/releases/focal/release-20220711/ubuntu-20.04-server-cloudimg-amd64.img",
	"debian":   "https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-generic-amd64.qcow2",
	"debian11": "https://cloud.debian.org/images/cloud/bullseye/20220711-1073/debian-11-genericcloud-amd64-20220711-1073.qcow2",
}
