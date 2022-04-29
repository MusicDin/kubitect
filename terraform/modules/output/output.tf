output "cluster" {

  value = {

    networkInterface = var.vm_network_interface

    ssh = {
      user = var.vm_user
      pkey = abspath(var.vm_ssh_private_key)
    }

    nodes = {
      loadBalancer = {
        vip       = length(var.lb_nodes) > 0 ? var.lb_vip : var.master_nodes[0].ip
        instances = var.lb_nodes
      }

      master = {
        instances = var.master_nodes
      }

      worker = {
        instances = var.worker_nodes
      }
    }
  }

}
