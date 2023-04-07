[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<div markdown="1" class="text-center">
# Hosts configuration
</div>

<div markdown="1" class="text-justify">

Defining hosts is an essential step when deploying a Kubernetes cluster with Kubitect. 
**Hosts represent the target servers** where the cluster will be deployed.

Every valid configuration must contain at least one host, which can be either local or remote.
However, you can add as many hosts as needed to support your cluster deployment. 


## Configuration

### Localhost

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

To configure a local host, you simply need to specify a host with the connection type set to `local`.

```yaml
hosts:
  - name: localhost # (1)!
    connection:
      type: local
``` 

1. Custom **unique** name of the host.

### Remote hosts

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

To configure a remote host, you need to set the connection type to `remote` and provide the IP address of the remote host, along with its SSH credentials.

```yaml
hosts:
  - name: my-remote-host
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.143 # (1)!
      ssh:
        keyfile: "~/.ssh/id_rsa_server1" # (2)!
```

1. IP address of the remote host.

2. Path to the **password-less** SSH key file required for establishing connection with the remote host. Default is `~/.ssh/id_rsa`.

#### Host's SSH port

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `22`

By default, SSH uses port `22`. If a host is running an SSH client on a different port, you can change the port for each host separately.

```yaml
hosts:
  - name: remote-host
    connection:
      type: remote
      ssh:
        port: 1234
```

#### Host verification (known SSH hosts)

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

By default, remote hosts are not verified in the known SSH hosts.
If you want to verify hosts, you can enable host verification for each host separately.

```yaml
hosts:
  - name: remote-host
    connection:
      type: remote
      ssh:
        verify: true
```

### Default host

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

If a host is specified as the default, all instances that do not point to a specific host are deployed to that default host. 
If no default host is specified, these instances are deployed on the first host in the list.

```yaml
hosts:
  - name: localhost
    connection:
      type: local
  - name: default-host
    default: true
    ...
```

### Main resource pool

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `/var/lib/libvirt/images/`

The main resource pool path specifies the location on the host where main virtual disks (volumes) are created for each node provisioned on that particular host. 
Because the main resource pool contains volumes on which the node's operating system and all required packages are installed, it's recommended that the main resource pool is created on fast storage devices, such as SSD disks.

```yaml
hosts:
  - name: host1 # (1)!
  - name: host2 
    mainResourcePoolPath: /mnt/ssd/kubitect/ # (2)!
```

1. Because the main resource pool path for this host is not set, the default path (`/var/lib/libvirt/images/`) is used.

2. The main resource pool path is set for this host, so the node's main disks are created in this location.

### Data resource pools

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

Data resource pools allow you to define additional resource pools, besides the required main resource pool. These pools can be used to attach additional virtual disks that can be used for various storage solutions, such as Rook or MinIO.

Multiple data resource pools can be defined on each host, and each pool must have a unique name on that host. The name of the data resource pool is used to associate the virtual disks defined in the node configuration with the actual data resource pool.

By default, the path of the data resources is set to `/var/lib/libvirt/images`, but it can be easily configured using the `path` property.

```yaml
hosts:
  - name: host1
    dataResourcePools:
      - name: rook-pool
        path: /mnt/hdd/kubitect/pools/
      - name: data-pool # (1)!
```

1. If the path of the resource pool is not specified, it will be created under the path `/var/lib/libvirt/images/`.

## Example usage

### Multiple hosts

Kubitect allows you to deploy a cluster on multiple hosts, which need to be specified in the configuration file.

```yaml
hosts:
  - name: localhost
    connection:
      type: local
  - name: remote-host-1
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.143
      ssh:
        port: 123
        keyfile: "~/.ssh/id_rsa_server1"
  - name: remote-host-2
    default: true
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.145
      ssh:
        keyfile: "~/.ssh/id_rsa_server2"
  ...
```

</div>
