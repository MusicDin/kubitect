# Requirements

## Local machine

The local machine is the machine on which the project is cloned. The following requirements must be met on the local machine:

+ [Git](https://git-scm.com/)
+ [Terraform](https://www.terraform.io/) >= 1.1.4
+ [Python](https://www.python.org/) >= 3.0
  - Project also uses [venv](https://docs.python.org/3/library/venv.html)
  
## Hosts

Hosts are physical servers running virtual machines that are part of the Kubernetes cluster. The local machine can also be a host.
Each host requires:

+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
  - Using *yum* or *apt* install following packages:
    + `qemu`
    + `qemu-kvm`
    + `libvirt-clients`
    + `libvirt-daemon`
    + `libvirt-daemon-system`
  - User needs to be int `kvm` and `libvirt` groups.
+ Password-less SSH keys (*Only if hosts are remote machines*).
