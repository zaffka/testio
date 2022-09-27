.PHONY: up kubeup istioup down build push apply remove addons

all: up build push apply

up: kubeup istioup

kubeup:
	./cluster_and_registry.sh

istioup:
	istioctl install --set profile=demo -y
	kubectl label namespace default istio-injection=enabled
	
down:
	kind delete cluster
	docker rm kind-registry -f

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.periodStr=5s" -o testio .
	docker build -t localhost:5001/testio:latest .

push:
	docker push localhost:5001/testio:latest

apply:
	kubectl apply -f testio-deployment.yaml
	kubectl apply -f testio-gateway.yaml

remove:
	kubectl delete -f testio-deployment.yaml
	kubectl delete -f testio-gateway.yaml

addons:
	kubectl apply -f istio_addons/
	kubectl rollout status deployment/kiali -n istio-system