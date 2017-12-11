= xenv

Xenv provides an executable environment for managing a process at
runtime. Whether that is managing environment variables, starting
sidecar / helper processes or performing small pre/post tasks, `xenv`
makes it possible to encapsulate complexity found in configuration
management, init systems and orchestration systems.

== Usage

The essence of `xenv` is the configuration file. It is a simple list
of entries that happen one after another.

Lets consider a gRPC service that runs in some production environment. In production we need a few things:

 - A frontend proxy that does TLS termination and some health checking.

 - Secrets from soemthing like [vault](https://www.vaultproject.io/)
   in order to connect to other services or get an auth token.

 - A configuration file and/or environment variables to turn on
   features and set functionality.

 - Standardized logging for sending logs to a centralized collector or local syslog.

 - Registration with a service discovery system or DNS.

The `xenv` executable allows you to consider all these sorts of
details when starting the process as opposed to depending on
configuration management or an orchestration system to provide the
necessary functionality. The benefit here is that you can begin to
test these variable conditions during development, encapsulate
assumptions and reliably debug the functionality because it is the
same code and not a feature of an upstream system like Kubernetes, CI
or a provisioned host.

Here is an example:
```
---
- task:
    name: load-mesh
	cmd: apt-get install -y nginx myorg-mesh-config

- service:
    name: nginx
    cmd: /usr/local/bin/myorg-start-nginx

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
