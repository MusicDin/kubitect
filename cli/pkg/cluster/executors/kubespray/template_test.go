package kubespray

import (
	modelconfig2 "github.com/MusicDin/kubitect/cli/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/cli/pkg/utils/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubesprayAllTemplate(t *testing.T) {
	nodes := modelconfig2.MockNodes(t)

	tpl := NewKubesprayAllTemplate(t.TempDir(), nodes)
	pop, err := template.Populate(tpl)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Contains(t, pop, "apiserver_loadbalancer_domain_name: \"192.168.113.200\"")
	assert.Contains(t, pop, "loadbalancer_apiserver:\n  address: \"192.168.113.200\"\n  port: 6443")
}

func TestKubesprayK8sClusterTemplate(t *testing.T) {
	cfg := modelconfig2.MockConfig(t)

	tpl := NewKubesprayK8sClusterTemplate(t.TempDir(), cfg)
	pop, err := template.Populate(tpl)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Contains(t, pop, "kube_version: v1.24.7")
	assert.Contains(t, pop, "kube_network_plugin: calico")
	assert.Contains(t, pop, "dns_mode: coredns")
	assert.Contains(t, pop, "auto_renew_certificates: false")
}

func TestKubesprayAddonsTemplate(t *testing.T) {
	tpl := NewKubesprayAddonsTemplate(t.TempDir(), "test: test")
	pop, err := template.Populate(tpl)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Equal(t, pop, "test: test")
}

func TestKubesprayEtcdTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	tpl := NewKubesprayEtcdTemplate(tmpDir)
	pop, err := template.Populate(tpl)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Contains(t, pop, "etcd_deployment_type: host")
}

func TestHostsTemplate(t *testing.T) {
	hosts := []modelconfig2.Host{
		modelconfig2.MockLocalHost(t, "local", true),
		modelconfig2.MockLocalHost(t, "localhost", false),
		modelconfig2.MockRemoteHost(t, "remote", false, false),
	}

	tpl := NewHostsTemplate(t.TempDir(), "~/.ssh/id_rsa", hosts)
	pop, err := template.Populate(tpl)

	expect := template.TrimTemplate(`
		all:
			hosts:
				local:
					ansible_connection: local
					ansible_host: localhost
				localhost:
					ansible_connection: local
					ansible_host: localhost
				remote:
					ansible_connection: ssh
					ansible_user: mocked-user
					ansible_host: 192.168.113.42
					ansible_port: 22
					ansible_private_key_file: ~/.ssh/id_rsa
			children:
				kubitect_hosts:
					hosts:
						local:
						localhost:
						remote:
	`)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Equal(t, expect, pop)
}

func TestNodesTemplate(t *testing.T) {
	nodes := modelconfig2.MockNodes(t)

	tpl := NewNodesTemplate(t.TempDir(), nodes, nodes)
	pop, err := template.Populate(tpl)

	expect := template.TrimTemplate(`
		all:
		  hosts:
		  	cls-lb-1:
		      ansible_host: 192.168.113.5
		      priority: 10
		  	cls-lb-2:
		      ansible_host: 192.168.113.6
		      priority: 200
				cls-master-1:
					ansible_host: 192.168.113.11
					node_labels:
						label-1: value-1
				cls-master-2:
					ansible_host: 192.168.113.12
					node_taints:
						- "taint1=value:NoSchedule"
				cls-master-3:
					ansible_host: 192.168.113.13
				cls-worker-1:
					ansible_host: 192.168.113.21
					node_labels:
						label-1: value-1
				cls-worker-2:
					ansible_host: 192.168.113.22
					node_taints:
						- "taint1=value:NoSchedule"
				cls-worker-3:
					ansible_host: 192.168.113.23
			children:
				haproxy:
					hosts:
						cls-lb-1:
						cls-lb-2:
				etcd:
					hosts:
						cls-master-1:
						cls-master-2:
						cls-master-3:
				k8s_cluster:
					children:
						kube_control_plane:
							hosts:
								cls-master-1:
								cls-master-2:
								cls-master-3:
						kube_node:
							hosts:
								cls-worker-1:
								cls-worker-2:
								cls-worker-3:
	`)

	assert.NoError(t, err)
	assert.NoError(t, tpl.Write())
	assert.Equal(t, expect, pop)
}

func TestNodesTemplate_NoWorkers(t *testing.T) {
	nodes := modelconfig2.MockNodes(t)
	nodes.Worker = modelconfig2.Worker{}

	tpl := NewNodesTemplate(t.TempDir(), nodes, nodes)
	pop, err := template.Populate(tpl)

	expect := template.TrimTemplate(`
		all:
		  hosts:
		  	cls-lb-1:
		      ansible_host: 192.168.113.5
		      priority: 10
		  	cls-lb-2:
		      ansible_host: 192.168.113.6
		      priority: 200
				cls-master-1:
					ansible_host: 192.168.113.11
					node_labels:
						label-1: value-1
				cls-master-2:
					ansible_host: 192.168.113.12
					node_taints:
						- "taint1=value:NoSchedule"
				cls-master-3:
					ansible_host: 192.168.113.13
			children:
				haproxy:
					hosts:
						cls-lb-1:
						cls-lb-2:
				etcd:
					hosts:
						cls-master-1:
						cls-master-2:
						cls-master-3:
				k8s_cluster:
					children:
						kube_control_plane:
							hosts:
								cls-master-1:
								cls-master-2:
								cls-master-3:
						kube_node:
							hosts:
								cls-master-1:
								cls-master-2:
								cls-master-3:
	`)

	assert.NoError(t, err)
	assert.Equal(t, expect, pop)
}
