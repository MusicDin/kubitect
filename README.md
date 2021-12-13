# terraform-kvm-kubespray
Set up HA Kubernetes cluster using KVM, Terraform and Kubespray.

## Requirements
+ [Git](https://git-scm.com/)
+ [Cloud-init](https://cloudinit.readthedocs.io/)
+ [Ansible](https://www.ansible.com/) >= 2.9
+ [Terraform](https://www.terraform.io/) >= 1.0.0
+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
+ [Libvirt provider](https://github.com/dmacvicar/terraform-provider-libvirt) >= 0.6.12
+ Internet connection


## Getting Started

*If you encounter any issues during the installation, please refer to the [troubleshooting](docs/troubleshooting.md) page first.*

Clone the project and move to the main directory:
```
git clone https://github.com/MusicDin/terraform-kvm-kubespray.git

cd terraform-kvm-kubespray
```

### SSH keys

Generate SSH keys that will be used to access created VMs:
```bash
ssh-keygen
```

Follow the instructions to generate SSH keys:
<pre>
Generating public/private rsa key pair.
Enter file in which to save the key (/home/<b>your_username</b>/.ssh/id_rsa): <b>[1]</b>
Enter passphrase (empty for no passphrase): <b>[2]</b>
Enter same passphrase again: <b>[2]</b>
...
</pre>

**[1]** You will be asked to enter file in which to save the key. Default is `/home/your_username/.ssh/id_rsa`.

**[2]** When asked to enter a password, press <kbd>ENTER</kbd> twice to skip setting a password.
**DO NOT** enter it, otherwise Terraform will fail to initialize a cluster.

Finally, you have to enter a location of SSH private key in the `vm_ssh_private_key` field in [terraform.tfvars](terraform.tfvars) file.


### Cluster setup

Change variables in [terraform.tfvars](terraform.tfvars) file to fit your needs.
Variables are set to work out of the box.
Only unset required variable is:
+ `vm_image_source` URL or path on the file system to OS image

**IMPORTANT:**
Review variables before initializing a cluster.

*Note: Script also supports deployment of [single node cluster](docs/examples/single-node-cluster.md).*

Execute terraform script:
```bash
# Initializes terraform project
terraform init

# Shows what is about to be done
terraform plan

# Runs/creates project
terraform apply
```

*Note: The installation process can take up to 20 minutes depending on the configuration.*

### Test cluster

All configuration files will be generated in `config/` directory,
and one of them will be `admin.conf` which is actually a `kubeconfig` file.

Test if the cluster works by displaying all cluster nodes:
```
kubectl --kubeconfig=config/admin.conf get nodes
```

## Cluster management

### Adding worker nodes to the cluster

In [terraform.tfvars](./terraform.tfvars) file add new worker node(s) in `worker_nodes` list.

Execute terraform script to add a worker (workers):
```
terraform apply -var 'action=add_worker'
```

### Removing worker nodes from the cluster

In [terraform.tfvars](./terraform.tfvars) file remove worker node(s) from `worker_nodes` list.

Execute terraform script to remove a worker (workers):
```
terraform apply -var 'action=remove_worker'
```
### Upgrading the cluster

In [terraform.tfvars](./terraform.tfvars) file modify:
  + `k8s_kubespray_version` and
  + `k8s_version`.

*Note: Before upgrading make sure [Kubespray](https://github.com/kubernetes-sigs/kubespray#supported-components) supports provided Kubernetes version.*

Execute terraform script to upgrade a cluster:
```
terraform apply -var 'action=upgrade'
```

**IMPORTANT**:
*Do not skip releases when upgrading--upgrade by one tag at a time.*
For more information read [Kubespray upgrades](https://github.com/kubernetes-sigs/kubespray/blob/master/docs/upgrades.md).

### Destroying the cluster

To destroy the cluster, simply run:
```
terraform destroy
```

## More documentation
+ [Load balancing](docs/load-balancer.md)
+ [Troubleshooting](docs/troubleshooting.md)
+ Examples:
  - [Cluster over bridged network](docs/examples/bridged-network.md)
  - [Single node cluster deployment](docs/examples/single-node-cluster.md)
  - [MetalLB configuration examples](docs/examples/metallb.md)

## Related projects

If you are interested in deploying a Kubernetes cluster on *vSphere* instead of *KVM* check out [this project](https://github.com/sguyennet/terraform-vsphere-kubespray).

## Having issues?

In case you have found a bug, or some unexpected behaviour please open an issue.

If you need anything else, you can contact me on GitHub.

## License

[Apache License 2.0](./LICENSE)
