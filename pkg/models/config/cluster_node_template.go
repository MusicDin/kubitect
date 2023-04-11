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
	CENTOS9  OSDistro = "centos9"
	ROCKY    OSDistro = "rocky"
)

func (d OSDistro) Validate() error {
	return v.Var(d, v.OneOf(UBUNTU, UBUNTU20, UBUNTU22, DEBIAN, DEBIAN11, CENTOS9, ROCKY))
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
