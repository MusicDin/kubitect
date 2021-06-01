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

variable "network_name" {
  type        = string
  description = "Network name in which VM resides"
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

variable "vm_index" {
  type        = number
  description = "Index of VM. Used to differentiate VMs of the same type."
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
  type        = string
  description = "Add virtual machine SSH known hosts"
}

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
}

#============================#
# Specific                   #
#============================#

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
}

variable "vm_ip" {
  type        = string
  description = "The IP address of the virtual machine"
}

