# terraform-kvm-kubespray
Set up Kubernetes cluster using KVM, Terraform and Kubespray

## Requirements
+ [Git](https://git-scm.com/) 
+ [Cloud-init](https://cloudinit.readthedocs.io/)
+ [Ansible](https://www.ansible.com/) >= v2.6
+ [Terraform](https://www.terraform.io/) **>= v0.13.x**
+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
+ Internet connection on machine that will run VMs and on VMs

*Note: for Terraform v0.12.x see [this branch](https://github.com/MusicDin/terraform-kvm-kubespray/tree/terraform-0.12).*

## Usage

### Create cluster

Move to main directory:
```
cd terraform-kvm-kubespray
```

Change variables to fit your needs:
```
nano terraform.tfvars
```
*Note: Variables are set to work out of the box. Only required variable that is not set is* `vm_image_source`.

Execute terraform script:
```bash
# Initializes terraform project
terraform init

# Shows what is about to be done
terraform plan

# Runs/creates project
terraform apply
```

*Note: Installation proccess can take up to 20 minutes.*

### Add worker to cluster

In [terraform.tfvars](./terraform.tfvars) file add *MAC* and *IP* address for new VM to `vm_worker_macs_ips`. 
  
Execute terraform script to add worker:
```
terraform apply -var 'action=add_worker'
```

### Remove worker from cluster

In [terraform.tfvars](./terraform.tfvars) file remove *MAC* and *IP* address of VM that is going to be deleted from `vm_worker_macs_ips`.

Execute terraform script to remove worker:
```
terraform apply -var 'action=remove_worker'
```
### Upgrade cluster

In [terraform.tfvars](./terraform.tfvars) file modify:
  + `k8s_kubespray_version` and
  + `k8s_version`.
  
Execute terraform script to upgrade cluster:
```
terraform apply -var 'action=upgrade'
```

## Credits

This script is modified to work with *KVM* instead of *vSphere*.

Script that works on *vSphere* can be found [here](https://github.com/sguyennet/terraform-vsphere-kubespray).

## Info/Issues

In case you have found a bug or dysfunctionality please open an issue.

If you need anything else contact me on GitHub.
