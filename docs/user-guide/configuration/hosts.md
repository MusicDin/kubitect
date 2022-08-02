[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<h1 align="center">Hosts configuration</h1>

Defining Kubitect hosts is esential. 
**Hosts represent the target servers** where the cluster will be deployed.
Every valid configuration must contain at least one host, but there can be as many hosts as needed.
The host can be either a local or remote server. 

## Configuration

### Localhost

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

When cluster is deployed on the server where the Kubitect command line tool is installed,
a host whose connection type is set to local needs to be specified. 
Such host is also refered to as localhost.

```yaml
hosts:
  - name: localhost # (1)
    connection:
      type: local
``` 

1. Custom **unique** name of the host.

### Remote hosts

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

When cluster is deployed on the remote host, the IP address of the remote host along with the SSH credentails needs to be specified for the host.

```yaml
hosts:
  - name: my-remote-host
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.143 # (1)
      ssh:
        keyfile: "~/.ssh/id_rsa_server1" # (2)
```

1. IP address of the remote host.

2. Path to the **password-less** SSH key file required for establishing connection with the remote host. Default is `~/.ssh/id_rsa`.

#### Host's SSH port

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `22`

By default, port `22` is used for SSH.
If host is running SSH client on a different port, it is possible to change it for each host separately.

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

By default remote hosts are not verified in the known SSH hosts.
If for any reason host verification is desired, you can enable it for each host separately.

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

If the host is specified as the default, all instances that do not point to a specific host are deployed to the default host. 
If the default host is not specified, these instances are deployed on the first host in the list.

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

The main resource pool path defines the location on the host where main disks (volumes) are created for each node provisioned on that particular host.
Because the main resource pool contains volumes on which the node's operating system and all required packages are installed, it is recommended that the main resource pool is created on fast storage devices such as SSD disks.
By default, main disk pool path is set to `/var/lib/libvirt/images/`.

```yaml
hosts:
  - name: host1 # (1)
  - name: host2 
    mainResourcePoolPath: /mnt/ssd/kubitect/ # (2)
```

1. Because the main resource pool path for this host is not set, the default path (`/var/lib/libvirt/images/`) is used.

2. The main resource pool path is set for this host, so the node's main disks are created in this location.

### Data resource pools

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]

Data resource pools define **additional resource pools** (*besides the required main resource pool*).
They are useful in situations where multiple disks need to be attached to the node's virtual machines.
For example, main disks contain the OS image and should be created on fast storage devices, while data resource pools can be used to attach additional virtual disks that can be created on slower storage devices such as HDDs.

Multiple data resource pools can be defined and each pool must contain a unique name (on a specific host) and a path under which it is created.
The name of the data resource pool is used to associate virtual disks defined in the node configuration with the actual data resource pool.


```yaml
hosts:
  - name: host1
    dataResourcePools:
      - name: rook-pool
        path: /mnt/hdd/kubitect/pools/
  - name: host2 
    dataResourcePools:
      - name: rook-pool
        path: /var/lib/libvirt/images/
```


## Example usage

### Multiple hosts

With Kubitect the cluster can be deployed on multiple hosts.
All hosts need to be specified in the configuration file.

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