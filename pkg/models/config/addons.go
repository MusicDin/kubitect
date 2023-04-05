package config

import v "github.com/MusicDin/kubitect/pkg/utils/validation"

type Addons struct {
	Kubespray string `yaml:"kubespray,omitempty" opt:"-"`
	Rook      Rook   `yaml:"rook,omitempty"`
}

func (a Addons) Validate() error {
	return v.Struct(&a,
		v.Field(&a.Rook),
	)
}

type Rook struct {
	Enabled      bool    `yaml:"enabled"`
	Version      Version `yaml:"version"`
	NodeSelector Labels  `yaml:"nodeSelector"`
}

func (r Rook) Validate() error {
	return v.Struct(&r,
		v.Field(&r.Version, v.OmitEmpty()),
		v.Field(&r.NodeSelector, v.OmitEmpty()),
	)
}
