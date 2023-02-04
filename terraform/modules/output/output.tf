output "nodes" {
  value = {
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
  
  description = "Nodes information after provisioning."
}
