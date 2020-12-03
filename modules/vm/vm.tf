# Creates volume for new virtual machine#
resource "libvirt_volume" "vm_volume" {
  name             = "${var.vm_name_prefix}-${var.vm_type}-${var.vm_index}.qcow2"
  pool             = var.resource_pool_name
  base_volume_id   = var.base_volume_id
  size             = var.vm_storage
  format           = "qcow2"
}

# Creates virtual machine
resource "libvirt_domain" "vm_domain" {

  # General configuration #
  name   = "${var.vm_name_prefix}-${var.vm_type}-${var.vm_index}"
  vcpu   = var.vm_cpu
  memory = var.vm_ram
  autostart = true

  cloudinit = var.cloud_init_id

  # Network configuration #
  network_interface {
    network_name   = var.vm_network_name
    mac            = var.vm_mac
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
       "while ! grep \"Cloud-init .* finished\" /var/log/cloud-init.log; do echo \"$(date -Ins) Waiting for cloud-init to finish\"; sleep 2; done"
    ]
  }

}


# Takes care of removing worker from cluster's configuration #
resource "null_resource" "remove_worker" {

  count = var.vm_type == "worker" ? 1 : 0

  triggers = {
    vm_user            = var.vm_user
    vm_ssh_private_key = var.vm_ssh_private_key
    vm_name_prefix     = var.vm_name_prefix
    vm_index           = var.vm_index
  }

  provisioner "local-exec" {
    when    = destroy
    command = "cd ansible/kubespray && virtualenv venv && source venv/bin/activate && pip install -r requirements.txt && ansible-playbook -i ../../config/hosts.ini -b --user=${self.triggers.vm_user} --private-key=${self.triggers.vm_ssh_private_key} -e \"node=$VM_NAME delete_nodes_confirmation=yes\" -v remove-node.yml"

    environment = {
      VM_NAME = "${self.triggers.vm_name_prefix}-worker-${count.index}"
    }

    on_failure = continue
  }

  provisioner "local-exec" {
    when = destroy
    command = "sed 's/${self.triggers.vm_name_prefix}-worker-{vm_index}$//' config/hosts.ini"
    on_failure = continue
  }

}
