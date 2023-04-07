[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

<div markdown="1" class="text-center">
# Cluster node template
</div>

<div markdown="1" class="text-justify">

The node template section of the cluster configuration **defines the properties of all nodes** in the cluster.
This includes the properties of the operating system (OS), DNS, and the virtual machine user.


## Configuration

### Virtual machine user

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `k8s`

The user property defines the name of the user created on each virtual machine. 
This user is used to access the virtual machines during cluster configuration. 
If you omit the user property, a user named `k8s` is created on all virtual machines. 
You can also use this user later to access each virtual machine via SSH.

```yaml
cluster:
  nodeTemplate:
    user: kubitect
```

### Operating system (OS)

#### OS distribution

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]
&ensp;
:octicons-file-symlink-file-24: Default: `ubuntu`

The operating system for virtual machines can be specified in the node template. 
Currently, you can configure either Ubuntu or Debian. 
By default, the Ubuntu distribution is installed on all virtual machines.
To use Debian instead, set the `os.distro` property to Debian.

```yaml
cluster:
  nodeTemplate:
    os:
      distro: debian # (1)!
```

1. By default, `ubuntu` is used.

The available operating system distribution presets are:

+ `ubuntu` - Latest Ubuntu 22.04 release. (default)
+ `ubuntu22` - Ubuntu 22.04 release as of *2023-03-02*.
+ `ubuntu20` - Ubuntu 20.04 release as of *2023-02-09*.
+ `debian` - Latest Debian 11 release.
+ `debian11` - Debian 11 release as of *2023-01-24*.

Note that Ubuntu images are downloaded from the [Ubuntu cloud image repository](https://cloud-images.ubuntu.com/) and Debian images are downloaded from the [Debian cloud image repository](https://cloud.debian.org/images/cloud/).

#### OS source

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

If the presets do not meet your needs, you can use a custom Ubuntu or Debian image by specifying the image source. 
The source of an image can be either a local path on your system or a URL pointing to the image download.

```yaml
cluster:
  nodeTemplate:
    os:
      distro: ubuntu
      source: https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img
```

#### Network interface

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

Generally, this setting does not have to be set, as Kubitect will correctly evaluate the network interface name to be used on each virtual machine.

If you want to instruct Kubitect to use a specific network interface on the virtual machine, you can set its name using the `os.networkInterface` property.

```yaml
cluster:
  nodeTemplate:
    os:
      networkInterface: ens3
```

### Custom DNS list

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

The configuration of Domain Name Servers (DNS) in the node template allows for customizing the DNS resolution of all virtual machines in the cluster. 
By default, the DNS list contains only the network gateway.

To add custom DNS servers, specify them using the `dns` property in the node template.

```yaml
cluster:
  nodeTemplate:
    dns: # (1)!
      - 1.1.1.1
      - 1.0.0.1
```

1. IP addresses `1.1.1.1` and `1.0.0.1` represent CloudFlare's primary and secondary public DNS resolvers, respectively.

### CPU mode

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]
&ensp;
:octicons-file-symlink-file-24: Default: `custom`

The `cpuMode` property in the node template can be used to configure a guest CPU to closely resemble the host CPU.

```yaml
cluster:
  nodeTemplate:
    cpuMode: host-passthrough
```

Currently, there are several CPU modes available:

- `custom` (default)
- `host-model`
- `host-passthrough`
- `maximum`

In short, the `host-model` mode uses the same CPU model as the host, while the `host-passthrough` mode provides full CPU feature set to the guest virtual machine, but may impact its live migration. 
The `maximum` mode selects the CPU with the most available features.
For a more detailed explanation of the available CPU modes and their usage, please refer to the [libvirt documentation](https://libvirt.org/formatdomain.html#cpu-model-and-topology).

!!! tip "Tip"

    The `host-model` and `host-passthrough` modes makes sense only when a virtual machine can run directly on the host CPUs (e.g. virtual machines of type *kvm*).
    The actual host CPU is irrelevant for virtual machines with emulated virtual CPUs (e.g. virtul machines of type *qemu*).

### SSH options

#### Custom SSH certificate

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

Kubitect automatically generates SSH certificates before deploying the cluster to ensure secure communication between nodes.
The generated certificates can be found in the `config/.ssh/` directory inside the cluster directory.

If you prefer to use a custom SSH certificate, you can specify the local path to the private key. 
Note that the public key must also be present in the same directory with the `.pub` suffix.


```yaml
cluster:
  nodeTemplate:
    ssh:
      privateKeyPath: "~/.ssh/id_rsa_test"
```

!!! warning "Important"

    SSH certificates must be **passwordless**, otherwise Kubespray will fail to configure the cluster.


#### Adding nodes to the known hosts

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

Kubitect allows you to add all created virtual machines to SSH known hosts and remove them once the cluster is destroyed.
To enable this behavior, set the `addToKnownHosts` property to true.

```yaml
cluster:
  nodeTemplate:
    ssh:
      addToKnownHosts: true
```

</div>
