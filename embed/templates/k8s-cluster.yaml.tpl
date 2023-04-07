##
# Kubesprays's source file (v2.17.1):
# https://github.com/kubernetes-sigs/kubespray/blob/v2.17.1/inventory/sample/group_vars/k8s_cluster/k8s-cluster.yml
##
---
auto_renew_certificates: {{ .Config.Kubernetes.Other.AutoRenewCertificates }}
cluster_name: cluster.local
dns_mode: {{ .Config.Kubernetes.DnsMode }}
kube_version: {{ .Config.Kubernetes.Version }}
kube_network_plugin: {{ .Config.Kubernetes.NetworkPlugin }}
kube_proxy_strict_arp: true
resolvconf_mode: host_resolvconf