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

variable "network_cidr4" {
  type        = string
  description = "Network CIDR (v4)"
}

variable "network_cidr6" {
  type        = string
  description = "Network CIDR (v6)"
  nullable    = true
}