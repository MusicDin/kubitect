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

// triggerEvents checks whether any events is triggered based on the provided
// changes and action. If any blocking event is triggered or some changes are
// not covered by any event, an error is thrown.
func triggerEvents(diff *cmp.DiffNode, action env.ApplyAction) []*OnChangeEvent {
	events := events(action)
	triggered := cmp.TriggerEvents(diff, events)
	nmc := cmp.NonMatchingChanges(diff, events)

	// Changes that are not covered by any event are automatically
	// considered disallowed (blocking).
	if len(nmc) > 0 {
		var paths []string

		for _, ch := range nmc {
			paths = append(paths, ch.Path)
		}

		triggered = append(triggered, &OnChangeEvent{
			msg:   "Disallowed changes.",
			paths: paths,
		})
	}

	return triggered
}

// events returns a copy of OnChangeEvent-s.
// Since OnChangeEvents has a setter method, each event
// must be a pointer.
func events(a env.ApplyAction) []*OnChangeEvent {
	var copy []OnChangeEvent

	switch a {
	case env.CREATE:
		copy = ModifyEvents
	case env.SCALE:
		copy = ScaleEvents
	case env.UPGRADE:
		copy = UpgradeEvents
	default:
		return nil
	}

	var events []*OnChangeEvent

	for _, e := range copy {
		events = append(events, &e)
	}

	return events
}

// Events
var (
	UpgradeEvents = []OnChangeEvent{
		{
			cType: OK,
			path:  "Kubernetes.Version",
		},
		{
			cType: OK,
			path:  "Kubernetes.Kubespray.Version",
		},
	}

	ScaleEvents = []OnChangeEvent{
		{
			cType:  OK,
			action: cmp.DELETE,
			path:   "Cluster.Nodes.Worker.Instances.*",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			path:   "Cluster.Nodes.Worker.Instances.*",
		},
		{
			cType:  OK,
			action: cmp.DELETE,
			path:   "Cluster.Nodes.LoadBalancer.Instances.*",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			path:   "Cluster.Nodes.LoadBalancer.Instances.*",
		},
	}

	ModifyEvents = []OnChangeEvent{
		// Warn data destructive host changes
		{
			cType:  WARN,
			action: cmp.MODIFY,
			path:   "Hosts.*.MainResourcePoolPath",
			msg:    "Changing main resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks.",
		},
		{
			cType:  WARN,
			action: cmp.DELETE,
			path:   "Hosts.*.DataResourcePools.*",
			msg:    "Removing data resource pool will destroy all the data on that location.",
		},
		{
			cType:  WARN,
			action: cmp.MODIFY,
			path:   "Hosts.*.DataResourcePools.*.Path",
			msg:    "Changing data resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks",
		},
		// Allow other host changes
		{
			cType: OK,
			path:  "Hosts",
		},
		// Prevent cluster network changes
		{
			cType: BLOCK,
			path:  "Cluster.Network",
			msg:   "Once the cluster is created, further changes to the network properties are not allowed. Such action may render the cluster unusable.",
		},
		// Prevent nodeTemplate changes
		{
			cType: BLOCK,
			path:  "Cluster.NodeTemplate",
			msg:   "Once the cluster is created, further changes to the nodeTemplate properties are not allowed. Such action may render the cluster unusable.",
		},
		// Prevent removing nodes
		{
			cType:  BLOCK,
			action: cmp.DELETE,
			paths: []string{
				"Cluster.Nodes.LoadBalancer.Instances.*",
				"Cluster.Nodes.Worker.Instances.*",
				"Cluster.Nodes.Master.Instances.*",
			},
			msg: "To remove existing nodes run apply command with '--action scale' flag.",
		},
		// Prevent adding nodes
		{
			cType:  BLOCK,
			action: cmp.CREATE,
			paths: []string{
				"Cluster.Nodes.LoadBalancer.Instances.*",
				"Cluster.Nodes.Worker.Instances.*",
				"Cluster.Nodes.Master.Instances.*",
			},
			msg: "To add new nodes run apply command with '--action scale' flag.",
		},
		// Prevent default CPU, RAM and main disk size changes
		{
			cType: BLOCK,
			paths: []string{
				"Cluster.Nodes.Worker.Default.CPU",
				"Cluster.Nodes.Worker.Default.RAM",
				"Cluster.Nodes.Worker.Default.MainDiskSize",
				"Cluster.Nodes.Master.Default.CPU",
				"Cluster.Nodes.Master.Default.RAM",
				"Cluster.Nodes.Master.Default.MainDiskSize",
				"Cluster.Nodes.LoadBalancer.Default.CPU",
				"Cluster.Nodes.LoadBalancer.Default.RAM",
				"Cluster.Nodes.LoadBalancer.Default.MainDiskSize",
			},
			msg: "Changing any default physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action may render the cluster unusable.",
		},
		// Prevent CPU, RAM and main disk size changes
		{
			cType:  BLOCK,
			action: cmp.MODIFY,
			paths: []string{
				"Cluster.Nodes.Worker.Instances.*.CPU",
				"Cluster.Nodes.Worker.Instances.*.RAM",
				"Cluster.Nodes.Worker.Instances.*.MainDiskSize",
				"Cluster.Nodes.Master.Instances.*.CPU",
				"Cluster.Nodes.Master.Instances.*.RAM",
				"Cluster.Nodes.Master.Instances.*.MainDiskSize",
				"Cluster.Nodes.LoadBalancer.Instances.*.CPU",
				"Cluster.Nodes.LoadBalancer.Instances.*.RAM",
				"Cluster.Nodes.LoadBalancer.Instances.*.MainDiskSize",
			},
			msg: "Changing any physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action will recreate the node.",
		},
		// Prevent IP and MAC changes
		{
			cType:  BLOCK,
			action: cmp.MODIFY,
			paths: []string{
				"Cluster.Nodes.Worker.Instances.*.IP",
				"Cluster.Nodes.Worker.Instances.*.MAC",
				"Cluster.Nodes.Master.Instances.*.IP",
				"Cluster.Nodes.Master.Instances.*.MAC",
				"Cluster.Nodes.LoadBalancer.Instances.*.IP",
				"Cluster.Nodes.LoadBalancer.Instances.*.MAC",
			},
			msg: "Changing IP or MAC address of the node is not allowed. Such action may render the cluster unusable.",
		},
		// Data disk changes
		{
			cType:  WARN,
			action: cmp.MODIFY,
			paths: []string{
				"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
				"Cluster.Nodes.Master.Instances.*.DataDisks.*",
			},
			msg: "Changing data disk properties, will recreate the disk (removing all of its content in the process).",
		},
		{
			cType:  WARN,
			action: cmp.DELETE,
			paths: []string{
				"Cluster.Nodes.Master.Instances.*.DataDisks.*",
				"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
			},
			msg: "One or more data disks will be removed.",
		},
		{
			cType:  OK,
			action: cmp.CREATE,
			paths: []string{
				"Cluster.Nodes.Master.Instances.*.DataDisks.*",
				"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
			},
		},
		// Prevent VIP changes
		{
			cType: BLOCK,
			path:  "Cluster.Nodes.LoadBalancer.VIP",
			msg:   "Once the cluster is created, changing virtual IP (VIP) is not allowed. Such action may render the cluster unusable.",
		},
		// Allow all other node properties to be changed
		{
			cType: OK,
			paths: []string{
				"Cluster.Nodes.Master.Instances.*",
				"Cluster.Nodes.Worker.Instances.*",
				"Cluster.Nodes.LoadBalancer.Instances.*",
			},
		},
		// Prevent k8s properties changes
		{
			cType: BLOCK,
			paths: []string{
				"Kubernetes.Version",
				"Kubernetes.Kubespray.Version",
			},
			msg: "Changing Kubernetes or Kubespray version is allowed only when upgrading the cluster.\nTo upgrade the cluster run apply command with '--action upgrade' flag.",
		},
		// Allow addons changes
		{
			cType: OK,
			path:  "Addons",
		},
		// Allow kubitect (project metadata) changes
		{
			cType: OK,
			path:  "Kubitect",
		},
	}
)
