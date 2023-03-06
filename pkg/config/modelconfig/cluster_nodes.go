package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/validation"
)

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
	Master       Master `yaml:"master"`
	Worker       Worker `yaml:"worker,omitempty"`
	LoadBalancer LB     `yaml:"loadBalancer,omitempty"`
}

func (n Nodes) Validate() error {
	defer validation.RemoveCustomValidator(LB_REQUIRED)

	validation.RegisterCustomValidator(LB_REQUIRED, n.isLBRequiredValidator())

	return validation.Struct(&n,
		validation.Field(&n.LoadBalancer),
		validation.Field(&n.Master),
		validation.Field(&n.Worker),
	)
}

// isLBRequired is a cross-validator that triggers an error when multiple master
// nodes are configured, but the load balancer is not.
func (n Nodes) isLBRequiredValidator() validation.Validator {
	if len(n.Master.Instances) <= 1 {
		return validation.None
	}

	if len(n.LoadBalancer.Instances) == 0 {
		return validation.Fail().Error("At least one load balancer instance is required when multiple master instances are configured.")
	}

	return validation.None
}

func (n Nodes) Instances() []Instance {
	var ins []Instance

	for _, i := range n.Master.Instances {
		ins = append(ins, i)
	}

	for _, i := range n.Worker.Instances {
		ins = append(ins, i)
	}

	for _, i := range n.LoadBalancer.Instances {
		ins = append(ins, i)
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
