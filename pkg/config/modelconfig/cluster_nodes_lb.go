package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
)

var (
	defaultVRID     = Uint8(51)
	defaultPriority = Uint8(10)
)

type LBDefault struct {
	CPU          VCpu `yaml:"cpu"`
	RAM          GB   `yaml:"ram"`
	MainDiskSize GB   `yaml:"mainDiskSize"`
}

func (def LBDefault) Validate() error {
	return v.Struct(&def,
		v.Field(&def.CPU),
		v.Field(&def.RAM),
		v.Field(&def.MainDiskSize),
	)
}

func (def *LBDefault) SetDefaults() {
	def.CPU = defaults.Default(def.CPU, defaultVCpu)
	def.RAM = defaults.Default(def.RAM, defaultRAM)
	def.MainDiskSize = defaults.Default(def.MainDiskSize, defaultMainDiskSize)
}

type LB struct {
	VIP             IPv4            `yaml:"vip,omitempty"`
	VirtualRouterId *Uint8          `yaml:"virtualRouterId,omitempty"`
	Default         LBDefault       `yaml:"default"`
	Instances       []LBInstance    `yaml:"instances,omitempty"`
	ForwardPorts    []LBPortForward `yaml:"forwardPorts,omitempty"`
}

func (lb LB) Validate() error {
	return v.Struct(&lb,
		v.Field(&lb.VIP,
			v.NotEmpty().When(len(lb.Instances) > 1).Error("Virtual IP (VIP) is required when multiple load balancer instances are configured."),
			v.OmitEmpty(),
			v.Custom(IP_IN_CIDR),
		),
		v.Field(&lb.VirtualRouterId),
		v.Field(&lb.Default),
		v.Field(&lb.Instances, v.UniqueField("Id")),
		v.Field(&lb.ForwardPorts),
	)
}

func (lb *LB) SetDefaults() {
	if len(lb.Instances) > 1 {
		lb.VirtualRouterId = defaults.Default(lb.VirtualRouterId, &defaultVRID)
	}

	for i := range lb.Instances {
		lb.Instances[i].CPU = defaults.Default(lb.Instances[i].CPU, lb.Default.CPU)
		lb.Instances[i].RAM = defaults.Default(lb.Instances[i].RAM, lb.Default.RAM)
		lb.Instances[i].MainDiskSize = defaults.Default(lb.Instances[i].MainDiskSize, lb.Default.MainDiskSize)

		if len(lb.Instances) > 1 {
			lb.Instances[i].Priority = defaults.Default(lb.Instances[i].Priority, &defaultPriority)
		}
	}
}

type LBPortForward struct {
	Name       string              `yaml:"name"`
	Port       Port                `yaml:"port"`
	TargetPort Port                `yaml:"targetPort,omitempty"`
	Target     LBPortForwardTarget `yaml:"target"`
}

func (pf LBPortForward) Validate() error {
	return v.Struct(&pf,
		v.Field(&pf.Name, v.NotEmpty(), v.AlphaNumericHypUS()),
		v.Field(&pf.Port, v.NotEmpty()),
		v.Field(&pf.TargetPort),
		v.Field(&pf.Target),
	)
}

func (pf *LBPortForward) SetDefaults() {
	pf.TargetPort = defaults.Default(pf.TargetPort, pf.Port)
	pf.Target = defaults.Default(pf.Target, WORKERS)
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
	Name         string `yaml:"name,omitempty" opt:"-"`
	Id           string `yaml:"id" opt:",id"`
	Host         string `yaml:"host,omitempty"`
	IP           IPv4   `yaml:"ip,omitempty"`
	MAC          MAC    `yaml:"mac,omitempty"`
	CPU          VCpu   `yaml:"cpu"`
	RAM          GB     `yaml:"ram"`
	MainDiskSize GB     `yaml:"mainDiskSize"`
	Priority     *Uint8 `yaml:"priority,omitempty"`
}

func (i LBInstance) GetTypeName() string {
	return "lb"
}

func (i LBInstance) GetID() string {
	return i.Id
}

func (i LBInstance) GetIP() IPv4 {
	return i.IP
}

func (i LBInstance) GetMAC() MAC {
	return i.MAC
}

func (i LBInstance) Validate() error {
	return v.Struct(&i,
		v.Field(&i.Id, v.NotEmpty(), v.AlphaNumericHypUS()),
		v.Field(&i.Host, v.OmitEmpty(), v.Custom(VALID_HOST)),
		v.Field(&i.IP, v.OmitEmpty(), v.Custom(IP_IN_CIDR)),
		v.Field(&i.MAC, v.OmitEmpty()),
		v.Field(&i.CPU),
		v.Field(&i.RAM),
		v.Field(&i.MainDiskSize),
		v.Field(&i.Priority),
	)
}
