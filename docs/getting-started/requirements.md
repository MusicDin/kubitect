<div markdown="1" class="text-center">
# Requirements
</div>

## Local machine

<div markdown="1" class="text-justify">

On the machine where the command-line tool (kubitect) is installed, the following requirements must be met:

+ [Git](https://git-scm.com/)
+ [Python](https://www.python.org/) >= 3.0
  - Python [virtualenv](https://virtualenv.pypa.io/en/latest/index.html)

## Hosts

A host is a physical server that can be either a local or remote machine.
Each host must have:

+ installed [libvirt](https://libvirt.org/) virtualization API and
+ installed hypervisor that is supported by libvirt

If the host is a remote server:

+ password-less SSH key to connect to the remote server

!!! quote ""

    ### Example - Install KVM

    For example, to install the [KVM](https://www.linux-kvm.org) (Kernel-based Virtual Machine) hypervisor and libvirt, use *yum* or *apt* to install the following packages:

    + `qemu-kvm`
    + `libvirt-clients`
    + `libvirt-daemon`
    + `libvirt-daemon-system`

    After installation, add user to the `kvm` and `libvirt` groups.

    ```sh
    sudo adduser $USER kvm
    sudo adduser $USER libvirt
    ```

</div>
