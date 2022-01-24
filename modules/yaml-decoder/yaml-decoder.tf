locals {
  yaml_config = yamldecode(var.yaml_config)
}

output "cluster_name" {
  value       = local.yaml_config.cluster.name
  description = "Cluster name"
}

output "name" {
  value       = ""
  sensitive   = true
  description = "description"
  depends_on  = []
}
