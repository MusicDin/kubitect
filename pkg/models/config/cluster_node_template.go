package config

import (
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
)

type NodeTemplate struct {
	User         User            `yaml:"user"`
	OS           OS              `yaml:"os"`
	SSH          NodeTemplateSSH `yaml:"ssh"`
	CpuMode      CpuMode         `yaml:"cpuMode,omitempty"`
	DNS          []IP            `yaml:"dns,omitempty"`
	UpdateOnBoot *bool           `yaml:"updateOnBoot"`
}

func (n NodeTemplate) Validate() error {
	return v.Struct(&n,
		v.Field(&n.User),
		v.Field(&n.OS),
		v.Field(&n.SSH),
		v.Field(&n.CpuMode),
		v.Field(&n.DNS),
	)
}

func (n *NodeTemplate) SetDefaults() {
	def := true

	n.User = defaults.Default(n.User, "k8s")
	n.CpuMode = defaults.Default(n.CpuMode, CUSTOM)
	n.UpdateOnBoot = defaults.Default(n.UpdateOnBoot, &def)
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
	s.Distro = defaults.Default(s.Distro, UBUNTU22)

	preset := env.ProjectOsPresets[string(s.Distro)]
	s.NetworkInterface = defaults.Default(s.NetworkInterface, OSNetworkInterface(preset.NetworkInterface))
	s.Source = defaults.Default(s.Source, OSSource(preset.Source))
}

type OSDistro string

const (
	UBUNTU22 OSDistro = "ubuntu22"
	UBUNTU20 OSDistro = "ubuntu20"
	DEBIAN11 OSDistro = "debian11"
	DEBIAN12 OSDistro = "debian12"
	CENTOS9  OSDistro = "centos9"
	ROCKY9   OSDistro = "rocky9"
)

func (d OSDistro) Validate() error {
	return v.Var(d, v.OneOf(UBUNTU20, UBUNTU22, DEBIAN11, DEBIAN12, CENTOS9, ROCKY9))
}

type OSNetworkInterface string

func (nic OSNetworkInterface) Validate() error {
	return v.Var(nic, v.AlphaNumeric(), v.MaxLen(16))
}

type OSSource string

func (os OSSource) Validate() error {
	return v.Var(os)
}

type NodeTemplateSSH struct {
	AddToKnownHosts bool `yaml:"addToKnownHosts"`
	PrivateKeyPath  File `yaml:"privateKeyPath,omitempty"`
}

func (ssh NodeTemplateSSH) Validate() error {
	return v.Struct(&ssh)
	// v.Field(&ssh.PrivateKeyPath, v.Skip()),
}

type CpuMode string

const (
	CUSTOM           CpuMode = "custom"
	HOST_MODEL       CpuMode = "host-model"
	HOST_PASSTHROUGH CpuMode = "host-passthrough"
	MAXIMUM          CpuMode = "maximum"
)

func (m CpuMode) Validate() error {
	return v.Var(m, v.OneOf(CUSTOM, HOST_MODEL, HOST_PASSTHROUGH, MAXIMUM))
}
