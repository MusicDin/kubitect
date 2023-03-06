---
hide:
  # - navigation
  - toc
---

<div markdown="1" class="text-center">
# Configuration reference
</div>

<div markdown="1" class="text-justify">

This document contains a reference of the Kubitect configuration file and documents all possible configuration properties.

The configuration sections are as follows:

+ `kubitect` - Project metadata.
+ `hosts` - A list of physical hosts (local or remote).
+ `cluster` - Configuration of the cluster infrastructure. Virtual machine properties, node types to install, and the host on which to install the nodes.
+ `kubernetes` - Kubernetes configuration.
+ `addons` - Configurable addons and applications.

Each configuration property is documented with 5 columns: Property name, description, type, default value and is the property required.

!!! note "Note"

    `[*]` annotates an array.

</div>

## *Kubitect* section

<table>
  <tbody>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Default value</th>
      <th>Required?</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>kubitect.url</code></td>
      <td>string</td>
      <td>https://github.com/MusicDin/kubitect</td>
      <td>No</td>
      <td>URL of the project's git repository.</td>
    </tr>
    <tr>
      <td><code>kubitect.version</code></td>
      <td>string</td>
      <td>
        <i>CLI tool version</i>
      </td>
      <td>No</td>
      <td>Version of the git repository. Can be a branch or a tag.</td>
    </tr>
  </tbody>
</table>

## *Hosts* section

<table>
  <tbody>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Default value</th>
      <th>Required?</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>hosts[*].connection.ip</code></td>
      <td>string</td>
      <td></td>
      <td>Yes, if <code>connection.type</code> is set to <code>remote</code></td>
      <td>IP address is used to SSH into the remote machine.</td>
    </tr>
    <tr>
      <td><code>hosts[*].connection.ssh.keyfile</code></td>
      <td>string</td>
      <td>~/.ssh/id_rsa</td>
      <td></td>
      <td>Path to the keyfile that is used to SSH into the remote machine</td>
    </tr>
    <tr>
      <td><code>hosts[*].connection.ssh.port</code></td>
      <td>number</td>
      <td>22</td>
      <td></td>
      <td>The port number of SSH protocol for remote machine.</td>
    </tr>
    <tr>
      <td><code>hosts[*].connection.ssh.verify</code></td>
      <td>boolean</td>
      <td>false</td>
      <td></td>
      <td>
        If true, the SSH host is verified, which means that the host must be present in the known SSH hosts.
      </td>
    </tr>
    <tr>
      <td><code>hosts[*].connection.type</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Possible values are:
        <ul>
            <li><code>local</code> or <code>localhost</code></li>
            <li><code>remote</code></li>
        </ul>
      </td>
    </tr>
    <tr>
      <td><code>hosts[*].connection.user</code></td>
      <td>string</td>
      <td></td>
      <td>Yes, if <code>connection.type</code> is set to <code>remote</code></td>
      <td>Username is used to SSH into the remote machine.</td>
    </tr>
    <tr>
      <td><code>hosts[*].dataResourcePools[*].name</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Name of the data resource pool. Must be unique within the same host.
        It is used to link virtual machine volumes to the specific resource pool.
      </td>
    </tr>
    <tr>
      <td><code>hosts[*].dataResourcePools[*].path</code></td>
      <td>string</td>
      <td>/var/lib/libvirt/images/</td>
      <td></td>
      <td>Host path to the location where data resource pool is created.</td>
    </tr>
    <tr>
      <td><code>hosts[*].default</code></td>
      <td>boolean</td>
      <td>false</td>
      <td></td>
      <td>
        Nodes where host is not specified will be installed on default host. 
        The first host in the list is used as a default host if none is marked as a default.
      </td>
    </tr>
    <tr>
      <td><code>hosts[*].name</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Custom server name used to link nodes with physical hosts.</td>
    </tr>
    <tr>
      <td><code>hosts[*].mainResourcePoolPath</code></td>
      <td>string</td>
      <td>/var/lib/libvirt/images/</td>
      <td></td>
      <td>Path to the resource pool used for main virtual machine volumes.</td>
    </tr>
  </tbody>
</table>


## *Cluster* section

<table>
  <tbody>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Default value</th>
      <th>Required?</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>cluster.name</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Custom cluster name that is used as a prefix for various cluster components.</td>
    </tr>
    <!-- Cluster network -->
    <tr>
      <td><code>cluster.network.bridge</code></td>
      <td>string</td>
      <td>virbr0</td>
      <td></td>
      <td>
        By default virbr0 is set as a name of virtual bridge.
        In case network mode is set to bridge, name of the preconfigured bridge needs to be set here.
      </td>
    </tr>
    <tr>
      <td><code>cluster.network.cidr</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Network cidr that contains network IP with network mask bits (IPv4/mask_bits).</td>
    </tr>
    <tr>
      <td><code>cluster.network.gateway</code></td>
      <td>string</td>
      <td><i>First client IP in network.</i></td>
      <td></td>
      <td>
        By default first client IP is taken as a gateway.
        If network cidr is set to 10.0.0.0/24 then gateway would be 10.0.0.1.
        Set gateway if it differs from default value.
      </td>
    </tr>
    <tr>
      <td><code>cluster.network.mode</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>
        Network mode. Possible values are:
        <ul>
          <li><code>nat</code> - Creates virtual local network.</li>
          <li><code>bridge</code> - Uses preconfigured bridge interface on the machine (Only bridge mode supports multiple hosts).</li>
          <li><code>route</code> - Creates virtual local network, but does not apply NAT.</li>
        </ul>
      </td>
    </tr>
    <!-- Cluster nodes (loadBalancer) -->
    <tr>
      <td><code>cluster.nodes.loadBalancer.default.cpu</code></td>
      <td>number</td>
      <td>2</td>
      <td></td>
      <td>Default number of vCPU allocated to a load balancer instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.default.mainDiskSize</code></td>
      <td>number</td>
      <td>32</td>
      <td></td>
      <td>Size of the main disk (in GiB) that is attached to a load balancer instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.default.ram</code></td>
      <td>number</td>
      <td>4</td>
      <td></td>
      <td>Default amount of RAM (in GiB) allocated to a load balancer instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.forwardPorts[*].name</code></td>
      <td>string</td>
      <td></td>
      <td>Yes, if port is configured</td>
      <td>Unique name of the forwarded port.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.forwardPorts[*].port</code></td>
      <td>number</td>
      <td></td>
      <td>Yes, if port is configured</td>
      <td>Incoming port is the port on which a load balancer listens for the incoming traffic.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.forwardPorts[*].targetPort</code></td>
      <td>number</td>
      <td><i>Incoming port value</i></td>
      <td></td>
      <td>Target port is the port on which a load balancer forwards traffic.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.forwardPorts[*].target</code></td>
      <td>string</td>
      <td>workers</td>
      <td></td>
      <td>
        Target is a group of nodes on which a load balancer forwards traffic.
        Possible targets are:
        <ul>
          <li><code>masters</code></li>
          <li><code>workers</code></li>
          <li><code>all</code></li>
        </ul>
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].cpu</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].host</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Name of the host on which the instance is deployed. 
        If the name is not specified, the instance is deployed on the default host.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].id</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>
        Unique identifier of a load balancer instance.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].ip</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        If an IP is set for an instance then the instance will use it as a static IP.
        Otherwise it will try to request an IP from a DHCP server.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].mac</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>MAC used by the instance. If it is not set, it will be generated.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].mainDiskSize</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].priority</code></td>
      <td>number</td>
      <td>10</td>
      <td></td>
      <td>
        Keepalived priority of the load balancer.
        A load balancer with the highest priority becomes the leader (active). 
        The priority can be set to any number between 0 and 255.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.instances[*].ram</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for the RAM for that instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.vip</code></td>
      <td>string</td>
      <td></td>
      <td>Yes, if more then one instance of load balancer is specified.</td>
      <td>
        Virtual IP (floating IP) is the static IP used by load balancers to provide a fail-over.
        Each load balancer still has its own IP beside the shared one.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.loadBalancer.virtualRouterId</code></td>
      <td>number</td>
      <td>51</td>
      <td></td>
      <td>
        Virtual router ID identifies the group of VRRP routers.
        It can be any number between 0 and 255 and should be unique among different clusters.
      </td>
    </tr>
    <!-- Cluster nodes (master) -->
    <tr>
      <td><code>cluster.nodes.master.default.cpu</code></td>
      <td>number</td>
      <td>2</td>
      <td></td>
      <td>Default number of vCPU allocated to a master node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.default.labels</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>
        Array of default node labels that are applied to all master nodes.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.default.mainDiskSize</code></td>
      <td>number</td>
      <td>32</td>
      <td></td>
      <td>Size of the main disk (in GiB) that is attached to a master node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.default.ram</code></td>
      <td>number</td>
      <td>4</td>
      <td></td>
      <td>Default amount of RAM (in GiB) allocated to a master node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.default.taints</code></td>
      <td>list</td>
      <td></td>
      <td></td>
      <td>
        List of default node taints that are applied to all master nodes.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].cpu</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].dataDisks[*].name</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>Name of the additional data disk that is attached to the master node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].dataDisks[*].pool</code></td>
      <td>string</td>
      <td>main</td>
      <td></td>
      <td>
        Name of the data resource pool where the additional data disk is created.
        Referenced resource pool must be configure on the same host.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].dataDisks[*].size</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Size of the additional data disk (in GiB) that is attached to the master node.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].host</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Name of the host on which the instance is deployed. 
        If the name is not specified, the instance is deployed on the default host.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].id</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Unique identifier of a master node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].ip</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        If an IP is set for an instance then the instance will use it as a static IP.
        Otherwise it will try to request an IP from a DHCP server.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].labels</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>
        Array of node labels that are applied to this specific master node.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].mac</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>MAC used by the instance. If it is not set, it will be generated.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].mainDiskSize</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].ram</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for the RAM for that instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.master.instances[*].taints</code></td>
      <td>list</td>
      <td></td>
      <td></td>
      <td>
        List of node taints that are applied to this specific master node.
      </td>
    </tr>
    <!-- Cluster nodes (worker) -->
    <tr>
      <td><code>cluster.nodes.worker.default.cpu</code></td>
      <td>number</td>
      <td>2</td>
      <td></td>
      <td>Default number of vCPU allocated to a worker node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.default.labels</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>Array of default node labels that are applied to all worker nodes.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.default.mainDiskSize</code></td>
      <td>number</td>
      <td>32</td>
      <td></td>
      <td>Size of the main disk (in GiB) that is attached to a worker node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.default.ram</code></td>
      <td>number</td>
      <td>4</td>
      <td></td>
      <td>Default amount of RAM (in GiB) allocated to a worker node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.default.taints</code></td>
      <td>list</td>
      <td></td>
      <td></td>
      <td>
        List of default node taints that are applied to all worker nodes.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].cpu</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].dataDisks[*].name</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>Name of the additional data disk that is attached to the worker node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].dataDisks[*].pool</code></td>
      <td>string</td>
      <td>main</td>
      <td></td>
      <td>
        Name of the data resource pool where the additional data disk is created.
        Referenced resource pool must be configure on the same host.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].dataDisks[*].size</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Size of the additional data disk (in GiB) that is attached to the worker node.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].host</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Name of the host on which the instance is deployed. 
        If the name is not specified, the instance is deployed on the default host.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].id</code></td>
      <td>string</td>
      <td></td>
      <td>Yes</td>
      <td>Unique identifier of a worker node.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].ip</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        If an IP is set for an instance then the instance will use it as a static IP.
        Otherwise it will try to request an IP from a DHCP server.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].labels</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>
        Array of node labels that are applied to this specific worker node.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].mac</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>MAC used by the instance. If it is not set, it will be generated.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].mainDiskSize</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for that specific instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].ram</code></td>
      <td>number</td>
      <td></td>
      <td></td>
      <td>Overrides a default value for the RAM for that instance.</td>
    </tr>
    <tr>
      <td><code>cluster.nodes.worker.instances[*].taints</code></td>
      <td>list</td>
      <td></td>
      <td></td>
      <td>
        List of node taints that are applied to this specific worker node.
      </td>
    </tr>
    <!-- Cluster node template -->
    <tr>
      <td><code>cluster.nodeTemplate.cpuMode</code></td>
      <td>string</td>
      <td>custom</td>
      <td></td>
      <td>
        Guest virtual machine CPU mode.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.dns</code></td>
      <td>list</td>
      <td>Value of <code>network.gateway</code></td>
      <td></td>
      <td>
        Custom DNS list used by all created virtual machines.
        If none is provided, network gateway is used.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.os.distro</code></td>
      <td>string</td>
      <td>ubuntu</td>
      <td></td>
      <td>
        Set OS distribution. Possible values are:
        <ul>
          <li><code>ubuntu</code></li>
          <li><code>debian</code></li>
          <li>
            <code>custom</code> - For all other distros
            <i>(for development only)</i>
          </li>
        </ul>
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.os.networkInterface</code></td>
      <td>string</td>
      <td>Depends on <code>os.distro</code></td>
      <td></td>
      <td>
        Network interface used by virtual machines to connect to the network.
        Network interface is preconfigured for each OS image (usually ens3 or eth0).
        By default, the value from distro preset (<i>/terraform/defaults.yaml</i>) is set, but can be overwritten if needed.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.os.source</code></td>
      <td>string</td>
      <td>Depends on <code>os.distro</code></td>
      <td></td>
      <td>
        Source of an OS image. 
        It can be either path on a local file system or an URL of the image.
        By default, the value from distro preset (<i>/terraform/defaults.yaml</i>)isset, but can be overwritten if needed.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.ssh.addToKnownHosts</code></td>
      <td>boolean</td>
      <td>true</td>
      <td></td>
      <td>
        If set to true, each virtual machine will be added to the known hosts on the machine where the project is being run.
        Note that all machines will also be removed from known hosts when destroying the cluster.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.ssh.privateKeyPath</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Path to private key that is later used to SSH into each virtual machine.
        On the same path with <code>.pub</code> prefix needs to be present public key.
        If this value is not set, SSH key will be generated in <code>./config/.ssh/</code> directory.
      </td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.updateOnBoot</code></td>
      <td>boolean</td>
      <td>true</td>
      <td></td>
      <td>If set to true, the operating system will be updated when it boots.</td>
    </tr>
    <tr>
      <td><code>cluster.nodeTemplate.user</code></td>
      <td>string</td>
      <td>k8s</td>
      <td></td>
      <td>User created on each virtual machine.</td>
    </tr>
  </tbody>
</table>


## *Kubernetes* section

<table>
  <tbody>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Default value</th>
      <th>Required?</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>kubernetes.dnsMode</code></td>
      <td>string</td>
      <td>coredns</td>
      <td></td>
      <td>
        DNS server used within a Kubernetes cluster. Possible values are: 
        <ul>
          <li><code>coredns</code></li>
        </ul>
      </td>
    </tr>
    <tr>
      <td><code>kubernetes.networkPlugin</code></td>
      <td>string</td>
      <td>calico</td>
      <td></td>
      <td>
        Network plugin used within a Kubernetes cluster. Possible values are: 
        <ul>
          <li><code>flannel</code></li>
          <li><code>weave</code></li>
          <li><code>calico</code></li>
          <li><code>cilium</code></li>
          <li><code>canal</code></li>
          <li><code>kube-router</code></li>
        </ul>
      </td>
    </tr>
    <tr>
      <td><code>kubernetes.other.autoRenewCertificates</code></td>
      <td>boolean</td>
      <td>false</td>
      <td></td>
      <td>
        When this property is set to true, control plane certificates are renewed first Monday of each month.
      </td>
    </tr>
    <tr>
      <td><code>kubernetes.other.copyKubeconfig</code></td>
      <td>boolean</td>
      <td>false</td>
      <td></td>
      <td>
        When this property is set to true, the kubeconfig of a new cluster is copied to the <code>~/.kube/config</code>.
        Please note that setting this property to true may cause the existing file at the destination to be overwritten.
      </td>
    </tr>
    <tr>
      <td><code>kubernetes.version</code></td>
      <td>string</td>
      <td>v1.25.6</td>
      <td></td>
      <td>Kubernetes version that will be installed.</td>
    </tr>
  </tbody>
</table>


## *Addons* section

<table>
  <tbody>
    <tr>
      <th>Name</th>
      <th>Type</th>
      <th>Default value</th>
      <th>Required?</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>addons.kubespray</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>
        Kubespray addons configuration.
      </td>
    </tr>
    <tr>
      <td><code>addons.rook.enabled</code></td>
      <td>boolean</td>
      <td>false</td>
      <td></td>
      <td>
        Enable Rook addon.
      </td>
    </tr>
    <tr>
      <td><code>addons.rook.nodeSelector</code></td>
      <td>dictionary</td>
      <td></td>
      <td></td>
      <td>
        Dictionary containing node labels ("key: value").
        Rook is deployed on the nodes that match all the given labels.
      </td>
    </tr>
    <tr>
      <td><code>addons.rook.version</code></td>
      <td>string</td>
      <td></td>
      <td></td>
      <td>
        Rook version.
        By default, the latest release version is used.
      </td>
    </tr>
  </tbody>
</table>
