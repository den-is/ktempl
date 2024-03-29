# ktempl - render Go templates with Nodes data from a Kubernetes

:warning: **Early stage of development! Use with caution! No warranties provided.** :construction:

ktempl renders template files formatted using Go [template][gotemplate] with data returned from a Kubernetes cluster.

Primary data for ktempl is a **list of nodes** returned by a query to a Kubernetes.
By default, ktempl fetches all nodes from the Kubernetes cluster.
You can limit/filter nodes list using `-l` and available kubernetes key=value labels.
Also you get list of nodes that host specific pods only by providing `-p` and `-l` to select specific pods.

Secondary data is a `key=value` data provided by user using `--set` arguments.

Optionally, after successful render, ktemp might execute command provided by user.

Initial motivation was to write piece of software which can be used as a glue between not flexible services, without built-in service discovery capabilities, and Kubernetes.
One such example is the Varnish HTTP accelerator, where you need to pass backend servers configuration, and reload Varnish daemon if configuration changes.

Inspired by many configuration management frameworks which operate templates to generate configs, e.g., [Helm][helm], [consul-template][consultemplate], [Ansible][ansibletemplate], [Puppet][puppettemplate], etc.

## Prerequisites

You should have a `kubeconfig` file with proper authentication and context details supplied.

User listed in the `kubeconfig` should be allowed to list nodes and pods.

As a minimum user should be allowed to list either nodes or pods.

Default kubeconfig location `~/.kube/config`

More details on how to obtain that file can be found [here][kubeconfigdoc].

## Install

Download from the [releases][releasespage] page. Example setup on linux machine:

```sh
curl -OL https://github.com/den-is/ktempl/releases/download/0.0.1/ktempl_0.0.1_`uname -s`_x86_64.tar.gz &&\
tar xzvf ktempl_0.0.1_`uname -s`_x86_64.tar.gz &&\
cp ktempl_0.0.1_`uname -s`_x86_64/ktempl /usr/local/bin
```

Compile it yourself. Minimum recommended Go version is 1.17+.

```sh
go get -u github.com/den-is/ktempl
```

## Template language and scope variables

Templates are formated using [Go template language][gotemplate] and support [Sprig template library][sprig].

Two main variables are passed to template for rendering:

- `.Nodes` - list of nodes returned by Kubernetes
- `.Values` - dictionary of key=values passed to ktempl using `--set` in the terminal or `values:` in the config file.

Each node in `.Nodes` has next fields

- `.mynode.Name`
- `.mynode.InternalIP`
- `.mynode.Annotations`
- `.mynode.Labels`

Example template file

```jinja
# Welcome to {{ .Values.title  | lower | repeat 5 }} <<<<
{{- range $i, $n := .Nodes }}
Node {{$n.Name }} has {{ $n.InternalIP }} IP and port {{ $.Values.port }}
{{- end }}
```

## Example usage

More examples [here](/examples/)

```sh
# query nodes with a specific label
ktempl -l disk=ssd -t myexamle.tpl

# query pods with a specific label
ktempl -p -l app=myapp1 -t myexamle.tpl

# extended example
ktempl -l app=stagingapps -t example.tpl -o output.conf --set port=32456 --exec="touch success_exec.txt"
```

## Config file

Optionally you can supply ktempl with a config file.
The default is `config.yaml` in the same directory as binary.
Or `/etc/ktempl/config.yaml`.
Or whatever you supply with `-c` command-line option.

## List of available configuration options

| Config file settings  |         CLI flags         | Description                                                              |
| ----------------------| --------------------------| ------------------------------------------------------------------------ |
| `kubeconfig`          | `-k, --kubeconfig`        | Path to kubeconfig                                                       |
| `pods`                | `-p, --pods`              | Query pods and get nodes they are running on                             |
| `namespace`           | `-n, --namespace`         | Kubernetes namespace where to look Pods for. Used with `-p`              |
| `selector`            | `-l, --selector`          | Kubernetes [label selectors][labelselectors] string                      |
| `template`            | `-t, --template`          | Path to template                                                         |
| `output`              | `-o, --output`            | Path where to put rendered results. default stdout                       |
| `permissions`         | `N/A`                     | Output file permissions. default 0644. should be in 4 digit format!      |
| `set`                 | `--set`                   | Additional key=values passed to template rendering engine                |
| `exec`                | `-e, --exec`              | Command to execute after successful template render                      |
| `log.file`            | `--log-file`              | Path to log file. Allowed values `disabled`, `stdout`, `stderr`, `path_to_log_file_dst`. Default `disabled`. |
| `log.level`           | `--log-level`             | Minimum log message level to log. Default `info`. Available levels by hierarchy: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic` |
| `N/A`                 | `-c, --config`            | Path to ktempl config file                                               |
| `daemon`              | `-d, --daemon`            | Run ktempl in service mode, rather than singleshot                       |
| `interval`            | `-i, --interval`          | Interval between polls. default 15s. Valid time units are "s", "m", "h". |
| `retries`             | `N/A`                     | _NOT YET IMPLEMENTED_ Number of retries to fetch data from Kubernetes    |
| `timeout`             | `N/A`                     | _NOT YET IMPLEMENTED_ ktempl operations timeout                          |

### Contact

Denis Iskandarov denis@cloudboom.io

[gotemplate]: https://golang.org/pkg/text/template/
[sprig]: http://masterminds.github.io/sprig/
[consultemplate]: https://github.com/hashicorp/consul-template
[helm]: https://helm.sh/
[ansibletemplate]: https://docs.ansible.com/ansible/latest/modules/template_module.html
[puppettemplate]: https://puppet.com/docs/puppet/latest/lang_template.html
[kubeconfigdoc]: https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/
[labelselectors]: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
[releasespage]: https://github.com/den-is/ktempl/releases
