package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LoadBalancerDefault struct {
	CPU          *CpuSize `yaml:"cpu,omitempty"  default:"2"`
	MainDiskSize *GB      `yaml:"mainDiskSize,omitempty"  default:"32"`
	RAM          *GB      `yaml:"ram,omitempty"  default:"4"`
}

func (l LoadBalancerDefault) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.CPU, validation.Required),
		validation.Field(&l.MainDiskSize, validation.Required),
		validation.Field(&l.RAM, validation.Required),
	)
}

type LoadBalancer struct {
	VIP             *IP                  `yaml:"vip,omitempty"`
	VirtualRouterId *LoadBalancerId      `yaml:"virtualRouterId,omitempty" default:"51"`
	Default         *LoadBalancerDefault `yaml:"default,omitempty"`

	ForwardPorts *[]ForwardPort          `yaml:"forwardPorts,omitempty"`
	Instances    *[]LoadBalancerInstance `yaml:"instances,omitempty"`
}

func (b LoadBalancer) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Default),
		validation.Field(&b.VirtualRouterId),
		validation.Field(&b.Instances),
		validation.Field(&b.ForwardPorts),
		validation.Field(&b.VIP, validation.Required.When(len(*b.Instances) > 0)),
	)
}

type LoadBalancerId uint8

func (i LoadBalancerId) Validate() error {
	return validation.Validate(int(i), validation.Min(0), validation.Max(255))
}
