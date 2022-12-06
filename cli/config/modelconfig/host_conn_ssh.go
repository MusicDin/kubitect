package modelconfig

import v "cli/lib/validation"

type ConnectionSSH struct {
	Keyfile *File `yaml:"keyfile"`
	Port    *Port `yaml:"port"`
	Verify  *bool `yaml:"verify"`
}

func (s ConnectionSSH) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Keyfile, v.Required().Error("Path to password-less private key of the remote host is required.")),
		v.Field(&s.Port),
	)
}
