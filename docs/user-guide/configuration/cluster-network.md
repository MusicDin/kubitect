[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<h1 align="center">Cluster network</h1>

This document describes **how to define the cluster network** in the Kubitect configuration.
It defines either the properties of the network to be created or the network to which the cluster nodes are to be assigned.


## Configuration

### Network mode

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

Kubitect supports two network modes. 
The first is the `nat` mode and the other is the `bridge` mode.


```yaml
cluster:
  network:
    mode: nat
```

#### NAT mode

In NAT (Network Address Translation) mode, the libvirt virtual network is created for the cluster.
This reduces the need for manual configurations, but is limited to one host (a single physical server).

#### Bridge mode

In bridge mode, a real host network device is shared with the virtual machines.
Therefore, each virtual machine can bind to any available IP address on the local network, just like a physical computer.
This approach makes the virtual machine visible on the network, which enables the creation of clusters across multiple physical servers.

The only requirement for using bridged networks is the preconfigured bridge interface on each target host.
Preconfiguring bridge interfaces is necessary because every environment is different.
For example, someone might use link aggregation (also known as link bonding or teaming), which cannot be detected automatically and therefore requires manual configuration.
The [Network bridge](../../examples/network-bridge.md) example describes how to create a bridge interface with netplan and configure Kubitect to use it.

!!! question "How to automate bridge interface creation?"

    If you have an idea or suggestion on how to automate the creation of  bridge interfaces, feel free to [open an issue](https://github.com/MusicDin/kubitect/issues/new) on GitHub.

### Network CIDR

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

Network CIDR (Classless Inter-Domain Routing) defines the network in a form of `<network_ip>/<network_prefix_bits>`. All IP addresses specified in the cluster section of the configuration must be in this network range (*this includes a network gateway, node instances, floating IP of the load balancer, etc.*).

When using NAT network mode, the network CIDR defines an unused private network that is created. 
In bridge mode, the network CIDR should specify the network to which the cluster belongs.


```yaml
cluster:
  network:
    cidr: 192.168.113.0/24 # (1)!
```

1.  In `nat` mode - Any unused private network within a local network.

    In `bridge` mode - A network to which the cluster belongs.


### Network gateway

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

The network gateway (or default gateway) defines the IP address of the router.

By default, it does not need to be specified because the first client IP in the network range is used as the gateway address.
If the gateway IP differs from this, it must be specified manually.

Also note that the gateway IP address must be within the specified network range.

```yaml
cluster:
  network:
    cidr: 10.10.0.0/20
    gateway: 10.10.0.230 # (1)!
```

1. If this option is omitted, `10.10.0.1` is used as the gateway IP (first client IP in the network range).


### Network bridge

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `virbr0`

The network bridge defines the bridge interface that virtual machines connect to.

When the NAT network mode is used, a virtual network bridge interface is created on the host.
Virtual bridges are usually prefixed with `vir` (example: `virbr44`).
If this option is omitted, the virtual bridge name is automatically determined by libvirt.
Otherwise, the specified name is used for the virtual bridge.

In the case of bridge network mode, the network bridge should be the name of the preconfigured bridge interface (example: `br0`).

```yaml
cluster:
  network:
    bridge: br0
```

## Example usage

### Virtual NAT network

If the cluster is created on a single host, the NAT network mode can be used.
In this case, only the CIDR of the new network needs to be specified in addition to the network mode.

```yaml
cluster:
  network:
    mode: nat
    cidr: 192.168.113.0/24
```

### Bridged network 

To make the cluster nodes visible on the local network as physical machines or to create the cluster across multiple hosts, bridge network mode must be used.
Also, the network CIDR of an existing network must be specified along with the preconfigured host bridge interface.

```yaml
cluster:
  network:
    mode: bridge 
    cidr: 10.10.64.0/24 
    bridge: br0 
```
