#======================================================================================
# General configuration
#======================================================================================

variable "action" {
  type        = string
  description = "Action (create, upgrade, scale)."
  default     = "create"
  nullable    = false

  validation {
    condition     = contains(["create", "upgrade", "scale"], var.action)
    error_message = "Variable 'action' is invalid. Possible values are: ['create', 'upgrade', 'scale']."
  }
}

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI."
  default     = "qemu:///system"
  nullable    = false
}

variable "hosts_mainResourcePoolPath" {
  type        = string
  description = "Path where main resource pool will be initialized."
  default     = "/var/lib/libvirt/images/"
  nullable    = false
}

variable "hosts_dataResourcePools" {
  type = list(object({
    name : string
    path : string
  }))
  description = "Location where main resource pool will be initialized."
  default     = []
  nullable    = false
}

#======================================================================================
# Cluster infrastructure configuration
#======================================================================================

variable "cluster_name" {
  type        = string
  description = "Cluster name used as a prefix for various cluster component names."
  nullable    = false
}

#================================
# Node template
#================================

variable "cluster_nodeTemplate_user" {
  type        = string
  description = "Username used to SSH to the virtual machines."
  nullable    = false
}

variable "cluster_nodeTemplate_ssh_privateKeyPath" {
  type        = string
  description = "Location of private SSH key that will be used for virtual machines."
}

variable "cluster_nodeTemplate_ssh_addToKnownHosts" {
  type        = bool
  description = "Add virtual machines to SSH known hosts."
  default     = true
  nullable    = false
}

variable "cluster_nodeTemplate_os_source" {
  type        = string
  description = "OS source, which can be path on host's filesystem or URL."
}

variable "cluster_nodeTemplate_os_networkInterface" {
  type        = string
  description = "Operating system (os) network interface, which is predefined for the os image."
}

variable "cluster_nodeTemplate_dns" {
  type        = list(string)
  description = "List of DNS servers used by virtual machines."
  default     = []
  nullable    = false
}

variable "cluster_nodeTemplate_updateOnBoot" {
  type        = bool
  description = "Update system on boot."
  default     = true
  nullable    = false
}

#================================
# Cluster network
#================================

variable "cluster_network_mode" {
  type        = string
  description = "Network mode."
  nullable    = false
}

variable "cluster_network_bridge" {
  type        = string
  description = "Network (virtual) bridge."
}

variable "cluster_network_gateway" {
  type        = string
  description = "Network gateway."
  default     = null
}

variable "cluster_network_cidr" {
  type        = string
  description = "Network CIDR."
}

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

variable "cluster_nodes_loadBalancer_vip" {
  type        = string
  description = "HAProxy load balancer virtual IP address (VIP)."
}

variable "cluster_nodes_loadBalancer_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the HAProxy load balancer."
}

variable "cluster_nodes_loadBalancer_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the HAProxy load balancer."
}

variable "cluster_nodes_loadBalancer_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the HAProxy load balancer."
}

variable "cluster_nodes_loadBalancer_instances" {
  type = list(object({
    id           = number
    host         = optional(string)
    mac          = optional(string)
    ip           = optional(string)
    cpu          = optional(number)
    ram          = optional(number)
    mainDiskSize = optional(number)
  }))
  description = "HAProxy load balancer node instances."
}

#======================================================================================
# Master node VMs parameters
#======================================================================================

variable "cluster_nodes_master_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the master node."
}

variable "cluster_nodes_master_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the master node."
}

variable "cluster_nodes_master_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the master node."
}

variable "cluster_nodes_master_default_dataDisks" {
  type = list(object({
    name : string
    pool : string
    size : number
  }))
  description = "List of additional data disks that are attached to the master node."
  default     = []
  nullable    = false
}

variable "cluster_nodes_master_instances" {
  type = list(object({
    id           = number
    host         = optional(string)
    mac          = optional(string)
    ip           = optional(string)
    cpu          = optional(number)
    ram          = optional(number)
    mainDiskSize = optional(number)
    dataDisks = optional(list(object({
      name : string
      pool : string
      size : number
    })))
  }))
  description = "Master node instances (control plane)"
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

variable "cluster_nodes_worker_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the worker node."
}

variable "cluster_nodes_worker_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the worker node."
}

variable "cluster_nodes_worker_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the worker node."
}

variable "cluster_nodes_worker_default_dataDisks" {
  type = list(object({
    name : string
    pool : string
    size : number
  }))
  description = "List of additional data disks that are attached to the worker node."
  default     = []
  nullable    = false
}

variable "cluster_nodes_worker_instances" {
  type = list(object({
    id           = number
    name         = optional(string)
    host         = optional(string)
    mac          = optional(string)
    ip           = optional(string)
    cpu          = optional(number)
    ram          = optional(number)
    mainDiskSize = optional(number)
    dataDisks = optional(list(object({
      name : string
      pool : string
      size : number
    })))
  }))
  description = "Worker node instances."
}

#======================================================================================
# Other internal variables
#======================================================================================

variable "node_types" {
  type = object({
    load_balancer = string
    master        = string
    worker        = string
  })
  description = "Node types."
}
