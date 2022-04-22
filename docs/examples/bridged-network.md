<h1 align="center">Setting up nodes over bridged network</h1>

This example shows how to configure a simple bridge interface using [netplan](https://netplan.io/).

## Step 1 - (Pre)configure the bridge on the host

In order to use the bridged network, bridge interface needs to be preconfigured on the host machine.

Create the bridge interface (`br0` in our case) by creating a file with the following content:
```yaml title="/etc/netplan/bridge0.yaml"
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0: {}       # (1)
  bridges:
    br0:           # (2)
      interfaces:
        - eth0
      dhcp4: true
      dhcp6: false
      addresses:   # (3)
        - 10.10.0.17
```

1. Existing ethernet interface to be enslaved.

2. Custom name of the bridge interface.

3. Optionally a static IP address can be set for the bridge interface.

!!! tip "Tip"

    See the official [netplan configuration examples](https://netplan.io/examples/) for more complex configurations.

Validate if the configuration is correctly parsed by netplan.
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
- `cluster.network.bridge` to the name of the bridge you have created (`br0` in our case) and
- `cluster.network.gateway` if the first host in `netwrok_cidr` is not a gateway.

```yaml
cluster:
  network:
    mode: "bridge"
    cidr: "10.10.13.0/24"
    gateway: "10.10.13.1"
    bridge: "br0"
...
```
