package modelconfig

import (
	"fmt"

	v "github.com/MusicDin/kubitect/cli/utils/validation"

	"github.com/MusicDin/kubitect/cli/utils/defaults"
)

type Connection struct {
	User User           `yaml:"user,omitempty"`
	IP   IPv4           `yaml:"ip,omitempty"`
	Type ConnectionType `yaml:"type"`
	SSH  ConnectionSSH  `yaml:"ssh,omitempty"`
}

func (c Connection) Validate() error {
	isRemote := (c.Type == REMOTE)
	reqForRemoteErr := fmt.Sprintf("Field '{.Field}' is required when connection type is set to '%s'.", REMOTE)

	return v.Struct(&c,
		v.Field(&c.Type, v.NotEmpty()),
		v.Field(&c.IP, v.Skip().When(!isRemote), v.NotEmpty().Error(reqForRemoteErr)),
		v.Field(&c.User, v.Skip().When(!isRemote), v.NotEmpty().Error(reqForRemoteErr)),
		v.Field(&c.SSH, v.Skip().When(!isRemote), v.NotEmpty().Error(reqForRemoteErr)),
	)
}

type ConnectionType string

const (
	LOCAL     ConnectionType = "local"
	LOCALHOST ConnectionType = "localhost" // equivalent to local
	REMOTE    ConnectionType = "remote"
)

func (t ConnectionType) Validate() error {
	return v.Var(t, v.OneOf(LOCALHOST, LOCAL, REMOTE))
}

type ConnectionSSH struct {
	Keyfile File `yaml:"keyfile,omitempty"`
	Port    Port `yaml:"port,omitempty"`
	Verify  bool `yaml:"verify,omitempty"`
}

func (s ConnectionSSH) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Keyfile, v.NotEmpty().Error("Path to password-less private key of the remote host is required.")),
		v.Field(&s.Port),
	)
}

func (s *ConnectionSSH) SetDefaults() {
	s.Port = defaults.Default(s.Port, Port(22))
}
