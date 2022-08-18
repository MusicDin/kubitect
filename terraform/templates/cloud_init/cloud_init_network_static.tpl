renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: false
    dhcp6: false
    addresses: [${vm_cidr}]
    gateway4: ${network_gateway}
    nameservers:
      addresses: [${vm_dns_list}]
%{ for interface in vm_extra_bridges ~}
  ${interface.networkInterface}:
    addresses: [${interface.ipCidr}]
%{ endfor ~}
