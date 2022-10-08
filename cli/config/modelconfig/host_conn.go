package modelconfig

import v "cli/validation"

type ConnectionType string

const (
	LOCALHOST ConnectionType = "localhost" // equivalent to local
	LOCAL     ConnectionType = "local"
	REMOTE    ConnectionType = "remote"
)

type Connection struct {
	IP   *IPv4           `yaml:"ip"`
	Type *ConnectionType `yaml:"type"`
	User *User           `yaml:"user"`
	SSH  *ConnectionSSH  `yaml:"ssh"`
}

func (c Connection) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Type, v.Required(), v.OneOf(LOCALHOST, LOCAL, REMOTE)),
		v.Field(&c.IP, v.Skip().When(*c.Type != REMOTE)),
		v.Field(&c.User, v.Skip().When(*c.Type != REMOTE)),
		v.Field(&c.SSH, v.Skip().When(*c.Type != REMOTE)),
	)
}
