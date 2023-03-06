package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/env"
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type NodeTemplate struct {
	User         User            `yaml:"user"`
	OS           OS              `yaml:"os"`
	SSH          NodeTemplateSSH `yaml:"ssh,omitempty"`
	CpuMode      CpuMode         `yaml:"cpuMode,omitempty"`
	DNS          []IP            `yaml:"dns,omitempty"`
	UpdateOnBoot bool            `yaml:"updateOnBoot"`
}

func (n NodeTemplate) Validate() error {
	return validation.Struct(&n,
		validation.Field(&n.User),
		validation.Field(&n.OS),
		validation.Field(&n.DNS),
		validation.Field(&n.SSH),
	)
}

func (n *NodeTemplate) SetDefaults() {
	n.User = defaults.Default(n.User, "k8s")
	n.CpuMode = defaults.Default(n.CpuMode, CUSTOM)
	n.UpdateOnBoot = defaults.Default(n.UpdateOnBoot, true)
}

type OS struct {
	Distro           OSDistro           `yaml:"distro"`
	NetworkInterface OSNetworkInterface `yaml:"networkInterface"`
	Source           OSSource           `yaml:"source"`
}

func (s OS) Validate() error {
	return validation.Struct(&s,
		validation.Field(&s.Distro),
		validation.Field(&s.NetworkInterface),
		validation.Field(&s.Source),
	)
}

func (s *OS) SetDefaults() {
	s.NetworkInterface = defaults.Default(s.NetworkInterface, OSNetworkInterface("ens3"))
	s.Distro = defaults.Default(s.Distro, UBUNTU)
	s.Source = defaults.Default(s.Source, OSSource(env.ProjectOsPresets[string(s.Distro)]))
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
	return validation.Var(d, validation.OneOf(UBUNTU, UBUNTU20, UBUNTU22, DEBIAN, DEBIAN11))
}

type OSNetworkInterface string

func (nic OSNetworkInterface) Validate() error {
	return validation.Var(nic, validation.AlphaNumeric(), validation.MaxLen(16))
}

type OSSource string

func (os OSSource) Validate() error {
	return validation.Var(os) // TODO: URL or FileExists
}

type NodeTemplateSSH struct {
	AddToKnownHosts bool `yaml:"addToKnownHosts"`
	PrivateKeyPath  File `yaml:"privateKeyPath"`
}

func (ssh NodeTemplateSSH) Validate() error {
	return validation.Struct(&ssh)
	// v.Field(&ssh.PrivateKeyPath, v.Skip()),
}

func (ssh *NodeTemplateSSH) SetDefaults() {
	ssh.AddToKnownHosts = defaults.Default(ssh.AddToKnownHosts, true)
}

type CpuMode string

const (
	CUSTOM      CpuMode = "custom"
	PASSTHROUGH CpuMode = "host-passthrough"
)

func (m CpuMode) Validate() error {
	return validation.Var(m, validation.OneOf(PASSTHROUGH, CUSTOM))
}
