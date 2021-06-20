#================================
# Network
#================================

resource "libvirt_network" "network" {
  name      = var.network_name
  mode      = var.network_mode
  bridge    = var.network_bridge
  addresses = [ var.network_cidr ]
  autostart = true

  dns {
    enabled    = true
    local_only = false
  }

  dhcp {
    enabled = true
  }
}