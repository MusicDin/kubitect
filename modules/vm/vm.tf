# Creates volume for new virtual machine #
resource "libvirt_volume" "vm_volume" {
  name           = "${var.vm_name_prefix}-${var.vm_type}-${var.vm_index}.qcow2"
  pool           = var.resource_pool_name
  base_volume_id = var.base_volume_id
  size           = var.vm_storage
  format         = "qcow2"
}

# Assigns static IP addresses to VM #
resource "null_resource" "vm_static_ip" {

  # Define triggers for on-destroy provisioner
  triggers = {
    libvirt_provider_uri = var.libvirt_provider_uri
    network_name         = var.network_name
    vm_mac               = var.vm_mac
    vm_ip                = var.vm_ip
  }

  provisioner "local-exec" {
    command     = "virsh --connect ${var.libvirt_provider_uri} net-update ${var.network_name} add ip-dhcp-host \"<host mac='${var.vm_mac}' ip='${var.vm_ip}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  provisioner "local-exec" {
    when        = destroy
    command     = "virsh --connect ${self.triggers.libvirt_provider_uri} net-update ${self.triggers.network_name} delete ip-dhcp-host \"<host mac='${self.triggers.vm_mac}' ip='${self.triggers.vm_ip}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
    on_failure  = continue
  }
}

# Creates virtual machine #
resource "libvirt_domain" "vm_domain" {

  # General configuration #
  name      = "${var.vm_name_prefix}-${var.vm_type}-${var.vm_index}"
  vcpu      = var.vm_cpu
  memory    = var.vm_ram
  autostart = true

  cloudinit = var.cloud_init_id

  # Network configuration #
  network_interface {
    network_name   = var.network_name
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
    command = "cd ansible/kubespray && virtualenv venv && . venv/bin/activate && pip install -r requirements.txt && ansible-playbook -i ../../config/hosts.ini -b --user=${self.triggers.vm_user} --private-key=${self.triggers.vm_ssh_private_key} -e \"node=$VM_NAME delete_nodes_confirmation=yes\" -v remove-node.yml"

    environment = {
      VM_NAME = "${self.triggers.vm_name_prefix}-worker-${count.index}"
    }

    on_failure = continue
  }

  provisioner "local-exec" {
    when       = destroy
    command    = "sed -i '/${self.triggers.vm_name_prefix}-worker-${self.triggers.vm_index}$/d' config/hosts.ini"
    on_failure = continue
  }

}

# Adds VM's SSH key to known hosts #
resource "null_resource" "ssh_known_hosts" {

  count = var.vm_ssh_known_hosts == "true" ? 1 : 0

  triggers = {
    vm_ip = var.vm_ip
  }

  provisioner "local-exec" {
    command = "sh ./scripts/filelock-exec.sh \"touch ~/.ssh/known_hosts && ssh-keygen -R ${var.vm_ip} && ssh-keyscan -t rsa ${var.vm_ip} | tee -a ~/.ssh/known_hosts && rm -f ~/.ssh/known_hosts.old\""
  }

  provisioner "local-exec" {
    when       = destroy
    command    = "sh ./scripts/filelock-exec.sh \"ssh-keygen -R ${self.triggers.vm_ip}\""
    on_failure = continue
  }

  depends_on = [libvirt_domain.vm_domain]
}
