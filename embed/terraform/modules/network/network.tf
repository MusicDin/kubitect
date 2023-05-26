#================================
# Network
#================================

resource "libvirt_network" "network" {
  name      = var.network_name
  mode      = var.network_mode
  bridge    = var.network_bridge
  autostart = true

  addresses = [
    var.network_cidr4, 
    var.network_cidr6,
  ]

  dns {
    enabled    = true
    local_only = false
  }

  dhcp {
    enabled = true
  }
}