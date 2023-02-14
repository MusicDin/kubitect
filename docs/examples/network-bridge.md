Bridged networks allow virtual machines to connect directly to the LAN.
To use Kubitect with bridged network mode, a bridge interface must be preconfigured on the host machine.
This example shows how to configure a simple bridge interface using [Netplan](https://netplan.io/).

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/network-bridge.png" 
    alt="NAT vs bridge network scheme"
    width="75%">
</div>

## Step 1 - (Pre)configure the bridge on the host

Before the network bridge can be created, a name of the host's network interface is required.
This interface will be used by the bridge.

To print the available network interfaces of the host, use the following command.
```sh
nmcli device | grep ethernet
```

Similarly to the previous command, network interfaces can be printed using `ifconfig` or `ip` commands.
Note that these commands output all interfaces, including virtual ones.
```sh
ifconfig -a
# or
ip a
```

Once you have obtained the name of the host's network interface (in our case `eth0`), you can create a bridge interface (in our case `br0`) by creating a file with the following content:
```yaml title="/etc/netplan/bridge0.yaml"
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0: {} # (1)
  bridges:
    br0: # (2)
      interfaces:
        - eth0
      dhcp4: true
      dhcp6: false
      addresses: # (3)
        - 10.10.0.17
```

1. Existing host's ethernet interface to be enslaved.

2. Custom name of the bridge interface.

3. Optionally a static IP address can be set for the bridge interface.

!!! tip "Tip"

    See the official [Netplan configuration examples](https://netplan.io/examples/) for more advance configurations.

Validate if the configuration is correctly parsed by Netplan.
```sh
sudo netplan generate
```

Apply the configuration.
```sh
sudo netplan apply
```

## Step 2 - Disable netfilter on the host

The final step is to prevent packets traversing the bridge from being sent to iptables for processing.
```sh
 cat >> /etc/sysctl.conf <<EOF
 net.bridge.bridge-nf-call-ip6tables = 0
 net.bridge.bridge-nf-call-iptables = 0
 net.bridge.bridge-nf-call-arptables = 0
 EOF
 
 sysctl -p /etc/sysctl.conf
```

!!! tip "Tip"

    For more information, see the [libvirt documentation](https://wiki.libvirt.org/page/Net.bridge.bridge-nf-call_and_sysctl.conf).

## Step 3 - Set up a cluster over bridged network

In the cluster configuration file, set the following variables:

- `cluster.network.mode` to `bridge`,
- `cluster.network.cidr` to the network CIDR of the LAN and
- `cluster.network.bridge` to the name of the bridge you have created (`br0` in our case)

```yaml
cluster:
  network:
    mode: bridge
    cidr: 10.10.13.0/24
    bridge: br0
...
```
