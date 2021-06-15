#cloud-config
users:
  - name: ${user}
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    shell: /bin/bash
    ssh_authorized_keys:
      - ${ssh_public_key}

packages:
  - qemu-guest-agent

runcmd:
  - [ systemctl, enable, qemu-guest-agent.service ]
  - [ systemctl, start, qemu-guest-agent.service ]