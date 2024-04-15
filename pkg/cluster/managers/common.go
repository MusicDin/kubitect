package managers

import (
	"fmt"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
)

type common struct {
	ClusterName       string
	ClusterPath       string
	SshPrivateKeyPath string
	ConfigDir         string
	CacheDir          string
	SharedDir         string
	Config            *config.Config
	InfraConfig       *infra.Config

	Ansible ansible.Ansible
}

func (e common) K8sVersion() string {
	return string(e.Config.Kubernetes.Version)
}

func (e common) SshUser() string {
	return string(e.Config.Cluster.NodeTemplate.User)
}

func (e common) SshPKey() string {
	return e.SshPrivateKeyPath
}

// extractRemovedNodes returns removed node instances extracted from the event changes.
func extractRemovedNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleDown {
			continue
		}

		node, ok := e.Change.ValueBefore.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}

// extractNewNodes returns new node instances extracted from the event changes.
func extractNewNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleUp {
			continue
		}

		node, ok := e.Change.ValueAfter.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}
