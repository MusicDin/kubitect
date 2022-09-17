package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Network struct {
	Bridge  *Bridge      `yaml:"bridge"`
	CIDR    *CIDR        `yaml:"cidr"`
	Gateway *Gateway     `yaml:"gateway"`
	Mode    *NetworkMode `yaml:"mode"`
}

func (n Network) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(n.CIDR, validation.Required),
		validation.Field(n.Bridge),
		validation.Field(n.Gateway),
		validation.Field(n.Mode, validation.Required),
	)
}

type Gateway string
type CIDR string

func (c CIDR) Validate() error {
	return validation.Validate(&c) // TODO: check CIDR
}

type Bridge string

func (b Bridge) Validate() error {
	return validation.Validate(&b, StringNotEmptyAlphaNumeric...)
}
