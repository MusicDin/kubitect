package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Instance struct {
	CPU          *CpuSize
	Host         *HostName
	Id           *InstanceId
	IP           *IP
	MAC          *MAC
	Labels       map[LabelKey]Label
	MainDiskSize *MB
	RAM          *MB
	Taints       []Taint
	DataDisks    []DataDisk
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
