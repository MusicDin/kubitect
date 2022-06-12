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
      addresses: [${vm_dns_list}]
    dhcp4: false
    dhcp6: false
