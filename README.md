# xenv

Xenv provides an executable environment for managing a process at
runtime. Whether that is managing environment variables, starting
sidecar / helper processes or performing small pre/post tasks, `xenv`
makes it possible to encapsulate complexity found in configuration
management, init systems and orchestration systems.

## Usage

The essence of `xenv` is the configuration file. It is a simple list
of entries that happen one after another.

Lets consider a gRPC service that runs in some production
environment. In production we need a few things:

 - A frontend proxy that does TLS termination and some health checking.

 - Secrets from soemthing like [vault](https://www.vaultproject.io/)
   in order to connect to other services or get an auth token.

 - A configuration file and/or environment variables to turn on
   features and set functionality.

 - Standardized logging for sending logs to a centralized collector or
   local syslog.

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

# Set up some environment variables. This structure is flattened using `_` between levels.
- env:
	# Set a single value
    foo: bar

	# Set a value to the result of a script. This doesn't use a shell!
    baz: '`cat baz.json | jq -r .baz`'

# Gather more environment data using a script that outputs JSON or YAML
- envscript: 'curl http://httpbin.org/ip'

# We can use the environment and write templates using Go's template
# syntax. This format is similar to consul-template.
- template: 'foo.conf.tmpl:/etc/foo.conf'

# You can also be explicit
- template:
    source: foo.conf.tmpl
	dest: /etc/foo.conf

# Use the the run command to call the command. You can also use `--`
# and then add the command to the `xenv` call like `xenv --config env.yml -- mysvc start`
- run_command

# We can run commands after the process exits such as cleaning up
# secret files or unregistering from service discovery.
- task:
  name: remove-config
  cmd: rm /etc/foo.conf
```
