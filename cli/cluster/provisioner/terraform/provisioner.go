package terraform

import (
	"cli/cluster/provisioner"
	"cli/config/modelconfig"
	"cli/env"
	"cli/tools/terraform"
	"path"
)

func NewTerraformProvisioner(
	clusterPath,
	sharedPath string,
	showPlan bool,
	hosts []modelconfig.Host,
) (
	provisioner.Provisioner,
	error,
) {
	tfVer := env.ConstTerraformVersion

	binDir := path.Join(sharedPath, "terraform", tfVer)
	projDir := path.Join(clusterPath, "terraform")

	err := NewMainTemplate(projDir, hosts).Write()
	if err != nil {
		return nil, err
	}

	return terraform.NewTerraform(tfVer, binDir, projDir, showPlan), nil
}
