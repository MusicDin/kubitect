package executors

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/utils/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubesprayTemplate_All(t *testing.T) {
	tpl := NewTemplate("kubespray/all.yaml", config.MockNodes(t))
	pop, err := template.Populate(tpl)

	require.NoError(t, err)
	require.NoError(t, tpl.Write(filepath.Join(t.TempDir(), "tpl")))
	assert.Contains(t, pop, "apiserver_loadbalancer_domain_name: \"192.168.113.200\"")
	assert.Contains(t, pop, "loadbalancer_apiserver:\n  address: \"192.168.113.200\"\n  port: 6443")
}

func TestKubesprayTemplate_K8sCluster(t *testing.T) {
	tpl := NewTemplate("kubespray/k8s-cluster.yaml", config.MockConfig(t))
	pop, err := template.Populate(tpl)

	require.NoError(t, err)
	require.NoError(t, tpl.Write(filepath.Join(t.TempDir(), "tpl")))
	assert.Contains(t, pop, fmt.Sprintf("kube_version: %s", env.ConstKubernetesVersion))
	assert.Contains(t, pop, "kube_network_plugin: calico")
	assert.Contains(t, pop, "dns_mode: coredns")
	assert.Contains(t, pop, "auto_renew_certificates: false")
}

func TestKubesprayTemplate_Etcd(t *testing.T) {
	tpl := NewTemplate("kubespray/etcd.yaml", "")
	pop, err := template.Populate(tpl)

	require.NoError(t, err)
	require.NoError(t, tpl.Write(filepath.Join(t.TempDir(), "tpl")))
	assert.Contains(t, pop, "etcd_deployment_type: host")
}

func TestKubesprayTemplate_Inventory(t *testing.T) {
	nodes := config.MockNodes(t)

	values := struct {
		ConfigNodes config.Nodes
		InfraNodes  config.Nodes
	}{
		ConfigNodes: nodes,
		InfraNodes:  nodes,
	}

	tpl := NewTemplate("kubespray/inventory.yaml", values)
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

	require.NoError(t, err)
	require.NoError(t, tpl.Write(filepath.Join(t.TempDir(), "tpl")))
	assert.Equal(t, expect, pop)
}

func TestTemplate_Inventory_NoWorkers(t *testing.T) {
	nodes := config.MockNodes(t)
	nodes.Worker = config.Worker{}

	values := struct {
		ConfigNodes config.Nodes
		InfraNodes  config.Nodes
	}{
		ConfigNodes: nodes,
		InfraNodes:  nodes,
	}

	tpl := NewTemplate("kubespray/inventory.yaml", values)
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

	require.NoError(t, err)
	assert.Equal(t, expect, pop)
}
