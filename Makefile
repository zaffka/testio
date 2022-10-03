.PHONY: up kubeup istioup down build push apply remove addons

all: up build push apply

up: kubeup istioup

kubeup:
	./cluster_and_registry.sh
	docker compose up -d

istioup:
	istioctl install --set profile=demo -y
	kubectl label namespace default istio-injection=enabled
	
down:
	kind delete cluster
	docker rm kind-registry -f
	docker compose down

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.periodStr=5s" -o testio .
	docker build -t localhost:5001/testio:latest .

push:
	docker push localhost:5001/testio:latest

apply:
	kubectl apply -f specs/testio-deployment.yaml
	kubectl apply -f specs/testio-gateway.yaml
	kubectl apply -f specs/dashboard.yaml
	kubectl apply -f specs/service-account.yaml

remove:
	kubectl delete -f specs/testio-deployment.yaml
	kubectl delete -f specs/testio-gateway.yaml

addons:
	kubectl apply -f specs/istio_addons/
	kubectl rollout status deployment/kiali -n istio-system

token:
	kubectl -n kubernetes-dashboard create token admin-user

migrate-up:
	migrate -path sql/migrations -database clickhouse://localhost:9000 up

migrate-down:
	migrate -path sql/migrations -database clickhouse://localhost:9000 down