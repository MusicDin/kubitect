<div markdown="1" class="text-center">
# Getting Started
</div>

<div markdown="1" class="text-justify">

In this **step-by-step** guide, you will learn how to prepare a custom cluster configuration file from scratch and use it to create a functional Kubernetes cluster consisting of a **one master and one worker node**.

<div class="text-center">
  <img
    class="mobile-w-75"
    src="/assets/images/topology-1m1w-base.png" 
    alt="Base scheme of the cluster with one master and one worker node"
    width="50%">
</div>

## Step 1 - Make sure all requirements are satisfied

For the successful installation of the Kubernetes cluster, some [requirements](../requirements) must be met.

## Step 2 - Create cluster configuration file

In the quick start you have created a very basic Kubernetes cluster from predefined cluster configuration file.
If configuration is not explicitly provided to the command-line tool using `--config` option, default cluster configuration file is used (`/examples/default-cluster.yaml`).

Now it's time to create your own cluster topology.

Before you begin create a new yaml file.
```sh
touch kubitect.yaml
```

## Step 3 - Prepare hosts configuration

In the cluster configuration file, we will first define hosts.
Hosts represent target servers.
A host can be either a local or a remote machine.

=== "Localhost"

    !!! quote ""

        If the cluster is set up on the same machine where the command line tool is installed, we specify a host whose connection type is set to `local`.

        ```yaml title="kubitect.yaml"
        hosts:
          - name: localhost # (1)!
            connection:
              type: local
        ```

        1. Custom **unique** name of the host.

=== "Remote host"

    !!! quote ""

        When cluster is deployed on the remote machine, the IP address of the remote machine along with the SSH credentails needs to be specified for the host.

        ```yaml title="kubitect.yaml"
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

        2. Path to the **password-less** SSH key file required for establishing connection with the remote host.

In this tutorial we will use only localhost.

## Step 4 - Define cluster infrastructure

The second part of the configuration file consists of the cluster infrastructure.
In this part, all virtual machines are defined along with their properties such as operating system, CPU cores, amount of RAM and so on.

For easier interpretation of the components that the final cluster will be made of, see the below image.

<div class="text-center">
  <img
    class="mobile-w-100"
    src="/assets/images/topology-1m1w-arch.png" 
    alt="Architecture of the cluster with one master and one worker node"
    width="75%">
</div>

Let's take a look at the following configuration:

```yaml title="kubitect.yaml"
cluster:
  name: local-k8s-cluster
  network:
    ...
  nodeTemplate:
    ...
  nodes:
    ...
```

We can see that the infrastructure configuration consists of the cluster name and 3 subsections:

- `cluster.name` is a cluster name that is used as a prefix for each resource created by Kubitect.
- `cluster.network` holds information about the network properties of the cluster.
- `cluster.nodeTemplate` contains properties that apply to all our nodes. For example, properties like operating system, SSH user, and SSH private key are the same for all our nodes.
- `cluster.nodes` subsection defines each node in our cluster.

Now that we have a general idea about the infrastructure configuration, we can look at each of these subsections in more detail.

### Step 4.1 - Cluster network

The cluster network subsection defines the network that our cluster will use.
Currently, two network modes are supported - NAT and bridge.

The `nat` network mode instructs Kubitect to create a virtual network that performs network address translation. This mode allows the use of IP address ranges that do not exist within our local area network (LAN).

The `bridge` network mode instructs Kubitect to use a predefined bridge interface.
In this mode, virtual machines can connect directly to LAN.
Use of this mode is mandatory when a cluster spreads over multipe hosts.

To keep this tutorial as simple as possible, we will use the NAT mode, as it does not require a preconfigured bridge interface.

```yaml title="kubitect.yaml"
cluster:
  ...
  network:
    mode: nat
    cidr: 192.168.113.0/24
```

The above configuration will instruct Kubitect to create a virtual network that uses `192.168.113.0/24` IP range.

### Step 4.2 - Node template

As mentioned earlier, the `nodeTemplate` subsection is used to define general properties of our nodes.

Required properties are:

+ `user` - the name of the user that will be created on all virtual machines and will also be used for SSH.

Besides the required properties, there are some potentially useful properties:

+ `os.distro` - defines the operating system for the nodes (currently ubuntu and debian are supported). By default, latest Ubuntu 22.04 release is used.
+ `ssh.addToKnownHosts` - if true, all virtual machines will be added to SSH known hosts. If you later destroy the cluster, these virtual machines will also be removed from the known hosts.
+ `updateOnBoot` - if true, all virtual machines are updated at the first boot.

Our `noteTemplate` subsection now looks like this:

```yaml title="kubitect.yaml"
cluster:
  ...
  nodeTemplate:
    user: k8s
    updateOnBoot: true
    ssh:
      addToKnownHosts: true
    os:
      distro: ubuntu22
```

### Step 4.3 - Cluster nodes

In the `nodes` subsection, we can define three types of nodes:

- `loadBalancer` nodes are internal load balancers used to expose the Kubernetes control plane at a single endpoint.
- `master` nodes are Kubernetes master nodes that also contain an etcd key-value store. Since etcd is present on these nodes, the number of master nodes must be odd. For more information, see [etcd FAQ](https://etcd.io/docs/v3.4/faq/#why-an-odd-number-of-cluster-members).
- `worker` nodes are the nodes on which your actual workload runs.

In this tutorial, we will use only one master node, so internal load balancers are not required. 

The easiest way to explain this part is to look at the actual configuration:

```yaml title="kubitect.yaml"
cluster:
  ...
  nodes:
    master:
      default: # (1)!
        ram: 4 # (2)!
        cpu: 2 # (3)!
        mainDiskSize: 32 # (4)!
      instances: # (5)!
        - id: 1 # (6)!
          ip: 192.168.113.10 # (7)!
    worker:
      default: 
        ram: 8
        cpu: 2
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.21
          ram: 4 # (8)!
```

1.  Default properties are applied to all nodes of the same type (in this case `master` nodes).
    They are especially useful, when multiple nodes of the same type are specified.

2.  Amount of RAM allocated to the master nodes (in GiB).

3.  Amount of vCPU allocated to the master nodes.

4.  Size of the virtual disk attached to each master node (in GiB).

5.  List of master node instances.

6.  Instance ID is the **only required field** that must be specified for each instance.

7.  Static IP address set for this particular instance.
    If the `ip` property is omitted, the DHCP lease is requested when the cluster is created.

8.  Since the amount of RAM (4 GiB) is specified for this particular instance, the default value (8 GiB) is overwritten.

### Step 4.4 - Kubernetes properties

The last part of the cluster configuration consists of the Kubernetes properties.
In this section we define the Kubernetes version, the DNS plugin and so on.
It is also important to check if Kubespray supports a specific Kubernetes version.

If you are using a custom Kubespray, you can also specify the URL to a custom Git repository.

```yaml title="kubitect.yaml"
kubernetes:
  version: v1.23.7
  networkPlugin: calico
  dnsMode: coredns
  kubespray:
    version: v2.19.0
```

## Step 5 - Create the cluster

!!! tip "Tip"

    If you encounter any issues during the installation, please refer to the [troubleshooting](../other/troubleshooting) page first.

Our final cluster configuration looks like this:

```yaml title="kubitect.yaml"
hosts:
  - name: localhost
    connection:
      type: local

cluster:
  name: local-k8s-cluster
  network:
    mode: nat
    cidr: 192.168.113.0/24
  nodeTemplate:
    user: k8s
    updateOnBoot: true
    ssh:
      addToKnownHosts: true
    os:
      distro: ubuntu22
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
        ram: 8
        cpu: 2
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.21
          ram: 4

kubernetes:
  version: v1.23.7
  networkPlugin: calico
  dnsMode: coredns
  kubespray:
    version: v2.19.0
```

Now create the cluster by applying your custom configuration using the Kubitect command line tool. 
Also, let's name our cluster `my-first-cluster`.
```sh
kubitect apply --cluster my-first-cluster --config kubitect.yaml
```

When the configuration is applied, the Terraform plan shows the changes Terraform wants to make to your infrastructure.
User confirmation of the plan is required before Kubitect begins creating the cluster.

!!! tip "Tip"

    To skip the user confirmation step, the flag `--auto-approve` can be used.

When the cluster is applied, it is created in Kubitect's home directory, which has the following structure.
```
~/.kubitect
   ├── clusters
   │   ├── default
   │   ├── my-first-cluster
   │   └── ...
   └── bin
       └── ...
```

All created clusters can be listed at any time using the following command.
```sh
kubitect list clusters

# Clusters:
#   - my-first-cluster (active)
```

## Step 6 - Test the cluster

After successful installation of the Kubernetes cluster, Kubeconfig is created in the cluster's directory.

To export the Kubeconfig to a separate file, run the following command.
```sh
kubitect export kubeconfig --cluster my-first-cluster > kubeconfig.yaml
```

Use the exported Kubeconfig to list all cluster nodes.
```sh
kubectl get nodes --kubeconfig kubeconfig.yaml
```

:clap: Congratulations, you have completed the *getting started* tutorial.

</div>
