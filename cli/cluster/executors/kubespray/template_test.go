package kubespray

import (
	"cli/config/modelconfig"
	"cli/utils/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostsTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	hosts := []modelconfig.Host{
		modelconfig.MockLocalHost(t, "local", true),
		modelconfig.MockLocalHost(t, "localhost", false),
		modelconfig.MockRemoteHost(t, "remote", false, false),
	}

	tpl := NewHostsTemplate(tmpDir, "~/.ssh/id_rsa", hosts)
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
	assert.Equal(t, expect, pop)
}

func TestNodesTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	nodes := modelconfig.MockNodes(t)

	tpl := NewNodesTemplate(tmpDir, nodes, nodes)
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
	assert.Equal(t, expect, pop)
}

func TestNodesTemplate_NoWorkers(t *testing.T) {
	tmpDir := t.TempDir()

	nodes := modelconfig.MockNodes(t)
	nodes.Worker = modelconfig.Worker{}

	tpl := NewNodesTemplate(tmpDir, nodes, nodes)
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
