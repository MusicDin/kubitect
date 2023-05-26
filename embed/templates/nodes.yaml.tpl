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
			ip: {{ .IP }}
			{{- if .IP6 }}
			ip6: {{ .IP6 }}
			{{- end }}
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
			ip: {{ .IP }}
			{{- if .IP6 }}
			ip6: {{ .IP6 }}
			{{- end }}
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