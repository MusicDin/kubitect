version: 2
ethernets:
  ${ network_interface }:
    dhcp4: %{ if vm_cidr4 == "" }true%{ else }false%{ endif }
    dhcp6: %{ if vm_cidr6 == "" }true%{ else }false%{ endif }
    addresses: [${ join(", ", compact([vm_cidr4, vm_cidr6])) }]
    # TODO: Custom gateway IPs
    nameservers:
      addresses: [${ vm_dns_list }]
