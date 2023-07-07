#================================
# Cloud-init
#================================

# Read SSH public key to inject it into cloud-init template. #
data "local_file" "ssh_public_key" {
  filename = "${var.vm_ssh_private_key}.pub"
}

# Initializes cloud-init disk for user data #
resource "libvirt_cloudinit_disk" "cloud_init" {
  name = "${var.vm_name}-cloud-init.iso"
  pool = var.main_resource_pool_name

  user_data = templatefile("./templates/cloud_init/cloud_init.tpl", {
    hostname       = var.vm_name
    user           = var.vm_user
    update         = var.vm_update
    ssh_public_key = data.local_file.ssh_public_key.content
  })

  network_config = templatefile(var.vm_ip != null
    ? "./templates/cloud_init/cloud_init_network_static.tpl"
    : "./templates/cloud_init/cloud_init_network_dhcp.tpl"
    , {
      network_interface = var.vm_network_interface
      network_bridge    = var.network_bridge
      network_gateway   = var.network_gateway
      vm_dns_list       = length(var.vm_dns) == 0 ? var.network_gateway : join(", ", var.vm_dns)
      vm_cidr           = var.vm_ip == null ? "" : "${var.vm_ip}/${split("/", var.network_cidr)[1]}"
  })
}

#================================
# VM
#================================

# Creates volume for new virtual machine #
resource "libvirt_volume" "vm_main_disk" {
  name           = "${var.vm_name}-main-disk"
  pool           = var.main_resource_pool_name
  base_volume_id = var.base_volume_id
  size           = var.vm_main_disk_size * pow(1024, 3) # GiB -> B
  format         = "qcow2"
}

# Creates volume for new virtual machine #
resource "libvirt_volume" "vm_data_disks" {
  for_each = { for disk in var.vm_data_disks : disk.name => disk }

  name = "${var.vm_name}-${each.key}-data-disk"
  pool = (each.value.pool == "main"
    ? var.main_resource_pool_name
    : "${var.cluster_name}-${each.value.pool}-data-resource-pool"
  )
  size = each.value.size * pow(1024, 3) # GiB -> B
}

# Creates virtual machine #
resource "libvirt_domain" "vm_domain" {

  # General configuration #
  name      = var.vm_name
  vcpu      = var.vm_cpu
  memory    = var.vm_ram * 1024 # GiB -> MiB
  autostart = true

  cpu {
    mode = var.vm_cpuMode
  }

  cloudinit = libvirt_cloudinit_disk.cloud_init.id

  qemu_agent = (var.network_mode == "bridge")

  # Network configuration #
  network_interface {
    network_id     = var.network_id
    mac            = var.vm_mac
    bridge         = var.network_bridge
    addresses      = var.network_mode == "nat" && var.vm_ip != null ? [var.vm_ip] : null
    wait_for_lease = true
  }

  # Storage configuration #
  dynamic "disk" {
    for_each = concat(
      [{ "id" : libvirt_volume.vm_main_disk.id }],
      [for disk in libvirt_volume.vm_data_disks : { "id" : disk.id }]
    )
    content {
      volume_id = disk.value.id
    }
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
      "while ! sudo grep \"Cloud-init .* finished\" /var/log/cloud-init.log; do echo \"Waiting for cloud-init to finish...\"; sleep 2; done"
    ]
  }
}

#================================
# SSH known hosts
#================================

# Adds VM's SSH key to known hosts #
resource "null_resource" "ssh_known_hosts" {

  count = var.vm_ssh_known_hosts ? 1 : 0

  triggers = {
    vm_ip = libvirt_domain.vm_domain.network_interface.0.addresses.0
  }

  provisioner "local-exec" {
    command = <<-EOF
      sh ./scripts/filelock-exec.sh \
        "touch $HOME/.ssh/known_hosts \
        && ssh-keygen -R $VM_IP \
        && ssh-keyscan -t rsa $VM_IP | tee -a $HOME/.ssh/known_hosts \
        && rm -f $HOME/.ssh/known_hosts.old"
    EOF

    environment = {
      HOME  = pathexpand("~")
      VM_IP = libvirt_domain.vm_domain.network_interface.0.addresses.0
    }

    quiet = true
  }

  provisioner "local-exec" {
    when       = destroy
    command    = "sh ./scripts/filelock-exec.sh \"ssh-keygen -R ${self.triggers.vm_ip}\""
    on_failure = continue
  }
}
