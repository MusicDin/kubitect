package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ForwardPort struct {
	Name       *ForwardPortName   `yaml:"name,omitempty"`
	Port       *Port              `yaml:"port,omitempty"`
	TargetPort *Port              `yaml:"targetPort,omitempty"`
	Target     *PortForwardTarget `yaml:"target,omitempty" default:"workers"`
}

func (f *ForwardPort) SetDefaults() {
	// Defaults to the incoming port value.
	if f.TargetPort == nil {
		f.TargetPort = f.Port
	}
}

func (f ForwardPort) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Port, validation.Required),
		validation.Field(&f.TargetPort),
		validation.Field(&f.Target),
	)
}

type ForwardPortName string

func (f ForwardPortName) Validate() error {
	return validation.Validate(string(f), StringNotEmptyAlphaNumericMinus...)
}
