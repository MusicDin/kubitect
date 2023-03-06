package modelconfig

import (
	"fmt"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/pkg/utils/validation"
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

	return validation.Struct(&c,
		validation.Field(&c.Type, validation.NotEmpty()),
		validation.Field(&c.IP, validation.Skip().When(!isRemote), validation.NotEmpty().Error(reqForRemoteErr)),
		validation.Field(&c.User, validation.Skip().When(!isRemote), validation.NotEmpty().Error(reqForRemoteErr)),
		validation.Field(&c.SSH, validation.Skip().When(!isRemote), validation.NotEmpty().Error(reqForRemoteErr)),
	)
}

type ConnectionType string

const (
	LOCAL     ConnectionType = "local"
	LOCALHOST ConnectionType = "localhost" // equivalent to local
	REMOTE    ConnectionType = "remote"
)

func (t ConnectionType) Validate() error {
	return validation.Var(t, validation.OneOf(LOCALHOST, LOCAL, REMOTE))
}

type ConnectionSSH struct {
	Keyfile File `yaml:"keyfile,omitempty"`
	Port    Port `yaml:"port,omitempty"`
	Verify  bool `yaml:"verify,omitempty"`
}

func (s ConnectionSSH) Validate() error {
	return validation.Struct(&s,
		validation.Field(&s.Keyfile, validation.NotEmpty().Error("Path to password-less private key of the remote host is required.")),
		validation.Field(&s.Port),
	)
}

func (s *ConnectionSSH) SetDefaults() {
	s.Port = defaults.Default(s.Port, Port(22))
}
