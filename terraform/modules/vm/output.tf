output "vm_info" {
  value = {
    id   = var.vm_id
    type = var.vm_type
    name = libvirt_domain.vm_domain.name,
    ip   = try(libvirt_domain.vm_domain.network_interface.0.addresses.0, null)
    dataDisks = [
      for disk in var.vm_data_disks : {
        name = disk.name
        size = disk.size
        pool = disk.pool
        dev  = trim(data.local_file.data_disks_mapping[disk.name].content, "\n")
    }]
  }
  description = "VM's info"
}