package config

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestKubernetesVersion(t *testing.T) {
	assert.NoError(t, KubernetesVersion("v1.24.0").Validate())
	assert.NoError(t, KubernetesVersion("v1.24.5").Validate())
	assert.NoError(t, KubernetesVersion("v1.25.0").Validate())
	assert.NoError(t, KubernetesVersion("v1.25.5").Validate())
	assert.NoError(t, KubernetesVersion("v1.26.0").Validate())
	assert.NoError(t, KubernetesVersion("v1.26.5").Validate())
	assert.ErrorContains(t, KubernetesVersion("v1.1.1").Validate(), "Unsupported Kubernetes version")
	assert.ErrorContains(t, KubernetesVersion("v1.26.100").Validate(), "Unsupported Kubernetes version")
}

func TestDnsMode(t *testing.T) {
	assert.Error(t, DnsMode("").Validate())
	assert.Error(t, DnsMode("wrong").Validate())
	assert.NoError(t, DnsMode("coredns").Validate())
}

func TestNetworkPlugin(t *testing.T) {
	assert.Error(t, NetworkPlugin("").Validate())
	assert.Error(t, NetworkPlugin("wrong").Validate())
	assert.NoError(t, NetworkPlugin("kube-router").Validate())
	assert.NoError(t, NetworkPlugin("flannel").Validate())
	assert.NoError(t, CALICO.Validate())
	assert.NoError(t, CILIUM.Validate())
}

func TestKubernetes_Empty(t *testing.T) {
	k8s := Kubernetes{}
	assert.ErrorContains(t, k8s.Validate(), "Field 'version' is required and cannot be empty.")
	assert.ErrorContains(t, k8s.Validate(), "Field 'dnsMode' is required and cannot be empty.")
	assert.ErrorContains(t, k8s.Validate(), "Field 'networkPlugin' is required and cannot be empty.")
	assert.NoError(t, defaults.Assign(&k8s).Validate())
}

func TestKubernetes_Valid(t *testing.T) {
	k := Kubernetes{
		Version:       env.ConstKubernetesVersion,
		DnsMode:       COREDNS,
		NetworkPlugin: CALICO,
		Other: Other{
			AutoRenewCertificates: true,
			CopyKubeconfig:        true,
		},
	}

	assert.NoError(t, k.Validate())
}

func TestDefault(t *testing.T) {
	k := Kubernetes{}
	assert.NoError(t, defaults.Set(&k))
	assert.Equal(t, COREDNS, k.DnsMode)
}
