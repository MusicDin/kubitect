renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: true
    nameservers:
      addresses: [${vm_dns_list}]
