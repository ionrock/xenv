= xenv

Xenv provides an executable environment for managing a process at
runtime. Whether that is managing environment variables, starting
sidecar / helper processes or performing small pre/post tasks, `xenv`
makes it possible to encapsulate complexity found in configuration
management, init systems and orchestration systems.

== Usage

The essence of `xenv` is the configuration file. It is a simple list
of entries that happen one after another.

```
---
- service:
    name: py-static-web
    cmd: python -m SimpleHTTPServer

- env:
    foo: bar
    baz: '`cat baz.json | jq -r .baz`'

- envscript: 'curl http://httpbin.org/ip'

- template: 'foo.conf.tmpl:/etc/foo.conf'

- template:
    source: foo.conf.tmpl
	dest: /etc/foo.conf

- run_command

- task:
    name: apt
    cmd: 'echo "apt get updating..."'
    dir: '/tmp'

- watch:
    file: /etc/cert/foo
    action:
      restart: py-static-web
```

There are two main elments that are at working.

  1. Configuration data composition.
  2. Process and task management.

The configuration data can be composed using the `env` and `envscript`
fields. These allow for you to build an environment that can be used
as-is or in by writing configuration files.

The process management allows you use processes as sidecars or helpers
to prepare the process to run.
