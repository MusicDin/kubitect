package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type WorkerDefault struct {
	CPU          *CpuSize            `yaml:"cpu,omitempty"  default:"2"`
	Labels       *map[LabelKey]Label `yaml:"labels,omitempty"`
	MainDiskSize *GB                 `yaml:"mainDiskSize,omitempty"  default:"32"`
	RAM          *GB                 `yaml:"ram,omitempty" default:"4"`
	Taints       *[]Taint            `yaml:"taints,omitempty"`
}

func (d WorkerDefault) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Labels),
		validation.Field(&d.Taints),
		validation.Field(&d.RAM),
		validation.Field(&d.MainDiskSize),
		validation.Field(&d.CPU),
	)
}

type LabelKey string // TODO: Check if correct type
type Label string    // TODO: Check if correct type
type Taint string    // TODO: Check if correct type
