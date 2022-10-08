package modelconfig

import v "cli/validation"

type ConnectionSSH struct {
	Keyfile *File `yaml:"keyfile"`
	Port    *Port `yaml:"port"`
	Verify  *bool `yaml:"verify"`
}

func (s ConnectionSSH) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Keyfile, v.OmitEmpty()),
		v.Field(&s.Port, v.OmitEmpty()),
	)
}
