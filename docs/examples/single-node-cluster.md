# Single node cluster

If you want to initialize a cluster with only one node,
specify single master node in cluster configuration file:

```yaml
# single-node.yaml
cluster:
  ...
  nodes:
    master:
      instances:
      - id: 1
        ip: "10.10.64.5"         # If omitted DHCP lease is requested
        mac: "52:54:00:00:00:40" # If omitted MAC address is generated
```

Do not forget to remove (or comment out) the worker and load balancer nodes.

Apply the cluster:
```sh
kubitect apply --config single-node.yaml
```

Your master node now also becomes a worker node.

> :scroll: **Note:** If you do not specify worker nodes, all master nodes also become worker nodes.
