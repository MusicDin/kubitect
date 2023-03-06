package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type Kubitect struct {
	Url     URL           `yaml:"url,omitempty"`
	Version MasterVersion `yaml:"version,omitempty"`
}

func (k Kubitect) Validate() error {
	return validation.Struct(&k,
		validation.Field(&k.Url, validation.OmitEmpty()),
		validation.Field(&k.Version, validation.OmitEmpty()),
	)
}
