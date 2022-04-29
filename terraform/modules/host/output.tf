output "nodes" {
  value = flatten([
    [for node in module.lb_module : node.vm_info],
    [for node in module.master_module : node.vm_info],
    [for node in module.worker_module : node.vm_info]
  ])
  description = "List of all nodes."
}