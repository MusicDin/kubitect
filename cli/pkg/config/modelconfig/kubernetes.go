package modelconfig

import (
	"fmt"
	"github.com/MusicDin/kubitect/cli/pkg/env"
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type Kubernetes struct {
	Version       KubernetesVersion `yaml:"version"`
	DnsMode       DnsMode           `yaml:"dnsMode"`
	NetworkPlugin NetworkPlugin     `yaml:"networkPlugin"`
	Other         Other             `yaml:"other"`
}

func (k Kubernetes) Validate() error {
	return validation.Struct(&k,
		validation.Field(&k.Version, validation.NotEmpty()),
		validation.Field(&k.DnsMode, validation.NotEmpty()),
		validation.Field(&k.NetworkPlugin, validation.NotEmpty()),
		validation.Field(&k.Other),
	)
}

func (k *Kubernetes) SetDefaults() {
	k.Version = defaults.Default(k.Version, env.ConstKubernetesVersion)
	k.DnsMode = defaults.Default(k.DnsMode, COREDNS)
	k.NetworkPlugin = defaults.Default(k.NetworkPlugin, CALICO)
}

type KubernetesVersion string

func (ver KubernetesVersion) Validate() error {
	var rs []string

	for _, r := range env.ProjectK8sVersions {
		regex := fmt.Sprintf("^%s.[0-9][0-9]?$", r)
		rs = append(rs, regex)
	}

	msg := fmt.Sprintf("Unsupported Kubernetes version (%s).", ver)
	msg += fmt.Sprintf("Supported versions are: %v", env.ProjectK8sVersions)

	return validation.Var(ver, validation.RegexAny(rs...).Error(msg))
}

type DnsMode string

const (
	COREDNS DnsMode = "coredns"
	KUBEDNS DnsMode = "kubedns"
)

func (m DnsMode) Validate() error {
	return validation.Var(m, validation.OneOf(COREDNS, KUBEDNS))
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
	return validation.Var(p, validation.OneOf(CALICO, CILIUM, CANAL, FLANNEL, WEAVE, KUBE_ROUTER))
}

type Other struct {
	AutoRenewCertificates bool `yaml:"autoRenewCertificates"`
	CopyKubeconfig        bool `yaml:"copyKubeconfig"`
}
