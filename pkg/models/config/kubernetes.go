package config

import (
	"fmt"
	"strings"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
)

type Kubernetes struct {
	Version       KubernetesVersion `yaml:"version"`
	Manager       KubernetesManager `yaml:"manager"`
	DnsMode       DnsMode           `yaml:"dnsMode"`
	NetworkPlugin NetworkPlugin     `yaml:"networkPlugin"`
	Other         Other             `yaml:"other"`
}

func (k Kubernetes) Validate() error {
	return v.Struct(&k,
		v.Field(&k.Version, v.NotEmpty(), v.VSemVer()),
		v.Field(&k.DnsMode, v.NotEmpty()),
		v.Field(&k.NetworkPlugin, v.NotEmpty()),
		v.Field(&k.Other),
	)
}

func (k *Kubernetes) SetDefaults() {
	k.Version = defaults.Default(k.Version, env.ConstKubernetesVersion)
	k.Manager = defaults.Default(k.Manager, ManagerKubespray)
	k.DnsMode = defaults.Default(k.DnsMode, COREDNS)
	k.NetworkPlugin = defaults.Default(k.NetworkPlugin, CALICO)
}

type KubernetesVersion string

func (ver KubernetesVersion) Validate() error {
	var err error

	version := strings.TrimPrefix(string(ver), "v")

	msg := fmt.Sprintf("Unsupported Kubernetes version (%s). ", ver)
	msg += fmt.Sprintf("Supported versions are:\n%s", strings.Join(env.ProjectK8sVersions, "\n"))

	for _, verRange := range env.ProjectK8sVersions {
		verRange = strings.ReplaceAll(verRange, " ", "")
		verRange = strings.ReplaceAll(verRange, "v", "")

		verRangeSplit := strings.Split(verRange, "-")
		min := verRangeSplit[0]
		max := verRangeSplit[1]

		// If validation passes, return nil
		err = v.Var(version, v.SemVerInRange(min, max).Error(msg))
		if err == nil {
			return nil
		}
	}

	return err
}

type KubernetesManager string

const (
	ManagerKubespray = "kubespray"
	ManagerK3s       = "k3s"
)

func (m KubernetesManager) Validate() error {
	return v.Var(m, v.OneOf(ManagerKubespray, ManagerK3s))
}

type DnsMode string

const (
	COREDNS DnsMode = "coredns"
)

func (m DnsMode) Validate() error {
	return v.Var(m, v.OneOf(COREDNS))
}

type NetworkPlugin string

const (
	CALICO      NetworkPlugin = "calico"
	CILIUM      NetworkPlugin = "cilium"
	FLANNEL     NetworkPlugin = "flannel"
	KUBE_ROUTER NetworkPlugin = "kube-router"
)

func (p NetworkPlugin) Validate() error {
	return v.Var(p, v.OneOf(CALICO, CILIUM, FLANNEL, KUBE_ROUTER))
}

type Other struct {
	AutoRenewCertificates bool `yaml:"autoRenewCertificates"`
	MergeKubeconfig       bool `yaml:"mergeKubeconfig"`
}
