#======================================================================================
# General configuration
#======================================================================================

# Provider's URI #
libvirt_provider_uri = "qemu:///system"

# Resource pool name #
libvirt_resource_pool_name = "k8s-resource-pool"

# Location where resource pool will be initialized (Path must contain "/" at the end) #
libvirt_resource_pool_location = "/var/lib/libvirt/pools/"


#======================================================================================
# Global VM configuration
#======================================================================================

# Username used to SSH to the VM #
vm_user = "k8s"

# Private SSH key location (for VMs) (example: ~/.ssh/id_rsa) #
vm_ssh_private_key = "~/.ssh/id_rsa"

# Add VMs to SSH known hosts #
vm_ssh_known_hosts = "true"

# Linux distribution that will be used on VMs (ubuntu, debian, centos, "") #
vm_distro = ""

# Source of linux image. It can be path to an image on host's filesystem or an URL #
vm_image_source = ""

# The prefix added to names of VMs #
vm_name_prefix = "k8s"

# Network interface used by VMs to connect to the network #
vm_network_interface = "ens3"


#======================================================================================
# Network configuration
#======================================================================================

# Network name #
network_name = "k8s-network"

# Network mode (nat, route) #
network_mode = "nat"

# Network (virtual) bridge #
network_bridge = "virbr1"

# Network MAC address #
network_mac = "52:54:00:4f:e3:88"

# Network gateway IP address #
network_gateway = "192.168.113.1"

# Bits used for network mask (example: 255.255.255.0 uses 24 bits for network) #
network_mask_bits = 24

# --- DHCP IP range --- #
# DHCP is used as network management protocol.
# VM's IP address is configured as static IP based on it's MAC.
# IP addresses have to be within DHCP IP range, otherwise IP won't be assigned to VM.

# DHCP IP start (from IP) #
network_dhcp_ip_start = "192.168.113.2"

# DHCP IP end (to IP) #
network_dhcp_ip_end = "192.168.113.254"


#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

# The number of vCPU allocated to the load balancer VM #
vm_lb_cpu = 1

# The amount of RAM allocated to the load balancer VM (in Megabytes - MB) #
vm_lb_ram = 2048

# The amount of disk allocated to the load balancer VM (in Bytes - B) #
vm_lb_storage = 16106127360

# MAC and IP addresses for load balancer VMs. #
vm_lb_macs_ips = {
  "52:54:00:00:00:05" = "192.168.113.5"
  "52:54:00:00:00:06" = "192.168.113.6"
}

# Floating IP address. #
# Note: Floating IP only applies if at least one load balancer is defined, #
# otherwise IP of the first master node will be used. #
vm_lb_vip = "192.168.113.200"


#======================================================================================
# Master node VMs parameters
#======================================================================================

# The number of vCPU allocated to the master VM #
vm_master_cpu = 2

# The amount of RAM allocated to the master VM (in Megabytes - MB) #
vm_master_ram = 2048

# The amount of disk allocated to the master VM (in Bytes - B) #
vm_master_storage = 16106127360

# MAC and IP addresses for master VMs. It is recommended to have at least 3 masters. #
# Also note that number of masters cannot be divisible by 2. #
vm_master_macs_ips = {
  "52:54:00:00:00:10" = "192.168.113.10"
  "52:54:00:00:00:11" = "192.168.113.11"
  "52:54:00:00:00:12" = "192.168.113.12"
}


#======================================================================================
# Worker node VMs parameters
#======================================================================================

# The number of vCPU allocated to the worker VM #
vm_worker_cpu = 4

# The amount of RAM allocated to the worker VM (in Megabytes - MB) #
vm_worker_ram = 8192

# The amount of disk allocated to the worker VM (in Bytes - B) #
vm_worker_storage = 32212254720

# MAC and IP addresses for worker VMs. #
vm_worker_macs_ips = {
  "52:54:00:00:00:40" = "192.168.113.100"
  "52:54:00:00:00:41" = "192.168.113.101"
  "52:54:00:00:00:42" = "192.168.113.102"
}

# Sets worker node's role label #
# Note: Leave empty ("") to not set the label. #
vm_worker_node_label = "node"


#======================================================================================
# General Kubernetes configuration
#======================================================================================

# The Git repository to clone Kubespray from #
k8s_kubespray_url = "https://github.com/kubernetes-sigs/kubespray.git"

# The version of Kubespray that will be used to deploy Kubernetes #
k8s_kubespray_version = "v2.16.0"

# The Kubernetes version that will be deployed #
k8s_version = "v1.20.7"

# The overlay network plugin used by the Kubernetes cluster (flannel/weave/calico/cilium/canal/kube-router) #
k8s_network_plugin = "calico"

# The DNS service used by Kubernetes cluster (coredns/kubedns)#
k8s_dns_mode = "coredns"

# Copies config file to ~/.kube directory #
# Note: Kubeconfig will be always available in config/admin.conf after installation #
k8s_copy_kubeconfig = "false"


#======================================================================================
# Kubespray addons
#======================================================================================

#=========================
# Custom addons
#=========================

# IMPORTANT: If custom addons are enabled, variables from other sections below
# will be ignored and addons from file path provided will be applied instead.

# Use custom addons.yml #
kubespray_custom_addons_enabled = "false"

# Path to custom addons.yml #
kubespray_custom_addons_path = "defaults/addons.yml"

#=========================
# General
#=========================

# Install Kubernetes dashboard #
k8s_dashboard_enabled = "false"

# Creates Kubernets dashboard RBAC token (dashboard needs to be enabled) #
k8s_dashboard_rbac_enabled = "false"
k8s_dashboard_rbac_user    = "admin"

# Install helm #
helm_enabled = "false"

#=========================
# Local path provisioner
#=========================

# Note: This is dynamic storage provisioner #

# Install Rancher's local path provisioner #
local_path_provisioner_enabled = "false"

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
metallb_enabled  = "false"

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
metallb_peers = [{
  peer_ip  = "192.168.113.1"
  peer_asn = 65000
  my_asn   = 65000
}]
