package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type DataDisk struct {
	Name *DiskName `yaml:"name,omitempty"`
	Pool *PoolName `yaml:"pool,omitempty"`
	Size *DiskSize `yaml:"size,omitempty"`
}

func (d DataDisk) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name),
		validation.Field(&d.Pool),
		validation.Field(&d.Size),
	)
}

type DiskName string

func (n DiskName) Validate() error {
	return validation.Validate(&n, StringNotEmptyAlphaNumeric...)
}

type PoolName string

func (p PoolName) Validate() error {
	return validation.Validate(&p) // TODO: IsValidPoolName
}

type DiskSize uint

func (d DiskSize) Validate() error {
	return validation.Validate(&d, validation.Min(1))
}
