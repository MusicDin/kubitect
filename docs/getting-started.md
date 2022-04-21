# Getting Started

In this step-by-step guide, you will learn how to prepare a custom cluster configuration file and use it to create a functional Kubernetes cluster consisting of a single master node and three worker nodes.

> :scroll: **Note:**
Detailed example and explanations of each possible configuration property can be found in the [configuration documentation](/docs/configuration.md).

> :scroll: **Note:**
For the successful installation of the Kubernetes cluster, some [requirements](/docs/requirements.md) must be met.

## Step 1 - Create cluster configuration file

In the quick start you have created a very basic Kubernetes cluster from predefined cluster configuration file.
If configuration is not explicitly provided to the command-line tool using `--config` option, default cluster configuration file is used ([/examples/default-cluster.yaml](/examples/default-cluster.yaml)).

Now it's time to create your own cluster topology.

Before you begin create a new yaml file.
```sh
touch tkk.yaml
```

## Step 2 - Prepare hosts configuration

In the cluster configuration file, we will first define hosts.
Hosts represent target servers.
A host can be either a local or a remote machine.

If the cluster is set up on the same machine where the command line tool is installed, we specify a host whose connection type is set to `local`.
```yaml
hosts:
  - name: localhost # Can be anything
    connection:
      type: local
```

When cluster is deployed on the remote machine, the IP address of the remote machine along with the SSH credentails needs to be specified for the host.
```yaml
hosts:
  - name: my-remote-host
    connection:
      type: remote
      user: myuser
      ip: 10.10.40.143 # IP address of the remote host
      ssh:
        keyfile: "~/.ssh/id_rsa_server1" # Password-less SSH key file
```

In this tutorial we will use only a localhost.

## Step 3 - Define cluster infrastructure

The second part of the configuration file consists of the cluster infrastructure.
In this part, all virtual machines are defined along with their properties such as operating system, CPU cores, amount of RAM and so on.

Let's take a look at the following configuration:
```yaml
cluster:
  name: "my-k8s-cluster"
  network:
    ...
  nodeTemplate:
    ...
  nodes:
    ...
```

We can see that the infrastructure configuration consists of the cluster name and 3 subsections:
- `cluster.name` is a cluster name that is used as a prefix for each resource created by *tkk*.
- `cluster.network` holds information about the network properties of the cluster.
- `cluster.nodeTemplate` contains properties that apply to all our nodes. For example, properties like operating system, SSH user, and SSH private key are the same for all our nodes.
- `cluster.nodes` subsection defines each node in our cluster.

Now that we have a general idea about the infrastructure configuration, we can look at each of these subsections in more detail.

### Step 3.1 - Cluster network

The cluster network subsection defines the network that our cluster will use.
Currently, two network modes are supported - NAT and bridge.

The `nat` network mode instructs *tkk* to create a virtual network that does network address translation. This mode allows us to use IP address ranges that do not exist within our local area network (LAN).

The `bridge` network mode instructs *tkk* to use a predefined bridge interface.
In this mode, virtual machines can connect directly to LAN.
Using this mode is mandatory when you set up a cluster that spreads over multipe hosts.

To keep this tutorial as simple as possible, we will use the NAT mode, as it does not require a preconfigured bridge interface.

```yaml
cluster:
  ...
  network:
    mode: "nat"
    cidr: "192.168.113.0/24"
```

The above configuration will instruct *tkk* to create a virtual network that uses `192.168.113.0/24` IP range.

### Step 3.2 - Node template

As mentioned earlier, the `nodeTemplate` subsection is used to define general properties of our nodes.

Required properties are:
+ `user` is the name of the user that will be created on all virtual machines and will also be used for SSH.
+ `image.distro` defines the type of the used operating system (ubuntu, debian, ...).
+ `image.source` defines the location of the OS image. It can be either a local file system path or an URL.

Besides the required properties, there are two potentially useful properties:
+ `ssh.addToKnownHosts` - if set to true, all virtual machines will be added to SSH known hosts. If you later destroy the cluster, these virtual machines will also be removed from the known hosts.
+ `updateOnBoot` - if set to true, all virtual machines are updated at the first boot.

Our `noteTemplate` subsection now looks like this:
```yaml
cluster:
  ...
  nodeTemplate:
    user: "k8s"
    ssh:
      addToKnownHosts: true
    image:
      distro: "ubuntu"
      source: "https://cloud-images.ubuntu.com/releases/focal/release-20220111/ubuntu-20.04-server-cloudimg-amd64.img"
    updateOnBoot: true
```

### Step 3.3 - Cluster nodes

In the `nodes` subsection, we can define three types of nodes:
- `loadBalancer` nodes are internal load balancers used to expose the Kubernetes control plane at a single endpoint.
- `master` nodes are Kubernetes master nodes that also contain an etcd key-value store. Since etcd is present on these nodes, the number of master nodes must be odd. For more information, see [etcd FAQ](https://etcd.io/docs/v3.4/faq/#why-an-odd-number-of-cluster-members).
- `worker` nodes are the nodes on which your actual workload runs.

In this tutorial, we will use only one master node, so internal load balancers are not required. 

The easiest way to explain this part is to look at the actual configuration:
```yaml
cluster:
  ...
  nodes:
    master:
      default: # Default properties of all master node instances
        ram: 4
        cpu: 2
        mainDiskSize: 32
      instances: # Master node instances
        - id: 1
          ip: 192.168.113.10
    worker:
      default:
        ram: 4
        cpu: 2
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.21
          cpu: 4  # Override default vCPU value for this node
          ram: 8  # Override default amount of RAM for this node
        - id: 7
          ip: 192.168.113.27
          mac: "52:54:00:00:00:42" # Specify MAC address for this node
        - id: 99
          # If ip property is omitted, node will request an IP address from the DHCP server.
          # If mac property is omitted, MAC address will be auto generated.
```

### Step 3.4 - Kubernetes properties

The last part of the cluster configuration consists of the Kubernetes properties.
In this section we define the Kubernetes version, the DNS plugin and so on.
It is also important to check if Kubespray supports a specific Kubernetes version.

If you are using a custom Kubespray, you can also specify the URL to a custom Git repository.

```yaml
kubernetes:
  version: "v1.22.6"
  networkPlugin: "calico"
  dnsMode: "coredns"
  kubespray:
    version: "v2.18.1"
    # url: URL to custom Kubespray git repository (default is: https://github.com/kubernetes-sigs/kubespray.git)
```

### Step 4 - Create the cluster

Our final cluster configuration looks like this:
```yaml
# tkk.yaml
hosts:
  - name: localhost
    connection:
      type: local

cluster:
  network:
    mode: "nat"
    cidr: "192.168.113.0/24"
  nodeTemplate:
    user: "k8s"
    ssh:
      addToKnownHosts: true
    image:
      distro: "ubuntu"
      source: "https://cloud-images.ubuntu.com/releases/focal/release-20220111/ubuntu-20.04-server-cloudimg-amd64.img"
    updateOnBoot: true
  nodes:
    master:
      default:
        ram: 4
        cpu: 2
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.10
    worker:
      default:
        ram: 4
        cpu: 2
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.21
          cpu: 4
          ram: 8
        - id: 7
          ip: 192.168.113.27
          mac: "52:54:00:00:00:42"
        - id: 99

kubernetes:
  version: "v1.22.6"
  networkPlugin: "calico"
  dnsMode: "coredns"
  kubespray:
    version: "v2.18.1"
```

Now create the cluster by applying your custom configuration using the *tkk* command line tool. Also, let's name our cluster `my-first-cluster`.
```
tkk apply --cluster my-first-cluster --config tkk.yaml
```

> :bulb: **Tip:** 
If you encounter any issues during the installation, please refer to the [troubleshooting](docs/troubleshooting.md) page first.

When the cluster is applied, it is created in *tkk* home directory, which has the following structure.
```
~/.tkk
   ├── clusters
   │   ├── default
   │   ├── my-first-cluster
   │   └── ...
   └── bin
       └── ...
```

### Step 5 - Test the cluster


Using *tkk* command line tool, list all created clusters.
```sh
tkk list clusters
```

After successful installation of the Kubernetes cluster, Kubeconfig will be created within cluster's directory.
To export the Kubeconfig into a custom file run the following command.
```
tkk export kubeconfig --cluster my-first-cluster > kubeconfig.yaml
```

Use the exported Kubeconfig to list all cluster nodes.
```
kubectl get nodes --kubeconfig kubeconfig.yaml
```

:clap: Congratulations, you have completed the *getting started* tutorial! 

## What's next?

+ [Learn how to manage created clusters](./cluster-management.md)
+ [See the configuration documentation](/docs/configuration.md)