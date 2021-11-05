# -------------------------- #
# Outputs                    #
# -------------------------- #

output "network_id" {
  value       = libvirt_network.network.id
  description = "Generated network id"
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

variable "network_cidr" {
  type        = string
  description = "Network CIDR"
}
