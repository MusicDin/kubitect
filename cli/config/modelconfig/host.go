package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Host struct {
	Connection           *Connection        `yaml:"connection"`
	DataResourcePools    []DataResourcePool `yaml:"dataResourcePools"`
	Default              bool               `yaml:"default"`
	Name                 *HostName          `yaml:"name"`
	MainResourcePoolPath *ResourcePath      `yaml:"mainResourcePoolPath"`
}

func (h Host) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(h.Name, validation.Required),
		validation.Field(h.Connection),
		validation.Field(&h.DataResourcePools),
		validation.Field(h.MainResourcePoolPath),
	)
}
