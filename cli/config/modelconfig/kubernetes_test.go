package modelconfig

import (
	"testing"

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
	ver := MasterVersion("master")
	url := URL("https://github.com/kubernetes-sigs/kubespray")

	ks1 := Kubespray{
		Version: &ver,
	}

	ks2 := Kubespray{
		Version: &ver,
		URL:     &url,
	}

	assert.EqualError(t, Kubespray{}.Validate(), "Field 'version' is required.")
	assert.NoError(t, ks1.Validate())
	assert.NoError(t, ks2.Validate())
}

func TestKubernetes_Empty(t *testing.T) {
	assert.ErrorContains(t, Kubernetes{}.Validate(), "Field 'version' is required.")
	assert.ErrorContains(t, Kubernetes{}.Validate(), "Field 'kubespray' is required and cannot be empty.")
}

func TestKubernetes_Valid(t *testing.T) {
	ver := Version("v1.2.3")
	mver := MasterVersion("v1.2.3")
	boolean := true
	dnsMode := COREDNS
	netPlugin := CALICO

	k := Kubernetes{
		Version:       &ver,
		DnsMode:       &dnsMode,
		NetworkPlugin: &netPlugin,
		Kubespray: Kubespray{
			Version: &mver,
		},
		Other: Other{
			AutoRenewCertificates: &boolean,
			CopyKubeconfig:        &boolean,
		},
	}

	assert.NoError(t, k.Validate())
}
