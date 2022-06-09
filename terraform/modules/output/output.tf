output "cluster" {
  value = {

    nodeTemplate = {
      user = var.vm_user
      ssh = {
        privateKeyPath = abspath(var.vm_ssh.privateKeyPath)
      }
      os = {
        distro           = var.vm_os.distro
        source           = var.vm_os.source
        networkInterface = var.vm_os.networkInterface
      }
    }

    nodes = {
      loadBalancer = {
        vip = (length(var.lb_nodes) == 0
          ? var.master_nodes[0].ip
          : (length(var.lb_nodes) > 1
            ? var.lb_vip
            : (var.lb_vip != null
              ? var.lb_vip
              : var.lb_nodes[0].ip
            )
          )
        )
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
  description = "Evaluated cluster section of the config."
}
