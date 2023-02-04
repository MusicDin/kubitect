package modelconfig

import (
	"testing"

	"cli/env"
	"cli/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestDnsMode(t *testing.T) {
	assert.Error(t, DnsMode("").Validate())
	assert.Error(t, DnsMode("wrong").Validate())
	assert.NoError(t, DnsMode("kubedns").Validate())
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

func TestKubespray(t *testing.T) {
	ks1 := Kubespray{
		Version: MasterVersion("master"),
	}

	ks2 := Kubespray{
		Version: MasterVersion("master"),
		URL:     URL(env.ConstKubesprayUrl),
	}

	assert.ErrorContains(t, Kubespray{}.Validate(), "Field 'version' must be a valid semantic version prefixed with 'v'")
	assert.NoError(t, ks1.Validate())
	assert.NoError(t, ks2.Validate())
}

func TestKubernetes_Empty(t *testing.T) {
	k8s := Kubernetes{}
	assert.ErrorContains(t, k8s.Validate(), "Field 'version' is required and cannot be empty.")
	assert.ErrorContains(t, k8s.Validate(), "Field 'dnsMode' is required and cannot be empty.")
	assert.ErrorContains(t, k8s.Validate(), "Field 'networkPlugin' is required and cannot be empty.")
	assert.ErrorContains(t, defaults.Assign(&k8s).Validate(), "Field 'kubespray' is required and cannot be empty.")
}

func TestKubernetes_Valid(t *testing.T) {
	k := Kubernetes{
		Version:       Version("v1.2.3"),
		DnsMode:       COREDNS,
		NetworkPlugin: CALICO,
		Kubespray: Kubespray{
			Version: MasterVersion("v1.2.3"),
		},
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

func TestDefault_Fail(t *testing.T) {
	dns := KUBEDNS
	k := Kubernetes{
		DnsMode: dns,
	}
	assert.NoError(t, defaults.Set(&k))
	assert.Equal(t, KUBEDNS, k.DnsMode)
}
