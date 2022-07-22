[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<h1 align="center">Cluster metadata</h1>

## Configuration

### Cluster name

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline:{ .icon-required } Required

The cluster name must be defined as part of the Kubitect configuration.
It will be used as a **prefix for all resources** created by Kubitect as part of this cluster.

```yaml
cluster:
  name: my-cluster
```

For example, the name of each virtual machine name is generated as `<cluster.name>-<node.type>-<node.instance.id>`.
For the master node with ID 1, this would result in `my-cluster-master-1`.

