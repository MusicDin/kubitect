package modelconfig

import v "cli/validation"

type Nodes struct {
	LoadBalancer *LB     `yaml:"loadBalancer"`
	Master       *Master `yaml:"master"`
	Worker       *Worker `yaml:"worker"`
}

func (n Nodes) Validate() error {
	defer v.RemoveCustomValidator(LB_REQUIRED)

	v.RegisterCustomValidator(LB_REQUIRED, n.isLBRequiredValidator())

	return v.Struct(&n,
		v.Field(&n.LoadBalancer),
		v.Field(&n.Master),
		v.Field(&n.Worker),
	)
}

// isLBRequired is a custom cross-validator that triggers an error when
// multiple master nodes are configured, but the load balancer is not.
func (n Nodes) isLBRequiredValidator() v.Validator {
	if n.Master == nil || n.Master.Instances == nil || len(*n.Master.Instances) == 0 {
		return v.None
	}

	return v.Required().When(len(*n.Master.Instances) > 1).Error("At least one load balancer instance is required when multiple master instances are configured.")
}
