[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<div markdown="1" class="text-center">
# Cluster network
</div>

<div markdown="1" class="text-justify">

Network section of the Kubitect configuration file defines the properties of the network to be created or the network to which the cluster nodes are to be assigned.

## Configuration

### Network mode

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

Kubitect supports two network modes: NAT and bridge.


```yaml
cluster:
  network:
    mode: nat
```

#### NAT mode

In NAT (Network Address Translation) mode, the libvirt virtual network is created for the cluster, which reduces the need for manual configurations. 
However, it's limited to a single host, i.e., a single physical server.

#### Bridge mode

In bridge mode, a real host network device is shared with the virtual machines, allowing each virtual machine to bind to any available IP address on the local network, just like a physical computer. 
This approach makes the virtual machine visible on the network, enabling the creation of clusters across multiple physical servers.

To use bridged networks, you need to preconfigure the bridge interface on each target host. 
This is necessary because each environment is unique. For instance, you might use link aggregation (also known as link bonding or teaming), which cannot be detected automatically and therefore requires manual configuration. 
The [Network bridge example](../../examples/network-bridge.md) provides instructions on how to create a bridge interface with netplan and configure Kubitect to use it.

### Network CIDR

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

The network CIDR (Classless Inter-Domain Routing) represents the network in the form of `<network_ip>/<network_prefix_bits>`.
All IP addresses specified in the cluster section of the configuration must be within this network range, including the network gateway, node instances, floating IP of the load balancer, and so on.

In NAT network mode, the network CIDR defines an unused private network that is created. In bridge mode, the network CIDR should specify the network to which the cluster belongs.


```yaml
cluster:
  network:
    cidr: 192.168.113.0/24 # (1)!
```

1.  In `nat` mode - Any unused private network within a local network.

    In `bridge` mode - A network to which the cluster belongs.


### Network gateway

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

The network gateway, also known as the default gateway, represents the IP address of the router. 
By default, it doesn't need to be specified, as the first client IP in the network range is used as the gateway address. However, if the gateway IP differs from this, it must be specified manually.

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

The network bridge determines the bridge interface that virtual machines connect to.

In NAT network mode, a virtual network bridge interface is created on the host. 
These bridges are usually prefixed with *vir*, such as *virbr44*. 
If you omit this option, the virtual bridge name is automatically determined by libvirt. 
Alternatively, you can specify the name to be used for the virtual bridge.

In bridge network mode, the network bridge should be the name of the preconfigured bridge interface, such as *br0*.

```yaml
cluster:
  network:
    bridge: br0
```

## Example usage

### Virtual NAT network

If the cluster is created on a single host, you can use the NAT network mode.
In this case, you only need to specify the CIDR of the new network in addition to the network mode.

```yaml
cluster:
  network:
    mode: nat
    cidr: 192.168.113.0/24
```

### Bridged network 

To make the cluster nodes visible on the local network as physical machines or to create the cluster across multiple hosts, you must use bridge network mode. 
Additionally, you need to specify the network CIDR of an existing network along with the preconfigured host bridge interface.

```yaml
cluster:
  network:
    mode: bridge 
    cidr: 10.10.64.0/24 
    bridge: br0 
```

</div>
