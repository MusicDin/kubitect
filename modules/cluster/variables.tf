variable "action" {
  type        = string
  description = "Action that needs to be done on cluster"
}

#======================================================================================
# Virtual machine variables
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

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
}

variable "vm_network_interface" {
  type        = string
  description = "VM network interface used for Keepalived"
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
  description = "The DNS service used by Kubernetes cluster"
}

#======================================================================================
# Other
#======================================================================================

variable "k8s_copy_kubeconfig" {
  type        = string
  description = "If enabled, kubeconfig (config/admin.conf) will be copied to ~/.kube directory"
}

#======================================================================================
# Kubespray addons
#======================================================================================

variable "kubespray_custom_addons_enabled" {
  type        = bool
  description = "If enabled, custom addons.yml will be used"
}

variable "kubespray_custom_addons_path" {
  type        = string
  description = "If custom addons are enabled, addons YAML file from this path will be used"
}

variable "k8s_dashboard_enabled" {
  type        = bool
  description = "Sets up Kubernetes dashboard if enabled"
}

variable "k8s_dashboard_rbac_enabled" {
  type        = bool
  description = "If enabled, Kubernetes dashboard service account will be created"
}

variable "k8s_dashboard_rbac_user" {
  type        = string
  description = "Kubernetes dashboard service account user"
}

variable "helm_enabled" {
  type        = bool
  description = "Sets up Helm if enabled"
}

variable "local_path_provisioner_enabled" {
  type        = bool
  description = "Sets up Rancher's local path provisioner if enabled"
}

variable "local_path_provisioner_version" {
  type        = string
  description = "Local path provisioner version"
}

variable "local_path_provisioner_namespace" {
  type        = string
  description = "Namespace in which local path provisioner will be installed"
}

variable "local_path_provisioner_storage_class" {
  type        = string
  description = "Local path provisioner storage class"
}

variable "local_path_provisioner_reclaim_policy" {
  type        = string
  description = "Local path provisioner reclaim policy"
}

variable "local_path_provisioner_claim_root" {
  type        = string
  description = "Local path provisioner claim root"
}

variable "metallb_enabled" {
  type        = bool
  description = "Sets up MetalLB if enabled"
}

variable "metallb_version" {
  type        = string
  description = "MetalLB version"
}

variable "metallb_port" {
  type        = number
  description = "Kubernetes MetalLB port"
}

variable "metallb_cpu_limit" {
  type        = string
  description = "MetalLB pod CPU limit"
}

variable "metallb_mem_limit" {
  type        = string
  description = "MetalLB pod memory (RAM) limit"
}

variable "metallb_protocol" {
  type        = string
  description = "MetalLB protocol"
}

variable "metallb_ip_range" {
  type        = string
  description = "IP range that MetalLB will use for services of type LoadBalancer"
}

variable "metallb_peers" {
  type = list(object({
    peer_ip  = string
    peer_asn = number
    my_asn   = number
  }))
  description = "List of MetalLB peers"
}
