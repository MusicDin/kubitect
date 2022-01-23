#================================
# Local variables
#================================

# Local variables used in many resources #
locals {
  extra_args = {
    debian = "--timeout 3000 --verbose --extra-vars 'ansible_become_method=su'"
    ubuntu = "--timeout 3000 --verbose"
    centos = "--timeout 3000 --verbose"
  }
  default_extra_args    = "--timeout 3000 --verbose"
  dashboard_namespace   = "kube-system"
  initial_vrrp_priority = 201
}


#======================================================================================
# Template files
#======================================================================================

# Kubespray all.yml template (Currently supports only 1 load balancer) #
data "template_file" "kubespray_all" {

  template = file("templates/kubespray/kubespray_all.tpl")

  vars = {
    loadbalancer_apiserver = (
      length(var.lb_nodes) > 0
      ? var.lb_vip
      : var.master_nodes[0].ip
    )
  }
}

# Kubespray k8s-cluster.yml template #
data "template_file" "kubespray_k8s_cluster" {

  template = file("templates/kubespray/kubespray_k8s_cluster.tpl")

  vars = {
    kube_version        = var.kubernetes_version
    kube_network_plugin = var.kubernetes_networkPlugin
    dns_mode            = var.kubernetes_dnsMode

    # If MetalLB is enable than strict ARP is set to true in k8s-cluster.yml
    kube_proxy_strict_arp = (var.kubernetes_kubespray_addons_enabled
      ? yamldecode(
        templatefile(var.kubernetes_kubespray_addons_configPath, {})
      )["metallb_enabled"]
      : false
    )

  }
}

# Kubespray etcd.yml template #
data "template_file" "kubespray_etcd" {
  template = file("templates/kubespray/kubespray_etcd.tpl")
}

# Load balancer hostname and ip list template #
data "template_file" "lb_hosts" {

  for_each = { for node in var.lb_nodes : node.name => node }

  template = file("templates/ansible/ansible_hosts.tpl")

  vars = {
    hostname    = each.value.name
    host_ip     = each.value.ip
    node_labels = ""
  }
}

# Master hostname and ip list template #
data "template_file" "master_hosts" {

  for_each = { for node in var.master_nodes : node.name => node }

  template = file("templates/ansible/ansible_hosts.tpl")

  vars = {
    hostname    = each.value.name
    host_ip     = each.value.ip
    node_labels = ""
  }
}

# Worker hostname and ip list template #
data "template_file" "worker_hosts" {

  for_each = { for node in var.worker_nodes : node.name => node }

  template = file("templates/ansible/ansible_hosts.tpl")

  vars = {
    hostname = each.value.name
    host_ip  = each.value.ip
    node_labels = (
      length(var.worker_node_label) > 0
      ? "node_labels=\"{'node-role.kubernetes.io/${var.worker_node_label}':''}\""
      : ""
    )
  }
}

# Template for hosts.ini file #
data "template_file" "kubespray_hosts" {
  template = file("templates/ansible/ansible_hosts_list.tpl")

  vars = {
    lb_hosts     = chomp(join("", [for tpl in data.template_file.lb_hosts : tpl.rendered]))
    master_hosts = chomp(join("", [for tpl in data.template_file.master_hosts : tpl.rendered]))
    worker_hosts = chomp(join("", [for tpl in data.template_file.worker_hosts : tpl.rendered]))
    lb_nodes     = join("\n", [for node in var.lb_nodes : node.name])
    master_nodes = join("\n", [for node in var.master_nodes : node.name])
    worker_nodes = (
      length(var.worker_nodes) > 0
      ? join("\n", [for node in var.worker_nodes : node.name])
      : join("\n", [for node in var.master_nodes : node.name])
    )
  }
}

# HAProxy template #
data "template_file" "haproxy" {
  template = file("templates/haproxy/haproxy.tpl")

  vars = {
    bind_ip = var.lb_vip
  }
}

# HAProxy server backend template #
data "template_file" "haproxy_backend" {

  for_each = { for node in var.master_nodes : node.name => node }

  template = file("templates/haproxy/haproxy_backend.tpl")

  vars = {
    server_name = each.value.name
    server_ip   = each.value.ip
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
  filename = "config/group_vars/k8s_cluster/k8s-cluster.yml"
}

# Create Kubespray etcd.yml configuration file from template #
resource "local_file" "kubespray_etcd" {
  content  = data.template_file.kubespray_etcd.rendered
  filename = "config/group_vars/etcd.yml"
}

# Create Kubespray addons.yml configuration file from template #
resource "local_file" "kubespray_addons" {

  count = var.kubernetes_kubespray_addons_enabled ? 1 : 0

  content  = templatefile(var.kubernetes_kubespray_addons_configPath, {})
  filename = "config/group_vars/k8s_cluster/addons.yml"
}

# Create hosts.ini file from template #
resource "local_file" "kubespray_hosts" {
  content  = data.template_file.kubespray_hosts.rendered
  filename = "config/hosts.ini"
}

# Create HAProxy configuration file from template #
resource "local_file" "haproxy" {
  content = join("",
    concat(
      [data.template_file.haproxy.rendered],
      [for tpl in data.template_file.haproxy_backend : tpl.rendered]
    )
  )

  filename = "config/haproxy/haproxy.cfg"
}

# Create Keepalived configuration file from template #
resource "local_file" "keepalived" {

  for_each = { for node in var.lb_nodes : node.id => node }

  content = templatefile("templates/keepalived/keepalived.tpl",
    {
      network_interface = var.vm_network_interface
      virtual_ip        = var.lb_vip
      priority          = local.initial_vrrp_priority - each.value.id
    }
  )
  filename = "config/keepalived/keepalived-${each.value.name}.cfg"
}


#======================================================================================
# Null resources - K8s cluster configuration using Kubespray
#======================================================================================

# Modify permissions on config directory #
resource "null_resource" "config_permissions" {
  provisioner "local-exec" {
    command = "chmod -R 700 config"
  }
}

# Clone Kubespray repository #
resource "null_resource" "kubespray_download" {
  provisioner "local-exec" {
    command = <<-EOF
              cd ansible
              rm -rf kubespray
              git clone --branch ${var.kubernetes_kubespray_version} ${var.kubernetes_kubespray_url}
              EOF
  }
}

# Execute create Kubernetes HAProxy playbook #
resource "null_resource" "haproxy_install" {

  count = var.action == "create" ? 1 : 0

  provisioner "local-exec" {
    command = <<-EOF
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible-playbook \
                --inventory config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                $EXTRA_ARGS \
                ansible/haproxy/haproxy.yml
              EOF

    environment = {
      SSH_USER        = var.vm_user
      SSH_PRIVATE_KEY = var.vm_ssh_private_key
      EXTRA_ARGS      = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
    }
  }

  # Requires hosts.ini, HAProxy and Keepalive configuration files created #
  depends_on = [
    local_file.kubespray_hosts,
    local_file.haproxy,
    local_file.keepalived
  ]
}

# Create Kubespray Ansible playbook #
resource "null_resource" "kubespray_create" {

  count = var.action == "create" ? 1 : 0

  provisioner "local-exec" {
    command = <<-EOF
              cd ansible/kubespray
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible-playbook \
                --inventory ../../config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                --extra-vars "kube_version=$K8S_VERSION" \
                $EXTRA_ARGS \
                cluster.yml
              EOF

    environment = {
      SSH_USER        = var.vm_user
      SSH_PRIVATE_KEY = var.vm_ssh_private_key
      K8S_VERSION     = var.kubernetes_version
      EXTRA_ARGS      = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
    }
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    local_file.kubespray_addons,
    null_resource.kubespray_download,
    null_resource.haproxy_install
  ]
}

# Execute scale Kubespray Ansible playbook #
resource "null_resource" "kubespray_add" {

  count = var.action == "add_worker" ? 1 : 0

  provisioner "local-exec" {
    command = <<-EOF
              cd ansible/kubespray
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible-playbook \
                --inventory ../../config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                --extra-vars "kube_version=$K8S_VERSION" \
                $EXTRA_ARGS \
                scale.yml
              EOF

    environment = {
      SSH_USER        = var.vm_user
      SSH_PRIVATE_KEY = var.vm_ssh_private_key
      K8S_VERSION     = var.kubernetes_version
      EXTRA_ARGS      = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
    }
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install
  ]
}

# Takes care of removing worker from cluster's configuration #
resource "null_resource" "kubespray_remove" {

  for_each = { for node in var.worker_nodes : node.name => node }

  triggers = {
    vm_name            = each.value.name
    vm_user            = var.vm_user
    vm_ssh_private_key = var.vm_ssh_private_key
    extra_args         = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<-EOF
              cd ansible/kubespray
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible-playbook \
                --inventory ../../config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                --extra-vars "node=$VM_NAME delete_nodes_confirmation=yes" \
                $EXTRA_ARGS \
                remove-node.yml
              EOF

    environment = {
      VM_NAME         = self.triggers.vm_name
      SSH_USER        = self.triggers.vm_user
      SSH_PRIVATE_KEY = self.triggers.vm_ssh_private_key
      EXTRA_ARGS      = self.triggers.extra_args
    }
    on_failure = continue
  }

  # Prevents node to be removed from inventory before it's removed from cluster #
  depends_on = [
    local_file.kubespray_hosts
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
    command = <<-EOF
              cd ansible
              rm -rf kubespray
              git clone --branch ${var.kubernetes_kubespray_version} ${var.kubernetes_kubespray_url}
              EOF
  }

  provisioner "local-exec" {
    command = <<-EOF
              cd ansible/kubespray
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible-playbook \
                --inventory ../../config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                --extra-vars "kube_version=$K8S_VERSION" \
                $EXTRA_ARGS \
                upgrade-cluster.yml
              EOF

    environment = {
      SSH_USER        = var.vm_user
      SSH_PRIVATE_KEY = var.vm_ssh_private_key
      K8S_VERSION     = var.kubernetes_version
      EXTRA_ARGS      = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
    }
  }

  depends_on = [
    local_file.kubespray_hosts,
    local_file.kubespray_all,
    local_file.kubespray_k8s_cluster,
    null_resource.kubespray_download,
    null_resource.haproxy_install
  ]
}

# Fetch the local admin.conf kubectl configuration file #
resource "null_resource" "fetch_kubeconfig" {

  provisioner "local-exec" {
    command = <<-EOF
              virtualenv -p python3 venv && . venv/bin/activate && pip3 install -r requirements.txt
              ansible \
                --inventory config/hosts.ini \
                --become \
                --user=$SSH_USER \
                --private-key=$SSH_PRIVATE_KEY \
                --module-name fetch \
                --args "src=/etc/kubernetes/admin.conf dest=config/admin.conf flat=yes" \
                $EXTRA_ARGS \
                $MASTER_NODE_NAME
              EOF

    environment = {
      MASTER_NODE_NAME = var.master_nodes[0].name
      SSH_USER         = var.vm_user
      SSH_PRIVATE_KEY  = var.vm_ssh_private_key
      EXTRA_ARGS       = lookup(local.extra_args, var.vm_distro, local.default_extra_args)
    }
  }

  # Cluster needs to be deployed before kubeconfig can be fetched
  depends_on = [null_resource.kubespray_create]
}

# Copy kubeconfig into ~/.kube directory #
resource "null_resource" "copy_kubeconfig" {

  count = var.kubernetes_other_copyKubeconfig ? 1 : 0

  provisioner "local-exec" {
    command = "mkdir -p ~/.kube && cp config/admin.conf ~/.kube/"
  }

  # Kubeconfig needs to be fetched before it can be copied
  depends_on = [null_resource.fetch_kubeconfig]
}

# Creates Kubernetes dashboard service account #
#resource "null_resource" "k8s_dashboard_rbac" {
#
#  count = (var.k8s_dashboard_enabled && var.k8s_dashboard_rbac_enabled) ? 1 : 0
#
#  provisioner "local-exec" {
#    command = "sh scripts/dashboard-rbac.sh ${var.k8s_dashboard_rbac_user} ${local.dashboard_namespace}"
#  }
#
#  # Kubeconfig needs to be ready when before script for creating service account is executed
#  depends_on = [null_resource.fetch_kubeconfig]
#}
