# Cluster management

Project currently supports the following actions that can be executed on the running Kubernetes cluster:
+ adding worker nodes,
+ removing worker nodes,
+ upgrading the cluster,
+ destroying the cluster.


## Add worker nodes to the cluster

In the configuration file (default is [cluster.yaml](/cluster.yml)) add new worker nodes to `cluster.nodes.worker.instances` list.
```yml
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

Apply the modified configuration using tkk tool to add new worker nodes:
```shell
tkk apply --config cluster.yaml --action add-worker
```


## Remove worker nodes from the cluster

In the configuration file (default is [cluster.yaml](/cluster.yml)) remove worker nodes from `cluster.nodes.worker.instances` list.
```yaml
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

Apply the modified configuration using `tkk` command to remove worker nodes:
```shell
tkk apply --config cluster.yaml --action remove_worker
```


## Upgrade the cluster

> :exclamation: **IMPORTANT:** *Do not skip releases when upgrading--upgrade by one tag at a time.*
> For more information read [Kubespray upgrades](https://github.com/kubernetes-sigs/kubespray/blob/master/docs/upgrades.md).

> :bulb: **Tip:** Before upgrading the cluster, make sure [Kubespray](https://github.com/kubernetes-sigs/kubespray#supported-components) supports specific Kubernetes version.

In the configuration file (default is [cluster.yaml](/cluster.yml)) set the following variables:
  + `kubernetes.version` and
  + `kubernetes.kubespray.version`.

For example:
```yaml
kubernetes:
  version: "v1.22.5" # Old value: "v1.21.6"
  ...
  kubespray:
    version: "v2.18.0" # Old value: "v2.17.1"
    ...
```

Execute terraform script to upgrade the cluster:
```shell
tkk apply --config cluster.yaml --action upgrade
```


## Destroy the cluster

To destroy the cluster, simply run:
```shell
tkk destroy
```