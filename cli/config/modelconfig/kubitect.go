package modelconfig

import v "cli/lib/validation"

type Kubitect struct {
	Url     *URL           `yaml:"url"`
	Version *MasterVersion `yaml:"version"`
}

func (k Kubitect) Validate() error {
	return v.Struct(&k,
		v.Field(&k.Url),
		v.Field(&k.Version),
	)
}
