# terraform-kvm-kubespray
Set up HA Kubernetes cluster using KVM, Terraform and Kubespray.

## Requirements
+ [Git](https://git-scm.com/) 
+ [Cloud-init](https://cloudinit.readthedocs.io/)
+ [Ansible](https://www.ansible.com/) >= v2.6
+ [Terraform](https://www.terraform.io/) **== v0.12.x**
+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
+ [Libvirt provider](https://github.com/dmacvicar/terraform-provider-libvirt) - Setup guide is provided in [docs](./docs/libvirt-provider-setup.md).
+ Internet connection on machine that will run VMs and on VMs

*Note: for Terraform **v0.13.x** see [this branch](https://github.com/MusicDin/terraform-kvm-kubespray/).*


## Getting Started

*If you run into any troubles during installation process, please check [troubleshooting](docs/troubleshooting.md) page first.*

### Libvirt provider

If you haven't yet, [install libvirt provider](docs/libvirt-provider-setup.md).

### SSH keys

Generate SSH keys, which will be used to access created VMs:
```bash
ssh-keygen
```

Follow the instructions to create SSH keys:
<pre>
Generating public/private rsa key pair.
Enter file in which to save the key (/home/<b>your_username</b>/.ssh/id_rsa): <b>[1]</b>
Enter passphrase (empty for no passphrase): <b>[2]</b>
Enter same passphrase again: <b>[2]</b>
...
</pre>

**[1]** You will be asked to enter file in which to save the key. Default is `/home/your_username/.ssh/id_rsa`.

**[2]** When asked to enter a password, press `ENTER` twice to skip setting a password. 
**DO NOT** enter it, otherwise Terraform will fail to initialize a cluster.

Finally, you have to enter a location of SSH private key in `vm_ssh_private_key` field in [terraform.tfvars](terraform.tfvars) file.
 

### Cluster setup

Clone project and move to main directory:
```
git clone https://github.com/MusicDin/terraform-kvm-kubespray.git --branch terraform-0.12

cd terraform-kvm-kubespray
```

Change variables in [terraform.tfvars](terraform.tfvars) file to fit your needs.
Variables are set to work out of the box. 
Only required variables that are not set are: 
+ `vm_image_source` URL or path on file system to OS image,
+ `vm_distro` a Linux distribution of OS image.

**IMPORTANT:** Review variables before initializing a cluster, as current configuration will create 8 VMs which are quite resource heavy!

Execute terraform script:
```bash
# Initializes terraform project
terraform init

# Shows what is about to be done
terraform plan

# Runs/creates project
terraform apply
```

*Note: Installation process can take up to 20 minutes based on a current configuration.*

### Test cluster

All configuration files will be generated in `config/` directory, 
and one of them will be `admin.conf` which is actually a `kubeconfig` file.
 
Test your cluster by displaying all cluster's nodes:
```
kubectl --kubeconfig=config/admin.conf get nodes
```

## Cluster management

### Add worker to the cluster

In [terraform.tfvars](./terraform.tfvars) file add:
  + *MAC* address for a new VM to `vm_worker_macs` and 
  + *IP* address for a new VM to `vm_worker_ips`.
  
*Note: MAC and IP addresses for certain VM have to be on same index in array.*

Execute terraform script to add a worker:
```
terraform apply -var 'action=add_worker'
```

### Remove worker from the cluster

In [terraform.tfvars](./terraform.tfvars) file delete:
  + *MAC* address of a current worker from `vm_worker_macs` and 
  + *IP* address of a current worker from `vm_worker_ips`.

Execute terraform script to remove a worker:
```
terraform apply -var 'action=remove_worker'
```

### Upgrade cluster

In [terraform.tfvars](./terraform.tfvars) file modify:
  + `k8s_kubespray_version` and
  + `k8s_version`.
  
*Note: Before upgrading make sure Kubespray supports provided Kubernetes version.*
  
Execute terraform script to upgrade a cluster:
```
terraform apply -var 'action=upgrade'
```

### Destroy cluster

To destroy the cluster, simply run:
```
terraform destroy
```

## More documentation
+ [Setup libvirt provider](docs/libvirt-provider-setup.md)
+ [Load balancing](docs/load-balancer.md)
+ [Troubleshooting](docs/troubleshooting.md)
+ Examples: 
    - [Load balancing to ingress controller](docs/examples/lb-and-ingress-controller.md)


## Related projects

If you are interested in installing kubernetes cluster on *vSphere* instead of *KVM* check [this project](https://github.com/sguyennet/terraform-vsphere-kubespray).

## Having issues?

In case you have found a bug, or some unexpected behaviour please open an issue.

If you need anything else, you can contact me on GitHub.

## License

[Apache License 2.0](./LICENSE)