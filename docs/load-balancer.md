# Load balancing

HAProxy load balancer is used to load balance traffic between master nodes.

If you would like to expose services of type [LoadBalancer](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) than check [MetalLB](https://metallb.universe.tf/) project. 

## Cluster without load balancer

If you decide to omit load balancer, all you have to do is to modify [terraform.tfvars](../terraform.tfvars) file.

First remove all load balancer's IP and MAC addresses:
```
vm_lb_macs = {}

vm_lb_ips = {}
``` 

Then set a floating IP to point on the master node:
<pre>
vm_lb_vip = "<b>master_node_IP</b>"
</pre>

## Cluster with load balancer(s)

*Note: This script supports up to 2 load balancers.*

Provide MAC and IP address for your load balancer(s) in [terraform.tfvars](../terraform.tfvars) file:
```
vm_lb_macs_ips = {
  "0" = "ip_for_lb_1"
  "1" = "ip_for_lb_2"
}

vm_lb_macs = {
  "0" = "mac_for_lb_1"
  "1" = "mac_for_lb_2"
}
``` 

Then set a floating IP that should not be taken by any other VM:
```
vm_lb_vip = "floating_ip"
```

## Modifying load balancer's configuration BEFORE cluster initialization

In order to have the same configuration on all of your load balancers, 
HAProxy configuration has to be modified before initialization.

To accomplish that, modify [haproxy.cfg](../templates/haproxy.tpl) file and put your custom configuration where 
comment `Place custom configurations here` is located. 

*For more information check [HAProxy documentation](https://cbonte.github.io/haproxy-dconv/)*.

## Modifying load balancer's configuration over SSH

After the cluster is all set up, SSH into it and modify it's configuration:
```bash
# SSH into load balancer
ssh <vm_user>@<vm_lb_ip> 

# Modify LB's configuration file (use your favorite editor)
nano /etc/haproxy/haproxy.cfg

# Test configuration
haproxy -f /etc/haproxy/haproxy.cfg -c

# Apply changes
haproxy -f /etc/haproxy/haproxy.cfg
``` 

*For more information check [HAProxy documentation](https://cbonte.github.io/haproxy-dconv/)*.