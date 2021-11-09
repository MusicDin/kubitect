# Single node cluster

If you would like to initialize a cluster with only one node than specify only one master node in [terraform.tfvars](../../terraform.tfvars):
```hcl
master_nodes = [
  {
    id  = 1     # Any positive number for id
    ip  = null  # Specific IP or null to auto generate it
    mac = null  # Specific MAC or null to auto generate it
  }
]
```

Don't forget to remove (or comment out) worker and load balancer nodes:
```hcl
lb_nodes = []
...

worker_nodes = []
```

Your master node will now also become a worker node.

*Note: If you do not specify worker nodes, all master nodes will also become worker nodes.*
