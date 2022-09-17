package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SSH struct {
	Keyfile *SSHKeyPath `yaml:"keyfile,omitempty"`
	Port    *Port       `yaml:"port,omitempty"`
	Verify  *bool       `yaml:"verify,omitempty"`
}

func (s SSH) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(s.Keyfile, validation.By(PathExists)),
		validation.Field(s.Port),
	)
}

type SSHKeyPath string
