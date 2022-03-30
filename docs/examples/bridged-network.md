# Setting up nodes that use network bridge

> :scroll: **Note:** This example uses `systemd-networkd` to set up the bridge,
but the same result can be achieved with many other approaches.

## (Pre)configure the bridge on the host

In order to use the bridged network, bridge needs to be preconfigured on the host machine.

Create the bridge interface (`br0` in our case) by creating the file `/etc/systemd/network/br0.netdev` 
with the following content:

```editorconfig
[NetDev]
Name=br0
Kind=bridge
```

Bind the ethernet interface (`eno1` in our case) to the bridge interface (`br0` in our case) 
by creating the file `/etc/systemd/network/br0-bind.network` with the following content:

```editorconfig
[Match]
Name=eno1

[Network]
Bridge=br0
```

### Wired adapter using DHCP

To use router's DHCP, instruct `systemd-networkd` to obtain an IPv4 DHCP lease through the bridge interface 
by creating the file `/etc/systemd/network/br0-dhcp.network` with the following content:

```editorconfig
[Match]
Name=br0

[Network]
DHCP=ipv4
```

### Wired adapter using a static IP

To provision virtual machines with static IPs, instruct `systemd-networkd` to obtain an IPv4 DHCP lease through the bridge interface
by creating the file `/etc/systemd/network/br0-static-ip.network` with the following content:

```editorconfig
[Match]
Name=br0

[Network]
Address=192.168.0.10/24
Gateway=192.168.0.1
DNS=192.168.0.1     # Router's DNS 
# DNS=8.8.8.8       # Additional DNS if required
```

In our case, the IP of the server is `192.168.0.10` in the subnet `192.168.0.0/24` 
and the IP of the router (gateway) is `192.168.0.1`.

---

After that, restart `systemd-networkd` and the bridge should be configured:
```sh
systemctl restart systemd-networkd
```

The final step is to disable netfilter on the bridge
(More information can be found [here](https://wiki.libvirt.org/page/Net.bridge.bridge-nf-call_and_sysctl.conf)):

```sh
 cat >> /etc/sysctl.conf <<EOF
 net.bridge.bridge-nf-call-ip6tables = 0
 net.bridge.bridge-nf-call-iptables = 0
 net.bridge.bridge-nf-call-arptables = 0
 EOF
 
 sysctl -p /etc/sysctl.conf
```

## Setting up a cluster over bridged network

In the config file (default is [cluster.yaml](/cluster.yaml)), set the following variables:
- `cluster.network.mode` to `bridge`,
- `cluster.network.bridge` to the name of the bridge you created (`br0` in our case) and
- `cluster.network.gateway` if the first host in `netwrok_cidr` is not a gateway.

```yaml
cluster:
  ...
  network:
    mode: "bridge"
    cidr: "10.10.13.0/24"
    gateway: "10.10.13.1"
    bridge: "br0"
  nodes:
    master:
      instances:
        - id: 1
    worker:
      instances:
        - id: 1
        - id: 2
        - id: 3
```
