output "vm_info" {
  value = {
    id   = var.vm_id
    type = var.vm_type
    name = libvirt_domain.vm_domain.name,
    ip   = try(libvirt_domain.vm_domain.network_interface.0.addresses.0, null)
    # TODO: Ensure IP list is longer then 2
    ip6  = try(libvirt_domain.vm_domain.network_interface.0.addresses.1, null)
    ips  = libvirt_domain.vm_domain.network_interface.0.addresses
    dataDisks = [
      for disk in var.vm_data_disks : {
        name = disk.name
        size = disk.size
        pool = disk.pool == null ? "main" : disk.pool
    }]
  }
  description = "VM's info"
}