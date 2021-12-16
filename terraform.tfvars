#======================================================================================
# General configuration
#======================================================================================

# Provider's URI #
libvirt_provider_uri = "qemu:///system"

# Resource pool name #
libvirt_resource_pool_name = "k8s-resource-pool"

# Location where resource pool will be initialized #
libvirt_resource_pool_location = "/var/lib/libvirt/pools/"


#======================================================================================
# Global VM configuration
#======================================================================================

# Username used to SSH to the VM #
vm_user = "k8s"

# Private SSH key location (for VMs) (example: ~/.ssh/id_rsa) #
vm_ssh_private_key = "~/.ssh/id_rsa"

# Add VMs to SSH known hosts #
vm_ssh_known_hosts = true

# Linux distribution that will be used on VMs (ubuntu, debian, centos, "") #
vm_distro = ""

# Source of linux image. It can be path to an image on host's filesystem or an URL #
vm_image_source = ""

# The prefix added to names of VMs #
vm_name_prefix = "k8s"

# Network interface used by VMs to connect to the network #
vm_network_interface = "ens3"

# If true, system will be updated and upgraded #
vm_update = false

#======================================================================================
# Network configuration
#======================================================================================

# Network name #
network_name = "k8s-network"

# Network mode (nat, route, bridge) #
network_mode = "nat"

# Network CIDR (example: 192.168.113.0/24) #
network_cidr = "192.168.113.0/24"

# Network (virtual) bridge #
# Note: For network mode 'bridge', bridge on host needs to preconfigured (example: br0) #
network_bridge = "virbr1"

# Network gateway (example: 192.168.113.1) #
# Note: If not provided, it will be calculated as first host in network CIDR. #
#       +-> first host of 192.168.113.0/24 is 192.168.113.1 #
#network_gateway = "192.168.113.1"

# Network DNS list (if empty, network gateway is set as a DNS) #
network_dns_list = [
  "1.1.1.1",
  "1.0.0.1"
]

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

# The default number of vCPU allocated to the load balancer VM #
lb_default_cpu = 1

# The default amount of RAM allocated to the load balancer VM [GiB] #
lb_default_ram = 2

# The default amount of disk allocated to the load balancer VM [GiB] #
lb_default_storage = 16

# HAProxy internal load balancer (iLB) nodes configuration #
lb_nodes = [
  {
    id  = 1
    ip  = "192.168.113.5"
    mac = "52:54:00:00:00:05"
  },
  {
    id  = 2
    ip  = "192.168.113.6"
    mac = "52:54:00:00:00:06"
  }
]

# Virtual/Floating IP address. #
# Note: Floating IP only applies if at least one load balancer is defined, #
# otherwise IP of the first master node will be used as control plane endpoint. #
lb_vip = "192.168.113.200"


#======================================================================================
# Master node VMs parameters
#======================================================================================

# The default number of vCPU allocated to the master VM #
master_default_cpu = 2

# The default amount of RAM allocated to the master VM [GiB]] #
master_default_ram = 2

# The default amount of disk allocated to the master VM [GiB] #
master_default_storage = 16

# Master nodes configuration #
# Note that number of masters cannot be divisible by 2. #
master_nodes = [
  {
    id  = 1
    ip  = "192.168.113.10"
    mac = "52:54:00:00:00:10"
  },
  {
    id  = 2
    ip  = "192.168.113.11"
    mac = "52:54:00:00:00:11"
  },
  {
    id  = 3
    ip  = "192.168.113.12"
    mac = "52:54:00:00:00:12"
  }
]


#======================================================================================
# Worker node VMs parameters
#======================================================================================

# The default number of vCPU allocated to the worker VM #
worker_default_cpu = 4

# The default amount of RAM allocated to the worker VM [GiB] #
worker_default_ram = 8

# The default amount of disk allocated to the worker VM [GiB] #
worker_default_storage = 32

# Sets worker node's role label #
# Note: Leave empty ("") to not set the label. #
worker_node_label = "node"

# Worker nodes configuration #
worker_nodes = [
  {
    id  = 1
    ip  = "192.168.113.100"
    mac = "52:54:00:00:00:40"
  },
  {
    # Example of optional MAC address
    id  = 2
    ip  = "192.168.113.101"
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
# General Kubernetes configuration
#======================================================================================

# The Git repository to clone Kubespray from #
k8s_kubespray_url = "https://github.com/kubernetes-sigs/kubespray.git"

# The version of Kubespray that will be used to deploy Kubernetes #
k8s_kubespray_version = "v2.17.1"

# The Kubernetes version that will be deployed #
k8s_version = "v1.21.6"

# The overlay network plugin used by the Kubernetes cluster (flannel/weave/calico/cilium/canal/kube-router) #
k8s_network_plugin = "calico"

# The DNS service used by Kubernetes cluster (coredns/kubedns)#
k8s_dns_mode = "coredns"

# Copies config file to ~/.kube directory #
# Note: Kubeconfig will be always available in config/admin.conf after installation #
k8s_copy_kubeconfig = false


#======================================================================================
# Kubespray addons
#======================================================================================

#=========================
# Custom addons
#=========================

# IMPORTANT: If custom addons are enabled, variables from other sections below
# will be ignored and addons from file path provided will be applied instead.

# Use custom addons.yml #
kubespray_custom_addons_enabled = false

# Path to custom addons.yml #
kubespray_custom_addons_path = "defaults/addons.yml"

#=========================
# General
#=========================

# Install Kubernetes dashboard #
k8s_dashboard_enabled = false

# Creates Kubernets dashboard RBAC token (dashboard needs to be enabled) #
k8s_dashboard_rbac_enabled = false
k8s_dashboard_rbac_user    = "admin"

# Install helm #
helm_enabled = false

#=========================
# Local path provisioner
#=========================

# Note: This is dynamic storage provisioner #

# Install Rancher's local path provisioner #
local_path_provisioner_enabled = false

# Version #
local_path_provisioner_version = "v0.0.19"

# Namespace in which provisioner will be installed #
local_path_provisioner_namespace = "local-path-provisioner"

# Storage class #
local_path_provisioner_storage_class = "local-storage"

# Reclaim policy (Delete/Retain) #
local_path_provisioner_reclaim_policy = "Delete"

# Claim root #
local_path_provisioner_claim_root = "/opt/local-path-provisioner/"

#=========================
# MetalLB
#=========================

# Install MetalLB #
metallb_enabled = false

# MetalLB version #
metallb_version = "v0.9.5"

# Kubernetes limits (1000m = 1 vCore) #
metallb_cpu_limit = "500m"
metallb_mem_limit = "500Mi"
metallb_port      = 7472

# MetalLB protocol (layer2/bgp) #
metallb_protocol = "layer2"

# IP range for services of type LoadBalancer #
metallb_ip_range = "192.168.113.241-192.168.113.254"

# MetalLB peers #
# Note: This variable will be applied only in 'bgp' mode #
metallb_peers = [
  {
    peer_ip  = "192.168.113.1"
    peer_asn = 65000
    my_asn   = 65000
  }
]
