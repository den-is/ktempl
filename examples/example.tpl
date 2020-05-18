# >>>> start of the template {{ .Values.title }} <<<<
{{- range $i, $n := .Nodes }}
Node {{$n.Name }} has {{ $n.InternalIP }} IP and port {{ $.Values.port }}
{{- end }}
