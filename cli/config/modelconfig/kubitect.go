package modelconfig

import v "github.com/MusicDin/kubitect/cli/utils/validation"

type Kubitect struct {
	Url     URL           `yaml:"url,omitempty"`
	Version MasterVersion `yaml:"version,omitempty"`
}

func (k Kubitect) Validate() error {
	return v.Struct(&k,
		v.Field(&k.Url, v.OmitEmpty()),
		v.Field(&k.Version, v.OmitEmpty()),
	)
}
