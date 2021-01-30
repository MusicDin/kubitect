# -------------------------- #
# Virsh specific variables   #
# -------------------------- #

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI"
}

# -------------------------- #
# Network specific variables #
# -------------------------- #

variable "network_name" {
  type        = string
  description = "Network name"
}

variable "network_forward_mode" {
  type        = string
  description = "Network forward mode"
}

variable "network_virtual_bridge" {
  type        = string
  description = "Network virtual bridge"
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

variable "network_dhcp_ip_start" {
  type        = string
  description = "DHCP IP range start"
}

variable "network_dhcp_ip_end" {
  type        = string
  description = "DHCP IP range end"
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
