package modelconfig

import (
	"cli/utils/defaults"
	v "cli/utils/validation"
)

type ConnectionSSH struct {
	Keyfile File `yaml:"keyfile"`
	Port    Port `yaml:"port"`
	Verify  bool `yaml:"verify"`
}

func (s ConnectionSSH) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Keyfile, v.NotEmpty().Error("Path to password-less private key of the remote host is required.")),
		v.Field(&s.Port),
	)
}

func (s *ConnectionSSH) SetDefaults() {
	s.Port = defaults.Default(s.Port, Port(22))
}
