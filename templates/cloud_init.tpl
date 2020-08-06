#cloud-config
users:
  - name: ${user}
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    shell: /bin/bash
    ssh_authorized_keys:
      - ${ssh_public_key}
# install packages
packages:
  - qemu-guest-agent
