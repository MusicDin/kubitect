# Single node cluster

If you would like to initialize a cluster with only one node than specify only one master node in [terraform.tfvars](../../terraform.tfvars):
```hcl
vm_master_macs_ips = {
  "mac_for_master_1" = "ip_for_master_1"
}
```

Don't forget to remove (or comment out) worker and load balancer nodes:
```hcl
vm_lb_macs_ips = {
}

...

vm_worker_macs_ips = {
}
```

Your master node will now also become a worker node.

*Note: If you don't specify any worker nodes then all master nodes will also become worker nodes.*
