// The env package provides constants for all other packages to consume,
// without creating import cycles.
//
// This package should not import any other packages.
package env

// Project related constants
const (
	ConstProjectUrl        = "https://github.com/MusicDin/kubitect"
	ConstProjectVersion    = "v3.0.1"
	ConstKubesprayUrl      = "https://github.com/kubernetes-sigs/kubespray"
	ConstKubesprayVersion  = "v2.21.0"
	ConstKubernetesVersion = "v1.25.6"
	ConstTerraformVersion  = "1.4.4"
)

// Defines applications that Kubitect depends on.
var ProjectRequiredApps = []string{
	"virtualenv",
	"python3",
	"git",
}

// Defines required files/directories that are copied from embedded
// resources, when cluster is created.
var ProjectRequiredFiles = []string{
	"ansible/",
	"terraform/",
}

// Defines options for "apply --action" command.
var ProjectApplyActions = [...]string{
	"create",
	"upgrade",
	"scale",
}

var ProjectOsPresets = map[string]string{
	"ubuntu":   "https://cloud-images.ubuntu.com/releases/jammy/release/ubuntu-22.04-server-cloudimg-amd64.img",
	"ubuntu22": "https://cloud-images.ubuntu.com/releases/jammy/release-20230302/ubuntu-22.04-server-cloudimg-amd64.img",
	"ubuntu20": "https://cloud-images.ubuntu.com/releases/focal/release-20230209/ubuntu-20.04-server-cloudimg-amd64.img",
	"debian":   "https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-generic-amd64.qcow2",
	"debian11": "https://cloud.debian.org/images/cloud/bullseye/20230124-1270/debian-11-genericcloud-amd64-20230124-1270.qcow2",
	"centos9":  "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20230405.1.x86_64.qcow2",
}

var ProjectK8sVersions = []string{
	"v1.23",
	"v1.24",
	"v1.25",
}
