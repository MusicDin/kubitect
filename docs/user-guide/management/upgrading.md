<div markdown="1" class="text-center">
# Upgrading the cluster
</div>

A running Kubernetes cluster can be upgraded to a higher version by increasing the Kubernetes version in the cluster's configuration file and reapplying it using the `upgrade` action.


## Export the cluster configuration

Exporting the current cluster configuration is optional, but strongly recommended to ensure that changes are made to the latest version of the configuration.
The cluster configuration file can be exported using the `export` command.

```sh
kubitect export config --cluster my-cluster > my-cluster.yaml
```


## Upgrade the cluster

!!! warning "Important"
    Do not skip Kubitect's minor releases when upgrading the cluster.

In the cluster configuration file, change the Kubernetes version.

Example:
```yaml title="cluster.yaml"
kubernetes:
  version: v1.22.5 # Old value: v1.21.6
  ...
```

Apply the modified configuration using `upgrade` action.
```sh
kubitect apply --config cluster.yaml --action upgrade
```

The cluster is upgraded using the *in-place* strategy, i.e., the nodes are upgraded one after the other, making each node unavailable for the duration of its upgrade.
