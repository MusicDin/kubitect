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

package_update: true
package_upgrade: ${update}

packages:
  - qemu-guest-agent

bootcmd:
  # Disable qemu-guest-agent to prevent reporting IP addresses
  # before cloud-init has configured the network.
  - cloud-init-per once disable-qemu-ga systemctl stop qemu-guest-agent.service

runcmd:
  - [ systemctl, enable, qemu-guest-agent.service ]
  - [ systemctl, start, qemu-guest-agent.service ]
