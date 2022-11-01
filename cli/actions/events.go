package actions

import (
	"cli/cmp"
	"cli/env"
)

type ChangeType string

const (
	OK ChangeType = "ok"

	// WARN change requires user permission to continue.
	WARN ChangeType = "warn"

	// BLOCK change prevents further actions on the cluster.
	BLOCK ChangeType = "block"
)

type OnChangeEvent struct {
	cType        ChangeType
	msg          string
	path         string
	paths        []string
	triggerPaths []string
	action       cmp.ActionType
}

func (e OnChangeEvent) Paths() []string {
	if len(e.path) > 0 {
		return []string{e.path}
	}

	return e.paths
}

func (e OnChangeEvent) Action() cmp.ActionType {
	return e.action
}

func (e *OnChangeEvent) TriggerPath(path string) {
	e.triggerPaths = append(e.triggerPaths, path)
}

func events(a env.ApplyAction) []*OnChangeEvent {
	switch a {
	case env.CREATE:
		return ModifyEvents
	case env.SCALE:
		return ScaleEvents
	case env.UPGRADE:
		return UpgradeEvents
	default:
		return nil
	}
}

var (
	UpgradeEvents = []*OnChangeEvent{
		{
			cType: OK,
			path:  "kubernetes.version",
		},
		{
			cType: OK,
			path:  "kubernetes.kubespray.version",
		},
	}

	ScaleEvents = []*OnChangeEvent{
		{
			cType:  OK,
			action: cmp.DELETE,
			path:   "cluster.nodes.worker.instances.[*]",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			path:   "cluster.nodes.worker.instances.[*]",
		},
		{
			cType:  OK,
			action: cmp.DELETE,
			path:   "cluster.nodes.loadBalancer.instances.[*]",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			path:   "cluster.nodes.loadBalancer.instances.[*]",
		},
	}

	ModifyEvents = []*OnChangeEvent{
		// Warn data destructive host changes
		{
			cType:  WARN,
			action: cmp.MODIFY,
			path:   "hosts[*].mainResourcePoolPath",
			msg:    "Changing main resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks.",
		},
		{
			cType:  WARN,
			action: cmp.MODIFY,
			path:   "hosts[*].dataResourcePools[*].path",
			msg:    "Changing data resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks",
		},
		{
			cType:  WARN,
			action: cmp.DELETE,
			path:   "hosts[*].dataResourcePools[*]",
			msg:    "Removing data resource pool will destroy all the data on that location.",
		},
		// Allow other host changes
		{
			cType: OK,
			path:  "hosts",
		},
		// Prevent cluster network changes
		{
			cType: BLOCK,
			path:  "cluster.network",
			msg:   "Once the cluster is created, further changes to the network properties are not allowed. Such action may render the cluster unusable.",
		},
		// Prevent nodeTemplate changes
		{
			cType: BLOCK,
			path:  "cluster.nodeTemplate",
			msg:   "Once the cluster is created, further changes to the nodeTemplate properties are not allowed. Such action may render the cluster unusable.",
		},
		// Prevent removing nodes
		{
			cType:  BLOCK,
			action: cmp.DELETE,
			paths: []string{
				"cluster.nodes.lb.instances.[*]",
				"cluster.nodes.worker.instances.[*]",
				"cluster.nodes.master.instances.[*]",
			},
			msg: "To remove existing nodes run apply command with '--action scale' flag.",
		},
		// Prevent adding nodes
		{
			cType:  BLOCK,
			action: cmp.CREATE,
			paths: []string{
				"cluster.nodes.lb.instances.[*]",
				"cluster.nodes.worker.instances.[*]",
				"cluster.nodes.master.instances.[*]",
			},
			msg: "To add new nodes run apply command with '--action scale' flag.",
		},
		// Prevent default CPU, RAM and main disk size changes
		{
			cType: BLOCK,
			paths: []string{
				"cluster.nodes.worker.default.cpu",
				"cluster.nodes.worker.default.ram",
				"cluster.nodes.worker.default.mainDiskSize",
				"cluster.nodes.master.default.cpu",
				"cluster.nodes.master.default.ram",
				"cluster.nodes.master.default.mainDiskSize",
				"cluster.nodes.loadBalancer.default.cpu",
				"cluster.nodes.loadBalancer.default.ram",
				"cluster.nodes.loadBalancer.default.mainDiskSize",
			},
			msg: "Changing any default physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action may render the cluster unusable.",
		},
		// Prevent CPU, RAM and main disk size changes
		{
			cType:  BLOCK,
			action: cmp.MODIFY,
			paths: []string{
				"cluster.nodes.worker.instances.[*].cpu",
				"cluster.nodes.worker.instances.[*].ram",
				"cluster.nodes.worker.instances.[*].mainDiskSize",
				"cluster.nodes.master.instances.[*].cpu",
				"cluster.nodes.master.instances.[*].ram",
				"cluster.nodes.master.instances.[*].mainDiskSize",
				"cluster.nodes.loadBalancer.instances.[*].cpu",
				"cluster.nodes.loadBalancer.instances.[*].ram",
				"cluster.nodes.loadBalancer.instances.[*].mainDiskSize",
			},
			msg: "Changing any physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action will recreate the node.",
		},
		// Prevent IP and MAC changes
		{
			cType:  BLOCK,
			action: cmp.MODIFY,
			paths: []string{
				"cluster.nodes.worker.instances.[*].ip",
				"cluster.nodes.worker.instances.[*].mac",
				"cluster.nodes.master.instances.[*].ip",
				"cluster.nodes.master.instances.[*].mac",
				"cluster.nodes.loadBalancer.instances.[*].ip",
				"cluster.nodes.loadBalancer.instances.[*].mac",
			},
			msg: "Changing IP or MAC address of the node is not allowed. Such action may render the cluster unusable.",
		},
		// Data disk changes
		{
			cType:  WARN,
			action: cmp.MODIFY,
			paths: []string{
				"cluster.nodes.worker.instances.[*].dataDisks.[*]",
				"cluster.nodes.master.instances.[*].dataDisks.[*]",
			},
			msg: "Changing data disk properties, will recreate the disk (removing all of its content in the process).",
		},
		{
			cType:  WARN,
			action: cmp.DELETE,
			paths: []string{
				"cluster.nodes.master.instances.[*].dataDisks.[*]",
				"cluster.nodes.worker.instances.[*].dataDisks.[*]",
			},
			msg: "One or more data disks will be removed.",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			paths: []string{
				"cluster.nodes.master.instances.[*].dataDisks.[*]",
				"cluster.nodes.worker.instances.[*].dataDisks.[*]",
			},
		},
		// Prevent VIP changes
		{
			cType: BLOCK,
			path:  "cluster.nodes.loadBalancer.vip",
			msg:   "Once the cluster is created, changing virtual IP (VIP) is not allowed. Such action may render the cluster unusable.",
		},
		// Allow all other node properties to be changed
		{
			cType: OK,
			paths: []string{
				"cluster.nodes.master.instances.[*]",
				"cluster.nodes.worker.instances.[*]",
				"cluster.nodes.loadBalancer.instances.[*]",
			},
		},
		// Prevent k8s properties changes
		{
			cType: BLOCK,
			paths: []string{
				"kubernetes.version",
				"kubernetes.kubespray.version",
			},
			msg: "Changing Kubernetes or Kubespray version is allowed only when upgrading the cluster.\nTo upgrade the cluster run apply command with '--action upgrade' flag.",
		},
		// Allow addons changes
		{
			cType: OK,
			path:  "addons",
		},
		// Allow kubitect (project metadata) changes
		{
			cType: OK,
			path:  "kubitect",
		},
	}
)
