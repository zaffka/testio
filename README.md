## Istio service mesh playground

### Prerequisites

- Kind tool installed to run local Kubernetes cluster https://kind.sigs.k8s.io/
- Istio installed https://istio.io/latest/docs/setup/getting-started
- Docker (+ compose plugin) https://docs.docker.com/engine/install/ubuntu/

### Usage

`make up` - to run k8s cluster and a local docker registry
`make apply` - to run a simple app within a cluster

### Envs

```
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}') \
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}') \
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}') \
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT
echo "$GATEWAY_URL"
```
