# Requirements

## Local machine

On the machine where the command-line tool is installed, the following requirements must be met:

+ [Git](https://git-scm.com/)
+ [Python](https://www.python.org/) >= 3.0
  - Python [virtualenv](https://docs.python.org/3/library/venv.html)

## Hosts

A host is a physical server that can be either a local or remote machine.
Each host must have **a hypervisor** installed along with the **[libvirt](https://libvirt.org/) virtualization API**.

If the host is a remote machine, pasword-less SSH keys are required to sucessfully connect to the remote hypervisor.

### Example - Install KVM

For example, to install the [KVM](https://www.linux-kvm.org) (Kernel Virtual Machine) hypervisor and libvirt, use *yum* or *apt* to install the following packages:
  + `qemu`
  + `qemu-kvm`
  + `libvirt-clients`
  + `libvirt-daemon`
  + `libvirt-daemon-system`

After installation, also add user to the `kvm` and `libvirt` groups.
