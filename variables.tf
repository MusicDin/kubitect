#======================================================================================
# Libvirt connection
#======================================================================================

variable "action" {
  description = "Which action have to be done on the cluster (create, add_worker, remove_worker, or upgrade)"
  default     = "create"
}

variable "libvirt_resource_pool_name" {
  description = "The libvirt resource pool name"
}

#======================================================================================
# Kubernetes infrastructure
#======================================================================================

#============================#
# General variables           #
#============================#

variable "vm_user" {
  description = "SSH user for VMs"
  default     = "user"
}

variable "vm_ssh_private_key" {
  description = "Location of private ssh key for VMs"
}

#variable "vm_privilege_password" {
#  description = "Sudo or su password for VMs privilege escalation (Don't set this in .tfvars as plain text!)"
#}

variable "vm_distro" {
  description = "Linux distribution used by VMs (currently not in use)"
  default     = "ubuntu"
}

variable "vm_image_name" {
  description = "Name of VM image (it has to be in downloads folder)"
}

variable "vm_name_prefix" {
  description = "Prefix added to names of VMs"
  default     = "vm"
}

variable "vm_network_name" {
  description = "Network used by VMs"
}

variable "vm_network_netmask" {
  description = "The netmask used to configure the network cards of the VMs"
}

variable "vm_network_gateway" {
  description = "The network gateway used by VMs"
}

variable "vm_domain" {
  description = "Domain name used by VMs"
}

# Should be set for lb, master and workers separatly! #
variable "vm_disk_size" {
  description = "Disk size in bytes (default: 15GB)"
  default = 16106127360
}

#============================#
# Load balancer variables    #
#============================#

variable "vm_lb_cpu" {
  description = "The number of vCPU allocated to the HAProxy load balancer"
  default     = "1"
}

variable "vm_lb_ram" {
  description = "The amount of RAM allocated to the HAProxy load balancer"
  default     = "4096"
}

variable "vm_lb_macs" {
  type = map(string)
  description = "The MAC addresses of HAProxy load balancer nodes"
}

variable "vm_lb_ips" {
  type = map(string)
  description = "The IP addresses of HAProxy load balancer nodes"
}

variable "vm_haproxy_vip" {
  description = "The IP address of the load balancer floating VIP"
}

#============================#
# Master nodes variables     #
#============================#

variable "vm_master_cpu" {
  description = "The number of vCPU allocated to the master node"
  default     = "1"
}

variable "vm_master_ram" {
  description = "The amount of RAM allocated to the master node"
  default     = "4096"
}

variable "vm_master_macs" {
  type        = map(string)
  description = "The MAC addresses of master nodes"
}


variable "vm_master_ips" {
  type        = map(string)
  description = "The IP addresses of master nodes"
}

#============================#
# Worker nodes variables     #
#============================#

variable "vm_worker_cpu" {
  description = "The number of vCPU allocated to the worker node"
  default     = "2"
}

variable "vm_worker_ram" {
  description = "The amount of RAM allocated to the worker node"
  default     = "8192"
}

variable "vm_worker_macs" {
  type        = map(string)
  description = "The MAC addresses of worker nodes"
}

variable "vm_worker_ips" {
  type        = map(string)
  description = "The IP addresses of worker nodes"
}

#======================================================================================
# General kubernetes (k8s) variables
#======================================================================================

variable "k8s_kubespray_url" {
  description = "The Git repository URL to clone Kubespray from"
}

variable "k8s_kubespray_version" {
  description = "The version of Kubespray that will be used to deploy Kubernetes"
}

variable "k8s_version" {
  description = "The version of Kuberenetes that will be deployed"
}

variable "k8s_network_plugin" {
  description = "The overlay network plugin used by Kubernetes cluster"
}

variable "k8s_dns_mode" {
  description = "The DNS service used by Kubernetes cluster (coredns/kubedns)"
}
