package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type Addons struct {
	Kubespray string `yaml:"kubespray,omitempty" opt:"-"`
	Rook      Rook   `yaml:"rook,omitempty"`
}

func (a Addons) Validate() error {
	return validation.Struct(&a,
		validation.Field(&a.Rook),
	)
}

type Rook struct {
	Enabled      bool    `yaml:"enabled"`
	Version      Version `yaml:"version"`
	NodeSelector Labels  `yaml:"nodeSelector"`
}

func (r Rook) Validate() error {
	return validation.Struct(&r,
		validation.Field(&r.Version, validation.OmitEmpty()),
		validation.Field(&r.NodeSelector, validation.OmitEmpty()),
	)
}
