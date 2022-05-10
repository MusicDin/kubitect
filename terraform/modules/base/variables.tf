variable "config" {
  type        = any
  description = "Cluster configuration file."
  nullable    = false
}

variable "infra_config" {
  type        = any
  description = "Infrastructure config created first initialization of the cluster."
}

variable "defaults_config" {
  type        = any
  description = "Configuration file containing default values for various fields, such as OS distributions."
  nullable    = false
}
