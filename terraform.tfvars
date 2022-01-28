#
# These variables are ignored by default, because currently the cluster.yml
# configuration is used. To use terraform.tfvars as an input, generate appropriate
# main.tf file by running:
#
#   "sh scripts/tkk.sh generate-tf"
#
# The limitation of this approach is that it is not possible to deploy
# Kubernetes cluster on multiple physical machines (servers) and is therefor
# only possible to setup the cluster on a local or a single remote machine.
#

#======================================================================================
# General configuration
#======================================================================================

# Provider's URI. #
libvirt_provider_uri = "qemu:///system"

# Location where resource pool will be initialized. #
libvirt_resource_pool_location = "/var/lib/libvirt/pools/"


#======================================================================================
# Cluster infrastructure configuration
#======================================================================================

# Cluster name, used as prefix for various component names. #
cluster_name = "k8s"


#================================
# Network configuration
#================================

# Network mode (nat, route, bridge). #
cluster_network_mode = "nat"

# Network CIDR (example: 10.0.13.0/24). #
cluster_network_cidr = "10.0.13.0/24"

# Network gateway (example: 10.0.13.1). #
# Note: If not provided, it will be calculated as a first host in the network range. #
#       +-> first host of 10.0.13.0/24 is 10.0.13.1 #
#network_gateway = "10.0.13.1"

# Network (virtual) bridge. #
# Note: For network mode 'bridge', bridge on host needs to preconfigured (example: br0) #
cluster_network_bridge = "virbr1"

# Network DNS list (if empty, network gateway is set as a DNS). #
cluster_network_dns = [
  "1.1.1.1",
  "1.0.0.1"
]


#================================
# Node template
#================================

# Username used to SSH to the VM. #
cluster_nodeTemplate_user = "k8s"

# Private SSH key location (for VMs) (example: ~/.ssh/id_rsa). #
cluster_nodeTemplate_ssh_privateKeyPath = "/home/dinmusic/.ssh/id_rsa_dev_k8s"

# Add VMs to SSH known hosts. #
cluster_nodeTemplate_ssh_addToKnownHosts = true

# Linux distribution that will be used on VMs (ubuntu, debian, centos, ""). #
cluster_nodeTemplate_image_distro = "ubuntu"

# Source of linux image. It can be path to an image on host's filesystem or an URL. #
cluster_nodeTemplate_image_source = "https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64-disk-kvm.img"

# Network interface used by VMs to connect to the network. #
cluster_nodeTemplate_networkInterface = "ens3"

# If true, system will be updated on boot. #
cluster_nodeTemplate_updateOnBoot = false


#======================================================================================
# HAProxy internal load balancer (iLB) nodes
#======================================================================================

# Virtual/Floating IP address. #
# Note: Floating IP only applies if at least one load balancer is defined, #
# otherwise IP of the first master node will be used as control plane endpoint. #
cluster_nodes_loadBalancer_vip = "10.0.13.200"

# The default number of vCPU allocated to the load balancer VM. #
cluster_nodes_loadBalancer_default_cpu = 1

# The default amount of RAM allocated to the load balancer VM [GiB]. #
cluster_nodes_loadBalancer_default_ram = 2

# The default amount of disk allocated to the load balancer VM [GiB]. #
cluster_nodes_loadBalancer_default_storage = 16

# HAProxy load balancer nodes configuration. #
cluster_nodes_loadBalancer_instances = [
  {
    id  = 1
    ip  = "10.0.13.5"
    mac = "52:54:00:00:00:05"
  },
  {
    id  = 2
    ip  = "10.0.13.6"
    mac = "52:54:00:00:00:06"
  }
]


#======================================================================================
# Master nodes (control plane)
#======================================================================================

# The default number of vCPU allocated to the master VM. #
cluster_nodes_master_default_cpu = 2

# The default amount of RAM allocated to the master VM [GiB]. #
cluster_nodes_master_default_ram = 2

# The default amount of disk allocated to the master VM [GiB]. #
cluster_nodes_master_default_storage = 16

# Master nodes configuration #
# Note that number of masters cannot be divisible by 2. #
cluster_nodes_master_instances = [
  {
    id  = 1
    ip  = "10.0.13.10"
    mac = "52:54:00:00:00:10"
  },
  {
    id  = 2
    ip  = "10.0.13.11"
    mac = "52:54:00:00:00:11"
  },
  {
    id  = 3
    ip  = "10.0.13.12"
    mac = "52:54:00:00:00:12"
  }
]


#======================================================================================
# Worker nodes
#======================================================================================

# The default number of vCPU allocated to the worker VM. #
cluster_nodes_worker_default_cpu = 4

# The default amount of RAM allocated to the worker VM [GiB]. #
cluster_nodes_worker_default_ram = 8

# The default amount of disk allocated to the worker VM [GiB]. #
cluster_nodes_worker_default_storage = 32

# Sets worker node's role label. #
# Note: Leave empty ("") to not set the label. #
cluster_nodes_worker_default_label = "node"

# Worker nodes configuration. #
cluster_nodes_worker_instances = [
  {
    id  = 1
    ip  = "10.0.13.100"
    mac = "52:54:00:00:00:40"
  },
  {
    # Example of optional MAC address
    id  = 2
    ip  = "10.0.13.101"
    mac = null
  },
  {
    # Example of optional IP and MAC addresses
    id  = 3
    ip  = null
    mac = null
  }
]


#======================================================================================
# Kubernetes and Kubespray configuration
#======================================================================================

# The Kubernetes version that will be deployed. #
kubernetes_version = "v1.21.6"

# The overlay network plugin used by the Kubernetes cluster (flannel/weave/calico/cilium/canal/kube-router). #
kubernetes_networkPlugin = "calico"

# The DNS service used by Kubernetes cluster (coredns/kubedns). #
kubernetes_dnsMode = "coredns"

# The Git repository to clone Kubespray from #
kubernetes_kubespray_url = "https://github.com/kubernetes-sigs/kubespray.git"

# The version of Kubespray that will be used to deploy Kubernetes #
kubernetes_kubespray_version = "v2.17.1"

# Enable Kubespray addons. #
kubernetes_kubespray_addons_enabled = false

# Path to Kubespray addons configuration file. #
kubernetes_kubespray_addons_configPath = "defaults/addons.yml"

# Copies config file to ~/.kube directory. #
# Note: Kubeconfig will be always available in config/admin.conf after installation. #
kubernetes_other_copyKubeconfig = false
