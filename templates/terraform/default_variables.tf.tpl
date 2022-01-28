#
# Note: All variables below have no desription and 'null' for default value.
#       Please see ./module/main/variables.tf for default values and descriptions. 
#       Reason behind that is that if YAML config is used all this file is mostly redundant.
#

#======================================================================================
# General configuration
#======================================================================================

variable "action" {
  type    = string
  default = null
}

variable "libvirt_provider_uri" {
  type    = string
  default = null
}

variable "libvirt_resource_pool_location" {
  type    = string
  default = null
}

#======================================================================================
# Cluster infrastructure configuration
#======================================================================================

variable "cluster_name" {
  type    = string
  default = null
}

#================================
# Node template
#================================

variable "cluster_nodeTemplate_user" {
  type    = string
  default = null
}

variable "cluster_nodeTemplate_ssh_privateKeyPath" {
  type    = string
  default = null
}

variable "cluster_nodeTemplate_ssh_addToKnownHosts" {
  type    = bool
  default = true
}

variable "cluster_nodeTemplate_image_distro" {
  type    = string
  default = null
}

variable "cluster_nodeTemplate_image_source" {
  type    = string
  default = null
}

variable "cluster_nodeTemplate_networkInterface" {
  type    = string
  default = null
}

variable "cluster_nodeTemplate_updateOnBoot" {
  type    = bool
  default = null
}

#================================
# Cluster network
#================================

variable "cluster_network_mode" {
  type    = string
  default = null
}

variable "cluster_network_cidr" {
  type    = string
  default = null
}

variable "cluster_network_gateway" {
  type    = string
  default = null
}

variable "cluster_network_bridge" {
  type    = string
  default = null
}

variable "cluster_network_dns" {
  type    = list(string)
  default = null
}

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

variable "cluster_nodes_loadBalancer_vip" {
  type    = string
  default = null
}

variable "cluster_nodes_loadBalancer_default_cpu" {
  type    = number
  default = null
}

variable "cluster_nodes_loadBalancer_default_ram" {
  type    = number
  default = null
}

variable "cluster_nodes_loadBalancer_default_storage" {
  type    = number
  default = null
}

variable "cluster_nodes_loadBalancer_instances" {
  type = list(object({
    id     = number
    mac    = string
    ip     = string
    server = optional(string)
  }))
  default = null
}

#======================================================================================
# Master node VMs parameters
#======================================================================================

variable "cluster_nodes_master_default_cpu" {
  type    = number
  default = null
}

variable "cluster_nodes_master_default_ram" {
  type    = number
  default = null
}

variable "cluster_nodes_master_default_storage" {
  type    = number
  default = null
}

variable "cluster_nodes_master_instances" {
  type = list(object({
    id     = number
    mac    = string
    ip     = string
    server = optional(string)
  }))
  default = null
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

variable "cluster_nodes_worker_default_cpu" {
  type    = number
  default = null
}

variable "cluster_nodes_worker_default_ram" {
  type    = number
  default = null
}

variable "cluster_nodes_worker_default_storage" {
  type    = number
  default = null
}

variable "cluster_nodes_worker_default_label" {
  type    = string
  default = null
}

variable "cluster_nodes_worker_instances" {
  type = list(object({
    id     = number
    mac    = string
    ip     = string
    server = optional(string)
  }))
  default = null
}

#======================================================================================
# General Kubernetes configuration
#======================================================================================

variable "kubernetes_version" {
  type    = string
  default = null
}

variable "kubernetes_networkPlugin" {
  type    = string
  default = null
}

variable "kubernetes_dnsMode" {
  type    = string
  default = null
}

variable "kubernetes_kubespray_url" {
  type    = string
  default = null
}

variable "kubernetes_kubespray_version" {
  type    = string
  default = null
}

variable "kubernetes_kubespray_addons_enabled" {
  type    = bool
  default = false
}

variable "kubernetes_kubespray_addons_configPath" {
  type    = string
  default = null
}

variable "kubernetes_other_copyKubeconfig" {
  type    = bool
  default = null
}
