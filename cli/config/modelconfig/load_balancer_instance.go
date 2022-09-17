package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type LoadBalancerInstance struct {
	CPU          *CpuSize    `yaml:"cpu"`
	Host         *HostName   `yaml:"host"`
	Id           *InstanceId `yaml:"id"`
	IP           *IP         `yaml:"ip"`
	MAC          *MAC        `yaml:"mac"`
	MainDiskSize *MB         `yaml:"mainDiskSize"`
	Priority     *Priority   `yaml:"priority"`
	RAM          *MB         `yaml:"ram"`
}

func (i LoadBalancerInstance) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(i.CPU),
		validation.Field(i.Id, validation.Required),
		validation.Field(i.Host), // TODO: Is valid Hostname?
		validation.Field(i.IP),   // TODO: Is withing CIDR?
		validation.Field(i.MainDiskSize),
		validation.Field(i.Priority),
		validation.Field(i.RAM),
	)
}

type Priority uint8

func (p Priority) Validate() error {
	return validation.Validate(&p, validation.Min(0), validation.Max(255))
}
