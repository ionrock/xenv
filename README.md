# xenv

[![](https://travis-ci.org/ionrock/xenv.svg?branch=master)](https://travis-ci.org/ionrock/xenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionrock/xenv)](https://goreportcard.com/report/github.com/ionrock/xenv)
[![GoDoc](https://godoc.org/github.com/ionrock/xenv?status.svg)](https://godoc.org/github.com/ionrock/xenv)


Xenv provides an executable environment that makes using configuration
consistent between development, CI/CD and production.

There is a common pattern that emerges when running applications.

 1. Configuration data is defined
 2. Configuration data is transformed and formatted
 3. Configuration data is applied to run (or restart) the application

For example, in Chef, you might define a set of attributes as your
configuration data, manipulate that data and write some configuration
files, and finally start your application or restart it when some new
data requires doing so.

[Confd](http://www.confd.io) and
[consul-template](https://github.com/hashicorp/consul-template) is
 another example where data is defined and stored in a key/value
 store, templates are written and the applications are restarted when
 necessary to get new configuration.

Where Xenv improves on these patterns is by providing a generic
mechanism that is not dependent on a specific configuration management
system (chef, puppet, ansible), key value store (etcd, consul) or
platform (kubernetes, mesos). Instead xenv makes it possible to create
interfaces that successfully hide the platform requirements in order
to simply the development and deployment process.

Xenv solves these sorts of issues by providing a dynamic, yet
consistent means of:

 - managing configuration data through environment variables and/or
   templates
 - stopping the process if configuration data changes (allowing k8s or
   systemd to restart the process)
 - pre/post tasks to implement service discovery, start sidecar
   processes, etc.


## Usage

The essence of `xenv` is the configuration file. It is a simple list
of entries that happen one after another.

Lets consider a gRPC service that runs in some production
environment. In production we need a few things:

 - A frontend proxy that does TLS termination and some health checking.

 - Secrets from something like [vault](https://www.vaultproject.io/)
   in order to connect to other services or get an auth token.

 - A configuration file and/or environment variables to turn on
   features and set functionality.

 - Registration with a service discovery system or DNS.

Here is an example:
```yaml
---

# Set up some environment variables. This structure is flattened using `_` between levels.
- env:
    # Set a single value
    foo: bar

    # Set a value to the result of a script.
    baz: '`cat baz.json | jq -r .baz`'

# Gather more environment data using a script that outputs JSON or
# YAML. A good example would be pulling secrets/certs from a secret store.
- envscript: 'curl http://httpbin.org/ip'

# We can use the environment and write templates using Go's template
# syntax. This format is similar to consul-template.
- template:
    template: foo.conf.tmpl
    target: /etc/foo.conf
    owner: nobody
    group: nobody
    mode: 0600

- task:
  name: start-envoy
  cmd: systemd-run --unit=myapp-envoy --property Restart=always -- envoy

- task:
  name: register-service
  cmd: svc-register.sh

# Anything defined in `post` will be called after the command exits,
# no matter the exit code.
- post:

  # We can run commands after the process exits such as cleaning up
  # secret files or unregistering from service discovery.
  - task:
      name: stop-envoy
      cmd: systemctl stop myapp-envoy.service
```

The actual command can be called with `xenv --config env.yml -- mysvc
start` in an init script, container `CMD`, CI pipeline or
orchestration system of choice.
