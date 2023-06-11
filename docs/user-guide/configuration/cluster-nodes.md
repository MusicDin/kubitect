[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0
[tag 2.3.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.3.0

<div markdown="1" class="text-center">
# Cluster nodes
</div>

<div markdown="1" class="text-justify">

## Background

Kubitect allows configuration of three distinct **node types**: worker nodes, master nodes (control plane), and load balancers.

### Worker nodes

Worker nodes in a Kubernetes cluster are responsible for executing the application workloads of the system. The addition of more worker nodes to the cluster enhances redundancy in case of worker node failure. However, allocating more resources to each worker node provides less overhead and more resources for the actual applications.

Kubitect does not offer automatic scaling of worker nodes based on resource demand. However, you can easily add or remove worker nodes by applying modified cluster configuration.

### Master nodes

The master node plays a vital role in a Kubernetes cluster as it manages the overall state of the system and coordinates the workloads running on the worker nodes.
Therefore, it is essential to **configure at least one master node for every cluster**.

Please note that Kubitect currently supports only a stacked control plane where etcd key-value stores are deployed on control plane nodes.
To ensure the best possible fault tolerance, it is important to configure an odd number of control plane nodes.
For more information, please refer to the [etcd FAQ](https://etcd.io/docs/v3.4/faq/#why-an-odd-number-of-cluster-members).

### Load balancer nodes

In a Kubernetes cluster with multiple control plane nodes, it is necessary to configure at least one load balancer.
A load balancer distributes incoming network traffic across multiple control plane nodes, ensuring the cluster operates normally even if any control plane node fails.

However, configuring only one load balancer represents a single point of failure for the cluster.
If it fails, incoming traffic will not be distributed to the control plane nodes, potentially resulting in downtime.
Therefore, configuring multiple load balancers is essential to ensure high availability for the cluster.

## Nodes configuration structure

The configuration structure for the nodes is as follows:

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

Each node type has two subsections: `default` and `instances`.
The instances subsection represents an array of actual nodes, while the default subsection provides the configuration that is applied to all instances of a particular node type.
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

Each node instance has a set of predefined properties that can be set to configure its behavior.
Some properties apply to all node types, while others are specific to a certain node type.
Properties that apply to all node types are referred to as *common properties*.

#### Instance ID

:material-tag-arrow-up-outline: [v2.3.0][tag 2.3.0]
&ensp;
:material-alert-circle-outline: Required

Each node in a cluster must have a unique identifier, or ID, that distinguishes it from other instances of the same node type.
The instance ID is used as a suffix for the name of each node, ensuring that each node has a unique name in the cluster.

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1
        - id: compute-1
        - id: 77
```

#### CPU

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `2` vCPU

The `cpu` property defines the amount of virtual CPU cores assigned to a node instance.
This property can be set for a specific instance, or as a default value for all instances of a certain node type.

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

1. Since the `cpu` property is not set for this instance, the default value (2) is used.

2. This instance has the `cpu` property set, and therefore the set value (4) overrides the default value (2).

If the property is not set at the instance level or as a default value, Kubitect uses its own default value (2).

```yaml
cluster:
  nodes:
    <node-type>:
      instances:
        - id: 1 # (1)!
```

1. Since the `cpu` property is not set at instance level or as a default value, Kubitect sets the value of the `cpu` property to **2 vCPU**.

#### RAM

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `4` GiB

The `ram` property defines the amount of RAM assigned to a node instance (in GiB).
This property can be set for a specific instance, or as a default value for all instances of a certain node type.

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

1. Since the `ram` property is not set for this instance, the default value (8 GiB) is used.

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

The `mainDiskSize` property defines the amount of disk space assigned to a node instance (in GiB).
This property can be set for a specific instance, or as a default value for all instances of a certain node type.

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

1. Since the `mainDiskSize` property is not set for this instance, the default value (128 GiB) is used.

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

Each node in a cluster can be assigned a static IP address to ensure a predictable and consistent IP address for the node.
If no IP address is set for a particular node, Kubitect will request a DHCP lease for that node.
Additionally, Kubitect checks whether all set IP addresses are within the defined network range, as explained in the [Network CIDR](../cluster-network/#network-cidr) section of the cluster network configuration.

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

The virtual machines created by Kubitect are assigned generated MAC addresses, but a custom MAC address can be set for a virtual machine if necessary.

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

By default, all instances in a cluster are deployed on the [default host](../hosts/#default-host).
However, by specifying a specific host for an instance, you can control where that instance is deployed

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

The following properties can only be configured for control plane or worker nodes.

#### Data disks

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]

By default, only a main disk (volume) is attached to each provisioned virtual machine. Since the main disk already contains an operating system, it may not be suitable for storing data, and additional disks may be required. For example, [Rook](https://rook.io/) can be easily configured to use all the empty disks attached to the virtual machine to form a storage cluster.

A name and size (in GiB) must be configured for each data disk. By default, data disks are created in the main resource pool. To create a data disk in a custom [data resource pool](../hosts/#data-resource-pools), you can set the pool property to the name of the desired data resource pool. Additionally, note that the data disk name must be unique among all data disks for a given instance.


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


#### Node labels

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

With node labels, you can help organize and manage your cluster by associating nodes with specific attributes or roles, and by grouping nodes for specific workloads or tasks.

Node labels are used to label actual Kubernetes nodes and can be set for a specific instance or as a default value for all instances.
It is important to note that labels set at the instance level are merged with the default labels.
However, if labels have the same key, then the labels set at the instance level take precedence over the default labels.

```yaml
cluster:
  nodes:
    <node-type>: # (1)!
      default:
        labels:
          key1: def-value-1
          key2: def-value-2
      instances:
        - id: 1
          labels: # (2)!
            key1: custom-value
        - id: 2
          labels: # (3)!
            key3: super-node
```

1. Node labels can only be applied to **worker** and **master** (control plane) nodes.


2.  Labels defined at the instance level take precedence over default labels.
    As a result, the following labels are applied to this instance:

    - `#!yaml key1: custom-value`
    - `#!yaml key2: def-value-2`

3.  Labels defined at the instance level are merged with default labels.
    As a result, the following labels are applied to this instance:

    - `#!yaml key1: def-value-1`
    - `#!yaml key2: def-value-2`
    - `#!yaml key3: super-node`


#### Node taints

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]

With node taints, you can limit which pods can be scheduled to run on a particular node, and help ensure that the workload running on that node is appropriate for its capabilities and resources.

Node taints are configured as a list of strings in the format `key=value:effect`.
Taints can be set for a specific instance or as a default value for all instances.
When taints are set for a particular instance, they are merged with the default taints, and any duplicate entries are removed.


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

The following properties can only be configured for load balancers.

#### Virtual IP address (VIP)

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

??? question "What is VIP? <i class="click-tip"></i>"

    Load balancers are responsible for distributing traffic to the control plane nodes.
    However, a single load balancer can cause issues if it fails.
    To avoid this, multiple load balancers can be configured with one as the primary, actively serving incoming traffic, while others act as secondary and take over the primary position only if the primary load balancer fails.
    If a secondary load balancer becomes primary, it should still be reachable via the same IP, which is referred to as a virtual or floating IP (VIP).

When multiple load balancers are configured, an unused IP address within the configured network must be specified as the VIP.


```yaml
cluster:
  nodes:
    loadBalancer:
      vip: 168.192.113.200
```

#### Virtual router ID (VRID)

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]
&ensp;
:octicons-file-symlink-file-24: Default: `51`

When a cluster is created with a VIP, Kubitect configures Virtual Router Redundancy Protocol (VRRP), which provides failover for load balancers.
Each VRRP group is identified by a virtual router ID (VRID), which can be any number between 0 and 255.
Since there can be only one master in each group, two groups cannot have the same ID.

By default, Kubitect sets the VRID to 51, but if you set up multiple clusters that use VIP, you must ensure that the VRID is different for each cluster.


```yaml
cluster:
  nodes:
    loadBalancer:
      vip: 168.192.113.200
      virtualRouterId: 30
```

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


By default, each configured load balancer has a port forwarding rule that distribute incoming traffic on port 6443 across the available control plane nodes.
However, Kubitect provides the flexibility to configure additional user-defined port forwarding rules.

The following properties can be configured for each rule:

+ `name` - A unique port identifier.
+ `port` - The incoming port on which the load balancer listens for traffic.
+ `targetPort` - The port to which traffic is forwarded by the load balancer.
+ `target` - The group of nodes to which traffic is directed. The possible targets are:
    - `masters` - control plane nodes
    - `workers` - worker nodes
    - `all` - worker and control plane nodes.

Every port forwarding rule must be configured with a unique `name` and `port`.
The name serves as a unique identifier for the rule, while the port specifies the incoming port on which the load balancer listens for traffic.

The `target` and `targetPort` configurations are optional.
If target port is not explicitly set, it will default to the same value as the incoming port.
Similarly, if target is not set, incoming traffic is automatically distributed across worker nodes.

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

By default, worker nodes in a Kubernetes cluster are not assigned any roles (`<none>`).
To set the role of all worker nodes in the cluster, the default label with the key `node-role.kubernetes.io/node` can be configured.

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

The roles of the nodes in a Kubernetes cluster can be viewed using `kubectl get nodes`.

```
NAME                   STATUS   ROLES                  AGE   VERSION
k8s-cluster-master-1   Ready    control-plane,master   19m   v1.26.5
k8s-cluster-worker-1   Ready    node                   19m   v1.26.5
k8s-cluster-worker-2   Ready    node                   19m   v1.26.5
```

### Load balance HTTP requests

Kubitect enables users to define custom port forwarding rules on load balancers.
For example, to distribute HTTP and HTTPS requests across all worker nodes, at least one load balancer must be specified and port forwarding must be configured as follows:

```yaml
cluster:
  nodes:
    loadBalancer:
      forwardPorts:
        - name: http
          port: 80
        - name: https
          port: 443
      instances:
        - id: 1
```

</div>
