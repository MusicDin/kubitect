<div markdown="1" class="text-center">
# Full (detailed) example
</div>

<div markdown="1" class="text-justify">
This document contains an example of Kubitect configuration.
Example covers all (*or most*) of the Kubitect properties.
This example is meant for users that learn the fastest from an example configuration.
</div>

```yaml
#
# The 'hosts' section contains data about the physical servers on which the
# Kubernetes cluster will be installed.
#
# For each host, a name and connection type must be specified. Only one host can
# have the connection type set to 'local' or 'localhost'.
#
# If the host is a remote machine, the path to the SSH key file must be specified.
# Note that connections to remote hosts support only passwordless certificates.
#
# The host can also be marked as default, i.e. if no specific host is specified
# for an instance (in the cluster.nodes section), it will be installed on a
# default host. If none of the hosts are marked as default, the first host in the
# list is used as the default host.
#
hosts:
  - name: localhost # (3)!
    default: true # (4)!
    connection:
      type: local # (5)!
  - name: remote-server-1
    connection:
      type: remote
      user: myuser # (6)!
      ip: 10.10.40.143 # (7)!
      ssh:
        port: 1234  # (8)!
        verify: true # (9)!
        keyfile: "~/.ssh/id_rsa_server1" # (10)!
  - name: remote-server-2
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.144
      ssh:
        keyfile: "~/.ssh/id_rsa_server2"
    mainResourcePoolPath: "/var/lib/libvirt/pools/" # (11)!
    dataResourcePools: # (12)!
      - name: data-pool # (13)!
        path: "/mnt/data/pool" # (14)!
      - name: backup-pool
        path: "/mnt/backup/pool"

#
# The 'cluster' section of the configuration contains general data about the
# cluster, the nodes that are part of the cluster, and the cluster's network.
#
cluster:
  name: my-k8s-cluster # (15)!
  network:
    mode: bridge # (16)!
    cidr: 10.10.64.0/24 # (17)!
    gateway: 10.10.64.1 # (18)!
    bridge: br0 # (19)!
  nodeTemplate:
    user: k8s
    ssh:
      privateKeyPath: "~/.ssh/id_rsa_test"
      addToKnownHosts: true
    os:
      distro: ubuntu22
      networkInterface: ens3 # (20)!
    dns: # (21)!
      - 1.1.1.1
      - 1.0.0.1
    updateOnBoot: true
  nodes:
    loadBalancer:
      vip: 10.10.64.200 # (22)!
      virtualRouterId: 13 # (23)!
      forwardPorts:
        - name: http
          port: 80
        - name: https
          port: 443
          target: all
        - name: sample
          port: 60000
          targetPort: 35000
      default: # (24)!
        ram: 4 # GiB
        cpu: 1 # vCPU
        mainDiskSize: 16 # GiB
      instances:
        - id: 1
          ip: 10.10.64.5 # (25)!
          mac: "52:54:00:00:00:40" # (26)!
          ram: 8 # (27)!
          cpu: 8 # (28)!
          host: remote-server-1 # (29)!
        - id: 2
          ip: 10.10.64.6
          mac: "52:54:00:00:00:41"
          host: remote-server-2
        - id: 3
          ip: 10.10.64.7
          mac: "52:54:00:00:00:42"
          # If host is not specifed, VM will be installed on the default host.
          # If default host is not specified, VM will be installed on the first
          # host in the list.
    master:
      default:
        ram: 8
        cpu: 2
        mainDiskSize: 256
      instances:
          # IMPORTANT: There should be odd number of master nodes.
        - id: 1
          host: remote-server-1
        - id: 2
          host: remote-server-2
        - id: 3
          host: localhost
    worker:
      default:
        ram: 16
        cpu: 4
        labels: # (30)!
          custom-label: "This is a custom default node label"
          node-role.kubernetes.io/node: # (31)!
      instances:
        - id: 1
          ip: 10.10.64.101
          cpu: 8
          ram: 64
          host: remote-server-1
        - id: 2
          ip: 10.10.64.102
          dataDisks: # (32)!
            - name: rook-disk # (33)!
              pool: data-pool # (34)!
              size: 128 # GiB
            - name: test-disk
              pool: data-pool
              size: 128
        - id: 3
          ip: 10.10.64.103
          ram: 64
          labels:
            custom-label: "Overwrite default node label" # (35)!
            instance-label: "Node label, only for this instance"
        - id: 4
          host: remote-server-2
        - id: 5

#
# The 'kubernetes' section contains Kubernetes related properties,
# such as version and network plugin.
#
kubernetes:
  version: v1.28.6
  networkPlugin: calico
  dnsMode: coredns # (36)!
  other:
    mergeKubeconfig: true

#
# The 'addons' section contains the configuration of the applications that
# will be installed on the Kubernetes cluster as part of the cluster setup.
#
addons:
  kubespray:
    # Sample Nginx ingress controller deployment
    ingress_nginx_enabled: true
    ingress_nginx_namespace: "ingress-nginx"
    ingress_nginx_insecure_port: 80
    ingress_nginx_secure_port: 443
    # Sample MetalLB deployment
    metallb_enabled: true
    metallb_speaker_enabled: true
    metallb_ip_range:
      - "10.10.9.201-10.10.9.254"
    metallb_pool_name: "default"
    metallb_auto_assign: true
    metallb_version: v0.12.1
    metallb_protocol: "layer2"
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

9.  If true, SSH host is verified. This means that the host must be present in the known SSH hosts.

10. Path to the **passwordless** SSH key used to connect to the remote host.

11. The path to the main resource pool defines where the virtual machine disk images are stored. These disks contain the virtual machine operating system, and therefore it is recommended to install them on SSD disks.

12. List of other data resource pools where virtual disks can be created.

13. Custom data resource pool name. Must be unique among all data resource pools on a specific host.

14. Path where data resource pool is created. All data disks linked to that resource pool will be created under this path.

15. Cluster name used as a prefix for the various components.

16. Network mode. Possible values are

    + `bridge` mode uses **predefined** bridge interface. This mode is mandatory for deployments across multiple hosts.
    + `nat` mode creates virtual network with IP range defined in `network.cidr`
    + `route`

17. Network CIDR represents the network IP together with the network mask.
    In `nat` mode, CIDR is used for the new network.
    In `bridge` mode, CIDR represents the current local area network (LAN).

18. The network gateway IP address.
    If omitted the first client IP from network CIDR is used as a gateway.

19. Bridge represents the bridge interface on the hosts.
    This field is mandatory if the network mode is set to `bridge`.
    If the network mode is set to `nat`, this field can be omitted.

20. Set custom DNS list for all nodes.
    If omitted, network gateway is also used as a DNS.

21. Specify the network interface used by the virtual machine. In general, this option can be omitted.

    If omitted, a network interface from distro preset (`/terraform/defaults.yaml`) is used.

22. Virtual (floating) IP shared between load balancers.

23. Virtual router ID that is set in Keepalived configuration when virtual IP is used.
    By default it is set to 51.
    If multiple clusters are created it must be ensured that it is unique for each cluster.

24. Default values apply for all virtual machines (VMs) of the same type.

25. Static IP address of the virtual machine.
    If omitted DHCP lease is requested.

26. Static MAC address.
    If omitted MAC address is generated.

27. Overrides default RAM value for this node.

28. Overrides default CPU value for this node.

29. Name of the host where instance should be created.
    If omitted the default host is used.

30. Default worker node labels.

31. Label sets worker nodes role to `node`.

32. Overrides default data disks for this node.

33. Custom data disk name. It must be unique among all data disks for a specific instance.

34. Resource pool name that must be defined on the host on which the instance will be deployed.

35. Node labels defined for specific instances take precedence over default labels with the same key, so this label overrides the default label.

36. Currently, the only DNS mode supported is CoreDNS.
