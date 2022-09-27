## Istio service mesh playground

### Prerequisites

- Kind tool installed to run local Kubernetes cluster https://kind.sigs.k8s.io/
- Istio installed https://istio.io/latest/docs/setup/getting-started
- Docker (+ compose plugin) https://docs.docker.com/engine/install/ubuntu/

### Usage

`make` - to run k8s cluster, local docker registry, build, tag and push docker image, apply specs to k8s cluster

Look into Makefile to see other recipes.

### Envs

```
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT
```
