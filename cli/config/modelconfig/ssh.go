package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SSH struct {
	Keyfile *SSHKeyPath
	Port    *Port
	Verify  bool
}

func (s SSH) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(s.Keyfile, validation.By(PathExists)),
		validation.Field(s.Port),
	)
}

type SSHKeyPath string
