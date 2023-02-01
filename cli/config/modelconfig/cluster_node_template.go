package modelconfig

import (
	"cli/utils/defaults"
	v "cli/utils/validation"
)

type NodeTemplate struct {
	User         User            `yaml:"user"`
	OS           OS              `yaml:"os"`
	SSH          NodeTemplateSSH `yaml:"ssh"`
	DNS          []IP            `yaml:"dns"`
	UpdateOnBoot bool            `yaml:"updateOnBoot"`
}

func (n NodeTemplate) Validate() error {
	return v.Struct(&n,
		v.Field(&n.User),
		v.Field(&n.OS),
		v.Field(&n.DNS),
		v.Field(&n.SSH),
	)
}

func (n *NodeTemplate) SetDefaults() {
	// TODO: Save these default values in env
	n.User = defaults.Default(n.User, "k8s")
	n.UpdateOnBoot = defaults.Default(n.UpdateOnBoot, true)
}

type OS struct {
	Distro           OSDistro           `yaml:"distro"`
	NetworkInterface OSNetworkInterface `yaml:"networkInterface"`
	Source           OSSource           `yaml:"source"`
}

func (s OS) Validate() error {
	return v.Struct(&s,
		v.Field(&s.Distro),
		v.Field(&s.NetworkInterface),
		v.Field(&s.Source),
	)
}

func (s *OS) SetDefaults() {
	s.Distro = defaults.Default(s.Distro, UBUNTU)

	// TODO: Evaluate as in Terraform base module
	s.NetworkInterface = defaults.Default(s.NetworkInterface, OSNetworkInterface("ens3"))
	// s.Source = defaults.Default(s.NetworkInterface, OSNetworkInterface("ens3"))
}

type OSDistro string

const (
	UBUNTU   OSDistro = "ubuntu"
	UBUNTU20 OSDistro = "ubuntu20"
	UBUNTU22 OSDistro = "ubuntu22"
	DEBIAN   OSDistro = "debian"
	DEBIAN11 OSDistro = "debian11"
)

func (d OSDistro) Validate() error {
	return v.Var(d, v.OneOf(UBUNTU, UBUNTU20, UBUNTU22, DEBIAN, DEBIAN11))
}

type OSNetworkInterface string

func (nic OSNetworkInterface) Validate() error {
	return v.Var(nic, v.AlphaNumeric(), v.MaxLen(16))
}

type OSSource string

func (os OSSource) Validate() error {
	return v.Var(os) // TODO: URL or FileExists
}

type NodeTemplateSSH struct {
	AddToKnownHosts bool `yaml:"addToKnownHosts"`
	PrivateKeyPath  File `yaml:"privateKeyPath"`
}

func (ssh NodeTemplateSSH) Validate() error {
	return v.Struct(&ssh,
		v.Field(&ssh.PrivateKeyPath, v.OmitEmpty()),
	)
}

func (ssh *NodeTemplateSSH) SetDefaults() {
	ssh.AddToKnownHosts = defaults.Default(ssh.AddToKnownHosts, true)
}
