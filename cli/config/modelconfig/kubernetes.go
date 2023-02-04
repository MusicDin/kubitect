package modelconfig

import (
	v "cli/utils/validation"

	"cli/utils/defaults"
)

type Kubernetes struct {
	Version       Version       `yaml:"version"`
	DnsMode       DnsMode       `yaml:"dnsMode"`
	NetworkPlugin NetworkPlugin `yaml:"networkPlugin"`
	Kubespray     Kubespray     `yaml:"kubespray,omitempty"`
	Other         Other         `yaml:"other"`
}

func (k Kubernetes) Validate() error {
	return v.Struct(&k,
		v.Field(&k.Version, v.NotEmpty()),
		v.Field(&k.DnsMode, v.NotEmpty()),
		v.Field(&k.NetworkPlugin, v.NotEmpty()),
		v.Field(&k.Kubespray, v.NotEmpty()),
		v.Field(&k.Other),
	)
}

func (k *Kubernetes) SetDefaults() {
	// TODO: Set default k8s version in env
	k.Version = defaults.Default(k.Version, "v1.23.0")
	k.DnsMode = defaults.Default(k.DnsMode, COREDNS)
	k.NetworkPlugin = defaults.Default(k.NetworkPlugin, CALICO)
}

type DnsMode string

const (
	COREDNS DnsMode = "coredns"
	KUBEDNS DnsMode = "kubedns"
)

func (m DnsMode) Validate() error {
	return v.Var(m, v.OneOf(COREDNS, KUBEDNS))
}

type NetworkPlugin string

const (
	CALICO      NetworkPlugin = "calico"
	CILIUM      NetworkPlugin = "cilium"
	CANAL       NetworkPlugin = "canal"
	FLANNEL     NetworkPlugin = "flannel"
	WEAVE       NetworkPlugin = "weave"
	KUBE_ROUTER NetworkPlugin = "kube-router"
)

func (p NetworkPlugin) Validate() error {
	return v.Var(p, v.OneOf(CALICO, CILIUM, CANAL, FLANNEL, WEAVE, KUBE_ROUTER))
}

type Other struct {
	AutoRenewCertificates bool `yaml:"autoRenewCertificates"`
	CopyKubeconfig        bool `yaml:"copyKubeconfig"`
}

type Kubespray struct {
	URL     URL           `yaml:"url,omitempty"`
	Version MasterVersion `yaml:"version,omitempty"`
}

func (k Kubespray) Validate() error {
	return v.Struct(&k,
		v.Field(&k.URL, v.OmitEmpty()),
		v.Field(&k.Version),
	)
}
