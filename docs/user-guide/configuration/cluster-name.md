[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<div markdown="1" class="text-center">
# Cluster metadata
</div>

<div markdown="1" class="text-justify">

## Configuration

### Cluster name

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline:{ .icon-required } Required

The cluster name must be defined in the Kubitect configuration, as it acts as a **prefix for all cluster resources**.

```yaml
cluster:
  name: my-cluster
```

For instance, each virtual machine name is generated as `<cluster.name>-<node.type>-<node.instance.id>`. 
Therefore, the name of the virtual machine for the worker node with ID `1` would be `my-cluster-master-1`.

!!! Note
    Cluster name cannot contain prefix `local`, as it is reserved for local clusters (created with `--local` flag).

</div>
