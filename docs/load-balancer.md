# Internal load balancing (iLB)

HAProxy load balancers are used to load balance traffic between master nodes.

If you want to deploy [LoadBalancer](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) type services, see the [MetalLB configuration example](examples/metallb.md) configuration example.

---

Specify load balancers in [terraform.tfvars](../terraform.tfvars) file:
```hcl
lb_nodes = [
  {
    id  = 1
    ip  = null  # Specific IP or null to retrieve it from router
    mac = null  # Specific MAC or null to auto generate it
  },
  {
    id  = 2
    ip  = "192.168.113.6"
    mac = "52:54:00:00:00:06"
  }
]
```
*Note: Load balancer IDs must be between 0 and 200, because priority is calculated from ID.*

Then set a floating IP that should not be taken by any other VM:
```hcl
lb_vip = "floating_ip"
```

## Cluster without internal load balancers

If you have only one master node, the internal load balancers are redundant.
In this case, remove all load balancer nodes in the [terraform.tfvars](../terraform.tfvars) file:
```hcl
lb_nodes = []
```

*Note: If multiple master nodes are specified, the IP of the first one is used as the cluster IP.*

## Changing the load balancer configuration BEFORE cluster initialization

To have the same configuration on all your load balancers,
HAProxy configuration has to be changed before initialization.

To do this, modify the [haproxy.tpl](../templates/haproxy/haproxy.tpl) file and insert your custom configuration where the
comment `Place custom configurations here` is located.


## Changing the load balancer configuration via SSH

After the cluster is set up, log in via SSH and change the configuration of the load balancer:
```bash
# SSH into load balancer
ssh <vm_user>@<lb_ip>

# Modify HAProxy configuration file (use your favorite editor)
nano /etc/haproxy/haproxy.cfg

# Test configuration
haproxy -f /etc/haproxy/haproxy.cfg -c

# Apply changes
sudo systemctl restart haproxy
```

*For more information, see the [HAProxy documentation](https://cbonte.github.io/haproxy-dconv/).*
