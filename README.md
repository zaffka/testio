## Istio service mesh playground

### Prerequisites

- Kind tool installed to run local Kubernetes cluster https://kind.sigs.k8s.io/
- Istio installed
- Docker (+ compose plugin)

### Usage

`make up` - to run k8s cluster and a local docker registry
`make apply` - to run a simple app within a cluster
