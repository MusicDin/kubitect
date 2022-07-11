output "os" {
  value       = local.os
  description = "Evaluated operating system (os) information."
}

output "ssh" {
  value       = local.ssh
  description = "Evaluated SSH private key location."
}

output "user" {
  value       = local.user
  description = "Evaluated VM user."
}