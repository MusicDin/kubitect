# ==================================== #
# Variables dependent on parent module #
# ==================================== #

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI"
}

variable "main_resource_pool_name" {
  type        = string
  description = "Main resource pool name"
}

variable "base_volume_id" {
  type        = string
  description = "Base image voulme ID"
}

variable "network_id" {
  type        = string
  description = "Id of the network in which VM resides"
}

variable "cluster_name" {
  type        = string
  description = "Cluster name"
}

# ==================================== #
# Network variables                    #
# ==================================== #

variable "network_mode" {
  type        = string
  description = "Network mode."
}

variable "network_bridge" {
  type        = string
  description = "Network bridge (used only when network mode is 'bridge')"
}

variable "network_gateway" {
  type        = string
  description = "Network gateway (used only when network mode is 'bridge')"
}

variable "network_cidr" {
  type        = string
  description = "Network CIDR"
}

# ==================================== #
# VM variables                         #
# ==================================== #

#============================#
# General                    #
#============================#

variable "vm_type" {
  type        = string
  description = "Vitual machine type."
}

variable "vm_network_interface" {
  type        = string
  description = "Network interface"
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

variable "vm_dns" {
  type        = list(string)
  description = "List of DNS servers used by VMs"
}

variable "vm_update" {
  type        = bool
  description = "Update system when ready"
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

variable "vm_cpuMode" {
  type        = string
  description = "The libvirt CPU emulation mode."
  default     = "custom"
  nullable    = false
}

variable "vm_cpu" {
  type        = number
  description = "The number of vCPU allocated to the virtual machine"
  default     = 2
  nullable    = false
}

variable "vm_ram" {
  type        = number
  description = "The amount of RAM allocated to the virtual machine"
  default     = 4
  nullable    = false
}

variable "vm_main_disk_size" {
  type        = number
  description = "The amount of main (os) disk (in GiB) allocated to the virtual machine"
  default     = 32
  nullable    = false
}

variable "vm_data_disks" {
  # If pool does not exist, null is passed.
  # If no additional disks need to be created, an empty array ([]) is passed.
  type = list(object({
    name : string
    size : number
    pool : string
  }))
  description = "Additional data disks attached to the virtual machine"
}

variable "vm_mac" {
  type        = string
  description = "The MAC address of the virtual machine"
}

variable "vm_ip" {
  type        = string
  description = "The IP address of the virtual machine"
}
