version: 2
ethernets:
  ${network_interface}:
    dhcp4: true
    dhcp6: false
    addresses: [${vm_cidr}]
    gateway4: ${network_gateway}
    nameservers:
      addresses: [${vm_dns_list}]
