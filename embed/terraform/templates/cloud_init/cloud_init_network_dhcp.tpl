version: 2
ethernets:
  ${network_interface}:
    dhcp4: true
    dhcp6: false
    nameservers:
      addresses: [${vm_dns_list}]
