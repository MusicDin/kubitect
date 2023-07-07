#======================================================================================
# Virtual machine configuration
#======================================================================================

variable "lb_vip" {
  type        = string
  description = "Load balancer virtual IP address (VIP)"
}

#======================================================================================
# Virtual machine instances
#======================================================================================

variable "worker_nodes" {
  type = list(object({
    id   = string
    name = string
    ip   = string
    ip6  = string
    ips  = list(string)
    dataDisks = list(object({
      name = string
      size = number
      pool = string
    }))
  }))
  description = "Worker nodes info"
}

variable "master_nodes" {
  type = list(object({
    id   = string
    name = string
    ip   = string
    ip6  = string
    ips  = list(string)
    dataDisks = list(object({
      name = string
      size = number
      pool = string
    }))
  }))
  description = "Master nodes info"
}

variable "lb_nodes" {
  type = list(object({
    id   = string
    name = string
    ip   = string
    ip6  = string
    ips  = list(string)
  }))
  description = "Load balancers info"
}
