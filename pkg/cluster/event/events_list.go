package event

import "github.com/MusicDin/kubitect/pkg/utils/cmp"

var UpgradeEvents = Events{
	{
		eType: OK,
		path:  "Kubernetes.Version",
	},
}

var ScaleEvents = Events{
	{
		eType:  SCALE_DOWN,
		action: cmp.DELETE,
		path:   "Cluster.Nodes.Worker.Instances.*",
	},
	{
		eType:  SCALE_UP,
		action: cmp.CREATE,
		path:   "Cluster.Nodes.Worker.Instances.*",
	},
	{
		eType:  SCALE_DOWN,
		action: cmp.DELETE,
		path:   "Cluster.Nodes.LoadBalancer.Instances.*",
	},
	{
		eType:  SCALE_UP,
		action: cmp.CREATE,
		path:   "Cluster.Nodes.LoadBalancer.Instances.*",
	},
}

var ModifyEvents = Events{
	// Warn data destructive host changes
	{
		eType:  WARN,
		action: cmp.MODIFY,
		path:   "Hosts.*.MainResourcePoolPath",
		msg:    "Changing main resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks.",
	},
	{
		eType:  WARN,
		action: cmp.DELETE,
		path:   "Hosts.*.DataResourcePools.*",
		msg:    "Removing data resource pool will destroy all the data on that location.",
	},
	{
		eType:  WARN,
		action: cmp.MODIFY,
		path:   "Hosts.*.DataResourcePools.*.Path",
		msg:    "Changing data resource pool location will trigger recreation of all resources bound to that resource pool, such as virtual machines and data disks",
	},
	// Allow other host changes
	{
		eType: OK,
		path:  "Hosts",
	},
	// Prevent cluster network changes
	{
		eType: BLOCK,
		path:  "Cluster.Network",
		msg:   "Once the cluster is created, further changes to the network properties are not allowed. Such action may render the cluster unusable.",
	},
	// Prevent nodeTemplate changes
	{
		eType: BLOCK,
		path:  "Cluster.NodeTemplate",
		msg:   "Once the cluster is created, further changes to the nodeTemplate properties are not allowed. Such action may render the cluster unusable.",
	},
	// Prevent removing nodes
	{
		eType:  BLOCK,
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
		eType:  BLOCK,
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
		eType: BLOCK,
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
		eType:  BLOCK,
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
		eType:  BLOCK,
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
		eType:  WARN,
		action: cmp.MODIFY,
		paths: []string{
			"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
			"Cluster.Nodes.Master.Instances.*.DataDisks.*",
		},
		msg: "Changing data disk properties, will recreate the disk (removing all of its content in the process).",
	},
	{
		eType:  WARN,
		action: cmp.DELETE,
		paths: []string{
			"Cluster.Nodes.Master.Instances.*.DataDisks.*",
			"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
		},
		msg: "One or more data disks will be removed.",
	},
	{
		eType:  OK,
		action: cmp.CREATE,
		paths: []string{
			"Cluster.Nodes.Master.Instances.*.DataDisks.*",
			"Cluster.Nodes.Worker.Instances.*.DataDisks.*",
		},
	},
	// Prevent VIP changes
	{
		eType: BLOCK,
		path:  "Cluster.Nodes.LoadBalancer.VIP",
		msg:   "Once the cluster is created, changing virtual IP (VIP) is not allowed. Such action may render the cluster unusable.",
	},
	// Allow all other node properties to be changed
	{
		eType: OK,
		paths: []string{
			"Cluster.Nodes.Master.Instances.*",
			"Cluster.Nodes.Worker.Instances.*",
			"Cluster.Nodes.LoadBalancer.Instances.*",
		},
	},
	// Prevent k8s properties changes
	{
		eType: BLOCK,
		path:  "Kubernetes.Version",
		msg:   "Changing Kubernetes is allowed only when upgrading the cluster.\nTo upgrade the cluster run apply command with '--action upgrade' flag.",
	},
	// Allow addons changes
	{
		eType: OK,
		path:  "Addons",
	},
	// Allow kubitect (project metadata) changes
	{
		eType: OK,
		path:  "Kubitect",
	},
}
