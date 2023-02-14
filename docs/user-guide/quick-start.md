### Step 1 - Create the cluster

Run the following command to create the default cluster.
Cluster will be created in `~/.kubitect/clusters/default/` directory.

```
kubitect apply
```

!!! note "Note"

    Using a `--cluster` option, you can provide custom cluster name.
    This way multiple clusters can be created.

### Step 2 - Export kubeconfig

After successful installation of the Kubernetes cluster, Kubeconfig will be created within cluster's directory.
To export the Kubeconfig into custom file run the following command.

```
kubitect export kubeconfig > kubeconfig.yaml
```

### Step 3 - Test the cluster

Test if the cluster works by displaying all cluster nodes.

```
kubectl get nodes --kubeconfig kubeconfig.yaml
```