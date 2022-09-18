package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type HostSSH struct {
	Keyfile *SSHKeyPath `yaml:"keyfile,omitempty" default:"~/.ssh/id_rsa"`
	Port    *Port       `yaml:"port,omitempty" default:"22"`
	Verify  *bool       `yaml:"verify,omitempty"  default:"false"`
}

func (s HostSSH) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Keyfile, validation.By(PathExists)),
		validation.Field(&s.Port),
	)
}

type SSHKeyPath string
