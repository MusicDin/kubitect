<div markdown="1" class="text-center">
# Requirements
</div>

<div markdown="1" class="text-justify">

On the local host (*where Kubitect command-line tool is installed*), the following requirements must be met:

:material-record-circle-outline: [Git](https://git-scm.com/)

:material-record-circle-outline: [Python](https://www.python.org/) >= 3.8

:material-record-circle-outline: Python [virtualenv](https://virtualenv.pypa.io/en/latest/index.html)

:material-record-circle-outline: Password-less SSH key for each **remote** host

<br/>

On hosts where a Kubernetes cluster will be deployed using Kubitect, the following requirements must be met:

:material-record-circle-outline: A [libvirt](https://libvirt.org/) virtualization API

:material-record-circle-outline: A running hypervisor that is supported by libvirt (e.g. KVM)

??? question "How to install KVM? <i class="click-tip"></i>"

    To install the [KVM](https://www.linux-kvm.org) (Kernel-based Virtual Machine) hypervisor and libvirt, use *apt* or *yum* to install the following packages:

    + `qemu-kvm`
    + `libvirt-clients`
    + `libvirt-daemon`
    + `libvirt-daemon-system`

    After the installation, add your user to the `kvm` group in order to access the kvm device:

    ```sh
    sudo usermod -aG kvm $USER
    ```

</div>
