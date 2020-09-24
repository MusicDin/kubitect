#======================================================================================
# Libvirt provider configuration
#======================================================================================

# Resource pool name #
libvirt_resource_pool_name = "k8s-resource-pool"

# Location where resource pool will be initialized (Path must contain "/" at the end) #
libvirt_resource_pool_location = "/var/lib/libvirt/pools/"

#======================================================================================
# Global VM parameters
#======================================================================================

# Username used to SSH to the VM #
vm_user = "k8s"

# Private SSH key (for VMs) location (example: ~/.ssh/id_rsa) #
vm_ssh_private_key = "~/.ssh/id_rsa"

# Linux distribution that will be used on VMs. Possible values are: [ubuntu, debian, centos] #
vm_distro = ""

# Source of linux image. It can be path to an image on host's filesystem or an URL #
vm_image_source = ""

# The prefix added to names of VMs #
vm_name_prefix = "k8s"

#======================================================================================
# Global network parameters
#======================================================================================

# Network KVM name used VMs #
network_name = "k8s-network"

# Network interface used for VMs (cloud-init) and Keepalived #
network_interface = "ens3"

# The network gateway used by VMs #
network_gateway = "192.168.113.1"

# Network mask #
network_mask = "255.255.255.0"

# Bits used for network mask (example: 255.255.255.0 uses 24 bits for network) #
network_mask_bits = 24


# --- NAT port range --- #

# NAT (Network Address Translation) port start (from port) #
network_nat_port_start = 1024

# NAT port end (to port) #
network_nat_port_end = 65535


# --- DHCP IP range --- #
# Based on VM's MAC and IP address, static IPs are configured in network, but all IPs 
# have to be within DHCP IP range.

# DHCP IP start (from IP) #
network_dhcp_ip_start = "192.168.113.2"

# DHCP IP end (to IP) #
network_dhcp_ip_end = "192.168.113.254"

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

# The number of vCPU allocated to the load balancer VM #
vm_lb_cpu = "1"

# The amount of RAM (in Megabytes) allocated to the load balancer VM #
vm_lb_ram = 2048

# The amount of disk (in Bytes) allocated to the load balancer VM #
vm_lb_storage = 16106127360

# The MAC addresses for load balancer VMs #
vm_lb_macs = {
  "0" = "52:54:00:00:00:05"
  "1" = "52:54:00:00:00:06"
}

# The IP addresses for load balancer VMs #
vm_lb_ips = {
  "0" = "192.168.113.5"
  "1" = "192.168.113.6"
}

# The floating IP address for load balancers #
vm_lb_vip = "192.168.113.200"


#======================================================================================
# Master node VMs parameters
#======================================================================================

# The number of vCPU allocated to the master VM #
vm_master_cpu = 1

# The amount of RAM (in Megabytes) allocated to the master VM #
vm_master_ram = 2048

# The amount of disk (in Bytes) allocated to the master VM #
vm_master_storage = 16106127360

# NOTE: Number of master nodes cannot be divisible by 2 #

# The MAC addresses for master VMs #
vm_master_macs = {
  "0" = "52:54:00:00:00:10",
  "1" = "52:54:00:00:00:11",
  "2" = "52:54:00:00:00:12"
}

# The IP addresses for master VMs #
vm_master_ips = {
  "0" = "192.168.113.10"
  "1" = "192.168.113.11"
  "2" = "192.168.113.12"
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

# The number of vCPU allocated to the worker VM #
vm_worker_cpu = 2

# The amount of RAM (in Megabytes) allocated to the worker VM #
vm_worker_ram = 8192

# The amount of disk (in Bytes) allocated to the worker VM #
vm_worker_storage = 16106127360

# The MAC addresses for worker VMs #
vm_worker_macs = {
  "0" = "52:54:00:00:00:40"
  "1" = "52:54:00:00:00:41"
  "2" = "52:54:00:00:00:42"
}

# The IP addresses for worker VMs #
vm_worker_ips = {
  "0" = "192.168.113.100"
  "1" = "192.168.113.101"
  "2" = "192.168.113.102"
}

#======================================================================================
# Kubernetes (k8s) parameters
#======================================================================================

# The Git repository to clone Kubespray from #
k8s_kubespray_url = "https://github.com/kubernetes-sigs/kubespray.git"

# The version of Kubespray that will be used to deploy Kubernetes #
k8s_kubespray_version = "v2.13.0"

# The Kubernetes version that will be deployed #
k8s_version = "v1.17.5"

# The overlay network plugin used by the Kubernetes cluster #
k8s_network_plugin = "calico"

# The DNS service used by Kubernetes cluster (coredns/kubedns)#
k8s_dns_mode = "coredns"
