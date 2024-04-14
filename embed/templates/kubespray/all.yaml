##
# Kubesprays's source file (v2.17.1):
# https://github.com/kubernetes-sigs/kubespray/blob/v2.17.1/inventory/sample/group_vars/all/all.yml
##
---
apiserver_loadbalancer_domain_name: "{{ .Values.LoadBalancer.VIP }}"
deploy_container_engine: true
etcd_kubeadm_enabled: false
loadbalancer_apiserver:
	address: "{{ .Values.LoadBalancer.VIP }}"
	port: 6443
## Upstream dns servers
# upstream_dns_servers:
#   - 8.8.8.8
#   - 8.8.4.4
