{{- $cfgNodes := .Values.ConfigNodes -}}
{{- $infNodes := .Values.InfraNodes -}}
all:
	hosts:
	{{- /* Load balancers */ -}}
	{{- range $infNodes.LoadBalancer.Instances }}
		{{- $i := $cfgNodes.LoadBalancer.Instances | select "Id" .Id | first }}
		{{ .Name }}:
			ansible_host: {{ .IP }}
			priority: {{ $i.Priority }}
	{{- end }}
	{{- /* Master nodes */ -}}
	{{- range $infNodes.Master.Instances }}
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
	{{- range $infNodes.Worker.Instances }}
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
			{{- range $infNodes.LoadBalancer.Instances }}
				{{ .Name }}:
			{{- end }}
		etcd:
			hosts:
			{{- range $infNodes.Master.Instances }}
				{{ .Name }}:
			{{- end }}
		k8s_cluster:
			children:
				kube_control_plane:
					hosts:
					{{- range $infNodes.Master.Instances }}
						{{ .Name }}:
					{{- end }}
				kube_node:
					hosts:
					{{- if $infNodes.Worker.Instances }}
						{{- range $infNodes.Worker.Instances }}
						{{ .Name }}:
						{{- end }}
					{{- else }}
						{{- /* No worker nodes -> masters also become workers */ -}}
						{{- range $infNodes.Master.Instances }}
						{{ .Name }}:
						{{- end }}
					{{- end }}
