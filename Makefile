.PHONY: up kubeup istioup down build push apply remove addons

all: up build push apply

up: kubeup dockerup istioup

dockerup:
	docker compose up -d

kubeup:
	./cluster_and_registry.sh

istioup:
	istioctl install --set profile=demo --set meshConfig.outboundTrafficPolicy.mode=ALLOW_ANY -y
	kubectl label namespace default istio-injection=enabled

down: clusterdown dockerdown

clusterdown:
	kind delete cluster
	docker rm kind-registry -f

dockerdown:
	docker compose down

build:
	$(eval GIT_BRANCH=$(shell git rev-parse --short HEAD))
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.ver=$(GIT_BRANCH)" -o testio .
	docker build -t localhost:5001/testio:latest .

push:
	docker push localhost:5001/testio:latest

apply:
	kubectl apply -f .k8s/specs/testio-deployment.yaml
	kubectl apply -f .k8s/specs/testio-gateway.yaml
	kubectl apply -f .k8s/specs/dashboard.yaml
	kubectl apply -f .k8s/specs/service-account.yaml

delapps:
	kubectl delete -f .k8s/specs/testio-deployment.yaml

remove: delapps
	kubectl delete -f .k8s/specs/testio-gateway.yaml
	kubectl delete -f .k8s/specs/dashboard.yaml
	kubectl delete -f .k8s/specs/service-account.yaml

addons:
	kubectl apply -f .k8s/istio_addons/
	kubectl rollout status deployment/kiali -n istio-system

token:
	kubectl -n kubernetes-dashboard create token admin-user

migrate-up:
	migrate -path sql/migrations -database clickhouse://localhost:9000 up

migrate-down:
	migrate -path sql/migrations -database clickhouse://localhost:9000 down

kiali:
	istioctl dashboard kiali