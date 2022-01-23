variable "action" {
  type        = string
  description = "Action that needs to be done on cluster"
}

#======================================================================================
# Virtual machine configuration
#======================================================================================

variable "lb_vip" {
  type        = string
  description = "Load balancer virtual IP address (VIP)"
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

variable "vm_network_interface" {
  type        = string
  description = "VM network interface used for Keepalived"
}

#======================================================================================
# Virtual machine instances
#======================================================================================

variable "worker_nodes" {
  type = list(object({
    name = string
    ip   = string
  }))
  description = "Worker nodes info"
}

variable "worker_node_label" {
  type        = string
  description = "Worker node role label"
}

variable "master_nodes" {
  type = list(object({
    name = string
    ip   = string
  }))
  description = "Master nodes info"
}

variable "lb_nodes" {
  type = list(object({
    id   = number
    name = string
    ip   = string
  }))
  description = "Load balancers info"
}

#======================================================================================
# Kubernetes and Kubespray variables
#======================================================================================

variable "kubernetes_version" {
  type        = string
  description = "The version of Kuberenetes that will be deployed"
}

variable "kubernetes_networkPlugin" {
  type        = string
  description = "The overlay network plugin used by Kubernetes cluster"
}

variable "kubernetes_dnsMode" {
  type        = string
  description = "The DNS service used by Kubernetes cluster"
}

variable "kubernetes_kubespray_url" {
  type        = string
  description = "The Git repository URL to clone Kubespray from"
}

variable "kubernetes_kubespray_version" {
  type        = string
  description = "The version of Kubespray that will be used to deploy Kubernetes"
}

variable "kubernetes_kubespray_addons_enabled" {
  type        = bool
  description = "If enabled, configured Kubespray addons will be applied."
  default     = false
}

variable "kubernetes_kubespray_addons_configPath" {
  type        = string
  description = "If Kubespray addons are enabled, addons configuration file on this path will be used."
  default     = ""
}

variable "kubernetes_other_copyKubeconfig" {
  type        = string
  description = "If enabled, kubeconfig (config/admin.conf) will be copied to ~/.kube directory"
}

#======================================================================================
# Kubernetes dashboard
#======================================================================================

#variable "k8s_dashboard_enabled" {
#  type        = bool
#  description = "Sets up Kubernetes dashboard if enabled"
#}

#variable "k8s_dashboard_rbac_enabled" {
#  type        = bool
#  description = "If enabled, Kubernetes dashboard service account will be created"
#}

#variable "k8s_dashboard_rbac_user" {
#  type        = string
#  description = "Kubernetes dashboard service account user"
#}
