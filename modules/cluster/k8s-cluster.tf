#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  extra_args = {
    debian = "-T 3000 -v -e 'ansible_become_method=su'"
    ubuntu = "-T 3000 -v"
    centos = "-T 3000 -v"
  }
  default_extra_args = "-T 3000 -v"
}

#================================
# Kubespray templates
#================================

# Kubespray all.yml template (Currently supports only 1 load balancer) #
data "template_file" "kubespray_all" {

  template = file("templates/kubespray_all.tpl")

  vars = {
    loadbalancer_apiserver = var.vm_lb_vip
  }
}

# Kubespray k8s-cluster.yml template #
data "template_file" "kubespray_k8s_cluster" {

  template = file("templates/kubespray_k8s_cluster.tpl")

  vars = {
    kube_version          = var.k8s_version
    kube_network_plugin   = var.k8s_network_plugin
    dns_mode              = var.k8s_dns_mode
    kube_proxy_strict_arp = yamldecode( var.kubespray_custom_addons_enabled == "false"
                                        ? data.template_file.kubespray_addons[0].rendered
                                        : data.template_file.kubespray_custom_addons[0].rendered )["metallb_enabled"]
  }

  # Correct addons template file has to be created before
  # 'metallb_enabled' value can be read from it
  depends_on = [
    data.template_file.kubespray_addons,
    data.template_file.kubespray_custom_addons
  ]
}

# Kubespray addons.yml template #
data "template_file" "kubespray_addons" {

  count = var.kubespray_custom_addons_enabled == "false" ? 1 : 0

  template = file("templates/kubespray_addons.tpl")

  vars = {
    dashboard_enabled = var.k8s_dashboard_enabled
    helm_enabled      = var.helm_enabled
    metallb_enabled   = var.metallb_enabled
    metallb_version   = var.metallb_version
    metallb_port      = var.metallb_port
    metallb_cpu_limit = var.metallb_cpu_limit
    metallb_mem_limit = var.metallb_mem_limit
    metallb_protocol  = var.metallb_protocol
    metallb_ip_range  = var.metallb_ip_range
  }
}

# Kubespray custom addons.yml #
data "template_file" "kubespray_custom_addons" {

  count = var.kubespray_custom_addons_enabled == "true" ? 1 : 0

  template = file(var.kubespray_custom_addons_path)
}

# Load balancer hostname and ip list template #
data "template_file" "lb_hosts" {

  count = length(var.vm_lb_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname    = "${var.vm_name_prefix}-lb-${count.index}"
    host_ip     = var.vm_lb_ips[count.index]
    node_labels = ""
  }
}

# Master hostname and ip list template #
data "template_file" "master_hosts" {

  count = length(var.vm_master_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname    = "${var.vm_name_prefix}-master-${count.index}"
    host_ip     = var.vm_master_ips[count.index]
    node_labels = ""
  }
}

# Worker hostname and ip list template #
data "template_file" "worker_hosts" {

  count = length(var.vm_worker_ips)

  template = file("templates/ansible_hosts.tpl")

  vars = {
    hostname    = "${var.vm_name_prefix}-worker-${count.index}"
    host_ip     = var.vm_worker_ips[count.index]
    node_labels = length(var.vm_worker_node_label) > 0 ? "node_labels=\"{'node-role.kubernetes.io/${var.vm_worker_node_label}':''}\"" : ""
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
    bind_ip = var.vm_lb_vip
  }
}

# HAProxy server backend template #
data "template_file" "haproxy_backend" {

  count = length(var.vm_master_ips)

  template = file("templates/haproxy_backend.tpl")

  vars = {
    prefix_server     = var.vm_name_prefix
    backend_server_ip = var.vm_master_ips[count.index]
    count             = count.index
  }
}

# Keepalived master template #
data "template_file" "keepalived_master" {
  template = file("templates/keepalived_master.tpl")

  vars = {
    network_interface = var.network_interface
    virtual_ip        = var.vm_lb_vip
  }
}

# Keepalived backup template #
data "template_file" "keepalived_backup" {
  template = file("templates/keepalived_backup.tpl")

  vars = {
    network_interface = var.network_interface
    virtual_ip        = var.vm_lb_vip
  }
}


#======================================================================================
# Local files
#======================================================================================

# Create Kubespray all.yml configuration file from template #
resource "local_file" "kubespray_all" {
  content  = data.template_file.kubespray_all.rendered
  filename = "config/group_vars/all/all.yml"
}

# Create Kubespray k8s-cluster.yml configuration file from template #
resource "local_file" "kubespray_k8s_cluster" {
  content  = data.template_file.kubespray_k8s_cluster.rendered
  filename = "config/group_vars/k8s-cluster/k8s-cluster.yml"
}

# Create Kubespray addons.yml configuration file from template #
resource "local_file" "kubespray_addons" {

  count = var.kubespray_custom_addons_enabled == "false" ? 1 : 0

  content = data.template_file.kubespray_addons[count.index].rendered
  filename = "config/group_vars/k8s-cluster/addons.yml"
}

# Copy custom Kubespray addons.yml configuration #
resource "local_file" "kubespray_custom_addons" {

  count = var.kubespray_custom_addons_enabled == "true" ? 1 : 0

  content = data.template_file.kubespray_custom_addons[count.index].rendered
  filename = "config/group_vars/k8s-cluster/addons.yml"
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

# Create keepalived backup configuration file from template #
resource "local_file" "keepalived_backup" {
  content  = data.template_file.keepalived_backup.rendered
  filename = "config/keepalived-backup.cfg"
}

#======================================================================================
# Null resources - K8s cluster configuration using Kubespray
#======================================================================================

# Modify permissions on config directory #
resource "null_resource" "config_permissions" {
  provisioner "local-exec" {
    command = "chmod -R 700 config"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    local_file.kubespray_addons,
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
  count = var.action == "create" ? 1 : 0

  provisioner "local-exec" {
    command = "cd ansible/haproxy && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} ${lookup(local.extra_args, var.vm_distro, local.default_extra_args)} haproxy.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.haproxy
  ]
}

# Create Kubespray Ansible playbook #
resource "null_resource" "kubespray_create" {

  count = var.action == "create" ? 1 : 0

  provisioner "local-exec" {
    command = "cd ansible/kubespray && virtualenv venv && . venv/bin/activate && pip install -r requirements.txt && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" ${lookup(local.extra_args, var.vm_distro, local.default_extra_args)} cluster.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.haproxy_install,
    null_resource.kubespray_download,
  ]
}

# Execute scale Kubespray Ansible playbook #
resource "null_resource" "kubespray_add" {
  count = var.action == "add_worker" ? 1 : 0

  provisioner "local-exec" {
    command = "cd ansible/kubespray && virtualenv venv && . venv/bin/activate && pip install -r requirements.txt && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" ${lookup(local.extra_args, var.vm_distro, local.default_extra_args)} scale.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install
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
    command = "cd ansible/kubespray && virtualenv venv && . venv/bin/activate && pip install -r requirements.txt && ansible-playbook -i ../../config/hosts.ini -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} -e \"kube_version=${var.k8s_version}\" ${lookup(local.extra_args, var.vm_distro, local.default_extra_args)} upgrade-cluster.yml"
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install
  ]
}


# Create the local admin.conf kubectl configuration file #
resource "null_resource" "kubectl_configuration" {
  provisioner "local-exec" {
    command = "ansible -i ${var.vm_master_ips[0]}, -b --user=${var.vm_user} --private-key=${var.vm_ssh_private_key} ${lookup(local.extra_args, var.vm_distro, local.default_extra_args)} -m fetch -a 'src=/etc/kubernetes/admin.conf dest=config/admin.conf flat=yes' all"
  }

  #  provisioner "local-exec" {
  #    command = "sed -i 's/lb-apiserver.kubernetes.local/${var.vm_lb_vip}/g' config/admin.conf && chmod 700 config/admin.conf"
  #  }

  #  provisioner "local-exec" {
  #    command = "chmod 600 config/admin.conf"
  #  }

  depends_on = [null_resource.kubespray_create]
}
