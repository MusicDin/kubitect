<div markdown="1" class="text-center">
# Scaling the cluster
</div>

<div markdown="1" class="text-justify">

Any cluster created with Kubitect can be subsequently scaled.
To do so, simply change the configuration and reapply it using the `scale` action.

!!! info "Info"

    Currently, only worker nodes and load balancers can be scaled.

## Export the cluster configuration

Exporting the current cluster configuration is optional, but strongly recommended to ensure that changes are made to the latest version of the configuration.
The cluster configuration file can be exported using the `export` command.

```sh
kubitect export config --cluster my-cluster > cluster.yaml
```

## Scale the cluster

In the configuration file, add new or remove existing nodes.

```yaml title="cluster.yaml"
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
kubitect apply --config cluster.yaml --action scale
```

As a result, the worker node with ID 2 is removed and the worker nodes with IDs 3 and 4 are added to the cluster.

</div>
