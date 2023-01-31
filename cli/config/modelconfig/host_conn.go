package modelconfig

import (
	v "cli/utils/validation"
	"fmt"
)

type ConnectionType string

const (
	LOCALHOST ConnectionType = "localhost" // equivalent to local
	LOCAL     ConnectionType = "local"
	REMOTE    ConnectionType = "remote"
)

func (t ConnectionType) Validate() error {
	return v.Var(t, v.OneOf(LOCALHOST, LOCAL, REMOTE))
}

type Connection struct {
	IP   *IPv4           `yaml:"ip"`
	Type *ConnectionType `yaml:"type"`
	User *User           `yaml:"user"`
	SSH  ConnectionSSH   `yaml:"ssh"`
}

func (c Connection) Validate() error {
	isRemote := (c.Type != nil && *c.Type == REMOTE)
	isRemoteErr := fmt.Sprintf("Field '{.Field}' is required when connection type is set to '%s'.", REMOTE)

	return v.Struct(&c,
		v.Field(&c.Type, v.Required()),
		v.Field(&c.IP, v.Required().When(isRemote).Error(isRemoteErr)),
		v.Field(&c.User, v.Required().When(isRemote).Error(isRemoteErr)),
		v.Field(&c.SSH,
			v.Skip().When(!isRemote),
			v.NotEmpty().When(isRemote).Error(isRemoteErr),
		),
	)
}
