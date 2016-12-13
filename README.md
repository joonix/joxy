# Joxy

Joxy is a containerized TCP proxy that terminates TLS using certificates automatically issued by [Let´s encrypt](https://letsencrypt.org/). This makes it easy to add encryption to any TCP listener which is normally not supported by traditional HTTP loadbalancers such as a websockets server.

## Dependencies

This project uses these tools and dependencies:

- [Go](golang.org) as well as [make](https://www.gnu.org/software/make/) for building the service
- [Godep](https://github.com/tools/godep) for maintaining Go dependencies
- [Kubernetes](kubernetes.io) available through [Container Engine](https://cloud.google.com/container-engine/)
- [Container Registry](https://cloud.google.com/container-registry/)

## Example

See the [Makefile](https://github.com/joonix/joxy/blob/master/Makefile) for an example container build. This assumes we have access to a google cloud registry. Change REGISTRY and IMAGE environment variables to match your own project.

See [kubernetes](https://github.com/joonix/joxy/tree/master/kubernetes) for an example deployment configuration.  
The deployment arguments must be changed:

- `domain`: the domain name that resolves to our service.
- `backend`: the non-TLS service we want to proxy to.

HTTP pprof is enabled and available on port 8080. When running on Kubernetes, you can forward this port to localhost

	kubectl port-forward joxy-4056657478-myuh6 8080:8080

You can then access `http://localhost:8080/debug/pprof/`, see [pprof documentation](https://golang.org/pkg/net/http/pprof/) for more information.

## Scalability

This service is dependent on running as a single instance due to how Let's encrypt challenging works. Future plan is to use a distributed lock through [etcd](https://github.com/coreos/etcd) to allow coordinating challenges when more than one instance is running.