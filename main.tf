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
  network_mode           = var.network_mode
  network_bridge         = var.network_bridge
  network_mac            = var.network_mac
  network_gateway        = var.network_gateway
  network_mask_bits      = var.network_mask_bits
  network_dhcp_ip_start  = var.network_dhcp_ip_start
  network_dhcp_ip_end    = var.network_dhcp_ip_end
}

# Create HAProxy load balancer #
module "lb_module" {
  source = "./modules/vm"

  count = length(var.vm_lb_macs_ips)

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  cloud_init_id        = libvirt_cloudinit_disk.cloud_init.id
  network_name         = var.network_name

  # Load balancer specific variables #
  vm_index           = count.index
  vm_type            = "lb"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_lb_cpu
  vm_ram             = var.vm_lb_ram
  vm_storage         = var.vm_lb_storage
  vm_mac             = keys(var.vm_lb_macs_ips)[count.index]
  vm_ip              = values(var.vm_lb_macs_ips)[count.index]

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
  source = "./modules/vm"

  count = length(var.vm_master_macs_ips)

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  cloud_init_id        = libvirt_cloudinit_disk.cloud_init.id
  network_name         = var.network_name

  # Master node specific variables #
  vm_index           = count.index
  vm_type            = "master"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_master_cpu
  vm_ram             = var.vm_master_ram
  vm_storage         = var.vm_master_storage
  vm_mac             = keys(var.vm_master_macs_ips)[count.index]
  vm_ip              = values(var.vm_master_macs_ips)[count.index]

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
  source = "./modules/vm"

  count = length(var.vm_worker_macs_ips)

  # Variables from general resources #
  libvirt_provider_uri = var.libvirt_provider_uri
  resource_pool_name   = libvirt_pool.resource_pool.name
  base_volume_id       = libvirt_volume.base_volume.id
  cloud_init_id        = libvirt_cloudinit_disk.cloud_init.id
  network_name         = var.network_name

  # Worker node specific variables #
  vm_index           = count.index
  vm_type            = "worker"
  vm_user            = var.vm_user
  vm_ssh_private_key = var.vm_ssh_private_key
  vm_ssh_known_hosts = var.vm_ssh_known_hosts
  vm_name_prefix     = var.vm_name_prefix
  vm_cpu             = var.vm_worker_cpu
  vm_ram             = var.vm_worker_ram
  vm_storage         = var.vm_worker_storage
  vm_mac             = keys(var.vm_worker_macs_ips)[count.index]
  vm_ip              = values(var.vm_worker_macs_ips)[count.index]

  # Dependancy takes care that resource pool is not removed before volumes are #
  # Also network must be created before VM is initialized #
  depends_on = [
    module.network_module,
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
  vm_network_interface = var.vm_network_interface

  # K8s cluster variables
  k8s_kubespray_url     = var.k8s_kubespray_url
  k8s_kubespray_version = var.k8s_kubespray_version
  k8s_version           = var.k8s_version
  k8s_network_plugin    = var.k8s_network_plugin
  k8s_dns_mode          = var.k8s_dns_mode

  # Other
  k8s_copy_kubeconfig        = var.k8s_copy_kubeconfig
  k8s_dashboard_rbac_enabled = var.k8s_dashboard_rbac_enabled
  k8s_dashboard_rbac_user    = var.k8s_dashboard_rbac_user

  # Kubespray addons
  kubespray_custom_addons_enabled       = var.kubespray_custom_addons_enabled
  kubespray_custom_addons_path          = var.kubespray_custom_addons_path
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

  # K8s cluster creation depends on network and all VMs
  depends_on = [
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
    network_interface = var.vm_network_interface
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
