#cloud-config
preserve_hostname: false
hostname: ${hostname}

users:
  - name: ${user}
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    shell: /bin/bash
    ssh_authorized_keys:
      - ${ssh_public_key}

package_upgrade: ${update}

packages:
  - qemu-guest-agent

runcmd:
  - [ systemctl, start, qemu-guest-agent.service ]
