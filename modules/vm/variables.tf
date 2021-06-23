# ==================================== #
# Outputs                              #
# ==================================== #

output "vm_info" {
  value = {
    name = libvirt_domain.vm_domain.name,
    ip   = libvirt_domain.vm_domain.network_interface.0.addresses.0
  }
  description = "VM's info containing it's name and an IP address"
}

# ==================================== #
# Variables dependent on parent module #
# ==================================== #

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI"
}

variable "resource_pool_name" {
  type        = string
  description = "Resource pool name"
}

variable "base_volume_id" {
  type        = string
  description = "Base image voulme ID"
}

variable "cloud_init_id" {
  type        = string
  description = "Cloud init disk ID"
}

variable "network_id" {
  type        = string
  description = "Id of the network in which VM resides"
}

# ==================================== #
# VM variables                         #
# ==================================== #

#============================#
# General                    #
#============================#

variable "vm_type" {
  type        = string
  description = "Possible virtual machine types are: [master, worker, lb]"

  validation {
    condition     = contains(["master", "worker", "lb"], var.vm_type)
    error_message = "Variable 'vm_type' is invalid.\nPossible values are: [\"master\", \"worker\", \"lb\"]."
  }
}

variable "vm_user" {
  type        = string
  description = "VM's SSH user"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private key for VM's SSH user"
}

variable "vm_ssh_known_hosts" {
  type        = bool
  description = "Add virtual machine SSH known hosts"
}

#============================#
# Specific                   #
#============================#

variable "vm_name" {
  type        = string
  description = "VM name"
}

variable "vm_id" {
  type        = number
  description = "Unique VM id used to differentiate VMs of the same type."
}

variable "vm_cpu" {
  type        = number
  description = "The number of vCPU allocated to the virtual machine"
}

variable "vm_ram" {
  type        = number
  description = "The amount of RAM allocated to the virtual machine"
}

variable "vm_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the virtual machine"
}

variable "vm_mac" {
  type        = string
  description = "The MAC address of the virtual machine"

  validation {
    condition = (
      var.vm_mac == null
      || can(regex("^([0-9A-Fa-f]){2}(:[0-9A-Fa-f]{2}){5}$", var.vm_mac))
    )
    error_message = "Invalid MAC address provided to VM.\nPlease check the following variables:\n - 'lb_nodes',\n - 'master_nodes',\n - 'worker_nodes'.\n\nNote that setting MAC to 'null' causes random valid MAC to be generated."
  }

}

variable "vm_ip" {
  type        = string
  description = "The IP address of the virtual machine"

  validation {
    condition = (
      var.vm_ip == null
      || can(regex("^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(.([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])){3}$", var.vm_ip))
    )
    error_message = "Invalid IP address provided to VM.\nPlease check the following variables:\n - 'lb_nodes',\n - 'master_nodes',\n - 'worker_nodes'.\n\nNote that setting IP to 'null' causes random valid IP to be generated."
  }
}

