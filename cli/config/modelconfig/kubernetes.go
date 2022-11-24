package modelconfig

import v "cli/lib/validation"

type Kubernetes struct {
	Version       *Version       `yaml:"version"`
	DnsMode       *DnsMode       `yaml:"dnsMode"`
	NetworkPlugin *NetworkPlugin `yaml:"networkPlugin"`
	Kubespray     Kubespray      `yaml:"kubespray"`
	Other         Other          `yaml:"other"`
}

func (k Kubernetes) Validate() error {
	return v.Struct(&k,
		v.Field(&k.Version, v.Required()),
		v.Field(&k.DnsMode),
		v.Field(&k.NetworkPlugin),
		v.Field(&k.Kubespray, v.NotEmpty()),
		v.Field(&k.Other),
	)
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
	AutoRenewCertificates *bool `yaml:"autoRenewCertificates"`
	CopyKubeconfig        *bool `yaml:"copyKubeconfig"`
}

type Kubespray struct {
	URL     *URL           `yaml:"url"`
	Version *MasterVersion `yaml:"version"`
}

func (k Kubespray) Validate() error {
	return v.Struct(&k,
		v.Field(&k.URL),
		v.Field(&k.Version, v.Required()),
	)
}
