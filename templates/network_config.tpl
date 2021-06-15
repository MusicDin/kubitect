<network>
  <name>${network_name}</name>
  <forward mode='${network_mode}'/>
  <bridge name='${network_bridge}' stp='on' delay='0'/>
  <mac address='${network_mac}'/>
  <ip address='${network_gateway}' netmask='${network_mask}'>
    <dhcp>
      <range start='${network_dhcp_ip_start}' end='${network_dhcp_ip_end}'/>
    </dhcp>
  </ip>
</network>

