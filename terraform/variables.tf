variable "action" {
  type    = string
  default = null
}

variable "config_path" {
  type        = string
  description = "Path to the cluster's YAML configuration file."
  default     = "../tkk.yaml"
}
