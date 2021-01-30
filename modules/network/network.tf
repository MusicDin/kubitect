#================================
# Network template
#================================

# Populate network config template file #
data "template_file" "network-config-tpl" {
  template = file("templates/network_config.tpl")

  vars = {
    network_name           = var.network_name
    network_forward_mode   = var.network_forward_mode
    network_virtual_bridge = var.network_virtual_bridge
    network_mac            = var.network_mac
    network_gateway        = var.network_gateway
    network_mask           = var.network_mask
    network_dhcp_ip_start  = var.network_dhcp_ip_start
    network_dhcp_ip_end    = var.network_dhcp_ip_end
  }
}

# Create network config file from template #
resource "local_file" "network-config-file" {
  content  = data.template_file.network-config-tpl.rendered
  filename = "config/network_config.xml"
}


#================================
# Network
#================================

# Let terraform manage the lifecycle of the network #
resource "null_resource" "network" {

  # Define triggers for destroy-time provisioners
  triggers = {
    libvirt_provider_uri = var.libvirt_provider_uri
    network_name         = var.network_name
  }

  # On terraform apply - Create network #
  provisioner "local-exec" {
    command     = "virsh --connect ${var.libvirt_provider_uri} net-define config/network_config.xml && virsh --connect ${var.libvirt_provider_uri} net-autostart ${var.network_name} && virsh --connect ${var.libvirt_provider_uri} net-start ${var.network_name}"
    interpreter = ["/bin/bash", "-c"]
  }

  # On terraform destroy - Destroy and undefine network #
  provisioner "local-exec" {
    when       = destroy
    command    = "virsh --connect ${self.triggers.libvirt_provider_uri} net-destroy ${self.triggers.network_name} && virsh --connect ${self.triggers.libvirt_provider_uri} net-undefine ${self.triggers.network_name}"
    on_failure = continue
  }

  # In order to create network configuration, config file must be first created #
  depends_on = [local_file.network-config-file]
}

# Assigns static IP addresses to load balancer VM depending on their MAC address #
resource "null_resource" "lb-static-ips" {

  for_each = var.vm_lb_macs_ips

  # On terraform apply - Add host
  provisioner "local-exec" {
    command     = "virsh --connect ${var.libvirt_provider_uri} net-update ${var.network_name} add ip-dhcp-host \"<host mac='${each.key}' ip='${each.value}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [null_resource.network]
}

# Assigns static IP addresses to master node VMs depending on their MAC address #
resource "null_resource" "master-static-ips" {

  for_each = var.vm_master_macs_ips

  # On terraform apply - Add hosts
  provisioner "local-exec" {
    command     = "virsh --connect ${var.libvirt_provider_uri} net-update ${var.network_name} add ip-dhcp-host \"<host mac='${each.key}' ip='${each.value}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [null_resource.network]
}

# Assigns static IP addresses to worker node VMs depending on their MAC address #
resource "null_resource" "worker-static-ips" {

  for_each = var.vm_worker_macs_ips

  # On terraform apply - Add hosts
  provisioner "local-exec" {
    command     = "virsh --connect ${var.libvirt_provider_uri} net-update ${var.network_name} add ip-dhcp-host \"<host mac='${each.key}' ip='${each.value}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [null_resource.network]
}
