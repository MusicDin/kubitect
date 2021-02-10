# MetalLB configuration examples

MetalLB allows you to expose your service as a type of `LoadBalancer`.

Official MetalLB documentation is available [here](https://www.metallb.org).

MetalLB can be configured in two modes:
- [Layer2](#layer2)
- [BGP](#bgp)

# Layer2

`layer2` mode is easier to configure. All you have to do is just configure the MetalLB section in [terraform.tfvar](/terraform.tfvars):
```hcl
metallb_enabled  = "true"
metallb_protocol = "layer2"
# MetalLB IP range should be defined from your network IP range
metallb_ip_range = "192.168.113.241-192.168.113.254" 
``` 

# BGP

To configure MetalLB in BGP mode, you will need BGP enabled router on your host.

Let's assume the following cluster configuration:
- Cluster's network: `192.168.113.1/24`
- Master node: `192.168.113.10/32`
- Worker node: `192.168.113.40/32`

### Routing - FRR

This example will show how to configure [FRR](https://frrouting.org/) on your host.

First [install FRR](http://docs.frrouting.org/en/latest/installation.html) on the host machine.

Then enable BGP daemon in `/etc/frr/daemons` by setting `bgpd` from *no* to *yes*:
```
bgpd=yes
```

Edit FRR configuration file `/etc/frr/frr.conf`:
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

*Note: It's recommended to use [vty shell](http://docs.frrouting.org/en/latest/vtysh.html) for editing FRR configuration.*

Now enable FRR on your host:
```
systemctl --now enable frr
```

### Firewall

This example shows how to configure [firewalld](https://firewalld.org/) to allow BGP traffic.

To allow our FRR router to communicate with cluster nodes, allow BGP in firewall:
```
firewall-cmd --permanent --add-service=bgp
firewall-cmd --reload
```

Alternatively you can also specify a zone on which BGP will be enabled:
```
firewall-cmd --permanent --zone=<ZONE> --add-service=bgp
```

### Configuring MetalLB

In [terraform.tfvars](/terraform.tfvars) configure MetalLB section:
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