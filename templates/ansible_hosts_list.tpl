[all]
${lb_hosts}
${master_hosts}
${worker_hosts}

[haproxy]
${lb_nodes}

[kube_control_plane]
${master_nodes}

[etcd]
${master_nodes}

[kube_node]
${worker_nodes}

[k8s_cluster:children]
kube_control_plane
kube_node
