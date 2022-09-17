package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LoadBalancerDefault struct {
	CPU          *CpuSize `yaml:"cpu"`
	MainDiskSize *MB      `yaml:"mainDiskSize"`
	RAM          *MB      `yaml:"ram"`
}

func (l LoadBalancerDefault) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(l.CPU, validation.Required),
		validation.Field(l.MainDiskSize, validation.Required),
		validation.Field(l.RAM, validation.Required),
	)
}

type LoadBalancer struct {
	VIP             *IP                  `yaml:"vip"`
	VirtualRouterId *LoadBalancerId      `yaml:"virtualRouterId"`
	Default         *LoadBalancerDefault `yaml:"default"`

	ForwardPorts []ForwardPort          `yaml:"forwardPorts"`
	Instances    []LoadBalancerInstance `yaml:"instances"`
}

type LoadBalancerId uint8

func (i LoadBalancerId) Validate() error {
	return validation.Validate(&i, validation.Min(0), validation.Max(255))
}

func (b LoadBalancer) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(b.Default),
		validation.Field(b.Instances),
		validation.Field(b.VIP, validation.Required.When(len(b.Instances) > 0)),
	)
}
