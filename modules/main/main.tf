#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  vm_type = {
    load_balancer = "lb"
    master        = "master"
    worker        = "worker"
  }
  resource_pool_name = "${var.cluster_name}-resource-pool"
  network_name       = "${var.cluster_name}-network"
  is_bridge          = (var.cluster_network_mode == "bridge")
}

#=====================================================================================
# Provider specific
#=====================================================================================

# Sets libvirt provider's uri #
provider "libvirt" {
  uri = var.libvirt_provider_uri
}

#======================================================================================
# General Resources
#======================================================================================

#================================
# Resource pool and base volume
#================================

# Creates a resource pool for Kubernetes VM volumes #
resource "libvirt_pool" "resource_pool" {
  name = local.resource_pool_name
  type = "dir"
  path = pathexpand("${trimsuffix(var.libvirt_resource_pool_location, "/")}/${local.resource_pool_name}")
}

# Creates base OS image for nodes in a cluster #
resource "libvirt_volume" "base_volume" {
  name   = "base_volume"
  pool   = local.resource_pool_name
  source = pathexpand(var.cluster_nodeTemplate_image_source)

  # Requires resource pool to be initialized #
  depends_on = [libvirt_pool.resource_pool]
}

#======================================================================================
# Modules
#======================================================================================

#================================
# Network
#================================

# Creates network #
module "network_module" {

  count = local.is_bridge ? 0 : 1

  source = "../network/"

  network_name   = local.network_name
  network_mode   = var.cluster_network_mode
  network_bridge = var.cluster_network_bridge
  network_cidr   = var.cluster_network_cidr
}

#================================
# Virtual machines
#================================

# Create HAProxy load balancer #
module "lb_module" {
  source = "../vm"

  for_each = { for node in var.cluster_nodes_loadBalancer_instances : node.id => node }

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  network_id           = local.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = local.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Load balancer specific variables #
  vm_name              = "${var.cluster_name}-${local.vm_type.load_balancer}-${each.value.id}"
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = pathexpand(var.cluster_nodeTemplate_ssh_privateKeyPath)
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = var.cluster_nodes_loadBalancer_default_cpu     #each.value.cpu != null ? each.value.cpu : var.cluster_nodes_loadBalancer_default_cpu
  vm_ram               = var.cluster_nodes_loadBalancer_default_ram     #each.value.ram != null ? each.value.ram : var.lb_default_ram
  vm_storage           = var.cluster_nodes_loadBalancer_default_storage #each.value.storage != null ? each.value.storage : var.lb_default_storage
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

# Creates master nodes #
module "master_module" {
  source = "../vm"

  for_each = { for node in var.cluster_nodes_master_instances : node.id => node }

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  network_id           = local.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = local.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Master node specific variables #
  vm_name              = "${var.cluster_name}-${local.vm_type.master}-${each.value.id}"
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = pathexpand(var.cluster_nodeTemplate_ssh_privateKeyPath)
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = var.cluster_nodes_master_default_cpu     #each.value.cpu != null ? each.value.cpu : var.cluster_nodes_master_default_cpu
  vm_ram               = var.cluster_nodes_master_default_ram     #each.value.ram != null ? each.value.ram : var.cluster_nodes_master_default_ram
  vm_storage           = var.cluster_nodes_master_default_storage #each.value.storage != null ? each.value.storage : var.cluster_nodes_master_default_storage
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

# Creates worker nodes #
module "worker_module" {
  source = "../vm"

  for_each = { for node in var.cluster_nodes_worker_instances : node.id => node }

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  network_id           = local.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = local.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Worker node specific variables #
  vm_name              = "${var.cluster_name}-${local.vm_type.worker}-${each.value.id}"
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = pathexpand(var.cluster_nodeTemplate_ssh_privateKeyPath)
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = var.cluster_nodes_worker_default_cpu     #each.value.cpu != null ? each.value.cpu : var.cluster_nodes_worker_default_cpu
  vm_ram               = var.cluster_nodes_worker_default_ram     #each.value.ram != null ? each.value.ram : var.cluster_nodes_worker_default_ram
  vm_storage           = var.cluster_nodes_worker_default_storage #each.value.storage != null ? each.value.storage : var.cluster_nodes_worker_default_storage
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

#================================
# Cluster
#================================

# Configures k8s cluster using Kubespray #
module "k8s_cluster" {
  source = "../cluster"

  action = var.action

  # VM variables #
  vm_user            = var.cluster_nodeTemplate_user
  vm_ssh_private_key = pathexpand(var.cluster_nodeTemplate_ssh_privateKeyPath)
  vm_distro          = var.cluster_nodeTemplate_image_distro
  vm_network_interface = local.is_bridge ? var.cluster_network_bridge : var.cluster_nodeTemplate_networkInterface
  
  worker_node_label    = var.cluster_nodes_worker_default_label
  lb_vip               = var.cluster_nodes_loadBalancer_vip
  lb_nodes             = [for node in module.lb_module : node.vm_info]
  master_nodes         = [for node in module.master_module : node.vm_info]
  worker_nodes         = [for node in module.worker_module : node.vm_info]

  # K8s cluster variables #
  kubernetes_version           = var.kubernetes_version
  kubernetes_networkPlugin    = var.kubernetes_networkPlugin
  kubernetes_dnsMode          = var.kubernetes_dnsMode
  kubernetes_kubespray_url     = var.kubernetes_kubespray_url
  kubernetes_kubespray_version = var.kubernetes_kubespray_version
  kubernetes_kubespray_addons_enabled       = false
  kubernetes_kubespray_addons_configPath          = ""
  kubernetes_other_copyKubeconfig        = var.kubernetes_other_copyKubeconfig

  # Other #
  
  k8s_dashboard_rbac_enabled = var.k8s_dashboard_rbac_enabled
  k8s_dashboard_rbac_user    = var.k8s_dashboard_rbac_user

  # Kubespray addons #
  
  k8s_dashboard_enabled                 = var.k8s_dashboard_enabled
  helm_enabled                          = var.helm_enabled
  local_path_provisioner_enabled        = var.local_path_provisioner_enabled
  local_path_provisioner_version        = var.local_path_provisioner_version
  local_path_provisioner_namespace      = var.local_path_provisioner_namespace
  local_path_provisioner_storage_class  = var.local_path_provisioner_storage_class
  local_path_provisioner_reclaim_policy = var.local_path_provisioner_reclaim_policy
  local_path_provisioner_claim_root     = var.local_path_provisioner_claim_root
  metallb_enabled                       = var.metallb_enabled
  metallb_version                       = var.metallb_version
  metallb_port                          = var.metallb_port
  metallb_cpu_limit                     = var.metallb_cpu_limit
  metallb_mem_limit                     = var.metallb_mem_limit
  metallb_protocol                      = var.metallb_protocol
  metallb_ip_range                      = var.metallb_ip_range
  metallb_peers                         = var.metallb_peers

  # K8s cluster creation depends on all VM modules #
  depends_on = [
    module.lb_module,
    module.worker_module,
    module.master_module
  ]
}