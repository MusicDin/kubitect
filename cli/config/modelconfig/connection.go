package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Connection struct {
	Type *ConnectionType `yaml:"type,omitempty"`
	User *UserString     `yaml:"user,omitempty"`
	IP   *IP             `yaml:"ip,omitempty"`
	SSH  *SSH            `yaml:"ssh,omitempty"`
}

func (c Connection) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(c.IP, validation.Required.When(*c.Type == remote_connection)),
		validation.Field(c.SSH),
		validation.Field(c.Type, validation.Required),
		validation.Field(c.User, validation.Required.When(*c.Type == remote_connection)),
	)
}

type UserString string

func (s UserString) Validate() error {
	return validation.Validate(&s, StringNotEmptyAlphaNumeric...)
}
