package modelconfig

import v "cli/validation"

type LBDefault struct {
	CPU          *VCpu `yaml:"cpu"`
	RAM          *GB   `yaml:"ram"`
	MainDiskSize *GB   `yaml:"mainDiskSize"`
}

func (def LBDefault) Validate() error {
	return v.Struct(&def,
		v.Field(&def.CPU),
		v.Field(&def.RAM),
		v.Field(&def.MainDiskSize),
	)
}

type LB struct {
	VIP             *IPv4            `yaml:"vip"`
	VirtualRouterId *Uint8           `yaml:"virtualRouterId"`
	Default         *LBDefault       `yaml:"default"`
	Instances       *[]LBInstance    `yaml:"instances"`
	ForwardPorts    *[]LBPortForward `yaml:"forwardPorts"`
}

func (lb LB) Validate() error {
	return v.Struct(&lb,
		v.Field(&lb.VIP,
			v.Required().When(lb.Instances != nil && len(*lb.Instances) > 1).Error("Virtual IP (VIP) is required when multiple load balancer instances are configured."),
			v.OmitEmpty(), v.Custom(IP_IN_CIDR),
		),
		v.Field(&lb.VirtualRouterId),
		v.Field(&lb.Default),
		v.Field(&lb.Instances, v.OmitEmpty(), v.UniqueField("Id")),
		v.Field(&lb.ForwardPorts),
	)
}

func (lb LB) IPs() []string {
	if lb.Instances == nil {
		return nil
	}

	var ips []string

	if lb.VIP != nil {
		ips = append(ips, string(*lb.VIP))
	}

	for _, i := range *lb.Instances {
		if i.IP != nil {
			ips = append(ips, string(*i.IP))
		}
	}

	return ips
}

func (lb LB) MACs() []string {
	if lb.Instances == nil {
		return nil
	}

	var macs []string

	for _, i := range *lb.Instances {
		if i.MAC != nil {
			macs = append(macs, string(*i.MAC))
		}
	}

	return macs
}

type LBPortForward struct {
	Name       *string              `yaml:"name"`
	Port       *Port                `yaml:"port"`
	TargetPort *Port                `yaml:"targetPort"`
	Target     *LBPortForwardTarget `yaml:"target"`
}

func (pf LBPortForward) Validate() error {
	return v.Struct(&pf,
		v.Field(&pf.Name, v.Required(), v.AlphaNumericHypUS()),
		v.Field(&pf.Port, v.Required()),
		v.Field(&pf.TargetPort),
		v.Field(&pf.Target),
	)
}

type LBPortForwardTarget string

const (
	WORKERS LBPortForwardTarget = "workers"
	MASTERS LBPortForwardTarget = "masters"
	ALL     LBPortForwardTarget = "all"
)

func (pft LBPortForwardTarget) Validate() error {
	return v.Var(pft, v.OmitEmpty(), v.OneOf(WORKERS, MASTERS, ALL))
}

type LBInstance struct {
	Id           *string `yaml:"id" opt:",id"`
	Host         *string `yaml:"host"`
	IP           *IPv4   `yaml:"ip"`
	MAC          *MAC    `yaml:"mac"`
	CPU          *VCpu   `yaml:"cpu"`
	RAM          *GB     `yaml:"ram"`
	MainDiskSize *GB     `yaml:"mainDiskSize"`
	Priority     *Uint8  `yaml:"priority"`
}

func (i LBInstance) Validate() error {
	return v.Struct(&i,
		v.Field(&i.Id, v.Required()),
		v.Field(&i.Host, v.OmitEmpty(), v.Custom(VALID_HOST)),
		v.Field(&i.IP, v.OmitEmpty(), v.Custom(IP_IN_CIDR)),
		v.Field(&i.MAC),
		v.Field(&i.CPU),
		v.Field(&i.RAM),
		v.Field(&i.MainDiskSize),
		v.Field(&i.Priority),
	)
}
