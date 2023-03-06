package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
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
	return v.Struct(&h,
		v.Field(&h.Name, v.NotEmpty(), v.AlphaNumericHypUS()),
		v.Field(&h.Connection),
		v.Field(&h.MainResourcePoolPath), // TODO: validate dir path which does not have to exist
		v.Field(&h.DataResourcePools, v.UniqueField("Name")),
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
	return v.Struct(&rp,
		v.Field(&rp.Name, v.NotEmpty(), v.AlphaNumericHyp()),
		v.Field(&rp.Path, v.NotEmpty()), // TODO: Valid file path. File does not need to exist.
	)
}

func (rp *DataResourcePool) SetDefaults() {
	rp.Path = defaults.Default(rp.Path, defaultResPoolPath)
}
