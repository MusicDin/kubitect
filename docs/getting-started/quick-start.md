<div markdown="1" class="text-center">
# Quick start
</div>

<div markdown="1" class="text-justify">

<div class="text-center">
  <img
    class="mobile-w-100"
    src="/assets/images/topology-1m1w-arch.png" 
    alt="Architecture of the cluster with one master and one worker node"
    width="75%">
</div>

### Step 1 - Create the cluster

Run the following command to apply the default cluster configuration, which creates a cluster with **one master and one worker node**.
Generated cluster configuration files will be stored in `~/.kubitect/clusters/default/` directory.

```
kubitect apply
```

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

</div>
