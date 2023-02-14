This example shows how to use Kubitect to set up a highly available cluster that spreads over 5 hosts.
Such topology provides redundancy in case any node or even host fails.

<div align=center>
  <img 
    src="/assets/images/topology-ha-arch.png" 
    alt="Architecture of the highly available cluster"
    width="100%">
</div>

## Step 1: Hosts configuration


!!! warning "Important"

    This example uses **preconfigured bridges on each host** to expose nodes on the local network.
    
    [Network bridge](../network-bridge) example shows how to configure a bridge interface using Netplan.

In this example, we deploy a Kubernetes cluster on 5 (remote) physical hosts.
The subnet of the local network is `10.10.0.0/20` and the gateway IP address is `10.10.0.1`.
Each host is connected to the same local network and has a pre-configured bridge interface `br0`.

In addition, a user `kubitect` is configured on each host, which is accessible via SSH with the same password-less certificate stored on our local machine under the path `~/.ssh/id_rsa_ha`.

Each host must be specified in the Kubitect configuration file.
In our case, the configurations of the hosts differ only by the name and IP address of the host.

```yaml title="ha.yaml" 
hosts:
  - name: host1
    connection:
      type: remote
      user: kubitect
      ip: 10.10.0.5
      ssh:
        keyfile: "~/.ssh/id_rsa_ha"
  - name: host2
    connection:
      type: remote
      user: kubitect
      ip: 10.10.0.6
      ssh:
        keyfile: "~/.ssh/id_rsa_ha"
  - name: host3
    connection:
      type: remote
      user: kubitect
      ip: 10.10.0.10
      ssh:
        keyfile: "~/.ssh/id_rsa_ha"
  - name: host4
    connection:
      type: remote
      user: kubitect
      ip: 10.10.0.11
      ssh:
        keyfile: "~/.ssh/id_rsa_ha"
  - name: host5
    connection:
      type: remote
      user: kubitect
      ip: 10.10.0.12
      ssh:
        keyfile: "~/.ssh/id_rsa_ha"
```

## Step 2: Network configuration

In the network configuration section, we specify the bridge interface that is preconfigured on each host and CIDR of our local network.

The following code snippet shows the network configuration for this example.

```yaml title="ha.yaml" 
cluster:
  network:
    mode: bridge
    cidr: 10.10.0.0/20
    bridge: br0
```

## Step 3: Load balancer configuration

Each working node stores only a single control plane IP address. 
By placing a load balancer in front of the control plane (as shown in the [Multi-master cluster](../multi-master-cluster) example), traffic can be distributed across all control plane nodes.

By having only a single load balancer in the cluster, the control plane may become unreachable if that load balancer fails.
This would cause the entire cluster to become unavailable.
To avoid this single point of failure, a failover (backup) load balancer can be configured.
Its main purpose is to serve incoming requests on the same virtual (shared) IP address if the primary load balancer fails, as shown in the following figure.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/topology-ha-base.png" 
    alt="Scheme of highly available topology"
    width="75%">
</div>


Failover is achieved using a virtual router redundancy protocol (VRRP).
In practice, each load balancer still has its own IP address, but the primary load balancer also serves requests on the virtual IP address, which is not bound to any network interface.
The primary load balancer periodically sends heartbeats to other load balancers (backups) to let them know it is still active.
If the backup load balancer does not receive a heartbeat within a certain period of time, it assumes that the primary load balancer is has failed.
The new primary load balancer is elected based on the priorities of the available load balancers.
Once the load balancer becomes primary, it starts serving requests on the same virtual IP address as the previous one.
This ensures that the requests are served through the same virtual IP address in case of a load balancer failure.

The following code snippet shows the configuration of two load balancers and virtual IP for their failover.

```yaml title="ha.yaml" 
cluster:
  nodes:
    loadBalancer:
      vip: 10.10.13.200
      instances:
        - id: 1
          ip: 10.10.13.201
          host: host1
        - id: 2
          ip: 10.10.13.202
          host: host2
```

Note that for each load balancer instance, the `host` property is set.
Its value is a name of the host on which a particular instance is to be deployed.

## Step 4: Nodes configuration

The configuration of the nodes is very simple and similar to the configuration of the load balancer instances.
Each instance is configured with an ID, an IP address, and a host affinity.

```yaml title="ha.yaml" 
cluster:
  nodes:
    master:
      instances:
        - id: 1
          ip: 10.10.13.10
          host: host3
        - id: 2
          ip: 10.10.13.11
          host: host4
        - id: 3
          ip: 10.10.13.12
          host: host5
    worker:
      instances:
        - id: 1
          ip: 10.10.13.20
          host: host3
        - id: 2
          ip: 10.10.13.21
          host: host4
        - id: 3
          ip: 10.10.13.22
          host: host5
```

### Step 4.1 (Optional): Data disks configuration

Kubitect creates a main (system) disk for each configured virtual machine (VM).
The main disk contains the VM's operating system along with all installed Kubernetes components.

The VM's storage can be expanded by creating additional disks, also called data disks.
This can be particularly useful when using storage solutions such as Rook.
For example, Rook can be configured to use empty disks on worker nodes to create reliable distributed storage. 

Data disks in Kubitect must be configured separately for each node instance.
They must also be connected to a resource pool, which can be either a main resource pool or a custom data resource pool.
In this example, we have defined a custom data resource pool named `data-pool` on each host running worker nodes.

```yaml title="ha.yaml" 
hosts:
  - name: host3
    ...
    dataResourcePools:
      - name: data-pool
        path: /mnt/libvirt/pools/

cluster:
  nodes:
    worker:
      - id: 1
        ip: 10.10.13.20
        host: host3
        dataDisks:
          - name: rook
            pool: data-pool
            size: 512 # GiB
```


??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="ha.yaml" 
    hosts:
      - name: host1
        connection:
          type: remote
          user: kubitect
          ip: 10.10.0.5
          ssh:
            keyfile: "~/.ssh/id_rsa_ha"
      - name: host2
        connection:
          type: remote
          user: kubitect
          ip: 10.10.0.6
          ssh:
            keyfile: "~/.ssh/id_rsa_ha"
      - name: host3
        connection:
          type: remote
          user: kubitect
          ip: 10.10.0.10
          ssh:
            keyfile: "~/.ssh/id_rsa_ha"
        dataResourcePools:
          - name: data-pool
            path: /mnt/libvirt/pools/
      - name: host4
        connection:
          type: remote
          user: kubitect
          ip: 10.10.0.11
          ssh:
            keyfile: "~/.ssh/id_rsa_ha"
        dataResourcePools:
          - name: data-pool
            path: /mnt/libvirt/pools/
      - name: host5
        connection:
          type: remote
          user: kubitect
          ip: 10.10.0.12
          ssh:
            keyfile: "~/.ssh/id_rsa_ha"
        dataResourcePools:
          - name: data-pool
            path: /mnt/libvirt/pools/

    cluster:
      name: kubitect-ha
      network:
        mode: bridge
        cidr: 10.10.0.0/20
        bridge: br0
      nodeTemplate:
        user: k8s
        updateOnBoot: true
        ssh:
          addToKnownHosts: true
        os:
          distro: ubuntu
      nodes:
        loadBalancer:
          vip: 10.10.13.200
          instances:
            - id: 1
              ip: 10.10.13.201
              host: host1
            - id: 2
              ip: 10.10.13.202
              host: host2
        master:
          instances:
            - id: 1
              ip: 10.10.13.10
              host: host3
            - id: 2
              ip: 10.10.13.11
              host: host4
            - id: 3
              ip: 10.10.13.12
              host: host5
        worker:
          instances:
            - id: 1
              ip: 10.10.13.20
              host: host3
              dataDisks:
                - name: rook
                  pool: data-pool
                  size: 512
            - id: 2
              ip: 10.10.13.21
              host: host4
              dataDisks:
                - name: rook
                  pool: data-pool
                  size: 512
            - id: 3
              ip: 10.10.13.22
              host: host5
              dataDisks:
                - name: rook
                  pool: data-pool
                  size: 512

    kubernetes:
      version: v1.23.7
      kubespray:
        version: v2.19.0
    ```

## Step 5: Applying the configuration

Apply the cluster configuration.
```sh
kubitect apply --config ha.yaml
```




