<h1 align="center">Full (detailed) example</h1>

```yaml
#
# In the 'kubitect' section, you can specify the target git project and version.
# This can be handy if you want to use a specific project version or if you
# want to point to your forked/cloned project.
#
# [!] Note that this is ignored if you use the --local option with the
#     actions of the CLI tool, since in this case you should be in the 
#     Git repository.
#
kubitect:
  url: "https://github.com/MusicDin/terraform-kvm-kubespray" # (1)
  version: "2.0.0"

#
# The "hosts" section contains data about the physical servers on which 
# the Kubernetes cluster will be installed.
#
# For each host, a name and connection type must be specified. Only one
# host can have the connection type set to 'local' or 'localhost'.
#
# If the host is a remote machine, SSH key file must be specified.
# [!] Note that the connection to the remote hosts supports only 
#     passwordless login (using only SSH keyfile).  
#
# The host can be marked also as default, which means that if an instance
# (in "cluster.nodes" section) does not have a particular host specified,
# it will be installed on a default host. If none of the hosts are marked 
# as default, the first one in the list is used as the default host.
# 
hosts:
  - name: localhost # (3)
    default: true # (4)
    connection:
      type: local # (5)
  - name: remote-server-1
    connection:
      type: remote
      user: myuser # (6)
      ip: 10.10.40.143 # (7)
      ssh:
        port: 1234  # (8)
        verify: false # (9)
        keyfile: "~/.ssh/id_rsa_server1" # (10)
  - name: remote-server-2
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.144
      ssh:
        keyfile: "~/.ssh/id_rsa_server2"
    mainResourcePoolPath: "/var/lib/libvirt/pools/" # (11)
    dataResourcePools: # (12)
      - name: data-pool
        path: "/mnt/data/pool"
      - name: backup-pool
        path: "/mnt/backup/pool"

#
# The "cluster" section of configuration contains general data about the cluster,
# nodes that are part of the cluster and cluster's network.
# 
cluster:
  name: "my-k8s-cluster" # (13)
  network:
    mode: bridge # (14)
    cidr: "10.10.64.0/24" # (15)
    gateway: 10.10.64.1 # (16)
    bridge: br0 # (17)
    dns: # (18)
      - 1.1.1.1
      - 1.0.0.1
  nodeTemplate:
    networkInterface: "ens3" # (19)
    user: "k8s"
    ssh:
      privateKeyPath: "~/.ssh/id_rsa_test"
      addToKnownHosts: true
    image:
      distro: "ubuntu"
      source: "https://cloud-images.ubuntu.com/releases/focal/release-20220111/ubuntu-20.04-server-cloudimg-amd64.img"
    updateOnBoot: true
  nodes:
    loadBalancer:
      vip: "10.10.64.200" # (20)
      default: # (21)
        ram: 4 # GiB
        cpu: 1 # vCPU
        mainDiskSize: 16 # GiB
      instances:
        - id: 1
          ip: 10.10.64.5 # (22)
          mac: "52:54:00:00:00:40" # (23)
          ram: 8 # (24)
          cpu: 8 # (25)
          host: remote-server-1 # (26)
        - id: 2
          ip: 10.10.64.6
          mac: "52:54:00:00:00:41"
          host: server2
        - id: 3
          ip: 10.10.64.7
          mac: "52:54:00:00:00:42"
          # If server is not specifed, VM will be installed on the default server.
          # If default server is not specified, VM will be installed on the first
          # server in the list.
    master:
      default:
        ram: 8
        cpu: 2
        mainDiskSize: 256
      instances:
          # IMPORTANT: There should be odd number of master nodes.
        - id: 1 # Node with generated MAC address, IP retrieved as an DHCP lease and default RAM and CPU.
          host: remote-server-1
        - id: 2
          host: remote-server-2
        - id: 3
          server: localhost
    worker:
      default:
        ram: 16
        cpu: 4
        label: node # (27)
        # Default dataDisks are NOT YET supported
        # dataDisks: # (29)
        #  - name: rook-disk # (30)
        #    pool: data-pool # (31)
        #    size: 128       # (32)
        #  - name: backup-disk
        #    pool: data-pool
        #    size: 512
      instances:
        - id: 1
          ip: 10.10.64.101
          cpu: 8
          ram: 64
          host: remote-server-1
        - id: 2
          ip: 10.10.64.102
          dataDisks: # (33)
            - name: rook-disk
              pool: data-pool
              size: 128
            - name: test-disk
              pool: data-pool
              size: 128
        - id: 3
          ip: 10.10.64.103
          ram: 64
        - id: 4
          host: remote-server-2
        - id: 5

#
# The "kubernetes" section specifies what version of Kubernetes and Kubespray
# should be used, which network plugin and dns server should be installed and
# whether or not to install a Kubespray addons.
#
kubernetes:
  version: "v1.21.6"
  networkPlugin: calico
  dnsMode: coredns
  kubespray:
    url: "https://github.com/kubernetes-sigs/kubespray.git"
    version: "v2.17.1"
  other:
    copyKubeconfig: false
```

1.  This allows you to set a custom URL that targets clone/fork of Kubitect project.

2.  Kubitect version.

3.  Custom host name. 
    It is used to link instances to the specific host.

4.  Makes the host a default host. 
    This means that if no host is specified for the node instance, the instance will be linked to the default host.

5.  Connection type can be either `local` or `remote`. 

    If it is set to *remote*, at least the following fields must be set:

    + `user`
    + `ip`
    + `ssh.keyfile`

6.  Remote host user that is used to connect to the remote hypervisor. 
    This user must be added in the `libvirt` group.

7.  IP address of the remote host.

8.  Overrides default SSH port (22).

9.  If set to false, host verification is skipped.

10. Path to the **passwordless** SSH key used to connect to the remote host.

11. The path to the main resource pool defines where the virtual machine disk images are stored. These disks contain the virtual machine operating system, and therefore it is recommended to install them on SSD disks.

12. List of other storage pools where virtual disks can be created.

13. Cluster name used as a prefix for the various components.

14. Network mode. Possible values are
    
    + `bridge` mode uses **predefined** bridge interface. This mode is mandatory for deployments across multiple hosts.
    + `nat` mode creates virtual network with IP range defined in `network.cidr`
    + `route`

15. Network CIDR represents the network IP together with the network mask. 
    In `nat` mode, CIDR is used for the new network.
    In `bridge` mode, CIDR represents the current local area network (LAN).

16. The network gateway IP address.
    If omitted the first client IP from network CIDR is used as a gateway.

17. Bridge represents the bridge interface on the hosts.
    This field is mandatory if the network mode is set to `bridge`.
    If the network mode is set to `nat`, this field can be omitted.

18. Set custom DNS for nodes. 
    If omitted, gateway is also used as DNS.

19. Specify the network interface used by the virtual machine. In general, this option can be omitted. 

    If you omit it, `ens3` is used for Ubuntu images and `eth0` for all other distributions.

20. Virtual (floating) IP shared between load balancers. 

21. Default values apply for all virtual machines (VMs) of the same type.

22. Static IP address of the virtual machine. 
    If omitted DHCP lease is requested.

23. Static MAC address. 
    If omitted MAC address is generated.

24. Overrides default RAM value for this node.

25. Overrides default CPU value for this node.

26. Name of the host where instance should be created.
    If omitted the default host is used.

27. Worker nodes label.

28. Overrides default data disks for this node.

29. Default data disks (attached to each worker node).

30. Unique data disk name.

31. Reference to the data resource pool that must exist on the same host as this node.

32. Size of the data disk in GiB. 
    Note that each node receives a data disk of a specific size.

33. Overrides default data disks for this node.