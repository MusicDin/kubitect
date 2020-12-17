# ==================================== #
# Variables dependent on parent module #
# ==================================== #

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

# ==================================== #
# VM variables                         #
# ==================================== #

#============================#
# General                    #
#============================#

variable "vm_type" {
  type        = string
  description = "Possible virtual machine types are: [master, worker, lb]"
}

variable "vm_index" {
  type        = number
  description = "Index of VM. Used to differentiate VMs of the same type."
}

variable "vm_network_name" {
  type        = string
  description = "Network name in which VM resides"
}

variable "vm_user" {
  type        = string
  description = "VM's SSH user"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private key for VM's SSH user"
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

