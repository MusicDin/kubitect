# Setting up a cluster using bridged network

*Note: This example uses `systemd-networkd` to set up the bridge,
but the same result can be achieved with many other approaches.*

## (Pre)configure the bridge on the host

In order to use the bridged network, bridge needs to be preconfigured on the host machine.

Create the bridge interface (`br0` in our case) by creating the file `/etc/systemd/network/br0.netdev` 
with the following content:
```
[NetDev]
Name=br0
Kind=bridge
```

Bind the ethernet interface (`eno1` in our case) to the bridge interface (`br0` in our case) 
by creating the file `/etc/systemd/network/br0-bind.network` with the following content:
```
[Match]
Name=eno1

[Network]
Bridge=br0
```

### Wired adapter using DHCP

To use router's DHCP, instruct `systemd-networkd` to obtain an IPv4 DHCP lease through the bridge interface 
by creating the file `/etc/systemd/network/br0-dhcp.network` with the following content:
```
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
```
systemctl restart systemd-networkd
```

The final step is to disable netfilter on the bridge
(More information can be found [here](https://wiki.libvirt.org/page/Net.bridge.bridge-nf-call_and_sysctl.conf)):
```bash
 cat >> /etc/sysctl.conf <<EOF
 net.bridge.bridge-nf-call-ip6tables = 0
 net.bridge.bridge-nf-call-iptables = 0
 net.bridge.bridge-nf-call-arptables = 0
 EOF
 
 sysctl -p /etc/sysctl.conf
```

## Setting up a cluster over bridged network

In the [terraform.tfvars](/terraform.tfvars) file, set the variables:
- `network_mode` to `bridge`,
- `network_bridge` to the name of the bridge you created (`br0` in our case) and
- `network_gateway` if the first host in `netwrok_cidr` is not a gateway.

```hcl
worker_nodes = [
  {
    id  = 1
    ip  = "192.168.0.25"  # Static IP
    mac = "52:54:00:00:00:40"
  },
  {
    id  = 2
    ip  = null            # DHCP lease
    mac = "52:54:00:00:00:41"
  }
]
```