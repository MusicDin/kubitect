#======================================================================================
# Libvirt connection
#======================================================================================

variable "action" {
  type        = string
  description = "Which action have to be done on the cluster (create, add_worker, remove_worker, or upgrade)"
  default     = "create"
}

variable "libvirt_resource_pool_name" {
  type        = string
  description = "The libvirt resource pool name"
}

#======================================================================================
# Network
#======================================================================================

variable "network_name" {
  type        = string
  description = "Network name"
  default     = "k8s-network"
}

variable "network_interface" {
  type        = string
  description = "Network interface used for VMs (cloud-init) and Keepalived"
  default     = "ens3"
}

variable "network_mac" {
  type        = string
  description = "Network MAC address"
  default     = "52:54:00:4f:e3:88"
}

variable "network_ip" {
  type        = string
  description = "Network gateway IP address"
  default     = "192.168.113.1"
}

variable "network_mask" {
  type        = string
  description = "Network mask"
  default     = "255.255.255.0"
}

variable "network_mask_bits" {
  type        = number
  description = "Bits used for network"
  default     = 24

  validation {
    condition     = var.network_mask_bits > 0 && var.network_mask_bits <= 32
    error_message = "Valid value for bits used for network is between 1 and 32. (Default value is 24)."
  }
}

variable "network_nat_port_start" {
  type        = number
  description = "NAT (Network Address Translation) port start (from port)"
  default     = 1024

  validation {
    condition     = var.network_nat_port_start > 0 && var.network_nat_port_start <= 65535
    error_message = "Valid NAT port value is between 1 and 65535."
  }
}

variable "network_nat_port_end" {
  type        = number
  description = "NAT port end (to port)"
  default     = 65535

  validation {
    condition     = var.network_nat_port_end > 0 && var.network_nat_port_end <= 65535
    error_message = "Valid NAT port value is between 1 and 65535."
  }

}

variable "network_dhcp_ip_start" {
  type        = string
  description = "DHCP IP range start"
  default     = "192.168.113.2"
}

variable "network_dhcp_ip_end" {

  type        = string
  description = "DHCP IP range end"
  default     = "192.168.113.254"
}

#======================================================================================
# Kubernetes infrastructure
#======================================================================================

#============================#
# General variables           #
#============================#

variable "vm_user" {
  type        = string
  description = "SSH user for VMs"
  default     = "user"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private ssh key for VMs"
}

variable "vm_distro" {
  type        = string
  description = "Linux distribution used on VMs. Possible values: [ubuntu, centos, debian]"
  default     = "N/A"
}

variable "vm_image_source" {
  type        = string
  description = "Image source, which can be path on host's filesystem or URL."
}

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
  default     = "vm"
}

#============================#
# Load balancer variables    #
#============================#

variable "vm_lb_cpu" {
  type        = number
  description = "The number of vCPU allocated to the HAProxy load balancer"
  default     = 1
}

variable "vm_lb_ram" {
  type        = number
  description = "The amount of RAM allocated to the HAProxy load balancer"
  default     = 4096
}

variable "vm_lb_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the HAProxy load balancer. Default: 15GB"
  default     = 16106127360
}

variable "vm_lb_macs_ips" {
  type = map(string)
  description = "MAC (key) and IP (value) addresses of HAProxy load balancer nodes"
}

variable "vm_lb_vip" {
  type        = string
  description = "The IP address of HAProxy load balancer floating VIP"
}

#============================#
# Master nodes variables     #
#============================#

variable "vm_master_cpu" {
  type        = number
  description = "The number of vCPU allocated to the master node"
  default     = 1
}

variable "vm_master_ram" {
  type        = number
  description = "The amount of RAM allocated to the master node"
  default     = 4096
}

variable "vm_master_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the master node. Default: 15GB"
  default     = 16106127360
}

variable "vm_master_macs_ips" {
  type        = map(string)
  description = "MAC and IP addresses of master nodes"
}

#============================#
# Worker nodes variables     #
#============================#

variable "vm_worker_cpu" {
  type        = number
  description = "The number of vCPU allocated to the worker node"
  default     = 2
}

variable "vm_worker_ram" {
  type        = number
  description = "The amount of RAM allocated to the worker node"
  default     = 8192
}
variable "vm_worker_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the worker node. Default: 30GB"
  default     = 32212254720
}
variable "vm_worker_macs_ips" {
  type        = map(string)
  description = "MAC and IP addresses of worker nodes"
}

#======================================================================================
# General kubernetes (k8s) variables
#======================================================================================

variable "k8s_kubespray_url" {
  type        = string
  description = "The Git repository URL to clone Kubespray from"
}

variable "k8s_kubespray_version" {
  type        = string
  description = "The version of Kubespray that will be used to deploy Kubernetes"
}

variable "k8s_version" {
  type        = string
  description = "The version of Kuberenetes that will be deployed"
}

variable "k8s_network_plugin" {
  type        = string
  description = "The overlay network plugin used by Kubernetes cluster"
}

variable "k8s_dns_mode" {
  type        = string
  description = "The DNS service used by Kubernetes cluster (coredns/kubedns)"
}
