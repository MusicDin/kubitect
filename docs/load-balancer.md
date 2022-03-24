# Internal load balancing (iLB)

Multiple master nodes ensure that services remain available if one or even more 
master nodes fail. Cluster has to be set up with an odd number of master nodes so 
that the quorum (the majority of master nodes) can be maintained if one or more 
masters fail. In the high-availability (HA) scenario, Kubernetes maintains a copy 
of the `etcd` databases on each master node, but holds elections for the `kube-controller` 
and `kube-scheduler` managers to avoid conflicts. This allows worker nodes to
communicate with any master node through a single endpoint provided by load balancers.


## Configure HAProxy load balancers

Specify load balancer instances in the configuration file (default is [cluster.yaml](/cluster.yaml)).
```yaml
cluster:
  ...
  nodes:
    ...
    loadBalancer:
      vip: 10.10.64.200 # Floating IP that should not be taken by any other device
      instances:
        - id: 1
        - id: 40
```
> :scroll: **Note:** Load balancers `id` must be a number between 0 and 200, because their fail-over priority is calculated from the `id`.

## Cluster without internal load balancers

If you have only one master node, the internal load balancers are redundant.
In this case, remove (or comment out) all load balancer nodes in the [cluster.yaml](/cluster.yaml) file:
```yaml
cluster:
  ...
  nodes:
    ...
    loadBalancer:
      ...
      instances:
        #- id: 1
        #- id: 40
```

> :scroll: **Note:** If multiple master nodes are specified, the IP of the first one is used as the cluster IP.