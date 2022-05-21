output "vm_info" {
  value = {
    id     = var.vm_id
    type   = var.vm_type
    labels = var.vm_labels
    name   = libvirt_domain.vm_domain.name,
    ip     = libvirt_domain.vm_domain.network_interface.0.addresses.0
  }
  description = "VM's info containing it's name and an IP address"
}