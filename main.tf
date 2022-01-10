#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  config = var.config_type == "yaml" ? yamldecode(file(pathexpand(var.config_path))) : null
}

#======================================================================================
# Modules
#======================================================================================

#================================
# YAML decoder
#================================

module "yaml_decoder" {
  source = "./modules/main"

  # Module contains provider and cannot be skipped using count
  #count = var.config_type ? 1 : 0

  # General
  action = var.action

  # Libvirt
  libvirt_provider_uri           = "qemu:///system" # TBD
  libvirt_resource_pool_location = "/var/lib/libvirt/pools/"


  # 
  cluster_name                             = local.config.cluster.name
  cluster_nodeTemplate_user                = local.config.cluster.nodeTemplate.user
  cluster_nodeTemplate_ssh_privateKeyPath  = local.config.cluster.nodeTemplate.ssh.privateKeyPath
  cluster_nodeTemplate_ssh_addToKnownHosts = local.config.cluster.nodeTemplate.ssh.addToKnownHosts
  cluster_nodeTemplate_image_distro        = local.config.cluster.nodeTemplate.image.distro
  cluster_nodeTemplate_image_source        = local.config.cluster.nodeTemplate.image.source


  # Global VM configuration
  cluster_nodeTemplate_networkInterface = var.cluster_nodeTemplate_networkInterface
  cluster_nodeTemplate_updateOnBoot     = var.cluster_nodeTemplate_updateOnBoot

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
  kubernetes_version              = var.kubernetes_version
  kubernetes_networkPlugin        = var.kubernetes_networkPlugin
  kubernetes_dnsMode              = var.kubernetes_dnsMode
  kubernetes_kubespray_url        = var.kubernetes_kubespray_url
  kubernetes_kubespray_version    = var.kubernetes_kubespray_version
  kubernetes_kubespray_addons_enabled = false # var.kubernetes_kubespray_addons_enabled
  kubernetes_kubespray_addons_configPath    = ""    # var.kubernetes_kubespray_addons_configPath
  kubernetes_other_copyKubeconfig = var.kubernetes_other_copyKubeconfig

  
}


#================================
# Original configuration
#================================

/*
module "tf_decoder" {
  source = "./modules/main"
  count = !var.config_type ? 1 : 0

  # General
  action = var.action
  config_type = var.config_type # for module count
  config_path = var.config_path # for decoding yaml

  # Libvirt
  libvirt_provider_uri = var.libvirt_provider_uri
  libvirt_resource_pool_location = var.libvirt_resource_pool_location

  # Global VM configuration
  cluster_nodeTemplate_networkInterface = var.cluster_nodeTemplate_networkInterface
  cluster_nodeTemplate_updateOnBoot     = var.cluster_nodeTemplate_updateOnBoot

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
  kubernetes_version              = var.kubernetes_version
  kubernetes_networkPlugin        = var.kubernetes_networkPlugin
  kubernetes_dnsMode              = var.kubernetes_dnsMode
  kubernetes_kubespray_url        = var.kubernetes_kubespray_url
  kubernetes_kubespray_version    = var.kubernetes_kubespray_version
  kubernetes_kubespray_addons_enabled = false # var.kubernetes_kubespray_addons_enabled
  kubernetes_kubespray_addons_configPath    = ""    # var.kubernetes_kubespray_addons_configPath
  kubernetes_other_copyKubeconfig = var.kubernetes_other_copyKubeconfig

}
*/

