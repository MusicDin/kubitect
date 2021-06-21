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

  validation {
    condition     = can(regex("^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(.([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])){3}/([1-9]|[1-2][0-9]|3[0-2])$", var.network_cidr))
    error_message = "Invalid network CIDR."
  }
}
