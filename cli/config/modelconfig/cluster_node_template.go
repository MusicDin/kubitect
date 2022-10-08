package modelconfig

import v "cli/validation"

type NodeTemplate struct {
	User         *User            `yaml:"user"`
	OS           *OS              `yaml:"os"`
	SSH          *NodeTemplateSSH `yaml:"ssh,omitempty"`
	DNS          *[]IP            `yaml:"dns"`
	UpdateOnBoot *bool            `yaml:"updateOnBoot"`
}

func (n NodeTemplate) Validate() error {
	return v.Struct(&n,
		v.Field(&n.User, v.OmitEmpty()),
		v.Field(&n.DNS, v.OmitEmpty()),
		v.Field(&n.SSH),
	)
}

type OS struct {
	Distro           *OSDistro           `yaml:"distro"`
	NetworkInterface *OSNetworkInterface `yaml:"networkInterface"`
	Source           *OSSource           `yaml:"source"`
}

func (s OS) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Distro),
		v.Field(&s.NetworkInterface), // TODO: depends on Distro
		v.Field(&s.Source),           // TODO: depends on Distro
	)
}

type OSDistro string

const (
	UBUNTU   OSDistro = "ubuntu"
	UBUNTU20 OSDistro = "ubuntu20"
	UBUNTU22 OSDistro = "ubuntu22"
	DEBIAN   OSDistro = "debian"
	DEBIAN11 OSDistro = "debian11"
)

func (distro OSDistro) Validate() error {
	return v.Var(distro, v.OmitEmpty(), v.OneOf(UBUNTU, UBUNTU20, UBUNTU22, DEBIAN, DEBIAN11))
}

type OSNetworkInterface string

func (nic OSNetworkInterface) Validate() error {
	return v.Var(nic, v.OmitEmpty(), v.AlphaNumeric())
}

type OSSource string

func (s OSSource) Validate() error {
	return v.Var(s, v.OmitEmpty()) // TODO: URL or FileExists
}

type NodeTemplateSSH struct {
	AddToKnownHosts *bool `yaml:"addToKnownHosts"`
	PrivateKeyPath  *File `yaml:"privateKeyPath"`
}

func (s NodeTemplateSSH) Validate() error {
	return v.Struct(&s,
		v.Field(&s.PrivateKeyPath, v.OmitEmpty()),
	)
}
