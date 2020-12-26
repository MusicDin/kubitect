#=====================================================================================
# Provider specific
#=====================================================================================

# Sets libvirt provider's uri #
provider "libvirt" {
  uri = var.libvirt_provider_uri
}

#======================================================================================
# Modules
#======================================================================================

# Creates network #
module "network_module" {
  source = "./modules/network/"

  libvirt_provider_uri   = var.libvirt_provider_uri
  network_name           = var.network_name
  network_mac            = var.network_mac
  network_gateway        = var.network_gateway
  network_mask           = var.network_mask
  network_nat_port_start = var.network_nat_port_start
  network_nat_port_end   = var.network_nat_port_end
  network_dhcp_ip_start  = var.network_dhcp_ip_start
  network_dhcp_ip_end    = var.network_dhcp_ip_end
  network_virtual_bridge = var.network_virtual_bridge
  vm_lb_macs_ips         = var.vm_lb_macs_ips
  vm_master_macs_ips     = var.vm_master_macs_ips
  vm_worker_macs_ips     = var.vm_worker_macs_ips
}

# Create HAProxy load balancer #
module "lb_module" {
  source = "./modules/vm"

  count = length(var.vm_lb_macs_ips)

  # Variables from general resources #
  resource_pool_name = libvirt_pool.resource_pool.name
  base_volume_id     = libvirt_volume.base_volume.id
  cloud_init_id      = libvirt_cloudinit_disk.cloud_init.id

  # Load balancer specific variables #
  vm_index           = count.index
  vm_type            = "lb"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_network_name    = var.network_name
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_lb_cpu
  vm_ram             = var.vm_lb_ram
  vm_storage         = var.vm_lb_storage
  vm_mac             = keys(var.vm_lb_macs_ips)[count.index]
  vm_ip              = values(var.vm_lb_macs_ips)[count.index]

  # Dependancy takes care that resource pool is not removed before volumes are #
  depends_on = [
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

# Creates master nodes #
module "master_module" {
  source = "./modules/vm"

  count = length(var.vm_master_macs_ips)

  # Variables from general resources #
  resource_pool_name = libvirt_pool.resource_pool.name
  base_volume_id     = libvirt_volume.base_volume.id
  cloud_init_id      = libvirt_cloudinit_disk.cloud_init.id

  # Master node specific variables #
  vm_index           = count.index
  vm_type            = "master"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_network_name    = var.network_name
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_master_cpu
  vm_ram             = var.vm_master_ram
  vm_storage         = var.vm_master_storage
  vm_mac             = keys(var.vm_master_macs_ips)[count.index]
  vm_ip              = values(var.vm_master_macs_ips)[count.index]

  # Dependancy takes care that resource pool is not removed before volumes are #
  depends_on = [
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

# Creates worker nodes #
module "worker_module" {
  source = "./modules/vm"

  count = length(var.vm_worker_macs_ips)

  # Variables from general resources #
  resource_pool_name = libvirt_pool.resource_pool.name
  base_volume_id     = libvirt_volume.base_volume.id
  cloud_init_id      = libvirt_cloudinit_disk.cloud_init.id

  # Worker node specific variables #
  vm_index           = count.index
  vm_type            = "worker"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_network_name    = var.network_name
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_worker_cpu
  vm_ram             = var.vm_worker_ram
  vm_storage         = var.vm_worker_storage
  vm_mac             = keys(var.vm_worker_macs_ips)[count.index]
  vm_ip              = values(var.vm_worker_macs_ips)[count.index]

  # Dependancy takes care that resource pool is not removed before volumes are #
  depends_on = [
    libvirt_pool.resource_pool,
    libvirt_volume.base_volume
  ]
}

# Configures k8s cluster using Kubespray #
module "k8s_cluster" {
  source = "./modules/cluster"

  action = var.action

  # VM variables
  vm_distro            = var.vm_distro
  vm_user              = var.vm_user
  vm_ssh_private_key   = var.vm_ssh_private_key
  vm_name_prefix       = var.vm_name_prefix
  vm_worker_ips        = values(var.vm_worker_macs_ips)
  vm_worker_node_label = var.vm_worker_node_label
  vm_master_ips        = values(var.vm_master_macs_ips)
  vm_lb_ips            = values(var.vm_lb_macs_ips)
  vm_lb_vip            = var.vm_lb_vip
  network_interface    = var.network_interface

  # K8s cluster variables
  k8s_kubespray_url     = var.k8s_kubespray_url
  k8s_kubespray_version = var.k8s_kubespray_version
  k8s_version           = var.k8s_version
  k8s_network_plugin    = var.k8s_network_plugin
  k8s_dns_mode          = var.k8s_dns_mode

  # K8s cluster creation depends on network and all VMs
  depends_on = [
    module.network_module,
    module.lb_module,
    module.worker_module,
    module.master_module
  ]
}


#======================================================================================
# General Resources
#======================================================================================

#================================
# Resource pool and base volume
#================================

# Creates a resource pool for Kubernetes VM volumes #
resource "libvirt_pool" "resource_pool" {
  name = var.libvirt_resource_pool_name
  type = "dir"
  path = "${var.libvirt_resource_pool_location}${var.libvirt_resource_pool_name}"
}

# Creates base OS image for nodes in a cluster #
resource "libvirt_volume" "base_volume" {
  name   = "base_volume"
  pool   = var.libvirt_resource_pool_name
  source = var.vm_image_source

  depends_on = [libvirt_pool.resource_pool]
}

#================================
# Cloud-init
#================================

# Public ssh key for vm (it is directly injected in cloud-init configuration) #
data "template_file" "public_ssh_key" {
  template = file("${var.vm_ssh_private_key}.pub")
}

# Network bridge configuration (for cloud-init) #
data "template_file" "cloud_init_network_tpl" {
  template = file("templates/cloud_init_network.tpl")

  vars = {
    network_interface = var.network_interface
  }
}

# Creates network bridge configuration file from template #
resource "local_file" "cloud_init_network_file" {
  content  = data.template_file.cloud_init_network_tpl.rendered
  filename = "config/cloud_init_network.cfg"
}

# Cloud-init configuration template #
data "template_file" "cloud_init_tpl" {
  template = file("templates/cloud_init.tpl")

  vars = {
    user           = var.vm_user
    ssh_public_key = data.template_file.public_ssh_key.rendered
  }
}

# Creates cloud-init configuration file from template #
resource "local_file" "cloud_init_file" {
  content  = data.template_file.cloud_init_tpl.rendered
  filename = "config/cloud_init.cfg"
}

# Initializes cloud-init disk for user data#
resource "libvirt_cloudinit_disk" "cloud_init" {
  name           = "cloud-init.iso"
  pool           = libvirt_pool.resource_pool.name
  user_data      = data.template_file.cloud_init_tpl.rendered
  network_config = data.template_file.cloud_init_network_tpl.rendered
}
