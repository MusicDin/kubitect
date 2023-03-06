package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockPKey(t *testing.T) File {
	path := filepath.Join(t.TempDir(), ".ssh", "id_rsa")

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	assert.NoError(t, err)

	_, err = os.Create(path)
	assert.NoError(t, err)

	return File(path)
}

func MockLocalHost(t *testing.T, name string, def bool) Host {
	host := Host{
		Name:    name,
		Default: def,
		Connection: Connection{
			Type: LOCAL,
		},
		DataResourcePools: []DataResourcePool{
			{
				Name: "pool1",
				Path: "/path1",
			},
			{
				Name: "pool2",
				Path: "/path2",
			},
		},
	}

	assert.NoError(t, defaults.Assign(&host).Validate())
	return host
}

func MockRemoteHost(t *testing.T, name string, def bool, verify bool) Host {
	host := Host{
		Name:    name,
		Default: def,
		Connection: Connection{
			Type: REMOTE,
			IP:   "192.168.113.42",
			User: "mocked-user",
			SSH: ConnectionSSH{
				Keyfile: MockPKey(t),
				Verify:  verify,
			},
		},
	}

	assert.NoError(t, defaults.Assign(&host).Validate())
	return host
}

func MockNodes(t *testing.T) Nodes {
	priority := Uint8(200)

	nodes := Nodes{
		LoadBalancer: LB{
			VIP: "192.168.113.200",
			Instances: []LBInstance{
				{Name: "cls-lb-1", Id: "1", IP: "192.168.113.5"},
				{Name: "cls-lb-2", Id: "2", IP: "192.168.113.6", Priority: &priority},
			},
		},
		Master: Master{
			Instances: []MasterInstance{
				{Name: "cls-master-1", Id: "1", IP: "192.168.113.11", Labels: Labels{"label-1": "value-1"}},
				{Name: "cls-master-2", Id: "2", IP: "192.168.113.12", Taints: []Taint{"taint1=value:NoSchedule"}},
				{Name: "cls-master-3", Id: "3", IP: "192.168.113.13"},
			},
		},
		Worker: Worker{
			Instances: []WorkerInstance{
				{Name: "cls-worker-1", Id: "1", IP: "192.168.113.21", Labels: Labels{"label-1": "value-1"}},
				{Name: "cls-worker-2", Id: "2", IP: "192.168.113.22", Taints: []Taint{"taint1=value:NoSchedule"}},
				{Name: "cls-worker-3", Id: "3", IP: "192.168.113.23"},
			},
		},
	}

	assert.NoError(t, defaults.Set(&nodes))
	return nodes
}

func MockConfig(t *testing.T) Config {
	cfg := Config{
		Hosts: []Host{
			MockLocalHost(t, "local", true),
			MockRemoteHost(t, "remote", false, false),
		},
		Cluster: Cluster{
			Name: "cluster-mock",
			Network: Network{
				CIDR: "192.168.113.0/24",
			},
			Nodes: Nodes{
				Master: Master{
					Instances: []MasterInstance{
						{
							Name: "cluster-mock-master-1",
							Id:   "1",
							IP:   "192.168.113.10",
						},
					},
				},
			},
		},
		Kubernetes: Kubernetes{
			Version: "v1.24.7",
		},
	}

	assert.NoError(t, defaults.Set(&cfg))
	return cfg
}
