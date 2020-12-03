<network connections='1'>
  <name>${network_name}</name>
  <forward mode='nat'>
    <nat>
      <port start='${network_nat_port_start}' end='${network_nat_port_end}'/>
    </nat>
  </forward>
  <bridge name='${network_virtual_bridge}' stp='on' delay='0'/>
  <mac address='${network_mac}'/>
  <ip address='${network_gateway}' netmask='${network_mask}'>
    <dhcp>
      <range start='${network_dhcp_ip_start}' end='${network_dhcp_ip_end}'/>
    </dhcp>
  </ip>
</network>

