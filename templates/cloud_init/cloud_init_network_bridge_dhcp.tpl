renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: false
    dhcp6: false
bridges:
  ${network_bridge}:
    interfaces: [${network_interface}]
    nameservers:
      addresses: [${vm_dns_list}]
    dhcp4: true
    dhcp6: false
