package event

import (
	"github.com/MusicDin/kubitect/pkg/utils/cmp"
)

var UpgradeRules = []Rule{
	{
		Type:            Allow,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("kubernetes.version"),
	},
	// Default rule.
	{
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("@"),
		Message:         "Change is not allowed. Upgrade action allows changing only 'kubernetes.version'.",
	},
}

var ScaleRules = []Rule{
	{
		Type:            Allow,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("cluster.nodes.worker.instances.@"),
		ActionType:      Action_ScaleDown,
	},
	{
		Type:            Allow,
		MatchChangeType: cmp.Create,
		MatchPath:       NewRulePath("cluster.nodes.worker.instances.@"),
		ActionType:      Action_ScaleUp,
	},
	{
		Type:            Allow,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("cluster.nodes.loadBalancer.instances.@"),
		ActionType:      Action_ScaleDown,
	},
	{
		Type:            Allow,
		MatchChangeType: cmp.Create,
		MatchPath:       NewRulePath("cluster.nodes.loadBalancer.instances.@"),
		ActionType:      Action_ScaleUp,
	},
	{
		Type:            Error,
		MatchChangeType: cmp.Create,
		MatchPath:       NewRulePath("cluster.nodes.master.instances.@"),
		Message:         "Currently, control plane cannot be scaled.",
	},
	{
		Type:            Allow,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("cluster.nodes.master.instances.@"),
		Message:         "Currently, control plane cannot be scaled.",
	},
	// Allow addition and deletion of hosts.
	{
		Type:            Allow,
		MatchChangeType: cmp.Create,
		MatchPath:       NewRulePath("hosts.@"),
	},
	{
		Type:            Allow,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("hosts.@"),
	},
	// Default rule.
	{
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("*"),
		Message:         "Change is not allowed. Scale action allows only addition and removal of worker and load balancer nodes.",
	},
}

var ModifyRules = []Rule{
	{
		// Warn about main resource pool path change (will replace the VM).
		Type:            Warn,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("hosts.*.mainResourcePoolPath"),
		Message:         "Changing main resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks.",
	},
	{
		// Warn about data resource pool removal (will destroy the pool).
		Type:            Warn,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("hosts.*.dataResourcePools.*"),
		Message:         "Removing data resource pool will destroy all the data on that location.",
	},
	{
		// Warn about data resource pool path change (will destroy the pool).
		Type:            Warn,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("hosts.*.dataResourcePools.*.path"),
		Message:         "Changing data resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks",
	},
	{
		// Allow other data resource pool changes.
		Type:            Allow,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("hosts.*.dataResourcePools.*"),
	},
	{
		// Prevent cluster network changes.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.network"),
		Message:         "Once the cluster is created, further changes to the network properties are not allowed. Such action may render the cluster unusable.",
	},
	{
		// Prevent nodeTemplate changes.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.nodeTemplate"),
		Message:         "Once the cluster is created, further changes to the nodeTemplate properties are not allowed. Such action may render the cluster unusable.",
	},
	{
		// Prevent removing nodes.
		Type:            Error,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.instances.@"),
		Message:         "To remove existing nodes run apply command with '--action scale' flag.",
	},
	{
		// Prevent adding nodes.
		Type:            Error,
		MatchChangeType: cmp.Create,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.instances.@"),
		Message:         "To add new nodes run apply command with '--action scale' flag.",
	},
	{
		// Prevent default cpu, ram and main disk size changes.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.default.{cpu, ram, mainDiskSize}"),
		Message:         "Changing any default physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action may render the cluster unusable.",
	},
	{
		// Prevent cpu, ram and main disk size changes.
		Type:            Error,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.instances.@.{cpu, ram, mainDiskSize}"),
		Message:         "Changing any physical properties of nodes (cpu, ram, mainDiskSize) is not allowed. Such action will recreate the node.",
	},
	{
		// Prevent IP and MAC changes.
		Type:            Error,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.instances.@.{ip, mac}"),
		Message:         "Changing IP or MAC address of the node is not allowed. Such action may render the cluster unusable.",
	},
	{
		// Warn about data disk changes.
		Type:            Warn,
		MatchChangeType: cmp.Modify,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker}.instances.*.dataDisks.*"),
		Message:         "Changing data disk properties, will recreate the disk (removing all of its content in the process).",
	},
	{
		// Warn about data disk removal.
		Type:            Warn,
		MatchChangeType: cmp.Delete,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker}.instances.*.dataDisks.*"),
		Message:         "One or more data disks will be removed.",
	},
	{
		// Allow changes to LB forward ports.
		Type:            Allow,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.nodes.loadBalancer.forwardPorts.*"),
	},
	{
		// Prevent VIP changes.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.nodes.loadBalancer.vip"),
		Message:         "Once the cluster is created, changing virtual IP (VIP) is not allowed. Such action may render the cluster unusable.",
	},
	{
		// Allow all other node properties to be changed.
		Type:            Allow,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("cluster.nodes.{master, worker, loadBalancer}.instances.*"),
	},
	{
		// Prevent k8s properties changes.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("kubernetes.version"),
		Message:         "Changing Kubernetes is allowed only when upgrading the cluster.\nTo upgrade the cluster run apply command with '--action upgrade' flag.",
	},
	{
		// Allow addons changes.
		Type:            Allow,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("addons"),
	},
	{
		// Default rule.
		Type:            Error,
		MatchChangeType: cmp.Any,
		MatchPath:       NewRulePath("@"),
		Message:         "Change is not allowed.",
	},
}
