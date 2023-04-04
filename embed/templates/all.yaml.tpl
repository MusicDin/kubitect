##
# Kubesprays's source file (v2.17.1):
# https://github.com/kubernetes-sigs/kubespray/blob/v2.17.1/inventory/sample/group_vars/all/all.yml
##
---
apiserver_loadbalancer_domain_name: "{{ .InfraNodes.LoadBalancer.VIP }}"
deploy_container_engine: true
etcd_kubeadm_enabled: false
loadbalancer_apiserver:
	address: "{{ .InfraNodes.LoadBalancer.VIP }}"
	port: 6443
# CoreDNS version 1.9.3 is not supported on Kubernetes version 1.23 - 1.25.
coredns_version: "v1.8.6"
## Upstream dns servers
# upstream_dns_servers:
#   - 8.8.8.8
#   - 8.8.4.4