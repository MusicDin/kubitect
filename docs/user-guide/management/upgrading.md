A running Kubernetes cluster can be upgraded to a higher version by increasing the Kubernetes version in the cluster's configuration file and reapplying it using the `upgrade` action.


## Export the cluster configuration

Exporting the current cluster configuration is optional, but strongly recommended to ensure that changes are made to the latest version of the configuration.
The cluster configuration file can be exported using the `export` command.

```sh
kubitect export config --cluster my-cluster > my-cluster.yaml
```


## Upgrade the cluster

!!! warning "Important"

    Do not skip releases when upgrading--upgrade by one tag at a time.

    > For more information read [Kubespray upgrades](https://github.com/kubernetes-sigs/kubespray/blob/master/docs/upgrades.md).

Before upgrading the cluster, ensure that [Kubespray](https://github.com/kubernetes-sigs/kubespray#supported-components) supports a specific Kubernetes version.

In the cluster configuration file, set the desired Kubernetes version and adjust the Kubespray version if necessary.

Example:
```yaml title="cluster.yaml"
kubernetes:
  version: v1.22.5 # Old value: v1.21.6
  ...
  kubespray:
    version: v2.18.0 # Old value: v2.17.1
    ...
```

Apply the modified configuration using `upgrade` action.
```sh
kubitect apply --config cluster.yaml --action upgrade
```

The cluster is upgraded using the *in-place* strategy, i.e., the nodes are upgraded one after the other, making each node unavailable for the duration of its upgrade.
