#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  main_resource_pool_name = "${var.cluster_name}-main-resource-pool"
  network_name            = "${var.cluster_name}-network"

  is_bridge = var.cluster_network_mode == "bridge"
}

#======================================================================================
# General Resources
#======================================================================================

#================================
# Resource pools and base volume
#================================

# Creates a resource pool for main (os) volumes #
resource "libvirt_pool" "main_resource_pool" {
  name = "${var.cluster_name}-main-resource-pool"
  type = "dir"
  path = pathexpand("${trimsuffix(var.hosts_mainResourcePoolPath, "/")}/${var.cluster_name}-main-resource-pool")
}

# Creates data resource pools #
resource "libvirt_pool" "data_resource_pools" {

  for_each = { for pool in var.hosts_dataResourcePools : pool.name => pool }

  name = "${var.cluster_name}-${each.key}-data-resource-pool"
  type = "dir"
  path = pathexpand("${trimsuffix(each.value.path, "/")}/${var.cluster_name}-data-resource-pool")
}

# Creates base OS image for nodes in a cluster #
resource "libvirt_volume" "base_volume" {
  name   = "base_volume"
  pool   = libvirt_pool.main_resource_pool.name
  source = pathexpand(var.cluster_nodeTemplate_os_source)

  # Requires resource pool to be initialized #
  depends_on = [libvirt_pool.main_resource_pool]
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
  cluster_name            = var.cluster_name
  libvirt_provider_uri    = var.libvirt_provider_uri
  main_resource_pool_name = libvirt_pool.main_resource_pool.name
  base_volume_id          = libvirt_volume.base_volume.id
  network_id              = local.is_bridge ? null : module.network_module.0.network_id

  # Network related variables
  network_mode    = var.cluster_network_mode
  network_bridge  = var.cluster_network_bridge
  network_gateway = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr    = var.cluster_network_cidr

  # Load balancer specific variables #
  vm_name              = "${var.cluster_name}-${var.node_types.load_balancer}-${each.value.id}"
  vm_type              = var.node_types.load_balancer
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_os_networkInterface
  vm_dns               = var.cluster_nodeTemplate_dns
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_loadBalancer_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_loadBalancer_default_ram
  vm_main_disk_size    = each.value.mainDiskSize != null ? each.value.mainDiskSize : var.cluster_nodes_loadBalancer_default_mainDiskSize
  vm_data_disks        = []
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.main_resource_pool,
    libvirt_pool.data_resource_pools,
    libvirt_volume.base_volume
  ]
}

# Creates master nodes #
module "master_module" {

  source = "../vm"

  for_each = { for node in var.cluster_nodes_master_instances : node.id => node }

  # Variables from general resources #
  cluster_name            = var.cluster_name
  libvirt_provider_uri    = var.libvirt_provider_uri
  main_resource_pool_name = libvirt_pool.main_resource_pool.name
  base_volume_id          = libvirt_volume.base_volume.id
  network_id              = local.is_bridge ? null : module.network_module.0.network_id

  # Network related variables
  network_mode    = var.cluster_network_mode
  network_bridge  = var.cluster_network_bridge
  network_gateway = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr    = var.cluster_network_cidr

  # Master node specific variables #
  vm_name              = "${var.cluster_name}-${var.node_types.master}-${each.value.id}"
  vm_type              = var.node_types.master
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_os_networkInterface
  vm_dns               = var.cluster_nodeTemplate_dns
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_master_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_master_default_ram
  vm_main_disk_size    = each.value.mainDiskSize != null ? each.value.mainDiskSize : var.cluster_nodes_master_default_mainDiskSize
  vm_data_disks        = each.value.dataDisks != null ? (length(setintersection(each.value.dataDisks.*.pool, keys(libvirt_pool.data_resource_pools))) == length(distinct(each.value.dataDisks.*.pool)) ? each.value.dataDisks : null) : []
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.main_resource_pool,
    libvirt_pool.data_resource_pools,
    libvirt_volume.base_volume
  ]
}

# Creates worker nodes #
module "worker_module" {

  source = "../vm"

  for_each = { for node in var.cluster_nodes_worker_instances : node.id => node }

  # Variables from general resources #
  cluster_name            = var.cluster_name
  libvirt_provider_uri    = var.libvirt_provider_uri
  main_resource_pool_name = libvirt_pool.main_resource_pool.name
  base_volume_id          = libvirt_volume.base_volume.id
  network_id              = local.is_bridge ? null : module.network_module.0.network_id

  # Network related variables
  network_mode    = var.cluster_network_mode
  network_bridge  = var.cluster_network_bridge
  network_gateway = var.cluster_network_gateway != null ? var.cluster_network_gateway : cidrhost(var.cluster_network_cidr, 1)
  network_cidr    = var.cluster_network_cidr

  # Worker node specific variables #
  vm_name              = each.value.name != null ? "${var.cluster_name}-${each.value.name}" : "${var.cluster_name}-${var.node_types.worker}-${each.value.id}"
  vm_type              = var.node_types.worker
  vm_user              = var.cluster_nodeTemplate_user
  vm_update            = var.cluster_nodeTemplate_updateOnBoot
  vm_ssh_private_key   = var.cluster_nodeTemplate_ssh_privateKeyPath
  vm_ssh_known_hosts   = var.cluster_nodeTemplate_ssh_addToKnownHosts
  vm_network_interface = var.cluster_nodeTemplate_os_networkInterface
  vm_dns               = var.cluster_nodeTemplate_dns
  vm_cpu               = each.value.cpu != null ? each.value.cpu : var.cluster_nodes_worker_default_cpu
  vm_ram               = each.value.ram != null ? each.value.ram : var.cluster_nodes_worker_default_ram
  vm_main_disk_size    = each.value.mainDiskSize != null ? each.value.mainDiskSize : var.cluster_nodes_worker_default_mainDiskSize
  vm_data_disks        = each.value.dataDisks != null ? (length(setintersection(each.value.dataDisks.*.pool, keys(libvirt_pool.data_resource_pools))) == length(distinct(each.value.dataDisks.*.pool)) ? each.value.dataDisks : null) : []
  vm_id                = each.value.id
  vm_mac               = each.value.mac
  vm_ip                = each.value.ip

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
    libvirt_pool.main_resource_pool,
    libvirt_pool.data_resource_pools,
    libvirt_volume.base_volume
  ]
}
