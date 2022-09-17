package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ForwardPort struct {
	Name       *ForwardPortName   `yaml:"name,omitempty"`
	Port       *Port              `yaml:"port,omitempty"`
	TargetPort *Port              `yaml:"targetPort,omitempty"`
	Target     *PortForwardTarget `yaml:"target,omitempty"`
}

func (f ForwardPort) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(f.Name),
		validation.Field(f.Port),
		validation.Field(f.TargetPort, validation.Required),
		validation.Field(f.Target),
	)
}

type ForwardPortName string

func (f ForwardPortName) Validate() error {
	return validation.Validate(&f, StringNotEmptyAlphaNumeric...)
}
