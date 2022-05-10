#======================================================================================
# Virtual machine configuration
#======================================================================================

variable "lb_vip" {
  type        = string
  description = "Load balancer virtual IP address (VIP)"
  default     = ""
  nullable    = false
}

variable "vm_user" {
  type        = string
  description = "SSH user for VMs"
}

variable "vm_ssh" {
  type = object({
    privateKeyPath = string
  })
  description = "Location of private ssh key for VMs"
}

variable "vm_os" {
  type = object({
    distro           = string
    source           = string
    networkInterface = string
  })
  description = "Operating system (os) information."
}

#======================================================================================
# Virtual machine instances
#======================================================================================

variable "worker_nodes" {
  type = list(object({
    id   = number
    name = string
    ip   = string
  }))
  description = "Worker nodes info"
}

variable "master_nodes" {
  type = list(object({
    id   = number
    name = string
    ip   = string
  }))
  description = "Master nodes info"

  validation {
    condition = (
      length(var.master_nodes) % 2 != 0
    )
    error_message = "Master nodes configuration is incorrect. Make sure that: \n - number of master nodes is odd (not divisible by 2)."
  }
}

variable "lb_nodes" {
  type = list(object({
    id   = number
    name = string
    ip   = string
  }))
  description = "Load balancers info"
}
