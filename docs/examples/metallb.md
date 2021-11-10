# MetalLB configuration examples

MetalLB allows you to expose services of type 'LoadBalancer'.

The official MetalLB documentation can be found [here](https://www.metallb.org).

MetalLB can be configured in two modes:
- [Layer2](#layer2)
- [BGP](#bgp)

# Layer2

`layer2` mode is easier to configure. 
You only need to configure the MetalLB section in [terraform.tfvar](/terraform.tfvars):
```hcl
metallb_enabled  = "true"
metallb_protocol = "layer2"
# MetalLB IP range should be defined from your network IP range
metallb_ip_range = "192.168.113.241-192.168.113.254" 
``` 

# BGP

To configure MetalLB in BGP mode, you need a BGP-capable router 
(virtual routers also work).

Let us assume the following cluster configuration:
- Cluster network: `192.168.113.0/24`
- Network gateway: `192.168.113.1/32` 
- Master node: `192.168.113.10/32`
- Worker node: `192.168.113.40/32`

### Routing - FRR

This example shows you how to configure [FRR](https://frrouting.org/) on your host.

First, [install FRR](http://docs.frrouting.org/en/latest/installation.html) on the host machine.

Then enable the BGP daemon by setting `bgpd` to *yes* in `/etc/frr/daemons`:
```
bgpd=yes
```

Edit the FRR configuration file `/etc/frr/frr.conf`:
```
!
frr defaults traditional
hostname localhost
log file /var/log/frr/frr.log
no ipv6 forwarding
!
router bgp 65000
 coalesce-time 1000
 neighbor 192.168.113.10 remote-as 65000
 neighbor 192.168.113.40 remote-as 65000
!
line vty
!
```

*Note: It is recommended to use the [vty shell](http://docs.frrouting.org/en/latest/vtysh.html) to edit the FRR configuration.*

Now enable FRR on your host:
```
systemctl --now enable frr
```

### Firewall

This example shows how to configure [firewalld](https://firewalld.org/) to allow BGP traffic.

To allow our FRR router to communicate with the cluster nodes, allow BGP in the firewall:
```
firewall-cmd --permanent --add-service=bgp
firewall-cmd --reload
```

Alternatively, you can specify a zone in which to allow BGP:
```
firewall-cmd --permanent --zone=<ZONE> --add-service=bgp
```

### Configure MetalLB

In the [terraform.tfvars](/terraform.tfvars) file, configure the MetalLB section:
```hcl
metallb_enabled  = "true"
metallb_protocol = "bgp"
# MetalLB IP range should be contained in our network space (192.168.113.0/24)
metallb_ip_range = "192.168.113.241-192.168.113.254"
metallb_peers = [{
  peer_ip  = "192.168.113.1"
  peer_asn = 65000
  my_asn   = 65000
}]
```