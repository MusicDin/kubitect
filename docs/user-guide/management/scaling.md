<h1 align="center">Scaling the cluster</h1>

Any cluster created with Kubitect can be subsequently scaled.
To do so, simply change the configuration and reapply it using the `scale` action.

## Export the cluster configuration

Exporting the current cluster configuration is optional, but strongly recommended to ensure that changes are made to the latest version of the configuration.
The cluster configuration file can be exported using the `export` command.

```sh
kubitect export config --cluster my-cluster > my-cluster.yaml
```

## Scale the cluster

!!! warning "Warning"

    Currently, only worker nodes can scale seamlessly. 
    The ultimate goal is to be able to scale every node type in the cluster with a single command.
    It is planned to address this issue in one of the following releases.

Kubitect supports the simultaneous addition and removal of multiple worker nodes.
In the configuration file, add new workers or remove/comment workers from the `cluster.nodes.worker.instances` list.

```yaml title="my-cluster.yaml"
cluster:
  ...
  nodes:
    ...
    worker:
      instances:
        - id: 1
        #- id: 2 # Worker node to be removed
        - id: 3 # New worker node
        - id: 4 # New worker node
```

Apply the modified configuration with action set to `scale`:
```sh
kubitect apply --config my-cluster.yaml --cluster my-cluster --action scale
```

As a result, the worker node with ID 2 is removed and the worker nodes with IDs 3 and 4 are added to the cluster.
