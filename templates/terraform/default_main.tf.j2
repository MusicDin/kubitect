#
# Configuration generated on: {{ now() }}
#

#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  internal = {
    is_bridge = (var.cluster_network_mode == "bridge")
    vm_types = {
      load_balancer = "lb"
      master        = "master"
      worker        = "worker"
    }
  }
}


#=====================================================================================
# Providers
#=====================================================================================
{# 
  Provider uri (localhost): "qemu:///system"
  Provider uri (remote): "qemu+ssh://<USER>@<IP>:<PORT>/system?keyfile=<PK_PATH>"
#}
provider "libvirt" {
  uri = var.libvirt_provider_uri
}

#======================================================================================
# Modules
#======================================================================================

module "main" {

  source = "./modules/main"

  # General
  action = var.action

  # Libvirt
  libvirt_resource_pool_location = var.libvirt_resource_pool_location

  # Global VM configuration
  cluster_name                             = var.cluster_name
  cluster_nodeTemplate_user                = var.cluster_nodeTemplate_user
  cluster_nodeTemplate_ssh_privateKeyPath  = var.cluster_nodeTemplate_ssh_privateKeyPath
  cluster_nodeTemplate_ssh_addToKnownHosts = var.cluster_nodeTemplate_ssh_addToKnownHosts
  cluster_nodeTemplate_image_distro        = var.cluster_nodeTemplate_image_distro
  cluster_nodeTemplate_image_source        = var.cluster_nodeTemplate_image_source
  cluster_nodeTemplate_networkInterface    = var.cluster_nodeTemplate_networkInterface
  cluster_nodeTemplate_updateOnBoot        = var.cluster_nodeTemplate_updateOnBoot

  # Network configuration
  cluster_network_mode    = var.cluster_network_mode
  cluster_network_cidr    = var.cluster_network_cidr
  cluster_network_gateway = var.cluster_network_gateway
  cluster_network_bridge  = var.cluster_network_bridge
  cluster_network_dns     = var.cluster_network_dns

  # HAProxy load balancer VMs parameters
  cluster_nodes_loadBalancer_vip             = var.cluster_nodes_loadBalancer_vip
  cluster_nodes_loadBalancer_default_cpu     = var.cluster_nodes_loadBalancer_default_cpu
  cluster_nodes_loadBalancer_default_ram     = var.cluster_nodes_loadBalancer_default_ram
  cluster_nodes_loadBalancer_default_storage = var.cluster_nodes_loadBalancer_default_storage
  cluster_nodes_loadBalancer_instances       = var.cluster_nodes_loadBalancer_instances

  # Master node VMs parameters
  cluster_nodes_master_default_cpu     = var.cluster_nodes_master_default_cpu
  cluster_nodes_master_default_ram     = var.cluster_nodes_master_default_ram
  cluster_nodes_master_default_storage = var.cluster_nodes_master_default_storage
  cluster_nodes_master_instances       = var.cluster_nodes_master_instances

  # Worker node VMs parameters
  cluster_nodes_worker_default_cpu     = var.cluster_nodes_worker_default_cpu
  cluster_nodes_worker_default_ram     = var.cluster_nodes_worker_default_ram
  cluster_nodes_worker_default_storage = var.cluster_nodes_worker_default_storage
  cluster_nodes_worker_default_label   = var.cluster_nodes_worker_default_label
  cluster_nodes_worker_instances       = var.cluster_nodes_worker_instances

  # Kubernetes & Kubespray
  kubernetes_version                     = var.kubernetes_version
  kubernetes_networkPlugin               = var.kubernetes_networkPlugin
  kubernetes_dnsMode                     = var.kubernetes_dnsMode
  kubernetes_kubespray_url               = var.kubernetes_kubespray_url
  kubernetes_kubespray_version           = var.kubernetes_kubespray_version
  kubernetes_kubespray_addons_enabled    = false # var.kubernetes_kubespray_addons_enabled
  kubernetes_kubespray_addons_configPath = ""    # var.kubernetes_kubespray_addons_configPath
  kubernetes_other_copyKubeconfig        = var.kubernetes_other_copyKubeconfig

  # Other
  internal = local.internal

  providers = {
    libvirt = libvirt
  }

}


#================================
# Cluster
#================================

# Configures k8s cluster using Kubespray #
module "k8s_cluster" {

  source = "./modules/cluster"

  action = var.action

  # VM variables #
  vm_user            = var.cluster_nodeTemplate_user
  vm_ssh_private_key = pathexpand(var.cluster_nodeTemplate_ssh_privateKeyPath)
  vm_distro          = var.cluster_nodeTemplate_image_distro
  vm_network_interface = (local.internal.is_bridge
    ? var.cluster_network_bridge
    : var.cluster_nodeTemplate_networkInterface
  )

  worker_node_label = var.cluster_nodes_worker_default_label
  lb_vip            = var.cluster_nodes_loadBalancer_vip
  lb_nodes = [
    for node in flatten(module.main) :
    node if node.type == local.internal.vm_types.load_balancer
  ]
  master_nodes = [
    for node in flatten(module.main) :
    node if node.type == local.internal.vm_types.master
  ]
  worker_nodes = [
    for node in flatten(module.main) :
    node if node.type == local.internal.vm_types.worker
  ]

  # K8s cluster variables #
  kubernetes_version                     = var.kubernetes_version
  kubernetes_networkPlugin               = var.kubernetes_networkPlugin
  kubernetes_dnsMode                     = var.kubernetes_dnsMode
  kubernetes_kubespray_url               = var.kubernetes_kubespray_url
  kubernetes_kubespray_version           = var.kubernetes_kubespray_version
  kubernetes_kubespray_addons_enabled    = false #var.kubernetes_kubespray_addons_enabled
  kubernetes_kubespray_addons_configPath = ""    #var.kubernetes_kubespray_addons_configPath
  kubernetes_other_copyKubeconfig        = var.kubernetes_other_copyKubeconfig

  # Other #
  #k8s_dashboard_rbac_enabled = var.k8s_dashboard_rbac_enabled
  #k8s_dashboard_rbac_user    = var.k8s_dashboard_rbac_user

  # K8s cluster creation depends on all VM modules #
  depends_on = [
    module.main,
  ]
} 