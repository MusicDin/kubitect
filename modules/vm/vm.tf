# Creates volume for new virtual machine #
resource "libvirt_volume" "vm_volume" {
  name           = "${var.vm_name}.qcow2"
  pool           = var.resource_pool_name
  base_volume_id = var.base_volume_id
  size           = var.vm_storage
  format         = "qcow2"
}

# Creates virtual machine #
resource "libvirt_domain" "vm_domain" {

  # General configuration #
  name      = var.vm_name
  vcpu      = var.vm_cpu
  memory    = var.vm_ram
  autostart = true

  cloudinit = var.cloud_init_id

  # Network configuration #
  network_interface {
    network_id     = var.network_id
    mac            = var.vm_mac
    addresses      = var.vm_ip == null ? null : [var.vm_ip]
    wait_for_lease = true
  }

  # Storage configuration #
  disk {
    volume_id = libvirt_volume.vm_volume.id
  }

  console {
    type        = "pty"
    target_type = "serial"
    target_port = "0"
  }

  console {
    type        = "pty"
    target_type = "virtio"
    target_port = "1"
  }

  graphics {
    type        = "spice"
    listen_type = "address"
    autoport    = true
  }

  # Connect to VM using SSH and wait until cloud-init finishes tasks #
  provisioner "remote-exec" {

    connection {
      host        = self.network_interface.0.addresses.0
      type        = "ssh"
      user        = var.vm_user
      private_key = file(var.vm_ssh_private_key)
    }

    inline = [
      "while ! sudo grep \"Cloud-init .* finished\" /var/log/cloud-init.log; do echo \"$(date -Ins) Waiting for cloud-init to finish\"; sleep 2; done"
    ]
  }
}

# Remove static IP address from network after destruction #
resource "null_resource" "remove_static_ip" {

  triggers = {
    libvirt_provider_uri = var.libvirt_provider_uri
    network_id           = libvirt_domain.vm_domain.network_interface.0.network_id
    vm_mac               = libvirt_domain.vm_domain.network_interface.0.mac
    vm_ip                = libvirt_domain.vm_domain.network_interface.0.addresses.0
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<-EOF
              virsh \
              --connect $URI \
              net-update $NETWORK_ID \
              delete ip-dhcp-host "<host mac='$VM_MAC' ip='$VM_IP'/>" \
              --live \
              --config \
              --parent-index 0
              EOF

    environment = {
      URI        = self.triggers.libvirt_provider_uri
      NETWORK_ID = self.triggers.network_id
      VM_MAC     = self.triggers.vm_mac
      VM_IP      = self.triggers.vm_ip
    }

    on_failure = continue
  }
}

# Adds VM's SSH key to known hosts #
resource "null_resource" "ssh_known_hosts" {

  count = var.vm_ssh_known_hosts ? 1 : 0

  triggers = {
    vm_ip = libvirt_domain.vm_domain.network_interface.0.addresses.0
  }

  provisioner "local-exec" {
    command = "sh ./scripts/filelock-exec.sh \"touch ~/.ssh/known_hosts && ssh-keygen -R ${libvirt_domain.vm_domain.network_interface.0.addresses.0} && ssh-keyscan -t rsa ${libvirt_domain.vm_domain.network_interface.0.addresses.0} | tee -a ~/.ssh/known_hosts && rm -f ~/.ssh/known_hosts.old\""
  }

  provisioner "local-exec" {
    when       = destroy
    command    = "sh ./scripts/filelock-exec.sh \"ssh-keygen -R ${self.triggers.vm_ip}\""
    on_failure = continue
  }
}
