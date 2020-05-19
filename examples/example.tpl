# >>>> start of the template {{ .Values.title  | lower | repeat 5 }} <<<<
{{- range $i, $n := .Nodes }}
Node {{$n.Name }} has {{ $n.InternalIP }} IP and port {{ $.Values.port }}
{{- end }}
