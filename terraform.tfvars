#======================================================================================
# Libvirt provider configuration
#======================================================================================

# Resource pool name #
libvirt_resource_pool_name = "k8s-resource-pool"

#======================================================================================
# Global virtual machines parameters
#======================================================================================

# Username used to SSH to the VM #
vm_user = "k8s"

# Private SSH key (for VMs) location (example: ~/.ssh/id_rsa) #
vm_ssh_private_key = "~/.ssh/id_rsa"

# Name of linux image. Image has to be in ./downloads/ folder. (example: ubuntu.img) #
vm_image_name = "bionic-server-cloudimg-amd64.img"

# The prefix added to names of VMs #
vm_name_prefix = "k8s"

# Network KVM name used VMs #
vm_network_name = "k8s-network"

# The netmask used to configure the network cards of the VMs (example: 24) #
vm_network_netmask = "24"

# The network gateway used by VMs #
vm_network_gateway = "192.168.113.1"

# The domain name used by VMs #7
vm_domain = "din-cloud.com"

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

# The number of vCPU allocated to the load balancer VM #
vm_lb_cpu = "1"

# The amount of RAM allocated to the load balancer VM #
vm_lb_ram = "4096"

# The MAC addresses for load balancer VMs. #
vm_lb_macs = {
  "0" = "52:54:00:00:00:05"
}

# The IP addresses for load balancer VMs. #
vm_lb_ips = {
  "0" = "192.168.113.5"
}

# The IP address of the load balancer floating VIP #
vm_haproxy_vip = "192.168.113.5"


#======================================================================================
# Master node VMs parameters
#======================================================================================

# The number of vCPU allocated to the master VM #
vm_master_cpu = "1"

# The amount of RAM allocated to the master VM #
vm_master_ram = "8192"

# The MAC addresses for master VMs. There MUST be exactly 3 inputs! #
vm_master_macs = {
  "0" = "52:54:00:00:00:10",
  "1" = "52:54:00:00:00:11",
  "2" = "52:54:00:00:00:12"
}

# The IP addresses for master VMs. There MUST be exactly 3 inputs! #
vm_master_ips = {
  "0" = "192.168.113.10"
  "1" = "192.168.113.11"
  "2" = "192.168.113.12"
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

# The number of vCPU allocated to the worker VM #
vm_worker_cpu = "2"

# The amount of RAM allocated to the worker VM #
vm_worker_ram = "8192"

# The MAC addresses for worker VMs. #
vm_worker_macs = {
  "0" = "52:54:00:00:00:40"
  "1" = "52:54:00:00:00:41"
  "2" = "52:54:00:00:00:42"
}

# The IP addresses for worker VMs. #
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
