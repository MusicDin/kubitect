renderer: networkd
version: 2
ethernets:
  ${network_interface}:
    dhcp4: true
    dhcp6: false
    nameservers:
      addresses: [${vm_dns_list}]
%{ for interface in vm_extra_bridges ~}
  ${interface.networkInterface}:
    addresses: [${interface.ipCidr}]
%{ endfor ~}
