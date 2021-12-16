renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: true
    nameservers:
      addresses: [${network_dns_list}]
