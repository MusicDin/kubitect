package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Instance struct {
	Id           *InstanceId         `yaml:"id,omitempty"`
	IP           *IP                 `yaml:"ip,omitempty"`
	MAC          *MAC                `yaml:"mac,omitempty"`
	CPU          *CpuSize            `yaml:"cpu,omitempty"`
	RAM          *MB                 `yaml:"ram,omitempty"`
	Host         *HostName           `yaml:"host,omitempty"`
	MainDiskSize *MB                 `yaml:"mainDiskSize,omitempty"`
	Labels       *map[LabelKey]Label `yaml:"labels,omitempty"`
	Taints       *[]Taint            `yaml:"taints,omitempty"`
	DataDisks    *[]DataDisk         `yaml:"dataDisks,omitempty"`
}

func (i Instance) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(i.CPU),
		validation.Field(i.Host), // TODO: Is valid Host?
		validation.Field(i.Id, validation.Required),
		validation.Field(i.IP), // TODO: Is within CIDR?
		validation.Field(i.Labels),
		validation.Field(i.Taints),
		validation.Field(i.MAC),
		validation.Field(i.MainDiskSize),
		validation.Field(i.RAM),
		validation.Field(i.DataDisks),
	)
}
