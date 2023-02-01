package modelconfig

import v "cli/utils/validation"

const (
	defaultVCpu         = VCpu(2)
	defaultRAM          = GB(4)
	defaultMainDiskSize = GB(32)
)

type Instance interface {
	GetTypeName() string
	GetID() string
	GetIP() IPv4
	GetMAC() MAC
}

type Nodes struct {
	LoadBalancer LB     `yaml:"loadBalancer"`
	Master       Master `yaml:"master"`
	Worker       Worker `yaml:"worker"`
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

// isLBRequired is a cross-validator that triggers an error when multiple master
// nodes are configured, but the load balancer is not.
func (n Nodes) isLBRequiredValidator() v.Validator {
	if len(n.Master.Instances) <= 1 {
		return v.None
	}

	if len(n.LoadBalancer.Instances) == 0 {
		return v.Fail().Error("At least one load balancer instance is required when multiple master instances are configured.")
	}

	return v.None
}

func (n Nodes) Instances() []Instance {
	var ins []Instance

	for _, i := range n.LoadBalancer.Instances {
		ins = append(ins, Instance(i))
	}

	for _, i := range n.Master.Instances {
		ins = append(ins, Instance(i))
	}

	for _, i := range n.Worker.Instances {
		ins = append(ins, Instance(i))
	}

	return ins
}

func (n Nodes) IPs() []string {
	var ips []string

	for _, i := range n.Instances() {
		ip := i.GetIP()
		if ip != "" {
			ips = append(ips, string(ip))
		}
	}

	return ips
}

func (n Nodes) MACs() []string {
	var macs []string

	for _, i := range n.Instances() {
		mac := i.GetMAC()

		if mac != "" {
			macs = append(macs, string(mac))
		}
	}

	return macs
}
