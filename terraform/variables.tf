variable "action" {
  type    = string
  default = null
}

variable "config_path" {
  type        = string
  description = "Path to the cluster's configuration file."
  default     = "./variables.yaml"

  validation {
    condition     = fileexists(var.config_path)
    error_message = "Cluster configuration file does not exist!"
  }
}

variable "infra_config_path" {
  type        = string
  description = "Path to the infrastructure configuration file."
  default     = "../config/infrastructure.yaml"
}
