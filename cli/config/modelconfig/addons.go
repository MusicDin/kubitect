package modelconfig

import v "cli/utils/validation"

type Addons struct {
	Kubespray *string `yaml:"kubespray" opt:"-"`
	Rook      *Rook   `yaml:"rook"`
}

func (a Addons) Validate() error {
	return v.Struct(&a,
		v.Field(&a.Rook),
	)
}

type Rook struct {
	Enabled      *bool    `yaml:"enabled"`
	NodeSelector *Labels  `yaml:"nodeSelector"`
	Version      *Version `yaml:"nodeSelector"`
}

func (r Rook) Validate() error {
	return v.Struct(&r,
		v.Field(&r.NodeSelector),
		v.Field(&r.Version),
	)
}
