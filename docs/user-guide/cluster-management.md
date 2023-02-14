Project currently supports the following actions that can be executed on the running Kubernetes cluster:

+ scaling the cluster
    - adding worker nodes,
    - removing worker nodes,
+ upgrading the cluster,
+ destroying the cluster.

!!! note "Note"
    Each action supports the `--cluster <cluster_name>` option, which allows you to execute the action on a specific cluster. 
    By default, all actions are executed on the `default` cluster, which corresponds to using the `--cluster default` option.

### Export cluster configuration file

Each action requires the cluster configuration file to be modified.
Cluster configuration file can be exported using `export` command of the `kubitect` tool.

```sh
kubitect export config > cluster.yaml
```

## Scale the cluster

### Add worker nodes to the cluster

In the configuration file add new worker nodes to `cluster.nodes.worker.instances` list.
```yaml title="cluster.yaml"
cluster:
  ...
  nodes:
    ...
    worker:
      instances:
        - id: 1
        - id: 2 # New worker node
        - id: 3 # New worker node
```

Apply the modified configuration using `kubitect` tool to add new worker nodes:
```sh
kubitect apply --config cluster.yaml --action scale
```


### Remove worker nodes from the cluster

In the configuration file remove worker nodes from `cluster.nodes.worker.instances` list.
```yaml title="cluster.yaml"
cluster:
  ...
  nodes:
    ...
    worker:
      instances:
        - id: 1
        #- id: 2
        #- id: 3
```

Apply the modified configuration using `kubitect` tool to remove worker nodes:
```sh
kubitect apply --config cluster.yaml --action scale
```


## Upgrade the cluster

!!! warning "Important"

    *Do not skip releases when upgrading--upgrade by one tag at a time.*
    > For more information read [Kubespray upgrades](https://github.com/kubernetes-sigs/kubespray/blob/master/docs/upgrades.md).

In the cluster configuration file set the following variables:
  + `kubernetes.version` and
  + `kubernetes.kubespray.version`.


!!! note "Note"
    Before upgrading the cluster, make sure that [Kubespray](https://github.com/kubernetes-sigs/kubespray#supported-components) supports a specific Kubernetes version.

Example:
```yaml title="cluster.yaml"
kubernetes:
  version: "v1.22.5" # Old value: "v1.21.6"
  ...
  kubespray:
    version: "v2.18.0" # Old value: "v2.17.1"
    ...
```

Apply the modified configuration using `kubitect` tool:
```sh
kubitect apply --config cluster.yaml --action upgrade
```


## Destroy the cluster

To destroy the cluster, simply run:
```sh
kubitect destroy
```