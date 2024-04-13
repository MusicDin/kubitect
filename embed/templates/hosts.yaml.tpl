all:
	hosts:
	{{- range .Values }}
		{{ .Name }}:
		{{- if eq .Connection.Type "remote" }}
			ansible_connection: ssh
			ansible_user: {{ .Connection.User }}
			ansible_host: {{ .Connection.IP }}
			ansible_port: {{ .Connection.SSH.Port }}
			ansible_private_key_file: {{ .Connection.SSH.Keyfile }}
		{{- else }}
			ansible_connection: local
			ansible_host: localhost
		{{- end }}
	{{- end }}
	children:
		kubitect_hosts:
			hosts:
			{{- range .Values }}
				{{ .Name }}:
			{{- end }}
