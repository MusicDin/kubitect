<h1 align=center>Kubitect</h1>

<p align=center><b>Kubitect provides a simple way to set up a highly available Kubernetes cluster across multiple hosts.</b></p>

---

## Quick start

### Step 1 - Download and install `kubitect` tool

**TODO**

```
wget -O kubitect 

sudo install kubitect /usr/local/bin
```

### Step 2 - Create the cluster

Run the following command to create the default cluster.
Cluster will be created in `~/.kubitect/clusters/default/` directory.

```
kubitect apply
```

> :scroll: **Note**:
Using a `--cluster` option, you can provide custom cluster name.
This way multiple clusters can be created.

### Step 3 - Export kubeconfig

After successful installation of the Kubernetes cluster, Kubeconfig will be created within cluster's directory.
To export the Kubeconfig into custom file run the following command.

```
kubitect export kubeconfig > kubeconfig.yaml
```

### Step 4 - Test the cluster

Test if the cluster works by displaying all cluster nodes.

```
kubectl get nodes --kubeconfig kubeconfig.yaml
```

## Documentation
+ [Getting started](docs/getting-started.md)
+ [Requirements](docs/requirements.md)
+ [Configuration](docs/configuration.md)
+ [Load balancing](docs/load-balancer.md)
+ [Troubleshooting](docs/troubleshooting.md)
+ Examples: 
  - [Cluster over bridged network](docs/examples/bridged-network.md)
  - [Single node cluster deployment](docs/examples/single-node-cluster.md)

## Having issues?

In case you have found a bug, or some unexpected behaviour please [open an issue](https://github.com/MusicDin/terraform-kvm-kubespray/issues/new).

If you need anything else, you can contact me on [:email: din.music@din-cloud.com](mailto:din.music@din-cloud.com).

## License

[Apache License 2.0](./LICENSE)
