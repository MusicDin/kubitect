package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Host struct {
	Name                 *HostName           `yaml:"name,omitempty"`
	Default              *bool               `yaml:"default,omitempty"  default:"false"`
	Connection           *Connection         `yaml:"connection,omitempty"`
	MainResourcePoolPath *ResourcePath       `yaml:"mainResourcePoolPath,omitempty"`
	DataResourcePools    *[]DataResourcePool `yaml:"dataResourcePools,omitempty"`
}

func (h Host) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Default),
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.Connection),
		validation.Field(&h.DataResourcePools),
		validation.Field(&h.MainResourcePoolPath),
	)
}
