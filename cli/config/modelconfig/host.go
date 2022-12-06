package modelconfig

import v "cli/lib/validation"

type Host struct {
	Name                 *string            `yaml:"name" opt:",id"`
	Default              *bool              `yaml:"default"`
	Connection           Connection         `yaml:"connection"`
	MainResourcePoolPath *string            `yaml:"mainResourcePoolPath"`
	DataResourcePools    []DataResourcePool `yaml:"dataResourcePools"`
}

func (h Host) Validate() error {
	return v.Struct(&h,
		v.Field(&h.Name, v.Required(), v.AlphaNumericHypUS()),
		v.Field(&h.Connection),
		v.Field(&h.MainResourcePoolPath), // TODO: validate dir path which does not have to exist
		v.Field(&h.DataResourcePools, v.UniqueField("Name")),
	)
}
