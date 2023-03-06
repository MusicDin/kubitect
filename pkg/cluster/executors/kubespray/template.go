package kubespray

import (
	"path"

	"github.com/MusicDin/kubitect/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/pkg/utils/template"
)

const groupVarsDir = "group_vars"

type KubesprayAllTemplate struct {
	InfraNodes modelconfig.Nodes
	configDir  string
}

func NewKubesprayAllTemplate(configDir string, infraNodes modelconfig.Nodes) KubesprayAllTemplate {
	return KubesprayAllTemplate{
		configDir:  configDir,
		InfraNodes: infraNodes,
	}
}

func (t KubesprayAllTemplate) Name() string {
	return "all.yaml"
}

func (t KubesprayAllTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "all", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayAllTemplate) Template() string {
	return template.TrimTemplate(`
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
		## Upstream dns servers
		# upstream_dns_servers:
		#   - 8.8.8.8
		#   - 8.8.4.4
	`)
}

type KubesprayK8sClusterTemplate struct {
	Config    modelconfig.Config
	configDir string
}

func NewKubesprayK8sClusterTemplate(configDir string, config modelconfig.Config) KubesprayK8sClusterTemplate {
	return KubesprayK8sClusterTemplate{
		configDir: configDir,
		Config:    config,
	}
}

func (t KubesprayK8sClusterTemplate) Name() string {
	return "k8s-cluster.yaml"
}

func (t KubesprayK8sClusterTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "k8s_cluster", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayK8sClusterTemplate) Delimiters() (string, string) {
	return "<<", ">>"
}

func (t KubesprayK8sClusterTemplate) Template() string {
	return template.TrimTemplate(`
		##
		# Kubesprays's source file (v2.17.1):
		# https://github.com/kubernetes-sigs/kubespray/blob/v2.17.1/inventory/sample/group_vars/k8s_cluster/k8s-cluster.yml
		##
		---
		auto_renew_certificates: << .Config.Kubernetes.Other.AutoRenewCertificates >>
		# TODO: Support custom DNS domain
		cluster_name: cluster.local
		dns_mode: << .Config.Kubernetes.DnsMode >>
		kube_version: << .Config.Kubernetes.Version >>
		kube_network_plugin: << .Config.Kubernetes.NetworkPlugin >>
		kube_proxy_strict_arp: true
		resolvconf_mode: host_resolvconf
	`)
}

type KubesprayAddonsTemplate struct {
	configDir string
	Addons    string
}

func NewKubesprayAddonsTemplate(configDir string, addons string) KubesprayAddonsTemplate {
	return KubesprayAddonsTemplate{
		configDir: configDir,
		Addons:    addons,
	}
}

func (t KubesprayAddonsTemplate) Name() string {
	return "addons.yaml"
}

func (t KubesprayAddonsTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "k8s_cluster", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayAddonsTemplate) Template() string {
	return "{{ .Addons }}"
}

type KubesprayEtcdTemplate struct {
	configDir string
}

func NewKubesprayEtcdTemplate(configDir string) KubesprayEtcdTemplate {
	return KubesprayEtcdTemplate{configDir}
}

func (t KubesprayEtcdTemplate) Name() string {
	return "etcd.yaml"
}

func (t KubesprayEtcdTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayEtcdTemplate) Template() string {
	return template.TrimTemplate(`
		##
		# Kubesprays's source file (v2.17.1):
		# https://github.com/kubernetes-sigs/kubespray/blob/v2.17.1/inventory/sample/group_vars/etcd.yml
		##
		---
		etcd_deployment_type: host
	`)
}

type HostsTemplate struct {
	configDir  string
	SshKeyFile string
	Hosts      []modelconfig.Host
}

func NewHostsTemplate(configDir, sshPrivateKeyPath string, hosts []modelconfig.Host) HostsTemplate {
	return HostsTemplate{
		configDir:  configDir,
		SshKeyFile: sshPrivateKeyPath,
		Hosts:      hosts,
	}
}

func (t HostsTemplate) Name() string {
	return "hosts.yaml"
}

func (t HostsTemplate) Write() error {
	dstPath := path.Join(t.configDir, t.Name())
	return template.Write(t, dstPath)
}

func (t HostsTemplate) Functions() map[string]interface{} {
	return map[string]interface{}{
		"isRemoteHost": isRemoteHost,
	}
}

// isRemoteHost returns true id host's connection type equals REMOTE.
func isRemoteHost(host modelconfig.Host) bool {
	return host.Connection.Type == modelconfig.REMOTE
}

func (t HostsTemplate) Template() string {
	return template.TrimTemplate(`
		{{- $pkPath := .SshKeyFile -}}
		all:
			hosts:
			{{- range .Hosts }}
				{{ .Name }}:
				{{- if isRemoteHost . }}
					ansible_connection: ssh
					ansible_user: {{ .Connection.User }}
					ansible_host: {{ .Connection.IP }}
					ansible_port: {{ .Connection.SSH.Port }}
					ansible_private_key_file: {{ $pkPath }}
				{{- else }}
					ansible_connection: local
					ansible_host: localhost
				{{- end }}
			{{- end }}
			children:
				kubitect_hosts:
					hosts:
					{{- range .Hosts }}
						{{ .Name }}:
					{{- end }}
	`)
}

type NodesTemplate struct {
	configDir   string
	ConfigNodes modelconfig.Nodes
	InfraNodes  modelconfig.Nodes
}

func NewNodesTemplate(configDir string, configNodes, infraNodes modelconfig.Nodes) NodesTemplate {
	return NodesTemplate{
		configDir:   configDir,
		ConfigNodes: configNodes,
		InfraNodes:  infraNodes,
	}
}

func (t NodesTemplate) Name() string {
	return "nodes.yaml"
}

func (t NodesTemplate) Write() error {
	dstPath := path.Join(t.configDir, t.Name())
	return template.Write(t, dstPath)
}

func (t NodesTemplate) Template() string {
	return template.TrimTemplate(`
		{{- $cfgNodes := .ConfigNodes -}}
		all:
			hosts:
			{{- /* Load balancers */ -}}
			{{- range .InfraNodes.LoadBalancer.Instances }}
				{{- $i := $cfgNodes.LoadBalancer.Instances | select "Id" .Id | first }}
				{{ .Name }}:
					ansible_host: {{ .IP }}
					priority: {{ $i.Priority }}
			{{- end }}
			{{- /* Master nodes */ -}}
			{{- range .InfraNodes.Master.Instances }}
				{{- $i := $cfgNodes.Master.Instances | select "Id" .Id | first }}
				{{ .Name }}:
					ansible_host: {{ .IP }}
					{{- if $i.Labels }}
					node_labels:
						{{- range $k, $v := $i.Labels }}
						{{ $k }}: {{ $v }}
						{{- end }}
					{{- end }}
					{{- if $i.Taints }}
					node_taints:
						{{- range $i.Taints }}
						- "{{ . }}"
						{{- end }}
					{{- end }}
			{{- end }}
			{{- /* Worker nodes */ -}}
			{{- range .InfraNodes.Worker.Instances }}
				{{- $i := $cfgNodes.Worker.Instances | select "Id" .Id | first }}
				{{ .Name }}:
					ansible_host: {{ .IP }}
					{{- if $i.Labels }}
					node_labels:
						{{- range $k, $v := $i.Labels }}
						{{ $k }}: {{ $v }}
						{{- end }}
					{{- end }}
					{{- if $i.Taints }}
					node_taints:
						{{- range $i.Taints }}
						- "{{ . }}"
						{{- end }}
					{{- end }}
			{{- end }}
			children:
				haproxy:
					hosts:
					{{- range .InfraNodes.LoadBalancer.Instances }}
						{{ .Name }}:
					{{- end }}
				etcd:
					hosts:
					{{- range .InfraNodes.Master.Instances }}
						{{ .Name }}:
					{{- end }}
				k8s_cluster:
					children:
						kube_control_plane:
							hosts:
							{{- range .InfraNodes.Master.Instances }}
								{{ .Name }}:
							{{- end }}
						kube_node:
							hosts:
							{{- if .InfraNodes.Worker.Instances }}
								{{- range .InfraNodes.Worker.Instances }}
								{{ .Name }}:
								{{- end }}
							{{- else }}
								{{- /* No worker nodes -> masters also become workers */ -}}
								{{- range .InfraNodes.Master.Instances }}
								{{ .Name }}:
								{{- end }}
							{{- end }}
	`)
}
