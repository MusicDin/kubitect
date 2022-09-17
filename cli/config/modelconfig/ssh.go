package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SSH struct {
	Keyfile *SSHKeyPath `yaml:"keyfile"`
	Port    *Port       `yaml:"port"`
	Verify  bool        `yaml:"verify"`
}

func (s SSH) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(s.Keyfile, validation.By(PathExists)),
		validation.Field(s.Port),
	)
}

type SSHKeyPath string
