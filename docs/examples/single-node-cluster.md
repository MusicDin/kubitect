# Single node cluster

If you want to initialize a cluster with only one node,
specify only one master node in [terraform.tfvars](../../terraform.tfvars):
```hcl
master_nodes = [
  {
    id  = 1     # Any positive number for id
    ip  = null  # Specific IP or null to auto generate it
    mac = null  # Specific MAC or null to auto generate it
  }
]
```

Do not forget to remove (or comment out) the worker and load balancer nodes:
```hcl
lb_nodes = []
...

worker_nodes = []
```

Your master node now also becomes a worker node.

*Note: If you do not specify worker nodes, all master nodes also become worker nodes.*
