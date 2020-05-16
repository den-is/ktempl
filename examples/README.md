# ktempl examples

Example configuration [file](./config.yaml)

Couple template examples:

## Simple example

Create template file for example `example.tpl`, with such [contents](./example.tpl):

```txt
# >>>> start of the template {{ .Values.title }} <<<<
{{- range $i, $n := .Nodes }}
Node {{$n.Name }} has {{ $n.InternalIP }} IP
{{- end }}
```

Execute ktempl:

```sh
# selects pods by app=myapp1 and get's nodes they are running on. outputs to stdout.
ktempl -p -l app=myapp1 -t example.tpl --set title=Hello
```

## Varnish

Create template `varnish.tpl` with contents listed in varnish example [template](./varnish.tpl):

```txt
import directors;
{{ range $i, $n := .Nodes }}
backend {{$n.Name}} {
    .host = "{{$n.InternalIP}}";
    .port = "{{ $.Values.port }}";
}{{end}}

sub vcl_init {
  new bar = directors.round_robin();
{{range $i, $n := .Nodes}}
  bar.add_backend({{$n.Name}});{{end}}
}

sub vcl_recv {
  set req.backend_hint = bar.backend();
}
```

```sh
# selects nodes using app=stagingapps label
ktempl -l app=stagingapps -t varnish.tpl  -o backend.conf --set port=32456 --exec="systemctl reload varnishd"
```
