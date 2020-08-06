# libvirt.tf

#======================================================================================
# Libvirt provider
#======================================================================================

provider "libvirt" {
  uri = "qemu:///system"
}

#======================================================================================
# Template files
#======================================================================================

#================================  
# Network template
#================================

# Network configuration #
data "template_file" "network_config" {
  template = file("network/k8s_network_config.cfg")
}

#================================  
# Cloud config template
#================================

# Public ssh key for vm (it is directly injected in cloud-init configuration) #
data "template_file" "public_ssh_key" {
  template = file("${var.vm_ssh_private_key}.pub")
}

# User data (cloud-init) configuration #
data "template_file" "user_data" {
  template = file("templates/cloud_init.tpl")
  
  vars = {
    user = var.vm_user
    ssh_public_key = data.template_file.public_ssh_key.rendered
  }
}

#================================  
# Kubespray templates 
#================================

# Kubespray all.yml template (Currently supports only 1 load balancer) #
data "template_file" "kubespray_all" {
  
  template = file("templates/kubespray_all.tpl")
  
  vars = {
    loadbalancer_apiserver = var.vm_haproxy_vip
  }
}

# Kubespray k8s-cluster.yml template #
data "template_file" "kubespray_k8s_cluster" {
  
  template = file("templates/kubespray_k8s_cluster.tpl")  
  
  vars = {
    kube_version = var.k8s_version
    kube_network_plugin = var.k8s_network_plugin
    dns_mode = var.k8s_dns_mode
  }
}

# Load balancer hostname and ip list template #
data "template_file" "lb_hosts" {

  count = length(var.vm_lb_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-lb-${count.index}"
    host_ip = lookup(var.vm_lb_ips, count.index)
  }
}

# Master hostname and ip list template #
data "template_file" "master_hosts" {

  count = length(var.vm_master_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-master-${count.index}"
    host_ip = lookup(var.vm_master_ips, count.index)
  }
}

# Worker hostname and ip list template #
data "template_file" "worker_hosts" {

  count = length(var.vm_worker_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-worker-${count.index}"
    host_ip = lookup(var.vm_worker_ips, count.index)
  }
}

# Load balancer hostname list template #
data "template_file" "lb_hosts_only" {

  count = length(var.vm_lb_ips)

  template = file("templates/ansible_hosts_list.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-lb-${count.index}"
  }
}

# Kubespray master hostname list template #
data "template_file" "master_hosts_only" {

  count = length(var.vm_master_ips)

  template = file("templates/ansible_hosts_list.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-master-${count.index}"
  }
}

# Kubespray worker hostname and ip list template #
data "template_file" "worker_hosts_only" {

  count = length(var.vm_worker_ips)

  template = file("templates/ansible_hosts_list.tpl")

  vars = {
    hostname = "${var.vm_name_prefix}-worker-${count.index}"
  }
}

# HAProxy template #
data "template_file" "haproxy" {
  template = file("templates/haproxy.tpl")

  vars = {
    bind_ip = var.vm_haproxy_vip
  }
}

# HAProxy server backend template #
data "template_file" "haproxy_backend" {
  
  count = length(var.vm_master_ips)
   
  template = file("templates/haproxy_backend.tpl")

  vars = {
    prefix_server     = var.vm_name_prefix
    backend_server_ip = lookup(var.vm_master_ips, count.index)
    count             = "${count.index}"
  }
}

# Keepalived master template #
data "template_file" "keepalived_master" {
  template = file("templates/keepalived_master.tpl")

  vars = {
    virtual_ip = var.vm_haproxy_vip
  }
}

# Keepalived slave template #
data "template_file" "keepalived_slave" {
  template = file("templates/keepalived_slave.tpl")

  vars = {
    virtual_ip = var.vm_haproxy_vip
  }
}


#======================================================================================
# Local files
#======================================================================================

# Create user data (cloud-init) file from template #
resource "local_file" "user_data" {
  content = data.template_file.user_data.rendered
  filename = "user_data/cloud_init.cfg"
}

# Create Kubespray all.yml configuration file from template #
resource "local_file" "kubespray_all" {
  content  = data.template_file.kubespray_all.rendered
  filename = "config/group_vars/all.yml" 
}

# Create Kubespray k8s-cluster.yml configuration file from template #
resource "local_file" "kubespray_k8s_cluster" {
  content  = data.template_file.kubespray_k8s_cluster.rendered
  filename = "config/group_vars/k8s-cluster.yml"
}

# Create Kubespray hosts.ini configuration file from template #
resource "local_file" "kubespray_hosts" {
  content  = "[all]\n${join("", data.template_file.lb_hosts.*.rendered)}${join("", data.template_file.master_hosts.*.rendered)}${join("", data.template_file.worker_hosts.*.rendered)}\n[haproxy]\n${join("", data.template_file.lb_hosts_only.*.rendered)}\n[kube-master]\n${join("", data.template_file.master_hosts_only.*.rendered)}\n[etcd]\n${join("", data.template_file.master_hosts_only.*.rendered)}\n[kube-node]\n${join("", data.template_file.worker_hosts_only.*.rendered)}\n[k8s-cluster:children]\nkube-master\nkube-node"
  filename = "config/hosts.ini"
}

# Create HAProxy configuration file from template #
resource "local_file" "haproxy" {
  content  = "${data.template_file.haproxy.rendered}${join("", data.template_file.haproxy_backend.*.rendered)}"
  filename = "config/haproxy.cfg"
}

# Create keepalived master configuration file from template #
resource "local_file" "keepalived_master" {
  content  = data.template_file.keepalived_master.rendered
  filename = "config/keepalived-master.cfg"
}

# Create keepalived slave configuration file from template #
resource "local_file" "keepalived_slave" {
  content  = data.template_file.keepalived_slave.rendered
  filename = "config/keepalived-slave.cfg"
}

#======================================================================================
# Null resources
#======================================================================================

#================================  
# Network 
#================================

# Let terraform manage the lifecycle of the network #
resource "null_resource" "network" {

  # On terraform apply - Create network #
  provisioner "local-exec" {
    command     = "virsh net-define network/network_config.xml && virsh net-autostart ${var.vm_network_name} && virsh net-start ${var.vm_network_name}"
    interpreter = ["/bin/bash", "-c"]
  }

  # On terraform destroy - Remove network #
  provisioner "local-exec" {
    when       = destroy
    command    = "virsh net-undefine ${var.vm_network_name} && virsh net-destroy ${var.vm_network_name}"
    on_failure = continue
  }
}

# Assigns static IP addresses to master node VMs depending on their MAC address #
resource "null_resource" "lb-static-ips" {

  count = length(var.vm_lb_macs)
  
  # On terraform apply - Add hosts
  provisioner "local-exec" {
    command     = "virsh net-update ${var.vm_network_name} add ip-dhcp-host \"<host mac='${var.vm_lb_macs[count.index]}' ip='${var.vm_lb_ips[count.index]}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
    #on_failure = continue
  }

  depends_on = [null_resource.network]
}

# Assigns static IP addresses to master node VMs depending on their MAC address #
resource "null_resource" "master-static-ips" {

  count = length(var.vm_master_macs)
  
  # On terraform apply - Add hosts
  provisioner "local-exec" {
    command     = "virsh net-update ${var.vm_network_name} add ip-dhcp-host \"<host mac='${var.vm_master_macs[count.index]}' ip='${var.vm_master_ips[count.index]}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [null_resource.network]
}

# Assigns static IP addresses to worker node VMs depending on their MAC address #
resource "null_resource" "worker-static-ips" {

  count = length(var.vm_worker_macs)
  
  # On terraform apply - Add hosts
  provisioner "local-exec" {
    command     = "virsh net-update ${var.vm_network_name} add ip-dhcp-host \"<host mac='${var.vm_worker_macs[count.index]}' ip='${var.vm_worker_ips[count.index]}'/>\" --live --config"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [null_resource.network]
}

#================================  
# Kubespray      
#================================  

# Modify permissions on config directory #
resource "null_resource" "config_permissions" {
  provisioner "local-exec" {
    command = "chmod -R 700 config"
  }

  depends_on = [
    local_file.kubespray_hosts, 
    local_file.kubespray_all, 
    local_file.kubespray_k8s_cluster, 
    null_resource.kubespray_download
  ]
}

# Clone Kubespray repository #
resource "null_resource" "kubespray_download" {
  provisioner "local-exec" {
    command = "cd ansible && rm -rf kubespray && git clone --branch ${var.k8s_kubespray_version} ${var.k8s_kubespray_url}"
  }
}

# Execute create Kubernetes HAProxy playbook #
resource "null_resource" "haproxy_install" {
  count = var.action == "create" ? 1: 0

#  provisioner "remote-exec" {
#    inline = [
#      "cloud-init status --wait"
#    ]
#  }
  
  provisioner "local-exec" {
    command = "cd ansible/haproxy && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -T 3000 -v haproxy.yml"
  }
  
  depends_on = [
    local_file.kubespray_hosts,
    local_file.haproxy,
    libvirt_domain.lb_nodes
  ]
}

# Create Kubespray Ansible playbook #
resource "null_resource" "kubespray_create" {
  
  count = var.action == "create" ? 1 : 0

  provisioner "local-exec" {
    command = "cd ansible/kubespray && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" -T 3000 -v cluster.yml"
  }

  depends_on = [
    local_file.kubespray_hosts, 
    local_file.kubespray_all, 
    local_file.kubespray_k8s_cluster, 
    null_resource.haproxy_install,
    null_resource.kubespray_download, 
    null_resource.worker-static-ips,
    null_resource.master-static-ips,
    null_resource.lb-static-ips,
    libvirt_domain.master_nodes, 
    libvirt_domain.worker_nodes, 
    libvirt_domain.lb_nodes
  ]
}

# Execte scale Kubespray Ansible playbook #
resource "null_resource" "kubespray_add" {
  count = var.action == "add_worker" ? 1 : 0

  provisioner "local-exec" {
    command = "cd ansible/kubespray && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" -T 3000 -v scale.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install,
    libvirt_domain.master_nodes,
    libvirt_domain.worker_nodes,
    libvirt_domain.lb_nodes
  ]
}

# Execute upgrade Kubespray Ansible playbook #
resource "null_resource" "kubespray_upgrade" {
  count = var.action == "upgrade" ? 1 : 0

  triggers = {
    ts = "$(timestamp())"
  }
  
  # Deletes old Kubespray and installs new one #
  provisioner "local-exec" {
    command = "cd ansible && rm -rf kubespray && git clone --branch ${var.k8s_kubespray_version} ${var.k8s_kubespray_url}"
  }

  provisioner "local-exec" {
    command = "cd ansible/kubespray && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" -T 3000 -v upgrade-cluster.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,     
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install,
    libvirt_domain.master_nodes,
    libvirt_domain.worker_nodes,
    libvirt_domain.lb_nodes
  ]
}


# Create the local admin.conf kubectl configuration file #
resource "null_resource" "kubectl_configuration" {
  provisioner "local-exec" {
    command = "ansible -i ${lookup(var.vm_master_ips, 0)}, -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -T 3000 -m fetch -a 'src=/etc/kubernetes/admin.conf dest=config/admin.conf flat=yes' all"
  }

#  provisioner "local-exec" {
#    command = "sed 's/lb-apiserver.kubernetes.local/${var.vm_haproxy_vip}/g' config/admin.conf | tee config/admin.conf.new $$ mv config/admin.conf.new config/admin.conf && chmod 700 config/admin.conf"
#  }

#  provisioner "local-exec" {
#    command = "chmod 600 config/admin.conf"
#  }

  depends_on = [null_resource.kubespray_create]
}

#======================================================================================
# Libvirt resources
#======================================================================================

# Create a resource pool for Kubernetes VMs #
resource "libvirt_pool" "resource_pool_base" {
  name = "${var.libvirt_resource_pool_name}-base"
  type = "dir"
  path = "/var/lib/libvirt/pools/${var.libvirt_resource_pool_name}-base"
}

# Create a resource pool for Kubernetes VMs #
resource "libvirt_pool" "resource_pool" {
  name = var.libvirt_resource_pool_name
  type = "dir"
  path = "/var/lib/libvirt/pools/${var.libvirt_resource_pool_name}"
}


# Create image volume for each load balancer node #
resource "libvirt_volume" "lb_volumes_base" {
  count  = length(var.vm_lb_ips)
  name   = "${var.vm_name_prefix}-lb-base-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool_base.name
  source = "${path.module}/downloads/${var.vm_image_name}"
  format = "qcow2"

  depends_on = [
    libvirt_pool.resource_pool_base
  ]
}

# Expand load balancer base volume #
resource "libvirt_volume" "lb_volumes" {
  count  = length(var.vm_lb_ips)
  name   = "${var.vm_name_prefix}-lb-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool.name
  base_volume_name = "${var.vm_name_prefix}-lb-base-${count.index}.qcow2"
  base_volume_pool = libvirt_pool.resource_pool_base.name
  size   = var.vm_disk_size
  format = "qcow2"

  depends_on = [
    libvirt_volume.lb_volumes_base, 
    libvirt_pool.resource_pool
  ]
}

# Create image volume for each master node #
resource "libvirt_volume" "master_volumes_base" {
  count  = length(var.vm_master_ips)
  name   = "${var.vm_name_prefix}-master-base-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool_base.name
  source = "${path.module}/downloads/${var.vm_image_name}"
  format = "qcow2"

  depends_on = [
    libvirt_pool.resource_pool_base 
  ]
}

# Expand master node base volume #
resource "libvirt_volume" "master_volumes" {
  count  = length(var.vm_master_ips)
  name   = "${var.vm_name_prefix}-master-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool.name
  base_volume_name = "${var.vm_name_prefix}-master-base-${count.index}.qcow2"
  base_volume_pool = libvirt_pool.resource_pool_base.name
  size   = var.vm_disk_size
  format = "qcow2"

  depends_on = [
    libvirt_volume.master_volumes_base, 
    libvirt_pool.resource_pool
  ]
}

# Create image volume for each worker node #
resource "libvirt_volume" "worker_volumes_base" {
  count  = length(var.vm_worker_ips)
  name   = "${var.vm_name_prefix}-worker-base-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool_base.name
  source = "${path.module}/downloads/${var.vm_image_name}"
  format = "qcow2"

  depends_on = [
    libvirt_pool.resource_pool_base 
  ]
}

# Expand worker node base volume #
resource "libvirt_volume" "worker_volumes" {
  count  = length(var.vm_worker_ips)
  name   = "${var.vm_name_prefix}-worker-${count.index}.qcow2"
  pool   = libvirt_pool.resource_pool.name
  base_volume_name = "${var.vm_name_prefix}-worker-base-${count.index}.qcow2"
  base_volume_pool = libvirt_pool.resource_pool_base.name
  size   = var.vm_disk_size
  format = "qcow2"

  depends_on = [
    libvirt_volume.worker_volumes_base, 
    libvirt_pool.resource_pool
  ]
}

# Creates disk for load balancer user_data #
resource "libvirt_cloudinit_disk" "cloud_init" {
  name           = "cloud-init.iso"
  pool           = libvirt_pool.resource_pool.name
  user_data      = data.template_file.user_data.rendered
  network_config = data.template_file.network_config.rendered
}

# Creates load balancer nodes
resource "libvirt_domain" "lb_nodes" {

  count  = length(var.vm_lb_macs)

  name   = "${var.vm_name_prefix}-lb-${count.index}"
  vcpu   = var.vm_lb_cpu
  memory = var.vm_lb_ram
  autostart = true
 
  cloudinit = libvirt_cloudinit_disk.cloud_init.id
 
  network_interface {
    network_name   = var.vm_network_name
    mac            = var.vm_lb_macs[count.index]
    wait_for_lease = true
  }

  disk {
    volume_id = element(libvirt_volume.lb_volumes.*.id, count.index)
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

# Creates master nodes
resource "libvirt_domain" "master_nodes" {

  count  = length(var.vm_master_macs) 

  name   = "${var.vm_name_prefix}-master-${count.index}"
  vcpu   = var.vm_master_cpu
  memory = var.vm_master_ram
  autostart = true
 
  cloudinit = libvirt_cloudinit_disk.cloud_init.id
 
  network_interface {
    network_name   = var.vm_network_name
    mac            = var.vm_master_macs[count.index]
    wait_for_lease = true
  }

  disk {
    volume_id = element(libvirt_volume.master_volumes.*.id, count.index)
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

  depends_on = [
    libvirt_domain.lb_nodes
  ]
}

# Creates worker nodes #
resource "libvirt_domain" "worker_nodes" {

  count  = length(var.vm_worker_macs)

  name   = "${var.vm_name_prefix}-worker-${count.index}"
  vcpu   = var.vm_worker_cpu
  memory = var.vm_worker_ram
  autostart = true
 
  cloudinit = libvirt_cloudinit_disk.cloud_init.id
 
  network_interface {
    network_name   = var.vm_network_name
    mac            = var.vm_worker_macs[count.index]
    wait_for_lease = true
  }

  disk {
    volume_id = element(libvirt_volume.worker_volumes.*.id, count.index)
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

  provisioner "local-exec" {
    when    = destroy
    command = "cd ansible/kubespray && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"node=$VM_NAME delete_nodes_confirmation=yes\" -v remove-node.yml" 

    environment = {
      VM_NAME = "${var.vm_name_prefix}-worker-${count.index}"
    }

    on_failure = continue
  }

  provisioner "local-exec" {
    when = destroy
    command = "sed 's/${var.vm_name_prefix}-worker-[0-9]*$//' config/hosts.ini"
  }
  
  depends_on = [
    libvirt_domain.master_nodes,
    local_file.kubespray_hosts,
    local_file.kubespray_k8s_cluster,
    local_file.kubespray_all
  ]
}

