output "nodes" {
  #value = [for node_lists in [module.lb_module, module.master_module, module.worker_module]: {
  #    for nodes in node_lists : {
  #        for node in nodes : node.vm_info => obj
  #    }
  #}]
  #value = module.lb_module
  #value = [for node in module.lb_module: node.vm_info]
  value = flatten([
    [for node in module.lb_module: node.vm_info], 
    [for node in module.master_module: node.vm_info],
    [for node in module.worker_module: node.vm_info]
  ])
  description = "List of all nodes."
}
/*
output "lb_nodes" {
  value = module.lb_module
  description = "List of load balancer nodes."
}

output "master_nodes" {
  value = module.master_module
  description = "List of master nodes."
}

output "worker_nodes" {
  value = module.worker_module
  description = "List of worker nodes."
}*/