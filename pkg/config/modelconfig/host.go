package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/pkg/utils/validation"
)

const (
	defaultResPoolPath = "/var/lib/libvirt/images/"
)

type Host struct {
	Name                 string             `yaml:"name" opt:",id"`
	Default              bool               `yaml:"default"`
	Connection           Connection         `yaml:"connection,omitempty"`
	MainResourcePoolPath string             `yaml:"mainResourcePoolPath"`
	DataResourcePools    []DataResourcePool `yaml:"dataResourcePools,omitempty"`
}

func (h Host) Validate() error {
	return validation.Struct(&h,
		validation.Field(&h.Name, validation.NotEmpty(), validation.AlphaNumericHypUS()),
		validation.Field(&h.Connection),
		validation.Field(&h.MainResourcePoolPath), // TODO: validate dir path which does not have to exist
		validation.Field(&h.DataResourcePools, validation.UniqueField("Name")),
	)
}

func (h *Host) SetDefaults() {
	h.MainResourcePoolPath = defaults.Default(h.MainResourcePoolPath, defaultResPoolPath)
}

type DataResourcePool struct {
	Name string `yaml:"name" opt:",id"`
	Path string `yaml:"path"`
}

func (rp DataResourcePool) Validate() error {
	return validation.Struct(&rp,
		validation.Field(&rp.Name, validation.NotEmpty(), validation.AlphaNumericHyp()),
		validation.Field(&rp.Path, validation.NotEmpty()), // TODO: Valid file path. File does not need to exist.
	)
}

func (rp *DataResourcePool) SetDefaults() {
	rp.Path = defaults.Default(rp.Path, defaultResPoolPath)
}
