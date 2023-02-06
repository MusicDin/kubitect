package kubespray

import (
	"cli/config/modelconfig"
	"cli/utils/template"
)

type HostsTemplate struct {
	projDir    string
	SshKeyFile string
	Hosts      []modelconfig.Host
}

func NewHostsTemplate(projectDir, sshPrivateKeyPath string, hosts []modelconfig.Host) HostsTemplate {
	return HostsTemplate{
		projDir:    projectDir,
		SshKeyFile: sshPrivateKeyPath,
		Hosts:      hosts,
	}
}

func (t HostsTemplate) Name() string {
	return "hosts.yaml"
}

func (t HostsTemplate) Functions() map[string]interface{} {
	return map[string]interface{}{
		"isRemoteHost": isRemoteHost,
	}
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
	ConfigNodes modelconfig.Nodes
	InfraNodes  modelconfig.Nodes
	projDir     string
}

func NewNodesTemplate(projectDir string, configNodes, infraNodes modelconfig.Nodes) NodesTemplate {
	return NodesTemplate{
		projDir:     projectDir,
		ConfigNodes: configNodes,
		InfraNodes:  infraNodes,
	}
}

func (t NodesTemplate) Name() string {
	return "nodes.yaml"
}

func (t NodesTemplate) Functions() map[string]interface{} {
	return map[string]interface{}{
		"isRemoteHost": isRemoteHost,
	}
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

// isRemoteHost returns true id host's connection type equals REMOTE.
func isRemoteHost(host modelconfig.Host) bool {
	return host.Connection.Type == modelconfig.REMOTE
}
