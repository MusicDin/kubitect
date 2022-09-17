package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Instance struct {
	CPU          *CpuSize           `yaml:"cpu"`
	Host         *HostName          `yaml:"host"`
	Id           *InstanceId        `yaml:"id"`
	IP           *IP                `yaml:"ip"`
	MAC          *MAC               `yaml:"mac"`
	Labels       map[LabelKey]Label `yaml:"labels"`
	MainDiskSize *MB                `yaml:"mainDiskSize"`
	RAM          *MB                `yaml:"ram"`
	Taints       []Taint            `yaml:"taints"`
	DataDisks    []DataDisk         `yaml:"dataDisks"`
}

func (i Instance) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(i.CPU),
		validation.Field(i.Host), // TODO: Is valid Host?
		validation.Field(i.Id, validation.Required),
		validation.Field(i.IP), // TODO: Is within CIDR?
		validation.Field(&i.Labels),
		validation.Field(&i.Taints),
		validation.Field(i.MAC),
		validation.Field(i.MainDiskSize),
		validation.Field(i.RAM),
		validation.Field(&i.DataDisks),
	)
}
