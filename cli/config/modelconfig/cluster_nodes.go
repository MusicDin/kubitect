package modelconfig

import v "cli/validation"

type NodeType interface {
	GetInstances() []Instance
}

type Instance interface {
	GetIP() *string
}

type Nodes struct {
	LoadBalancer *LB     `yaml:"loadBalancer"`
	Master       *Master `yaml:"master"`
	Worker       *Worker `yaml:"worker"`
}

func (n Nodes) Validate() error {
	defer v.RemoveCustomValidator(LB_REQUIRED)

	v.RegisterCustomValidator(LB_REQUIRED, n.isLBRequiredValidator())

	return v.Struct(&n,
		v.Field(&n.LoadBalancer, v.OmitEmpty(), v.UniqueField("Name")),
		v.Field(&n.Master, v.UniqueField("Name")),
		v.Field(&n.Worker, v.OmitEmpty(), v.UniqueField("Name")),
	)
}

// isLBRequired is a cross-validator that triggers an error when multiple master
// nodes are configured, but the load balancer is not.
func (n Nodes) isLBRequiredValidator() v.Validator {
	if n.Master == nil || n.Master.Instances == nil || len(*n.Master.Instances) <= 1 {
		return v.None
	}

	if n.LoadBalancer == nil || n.LoadBalancer.Instances == nil || len(*n.LoadBalancer.Instances) == 0 {
		return v.Fail().Error("At least one load balancer instance is required when multiple master instances are configured.")
	}

	return v.None
}

func (n Nodes) IPs() []string {
	var ips []string

	if n.LoadBalancer != nil {
		ips = append(ips, n.LoadBalancer.IPs()...)
	}

	if n.Master != nil {
		ips = append(ips, n.Master.IPs()...)
	}

	if n.Worker != nil {
		ips = append(ips, n.Worker.IPs()...)
	}

	return ips
}

func (n Nodes) MACs() []string {
	var macs []string

	if n.LoadBalancer != nil {
		macs = append(macs, n.LoadBalancer.MACs()...)
	}

	if n.Master != nil {
		macs = append(macs, n.Master.MACs()...)
	}

	if n.Worker != nil {
		macs = append(macs, n.Worker.MACs()...)
	}

	return macs
}
