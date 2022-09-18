package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type LoadBalancerInstance struct {
	CPU          *CpuSize    `yaml:"cpu,omitempty"`
	Host         *HostName   `yaml:"host,omitempty"`
	Id           *InstanceId `yaml:"id,omitempty"`
	IP           *IP         `yaml:"ip,omitempty"`
	MAC          *MAC        `yaml:"mac,omitempty"`
	MainDiskSize *MB         `yaml:"mainDiskSize,omitempty"`
	Priority     *Priority   `yaml:"priority,omitempty"`
	RAM          *MB         `yaml:"ram,omitempty"`
}

func (i LoadBalancerInstance) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.CPU),
		validation.Field(&i.Id, validation.Required),
		validation.Field(&i.Host), // TODO: Is valid Hostname?
		validation.Field(&i.IP),   // TODO: Is withing CIDR?
		validation.Field(&i.MainDiskSize),
		validation.Field(&i.Priority),
		validation.Field(&i.RAM),
	)
}

type Priority uint8

func (p Priority) Validate() error {
	return validation.Validate(uint8(p), validation.Min(0), validation.Max(255))
}
