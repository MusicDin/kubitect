#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  resource_pool_name = "${var.cluster_name}-resource-pool"
  network_name       = "${var.cluster_name}-network"
}

#=====================================================================================
# Provider specific
#=====================================================================================

# Sets libvirt provider's uri #
#provider "libvirt" {
#  uri = var.libvirt_provider_uri
#}

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

  count = var.internal.is_bridge ? 0 : 1

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
  network_id           = var.internal.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = var.internal.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Load balancer specific variables #
  vm_name              = "${var.cluster_name}-${var.internal.vm_types.load_balancer}-${each.value.id}"
  vm_type              = var.internal.vm_types.load_balancer
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_loadBalancer_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_loadBalancer_default_ram
  vm_storage           = each.value.storage != null ? each.value.storage : var.cluster_nodes_loadBalancer_default_storage
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
  network_id           = var.internal.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = var.internal.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Master node specific variables #
  vm_name              = "${var.cluster_name}-${var.internal.vm_types.master}-${each.value.id}"
  vm_type              = var.internal.vm_types.master
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_master_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_master_default_ram
  vm_storage           = each.value.storage != null ? each.value.storage : var.cluster_nodes_master_default_storage
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
  network_id           = var.internal.is_bridge ? null : module.network_module.0.network_id

  is_bridge        = var.internal.is_bridge
  network_bridge   = var.cluster_network_bridge
  network_gateway  = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr     = var.cluster_network_cidr
  network_dns_list = var.cluster_network_dns

  # Worker node specific variables #
  vm_name              = "${var.cluster_name}-${var.internal.vm_types.worker}-${each.value.id}"
  vm_type              = var.internal.vm_types.worker
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_networkInterface
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_worker_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_worker_default_ram
  vm_storage           = each.value.storage != null ? each.value.storage : var.cluster_nodes_worker_default_storage
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
