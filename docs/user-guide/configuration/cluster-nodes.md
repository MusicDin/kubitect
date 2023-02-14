[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

## Nodes configuration structure

Cluster's nodes configuration consists of three **node types**:

- Master nodes (control plane)
- Worker nodes
- Load balancers

For any cluster deployment, **at least one master node needs to be configured**.
Configuring only one master node produces a single-node cluster.
In most cases, a multi-node cluster is desired and therefore worker nodes should be configured as well.

If the control plane of the cluster contains multiple nodes, at least one load balancer must be configured.
Such topology allows the cluster to continue operating normally if any control plane node fails.
In addition, configuring multiple load balancers provides failover in case the primary load balancer fails.

Kubitect currently supports only stacked control plane, which means that etcd key-value stores are deployed on control plane nodes.
Since an etcd cluster requires a majority "(n/2) + 1" of nodes to agree to a change in the cluster, an odd number of nodes (1, 3, 5, ...) provides the best fault tolerance. 
For example, in control planes with 3 nodes, 2 nodes represent the majority, giving a fault tolerance of 1 node. 
In control planes with 4 nodes, the majority is 3 nodes, which provides the same fault tolerance.
For this reason, Kubitect prevents deployment of the cluster whose control plane contains an even number of nodes.

The nodes configuration structure is the following:
```yaml
cluster:
  nodes:
    masters:
      ...
    workers:
      ...
    loadBalancers:
      ...
```

Each node type has two subsections, `default` and `instances`.
Instances represent an array of actual nodes, while defaults provide the configuration that is applied to all instances of a certain node type.
Each default value can also be overwritten by setting the same property for a specific instance.

```yaml
cluster:
  nodes:
    <node-type>:
      default:
        ...
      instances:
        ...
```

## Configuration

### Common node properties

For each instance there is a set of predefined properties that can be set.
Some properties apply for all node types, while some properties are specific for a certain node type.
Properties that apply for all node types, are referred to as *common properties*.

#### Instance ID

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

Each instance in the cluster must have a ID, which can be any positive number and must be unique among all instances within the same node type.
The instance ID is used as a suffix for the name of each node.

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1
        - id: 2
        - id: 77
```

#### CPU

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `2` vCPU

The `cpu` property defines an amount of vCPU cores assigned to the virtual machine.
It can be set for a specific instance or as a default value for all instances.

```yaml
cluster:
  nodes:
    <node-type>:
      default:
        cpu: 2
      instances:
        - id: 1 # (1)!
        - id: 2
          cpu: 4 # (2)!
```

1. Since the `cpu` property is not set for this instance, the default value is used (2).

2. This instance has the `cpu` property set, and therefore the set value (4) overrides the default value (2).

If the property is not set at the instance level or as a default value, Kubitect uses its own default value (2 vCPU).

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1 # (1)!
```

1. Since the 'cpu' property is not set at instance level or as a default value, Kubitect sets the value of the 'cpu' property to **2 vCPU**.

#### RAM 

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `4` GiB

The `ram` property defines an amount of RAM assigned to the virtual machine (in GiB).
It can be set for a specific instance or as a default value for all instances.

```yaml
cluster:
  nodes:
    <node-type>:
      default:
        ram: 8
      instances:
        - id: 1 # (1)!
        - id: 2
          ram: 16 # (2)!
```

1. Since the `ram` property is not set for this instance, the default value is used (8 GiB).

2. This instance has the `ram` property set, and therefore the set value (16 GiB) overrides the default value (8 GiB).

If the property is not set at the instance level or as a default value, Kubitect uses its own default value (4 GiB).

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1 # (1)!
```

1. Since the `ram` property is not set at instance level or as a default value, Kubitect sets the value of the `ram` property to **4 GiB**.

#### Main disk size

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `32` GiB

The `mainDiskSize` property defines an amount of space assigned to the virtual machine (in GiB).
It can be set for a specific instance or as a default value for all instances.

```yaml
cluster:
  nodes:
    <node-type>:
      default:
        mainDiskSize: 128
      instances:
        - id: 1 # (1)!
        - id: 2
          mainDiskSize: 256 # (2)!
```

1. Since the `mainDiskSize` property is not set for this instance, the default value is used (128 GiB).

2. This instance has the `mainDiskSize` property set, so therefore the set value (256 GiB) overrides the default value (128 GiB).

If the property is not set at the instance level or as a default value, Kubitect uses its own default value (32 GiB).

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1 # (1)!
```

1. Since the `mainDiskSize` property is not set at instance level or as a default value, Kubitect sets the value of the `mainDiskSize` property to **32 GiB**.

#### IP address

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

For each node a static IP address can be set.
f no IP address is set for a particular node, a DHCP lease is requested. Kubitect also checks whether all set IP addresses are within the defined network range (see [Network CIDR](../cluster-network/#network-cidr)).

```yaml
cluster:
  network:
    mode: nat
    cidr: 192.168.113.0/24
  nodes:
    <node-type>:
      instances:
        - id: 1
          ip: 192.168.113.5 # (1)!
        - id: 2 # (2)!
```

1. A static IP (`192.168.113.5`) is set for this instance.

2. Since no IP address is defined for this instance, a DHCP lease is requested.

#### MAC address

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

By default, MAC addresses are generated for each virtual machine created, but a custom MAC address can also be set.

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1
          mac: "52:54:00:00:13:10" # (1)!
        - id: 2 # (2)!
```

1. A custom MAC address (`52:54:00:00:13:10`) is set for this instance.

2. Since no MAC address is defined for this instance, the MAC address is generated during cluster creation.

#### Host affinity

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

By default, all instances are deployed on the *[*default host*](../hosts/#default-host)*.
Kubitect can be instructed to deploy the instance on a specific host by specifying the name of the host in the instance configuration.

```yaml

hosts:
  - name: host1
    ...
  - name: host2
    default: true
    ...

cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1
          host: host1 # (1)!
        - id: 2 # (2)!
```

1. The instance is deployed on `host1`.

2. Since no host is specified, the instance is deployed on the default host (`host2`).

### Control plane and worker node properties

The following properties can be configured only for control plane or worker nodes.

#### Data disks

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]

By default, only a main disk (volume) is attached to each provisioned virtual machine.
Since the main disk already contains an operating system, so it may not be suitable for storing data.
Therefore, additional disks might be required.
For example, a [Rook](https://rook.io/) can be easily configured to use all the empty disks attached to the virtual machine to form a storage cluster.

A name and size (in GiB) must be configured for each data disk.
By default, data disks are created in the main resource pool.
To create a data disk in a custom [data resource pool](../hosts/#data-resource-pools), the pool property can be set to the name of the desired data resource pool.
Also note that the data disk name must be unique among all data disks for a given instance.

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1
          dataDisks:
            - name: data-volume
              pool: main # (1)!
              size: 256
            - name: rook-volume
              pool: rook-pool # (2)!
              size: 512
```

1. When `pool` property is omitted or set to `main`, the data disk is created in the main resource pool.

2. Custom [data resource pool](../hosts/#data-resource-pools) must be configured in the hosts section.

!!! note "Note"

    Default data disks are currently not supported.


#### Node labels

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

Node labels are configured as a dictionary of key-value pairs.
They are used to label actual Kubernetes nodes, and therefore can only be applied to control plane (master) and worker nodes.

They can be set for a specific instance or as a default value for all instances.
Labels set for a specific instance are merged with the default labels.
However, labels set at the instance level take precedence over default labels.

```yaml
cluster:
  nodes:
    <node-type>: # (1)!
      default:
        labels:
          label-key-1: def-label-value-1
          label-key-2: def-label-value-2
      instances:
        - id: 1
          labels: # (2)!
            label-key-3: instance-label-value-3
        - id: 2
          labels: # (3)!
            label-key-1: new-label-value-1
```

1. Node labels can only be applied to **control plane** (master) and **worker** nodes.

2.  Labels defined at the instance level are merged with default labels.
    As a result, the following labels are applied to this instance:

    - `#!yaml label-key-1: def-label-value-1`
    - `#!yaml label-key-2: def-label-value-2`
    - `#!yaml label-key-3: instance-label-value-3`

3.  Labels defined at the instance level take precedence over default labels.
    As a result, the following labels are applied to this instance:

    - `#!yaml label-key-1: new-label-value-1`
    - `#!yaml label-key-2: def-label-value-2`

#### Node taints

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]

Node taints are configured as a list of strings in the format `key=value:effect`.
Similar to node labels, taints can only be applied to control plane (master) and worker nodes.

Taints can be set for a specific instance or as a default value for all instances.
Taints set for a particular instance are merged with the default taints and duplicate entries are removed.

```yaml
cluster:
  nodes:
    <node-type>: # (1)!
      default:
        taints:
          - "key1=value1:NoSchedule"
      instances:
        - id: 1
          taints:
            - "key2=value2:NoExecute"
```

1. Node taints can only be applied to **control plane** (master) and **worker** nodes.

### Load balancer properties

The following properties can be configured only for load balancers.

#### Virtual IP address (VIP)

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

Load balancers distribute traffic directed to the control plane across all master nodes.
Nevertheless, a load balancer can fail and make the control plane unreachable.
To avoid such a situation, multiple load balancers can be configured.
They work on the failover principle, i.e. one of them is primary and actively serves incoming traffic, while others are secondary and take over the primary position only if the primary load balancer fails.
If one of the secondary load balancers becomes primary, it should still be reachable via the same IP.
This IP is usually referred to as a virtual or floating IP (VIP).

VIP must be specified if multiple load balancers are configured.
It must also be an unused host IP address within the configured network.


```yaml
cluster:
  nodes:
    loadBalancer:
      vip: 168.192.113.200
```

#### Virtual router ID

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]
&ensp;
:octicons-file-symlink-file-24: Default: `51`

When a cluster is created with a virtual IP (VIP) set, Kubitect configures the virtual router redundancy protocol (VRRP), which provides failover for load balancers.
A virtual router ID (VRID) identifies the group of VRRP routers.
Each group has its own ID.
Since there can be only one master in each group, two groups cannot have the same ID.

The virtual router ID can be any number between 0 and 255.
By default, Kubitect sets the virtual router ID to `51`.
If you set up multiple clusters that use VIP, you must ensure that the virtual router ID is different for each cluster.


```yaml
cluster:
  nodes:
    loadBalancer:
      vip: 168.192.113.200
      virtualRouterId: 30 # (1)!
```

1. If the virtual IP (VIP) is not set, the virtual router ID is ignored.

#### Priority

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]
&ensp;
:octicons-file-symlink-file-24: Default: `10`

Each load balancer has a priority that is used to select a primary load balancer.
The one with the highest priority becomes the primary and all others become secondary.
If the primary load balancer fails, the next one with the highest priority takes over.
If two load balancers have the same priority, the one with the higher sum of IP address digits is selected.

The priority can be any number between 0 and 255.
The default priority is 10.

```yaml
cluster:
  nodes:
    loadBalancer:
      instances:
        - id: 1 # (1)!
        - id: 2 
          priority: 200 # (2)!
```

1. Since the load balancer priority for this instance is not specified, it is set to 10.

2. Since this load balancer instance has the highest priority (200 > 10), it becomes the primary load balancer.


#### Port forwarding

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

By default, all configured load balancers distribute incoming traffic on port 6443 across all control plane nodes.
Kubitect allows additional user-defined ports to be configured.

The following properties can be configured for each port:

+ `name` - Name is a unique port identifier.
+ `port` - Incoming port is a port on which the load balancer listens for incoming traffic.
+ `targetPort` - Target port is a port where traffic is forwarded by the load balancer.
+ `target` - Target is a group of nodes to which traffic is forwarded. Possible targets are:
    - `masters` - control plane nodes
    - `workers` - worker nodes 
    - `all` - control plane and worker nodes.

A unique name and a unique incoming port must be configured for each port.
The configuration of target and target port is optional.
If target port is not configured, it is set to the same value as the incoming port.
If target is not configured, incoming traffic is distributed across worker nodes by default.

```yaml
cluster:
  nodes:
    loadBalancer:
      forwardPorts:
        - name: https
          port: 443 # (1)!
          targetPort: 31200 # (2)!
          target: all # (3)!
```

1.  Incoming port is the port on which a load balancer listens for incoming traffic.
    It can be any number between 1 and 65353, excluding ports 6443 (Kubernetes API server) and 22 (SSH).

2.  Target port is the port on which the traffic is forwarded. 
    By default, it is set to the same value as the incoming port.

3.  Target represents a group of nodes to which incoming traffic is forwarded.
    Possible values are:

    - `masters`
    - `workers`
    - `all`

    If the target is not configured, it defaults to the `workers`.



## Example usage

### Set a role to all worker nodes

By default worker nodes have no roles (`<none>`).
For example, to set `node` role to all worker nodes in the cluster, set default label with key `node-role.kubernetes.io/node`.

```yaml
cluster:
  nodes:
    worker:
      default:
        labels:
          node-role.kubernetes.io/node: # (1)!
      instances:
        ...
```

1. If the label value is omitted, `null` is set as the label value.

Node roles can be seen by listing cluster nodes with `kubectl`.

```
NAME                         STATUS   ROLES                  AGE   VERSION
local-k8s-cluster-master-1   Ready    control-plane,master   19m   v1.22.6
local-k8s-cluster-worker-1   Ready    node                   19m   v1.22.6
local-k8s-cluster-worker-2   Ready    node                   19m   v1.22.6
```

### Load balance HTTP requests

Kubitect allows users to define custom port forwarding on load balancers.
For example, to distribute HTTP and HTTPS requests across all worker nodes, at least one load balancer has to be specified and port forwarding must be configured, as shown in the sample configuration below.

```yaml
cluster:
  nodes:
    loadBalancer:
      forwardPorts:
        - name: http
          port: 80
        - name: https
          port: 443
          target: all # (1)!
      instances:
        - id: 1
```

1. When the target is set to `all`, load balancers distribute traffic across all nodes (master and worker nodes).