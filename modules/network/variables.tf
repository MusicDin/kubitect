# -------------------------- #
# Network specific variables #
# -------------------------- #

variable "network_name" {
  type        = string
  description = "Network name"
}

variable "network_mac" {
  type        = string
  description = "Network MAC address"
}

variable "network_gateway" {
  type        = string
  description = "Network gateway IP address"
}

variable "network_mask" {
  type        = string
  description = "Network mask"
}

variable "network_nat_port_start" {
  type        = string
  description = "NAT (Network Address Translation) port start (from port)"
}

variable "network_nat_port_end" {
  type        = string
  description = "NAT port end (to port)"
}

variable "network_dhcp_ip_start" {
  type        = string
  description = "DHCP IP range start"
}

variable "network_dhcp_ip_end" {
  type        = string
  description = "DHCP IP range end"
}

variable "network_virtual_brdige" {
  type        = string
  description = "Network virtual bridge"
}

# -------------------------- #
# VM specific variables      #
# -------------------------- #

variable "vm_lb_macs_ips" {
  type        = map(string)
  description = "Map of MACs and IPs for load balancers"
}

variable "vm_master_macs_ips" {
  type        = map(string)
  description = "Map of MACs and IPs for master nodes"
}

variable "vm_worker_macs_ips" {
  type        = map(string)
  description = "Map of MACs and IPs for worker nodes"
}
