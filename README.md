# terraform-kvm-kubespray

> Set up highly available (HA) Kubernetes cluster using KVM, Terraform and Kubespray.

## Project goal

The goal of this project is to provide a simple way to set up a highly available Kubernetes cluster and to enable the setup of a cluster on multiple physical machines (hosts).

## Quick start

### Step 1 - Download and install `tkk` tool

**TODO**

```
wget -O tkk https://raw.githubusercontent.com/MusicDin/terraform-kvm-kubespray/<BRANCH>/scripts/tkk.sh

sudo install tkk /usr/local/bin
```

### Step 2 - Create the cluster

Run the following command to create the default cluster.
Cluster will be created in `~/.tkk/clusters/default/` directory.

```
tkk apply
```

> :scroll: **Note**:
Using a `--cluster` option, you can provide custom cluster name.
This way multiple clusters can be created.

### Step 3 - Export kubeconfig

After successful installation of the Kubernetes cluster, Kubeconfig will be created within cluster's directory.
To export the Kubeconfig into custom file run the following command.

```
tkk export kubeconfig > kubeconfig.yaml
```

### Step 4 - Test the cluster

Test if the cluster works by displaying all cluster nodes.

```
kubectl get nodes --kubeconfig kubeconfig.yaml
```

### What's next?

By following the quick start you have created a *default* cluster.
If custom configuration file is not passed to the apply command then  [default configuration file](/examples/localhost_1-worker_1-master.yaml) is used.

This configuration defines a simple local cluster that consists of 1 master and 1 worker node which is mostly used to test if *tkk* works on your setup.

Now it's time to prepare your own cluster.


## Documentation
+ [Getting started](docs/getting-started.md)
+ [Requirements](docs/requirements.md)
+ [Configuration](docs/configuration.md)
+ [Load balancing](docs/load-balancer.md)
+ [Troubleshooting](docs/troubleshooting.md)
+ Examples: 
  - [Cluster over bridged network](docs/examples/bridged-network.md)
  - [Single node cluster deployment](docs/examples/single-node-cluster.md)


## Related projects

If you are interested in deploying a Kubernetes cluster on *vSphere* instead of *KVM* check out [this project](https://github.com/sguyennet/terraform-vsphere-kubespray).

## Having issues?

In case you have found a bug, or some unexpected behaviour please [open an issue](https://github.com/MusicDin/terraform-kvm-kubespray/issues/new).

If you need anything else, you can contact me on [:email: din.music@din-cloud.com](mailto:din.music@din-cloud.com).

## License

[Apache License 2.0](./LICENSE)
