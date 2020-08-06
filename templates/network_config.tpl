<network connections='1'>
  <name>k8s-network</name>
  <uuid>2480ea41-97ed-4a51-b356-0c6114ef4a5c</uuid>
  <forward mode='nat'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='virbr1' stp='on' delay='0'/>
  <mac address='52:54:00:4f:e3:99'/>
  <ip address='192.168.113.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.113.2' end='192.168.113.254'/>
    </dhcp>
  </ip>
</network>

