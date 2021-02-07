# Load balancing

HAProxy load balancer is used to load balance traffic between master nodes.

If you would like to expose services of type [LoadBalancer](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) than check [MetalLB](https://metallb.universe.tf/) project.

## Cluster without load balancer

If you decide to omit load balancer, all you have to do is to modify [terraform.tfvars](../terraform.tfvars) file.

Remove all load balancers IP and MAC addresses:
```hcl
vm_lb_macs_ips = {}
```

*Note: If there is more master nodes specified, IP of the first one will be used for a cluster IP.*

## Cluster with load balancer(s)

*Note: This script supports up to 2 HAProxy load balancers.*

Provide a MAC and IP address for each load balancer in [terraform.tfvars](../terraform.tfvars) file:
```hcl
vm_lb_macs_ips = {
  "mac_for_lb_1" = "ip_for_lb_1"
  "mac_for_lb_2" = "ip_for_lb_2"
}
```

Then set a floating IP that should not be taken by any other VM:
```hcl
vm_lb_vip = "floating_ip"
```

## Modifying load balancer's configuration BEFORE cluster initialization

In order to have the same configuration on all of your load balancers,
HAProxy configuration has to be modified before initialization.

To accomplish that, modify [haproxy.cfg](../templates/haproxy.tpl) file and put your custom configuration where
comment `Place custom configurations here` is located.

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
sudo systemctl restart haproxy
```

*For more information check [HAProxy documentation](https://cbonte.github.io/haproxy-dconv/)*.
