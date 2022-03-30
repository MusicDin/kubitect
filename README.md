# terraform-kvm-kubespray
> Set up highly available (HA) Kubernetes cluster using KVM, Terraform and Kubespray.

## Project goal

The goal of this project is to provide an easy way to set up highly 
available Kubernetes cluster and to allow setting up a cluster 
on multiple physical machines (hosts) that are running KVM.


## Documentation
+ [Requirements](docs/requirements.md)
+ [Getting started](docs/getting-started.md)
+ [Configuration](docs/configuration.md)
+ [Load balancing](docs/load-balancer.md)
+ [Troubleshooting](docs/troubleshooting.md)
+ Tutorial:
  -[From zero 2 hero]() - Tutorial how to setup a HA cluster step by step.
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
