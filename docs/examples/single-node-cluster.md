# Single node cluster

If you would like to initialize a cluster with only one node than specify only master node in [terraform.tfvars](../../terraform.tfvars):
```hcl
vm_master_macs_ips = {
  "mac_for_master_1" = "ip_for_master_1"
}
```

Your master node will also become a worker node.

*Note: If you don't specify any worker node all master nodes will also become worker nodes.*