variable "action" {
  type = string
  description = "Action that needs to be done on cluster"
}

#======================================================================================
# Virtual machine variables 
#======================================================================================

variable "vm_worker_ips" {
  type        = list
  description = "IP addresses of worker nodes"
}

variable "vm_master_ips" {
  type        = list
  description = "IP addresses of master nodes"
}

variable "vm_lb_ips" {
  type        = list
  description = "IP addresses of load balancer VMs"
}

variable "vm_lb_vip" {
  type        = string
  description = "Floating virtual IP of the load balancer"
}

variable "vm_distro" {
  type        = string
  description = "Linux distribution used for VMs"
}

variable "vm_user" {
  type        = string
  description = "SSH user for VMs"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private ssh key for VMs"
}

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
}

variable "network_interface" {
  type        = string
  description = "Network interface used for Keepalived"
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
