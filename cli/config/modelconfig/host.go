package modelconfig

import (
	"cli/utils/defaults"
	v "cli/utils/validation"
)

type Host struct {
	Name                 string             `yaml:"name" opt:",id"`
	Default              bool               `yaml:"default"`
	Connection           Connection         `yaml:"connection"`
	MainResourcePoolPath string             `yaml:"mainResourcePoolPath"`
	DataResourcePools    []DataResourcePool `yaml:"dataResourcePools"`
}

func (h Host) Validate() error {
	return v.Struct(&h,
		v.Field(&h.Name, v.NotEmpty(), v.AlphaNumericHypUS()),
		v.Field(&h.Connection),
		v.Field(&h.MainResourcePoolPath), // TODO: validate dir path which does not have to exist
		v.Field(&h.DataResourcePools, v.UniqueField("Name")),
	)
}

func (h *Host) SetDefaults() {
	h.MainResourcePoolPath = defaults.Default(h.MainResourcePoolPath, "/var/lib/libvirt/images/")
}
