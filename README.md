# terraform-kvm-kubespray
Set up Kubernetes cluster using KVM, Terraform and Kubespray

## Requirements
+ [Git](https://git-scm.com/) 
+ [Ansible](https://www.ansible.com/) >= v2.6
+ [Terraform](https://www.terraform.io/) == v0.12.x
+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
+ Internet connection on machine that will run VMs and on VMs

*Note: for Terraform **v0.13.x** see [this branch](https://github.com/MusicDin/terraform-kvm-kubespray/).*

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

In [terraform.tfvars](./terraform.tfvars) file add:
  + *MAC* address for new VM to `vm_worker_macs` and 
  + *IP* address for new VM to `vm_worker_ips`.
  
*Note: MAC and IP addresses for certain VM have to be on same index in array.*

Execute terraform script to add worker:
```
terraform apply -var 'action=add_worker'
```

### Remove worker from cluster

+ TBD - (Not tested yet)

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

In case you found a bug or dysfunctionality please open an issue.

If you need anything else contact me on GitHub.
