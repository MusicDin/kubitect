<div markdown="1" class="text-center">
# Destroying the cluster
</div>

<div markdown="1" class="text-justify">

The destruction of the cluster consists of 2 parts.
The first part destroys all active cluster components, while the second part deletes all configuration files for a given cluster.

## Destroy the cluster

To destroy a specific cluster, simply run the destroy command, specifying the name of the cluster to be destroyed.

```sh
kubitect destroy --cluster my-cluster
```

This command destroys all active cluster components, including virtual machines, virtual networks, and resource pools, while preserving important configuration files, such as the Kubitect configuration.
Therefore, the cluster configuration file can be retrieved even after the cluster is destroyed.

## Purge the cluster

A purge command deletes all files associated with the cluster, including initialized virtual environments and configuration files.
It can be executed only if the cluster is already destroyed.

```sh
kubitect purge --cluster my-cluster
```

</div>
