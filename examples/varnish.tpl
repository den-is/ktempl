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
