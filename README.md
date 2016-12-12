# Joxy

Joxy is a containerized TCP proxy that terminates TLS using certificates automatically issued by [Let´s encrypt](https://letsencrypt.org/). This makes it easy to add encryption to any TCP listener which is normally not supported by traditional HTTP loadbalancers such as a websockets server.

## Example

See the [Makefile](https://github.com/joonix/joxy/blob/master/Makefile) for an example container build. This assumes we have access to a google cloud registry. Change REGISTRY and IMAGE variables to match your own project.

See [kubernetes](https://github.com/joonix/joxy/tree/master/kubernetes) for an example deployment configuration.  
The deployment arguments must be changed:

- `domain`: the domain name that resolves to our service.
- `backend`: the non-TLS service we want to proxy to.