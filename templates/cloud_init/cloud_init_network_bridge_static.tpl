renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: false
    dhcp6: false
bridges:
  ${network_bridge}:
    interfaces: [${network_interface}]
    addresses: [${vm_cidr}]
    gateway4: ${network_gateway}
    nameservers:
      addresses: [${network_gateway}, 8.8.8.8]
    dhcp4: false
    dhcp6: false
