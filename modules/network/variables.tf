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

variable "network_mode" {
  type        = string
  description = "Network mode"
}

variable "network_bridge" {
  type        = string
  description = "Network (virtual) bridge"
}

variable "network_mac" {
  type        = string
  description = "Network MAC address"
}

variable "network_gateway" {
  type        = string
  description = "Network gateway IP address"
}

variable "network_mask_bits" {
  type        = number
  description = "Network mask bits"
}

variable "network_dhcp_ip_start" {
  type        = string
  description = "DHCP IP range start"
}

variable "network_dhcp_ip_end" {
  type        = string
  description = "DHCP IP range end"
}
