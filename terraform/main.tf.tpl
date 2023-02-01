{{- $hosts := .Hosts -}}
{{- $defHost := defaultHost $hosts -}}

#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  # Configuration files
  # Infra config can be null, as it is created after first initialization
  # of the cluster.
  config          = yamldecode(file(var.config_path))
  defaults_config = yamldecode(file(var.defaults_config_path))
  infra_config    = try(yamldecode(file(var.infra_config_path)), null)

  node_types = {
    load_balancer = "lb"
    master        = "master"
    worker        = "worker"
  }
}


#=====================================================================================
# Providers
#=====================================================================================
{{ range .Hosts }}
provider "libvirt" {
  alias = "{{ .Name }}"
  uri   = "{{ hostUri .}}"
}
{{- end }}


#======================================================================================
# Modules
#======================================================================================

module "base" {
  source = "./modules/base"

  # Configuration files
  config          = local.config
  infra_config    = local.infra_config
  defaults_config = local.defaults_config
}
{{- range .Hosts }}
  {{- $defSelector := "" }}
  {{- if eq .Name $defHost.Name }}
    {{- $defSelector = " || try(node.host, null) == null" }}
  {{- end }}

module "host_{{ .Name }}" {
  source = "./modules/host"

  # General
  action = var.action

  # Resource pools
  hosts_mainResourcePoolPath = "{{ .MainResourcePoolPath }}"
  hosts_dataResourcePools    = try(local.config.hosts[index(local.config.hosts.*.name, "{{ .Name }}")].dataResourcePools, null)

  # Cluster name and node template
  cluster_name                             = try(local.config.cluster.name, null)
  cluster_nodeTemplate_user                = module.base.user
  cluster_nodeTemplate_ssh_privateKeyPath  = module.base.ssh.privateKeyPath
  cluster_nodeTemplate_ssh_addToKnownHosts = try(local.config.cluster.nodeTemplate.ssh.addToKnownHosts, null)
  cluster_nodeTemplate_os_source           = module.base.os.source
  cluster_nodeTemplate_os_networkInterface = module.base.os.networkInterface
  cluster_nodeTemplate_dns                 = try(local.config.cluster.nodeTemplate.dns, null)
  cluster_nodeTemplate_updateOnBoot        = try(local.config.cluster.nodeTemplate.updateOnBoot, null)
  cluster_nodeTemplate_cpuMode             = try(local.config.cluster.nodeTemplate.cpuMode, null)

  # Network configuration
  cluster_network_mode    = try(local.config.cluster.network.mode, null)
  cluster_network_cidr    = try(local.config.cluster.network.cidr, null)
  cluster_network_gateway = try(local.config.cluster.network.gateway, null)
  cluster_network_bridge  = try(local.config.cluster.network.bridge, null)

  # HAProxy load balancer VMs parameters
  cluster_nodes_loadBalancer_vip                  = try(local.config.cluster.nodes.loadBalancer.vip, null)
  cluster_nodes_loadBalancer_default_cpu          = try(local.config.cluster.nodes.loadBalancer.default.cpu, null)
  cluster_nodes_loadBalancer_default_ram          = try(local.config.cluster.nodes.loadBalancer.default.ram, null)
  cluster_nodes_loadBalancer_default_mainDiskSize = try(local.config.cluster.nodes.loadBalancer.default.mainDiskSize, null)
  cluster_nodes_loadBalancer_instances = [
    for node in try(flatten([local.config.cluster.nodes.loadBalancer.instances]), []) : node
    if node != null && (try(node.host, null) == "{{ .Name }}"{{ $defSelector }})
  ]

  # Master node VMs parameters
  cluster_nodes_master_default_cpu          = try(local.config.cluster.nodes.master.default.cpu, null)
  cluster_nodes_master_default_ram          = try(local.config.cluster.nodes.master.default.ram, null)
  cluster_nodes_master_default_mainDiskSize = try(local.config.cluster.nodes.master.default.mainDiskSize, null)
  cluster_nodes_master_instances = [
    for node in try(flatten([local.config.cluster.nodes.master.instances]), []) : node
    if node != null && (try(node.host, null) == "{{ .Name }}"{{ $defSelector }})
  ]

  # Worker node VMs parameters
  cluster_nodes_worker_default_cpu          = try(local.config.cluster.nodes.worker.default.cpu, null)
  cluster_nodes_worker_default_ram          = try(local.config.cluster.nodes.worker.default.ram, null)
  cluster_nodes_worker_default_mainDiskSize = try(local.config.cluster.nodes.worker.default.mainDiskSize, null)
  cluster_nodes_worker_instances = [
    for node in try(flatten([local.config.cluster.nodes.worker.instances]), []) : node
    if node != null && (try(node.host, null) == "{{ .Name }}"{{ $defSelector }})
  ]

  # Other
  node_types = local.node_types

  providers = {
    libvirt = libvirt.{{ .Name }}
  }
}
{{- end }}


#================================
# Infrastructure output
#================================

{{- $modules := list }}
{{- range .Hosts }}
  {{- $modules = deref .Name | printf "module.host_%s.nodes" | append $modules  }}
{{- end }}
{{- $modules = $modules | join ", " }}

# Outputs evaluated cluster information #
module "output" {
  source = "./modules/output"

  # VM variables #
  vm_user = module.base.user
  vm_ssh  = module.base.ssh
  vm_os   = module.base.os
  lb_vip  = try(local.config.cluster.nodes.loadBalancer.vip, null)

  lb_nodes = [
    for node in flatten([{{ $modules }}]) :
    node if node.type == local.node_types.load_balancer
  ]

  master_nodes = [
    for node in flatten([{{ $modules }}]) :
    node if node.type == local.node_types.master
  ]

  worker_nodes = [
    for node in flatten([{{ $modules }}]) :
    node if node.type == local.node_types.worker
  ]

  # K8s cluster creation depends on all VM modules #
  depends_on = [
  {{- range .Hosts }}
    module.host_{{ .Name }},
  {{- end }}
  ]
}

# Creates a configuration file that contains evaluated
# information about the created infrastructure.
resource "local_file" "output" {
  content         = replace(yamlencode(module.output), "/((?:^|\n)[\\s-]*)\"([\\w-]+)\":/", "$1$2:")
  filename        = var.infra_config_path
  file_permission = 0600
}