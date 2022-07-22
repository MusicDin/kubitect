<h1 align="center">Creating the cluster</h1>

With Kubitect, clusters are created by applying the cluster configuration to the Kubitect CLI tool.
If no cluster configuration is not specified, the default configuration is applied, as described in the [Quick start](../../../getting-started/quick-start) guide.

## Applying default configuration

To create a cluster with the default configuration, simply run the following command.
```sh
kubitect apply
```

## Applying custom configuration

To create a cluster with the custom configuration file, run the apply command with the `--config` flag.
```sh
kubitect apply --flag <PathToClusterConfig>
```

## Specifying cluster directory name

The configuration files for each cluster created with Kubitect are generated under the path `~/.kubitect/clusters/<ClusterName>`.
If no cluster name is specified, cluster name `default` is used.

The name of the cluster can be specified with the flag `--cluster`.
```sh
kubitect apply --cluster <ClusterName>
```