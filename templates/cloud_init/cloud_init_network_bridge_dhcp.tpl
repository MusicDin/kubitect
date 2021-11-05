version: 2
ethernets:
  ${network_interface}:
    dhcp4: false
    dhcp6: false
bridges:
  ${network_bridge}:
    interfaces: [${network_interface}]
    dhcp4: true
    dhcp6: false
    parameters:
      # STP - Spanning Tree Protocol
      stp: false